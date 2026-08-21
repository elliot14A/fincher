package vendors_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestVendors_HTTP_Lifecycle(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	vendor := models.Vendor{
		Base: models.Base{
			ID: "vendor_a",
			Metadata: map[string]any{
				"tier": "P0",
			},
		},
		Name:      "Vendor A",
		Specialty: "AUDIO_DUBBING",
	}

	// 1. POST /api/vendors
	body, _ := json.Marshal(vendor)
	req := httptest.NewRequest(http.MethodPost, "/api/vendors", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	// 2. GET /api/vendors/:id
	req = httptest.NewRequest(http.MethodGet, "/api/vendors/vendor_a", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. GET /api/vendors?specialty=AUDIO_DUBBING
	req = httptest.NewRequest(http.MethodGet, "/api/vendors?specialty=AUDIO_DUBBING", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 4. PATCH /api/vendors/:id
	newName := "Vendor A International"
	patchReq := models.UpdateVendorInput{Name: &newName}
	body, _ = json.Marshal(patchReq)
	req = httptest.NewRequest(http.MethodPatch, "/api/vendors/vendor_a", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on patch, got %d", rec.Code)
	}

	// 5. DELETE /api/vendors/:id
	req = httptest.NewRequest(http.MethodDelete, "/api/vendors/vendor_a", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
