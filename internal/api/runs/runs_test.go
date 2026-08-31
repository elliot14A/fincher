package runs_test

import (
	"bytes"
	"context"
	"encoding/json"
	"iter"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/adk/v2/model"
	genai "google.golang.org/genai"

	"github.com/elliot14A/fincher/internal/api"
	"github.com/elliot14A/fincher/internal/api/runs"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursoruns "github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func setupTestServer(t *testing.T) (*api.Server, *ent.Client) {
	client := tursotest.NewMemoryClient(t)
	server := api.NewServer(client)
	return server, client
}

func TestRuns_HTTP_Lifecycle(t *testing.T) {
	server, client := setupTestServer(t)
	ctx := context.Background()
	e := server.Router()

	// Seed a test run
	now := time.Now().UTC()
	seededRunRes := tursoruns.CreateRun(ctx, client, &models.Run{
		Base: models.Base{
			ID: "run-seeded-1",
		},
		TitleSlug: "eclipse",
		Trigger:   "incident",
		Status:    models.RunStatusRunning,
		StartedAt: now,
	})
	if seededRunRes.IsErr() {
		t.Fatalf("failed to seed run: %v", seededRunRes.Error())
	}

	// 1. GET /api/runs (List)
	t.Run("GET /api/runs returns paginated list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/runs", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got: %d (%s)", rec.Code, rec.Body.String())
		}

		var page models.RunPaginationResult
		if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
			t.Fatalf("failed to parse pagination result: %v", err)
		}

		if page.TotalItems < 1 {
			t.Errorf("expected at least 1 total item, got: %d", page.TotalItems)
		}
		if len(page.Items) < 1 {
			t.Errorf("expected at least 1 item in list, got: %d", len(page.Items))
		}
	})

	// 2. GET /api/runs with filter
	t.Run("GET /api/runs?wf=incident filters correctly", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/runs?wf=incident&status=RUNNING", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got: %d", rec.Code)
		}

		var page models.RunPaginationResult
		if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
			t.Fatalf("failed to parse: %v", err)
		}
		if page.TotalItems != 1 {
			t.Errorf("expected 1 filtered item, got: %d", page.TotalItems)
		}
	})

	// 3. GET /api/runs/:id (Found)
	t.Run("GET /api/runs/:id returns run details", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/runs/run-seeded-1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got: %d", rec.Code)
		}

		var r models.Run
		if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
			t.Fatalf("failed to parse run: %v", err)
		}
		if r.ID != "run-seeded-1" {
			t.Errorf("expected run ID run-seeded-1, got: %s", r.ID)
		}
		if r.TitleSlug != "eclipse" {
			t.Errorf("expected title_slug eclipse, got: %s", r.TitleSlug)
		}
	})

	// 4. GET /api/runs/:id (Not Found)
	t.Run("GET /api/runs/:id returns 404 for unknown ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/runs/run-non-existent", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got: %d", rec.Code)
		}
	})

	// 5. POST /api/runs returns 503 when model is nil, 201 when set
	t.Run("POST /api/runs guards model presence and creates run", func(t *testing.T) {
		body, _ := json.Marshal(runs.CreateRunRequest{
			TitleSlug:          "eclipse",
			HoursUntilPremiere: 48,
			Component:          "AUDIO",
		})

		// Before model is set: 503 Service Unavailable
		req503 := httptest.NewRequest(http.MethodPost, "/api/runs?wf=incident", bytes.NewReader(body))
		req503.Header.Set("Content-Type", "application/json")
		rec503 := httptest.NewRecorder()
		e.ServeHTTP(rec503, req503)
		if rec503.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected 503 Service Unavailable when model is nil, got: %d", rec503.Code)
		}

		// Inject mock model
		mock := &mockLLM{
			responses: []string{
				`{"actionable":false,"anomaly_type":"BENIGN_TELEMETRY","severity":"INFO","rationale":"Routine heartbeat ping"}`,
			},
		}
		server.SetModel(mock)

		// After model is set: 201 Created
		req201 := httptest.NewRequest(http.MethodPost, "/api/runs?wf=incident", bytes.NewReader(body))
		req201.Header.Set("Content-Type", "application/json")
		rec201 := httptest.NewRecorder()
		e.ServeHTTP(rec201, req201)

		if rec201.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got: %d (%s)", rec201.Code, rec201.Body.String())
		}

		var created models.Run
		if err := json.Unmarshal(rec201.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode created run: %v", err)
		}
		if created.Trigger != "incident" {
			t.Errorf("expected trigger incident, got: %s", created.Trigger)
		}
		if created.Status != models.RunStatusRunning {
			t.Errorf("expected status RUNNING, got: %s", created.Status)
		}
	})

	// 6. POST /api/runs?wf=allocation creates allocation run
	t.Run("POST /api/runs creates allocation run", func(t *testing.T) {
		body, _ := json.Marshal(runs.CreateRunRequest{
			TitleSlug:          "eclipse",
			HoursUntilPremiere: 72,
			Component:          "SUBTITLE",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/runs?wf=allocation", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got: %d (%s)", rec.Code, rec.Body.String())
		}

		var created models.Run
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode created run: %v", err)
		}
		if created.Trigger != "allocation" {
			t.Errorf("expected trigger allocation, got: %s", created.Trigger)
		}
	})

	// 7. Idempotent run creation for duplicate event IDs
	t.Run("POST /api/runs is idempotent for duplicate event IDs", func(t *testing.T) {
		body, _ := json.Marshal(runs.CreateRunRequest{
			TitleSlug: "eclipse",
			Event: &models.Event{
				ID:       "idempotent-event-1",
				Type:     "vendor.ping",
				Subject:  "eclipse",
				Time:     time.Now().UTC(),
				Severity: models.SeverityInfo,
			},
		})
		// First call: 201 Created
		req1 := httptest.NewRequest(http.MethodPost, "/api/runs?wf=incident", bytes.NewReader(body))
		req1.Header.Set("Content-Type", "application/json")
		rec1 := httptest.NewRecorder()
		e.ServeHTTP(rec1, req1)
		if rec1.Code != http.StatusCreated {
			t.Fatalf("expected 201 on first call, got: %d", rec1.Code)
		}

		// Second call with same event ID: 200 OK (idempotent)
		req2 := httptest.NewRequest(http.MethodPost, "/api/runs?wf=incident", bytes.NewReader(body))
		req2.Header.Set("Content-Type", "application/json")
		rec2 := httptest.NewRecorder()
		e.ServeHTTP(rec2, req2)
		if rec2.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on duplicate event call, got: %d", rec2.Code)
		}
	})

	// 8. GET /api/runs/:id/stream (SSE)
	t.Run("GET /api/runs/:id/stream streams events and completes", func(t *testing.T) {
		doneTime := time.Now().UTC()
		tursoruns.UpdateRunStatus(ctx, client, "run-seeded-1", models.RunStatusCompleted, &doneTime, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/runs/run-seeded-1/stream", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK for SSE stream, got: %d", rec.Code)
		}

		bodyStr := rec.Body.String()
		if !strings.Contains(bodyStr, "event: update") {
			t.Errorf("expected SSE body to contain 'event: update', got: %s", bodyStr)
		}
		if !strings.Contains(bodyStr, "event: done") {
			t.Errorf("expected SSE body to contain 'event: done', got: %s", bodyStr)
		}
	})

	// 9. End-to-end execution with SetModel and verify step persistence
	t.Run("POST /api/runs executes with injected model and persists triage step", func(t *testing.T) {
		mock := &mockLLM{
			responses: []string{
				`{"actionable":false,"anomaly_type":"BENIGN_TELEMETRY","severity":"INFO","rationale":"Routine heartbeat ping"}`,
			},
		}
		server.SetModel(mock)

		body, _ := json.Marshal(runs.CreateRunRequest{
			TitleSlug:          "eclipse",
			HoursUntilPremiere: 48,
			Event: &models.Event{
				ID:       "ping-event-stage-test",
				Type:     "vendor.ping",
				Subject:  "eclipse",
				Time:     time.Now().UTC(),
				Severity: models.SeverityInfo,
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/runs?wf=incident", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got: %d", rec.Code)
		}

		// Wait briefly for background goroutine to process benign event
		time.Sleep(100 * time.Millisecond)

		getRes := tursoruns.GetRun(ctx, client, "run-ping-event-stage-test")
		if getRes.IsErr() {
			t.Fatalf("failed to fetch run: %v", getRes.Error())
		}
		r := getRes.Unwrap()
		if r.Status != models.RunStatusCompleted {
			t.Errorf("expected status COMPLETED, got: %s", r.Status)
		}
		if len(r.Steps) != 1 {
			t.Fatalf("expected 1 triage step, got: %d", len(r.Steps))
		}
		if r.Steps[0].Name != "triage_judge" {
			t.Errorf("expected step name triage_judge, got: %s", r.Steps[0].Name)
		}
		if len(r.Results) != 1 {
			t.Fatalf("expected 1 triage WfResult, got: %d", len(r.Results))
		}
		if r.Results[0].Outcome != "FILTERED" {
			t.Errorf("expected outcome FILTERED, got: %s", r.Results[0].Outcome)
		}
	})

	// 10. POST /api/runs computes HoursUntilPremiere dynamically from title slug
	t.Run("POST /api/runs computes HoursUntilPremiere from Title slug", func(t *testing.T) {
		premiere := time.Now().UTC().Add(10 * time.Hour)
		tursotitles.Create(ctx, client, &models.Title{
			Base: models.Base{
				ID: "title-avatar",
			},
			Name:                 "Avatar Fire and Ash",
			Slug:                 "avatar-fire-ash",
			Type:                 models.TitleTypeFeature,
			PremiereDate:         premiere,
			Territories:          50,
			CurrentMasterVersion: "V01",
			OverallStatus:        models.StatusProcessing,
		})

		var mu sync.Mutex
		var capturedPrompt string
		capturingMock := &mockLLM{
			capturePrompt: func(p string) {
				mu.Lock()
				capturedPrompt = p
				mu.Unlock()
			},
			responses: []string{
				`{"actionable":false,"anomaly_type":"BENIGN_TELEMETRY","severity":"INFO","rationale":"Routine heartbeat ping"}`,
			},
		}
		server.SetModel(capturingMock)

		body, _ := json.Marshal(runs.CreateRunRequest{
			TitleSlug: "avatar-fire-ash",
			Event: &models.Event{
				ID:       "avatar-ping-1",
				Type:     "vendor.ping",
				Subject:  "avatar-fire-ash",
				Time:     time.Now().UTC(),
				Severity: models.SeverityInfo,
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/runs?wf=incident", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got: %d (%s)", rec.Code, rec.Body.String())
		}

		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		prompt := capturedPrompt
		mu.Unlock()

		// Confirm that the triage prompt received the computed hours (~10h), NOT the hardcoded 48.0h!
		if !strings.Contains(prompt, "Hours Until Premiere: 9.") && !strings.Contains(prompt, "Hours Until Premiere: 10.") {
			t.Errorf("expected prompt to contain ~10 hours until premiere, got: %s", prompt)
		}
	})
}

type mockLLM struct {
	mu            sync.Mutex
	responses     []string
	callIndex     int
	capturePrompt func(string)
}

func (m *mockLLM) Name() string {
	return "mock-gemini"
}

func (m *mockLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		m.mu.Lock()
		if m.capturePrompt != nil && req != nil && len(req.Contents) > 0 {
			var sb strings.Builder
			for _, c := range req.Contents {
				for _, p := range c.Parts {
					sb.WriteString(p.Text)
				}
			}
			m.capturePrompt(sb.String())
		}

		respText := ""
		if len(m.responses) > 0 {
			if m.callIndex < len(m.responses) {
				respText = m.responses[m.callIndex]
				m.callIndex++
			} else {
				respText = m.responses[len(m.responses)-1]
			}
		}
		m.mu.Unlock()

		resp := &model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{Text: respText},
				},
			},
		}
		yield(resp, nil)
	}
}
