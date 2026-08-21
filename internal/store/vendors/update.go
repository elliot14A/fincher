package vendors

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Update modifies an existing vendor.
func Update(ctx context.Context, client *ent.Client, id string, input *models.UpdateVendorInput) domainerrors.Result[*models.Vendor] {
	if err := input.Validate(); err != nil {
		return domainerrors.Err[*models.Vendor](turso.NewError("vendors.Update", domainerrors.CodeInvalidInput, "invalid vendor update input", err))
	}

	builder := client.Vendor.UpdateOneID(id)

	if input.Name != nil {
		builder.SetName(*input.Name)
	}
	if input.Specialty != nil {
		builder.SetSpecialty(*input.Specialty)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Vendor](turso.MapEntError("vendors.Update", "vendor", id, err))
	}

	return domainerrors.Ok(toDomain(updated))
}
