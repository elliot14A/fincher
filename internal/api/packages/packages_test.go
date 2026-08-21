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
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
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

	// 1. Seed Title
	tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-eclipse"},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})

	// 2. Seed Vendor
	tursovendors.Create(ctx, client, &models.Vendor{
		Base:      models.Base{ID: "vendor_a"},
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
		Base: models.Base{
			ID: "pkg-eclipse-es-audio",
			Metadata: map[string]any{
				"codec": "AAC",
			},
		},
		TitleID:                  "title-eclipse",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V13",
		RedeliveryCount:          0,
		Status:                   models.PackageStatusPending,
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

	// 2. GET /packages/:id - assert response fields
	req = httptest.NewRequest(http.MethodGet, "/packages/pkg-eclipse-es-audio", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var fetched models.Package
	if err := json.Unmarshal(rec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("failed to parse get response: %v", err)
	}
	if fetched.Component != models.ComponentAudio || fetched.Language != "es" || fetched.VendorID != "vendor_a" {
		t.Errorf("unexpected fetched package fields: %+v", fetched)
	}

	// 3. GET /packages?title_id=title-eclipse&component=AUDIO - assert list length
	req = httptest.NewRequest(http.MethodGet, "/packages?title_id=title-eclipse&component=AUDIO", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var list []models.Package
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Errorf("expected 1 package in list, got %d", len(list))
	}

	// 4. PATCH /packages/:id - verify partial update preserves other fields
	newStatus := models.PackageStatusValid
	patchReq := models.UpdatePackageInput{Status: &newStatus}
	body, _ = json.Marshal(patchReq)
	req = httptest.NewRequest(http.MethodPatch, "/packages/pkg-eclipse-es-audio", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on patch, got %d", rec.Code)
	}
	var patched models.Package
	if err := json.Unmarshal(rec.Body.Bytes(), &patched); err != nil {
		t.Fatalf("failed to parse patch response: %v", err)
	}
	if patched.Status != models.PackageStatusValid {
		t.Errorf("expected patched status VALID, got: %s", patched.Status)
	}
	if patched.Component != models.ComponentAudio || patched.Language != "es" || patched.VendorID != "vendor_a" {
		t.Errorf("expected PATCH to preserve untouched fields (Component, Language, VendorID): %+v", patched)
	}

	// 5. DELETE /packages/:id
	req = httptest.NewRequest(http.MethodDelete, "/packages/pkg-eclipse-es-audio", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
