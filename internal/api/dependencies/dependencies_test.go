package dependencies_test

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
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
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
	_ = tursotitles.Create(ctx, client, &models.Title{
		ID:                   "title-eclipse",
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})

	// 2. Seed Vendor
	_ = tursovendors.Create(ctx, client, &models.Vendor{
		ID:        "vendor_a",
		Name:      "Vendor A",
		Specialty: "AUDIO_DUBBING",
	})

	// 3. Seed Packages
	_ = tursopackages.Create(ctx, client, &models.Package{
		ID:                       "pkg-video-ov",
		TitleID:                  "title-eclipse",
		Component:                models.ComponentVideo,
		Language:                 "ov",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V13",
		Status:                   models.PackageStatusValid,
	})

	_ = tursopackages.Create(ctx, client, &models.Package{
		ID:                       "pkg-audio-es",
		TitleID:                  "title-eclipse",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V13",
		Status:                   models.PackageStatusValid,
	})

	server := api.NewServer(client)
	return server, client
}

func TestDependencies_HTTP_Lifecycle(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	dep := models.Dependency{
		ID:             "dep-video-audio",
		ParentID:       "pkg-video-ov",
		ChildID:        "pkg-audio-es",
		DependencyType: models.DependencyMasterDerivation,
	}

	// 1. POST /dependencies
	body, _ := json.Marshal(dep)
	req := httptest.NewRequest(http.MethodPost, "/dependencies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	// 2. GET /dependencies?parent_id=pkg-video-ov
	req = httptest.NewRequest(http.MethodGet, "/dependencies?parent_id=pkg-video-ov", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. GET /dependencies/graph/title-eclipse
	req = httptest.NewRequest(http.MethodGet, "/dependencies/graph/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on graph, got %d", rec.Code)
	}

	var graphRes models.LineageGraph
	_ = json.Unmarshal(rec.Body.Bytes(), &graphRes)
	if len(graphRes.Roots) != 1 || graphRes.Roots[0].PackageID != "pkg-video-ov" {
		t.Errorf("unexpected lineage graph response: %+v", graphRes)
	}

	// 4. DELETE /dependencies/:id
	req = httptest.NewRequest(http.MethodDelete, "/dependencies/dep-video-audio", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 on delete, got %d", rec.Code)
	}
}
