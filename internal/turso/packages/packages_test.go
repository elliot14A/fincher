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

	// Seed Title
	titleRes := titles.Create(ctx, client, &models.Title{
		ID:                   "title-eclipse",
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})
	if titleRes.IsErr() {
		t.Fatalf("failed to seed title: %v", titleRes.Error())
	}

	// Seed Vendor
	vRes := vendors.Create(ctx, client, &models.Vendor{
		ID:        "vendor_a",
		Name:      "Vendor A Dubs",
		Specialty: "AUDIO_DUBBING",
	})
	if vRes.IsErr() {
		t.Fatalf("failed to seed vendor: %v", vRes.Error())
	}

	return client
}

func TestPackages_CRUD_And_Staleness(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	pkg1 := &models.Package{
		ID:                       "pkg-eclipse-es-audio",
		TitleID:                  "title-eclipse",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V12", // Stale against Title's V13!
		Status:                   models.PackageStatusValid,
	}

	// 1. Create
	createRes := packages.Create(ctx, client, pkg1)
	if createRes.IsErr() {
		t.Fatalf("failed to create package: %v", createRes.Error())
	}

	// 2. Get & Verify Staleness check
	getRes := packages.Get(ctx, client, "pkg-eclipse-es-audio")
	if getRes.IsErr() {
		t.Fatalf("failed to get package: %v", getRes.Error())
	}
	pkg := getRes.Unwrap()
	title := titles.Get(ctx, client, pkg.TitleID).Unwrap()

	// Assert staleness against current master version
	if !pkg.IsStaleAgainst(title.CurrentMasterVersion) {
		t.Errorf("expected package derived from V12 to be stale against title master %s", title.CurrentMasterVersion)
	}

	// Update to V13 cut and verify staleness resolves
	v13Master := "V13"
	upRes := packages.Update(ctx, client, pkg.ID, &models.UpdatePackageInput{
		DerivedFromMasterVersion: &v13Master,
	})
	if upRes.IsErr() {
		t.Fatalf("failed to update package master version: %v", upRes.Error())
	}
	if upRes.Unwrap().IsStaleAgainst(title.CurrentMasterVersion) {
		t.Errorf("package matching title master version should not be stale")
	}

	// Now deliver a new master V14 to the title
	_ = masters.Create(ctx, client, &models.Master{
		ID:                "master-eclipse-v14",
		TitleID:           "title-eclipse",
		Version:           "V14",
		SupersedesVersion: "V13",
	})
	updatedTitle := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if !upRes.Unwrap().IsStaleAgainst(updatedTitle.CurrentMasterVersion) {
		t.Errorf("package should be stale against newly delivered V14 master cut")
	}

	// 3. List
	listRes := packages.List(ctx, client, packages.ListFilter{
		TitleID:   domainerrors.Some("title-eclipse"),
		VendorID:  domainerrors.Some("vendor_a"),
		Component: domainerrors.Some(models.ComponentAudio),
		Status:    domainerrors.None[models.PackageStatus](),
	})
	if listRes.IsErr() {
		t.Fatalf("failed to list packages: %v", listRes.Error())
	}
	if len(listRes.Unwrap()) != 1 {
		t.Errorf("expected 1 package, got %d", len(listRes.Unwrap()))
	}

	// 4. Invalidate
	invalidStatus := models.PackageStatusInvalidated
	invRes := packages.Update(ctx, client, "pkg-eclipse-es-audio", &models.UpdatePackageInput{
		Status: &invalidStatus,
	})
	if invRes.IsErr() {
		t.Fatalf("failed to update package: %v", invRes.Error())
	}
	if invRes.Unwrap().Status != models.PackageStatusInvalidated {
		t.Errorf("expected status INVALIDATED, got %s", invRes.Unwrap().Status)
	}

	// 5. Delete
	delRes := packages.Delete(ctx, client, "pkg-eclipse-es-audio")
	if delRes.IsErr() {
		t.Fatalf("failed to delete package: %v", delRes.Error())
	}
}

func TestPackages_FK_Constraint(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// Negative Test 1: Invalid Title ID
	orphanTitlePkg := &models.Package{
		ID:                       "pkg-orphan-title",
		TitleID:                  "non-existent-title",
		Component:                models.ComponentVideo,
		Language:                 "en",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusPending,
	}
	res1 := packages.Create(ctx, client, orphanTitlePkg)
	if res1.IsOk() {
		t.Fatalf("expected package with invalid title_id to fail")
	}
	if domErr := res1.Error().(*domainerrors.DomainError); domErr.Code != domainerrors.CodeInvalidInput {
		t.Errorf("expected CodeInvalidInput, got %v", domErr.Code)
	}

	// Negative Test 2: Invalid Vendor ID
	orphanVendorPkg := &models.Package{
		ID:                       "pkg-orphan-vendor",
		TitleID:                  "title-eclipse",
		Component:                models.ComponentVideo,
		Language:                 "en",
		Version:                  "v1",
		VendorID:                 "non-existent-vendor",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusPending,
	}
	res2 := packages.Create(ctx, client, orphanVendorPkg)
	if res2.IsOk() {
		t.Fatalf("expected package with invalid vendor_id to fail")
	}
	if domErr := res2.Error().(*domainerrors.DomainError); domErr.Code != domainerrors.CodeInvalidInput {
		t.Errorf("expected CodeInvalidInput, got %v", domErr.Code)
	}
}
