package masters

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Get fetches a single master by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Master] {
	m, err := client.Master.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Master](turso.MapEntError("masters.Get", "master", id, err))
	}
	return domainerrors.Ok(toDomain(m))
}
