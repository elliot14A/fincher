package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	ch "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/config"
	"github.com/elliot14A/fincher/internal/scheduler"
	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

type RunnerResult struct {
	RunID             string          `json:"run_id"`
	ExecutedActions   []models.Action `json:"executed_actions"`
	Artifacts         []models.Action `json:"artifacts"`
	SkippedActions    []models.Action `json:"skipped_actions,omitempty"`
	DownstreamEmitted int             `json:"downstream_emitted"`
}

type SchedulerInterface interface {
	ScheduleTask(
		kind scheduler.TaskKind,
		targetID, titleSlug, vendorID string,
		component models.ComponentType,
		forceOutcome string,
		turnaroundHours float64,
		onComplete func(t *scheduler.Task),
	) (*scheduler.Task, error)
	CancelTasksForTitle(titleSlug string) int
	DecideOutcome(force string, component models.ComponentType) scheduler.QCOutcome
}

type RunnerDeps struct {
	TursoClient        *ent.Client
	ClickHouse         *sql.DB
	Scheduler          SchedulerInterface
	OnScheduleComplete func(event models.Event)
}

func RunActionPlan(
	ctx context.Context,
	tursoClient *ent.Client,
	chDB *sql.DB,
	runID string,
	stepID string,
	plan *models.ActionPlan,
) domainerrors.Result[*RunnerResult] {
	return RunActionPlanWithDeps(ctx, RunnerDeps{
		TursoClient: tursoClient,
		ClickHouse:  chDB,
	}, runID, stepID, plan)
}

