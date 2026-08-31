package runs

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entrun "github.com/elliot14A/fincher/internal/turso/ent/run"
	entstep "github.com/elliot14A/fincher/internal/turso/ent/step"
	entresult "github.com/elliot14A/fincher/internal/turso/ent/wfresult"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// GetRun fetches a single workflow Run by ID, eagerly loading its steps and results.
func GetRun(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Run] {
	r, err := client.Run.
		Query().
		Where(entrun.IDEQ(id)).
		WithSteps(func(sq *ent.StepQuery) {
			sq.Order(ent.Asc(entstep.FieldStartedAt)).
				WithResults(func(rq *ent.WfResultQuery) {
					rq.Order(ent.Asc(entresult.FieldCreatedAt))
				})
		}).
		WithResults(func(rq *ent.WfResultQuery) {
			rq.Order(ent.Asc(entresult.FieldCreatedAt))
		}).
		Only(ctx)

	if err != nil {
		return domainerrors.Err[*models.Run](turso.MapEntError("runs.GetRun", "run", id, err))
	}

	return domainerrors.Ok(toDomain(r))
}
