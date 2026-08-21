package packages

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Get fetches a single media package by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Package] {
	p, err := client.MediaPackage.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Package](turso.MapEntError("packages.Get", "package", id, err))
	}
	return domainerrors.Ok(toDomain(p))
}
