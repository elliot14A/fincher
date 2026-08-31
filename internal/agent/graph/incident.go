package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/agent/scheduler"
	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// DefaultMaxRemediationAttempts is the maximum number of revisions allowed before human escalation.
const DefaultMaxRemediationAttempts = 3

// ExecuteIncident runs the multi-stage incident investigation pipeline:
//   - Stage 1: Triage Judge (screens routine vs actionable events)
//   - Stage 2: Context Gathering (queries blast radius, ClickHouse analytics, vendor candidates)
//   - Stage 3: Remediation & Policy Verification Loop (proposes and gates actions up to MaxAttempts)
//   - Stage 4: Execution / Escalation (applies state mutations or marks for human review)
//
// Each stage persists Step and WfResult rows into Turso so the live SSE console streams real progress.
func ExecuteIncident(ctx context.Context, deps IncidentGraphDeps, input IncidentInput) (*IncidentOutput, error) {
	if input.Event == nil {
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInvalidInput, "input event cannot be nil", nil)
	}

	// 1. Resolve Run ID and ensure root Run record exists
	runID := input.RunID
	if runID == "" {
		if input.Event.ID != "" {
			runID = fmt.Sprintf("run-%s", input.Event.ID)
		} else {
			runID = "run-" + uuid.NewString()[:8]
		}
	}

	if deps.TursoClient != nil {
		existingRun := runs.GetRun(ctx, deps.TursoClient, runID)
		if existingRun.IsErr() {
			rRes := runs.CreateRun(ctx, deps.TursoClient, &models.Run{
				Base:      models.Base{ID: runID},
				TitleSlug: input.Event.Subject,
				Trigger:   input.Event.Type,
				Status:    models.RunStatusRunning,
				StartedAt: time.Now().UTC(),
			})
			if rRes.IsErr() {
				logger.Error("failed to persist initial run in turso", "run_id", runID, "error", rRes.Error())
			}
		}
	}

	// 2. Stage 1: Triage Judge
	triageStepID := fmt.Sprintf("step-%s-triage", runID)
	if deps.TursoClient != nil {
		sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
			Base:      models.Base{ID: triageStepID},
			RunID:     runID,
			Name:      "triage_judge",
			Status:    models.StepStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if sRes.IsErr() {
			logger.Error("failed to persist triage step", "step_id", triageStepID, "error", sRes.Error())
		}
	}

	hoursUntilPremiere := tursotitles.ResolveHoursUntilPremiere(ctx, deps.TursoClient, input.Event.Subject, input.HoursUntilPremiere)

	filterRes := agent.FilterEvent(ctx, deps.Model, input.Event, hoursUntilPremiere)
	now := time.Now().UTC()
	if filterRes.IsErr() {
		if deps.TursoClient != nil {
			runs.UpdateStepStatus(ctx, deps.TursoClient, triageStepID, models.StepStatusFailed, &now, map[string]any{"error": filterRes.Error().Error()})
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		}
		return nil, filterRes.Error()
	}
	filterDecision := filterRes.Unwrap()

	outcome := "ACTIONABLE"
	if !filterDecision.Actionable {
		outcome = "FILTERED"
	}
	if deps.TursoClient != nil {
		runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
			Base:      models.Base{ID: fmt.Sprintf("res-%s-triage", runID)},
			RunID:     runID,
			StepID:    triageStepID,
			Judge:     "triage_judge",
			Outcome:   outcome,
			Rationale: filterDecision.Rationale,
			Attempt:   1,
		})
		runs.UpdateStepStatus(ctx, deps.TursoClient, triageStepID, models.StepStatusCompleted, &now, map[string]any{
			"actionable": filterDecision.Actionable,
			"severity":   string(filterDecision.Severity),
			"rationale":  filterDecision.Rationale,
		})
	}

	if !filterDecision.Actionable {
		if deps.TursoClient != nil {
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, map[string]any{
				"decision":  "FILTERED",
				"rationale": filterDecision.Rationale,
			})
		}
		return &IncidentOutput{
			Actionable: false,
			Decision:   "FILTERED",
			Rationale:  filterDecision.Rationale,
			Attempts:   0,
		}, nil
	}

	// 3. Stage 2: Context Gathering
	contextStepID := fmt.Sprintf("step-%s-context", runID)
	if deps.TursoClient != nil {
		sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
			Base:      models.Base{ID: contextStepID},
			RunID:     runID,
			Name:      "context_gathering",
			Status:    models.StepStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if sRes.IsErr() {
			logger.Error("failed to persist context step", "step_id", contextStepID, "error", sRes.Error())
		}
	}

	packageID, _ := input.Event.Data["package_id"].(string)
	vendorID, _ := input.Event.Data["vendor_id"].(string)
	component, _ := input.Event.Data["component"].(string)

	impact, err := tools.FetchDeliveryImpact(ctx, deps.TursoClient, tools.DeliveryImpactArgs{
		PackageID:          packageID,
		HoursUntilPremiere: hoursUntilPremiere,
	})
	if err != nil {
		now = time.Now().UTC()
		if deps.TursoClient != nil {
			runs.UpdateStepStatus(ctx, deps.TursoClient, contextStepID, models.StepStatusFailed, &now, map[string]any{"error": err.Error()})
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		}
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInternal, "failed to gather delivery impact", err)
	}

	analytics, err := tools.FetchAnalytics(ctx, deps.ClickHouse, tools.AnalyticsArgs{
		VendorID:  vendorID,
		TitleSlug: input.Event.Subject,
		Component: component,
	})
	if err != nil {
		now = time.Now().UTC()
		if deps.TursoClient != nil {
			runs.UpdateStepStatus(ctx, deps.TursoClient, contextStepID, models.StepStatusFailed, &now, map[string]any{"error": err.Error()})
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		}
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInternal, "failed to gather historical analytics", err)
	}

	candidates, err := tools.FetchVendorCandidates(ctx, deps.TursoClient, deps.ClickHouse, tools.VendorCandidatesArgs{
		Component: component,
		Specialty: component,
	})
	if err != nil {
		now = time.Now().UTC()
		if deps.TursoClient != nil {
			runs.UpdateStepStatus(ctx, deps.TursoClient, contextStepID, models.StepStatusFailed, &now, map[string]any{"error": err.Error()})
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		}
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInternal, "failed to gather vendor candidates", err)
	}

	now = time.Now().UTC()
	if deps.TursoClient != nil {
		runs.UpdateStepStatus(ctx, deps.TursoClient, contextStepID, models.StepStatusCompleted, &now, map[string]any{
			"deliveries_on_hold": len(impact.AffectedDeliveries),
			"candidates_count":   len(candidates),
		})
	}

	// 4. Stage 3: Remediation Planning & Policy Verification Loop
	maxAttempts := deps.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = DefaultMaxRemediationAttempts
	}

	remediationStepID := fmt.Sprintf("step-%s-remediation", runID)
	if deps.TursoClient != nil {
		sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
			Base:      models.Base{ID: remediationStepID},
			RunID:     runID,
			Name:      "remediation_loop",
			Status:    models.StepStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if sRes.IsErr() {
			logger.Error("failed to persist remediation step", "step_id", remediationStepID, "error", sRes.Error())
		}
	}

	var finalPlan *models.ActionPlan
	var lastVerification *agent.VerificationResult
	feedback := ""

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		planRes := agent.PlanRemediation(ctx, deps.Model, input.Event, impact, analytics, candidates, feedback)
		if planRes.IsErr() {
			now = time.Now().UTC()
			if deps.TursoClient != nil {
				runs.UpdateStepStatus(ctx, deps.TursoClient, remediationStepID, models.StepStatusFailed, &now, map[string]any{"error": planRes.Error().Error()})
				runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
			}
			return nil, planRes.Error()
		}
		currentPlan := planRes.Unwrap()

		verifyRes := agent.VerifyPlan(currentPlan, impact, candidates, attempt)
		if verifyRes.IsErr() {
			now = time.Now().UTC()
			if deps.TursoClient != nil {
				runs.UpdateStepStatus(ctx, deps.TursoClient, remediationStepID, models.StepStatusFailed, &now, map[string]any{"error": verifyRes.Error().Error()})
				runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
			}
			return nil, verifyRes.Error()
		}
		lastVerification = verifyRes.Unwrap()

		if deps.TursoClient != nil {
			runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
				Base:      models.Base{ID: fmt.Sprintf("res-%s-verify-%d", runID, attempt)},
				RunID:     runID,
				StepID:    remediationStepID,
				Judge:     "policy_verifier",
				Outcome:   string(lastVerification.Decision),
				Rationale: lastVerification.Rationale,
				Attempt:   attempt,
			})
		}

		if lastVerification.Decision == agent.DecisionApproved {
			finalPlan = currentPlan
			break
		}

		if lastVerification.Decision == agent.DecisionEscalate {
			break
		}

		feedback = lastVerification.Rationale
	}

	now = time.Now().UTC()
	if deps.TursoClient != nil {
		runs.UpdateStepStatus(ctx, deps.TursoClient, remediationStepID, models.StepStatusCompleted, &now, map[string]any{
			"attempts": lastVerification.Attempt,
			"decision": string(lastVerification.Decision),
		})
	}

	// 5. Stage 4: Execution / Escalation
	if lastVerification.Decision == agent.DecisionApproved && finalPlan != nil {
		executorStepID := fmt.Sprintf("step-%s-executor", runID)
		if deps.TursoClient != nil {
			sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
				Base:      models.Base{ID: executorStepID},
				RunID:     runID,
				Name:      "remediation_executor",
				Status:    models.StepStatusRunning,
				StartedAt: time.Now().UTC(),
			})
			if sRes.IsErr() {
				logger.Error("failed to persist executor step", "step_id", executorStepID, "error", sRes.Error())
			}
		}

		var scheduleTask func(kind, targetID, titleSlug, vendorID string, turnaroundHours float64, onComplete func()) error
		if deps.Scheduler != nil {
			scheduleTask = func(kind, targetID, titleSlug, vendorID string, turnaroundHours float64, onComplete func()) error {
				_, err := deps.Scheduler.ScheduleTask(
					scheduler.TaskKind(kind),
					targetID,
					titleSlug,
					vendorID,
					turnaroundHours,
					func(t *scheduler.Task) {
						if onComplete != nil {
							onComplete()
						}
					},
				)
				return err
			}
		}

		execRes := agent.RunActionPlanWithDeps(ctx, agent.RunnerDeps{
			TursoClient:        deps.TursoClient,
			ClickHouse:         deps.ClickHouse,
			ScheduleTask:       scheduleTask,
			OnScheduleComplete: deps.OnScheduleComplete,
		}, runID, executorStepID, finalPlan)
		if execRes.IsErr() {
			now = time.Now().UTC()
			if deps.TursoClient != nil {
				runs.UpdateStepStatus(ctx, deps.TursoClient, executorStepID, models.StepStatusFailed, &now, map[string]any{"error": execRes.Error().Error()})
				runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
			}
			return nil, execRes.Error()
		}
		result := execRes.Unwrap()

		now = time.Now().UTC()
		if deps.TursoClient != nil {
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, map[string]any{
				"decision": string(agent.DecisionApproved),
				"attempts": lastVerification.Attempt,
			})
		}

		return &IncidentOutput{
			Actionable:   true,
			Decision:     agent.DecisionApproved,
			Rationale:    lastVerification.Rationale,
			ActionPlan:   finalPlan,
			RunnerResult: result,
			Attempts:     lastVerification.Attempt,
		}, nil
	}

	now = time.Now().UTC()
	if deps.TursoClient != nil {
		runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusEscalated, &now, map[string]any{
			"decision":  string(lastVerification.Decision),
			"rationale": lastVerification.Rationale,
			"attempts":  lastVerification.Attempt,
		})
	}

	return &IncidentOutput{
		Actionable: true,
		Decision:   lastVerification.Decision,
		Rationale:  lastVerification.Rationale,
		ActionPlan: finalPlan,
		Attempts:   lastVerification.Attempt,
	}, nil
}
