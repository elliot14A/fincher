package packages

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	entmediapackage "github.com/elliot14A/fincher/pkg/ent/mediapackage"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Update modifies an existing package.
func Update(ctx context.Context, client *ent.Client, id string, input *models.UpdatePackageInput) domainerrors.Result[*models.Package] {
	if err := input.Validate(); err != nil {
		return domainerrors.Err[*models.Package](turso.NewError("packages.Update", domainerrors.CodeInvalidInput, "invalid package update input", err))
	}

	builder := client.MediaPackage.UpdateOneID(id)

	if input.Component != nil {
		builder.SetComponent(entmediapackage.Component(*input.Component))
	}
	if input.Language != nil {
		builder.SetLanguage(*input.Language)
	}
	if input.Version != nil {
		builder.SetVersion(*input.Version)
	}
	if input.VendorID != nil {
		builder.SetVendorID(*input.VendorID)
	}
	if input.DerivedFromMasterVersion != nil {
		builder.SetDerivedFromMasterVersion(*input.DerivedFromMasterVersion)
	}
	if input.RedeliveryCount != nil {
		builder.SetRedeliveryCount(*input.RedeliveryCount)
	}
	if input.Status != nil {
		builder.SetStatus(entmediapackage.Status(*input.Status))
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Package](turso.MapEntError("packages.Update", "package", id, err))
	}

	return domainerrors.Ok(toDomain(updated))
}
