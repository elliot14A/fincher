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

	now := time.Now().UTC().Truncate(time.Second)
	premiere := now.Add(48 * time.Hour)

	title := models.Title{
		Base: models.Base{
			ID: "title-eclipse",
			Metadata: map[string]any{
				"poster_url": "https://cdn.lume.stream/posters/eclipse.jpg",
			},
		},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         premiere,
		Territories:          40,
		CurrentMasterVersion: "V12",
		OverallStatus:        models.StatusAtRisk,
	}

	// 1. POST /api/titles
	body, _ := json.Marshal(title)
	req := httptest.NewRequest(http.MethodPost, "/api/titles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	var created models.Title
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if created.ID != "title-eclipse" || created.Name != "Eclipse" {
		t.Errorf("unexpected created title data: %+v", created)
	}

	// 2. GET /api/titles/:id
	req = httptest.NewRequest(http.MethodGet, "/api/titles/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var fetched models.Title
	if err := json.Unmarshal(rec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("failed to parse get response: %v", err)
	}
	if fetched.Name != "Eclipse" || fetched.Territories != 40 || fetched.CurrentMasterVersion != "V12" {
		t.Errorf("unexpected fetched title fields: %+v", fetched)
	}

	// 3. GET /api/titles (paginated list)
	req = httptest.NewRequest(http.MethodGet, "/api/titles?page=1&limit=10", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on list, got %d", rec.Code)
	}

	var result models.PaginationResult[models.Title]
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal pagination result: %v", err)
	}
	if len(result.Items) != 1 || result.TotalItems != 1 {
		t.Errorf("expected 1 title in list, got %d (total: %d)", len(result.Items), result.TotalItems)
	}

	// 4. PATCH /api/titles/:id
	newStatus := models.StatusHold
	patchReq := models.UpdateTitleInput{OverallStatus: &newStatus}
	body, _ = json.Marshal(patchReq)
	req = httptest.NewRequest(http.MethodPatch, "/api/titles/title-eclipse", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on patch, got %d", rec.Code)
	}
	var patched models.Title
	if err := json.Unmarshal(rec.Body.Bytes(), &patched); err != nil {
		t.Fatalf("failed to parse patch response: %v", err)
	}
	if patched.OverallStatus != models.StatusHold {
		t.Errorf("expected patched status HOLD, got: %s", patched.OverallStatus)
	}

	// 5. DELETE /api/titles/:id
	req = httptest.NewRequest(http.MethodDelete, "/api/titles/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 on delete, got %d", rec.Code)
	}

	// 6. GET /api/titles/:id -> 404
	req = httptest.NewRequest(http.MethodGet, "/api/titles/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", rec.Code)
	}
}
