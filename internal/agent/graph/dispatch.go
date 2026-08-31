package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/recovery"
)

// markRunFailedOnPanic returns a panic handler closure that transitions a non-terminal run to FAILED.
// If the run has already reached a terminal state (COMPLETED, FAILED, ESCALATED), the update is skipped
// to prevent clobbering valid terminal outcomes or error metadata.
func markRunFailedOnPanic(client *ent.Client, runID, titleSlug string) func(r any, stack string) {
	return func(r any, stack string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		now := time.Now().UTC()
		updRes := runs.UpdateRunStatusIfRunning(bgCtx, client, runID, models.RunStatusFailed, &now, map[string]any{
			"panic": fmt.Sprint(r),
		})
		if updRes.IsErr() {
			logger.Error("dispatch: failed to mark run failed after panic",
				"run_id", runID,
				"title_slug", titleSlug,
				"error", updRes.Error(),
			)
		} else {
			runObj := updRes.Unwrap()
			if runObj.Status != models.RunStatusFailed {
				logger.Warn("dispatch: run already in terminal state, skipping panic failure transition",
					"run_id", runID,
					"title_slug", titleSlug,
					"current_status", string(runObj.Status),
					"panic", fmt.Sprint(r),
				)
			}
		}
	}
}

// dispatchAsyncWorkflow performs idempotent run initialization, creates root Run row, and launches
// the workflow in a panic-protected background goroutine.
func dispatchAsyncWorkflow(
	ctx context.Context,
	client *ent.Client,
	runID, titleSlug, trigger string,
	workflowFn func(bgCtx context.Context) error,
) (*models.Run, bool, error) {
	// Idempotency check: if run already exists, return existing run and do not re-dispatch
	existing := runs.GetRun(ctx, client, runID)
	if existing.IsOk() {
		return existing.Unwrap(), false, nil
	}

	// Persist initial Run row upfront with RUNNING status
	initialRun, err := ensureRun(ctx, client, runID, titleSlug, trigger)
	if err != nil {
		return nil, false, err
	}

	safeName := fmt.Sprintf("dispatch.%s.run=%s.title=%s", trigger, runID, titleSlug)
	recovery.SafeGo(safeName, func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := workflowFn(bgCtx); err != nil {
			logger.Error(fmt.Sprintf("%s workflow execution failed", trigger),
				"run_id", runID,
				"title_slug", titleSlug,
				"error", err,
			)
		}
	}, markRunFailedOnPanic(client, runID, titleSlug))

	return initialRun, true, nil
}

// DispatchIncident handles idempotent initialization and asynchronous background execution of the incident workflow.
func DispatchIncident(ctx context.Context, deps IncidentGraphDeps, input IncidentInput) (*models.Run, bool, error) {
	if input.Event == nil {
		return nil, false, domainerrors.NewWithOp("graph.DispatchIncident", domainerrors.CodeInvalidInput, "input event cannot be nil", nil)
	}
	if deps.TursoClient == nil {
		return nil, false, domainerrors.NewWithOp("graph.DispatchIncident", domainerrors.CodeInvalidInput, "turso client is required for incident workflow", nil)
	}

	runID := input.RunID
	if runID == "" {
		if input.Event.ID != "" {
			runID = "run-" + input.Event.ID
		} else {
			runID = "run-" + uuid.NewString()[:8]
		}
	}
	input.RunID = runID

	titleSlug := input.Event.Subject
	if titleSlug == "" {
		titleSlug = models.DefaultTitleAgnosticSentinel
	}

	if input.HoursUntilPremiere <= 0 {
		input.HoursUntilPremiere = tursotitles.ResolveHoursUntilPremiere(ctx, deps.TursoClient, titleSlug, 0)
	}

	if deps.MaxAttempts <= 0 {
		deps.MaxAttempts = DefaultMaxRemediationAttempts
	}

	return dispatchAsyncWorkflow(ctx, deps.TursoClient, runID, titleSlug, "incident", func(bgCtx context.Context) error {
		_, err := ExecuteIncident(bgCtx, deps, input)
		return err
	})
}

// DispatchAllocation handles idempotent initialization and asynchronous background execution of the vendor allocation workflow.
func DispatchAllocation(ctx context.Context, deps AllocationGraphDeps, input AllocationInput) (*models.Run, bool, error) {
	if deps.TursoClient == nil {
		return nil, false, domainerrors.NewWithOp("graph.DispatchAllocation", domainerrors.CodeInvalidInput, "turso client is required for allocation workflow", nil)
	}

	if input.TitleSlug == "" {
		input.TitleSlug = models.DefaultTitleAgnosticSentinel
	}
	if input.Component == "" {
		input.Component = "AUDIO"
	}

	runID := input.RunID
	if runID == "" {
		runID = "run-" + uuid.NewString()[:8]
	}
	input.RunID = runID

	if input.HoursUntilPremiere <= 0 {
		input.HoursUntilPremiere = tursotitles.ResolveHoursUntilPremiere(ctx, deps.TursoClient, input.TitleSlug, 0)
	}

	return dispatchAsyncWorkflow(ctx, deps.TursoClient, runID, input.TitleSlug, "allocation", func(bgCtx context.Context) error {
		_, err := ExecuteAllocation(bgCtx, deps, input)
		return err
	})
}

// DispatchResolution handles idempotent initialization and background execution of the closed-loop resolution workflow.
func DispatchResolution(ctx context.Context, deps ResolutionDeps, input ResolutionInput) (*models.Run, bool, error) {
	if deps.TursoClient == nil {
		return nil, false, domainerrors.NewWithOp("graph.DispatchResolution", domainerrors.CodeInvalidInput, "turso client is required for resolution workflow", nil)
	}

	runID := input.RunID
	if runID == "" {
		if input.Event != nil && input.Event.ID != "" {
			runID = "run-" + input.Event.ID
		} else {
			runID = "run-" + uuid.NewString()[:8]
		}
	}
	input.RunID = runID

	titleSlug := input.TitleSlug
	if titleSlug == "" && input.Event != nil {
		titleSlug = input.Event.Subject
	}
	if titleSlug == "" {
		titleSlug = models.DefaultTitleAgnosticSentinel
	}

	return dispatchAsyncWorkflow(ctx, deps.TursoClient, runID, titleSlug, "resolution", func(bgCtx context.Context) error {
		_, err := ExecuteResolution(bgCtx, deps, input)
		return err
	})
}
