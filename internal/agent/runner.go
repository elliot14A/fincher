package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/agent/scheduler"
	ch "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/config"
	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/runs"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
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
		kind scheduler.TaskKind,
		targetID, titleSlug, vendorID string,
		component models.ComponentType,
		forceOutcome string,
		turnaroundHours float64,
		onComplete func(t *scheduler.Task),
	) (*scheduler.Task, error)
	DecideOutcome(force string, component models.ComponentType) scheduler.QCOutcome
}

// RunnerDeps supplies dependencies to ActionPlan execution.
type RunnerDeps struct {
	TursoClient        *ent.Client
	ClickHouse         *sql.DB
	Scheduler          SchedulerInterface
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
		return domainerrors.Err[*RunnerResult](fmt.Errorf("plan cannot be nil"))
	}
	if deps.TursoClient == nil {
		return domainerrors.Err[*RunnerResult](fmt.Errorf("turso client cannot be nil"))
	}

	executed := make([]models.Action, 0, len(plan.Actions))
	artifacts := make([]models.Action, 0)
	var downstreamEvents []models.Event

	for _, action := range plan.Actions {
		switch action.Type {
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
				return domainerrors.Err[*RunnerResult](updRes.Error())
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
					updRes := tursopackages.Update(ctx, deps.TursoClient, pkgID, &models.UpdatePackageInput{
						VendorID: &newVendorID,
					})
					if updRes.IsErr() {
						return domainerrors.Err[*RunnerResult](updRes.Error())
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

			// Schedule background repair task if scheduler is attached
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
					func(t *scheduler.Task) {
						var outcome scheduler.QCOutcome
						if deps.Scheduler != nil {
							outcome = deps.Scheduler.DecideOutcome(t.ForceOutcome, t.Component)
						} else {
							outcome = scheduler.QCOutcomePass
						}

						if outcome == scheduler.QCOutcomePass {
							// PASS path: mark package valid and emit completed QC
							validStatus := models.PackageStatusValid
							updPkgRes := tursopackages.Update(ctx, deps.TursoClient, t.TargetID, &models.UpdatePackageInput{
								Status: &validStatus,
							})
							if updPkgRes.IsErr() {
								logger.Error("runner: failed to update package status to VALID on QC pass",
									"run_id", runID,
									"package_id", t.TargetID,
									"task_id", t.ID,
									"error", updPkgRes.Error(),
								)
							}

							qcEvent := models.Event{
								ID:              fmt.Sprintf("evt-qc-%s", t.ID),
								Source:          "fincher/qc.agent",
								Type:            models.TypeQCInspectionCompleted,
								Subject:         t.TitleSlug,
								Time:            time.Now().UTC(),
								Severity:        models.SeverityInfo,
								DataContentType: "application/json",
								Data: map[string]any{
									"package_id": t.TargetID,
									"status":     "PASSED",
									"vendor_id":  t.VendorID,
								},
							}
							if deps.OnScheduleComplete != nil {
								deps.OnScheduleComplete(qcEvent)
							}
							return
						}

						// FAIL path: inspect current redelivery attempts
						currentRedelivery := 0
						pRes := tursopackages.Get(ctx, deps.TursoClient, t.TargetID)
						if pRes.IsOk() {
							currentRedelivery = pRes.Unwrap().RedeliveryCount
						}

						if currentRedelivery >= config.MaxRedeliveryAttempts {
							// FAIL, cap exceeded -> ESCALATE to SLA breach
							slaEvent := models.Event{
								ID:              fmt.Sprintf("evt-sla-breach-%s", t.ID),
								Source:          "fincher/qc.agent",
								Type:            models.TypeVendorSLABreach,
								Subject:         t.TitleSlug,
								Time:            time.Now().UTC(),
								Severity:        models.SeverityCritical,
								DataContentType: "application/json",
								Data: map[string]any{
									"package_id":       t.TargetID,
									"vendor_id":        t.VendorID,
									"reason":           "redelivery_cap_exceeded",
									"redelivery_count": currentRedelivery,
								},
							}
							if deps.OnScheduleComplete != nil {
								deps.OnScheduleComplete(slaEvent)
							}
							return
						}

						// FAIL, under cap -> increment RedeliveryCount and emit domain-scoped defect
						newRedelivery := currentRedelivery + 1
						updRedelivRes := tursopackages.Update(ctx, deps.TursoClient, t.TargetID, &models.UpdatePackageInput{
							RedeliveryCount: &newRedelivery,
						})
						if updRedelivRes.IsErr() {
							logger.Error("runner: failed to update package redelivery count",
								"run_id", runID,
								"package_id", t.TargetID,
								"task_id", t.ID,
								"new_count", newRedelivery,
								"error", updRedelivRes.Error(),
							)
						}

						defectEventType, defectSeverity := scheduler.DefectEventTypeFor(t.Component)
						defectData := map[string]any{
							"package_id":       t.TargetID,
							"vendor_id":        t.VendorID,
							"defect_type":      "REPAIR_INSPECTION_FAILED",
							"redelivery_count": newRedelivery,
						}
						if t.Component == models.ComponentAudio {
							defectData["drift_ms"] = 110.0
							defectData["defect_type"] = "AUDIO_SYNC_DRIFT"
						}

						defectEvent := models.Event{
							ID:              fmt.Sprintf("evt-defect-%s-%d", t.ID, newRedelivery),
							Source:          "fincher/qc.agent",
							Type:            defectEventType,
							Subject:         t.TitleSlug,
							Time:            time.Now().UTC(),
							Severity:        defectSeverity,
							DataContentType: "application/json",
							Data:            defectData,
						}
						if deps.OnScheduleComplete != nil {
							deps.OnScheduleComplete(defectEvent)
						}
					},
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
		DownstreamEmitted: len(downstreamEvents),
	})
}
