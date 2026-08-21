package masters

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// Delete removes a master by ID.
func Delete(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	err := client.Master.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return domainerrors.Err[bool](turso.MapEntError("masters.Delete", "master", id, err))
	}
	return domainerrors.Ok(true)
}
