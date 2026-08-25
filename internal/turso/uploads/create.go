package uploads

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create persists a binary upload record in SQLite.
func Create(ctx context.Context, client *ent.Client, u *models.Upload) domainerrors.Result[*models.Upload] {
	if err := u.Validate(); err != nil {
		return domainerrors.Err[*models.Upload](turso.NewError("uploads.Create", domainerrors.CodeInvalidInput, "invalid upload data", err))
	}

	created, err := client.Upload.Create().
		SetID(u.ID).
		SetFilename(u.Filename).
		SetMimeType(u.MimeType).
		SetData(u.Data).
		SetSizeBytes(u.SizeBytes).
		Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Upload](turso.MapEntError("uploads.Create", "upload", u.ID, err))
	}

	return domainerrors.Ok(&models.Upload{
		ID:        created.ID,
		Filename:  created.Filename,
		MimeType:  created.MimeType,
		Data:      created.Data,
		SizeBytes: created.SizeBytes,
		CreatedAt: created.CreatedAt,
	})
}
