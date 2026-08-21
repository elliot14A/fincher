package packages

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	entmediapackage "github.com/elliot14A/fincher/pkg/ent/mediapackage"
	"github.com/elliot14A/fincher/pkg/turso"
)

// ListFilter options for querying packages.
type ListFilter struct {
	TitleID   domainerrors.Option[string]
	VendorID  domainerrors.Option[string]
	Component domainerrors.Option[models.ComponentType]
	Status    domainerrors.Option[models.PackageStatus]
}

// List fetches media packages matching optional filters.
func List(ctx context.Context, client *ent.Client, filter ListFilter) domainerrors.Result[[]*models.Package] {
	query := client.MediaPackage.Query().Order(ent.Asc(entmediapackage.FieldCreatedAt))

	if filter.TitleID.IsSome() {
		query = query.Where(entmediapackage.TitleIDEQ(filter.TitleID.Unwrap()))
	}
	if filter.VendorID.IsSome() {
		query = query.Where(entmediapackage.VendorIDEQ(filter.VendorID.Unwrap()))
	}
	if filter.Component.IsSome() {
		query = query.Where(entmediapackage.ComponentEQ(entmediapackage.Component(filter.Component.Unwrap())))
	}
	if filter.Status.IsSome() {
		query = query.Where(entmediapackage.StatusEQ(entmediapackage.Status(filter.Status.Unwrap())))
	}

	packagesList, err := query.All(ctx)
	if err != nil {
		return domainerrors.Err[[]*models.Package](turso.NewError("packages.List", domainerrors.CodeInternal, "failed to query packages", err))
	}

	return domainerrors.Ok(toDomainList(packagesList))
}
