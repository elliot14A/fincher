package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// DefaultMaxRemediationAttempts caps the verifier self-correction loop.
const DefaultMaxRemediationAttempts = 3

func failIncidentStage(ctx context.Context, client *ent.Client, runID, titleSlug, stepID, stageName string, err error) {
	now := time.Now().UTC()
	stepUpd := runs.UpdateStepStatus(ctx, client, stepID, models.StepStatusFailed, &now, map[string]any{"error": err.Error()})
	if stepUpd.IsErr() {
		logger.Warn("incident: failed to update step status on stage failure",
			"run_id", runID,
			"title_slug", titleSlug,
			"step_id", stepID,
			"stage", stageName,
			"error", stepUpd.Error(),
		)
	}
	runUpd := runs.UpdateRunStatus(ctx, client, runID, models.RunStatusFailed, &now, nil)
	if runUpd.IsErr() {
		logger.Warn("incident: failed to update run status on stage failure",
			"run_id", runID,
			"title_slug", titleSlug,
			"stage", stageName,
			"error", runUpd.Error(),
		)
	}
}

// ExecuteIncident runs the 4-stage Multi-Agent incident investigation and remediation graph:
//  1. Stage 1: Triage Judge (Filter benign/routine events vs actionable anomalies)
//  2. Stage 2: Context Gathering (Delivery impact, ClickHouse vendor analytics, vendor candidates, title projection)
//  3. Stage 3: Remediation Planning & Policy Verification Loop (bounded self-correction up to maxAttempts)
//  4. Stage 4: Execution or Escalation (software executor mutates SQLite state & emits downstream events, or escalates to operator)
func ExecuteIncident(ctx context.Context, deps IncidentGraphDeps, input IncidentInput) (*IncidentOutput, error) {
	if input.Event == nil {
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInvalidInput, "input event cannot be nil", nil)
	}
	if deps.TursoClient == nil {
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	if deps.Model == nil {
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil)
	}

	titleSlug := input.Event.Subject
	if titleSlug == "" {
		titleSlug = models.DefaultTitleAgnosticSentinel
	}

	// 1. Resolve Run ID and guarantee root Run record exists in Turso
	runID := input.RunID
	if runID == "" {
		if input.Event.ID != "" {
			runID = "run-" + input.Event.ID
		} else {
			runID = "run-" + uuid.NewString()[:8]
		}
	}

	existingRun := runs.GetRun(ctx, deps.TursoClient, runID)
	if existingRun.IsErr() {
		rRes := runs.CreateRun(ctx, deps.TursoClient, &models.Run{
			Base:      models.Base{ID: runID},
			TitleSlug: titleSlug,
			Trigger:   "incident",
			Status:    models.RunStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if rRes.IsErr() {
			logger.Error("incident: failed to persist initial run in turso",
				"run_id", runID,
				"title_slug", titleSlug,
				"event_id", input.Event.ID,
				"error", rRes.Error(),
			)
		}
	}

	// 2. Stage 1: Triage Judge
	triageStepID := fmt.Sprintf("step-%s-triage", runID)
	sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
		Base:      models.Base{ID: triageStepID},
		RunID:     runID,
		Name:      "triage_judge",
		Status:    models.StepStatusRunning,
		StartedAt: time.Now().UTC(),
	})
	if sRes.IsErr() {
		logger.Error("incident: failed to persist triage step",
			"run_id", runID,
			"title_slug", titleSlug,
			"step_id", triageStepID,
			"error", sRes.Error(),
		)
	}

	hoursUntilPremiere := tursotitles.ResolveHoursUntilPremiere(ctx, deps.TursoClient, input.Event.Subject, input.HoursUntilPremiere)

	filterRes := agent.FilterEvent(ctx, deps.Model, input.Event, hoursUntilPremiere)
	now := time.Now().UTC()
	if filterRes.IsErr() {
		failIncidentStage(ctx, deps.TursoClient, runID, titleSlug, triageStepID, "triage_judge", filterRes.Error())
		return nil, filterRes.Error()
	}
	filterDecision := filterRes.Unwrap()

	outcome := "ACTIONABLE"
	if !filterDecision.Actionable {
		outcome = "FILTERED"
	}
	resRes := runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
		Base:      models.Base{ID: fmt.Sprintf("res-%s-triage", runID)},
		RunID:     runID,
		StepID:    triageStepID,
		Judge:     "triage_judge",
		Outcome:   outcome,
		Rationale: filterDecision.Rationale,
		Attempt:   1,
	})
	if resRes.IsErr() {
		logger.Warn("incident: failed to record triage wf_result", "run_id", runID, "step_id", triageStepID, "error", resRes.Error())
	}

	runs.UpdateStepStatus(ctx, deps.TursoClient, triageStepID, models.StepStatusCompleted, &now, map[string]any{
		"actionable": filterDecision.Actionable,
		"severity":   string(filterDecision.Severity),
		"rationale":  filterDecision.Rationale,
	})

	if !filterDecision.Actionable {
		runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, map[string]any{
			"decision":  "FILTERED",
			"rationale": filterDecision.Rationale,
		})
		return &IncidentOutput{
			Actionable: false,
			Decision:   "FILTERED",
			Rationale:  filterDecision.Rationale,
			Attempts:   0,
		}, nil
	}

	// 3. Stage 2: Context Gathering
	contextStepID := fmt.Sprintf("step-%s-context", runID)
	cRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
		Base:      models.Base{ID: contextStepID},
		RunID:     runID,
		Name:      "context_gathering",
		Status:    models.StepStatusRunning,
		StartedAt: time.Now().UTC(),
	})
	if cRes.IsErr() {
		logger.Error("incident: failed to persist context step",
			"run_id", runID,
			"title_slug", titleSlug,
			"step_id", contextStepID,
			"error", cRes.Error(),
		)
	}

	packageID, _ := input.Event.Data["package_id"].(string)
	vendorID, _ := input.Event.Data["vendor_id"].(string)
	component, _ := input.Event.Data["component"].(string)

	impact, err := tools.FetchDeliveryImpact(ctx, deps.TursoClient, tools.DeliveryImpactArgs{
		PackageID:          packageID,
		HoursUntilPremiere: hoursUntilPremiere,
	})
	if err != nil {
		failIncidentStage(ctx, deps.TursoClient, runID, titleSlug, contextStepID, "context_delivery_impact", err)
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInternal, "failed to gather delivery impact", err)
	}

	var analytics *models.AnalyticsSummary
	if deps.ClickHouse != nil {
		analytics, err = tools.FetchAnalytics(ctx, deps.ClickHouse, tools.AnalyticsArgs{
			VendorID:  vendorID,
			TitleSlug: input.Event.Subject,
			Component: component,
		})
		if err != nil {
			failIncidentStage(ctx, deps.TursoClient, runID, titleSlug, contextStepID, "context_historical_analytics", err)
			return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInternal, "failed to gather historical analytics", err)
		}
	} else {
		analytics = &models.AnalyticsSummary{
			VendorHistoricalAccuracy: models.UnmeasuredHistoricalAccuracy,
		}
	}

	candidates, err := tools.FetchVendorCandidates(ctx, deps.TursoClient, deps.ClickHouse, tools.VendorCandidatesArgs{
		Component: component,
		Specialty: component,
	})
	if err != nil {
		failIncidentStage(ctx, deps.TursoClient, runID, titleSlug, contextStepID, "context_vendor_candidates", err)
		return nil, domainerrors.NewWithOp("graph.ExecuteIncident", domainerrors.CodeInternal, "failed to gather vendor candidates", err)
	}

	projection, err := tools.GetTitleReadyProjection(ctx, deps.TursoClient, input.Event.Subject)
	if err != nil {
		logger.Warn("incident: could not fetch title readiness projection",
			"run_id", runID,
			"title_slug", input.Event.Subject,
			"error", err,
		)
	}

	now = time.Now().UTC()
	contextMeta := map[string]any{
		"deliveries_on_hold": len(impact.AffectedDeliveries),
		"candidates_count":   len(candidates),
	}
	if projection != nil {
		contextMeta["projection"] = projection
	}
	runs.UpdateStepStatus(ctx, deps.TursoClient, contextStepID, models.StepStatusCompleted, &now, contextMeta)

	// 4. Stage 3: Remediation Planning & Policy Verification Loop
	maxAttempts := deps.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = DefaultMaxRemediationAttempts
	}

	remediationStepID := fmt.Sprintf("step-%s-remediation", runID)
	rStepRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
		Base:      models.Base{ID: remediationStepID},
		RunID:     runID,
		Name:      "remediation_loop",
		Status:    models.StepStatusRunning,
		StartedAt: time.Now().UTC(),
	})
	if rStepRes.IsErr() {
		logger.Error("incident: failed to persist remediation step",
			"run_id", runID,
			"title_slug", titleSlug,
			"step_id", remediationStepID,
			"error", rStepRes.Error(),
		)
	}

	var finalPlan *models.ActionPlan
	var lastVerification *agent.VerificationResult
	feedback := ""

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		planRes := agent.PlanRemediation(ctx, deps.Model, input.Event, impact, analytics, candidates, projection, feedback)
		if planRes.IsErr() {
			failIncidentStage(ctx, deps.TursoClient, runID, titleSlug, remediationStepID, "remediation_plan", planRes.Error())
			return nil, planRes.Error()
		}
		currentPlan := planRes.Unwrap()

		verifyRes := agent.VerifyPlan(currentPlan, impact, candidates, attempt)
		if verifyRes.IsErr() {
			failIncidentStage(ctx, deps.TursoClient, runID, titleSlug, remediationStepID, "remediation_verify", verifyRes.Error())
			return nil, verifyRes.Error()
		}
		lastVerification = verifyRes.Unwrap()

		verifyResultRes := runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
			Base:      models.Base{ID: fmt.Sprintf("res-%s-verify-%d", runID, attempt)},
			RunID:     runID,
			StepID:    remediationStepID,
			Judge:     "policy_verifier",
			Outcome:   string(lastVerification.Decision),
			Rationale: lastVerification.Rationale,
			Attempt:   attempt,
		})
		if verifyResultRes.IsErr() {
			logger.Warn("incident: failed to record policy_verifier result", "run_id", runID, "attempt", attempt, "error", verifyResultRes.Error())
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
	runs.UpdateStepStatus(ctx, deps.TursoClient, remediationStepID, models.StepStatusCompleted, &now, map[string]any{
		"attempts": lastVerification.Attempt,
		"decision": string(lastVerification.Decision),
	})

	// 5. Stage 4: Execution / Escalation
	if lastVerification.Decision == agent.DecisionApproved && finalPlan != nil {
		executorStepID := fmt.Sprintf("step-%s-executor", runID)
		exStepRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
			Base:      models.Base{ID: executorStepID},
			RunID:     runID,
			Name:      "remediation_executor",
			Status:    models.StepStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if exStepRes.IsErr() {
			logger.Error("incident: failed to persist executor step",
				"run_id", runID,
				"title_slug", titleSlug,
				"step_id", executorStepID,
				"error", exStepRes.Error(),
			)
		}

		var sched agent.SchedulerInterface
		if deps.Scheduler != nil {
			sched = deps.Scheduler
		}

		execRes := agent.RunActionPlanWithDeps(ctx, agent.RunnerDeps{
			TursoClient:        deps.TursoClient,
			ClickHouse:         deps.ClickHouse,
			Scheduler:          sched,
			OnScheduleComplete: deps.OnScheduleComplete,
		}, runID, executorStepID, finalPlan)
		if execRes.IsErr() {
			failIncidentStage(ctx, deps.TursoClient, runID, titleSlug, executorStepID, "remediation_executor", execRes.Error())
			return nil, execRes.Error()
		}
		result := execRes.Unwrap()

		now = time.Now().UTC()
		runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, map[string]any{
			"decision": string(agent.DecisionApproved),
			"attempts": lastVerification.Attempt,
		})

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
	logger.Warn("incident: workflow escalated to human operator",
		"run_id", runID,
		"title_slug", titleSlug,
		"decision", string(lastVerification.Decision),
		"attempts", lastVerification.Attempt,
		"rationale", lastVerification.Rationale,
	)
	updRunRes := runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusEscalated, &now, map[string]any{
		"decision":  string(lastVerification.Decision),
		"rationale": lastVerification.Rationale,
		"attempts":  lastVerification.Attempt,
	})
	if updRunRes.IsErr() {
		logger.Error("incident: failed to persist escalated run status",
			"run_id", runID,
			"title_slug", titleSlug,
			"error", updRunRes.Error(),
		)
	}

	return &IncidentOutput{
		Actionable: true,
		Decision:   lastVerification.Decision,
		Rationale:  lastVerification.Rationale,
		ActionPlan: finalPlan,
		Attempts:   lastVerification.Attempt,
	}, nil
}
