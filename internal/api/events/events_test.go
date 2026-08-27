package events_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/elliot14A/fincher/internal/api/events"
	"github.com/elliot14A/fincher/internal/clickhouse"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestEvents_BatchIngestion(t *testing.T) {
	conn, err := clickhouse.Open("127.0.0.1:9000")
	if err != nil {
		t.Skip("skipping integration: clickhouse connection failed:", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := clickhouse.AutoMigrate(ctx, conn); err != nil {
		t.Fatalf("running migrations: %v", err)
	}

	e := echo.New()
	g := e.Group("/api/events")
	events.RegisterRoutes(g, conn)

	titleSlug := "batch-title-" + uuid.NewString()[:8]

	// 1. Successful batch ingestion (array of 3 events)
	batch := []models.Event{
		{
			ID:      uuid.NewString(),
			Type:    models.TypeVendorHeartbeat,
			Source:  "vendor.agent",
			Subject: "GLOBAL",
			Time:    time.Now().UTC(),
		},
		{
			ID:      uuid.NewString(),
			Type:    models.TypeQCInspectionCompleted,
			Source:  "qc.agent.audio",
			Subject: titleSlug,
			Time:    time.Now().UTC(),
			Data: map[string]any{
				"status": "PASSED",
			},
		},
		{
			ID:      uuid.NewString(),
			Type:    models.TypeAudioSyncDriftDetected,
			Source:  "qc.agent.audio",
			Subject: titleSlug,
			Time:    time.Now().UTC(),
			Data: map[string]any{
				"drift_ms": 140.5,
			},
		},
	}

	body, err := json.Marshal(batch)
	if err != nil {
		t.Fatalf("marshaling batch: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got: %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("parsing response: %v", err)
	}

	if resp.Status != "ingested" {
		t.Errorf("expected status 'ingested', got: %s", resp.Status)
	}
	if resp.Count != 3 {
		t.Errorf("expected count 3, got: %d", resp.Count)
	}

	// 2. Empty array rejected
	reqEmpty := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader([]byte("[]")))
	reqEmpty.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recEmpty := httptest.NewRecorder()

	e.ServeHTTP(recEmpty, reqEmpty)

	if recEmpty.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on empty array, got: %d", recEmpty.Code)
	}

	// 3. Single event object (non-array) rejected
	singleEvent := models.Event{
		ID:      uuid.NewString(),
		Type:    models.TypeVendorHeartbeat,
		Source:  "vendor.agent",
		Subject: "GLOBAL",
		Time:    time.Now().UTC(),
	}
	singleBody, _ := json.Marshal(singleEvent)
	reqSingle := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(singleBody))
	reqSingle.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recSingle := httptest.NewRecorder()

	e.ServeHTTP(recSingle, reqSingle)

	if recSingle.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on single non-array object, got: %d", recSingle.Code)
	}

	// 4. Malformed JSON rejected
	reqBad := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader([]byte("not-json")))
	reqBad.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recBad := httptest.NewRecorder()

	e.ServeHTTP(recBad, reqBad)

	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on malformed body, got: %d", recBad.Code)
	}

	// 5. Partial failure rejected upfront (validate-all before writing any)
	invalidBatch := []models.Event{
		{
			ID:      uuid.NewString(),
			Type:    models.TypeVendorHeartbeat,
			Source:  "vendor.agent",
			Subject: "GLOBAL",
			Time:    time.Now().UTC(),
		},
		{
			ID:     uuid.NewString(),
			Type:   "", // Invalid: missing required type
			Source: "vendor.agent",
		},
	}
	invalidBody, _ := json.Marshal(invalidBatch)
	reqInvalid := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(invalidBody))
	reqInvalid.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recInvalid := httptest.NewRecorder()

	e.ServeHTTP(recInvalid, reqInvalid)

	if recInvalid.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on invalid event in batch, got: %d", recInvalid.Code)
	}

	// Verify first event was NOT written to ClickHouse
	var count int
	_ = conn.QueryRowContext(ctx, "select count() from fincher.events where id = ?", invalidBatch[0].ID).Scan(&count)
	if count != 0 {
		t.Fatalf("expected 0 events written when batch fails validation, got: %d", count)
	}

	// Cleanup test rows
	for _, ev := range batch {
		_, _ = conn.ExecContext(ctx, "alter table fincher.events delete where id = ?", ev.ID)
	}
}
