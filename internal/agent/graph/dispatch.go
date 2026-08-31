package graph

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// DispatchIncident handles idempotent initialization and asynchronous background execution of the incident workflow.
// It persists the root Run row upfront (with RunStatusRunning) so SSE stream clients can subscribe immediately.
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

	// Idempotency check: if run already exists, return existing run and do not re-dispatch
	existing := runs.GetRun(ctx, deps.TursoClient, runID)
	if existing.IsOk() {
		return existing.Unwrap(), false, nil
	}

	// Resolve urgency deadline
	if input.HoursUntilPremiere <= 0 {
		input.HoursUntilPremiere = tursotitles.ResolveHoursUntilPremiere(ctx, deps.TursoClient, titleSlug, 0)
	}

	// Persist initial Run row upfront with RUNNING status
	initialRun, err := ensureRun(ctx, deps.TursoClient, runID, titleSlug, "incident")
	if err != nil {
		return nil, false, err
	}

	// Launch workflow in background goroutine
	go func(in IncidentInput) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if deps.MaxAttempts <= 0 {
			deps.MaxAttempts = DefaultMaxRemediationAttempts
		}

		_, err := ExecuteIncident(bgCtx, deps, in)
		if err != nil {
			logger.Error("incident workflow execution failed", "run_id", in.RunID, "error", err)
		}
	}(input)

	return initialRun, true, nil
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

	existing := runs.GetRun(ctx, deps.TursoClient, runID)
	if existing.IsOk() {
		return existing.Unwrap(), false, nil
	}

	if input.HoursUntilPremiere <= 0 {
		input.HoursUntilPremiere = tursotitles.ResolveHoursUntilPremiere(ctx, deps.TursoClient, input.TitleSlug, 0)
	}

	initialRun, err := ensureRun(ctx, deps.TursoClient, runID, input.TitleSlug, "allocation")
	if err != nil {
		return nil, false, err
	}

	go func(in AllocationInput) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		_, err := ExecuteAllocation(bgCtx, deps, in)
		if err != nil {
			logger.Error("allocation workflow execution failed", "run_id", in.RunID, "error", err)
		}
	}(input)

	return initialRun, true, nil
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

	existing := runs.GetRun(ctx, deps.TursoClient, runID)
	if existing.IsOk() {
		return existing.Unwrap(), false, nil
	}

	initialRun, err := ensureRun(ctx, deps.TursoClient, runID, titleSlug, "resolution")
	if err != nil {
		return nil, false, err
	}

	go func(in ResolutionInput) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		_, err := ExecuteResolution(bgCtx, deps, in)
		if err != nil {
			logger.Error("resolution workflow execution failed", "run_id", in.RunID, "error", err)
		}
	}(input)

	return initialRun, true, nil
}
