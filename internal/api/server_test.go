package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elliot14A/fincher/internal/api"
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
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

func TestServer_HealthAndOpenAPI(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 for /health, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("expected health response to contain status:ok, got %s", rec.Body.String())
	}

	reqOpenAPI := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	recOpenAPI := httptest.NewRecorder()
	e.ServeHTTP(recOpenAPI, reqOpenAPI)
	if recOpenAPI.Code != http.StatusOK {
		t.Errorf("expected status 200 for /openapi.json, got %d", recOpenAPI.Code)
	}
	if !strings.Contains(recOpenAPI.Body.String(), `"swagger"`) || !strings.Contains(recOpenAPI.Body.String(), `Fincher Media Delivery Operations API`) {
		t.Errorf("expected valid swagger JSON, got %s", recOpenAPI.Body.String())
	}

	reqAPIHealth := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	recAPIHealth := httptest.NewRecorder()
	e.ServeHTTP(recAPIHealth, reqAPIHealth)
	if recAPIHealth.Code != http.StatusOK {
		t.Errorf("expected status 200 for /api/health, got %d", recAPIHealth.Code)
	}
}

func TestServer_EmbeddedWebMPA(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	reqRoot := httptest.NewRequest(http.MethodGet, "/", nil)
	recRoot := httptest.NewRecorder()
	e.ServeHTTP(recRoot, reqRoot)
	if recRoot.Code != http.StatusOK {
		t.Errorf("expected status 200 for /, got %d", recRoot.Code)
	}
	if !strings.Contains(recRoot.Body.String(), `<!doctype html>`) && !strings.Contains(recRoot.Body.String(), `id="app"`) {
		t.Errorf("expected index.html content for /, got: %s", recRoot.Body.String())
	}

	reqTitles := httptest.NewRequest(http.MethodGet, "/titles", nil)
	recTitles := httptest.NewRecorder()
	e.ServeHTTP(recTitles, reqTitles)
	if recTitles.Code != http.StatusOK {
		t.Errorf("expected status 200 for /titles, got %d", recTitles.Code)
	}
	if !strings.Contains(recTitles.Body.String(), `<!doctype html>`) && !strings.Contains(recTitles.Body.String(), `id="app"`) {
		t.Errorf("expected HTML content for /titles, got: %s", recTitles.Body.String())
	}

	reqPostTitles := httptest.NewRequest(http.MethodPost, "/titles", nil)
	recPostTitles := httptest.NewRecorder()
	e.ServeHTTP(recPostTitles, reqPostTitles)
	if recPostTitles.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for POST /titles, got %d", recPostTitles.Code)
	}
	if strings.Contains(recPostTitles.Body.String(), `<!doctype html>`) {
		t.Errorf("expected non-GET request to not receive HTML shell body")
	}
}
