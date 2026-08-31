package titles_test

import (
	"bytes"
	"context"
	"encoding/json"
	"iter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/elliot14A/fincher/internal/api"
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

type mockLLM struct {
	response string
}

func (m *mockLLM) Name() string { return "mock-llm" }

func (m *mockLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		resp := &model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{Text: m.response},
				},
			},
		}
		yield(resp, nil)
	}
}

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

	// Verify bundled master was created in Turso
	mList, err := client.Master.Query().Where().All(context.Background())
	if err != nil {
		t.Fatalf("failed to query masters: %v", err)
	}
	if len(mList) != 1 {
		t.Fatalf("expected 1 bundled master in database, got: %d", len(mList))
	}
	if mList[0].Version != "V12" {
		t.Errorf("expected bundled master version V12, got: %s", mList[0].Version)
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

	// 5. DELETE /api/titles/:id -> 409 Conflict (blocked by bundled master)
	req = httptest.NewRequest(http.MethodDelete, "/api/titles/title-eclipse", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409 conflict when deleting title with master, got %d", rec.Code)
	}

	// Clean up bundled master and retry delete
	_ = client.Master.DeleteOneID("mst-eclipse-v12").Exec(context.Background())

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

func TestTitles_Onboarding_WithAllocation(t *testing.T) {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open memory ent client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to run ent automigration: %v", err)
	}

	// Seed vendors for allocation
	_ = vendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vnd-technicolor"},
		Name:            "Technicolor",
		Components:      []string{"VIDEO"},
		Markets:         []string{},
		HourlyRateUSD:   185.0,
		TurnaroundHours: 16,
	})
	_ = vendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vnd-deluxe"},
		Name:            "Deluxe Media",
		Components:      []string{"AUDIO", "SUBTITLE"},
		Markets:         []string{"en-US", "de-DE"},
		HourlyRateUSD:   200.0,
		TurnaroundHours: 12,
	})

	mockPlan := `{
		"assignments": [
			{
				"component": "VIDEO",
				"market": "",
				"language": "en-US",
				"winner_vendor_id": "vnd-technicolor",
				"winner_vendor_name": "Technicolor",
				"hourly_rate_usd": 185.0,
				"turnaround_hours": 16,
				"rationale": "Sole video provider."
			},
			{
				"component": "AUDIO",
				"market": "en-US",
				"language": "en-US",
				"winner_vendor_id": "vnd-deluxe",
				"winner_vendor_name": "Deluxe Media",
				"hourly_rate_usd": 200.0,
				"turnaround_hours": 12,
				"rationale": "High quality audio."
			},
			{
				"component": "SUBTITLE",
				"market": "en-US",
				"language": "en-US",
				"winner_vendor_id": "vnd-deluxe",
				"winner_vendor_name": "Deluxe Media",
				"hourly_rate_usd": 200.0,
				"turnaround_hours": 12,
				"rationale": "Rapid subtitle turnaround."
			}
		],
		"overall_summary": "Staffing plan complete."
	}`

	mock := &mockLLM{response: mockPlan}
	server := api.NewServer(client)
	server.SetModel(mock)
	e := server.Router()

	titlePayload := models.Title{
		Base: models.Base{
			ID: "title-avatar-fire-ash",
			Metadata: map[string]any{
				"markets": []string{"en-US"},
			},
		},
		Name:                 "Avatar: Fire and Ash",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(72 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusOnTrack,
	}

	body, _ := json.Marshal(titlePayload)
	req := httptest.NewRequest(http.MethodPost, "/api/titles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got: %d (%s)", rec.Code, rec.Body.String())
	}

	// 1. Verify Master V01 was bundled
	mList, err := client.Master.Query().Where().All(ctx)
	if err != nil || len(mList) != 1 {
		t.Fatalf("expected 1 bundled master, got: %d (err: %v)", len(mList), err)
	}

	// Wait briefly for background allocation run
	time.Sleep(150 * time.Millisecond)

	// 2. Verify single allocation Run was dispatched in Turso
	runsList, err := client.Run.Query().Where().All(ctx)
	if err != nil || len(runsList) != 1 {
		t.Fatalf("expected exactly 1 allocation run dispatched, got: %d (err: %v)", len(runsList), err)
	}
	if runsList[0].Trigger != "allocation" {
		t.Errorf("expected trigger 'allocation', got: %s", runsList[0].Trigger)
	}

	// 3. Verify Run has Steps and Results
	stepsList, _ := client.Step.Query().Where().All(ctx)
	if len(stepsList) != 2 {
		t.Fatalf("expected 2 steps in run, got: %d", len(stepsList))
	}
	resultsList, _ := client.WfResult.Query().Where().All(ctx)
	if len(resultsList) != 3 {
		t.Fatalf("expected 3 assignment results, got: %d", len(resultsList))
	}
}
