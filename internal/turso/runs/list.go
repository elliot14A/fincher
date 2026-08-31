package runs

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entrun "github.com/elliot14A/fincher/internal/turso/ent/run"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// ListFilter holds optional query filters for listing workflow runs.
type ListFilter struct {
	TitleSlug domainerrors.Option[string]
	Workflow  domainerrors.Option[string]
	Status    domainerrors.Option[models.RunStatus]
}

// ListRuns fetches paginated workflow runs, optionally filtered by title slug, workflow trigger, and status.
func ListRuns(ctx context.Context, client *ent.Client, filter ListFilter, p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Run]] {
	query := client.Run.Query()

	if filter.TitleSlug.IsSome() {
		query = query.Where(entrun.TitleSlugEQ(filter.TitleSlug.Unwrap()))
	}
	if filter.Workflow.IsSome() {
		query = query.Where(entrun.TriggerEQ(filter.Workflow.Unwrap()))
	}
	if filter.Status.IsSome() {
		query = query.Where(entrun.StatusEQ(entrun.Status(filter.Status.Unwrap())))
	}
	if p.Search != "" {
		query = query.Where(entrun.TriggerContainsFold(p.Search))
	}

	query = query.Order(turso.OrderBy(p, ent.Desc(entrun.FieldStartedAt), ent.Asc(entrun.FieldStartedAt)))

	return turso.Paginate(
		ctx,
		"runs.ListRuns",
		p,
		query.Count,
		func(ctx context.Context, limit, offset int) ([]*ent.Run, error) {
			return query.
				Limit(limit).
				Offset(offset).
				WithSteps().
				WithResults().
				All(ctx)
		},
		toDomainList,
	)
}
