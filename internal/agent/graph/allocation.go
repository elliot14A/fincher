package graph

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/runs"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// ExecuteAllocation evaluates candidate vendors and creates a holistic staffing plan across all requirements
// based on turnaround feasibility, quality floors, and commercial rates.
//
// It persists candidate_gathering and vendor_selection Step and WfResult rows into Turso.
func ExecuteAllocation(ctx context.Context, deps AllocationGraphDeps, input AllocationInput) (*AllocationOutput, error) {
	if input.TitleSlug == "" {
		return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInvalidInput, "title slug cannot be empty", nil)
	}
	if len(input.Requirements) == 0 {
		return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInvalidInput, "requirements cannot be empty", nil)
	}
	if deps.TursoClient == nil {
		return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	if deps.Model == nil {
		return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil)
	}

	runID := input.RunID
	if runID == "" {
		runID = "run-" + uuid.NewString()[:8]
	}

	existingRun := runs.GetRun(ctx, deps.TursoClient, runID)
	if existingRun.IsErr() {
		rRes := runs.CreateRun(ctx, deps.TursoClient, &models.Run{
			Base:      models.Base{ID: runID},
			TitleSlug: input.TitleSlug,
			Trigger:   "allocation",
			Status:    models.RunStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if rRes.IsErr() {
			logger.Error("allocation: failed to persist initial run in turso",
				"run_id", runID,
				"title_slug", input.TitleSlug,
				"error", rRes.Error(),
			)
		}
	}

	candStepID := fmt.Sprintf("step-%s-candidates", runID)
	sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
		Base:      models.Base{ID: candStepID},
		RunID:     runID,
		Name:      "candidate_gathering",
		Status:    models.StepStatusRunning,
		StartedAt: time.Now().UTC(),
	})
	if sRes.IsErr() {
		logger.Error("allocation: failed to persist candidate_gathering step",
			"run_id", runID,
			"title_slug", input.TitleSlug,
			"step_id", candStepID,
			"error", sRes.Error(),
		)
	}

	type pairKey struct {
		Component string
		Market    string
	}
	seenPairs := make(map[pairKey]bool)
	candidatesByRequirement := make(map[string][]models.VendorCandidate)
	candidateCounts := make(map[string]int)

	for _, req := range input.Requirements {
		normComp := strings.ToUpper(strings.TrimSpace(req.Component))
		normMarket := strings.TrimSpace(req.Market)
		pk := pairKey{Component: normComp, Market: normMarket}
		key := fmt.Sprintf("%s|%s", normComp, normMarket)

		if seenPairs[pk] {
			continue
		}
		seenPairs[pk] = true

		cands, err := tools.FetchVendorCandidates(ctx, deps.TursoClient, deps.ClickHouse, tools.VendorCandidatesArgs{
			Component: normComp,
			Market:    normMarket,
		})
		now := time.Now().UTC()
		if err != nil {
			updStepRes := runs.UpdateStepStatus(ctx, deps.TursoClient, candStepID, models.StepStatusFailed, &now, map[string]any{"error": err.Error()})
			if updStepRes.IsErr() {
				logger.Warn("allocation: failed to update step status to failed", "run_id", runID, "step_id", candStepID, "error", updStepRes.Error())
			}
			updRunRes := runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
			if updRunRes.IsErr() {
				logger.Warn("allocation: failed to update run status to failed", "run_id", runID, "error", updRunRes.Error())
			}
			return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInternal, "failed to gather vendor candidates", err)
		}

		candidatesByRequirement[key] = cands
		candidateCounts[key] = len(cands)
	}

	now := time.Now().UTC()
	runs.UpdateStepStatus(ctx, deps.TursoClient, candStepID, models.StepStatusCompleted, &now, map[string]any{
		"requirements_count": len(input.Requirements),
		"candidate_pools":    candidateCounts,
	})

	selectStepID := fmt.Sprintf("step-%s-selection", runID)
	selStepRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
		Base:      models.Base{ID: selectStepID},
		RunID:     runID,
		Name:      "vendor_selection",
		Status:    models.StepStatusRunning,
		StartedAt: time.Now().UTC(),
	})
	if selStepRes.IsErr() {
		logger.Error("allocation: failed to persist vendor_selection step",
			"run_id", runID,
			"title_slug", input.TitleSlug,
			"step_id", selectStepID,
			"error", selStepRes.Error(),
		)
	}

	planRes := agent.SelectVendorsForPlan(ctx, deps.Model, input.TitleSlug, input.Requirements, candidatesByRequirement, input.HoursUntilPremiere)
	now = time.Now().UTC()
	if planRes.IsErr() {
		updStepRes := runs.UpdateStepStatus(ctx, deps.TursoClient, selectStepID, models.StepStatusFailed, &now, map[string]any{"error": planRes.Error().Error()})
		if updStepRes.IsErr() {
			logger.Warn("allocation: failed to update step status to failed", "run_id", runID, "step_id", selectStepID, "error", updStepRes.Error())
		}
		updRunRes := runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		if updRunRes.IsErr() {
			logger.Warn("allocation: failed to update run status to failed", "run_id", runID, "error", updRunRes.Error())
		}
		return nil, planRes.Error()
	}
	plan := planRes.Unwrap()

	assignedVendors := make([]string, 0, len(plan.Assignments))
	for i, assignment := range plan.Assignments {
		assignedVendors = append(assignedVendors, assignment.WinnerVendorID)
		resRes := runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
			Base:      models.Base{ID: fmt.Sprintf("res-%s-sel-%d", runID, i+1)},
			RunID:     runID,
			StepID:    selectStepID,
			Judge:     "vendor_selector",
			Outcome:   assignment.WinnerVendorID,
			Rationale: fmt.Sprintf("[%s/%s] %s", assignment.Component, assignment.Market, assignment.Rationale),
			Attempt:   1,
		})
		if resRes.IsErr() {
			logger.Warn("allocation: failed to record wf_result", "run_id", runID, "step_id", selectStepID, "index", i, "error", resRes.Error())
		}
	}

	runs.UpdateStepStatus(ctx, deps.TursoClient, selectStepID, models.StepStatusCompleted, &now, map[string]any{
		"assignments_count": len(plan.Assignments),
		"overall_summary":   plan.OverallSummary,
	})
	runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, map[string]any{
		"assignments_count": len(plan.Assignments),
		"vendors":           assignedVendors,
	})

	var firstDecision *agent.SelectionDecision
	if len(plan.Assignments) > 0 {
		first := plan.Assignments[0]
		firstDecision = &agent.SelectionDecision{
			WinnerVendorID:   first.WinnerVendorID,
			WinnerVendorName: first.WinnerVendorName,
			HourlyRateUSD:    first.HourlyRateUSD,
			TurnaroundHours:  first.TurnaroundHours,
			Rationale:        first.Rationale,
		}
	}

	return &AllocationOutput{
		Plan:     plan,
		Decision: firstDecision,
	}, nil
}
