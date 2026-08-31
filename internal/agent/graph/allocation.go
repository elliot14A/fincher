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

	runID := input.RunID
	if runID == "" {
		runID = "run-" + uuid.NewString()[:8]
	}

	if deps.TursoClient != nil {
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
				logger.Error("failed to persist initial allocation run in turso", "run_id", runID, "error", rRes.Error())
			}
		}
	}

	// Stage 1: Candidate Gathering
	candStepID := fmt.Sprintf("step-%s-candidates", runID)
	if deps.TursoClient != nil {
		sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
			Base:      models.Base{ID: candStepID},
			RunID:     runID,
			Name:      "candidate_gathering",
			Status:    models.StepStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if sRes.IsErr() {
			logger.Error("failed to persist candidate_gathering step", "step_id", candStepID, "error", sRes.Error())
		}
	}

	candidates, err := tools.FetchVendorCandidates(ctx, deps.TursoClient, deps.ClickHouse, tools.VendorCandidatesArgs{
		Component: input.Component,
		Specialty: input.Component,
	})
	now := time.Now().UTC()
	if err != nil {
		if deps.TursoClient != nil {
			runs.UpdateStepStatus(ctx, deps.TursoClient, candStepID, models.StepStatusFailed, &now, map[string]any{"error": err.Error()})
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		}
		return nil, domainerrors.NewWithOp("graph.ExecuteAllocation", domainerrors.CodeInternal, "failed to gather vendor candidates", err)
	}

	if deps.TursoClient != nil {
		runs.UpdateStepStatus(ctx, deps.TursoClient, candStepID, models.StepStatusCompleted, &now, map[string]any{
			"candidates_count": len(candidates),
		})
	}

	// Stage 2: Vendor Selection
	selectStepID := fmt.Sprintf("step-%s-selection", runID)
	if deps.TursoClient != nil {
		sRes := runs.CreateStep(ctx, deps.TursoClient, &models.Step{
			Base:      models.Base{ID: selectStepID},
			RunID:     runID,
			Name:      "vendor_selection",
			Status:    models.StepStatusRunning,
			StartedAt: time.Now().UTC(),
		})
		if sRes.IsErr() {
			logger.Error("failed to persist vendor_selection step", "step_id", selectStepID, "error", sRes.Error())
		}
	}

	decisionRes := agent.SelectVendor(ctx, deps.Model, input.TitleSlug, input.Component, candidates, input.HoursUntilPremiere)
	now = time.Now().UTC()
	if decisionRes.IsErr() {
		if deps.TursoClient != nil {
			runs.UpdateStepStatus(ctx, deps.TursoClient, selectStepID, models.StepStatusFailed, &now, map[string]any{"error": decisionRes.Error().Error()})
			runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusFailed, &now, nil)
		}
		return nil, decisionRes.Error()
	}
	decision := decisionRes.Unwrap()

	if deps.TursoClient != nil {
		runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
			Base:      models.Base{ID: fmt.Sprintf("res-%s-selection", runID)},
			RunID:     runID,
			StepID:    selectStepID,
			Judge:     "vendor_selector",
			Outcome:   decision.WinnerVendorID,
			Rationale: decision.Rationale,
			Attempt:   1,
		})
		runs.UpdateStepStatus(ctx, deps.TursoClient, selectStepID, models.StepStatusCompleted, &now, map[string]any{
			"winner_vendor_id": decision.WinnerVendorID,
			"hourly_rate_usd":  decision.HourlyRateUSD,
			"turnaround_hours": decision.TurnaroundHours,
		})
		runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &now, map[string]any{
			"winner_vendor_id": decision.WinnerVendorID,
		})
	}

	return &AllocationOutput{
		Decision: decision,
	}, nil
}
