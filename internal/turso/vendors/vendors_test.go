package vendors_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/titles"
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

	v1 := &models.Vendor{
		ID:        "vendor_a",
		Name:      "Vendor A Audio Labs",
		Specialty: "AUDIO_DUBBING",
	}
	v2 := &models.Vendor{
		ID:        "vendor_b",
		Name:      "Vendor B Subtitles",
		Specialty: "SUBTITLES",
	}

	// 1. Create
	create1 := vendors.Create(ctx, client, v1)
	if create1.IsErr() {
		t.Fatalf("failed to create vendor 1: %v", create1.Error())
	}
	create2 := vendors.Create(ctx, client, v2)
	if create2.IsErr() {
		t.Fatalf("failed to create vendor 2: %v", create2.Error())
	}

	// 2. Get
	getRes := vendors.Get(ctx, client, "vendor_a")
	if getRes.IsErr() {
		t.Fatalf("failed to get vendor: %v", getRes.Error())
	}
	if getRes.Unwrap().Name != "Vendor A Audio Labs" {
		t.Errorf("unexpected vendor name: %s", getRes.Unwrap().Name)
	}

	// 3. List with filter
	listRes := vendors.List(ctx, client, domainerrors.Some("AUDIO_DUBBING"))
	if listRes.IsErr() {
		t.Fatalf("failed to list vendors: %v", listRes.Error())
	}
	if len(listRes.Unwrap()) != 1 {
		t.Errorf("expected 1 audio vendor, got %d", len(listRes.Unwrap()))
	}

	// 4. Update
	newName := "Vendor A Global Dubbing House"
	upRes := vendors.Update(ctx, client, "vendor_a", &models.UpdateVendorInput{
		Name: &newName,
	})
	if upRes.IsErr() {
		t.Fatalf("failed to update vendor: %v", upRes.Error())
	}
	if upRes.Unwrap().Name != newName {
		t.Errorf("updated name not persisted: %s", upRes.Unwrap().Name)
	}

	// 5. Delete
	delRes := vendors.Delete(ctx, client, "vendor_b")
	if delRes.IsErr() {
		t.Fatalf("failed to delete vendor: %v", delRes.Error())
	}
}

func TestVendors_FK_DeleteBlockedByDependents(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Seed Title
	_ = titles.Create(ctx, client, &models.Title{
		ID:                   "title-eclipse",
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})

	// 2. Seed Vendor
	_ = vendors.Create(ctx, client, &models.Vendor{
		ID:        "vendor_locked",
		Name:      "Vendor Locked",
		Specialty: "SUBTITLES",
	})

	// 3. Seed Package dependent on vendor_locked
	_ = packages.Create(ctx, client, &models.Package{
		ID:                       "pkg-dep",
		TitleID:                  "title-eclipse",
		Component:                models.ComponentSubtitle,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_locked",
		DerivedFromMasterVersion: "V13",
		Status:                   models.PackageStatusValid,
	})

	// 4. Deleting vendor_locked must fail with CodeConflict (409) because dependent Package references it
	delRes := vendors.Delete(ctx, client, "vendor_locked")
	if delRes.IsOk() {
		t.Fatalf("expected vendor delete to be blocked by dependent package")
	}

	domErr, ok := delRes.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeConflict {
		t.Errorf("expected CodeConflict (409) for delete blocked by dependents, got: %v", delRes.Error())
	}
}
