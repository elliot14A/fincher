package uploads

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// Delete removes an upload record by ID.
func Delete(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	if err := client.Upload.DeleteOneID(id).Exec(ctx); err != nil {
		return domainerrors.Err[bool](turso.MapEntError("uploads.Delete", "upload", id, err))
	}
	return domainerrors.Ok(true)
}
