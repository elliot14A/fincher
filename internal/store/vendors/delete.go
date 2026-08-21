package vendors

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Delete removes a vendor by ID.
func Delete(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	err := client.Vendor.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return domainerrors.Err[bool](turso.MapEntError("vendors.Delete", "vendor", id, err))
	}
	return domainerrors.Ok(true)
}
