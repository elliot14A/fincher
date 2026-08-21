package deliveries

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entdelivery "github.com/elliot14A/fincher/internal/turso/ent/delivery"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create inserts a new delivery.
func Create(ctx context.Context, client *ent.Client, d *models.Delivery) domainerrors.Result[*models.Delivery] {
	if err := d.Validate(); err != nil {
		return domainerrors.Err[*models.Delivery](turso.NewError("deliveries.Create", domainerrors.CodeInvalidInput, "invalid delivery data", err))
	}

	builder := client.Delivery.Create().
		SetID(d.ID).
		SetTitleID(d.TitleID).
		SetCountry(d.Country).
		SetStatus(entdelivery.Status(d.Status)).
		SetTargetDate(d.TargetDate)

	if d.Metadata != nil {
		builder.SetMetadata(d.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Delivery](turso.MapEntError("deliveries.Create", "delivery", d.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}
