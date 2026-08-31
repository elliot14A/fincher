package packages_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/masters"
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

	// 1. Seed Title (at V12 initially)
	titleRes := titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-eclipse"},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V12",
		OverallStatus:        models.StatusAtRisk,
	})
	if titleRes.IsErr() {
		t.Fatalf("failed to seed title: %v", titleRes.Error())
	}

	// 2. Seed Vendor
	vendorRes := vendors.Create(ctx, client, &models.Vendor{
		Base:       models.Base{ID: "vendor_a"},
		Name:       "Vendor A",
		Components: []string{"AUDIO"},
		Markets:    []string{"en-US"},
	})
	if vendorRes.IsErr() {
		t.Fatalf("failed to seed vendor: %v", vendorRes.Error())
	}

	return client
}

func TestPackages_CRUD_And_Staleness(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create Package derived from V12
	pkg := &models.Package{
		Base: models.Base{
			ID: "pkg-eclipse-es-audio",
			Metadata: map[string]any{
				"sample_rate": "48kHz",
				"channels":    "5.1",
			},
		},
		TitleID:                  "title-eclipse",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V12",
		RedeliveryCount:          1,
		Status:                   models.PackageStatusValid,
	}

	createRes := packages.Create(ctx, client, pkg)
	if createRes.IsErr() {
		t.Fatalf("failed to create package: %v", createRes.Error())
	}
	created := createRes.Unwrap()

	// 2. Initially title is V12 -> Package derived from V12 is NOT stale
	title := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if created.IsStaleAgainst(title.CurrentMasterVersion) {
		t.Fatalf("expected package derived from V12 to not be stale against title V12")
	}

	// 3. New Master cut V13 is cut -> Title's CurrentMasterVersion syncs to V13
	m13Res := masters.Create(ctx, client, &models.Master{
		ID:                "master-eclipse-v13",
		TitleID:           "title-eclipse",
		Version:           "V13",
		SupersedesVersion: "V12",
	})
	if m13Res.IsErr() {
		t.Fatalf("failed to create master V13: %v", m13Res.Error())
	}

	// Re-fetch title and verify package is now stale
	title = titles.Get(ctx, client, "title-eclipse").Unwrap()
	if !created.IsStaleAgainst(title.CurrentMasterVersion) {
		t.Fatalf("expected package derived from V12 to be stale against title newly synced to V13")
	}

	// 4. Update Package to V13 cut and RE_QC_PENDING
	newMasterVer := "V13"
	newStatus := models.PackageStatusReQCPending
	upRes := packages.Update(ctx, client, "pkg-eclipse-es-audio", &models.UpdatePackageInput{
		DerivedFromMasterVersion: &newMasterVer,
		Status:                   &newStatus,
		Metadata: map[string]any{
			"sample_rate": "48kHz",
			"channels":    "7.1",
		},
	})
	if upRes.IsErr() {
		t.Fatalf("failed to update package: %v", upRes.Error())
	}
	updated := upRes.Unwrap()
	if updated.DerivedFromMasterVersion != "V13" || updated.Status != models.PackageStatusReQCPending {
		t.Fatalf("unexpected updated package: %+v", updated)
	}

	// 5. Get Package
	getRes := packages.Get(ctx, client, "pkg-eclipse-es-audio")
	if getRes.IsErr() {
		t.Fatalf("failed to get package: %v", getRes.Error())
	}

	// 6. List Packages by Title and Component with Pagination
	comp := models.ComponentAudio
	p := models.NewPagination(1, 10, "asc", "")
	listRes := packages.List(ctx, client, packages.ListFilter{
		TitleID:   domainerrors.Some("title-eclipse"),
		Component: domainerrors.Some(comp),
		VendorID:  domainerrors.None[string](),
		Status:    domainerrors.None[models.PackageStatus](),
	}, p)
	if listRes.IsErr() {
		t.Fatalf("failed to list packages: %v", listRes.Error())
	}
	res := listRes.Unwrap()
	if len(res.Items) != 1 || res.TotalItems != 1 {
		t.Fatalf("expected 1 package, got %d (total: %d)", len(res.Items), res.TotalItems)
	}

	// 7. Delete Package
	delRes := packages.Delete(ctx, client, "pkg-eclipse-es-audio")
	if delRes.IsErr() {
		t.Fatalf("failed to delete package: %v", delRes.Error())
	}
}

func TestPackages_FK_Constraint(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// Orphan Title ID
	orphanPkg1 := &models.Package{
		Base:                     models.Base{ID: "pkg-orphan-title"},
		TitleID:                  "non-existent-title",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V12",
		Status:                   models.PackageStatusValid,
	}
	res1 := packages.Create(ctx, client, orphanPkg1)
	if res1.IsOk() {
		t.Fatalf("expected orphan title FK error, but creation succeeded")
	}

	// Orphan Vendor ID
	orphanPkg2 := &models.Package{
		Base:                     models.Base{ID: "pkg-orphan-vendor"},
		TitleID:                  "title-eclipse",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "non-existent-vendor",
		DerivedFromMasterVersion: "V12",
		Status:                   models.PackageStatusValid,
	}
	res2 := packages.Create(ctx, client, orphanPkg2)
	if res2.IsOk() {
		t.Fatalf("expected orphan vendor FK error, but creation succeeded")
	}
}
