package dependencies

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Get fetches a single dependency by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Dependency] {
	d, err := client.Dependency.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Dependency](turso.MapEntError("dependencies.Get", "dependency", id, err))
	}
	return domainerrors.Ok(toDomain(d))
}
