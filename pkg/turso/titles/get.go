package titles

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Get fetches a single title by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*ent.Title] {
	t, err := client.Title.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*ent.Title](turso.MapEntError("titles.Get", "title", id, err))
	}
	return domainerrors.Ok(t)
}
