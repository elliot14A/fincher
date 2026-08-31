package runs

import (
	"context"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/ent/run"
	"github.com/elliot14A/fincher/internal/turso/ent/step"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// UpdateRunStatus updates the status, completion time, and metadata of a workflow Run.
func UpdateRunStatus(ctx context.Context, client *ent.Client, id string, status models.RunStatus, endedAt *time.Time, metadata map[string]any) domainerrors.Result[*models.Run] {
	if client == nil {
		return domainerrors.Err[*models.Run](turso.NewError("runs.UpdateRunStatus", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil))
	}

	builder := client.Run.UpdateOneID(id).
		SetStatus(run.Status(status))

	if endedAt != nil {
		builder.SetEndedAt(*endedAt)
	}
	if metadata != nil {
		builder.SetMetadata(metadata)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Run](turso.MapEntError("runs.UpdateRunStatus", "run", id, err))
	}

	return domainerrors.Ok(toDomain(updated))
}

// UpdateRunStatusIfRunning updates a run's status atomically only if it is currently in RUNNING status.
// If the run has already reached a terminal state (COMPLETED, FAILED, ESCALATED), it returns Ok with the existing state without mutating it.
func UpdateRunStatusIfRunning(ctx context.Context, client *ent.Client, id string, status models.RunStatus, endedAt *time.Time, metadata map[string]any) domainerrors.Result[*models.Run] {
	if client == nil {
		return domainerrors.Err[*models.Run](turso.NewError("runs.UpdateRunStatusIfRunning", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil))
	}

	r, err := client.Run.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Run](turso.MapEntError("runs.UpdateRunStatusIfRunning", "run", id, err))
	}

	if r.Status == run.StatusCOMPLETED || r.Status == run.StatusFAILED || r.Status == run.StatusESCALATED {
		return domainerrors.Ok(toDomain(r))
	}

	builder := client.Run.UpdateOneID(id).
		Where(run.StatusEQ(run.StatusRUNNING)).
		SetStatus(run.Status(status))

	if endedAt != nil {
		builder.SetEndedAt(*endedAt)
	}
	if metadata != nil {
		builder.SetMetadata(metadata)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			current, getErr := client.Run.Get(ctx, id)
			if getErr == nil {
				return domainerrors.Ok(toDomain(current))
			}
		}
		return domainerrors.Err[*models.Run](turso.MapEntError("runs.UpdateRunStatusIfRunning", "run", id, err))
	}

	return domainerrors.Ok(toDomain(updated))
}

// UpdateStepStatus updates the status, completion time, and metadata of a workflow Step.
func UpdateStepStatus(ctx context.Context, client *ent.Client, id string, status models.StepStatus, endedAt *time.Time, metadata map[string]any) domainerrors.Result[*models.Step] {
	if client == nil {
		return domainerrors.Err[*models.Step](turso.NewError("runs.UpdateStepStatus", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil))
	}

	builder := client.Step.UpdateOneID(id).
		SetStatus(step.Status(status))

	if endedAt != nil {
		builder.SetEndedAt(*endedAt)
	}
	if metadata != nil {
		builder.SetMetadata(metadata)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Step](turso.MapEntError("runs.UpdateStepStatus", "step", id, err))
	}

	res := toDomainStep(updated)
	return domainerrors.Ok(&res)
}
