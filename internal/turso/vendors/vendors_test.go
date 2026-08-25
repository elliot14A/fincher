package vendors_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/uploads"
	"github.com/elliot14A/fincher/internal/turso/vendors"
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
		t.Fatalf("failed to automigrate: %v", err)
	}

	return client
}

func TestVendors_CRUD(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create Vendor with Metadata
	v1 := &models.Vendor{
		Base: models.Base{
			ID: "vendor_a",
			Metadata: map[string]any{
				"contact_email": "ops@vendor-a.com",
				"sla_tier":      "tier_1",
			},
		},
		Name:      "Vendor A",
		Specialty: "AUDIO_DUBBING",
	}

	createRes := vendors.Create(ctx, client, v1)
	if createRes.IsErr() {
		t.Fatalf("failed to create vendor: %v", createRes.Error())
	}
	created := createRes.Unwrap()
	if created.Metadata["sla_tier"] != "tier_1" {
		t.Errorf("expected sla_tier in metadata, got: %v", created.Metadata["sla_tier"])
	}

	// 2. Get Vendor
	getRes := vendors.Get(ctx, client, "vendor_a")
	if getRes.IsErr() {
		t.Fatalf("failed to get vendor: %v", getRes.Error())
	}

	// 3. List Vendors with Pagination
	p := models.NewPagination(1, 10, "asc", "")
	listRes := vendors.List(ctx, client, domainerrors.Some("AUDIO_DUBBING"), p)
	if listRes.IsErr() {
		t.Fatalf("failed to list vendors: %v", listRes.Error())
	}
	res := listRes.Unwrap()
	if len(res.Items) != 1 || res.TotalItems != 1 {
		t.Fatalf("expected 1 vendor, got %d (total: %d)", len(res.Items), res.TotalItems)
	}

	// 4. Update Vendor
	newName := "Vendor A International"
	upRes := vendors.Update(ctx, client, "vendor_a", &models.UpdateVendorInput{
		Name: &newName,
		Metadata: map[string]any{
			"contact_email": "ops@vendor-a.com",
			"sla_tier":      "tier_1_premium",
		},
	})
	if upRes.IsErr() {
		t.Fatalf("failed to update vendor: %v", upRes.Error())
	}
	if upRes.Unwrap().Name != newName {
		t.Errorf("expected updated name %s, got %s", newName, upRes.Unwrap().Name)
	}

	// 5. Delete Vendor
	delRes := vendors.Delete(ctx, client, "vendor_a")
	if delRes.IsErr() {
		t.Fatalf("failed to delete vendor: %v", delRes.Error())
	}
}

func TestVendors_FK_DeleteBlockedByDependents(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	_ = titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-eclipse"},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(24 * time.Hour),
		Territories:          10,
		CurrentMasterVersion: "V1",
		OverallStatus:        models.StatusOnTrack,
	})

	_ = vendors.Create(ctx, client, &models.Vendor{
		Base:      models.Base{ID: "vendor_locked"},
		Name:      "Vendor Locked",
		Specialty: "SUBTITLES",
	})

	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-locked"},
		TitleID:                  "title-eclipse",
		Component:                models.ComponentSubtitle,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_locked",
		DerivedFromMasterVersion: "V1",
		Status:                   models.PackageStatusValid,
	})

	delRes := vendors.Delete(ctx, client, "vendor_locked")
	if delRes.IsOk() {
		t.Fatalf("expected vendor delete to fail due to dependent packages")
	}

	domErr, ok := delRes.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeConflict {
		t.Fatalf("expected CodeConflict (409) for delete blocked by dependent packages, got: %v", delRes.Error())
	}
}

func TestVendors_CascadeUploadDelete(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create upload blob
	upRes := uploads.Create(ctx, client, &models.Upload{
		ID:        "upload-logo-456",
		Filename:  "logo.png",
		MimeType:  "image/png",
		SizeBytes: 4,
		Data:      []byte{0x89, 0x50, 0x4E, 0x47},
	})
	if upRes.IsErr() {
		t.Fatalf("failed to create upload: %v", upRes.Error())
	}

	// 2. Create vendor referencing this upload
	vRes := vendors.Create(ctx, client, &models.Vendor{
		Base: models.Base{
			ID: "vendor-with-avatar",
			Metadata: map[string]any{
				"avatar_url": "/api/uploads/upload-logo-456",
			},
		},
		Name:      "Logo Studio",
		Specialty: "AUDIO_DUBBING",
	})
	if vRes.IsErr() {
		t.Fatalf("failed to create vendor: %v", vRes.Error())
	}

	// 3. Delete vendor
	delRes := vendors.Delete(ctx, client, "vendor-with-avatar")
	if delRes.IsErr() {
		t.Fatalf("failed to delete vendor: %v", delRes.Error())
	}

	// 4. Verify that upload was cascaded and deleted
	getUpRes := uploads.Get(ctx, client, "upload-logo-456")
	if getUpRes.IsOk() {
		t.Fatalf("expected upload blob to be deleted along with vendor, but it still exists")
	}
}
