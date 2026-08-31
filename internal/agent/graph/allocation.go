package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/runs"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// ExecuteAllocation evaluates candidate vendors and selects the optimal partner based on
// SLA turnaround feasibility, quality floor constraints, and commercial rate.
//
// It persists candidate_gathering and vendor_selection Step and WfResult rows into Turso.
func ExecuteAllocation(ctx context.Context, deps AllocationGraphDeps, input AllocationInput) (*AllocationOutput, error) {
	if input.TitleSlug == "" {
		return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInvalidInput, "title slug cannot be empty", nil)
	}
	if input.Component == "" {
		return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInvalidInput, "component cannot be empty", nil)
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

	// Stage 1: Candidate Gathering
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

	candidates, err := tools.FetchVendorCandidates(ctx, deps.TursoClient, deps.ClickHouse, tools.VendorCandidatesArgs{
		Component: input.Component,
		Specialty: input.Component,
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

	runs.UpdateStepStatus(ctx, deps.TursoClient, candStepID, models.StepStatusCompleted, &now, map[string]any{
		"candidates_count": len(candidates),
	})

	// Stage 2: Vendor Selection
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

	decisionRes := agent.SelectVendor(ctx, deps.Model, input.TitleSlug, input.Component, candidates, input.HoursUntilPremiere)
	now = time.Now().UTC()
	if decisionRes.IsErr() {
		updStepRes := runs.UpdateStepStatus(ctx, deps.TursoClient, selectStepID, models.StepStatusFailed, &now, map[string]any{"error": decisionRes.Error().Error()})
		if updStepRes.IsErr() {
			logger.Warn("allocation: failed to update step status to failed", "run_id", runID, "step_id", selectStepID, "error", updStepRes.Error())
		}
		updRunRes := runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		if updRunRes.IsErr() {
			logger.Warn("allocation: failed to update run status to failed", "run_id", runID, "error", updRunRes.Error())
		}
		return nil, decisionRes.Error()
	}
	decision := decisionRes.Unwrap()

	resRes := runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
		Base:      models.Base{ID: fmt.Sprintf("res-%s-selection", runID)},
		RunID:     runID,
		StepID:    selectStepID,
		Judge:     "vendor_selector",
		Outcome:   decision.WinnerVendorID,
		Rationale: decision.Rationale,
		Attempt:   1,
	})
	if resRes.IsErr() {
		logger.Warn("allocation: failed to record wf_result", "run_id", runID, "step_id", selectStepID, "error", resRes.Error())
	}

	runs.UpdateStepStatus(ctx, deps.TursoClient, selectStepID, models.StepStatusCompleted, &now, map[string]any{
		"winner_vendor_id": decision.WinnerVendorID,
		"hourly_rate_usd":  decision.HourlyRateUSD,
		"turnaround_hours": decision.TurnaroundHours,
	})
	runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, map[string]any{
		"winner_vendor_id": decision.WinnerVendorID,
	})

	return &AllocationOutput{
		Decision: decision,
	}, nil
}
