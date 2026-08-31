package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	ch "github.com/elliot14A/fincher/internal/clickhouse/events"
	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/runs"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// RunnerResult records the execution outcome of an ActionPlan.
type RunnerResult struct {
	RunID             string          `json:"run_id"`
	ExecutedActions   []models.Action `json:"executed_actions"`
	Artifacts         []models.Action `json:"artifacts"`
	DownstreamEmitted int             `json:"downstream_emitted"`
}

// SchedulerInterface defines the scheduler contract for background simulation tasks.
type SchedulerInterface interface {
	ScheduleTask(
		kind string,
		targetID, titleSlug, vendorID string,
		turnaroundHours float64,
		onComplete func(t any),
	) (any, error)
}

// RunnerDeps supplies dependencies to ActionPlan execution.
type RunnerDeps struct {
	TursoClient        *ent.Client
	ClickHouse         *sql.DB
	ScheduleTask       func(kind, targetID, titleSlug, vendorID string, turnaroundHours float64, onComplete func()) error
	OnScheduleComplete func(event models.Event)
}

// RunActionPlan executes an approved ActionPlan using turso domain actions and emits downstream events to ClickHouse.
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

// RunActionPlanWithDeps executes an approved ActionPlan using provided dependencies and schedules background repairs if configured.
func RunActionPlanWithDeps(
	ctx context.Context,
	deps RunnerDeps,
	runID string,
	stepID string,
	plan *models.ActionPlan,
) domainerrors.Result[*RunnerResult] {
	if plan == nil {
		return domainerrors.Err[*RunnerResult](NewError("agent.RunActionPlan", domainerrors.CodeInvalidInput, "plan cannot be nil", nil))
	}

	executed := make([]models.Action, 0, len(plan.Actions))
	artifacts := make([]models.Action, 0)
	var downstreamEvents []models.Event

	for _, action := range plan.Actions {
		switch action.Type {
		case models.ActionHoldDelivery:
			if deps.TursoClient != nil {
				holdStatus := models.DeliveryStatusHold
				updRes := tursodeliveries.Update(ctx, deps.TursoClient, action.TargetID, &models.UpdateDeliveryInput{
					Status: &holdStatus,
				})
				if updRes.IsErr() {
					return domainerrors.Err[*RunnerResult](updRes.Error())
				}
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
			if deps.TursoClient != nil {
				readyStatus := models.DeliveryStatusReadyToShip
				updRes := tursodeliveries.Update(ctx, deps.TursoClient, action.TargetID, &models.UpdateDeliveryInput{
					Status: &readyStatus,
				})
				if updRes.IsErr() {
					return domainerrors.Err[*RunnerResult](updRes.Error())
				}
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
					if deps.TursoClient != nil {
						updRes := tursopackages.Update(ctx, deps.TursoClient, pkgID, &models.UpdatePackageInput{
							VendorID: &newVendorID,
						})
						if updRes.IsErr() {
							return domainerrors.Err[*RunnerResult](updRes.Error())
						}
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
					"vendor_id":   action.TargetID,
					"package_id":  targetPkgID,
					"reason":      action.Reason,
				},
			})

			// Schedule background repair task if scheduler hook is attached
			if deps.ScheduleTask != nil && targetPkgID != "" {
				turnaroundHours := 12.0
				if deps.TursoClient != nil {
					vRes := tursovendors.Get(ctx, deps.TursoClient, newVendorID)
					if vRes.IsOk() {
						turnaroundHours = float64(vRes.Unwrap().TurnaroundHours)
					}
				}
				_ = deps.ScheduleTask("package", targetPkgID, plan.TitleSlug, newVendorID, turnaroundHours, func() {
					qcEvent := models.Event{
						ID:              "evt-qc-" + uuid.NewString()[:8],
						Source:          "fincher/qc.agent",
						Type:            models.TypeQCInspectionCompleted,
						Subject:         plan.TitleSlug,
						Time:            time.Now().UTC(),
						Severity:        models.SeverityInfo,
						DataContentType: "application/json",
						Data: map[string]any{
							"package_id":       targetPkgID,
							"vendor_id":        newVendorID,
							"status":           "PASSED",
							"turnaround_hours": turnaroundHours,
						},
					}
					if deps.OnScheduleComplete != nil {
						deps.OnScheduleComplete(qcEvent)
					}
				})
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
	if stepID != "" && deps.TursoClient != nil {
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

	if runID != "" && deps.TursoClient != nil {
		runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, nil)
	}

	return domainerrors.Ok(&RunnerResult{
		RunID:             runID,
		ExecutedActions:   executed,
		Artifacts:         artifacts,
		DownstreamEmitted: len(downstreamEvents),
	})
}
