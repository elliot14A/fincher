package vendors

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Get fetches a single vendor by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Vendor] {
	v, err := client.Vendor.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Vendor](turso.MapEntError("vendors.Get", "vendor", id, err))
	}
	return domainerrors.Ok(toDomain(v))
}
