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

	tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-eclipse"},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})

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

	// 1. POST /api/packages
	body, _ := json.Marshal(pkg)
	req := httptest.NewRequest(http.MethodPost, "/api/packages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	// 2. GET /api/packages/:id
	req = httptest.NewRequest(http.MethodGet, "/api/packages/pkg-eclipse-es-audio", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. GET /api/packages?title_id=title-eclipse&component=AUDIO
	req = httptest.NewRequest(http.MethodGet, "/api/packages?title_id=title-eclipse&component=AUDIO", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var listRes models.PaginationResult[models.Package]
	if err := json.Unmarshal(rec.Body.Bytes(), &listRes); err != nil {
		t.Fatalf("failed to unmarshal pagination result: %v", err)
	}
	if len(listRes.Items) != 1 || listRes.TotalItems != 1 {
		t.Errorf("expected 1 package in list, got %d (total: %d)", len(listRes.Items), listRes.TotalItems)
	}

	// 4. PATCH /api/packages/:id
	newStatus := models.PackageStatusValid
	patchReq := models.UpdatePackageInput{Status: &newStatus}
	body, _ = json.Marshal(patchReq)
	req = httptest.NewRequest(http.MethodPatch, "/api/packages/pkg-eclipse-es-audio", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on patch, got %d", rec.Code)
	}

	// 5. DELETE /api/packages/:id
	req = httptest.NewRequest(http.MethodDelete, "/api/packages/pkg-eclipse-es-audio", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
