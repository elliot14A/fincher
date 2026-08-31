package packages

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entmediapackage "github.com/elliot14A/fincher/internal/turso/ent/mediapackage"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// ListFilter options for querying packages.
type ListFilter struct {
	TitleID   domainerrors.Option[string]
	VendorID  domainerrors.Option[string]
	Component domainerrors.Option[models.ComponentType]
	Status    domainerrors.Option[models.PackageStatus]
	Market    domainerrors.Option[string]
}

// List fetches paginated media packages matching optional filters.
func List(ctx context.Context, client *ent.Client, filter ListFilter, p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Package]] {
	query := client.MediaPackage.Query()

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
	if filter.Market.IsSome() {
		query = query.Where(entmediapackage.MarketEQ(filter.Market.Unwrap()))
	}
	if p.Search != "" {
		query = query.Where(entmediapackage.LanguageContainsFold(p.Search))
	}

	query = query.Order(turso.OrderBy(p, ent.Asc(entmediapackage.FieldCreatedAt), ent.Desc(entmediapackage.FieldCreatedAt)))

	return turso.Paginate(
		ctx,
		"packages.List",
		p,
		query.Count,
		func(ctx context.Context, limit, offset int) ([]*ent.MediaPackage, error) {
			return query.Limit(limit).Offset(offset).All(ctx)
		},
		toDomainList,
	)
}
