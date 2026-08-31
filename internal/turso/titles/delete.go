package titles

import (
	"context"
	"strings"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// Delete removes a title by ID and cascades poster upload cleanup if present.
func Delete(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	t, err := client.Title.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[bool](turso.MapEntError("titles.Delete", "title", id, err))
	}

	posterURL := ""
	if t.Metadata != nil {
		if u, ok := t.Metadata["poster_url"].(string); ok {
			posterURL = u
		} else if u, ok := t.Metadata["avatar_url"].(string); ok {
			posterURL = u
		}
	}

	if err := client.Title.DeleteOneID(id).Exec(ctx); err != nil {
		return domainerrors.Err[bool](turso.MapEntError("titles.Delete", "title", id, err))
	}

	// Clean up associated internal upload blob if referenced
	if strings.HasPrefix(posterURL, "/api/uploads/") {
		uploadID := strings.TrimPrefix(posterURL, "/api/uploads/")
		_ = client.Upload.DeleteOneID(uploadID).Exec(ctx)
	}

	return domainerrors.Ok(true)
}
