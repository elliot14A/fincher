package uploads_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/uploads"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func setupTestDB(t *testing.T) *ent.Client {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to run schema automigrations: %v", err)
	}

	return client
}

func TestUploads_CRUD(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	rawBytes := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG magic header

	// 1. Create Upload
	u := &models.Upload{
		ID:        "upload-test-1",
		Filename:  "avatar.png",
		MimeType:  "image/png",
		Data:      rawBytes,
		SizeBytes: int64(len(rawBytes)),
		CreatedAt: time.Now().UTC(),
	}

	res := uploads.Create(ctx, client, u)
	if res.IsErr() {
		t.Fatalf("failed to create upload: %v", res.Error())
	}
	created := res.Unwrap()
	if created.ID != "upload-test-1" || len(created.Data) != len(rawBytes) {
		t.Fatalf("unexpected created upload: %+v", created)
	}

	// 2. Get Upload
	getRes := uploads.Get(ctx, client, "upload-test-1")
	if getRes.IsErr() {
		t.Fatalf("failed to get upload: %v", getRes.Error())
	}
	got := getRes.Unwrap()
	if got.Filename != "avatar.png" || got.MimeType != "image/png" {
		t.Fatalf("unexpected upload metadata: %+v", got)
	}

	// 3. Delete Upload
	delRes := uploads.Delete(ctx, client, "upload-test-1")
	if delRes.IsErr() {
		t.Fatalf("failed to delete upload: %v", delRes.Error())
	}

	// 4. Verify Not Found
	notFoundRes := uploads.Get(ctx, client, "upload-test-1")
	if notFoundRes.IsOk() {
		t.Fatalf("expected upload to be deleted, got: %+v", notFoundRes.Unwrap())
	}
	domErr, ok := notFoundRes.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeNotFound {
		t.Fatalf("expected CodeNotFound, got: %v", notFoundRes.Error())
	}
}

func TestUploads_SizeValidation(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	oversizedData := make([]byte, models.MaxUploadSizeBytes+1)
	u := &models.Upload{
		ID:        "upload-oversized",
		Filename:  "huge.png",
		MimeType:  "image/png",
		Data:      oversizedData,
		SizeBytes: int64(len(oversizedData)),
		CreatedAt: time.Now().UTC(),
	}

	res := uploads.Create(ctx, client, u)
	if res.IsOk() {
		t.Fatalf("expected oversized upload to fail validation")
	}
	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeInvalidInput {
		t.Fatalf("expected CodeInvalidInput, got: %v", res.Error())
	}
}
