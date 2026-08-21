package deliveries_test

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

func TestDeliveries_HTTP_Lifecycle(t *testing.T) {
	server, client := setupTestServer(t)
	defer client.Close()

	e := server.Router()

	del := models.Delivery{
		ID:         "del-eclipse-es",
		TitleID:    "title-eclipse",
		Country:    "ES",
		Status:     models.DeliveryStatusPending,
		TargetDate: time.Now().Add(24 * time.Hour),
	}

	// 1. POST /deliveries
	body, _ := json.Marshal(del)
	req := httptest.NewRequest(http.MethodPost, "/deliveries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", rec.Code, rec.Body.String())
	}

	// 2. GET /deliveries/:id
	req = httptest.NewRequest(http.MethodGet, "/deliveries/del-eclipse-es", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. GET /deliveries?title_id=title-eclipse&country=ES
	req = httptest.NewRequest(http.MethodGet, "/deliveries?title_id=title-eclipse&country=ES", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var listRes []models.Delivery
	_ = json.Unmarshal(rec.Body.Bytes(), &listRes)
	if len(listRes) != 1 {
		t.Errorf("expected 1 delivery, got %d", len(listRes))
	}

	// 4. PATCH /deliveries/:id
	newStatus := models.DeliveryStatusReadyToShip
	patchReq := models.UpdateDeliveryInput{Status: &newStatus}
	body, _ = json.Marshal(patchReq)
	req = httptest.NewRequest(http.MethodPatch, "/deliveries/del-eclipse-es", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on patch, got %d", rec.Code)
	}

	// 5. DELETE /deliveries/:id
	req = httptest.NewRequest(http.MethodDelete, "/deliveries/del-eclipse-es", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
