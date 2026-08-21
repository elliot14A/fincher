package masters_test

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

	server := api.NewServer(client)
	return server, client
}

func TestMasters_HTTP_Lifecycle(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	master := models.Master{
		ID:                "master-eclipse-v13",
		TitleID:           "title-eclipse",
		Version:           "V13",
		SupersedesVersion: "V12",
	}

	// 1. POST /masters
	body, _ := json.Marshal(master)
	req := httptest.NewRequest(http.MethodPost, "/masters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	// 2. GET /masters/:id
	req = httptest.NewRequest(http.MethodGet, "/masters/master-eclipse-v13", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. GET /masters?title_id=title-eclipse
	req = httptest.NewRequest(http.MethodGet, "/masters?title_id=title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 4. DELETE /masters/:id
	req = httptest.NewRequest(http.MethodDelete, "/masters/master-eclipse-v13", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
