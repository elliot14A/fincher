package packages

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Get fetches a single media package by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Package] {
	p, err := client.MediaPackage.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Package](turso.MapEntError("packages.Get", "package", id, err))
	}
	return domainerrors.Ok(toDomain(p))
}
