package vendors

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	entvendor "github.com/elliot14A/fincher/pkg/ent/vendor"
	"github.com/elliot14A/fincher/pkg/turso"
)

// List fetches all vendors, optionally filtered by specialty.
func List(ctx context.Context, client *ent.Client, specialtyFilter domainerrors.Option[string]) domainerrors.Result[[]*models.Vendor] {
	query := client.Vendor.Query().Order(ent.Asc(entvendor.FieldName))

	if specialtyFilter.IsSome() {
		query = query.Where(entvendor.SpecialtyEQ(specialtyFilter.Unwrap()))
	}

	vendorsList, err := query.All(ctx)
	if err != nil {
		return domainerrors.Err[[]*models.Vendor](turso.NewError("vendors.List", domainerrors.CodeInternal, "failed to query vendors", err))
	}

	return domainerrors.Ok(toDomainList(vendorsList))
}
