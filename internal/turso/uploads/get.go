package uploads

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Get fetches a binary upload by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Upload] {
	record, err := client.Upload.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Upload](turso.MapEntError("uploads.Get", "upload", id, err))
	}

	return domainerrors.Ok(&models.Upload{
		ID:        record.ID,
		Filename:  record.Filename,
		MimeType:  record.MimeType,
		Data:      record.Data,
		SizeBytes: record.SizeBytes,
		CreatedAt: record.CreatedAt,
	})
}
