package deliveries

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entdelivery "github.com/elliot14A/fincher/internal/turso/ent/delivery"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Update modifies an existing delivery.
func Update(ctx context.Context, client *ent.Client, id string, input *models.UpdateDeliveryInput) domainerrors.Result[*models.Delivery] {
	if err := input.Validate(); err != nil {
		return domainerrors.Err[*models.Delivery](turso.NewError("deliveries.Update", domainerrors.CodeInvalidInput, "invalid delivery update input", err))
	}

	builder := client.Delivery.UpdateOneID(id)

	if input.Country != nil {
		builder.SetCountry(*input.Country)
	}
	if input.Status != nil {
		builder.SetStatus(entdelivery.Status(*input.Status))
	}
	if input.TargetDate != nil {
		builder.SetTargetDate(*input.TargetDate)
	}
	if input.Metadata != nil {
		builder.SetMetadata(input.Metadata)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Delivery](turso.MapEntError("deliveries.Update", "delivery", id, err))
	}

	return domainerrors.Ok(toDomain(updated))
}
