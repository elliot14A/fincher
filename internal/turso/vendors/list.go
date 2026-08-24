package vendors

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entvendor "github.com/elliot14A/fincher/internal/turso/ent/vendor"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List fetches paginated vendors, optionally filtered by specialty and search term.
func List(ctx context.Context, client *ent.Client, specialtyFilter domainerrors.Option[string], p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Vendor]] {
	query := client.Vendor.Query()

	if specialtyFilter.IsSome() {
		query = query.Where(entvendor.SpecialtyEQ(specialtyFilter.Unwrap()))
	}
	if p.Search != "" {
		query = query.Where(entvendor.NameContainsFold(p.Search))
	}

	query = query.Order(turso.OrderBy(p, ent.Asc(entvendor.FieldName), ent.Desc(entvendor.FieldName)))

	return turso.Paginate(
		ctx,
		"vendors.List",
		p,
		query.Count,
		func(ctx context.Context, limit, offset int) ([]*ent.Vendor, error) {
			return query.Limit(limit).Offset(offset).All(ctx)
		},
		toDomainList,
	)
}
