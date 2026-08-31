package events_test

import (
	"bytes"
	"context"
	"encoding/json"
	"iter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"
	genai "google.golang.org/genai"

	"github.com/elliot14A/fincher/internal/agent/scheduler"
	"github.com/elliot14A/fincher/internal/api/events"
	"github.com/elliot14A/fincher/internal/clickhouse"
	tursoruns "github.com/elliot14A/fincher/internal/turso/runs"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

type mockLLM struct {
	responses []string
	callIndex int
}

func (m *mockLLM) Name() string {
	return "mock-gemini"
}

func (m *mockLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		var text string
		if m.callIndex < len(m.responses) {
			text = m.responses[m.callIndex]
			m.callIndex++
		} else {
			text = `{"actionable":false,"anomaly_type":"BENIGN_TELEMETRY","severity":"INFO","rationale":"Default mock response"}`
		}

		resp := &model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{Text: text},
				},
			},
		}
		yield(resp, nil)
	}
}

func TestEvents_BatchIngestion_And_Routing(t *testing.T) {
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

	client := tursotest.NewMemoryClient(t)

	mock := &mockLLM{
		responses: []string{
			`{"actionable":false,"anomaly_type":"BENIGN_TELEMETRY","severity":"INFO","rationale":"Routine heartbeat ping"}`,
			`{"selected_vendor_id":"vendor-deluxe","confidence":0.95,"rationale":"Fast turnaround selected"}`,
		},
	}

	sched := scheduler.NewScheduler(time.Second)
	e := echo.New()
	g := e.Group("/api/events")
	events.RegisterRoutes(g, conn, client, func() model.LLM { return mock }, sched)

	titleSlug := "batch-title-" + uuid.NewString()[:8]
	anomalyEventID := uuid.NewString()
	allocEventID := uuid.NewString()
	qcEventID := uuid.NewString()

	// Seed vendor, title, and package so background workflows resolve cleanly
	_ = client.Vendor.Create().
		SetID("vendor-deluxe").
		SetName("Deluxe Media").
		SetSpecialty("AUDIO_DUBBING").
		SetHourlyRateUsd(150.0).
		SetTurnaroundHours(24).
		SaveX(context.Background())

	titleRecord := client.Title.Create().
		SetID("title-" + titleSlug).
		SetName("Batch Title Test").
		SetSlug(titleSlug).
		SetType("FEATURE").
		SetPremiereDate(time.Now().UTC().Add(48 * time.Hour)).
		SetTerritories(10).
		SetCurrentMasterVersion("V01").
		SetOverallStatus("HOLD").
		SaveX(context.Background())

	_ = client.MediaPackage.Create().
		SetID("pkg-test-1").
		SetTitleID(titleRecord.ID).
		SetComponent("AUDIO").
		SetLanguage("en-US").
		SetMarket("US").
		SetVersion("V01").
		SetVendorID("vendor-deluxe").
		SetDerivedFromMasterVersion("V01").
		SetStatus("RE_QC_PENDING").
		SaveX(context.Background())

	// 1. Successful batch ingestion with routing (4 events: telemetry, anomaly, allocation, clean QC)
	batch := []models.Event{
		{
			ID:      uuid.NewString(),
			Type:    models.TypeVendorHeartbeat,
			Source:  "vendor.agent",
			Subject: "GLOBAL",
			Time:    time.Now().UTC(),
		},
		{
			ID:       anomalyEventID,
			Type:     models.TypeAudioSyncDriftDetected,
			Source:   "qc.agent.audio",
			Subject:  titleSlug,
			Time:     time.Now().UTC(),
			Severity: models.SeverityCritical,
			Data: map[string]any{
				"drift_ms": 140.5,
			},
		},
		{
			ID:       allocEventID,
			Type:     models.TypePackageRequired,
			Source:   "delivery.planner",
			Subject:  titleSlug,
			Time:     time.Now().UTC(),
			Severity: models.SeverityInfo,
			Data: map[string]any{
				"component": "AUDIO",
			},
		},
		{
			ID:       qcEventID,
			Type:     models.TypeQCInspectionCompleted,
			Source:   "qc.agent.audio",
			Subject:  titleSlug,
			Time:     time.Now().UTC(),
			Severity: models.SeverityInfo,
			Data: map[string]any{
				"package_id": "pkg-test-1",
				"status":     "PASSED",
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

	var resp models.EventBatchResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("parsing response: %v", err)
	}

	if resp.Status != "ingested" {
		t.Errorf("expected status 'ingested', got: %s", resp.Status)
	}
	if resp.Count != 4 {
		t.Errorf("expected count 4, got: %d", resp.Count)
	}
	if len(resp.RunIDs) != 3 {
		t.Errorf("expected 3 triggered run IDs (1 anomaly, 1 allocation, 1 resolution), got: %d (%v)", len(resp.RunIDs), resp.RunIDs)
	}

	// Verify background runs were initialized in Turso
	time.Sleep(150 * time.Millisecond)

	anomalyRun := tursoruns.GetRun(context.Background(), client, "run-"+anomalyEventID)
	if anomalyRun.IsErr() {
		t.Errorf("expected anomaly run 'run-%s' to exist in Turso: %v", anomalyEventID, anomalyRun.Error())
	}

	allocRun := tursoruns.GetRun(context.Background(), client, "run-"+allocEventID)
	if allocRun.IsErr() {
		t.Errorf("expected alloc run 'run-%s' to exist in Turso: %v", allocEventID, allocRun.Error())
	}

	qcRun := tursoruns.GetRun(context.Background(), client, "run-"+qcEventID)
	if qcRun.IsErr() {
		t.Errorf("expected resolution run 'run-%s' to exist in Turso: %v", qcEventID, qcRun.Error())
	}

	// 1b. Idempotency test: Re-ingesting the identical batch returns 201 ingested with the existing run IDs preserved
	reqDup := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(body))
	reqDup.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recDup := httptest.NewRecorder()
	e.ServeHTTP(recDup, reqDup)
	if recDup.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on re-delivery: %d", recDup.Code)
	}
	var dupResp models.EventBatchResponse
	_ = json.Unmarshal(recDup.Body.Bytes(), &dupResp)
	if len(dupResp.RunIDs) != 3 {
		t.Errorf("expected 3 existing run IDs preserved on duplicate event batch, got %d (%v)", len(dupResp.RunIDs), dupResp.RunIDs)
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