func RunActionPlanWithDeps(
	ctx context.Context,
	deps RunnerDeps,
	runID string,
	stepID string,
	plan *models.ActionPlan,
) domainerrors.Result[*RunnerResult] {
	if plan == nil {
		return domainerrors.Err[*RunnerResult](fmt.Errorf("plan cannot be nil"))
	}
	if deps.TursoClient == nil {
		return domainerrors.Err[*RunnerResult](fmt.Errorf("turso client cannot be nil"))
	}

	executed := make([]models.Action, 0, len(plan.Actions))
	artifacts := make([]models.Action, 0)
	skipped := make([]models.Action, 0)
	var downstreamEvents []models.Event

	for _, action := range plan.Actions {
		switch action.Type {
		case models.ActionHoldTitle:
			titleID := action.TargetID
			if resolved := tursotitles.FindByIDOrSlug(ctx, deps.TursoClient, action.TargetID); resolved.IsOk() {
				titleID = resolved.Unwrap().ID
			}
			overdueStatus := models.StatusOverdue
			updRes := tursotitles.Update(ctx, deps.TursoClient, titleID, &models.UpdateTitleInput{
				OverallStatus: &overdueStatus,
			})
			if updRes.IsErr() {
				logger.Warn("runner: skipping HOLD_TITLE for unresolvable title; continuing plan",
					"run_id", runID, "target_id", action.TargetID, "error", updRes.Error())
				skipped = append(skipped, action)
				continue
			}
			executed = append(executed, action)

			downstreamEvents = append(downstreamEvents, models.Event{
				ID:              "evt-" + titleID + "-hold",
				Source:          "fincher/runner",
				Type:            "fincher.title.held",
				Subject:         plan.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityCritical,
				DataContentType: "application/json",
				Data: map[string]any{
					"title_id": titleID,
					"reason":   action.Reason,
					"status":   "OVERDUE",
				},
			})

		case models.ActionHoldDelivery:
			holdStatus := models.DeliveryStatusHold
			updRes := tursodeliveries.Update(ctx, deps.TursoClient, action.TargetID, &models.UpdateDeliveryInput{
				Status: &holdStatus,
			})
			if updRes.IsErr() {
				return domainerrors.Err[*RunnerResult](updRes.Error())
			}
			executed = append(executed, action)

			downstreamEvents = append(downstreamEvents, models.Event{
				ID:              "evt-" + action.TargetID + "-hold",
				Source:          "fincher/runner",
				Type:            models.TypeDeliveryHeld,
				Subject:         plan.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityWarn,
				DataContentType: "application/json",
				Data: map[string]any{
					"delivery_id": action.TargetID,
					"reason":      action.Reason,
				},
			})

		case models.ActionReleaseDelivery:
			readyStatus := models.DeliveryStatusReadyToShip
			updRes := tursodeliveries.Update(ctx, deps.TursoClient, action.TargetID, &models.UpdateDeliveryInput{
				Status: &readyStatus,
			})
			if updRes.IsErr() {
				logger.Warn("runner: skipping RELEASE_DELIVERY for unresolvable target; continuing plan",
					"run_id", runID,
					"target_id", action.TargetID,
					"error", updRes.Error(),
				)
				skipped = append(skipped, action)
				continue
			}
			executed = append(executed, action)

			downstreamEvents = append(downstreamEvents, models.Event{
				ID:              "evt-" + action.TargetID + "-rel",
				Source:          "fincher/runner",
				Type:            models.TypeDeliveryReleased,
				Subject:         plan.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityInfo,
				DataContentType: "application/json",
				Data: map[string]any{
					"delivery_id": action.TargetID,
					"reason":      action.Reason,
				},
			})

		case models.ActionReassignVendor:
			newVendorID := action.TargetID
			var targetPkgID string
			if pkgIDVal, ok := action.Payload["package_id"]; ok {
				if pkgID, isStr := pkgIDVal.(string); isStr && pkgID != "" {
					targetPkgID = pkgID
					reassignInput := &models.UpdatePackageInput{VendorID: &newVendorID}
					if pRes := tursopackages.Get(ctx, deps.TursoClient, pkgID); pRes.IsOk() {
						if tRes := tursotitles.Get(ctx, deps.TursoClient, pRes.Unwrap().TitleID); tRes.IsOk() {
							activeMaster := tRes.Unwrap().CurrentMasterVersion
							if activeMaster != "" {
								reassignInput.DerivedFromMasterVersion = &activeMaster
							}
						}
					}
					updRes := tursopackages.Update(ctx, deps.TursoClient, pkgID, reassignInput)
					if updRes.IsErr() {
						logger.Warn("runner: skipping REASSIGN_VENDOR for unresolvable package; continuing plan",
							"run_id", runID,
							"package_id", pkgID,
							"vendor_id", newVendorID,
							"error", updRes.Error(),
						)
						skipped = append(skipped, action)
						continue
					}
				}
			}
			executed = append(executed, action)

			downstreamEvents = append(downstreamEvents, models.Event{
				ID:              "evt-" + action.TargetID + "-assign",
				Source:          "fincher/runner",
				Type:            models.TypeVendorAssigned,
				Subject:         plan.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityInfo,
				DataContentType: "application/json",
				Data: map[string]any{
					"vendor_id":  action.TargetID,
					"package_id": targetPkgID,
					"reason":     action.Reason,
				},
			})

			if deps.Scheduler != nil && targetPkgID != "" {
				turnaroundHours := config.DefaultTurnaroundHours
				var pkgComponent models.ComponentType
				pRes := tursopackages.Get(ctx, deps.TursoClient, targetPkgID)
				if pRes.IsOk() {
					pkgComponent = pRes.Unwrap().Component
				}

				if vHoursVal, ok := action.Payload["turnaround_hours"]; ok {
					if vHours, isFloat := vHoursVal.(float64); isFloat && vHours > 0 {
						turnaroundHours = vHours
					} else if vHoursInt, isInt := vHoursVal.(int); isInt && vHoursInt > 0 {
						turnaroundHours = float64(vHoursInt)
					}
				}
				if pkgComponent == "" && action.Payload != nil {
					if cVal, ok := action.Payload["component"].(string); ok && cVal != "" {
						pkgComponent = models.ComponentType(cVal)
					}
				}

				forceOutcome := ""
				if action.Payload != nil {
					if fVal, ok := action.Payload["force_outcome"].(string); ok && fVal != "" {
						forceOutcome = fVal
					}
				}

				_, schedErr := deps.Scheduler.ScheduleTask(
					scheduler.TaskKindPackage,
					targetPkgID,
					plan.TitleSlug,
					newVendorID,
					pkgComponent,
					forceOutcome,
					turnaroundHours,
					BuildQCCompletionCallback(ctx, QCScheduleDeps{
						TursoClient:        deps.TursoClient,
						Scheduler:          deps.Scheduler,
						OnScheduleComplete: deps.OnScheduleComplete,
						LogContext:         runID,
					}),
				)
				if schedErr != nil {
					logger.Error("runner: failed to schedule repair task",
						"run_id", runID,
						"title_slug", plan.TitleSlug,
						"package_id", targetPkgID,
						"vendor_id", newVendorID,
						"error", schedErr,
					)
				}
			}

		case models.ActionEmailVendor:
			if action.Payload == nil {
				action.Payload = make(map[string]any)
			}
			dispatchID := "msg-email-" + uuid.NewString()[:8]
			action.Payload["dispatch_id"] = dispatchID
			action.Payload["status"] = "DELIVERED"
			action.Payload["dispatched_at"] = time.Now().UTC().Format(time.RFC3339)
			artifacts = append(artifacts, action)
			executed = append(executed, action)

			downstreamEvents = append(downstreamEvents, models.Event{
				ID:              "evt-" + dispatchID,
				Source:          "fincher/runner",
				Type:            models.TypeVendorEmailed,
				Subject:         plan.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityInfo,
				DataContentType: "application/json",
				Data:            action.Payload,
			})

		case models.ActionNotifyStakeholders:
			if action.Payload == nil {
				action.Payload = make(map[string]any)
			}
			dispatchID := "msg-slack-" + uuid.NewString()[:8]
			action.Payload["dispatch_id"] = dispatchID
			action.Payload["channel"] = "#ops-war-room"
			action.Payload["status"] = "DELIVERED"
			action.Payload["dispatched_at"] = time.Now().UTC().Format(time.RFC3339)
			artifacts = append(artifacts, action)
			executed = append(executed, action)

			downstreamEvents = append(downstreamEvents, models.Event{
				ID:              "evt-" + dispatchID,
				Source:          "fincher/runner",
				Type:            models.TypeStakeholdersNotified,
				Subject:         plan.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityInfo,
				DataContentType: "application/json",
				Data:            action.Payload,
			})

		case models.ActionPostSocialUpdate:
			if action.Payload == nil {
				action.Payload = make(map[string]any)
			}
			postID := "post-x-" + uuid.NewString()[:8]
			action.Payload["post_id"] = postID
			action.Payload["platform"] = "x/twitter"
			action.Payload["status"] = "PUBLISHED"
			action.Payload["dispatched_at"] = time.Now().UTC().Format(time.RFC3339)
			artifacts = append(artifacts, action)
			executed = append(executed, action)

			downstreamEvents = append(downstreamEvents, models.Event{
				ID:              "evt-" + postID,
				Source:          "fincher/runner",
				Type:            models.TypeSocialPosted,
				Subject:         plan.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityInfo,
				DataContentType: "application/json",
				Data:            action.Payload,
			})
		}
	}

	now := time.Now().UTC()
	if stepID != "" {
		artifactsJSON, _ := json.Marshal(artifacts)
		runs.UpdateStepStatus(ctx, deps.TursoClient, stepID, models.StepStatusCompleted, &now, map[string]any{
			"artifacts_count": len(artifacts),
			"artifacts_json":  string(artifactsJSON),
		})
	}

	if deps.ClickHouse != nil && len(downstreamEvents) > 0 {
		batchRes := ch.InsertBatch(ctx, deps.ClickHouse, downstreamEvents)
		if batchRes.IsErr() {
			return domainerrors.Err[*RunnerResult](batchRes.Error())
		}
	}

	if runID != "" {
		runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, nil)
	}

	return domainerrors.Ok(&RunnerResult{
		RunID:             runID,
		ExecutedActions:   executed,
		Artifacts:         artifacts,
		SkippedActions:    skipped,
		DownstreamEmitted: len(downstreamEvents),
	})
}
