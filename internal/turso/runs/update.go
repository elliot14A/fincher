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

// UpdateStepStatus updates the status, completion time, and metadata of a workflow Step.
func UpdateStepStatus(ctx context.Context, client *ent.Client, id string, status models.StepStatus, endedAt *time.Time, metadata map[string]any) domainerrors.Result[*models.Step] {
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
