package packages

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entmediapackage "github.com/elliot14A/fincher/internal/turso/ent/mediapackage"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create inserts a new media package.
func Create(ctx context.Context, client *ent.Client, p *models.Package) domainerrors.Result[*models.Package] {
	if err := p.Validate(); err != nil {
		return domainerrors.Err[*models.Package](turso.NewError("packages.Create", domainerrors.CodeInvalidInput, "invalid package data", err))
	}

	builder := client.MediaPackage.Create().
		SetID(p.ID).
		SetTitleID(p.TitleID).
		SetComponent(entmediapackage.Component(p.Component)).
		SetLanguage(p.Language).
		SetVersion(p.Version).
		SetVendorID(p.VendorID).
		SetDerivedFromMasterVersion(p.DerivedFromMasterVersion).
		SetRedeliveryCount(p.RedeliveryCount).
		SetStatus(entmediapackage.Status(p.Status))

	if p.Metadata != nil {
		builder.SetMetadata(p.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Package](turso.MapEntError("packages.Create", "package", p.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}
