package packages_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/api"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
	tursotitles "github.com/elliot14A/fincher/pkg/turso/titles"
	tursovendors "github.com/elliot14A/fincher/pkg/turso/vendors"
)

func setupTestServer(t *testing.T) (*api.Server, *ent.Client) {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open memory ent client: %v", err)
	}

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to run ent automigration: %v", err)
	}

	// Seed Title
	tursotitles.Create(ctx, client, &models.Title{
		ID:                   "title-eclipse",
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})

	// Seed Vendor
	tursovendors.Create(ctx, client, &models.Vendor{
		ID:        "vendor_a",
		Name:      "Vendor A",
		Specialty: "AUDIO_DUBBING",
	})

	server := api.NewServer(client)
	return server, client
}

func TestPackages_HTTP_Lifecycle(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	pkg := models.Package{
		ID:                       "pkg-eclipse-es-audio",
		TitleID:                  "title-eclipse",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V12",
		Status:                   models.PackageStatusValid,
	}

	// 1. POST /packages
	body, _ := json.Marshal(pkg)
	req := httptest.NewRequest(http.MethodPost, "/packages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	// 2. GET /packages/:id
	req = httptest.NewRequest(http.MethodGet, "/packages/pkg-eclipse-es-audio", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var getRes models.Package
	_ = json.Unmarshal(rec.Body.Bytes(), &getRes)
	if getRes.Language != "es" || getRes.DerivedFromMasterVersion != "V12" {
		t.Errorf("unexpected package data: %+v", getRes)
	}

	// 3. GET /packages?title_id=title-eclipse&component=AUDIO
	req = httptest.NewRequest(http.MethodGet, "/packages?title_id=title-eclipse&component=AUDIO", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var listRes []models.Package
	_ = json.Unmarshal(rec.Body.Bytes(), &listRes)
	if len(listRes) != 1 {
		t.Errorf("expected 1 package, got %d", len(listRes))
	}

	// 4. PATCH /packages/:id (Invalidate)
	invStatus := models.PackageStatusInvalidated
	patchReq := models.UpdatePackageInput{Status: &invStatus}
	body, _ = json.Marshal(patchReq)
	req = httptest.NewRequest(http.MethodPatch, "/packages/pkg-eclipse-es-audio", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on patch, got %d", rec.Code)
	}

	// 5. DELETE /packages/:id
	req = httptest.NewRequest(http.MethodDelete, "/packages/pkg-eclipse-es-audio", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
