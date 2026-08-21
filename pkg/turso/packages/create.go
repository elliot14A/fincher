package packages

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	entmediapackage "github.com/elliot14A/fincher/pkg/ent/mediapackage"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Create inserts a new media package.
func Create(ctx context.Context, client *ent.Client, p *models.Package) domainerrors.Result[*models.Package] {
	if err := p.Validate(); err != nil {
		return domainerrors.Err[*models.Package](turso.NewError("packages.Create", domainerrors.CodeInvalidInput, "invalid package data", err))
	}

	created, err := client.MediaPackage.Create().
		SetID(p.ID).
		SetTitleID(p.TitleID).
		SetComponent(entmediapackage.Component(p.Component)).
		SetLanguage(p.Language).
		SetVersion(p.Version).
		SetVendorID(p.VendorID).
		SetDerivedFromMasterVersion(p.DerivedFromMasterVersion).
		SetRedeliveryCount(p.RedeliveryCount).
		SetStatus(entmediapackage.Status(p.Status)).
		Save(ctx)

	if err != nil {
		return domainerrors.Err[*models.Package](turso.MapEntError("packages.Create", "package", p.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}
