package packages

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Delete removes a media package by ID.
func Delete(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	err := client.MediaPackage.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return domainerrors.Err[bool](turso.MapEntError("packages.Delete", "package", id, err))
	}
	return domainerrors.Ok(true)
}
