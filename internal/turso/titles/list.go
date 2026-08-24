package titles

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	enttitle "github.com/elliot14A/fincher/internal/turso/ent/title"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List fetches paginated titles, optionally filtered by overall status and search term.
func List(ctx context.Context, client *ent.Client, statusFilter domainerrors.Option[models.TitleStatus], p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Title]] {
	query := client.Title.Query()

	if statusFilter.IsSome() {
		query = query.Where(enttitle.OverallStatusEQ(enttitle.OverallStatus(statusFilter.Unwrap())))
	}
	if p.Search != "" {
		query = query.Where(enttitle.NameContainsFold(p.Search))
	}

	query = query.Order(turso.OrderBy(p, ent.Asc(enttitle.FieldPremiereDate), ent.Desc(enttitle.FieldPremiereDate)))

	return turso.Paginate(
		ctx,
		"titles.List",
		p,
		query.Count,
		func(ctx context.Context, limit, offset int) ([]*ent.Title, error) {
			return query.Limit(limit).Offset(offset).All(ctx)
		},
		toDomainList,
	)
}
