package runs

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/ent/run"
	"github.com/elliot14A/fincher/internal/turso/ent/step"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// CreateRun inserts a new workflow run.
func CreateRun(ctx context.Context, client *ent.Client, r *models.Run) domainerrors.Result[*models.Run] {
	if r.ID == "" {
		r.ID = "run-" + uuid.NewString()[:8]
	}
	if r.StartedAt.IsZero() {
		r.StartedAt = time.Now().UTC()
	}

	if err := r.Validate(); err != nil {
		return domainerrors.Err[*models.Run](turso.NewError("runs.CreateRun", domainerrors.CodeInvalidInput, "invalid run data", err))
	}

	builder := client.Run.Create().
		SetID(r.ID).
		SetTitleSlug(r.TitleSlug).
		SetTrigger(r.Trigger).
		SetStatus(run.Status(r.Status)).
		SetStartedAt(r.StartedAt)

	if r.Metadata != nil {
		builder.SetMetadata(r.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Run](turso.MapEntError("runs.CreateRun", "run", r.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}

// CreateStep inserts an individual workflow execution step.
func CreateStep(ctx context.Context, client *ent.Client, s *models.Step) domainerrors.Result[*models.Step] {
	if s.ID == "" {
		s.ID = "step-" + uuid.NewString()[:8]
	}
	if s.StartedAt.IsZero() {
		s.StartedAt = time.Now().UTC()
	}

	if err := s.Validate(); err != nil {
		return domainerrors.Err[*models.Step](turso.NewError("runs.CreateStep", domainerrors.CodeInvalidInput, "invalid step data", err))
	}

	builder := client.Step.Create().
		SetID(s.ID).
		SetRunID(s.RunID).
		SetName(s.Name).
		SetStatus(step.Status(s.Status)).
		SetStartedAt(s.StartedAt)

	if s.Metadata != nil {
		builder.SetMetadata(s.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Step](turso.MapEntError("runs.CreateStep", "step", s.ID, err))
	}

	res := toDomainStep(created)
	return domainerrors.Ok(&res)
}

// CreateResult inserts a decision, verdict, or evaluation outcome.
func CreateResult(ctx context.Context, client *ent.Client, res *models.WfResult) domainerrors.Result[*models.WfResult] {
	if res.ID == "" {
		res.ID = "res-" + uuid.NewString()[:8]
	}
	if res.Attempt == 0 {
		res.Attempt = 1
	}

	if err := res.Validate(); err != nil {
		return domainerrors.Err[*models.WfResult](turso.NewError("runs.CreateResult", domainerrors.CodeInvalidInput, "invalid result data", err))
	}

	builder := client.WfResult.Create().
		SetID(res.ID).
		SetRunID(res.RunID).
		SetJudge(res.Judge).
		SetOutcome(res.Outcome).
		SetRationale(res.Rationale).
		SetAttempt(res.Attempt)

	if res.StepID != "" {
		builder.SetStepID(res.StepID)
	}
	if res.Metadata != nil {
		builder.SetMetadata(res.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.WfResult](turso.MapEntError("runs.CreateResult", "wf_result", res.ID, err))
	}

	domRes := toDomainResult(created)
	return domainerrors.Ok(&domRes)
}
