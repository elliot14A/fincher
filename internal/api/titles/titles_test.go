package titles_test

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

	server := api.NewServer(client)
	return server, client
}

func TestTitles_HTTP_Lifecycle(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	premiere := now.Add(48 * time.Hour)

	title := models.Title{
		ID:                   "title-eclipse",
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         premiere,
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	}

	// 1. POST /titles
	body, _ := json.Marshal(title)
	req := httptest.NewRequest(http.MethodPost, "/titles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	// 2. GET /titles/:id
	req = httptest.NewRequest(http.MethodGet, "/titles/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var getRes models.Title
	_ = json.Unmarshal(rec.Body.Bytes(), &getRes)
	if getRes.Name != "Eclipse" || getRes.Territories != 40 {
		t.Errorf("unexpected fetched title: %+v", getRes)
	}

	// 3. GET /titles (List)
	req = httptest.NewRequest(http.MethodGet, "/titles", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var listRes []models.Title
	_ = json.Unmarshal(rec.Body.Bytes(), &listRes)
	if len(listRes) != 1 {
		t.Errorf("expected 1 title in list, got %d", len(listRes))
	}

	// 4. PATCH /titles/:id (Partial Update)
	patchStatus := models.StatusHold
	patchReq := models.UpdateTitleInput{
		OverallStatus: &patchStatus,
	}
	body, _ = json.Marshal(patchReq)
	req = httptest.NewRequest(http.MethodPatch, "/titles/title-eclipse", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on partial update, got %d. body: %s", rec.Code, rec.Body.String())
	}

	var patchRes models.Title
	_ = json.Unmarshal(rec.Body.Bytes(), &patchRes)
	if patchRes.OverallStatus != models.StatusHold || patchRes.Name != "Eclipse" || patchRes.Territories != 40 {
		t.Errorf("partial update did not retain original fields: %+v", patchRes)
	}

	// 5. DELETE /titles/:id
	req = httptest.NewRequest(http.MethodDelete, "/titles/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 on delete, got %d", rec.Code)
	}

	// 6. GET /titles/:id after delete (404)
	req = httptest.NewRequest(http.MethodGet, "/titles/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", rec.Code)
	}
}
