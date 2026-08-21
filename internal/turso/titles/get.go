package titles

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Get fetches a single title by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Title] {
	t, err := client.Title.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Title](turso.MapEntError("titles.Get", "title", id, err))
	}
	return domainerrors.Ok(toDomain(t))
}
