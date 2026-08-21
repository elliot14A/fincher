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

// List fetches deliveries matching optional filters.
func List(ctx context.Context, client *ent.Client, filter ListFilter) domainerrors.Result[[]*models.Delivery] {
	query := client.Delivery.Query().Order(ent.Asc(entdelivery.FieldTargetDate))

	if filter.TitleID.IsSome() {
		query = query.Where(entdelivery.TitleIDEQ(filter.TitleID.Unwrap()))
	}
	if filter.Country.IsSome() {
		query = query.Where(entdelivery.CountryEQ(filter.Country.Unwrap()))
	}
	if filter.Status.IsSome() {
		query = query.Where(entdelivery.StatusEQ(entdelivery.Status(filter.Status.Unwrap())))
	}

	deliveriesList, err := query.All(ctx)
	if err != nil {
		return domainerrors.Err[[]*models.Delivery](turso.NewError("deliveries.List", domainerrors.CodeInternal, "failed to query deliveries", err))
	}

	return domainerrors.Ok(toDomainList(deliveriesList))
}
