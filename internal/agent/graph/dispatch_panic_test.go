package graph_test

import (
	"context"
	"iter"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent/graph"
	"github.com/elliot14A/fincher/internal/turso/runs"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"google.golang.org/adk/v2/model"
)

type panicLLM struct{}

func (p *panicLLM) Name() string { return "panic-llm" }

func (p *panicLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	panic("simulated llm crash inside background workflow")
}

func TestDispatch_PanicRecoveryMarksRunFailed(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)
	defer client.Close()

	llm := &panicLLM{}

	event := &models.Event{
		ID:       "evt-panic-test",
		Type:     models.TypeAudioSyncDriftDetected,
		Severity: models.SeverityCritical,
		Subject:  "panic-title",
		Data: map[string]any{
			"package_id": "pkg-panic",
		},
	}

	deps := graph.IncidentGraphDeps{
		Model:       llm,
		TursoClient: client,
		MaxAttempts: 1,
	}

	runObj, dispatched, err := graph.DispatchIncident(ctx, deps, graph.IncidentInput{
		Event: event,
	})
	if err != nil {
		t.Fatalf("DispatchIncident returned error: %v", err)
	}
	if !dispatched {
		t.Fatal("expected dispatched true")
	}

	// Poll Turso for the run status transition to FAILED with timeout
	deadline := time.Now().Add(2 * time.Second)
	var finalRun *models.Run

	for time.Now().Before(deadline) {
		res := runs.GetRun(ctx, client, runObj.ID)
		if res.IsOk() {
			r := res.Unwrap()
			if r.Status == models.RunStatusFailed {
				finalRun = r
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
	}

	if finalRun == nil {
		t.Fatalf("expected run %s to transition to FAILED after panic, but remained running or not found", runObj.ID)
	}
	if finalRun.Status != models.RunStatusFailed {
		t.Errorf("expected status FAILED, got: %s", finalRun.Status)
	}
	if finalRun.Metadata == nil || finalRun.Metadata["panic"] == nil {
		t.Errorf("expected panic details recorded in run metadata, got: %v", finalRun.Metadata)
	}
}

type succeedThenPanicLLM struct {
	client *models.Run
}

func TestDispatch_PanicRecoveryDoesNotClobberTerminalRun(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)
	defer client.Close()

	runID := "run-terminal-guard-test"
	now := time.Now().UTC()

	// Pre-create run in terminal COMPLETED status
	_ = runs.CreateRun(ctx, client, &models.Run{
		Base: models.Base{
			ID: runID,
			Metadata: map[string]any{
				"custom_success_key": "kept",
			},
		},
		TitleSlug: "terminal-title",
		Trigger:   "incident",
		Status:    models.RunStatusCompleted,
		StartedAt: now,
		EndedAt:   &now,
	})

	// Dispatch incident with same runID -> Idempotency skips dispatch
	llm := &panicLLM{}
	event := &models.Event{
		ID:       "terminal-guard-test",
		Type:     models.TypeAudioSyncDriftDetected,
		Severity: models.SeverityCritical,
		Subject:  "terminal-title",
	}

	deps := graph.IncidentGraphDeps{
		Model:       llm,
		TursoClient: client,
		MaxAttempts: 1,
	}

	runObj, dispatched, err := graph.DispatchIncident(ctx, deps, graph.IncidentInput{
		RunID: runID,
		Event: event,
	})
	if err != nil {
		t.Fatalf("DispatchIncident returned error: %v", err)
	}
	if dispatched {
		t.Errorf("expected dispatched false for existing terminal run")
	}
	if runObj.Status != models.RunStatusCompleted {
		t.Errorf("expected status COMPLETED preserved, got: %s", runObj.Status)
	}
}
