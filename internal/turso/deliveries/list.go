package deliveries

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entdelivery "github.com/elliot14A/fincher/internal/turso/ent/delivery"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// ListFilter holds optional query params for listing deliveries.
type ListFilter struct {
	TitleID domainerrors.Option[string]
	Country domainerrors.Option[string]
	Status  domainerrors.Option[models.DeliveryStatus]
}

// List fetches paginated deliveries matching optional filters.
func List(ctx context.Context, client *ent.Client, filter ListFilter, p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Delivery]] {
	query := client.Delivery.Query()

	if filter.TitleID.IsSome() {
		query = query.Where(entdelivery.TitleIDEQ(filter.TitleID.Unwrap()))
	}
	if filter.Country.IsSome() {
		query = query.Where(entdelivery.CountryEQ(filter.Country.Unwrap()))
	}
	if filter.Status.IsSome() {
		query = query.Where(entdelivery.StatusEQ(entdelivery.Status(filter.Status.Unwrap())))
	}
	if p.Search != "" {
		query = query.Where(entdelivery.CountryContainsFold(p.Search))
	}

	query = query.Order(turso.OrderBy(p, ent.Asc(entdelivery.FieldTargetDate), ent.Desc(entdelivery.FieldTargetDate)))

	return turso.Paginate(
		ctx,
		"deliveries.List",
		p,
		query.Count,
		func(ctx context.Context, limit, offset int) ([]*ent.Delivery, error) {
			return query.Limit(limit).Offset(offset).All(ctx)
		},
		toDomainList,
	)
}
