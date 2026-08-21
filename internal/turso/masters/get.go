package masters

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Get fetches a single master by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Master] {
	m, err := client.Master.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Master](turso.MapEntError("masters.Get", "master", id, err))
	}
	return domainerrors.Ok(toDomain(m))
}
