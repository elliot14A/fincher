package masters

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entmaster "github.com/elliot14A/fincher/internal/turso/ent/master"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List fetches paginated masters, optionally filtered by title_id.
func List(ctx context.Context, client *ent.Client, titleIDFilter domainerrors.Option[string], p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Master]] {
	query := client.Master.Query()

	if titleIDFilter.IsSome() {
		query = query.Where(entmaster.TitleIDEQ(titleIDFilter.Unwrap()))
	}
	if p.Search != "" {
		query = query.Where(entmaster.VersionContainsFold(p.Search))
	}

	query = query.Order(turso.OrderBy(p, ent.Asc(entmaster.FieldCreatedAt), ent.Desc(entmaster.FieldCreatedAt)))

	return turso.Paginate(
		ctx,
		"masters.List",
		p,
		query.Count,
		func(ctx context.Context, limit, offset int) ([]*ent.Master, error) {
			return query.Limit(limit).Offset(offset).All(ctx)
		},
		toDomainList,
	)
}
