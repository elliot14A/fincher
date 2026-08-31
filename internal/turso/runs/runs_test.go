package runs_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/runs"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func setupTestDB(t *testing.T) *ent.Client {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to run schema automigrations: %v", err)
	}

	return client
}

func TestRuns_Lifecycle(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create Run
	runInput := &models.Run{
		Base: models.Base{
			ID:       "run-123",
			Metadata: map[string]any{"source_event_id": "evt-abc"},
		},
		TitleSlug: "eclipse",
		Trigger:   "ANOMALY_SIGNAL",
		Status:    models.RunStatusRunning,
		StartedAt: time.Now().UTC(),
	}

	createRunRes := runs.CreateRun(ctx, client, runInput)
	if createRunRes.IsErr() {
		t.Fatalf("creating run failed: %v", createRunRes.Error())
	}
	createdRun := createRunRes.Unwrap()
	if createdRun.ID != "run-123" {
		t.Fatalf("expected ID run-123, got: %s", createdRun.ID)
	}
	if createdRun.TitleSlug != "eclipse" {
		t.Fatalf("expected title slug eclipse, got: %s", createdRun.TitleSlug)
	}

	// 2. Create Step
	stepInput := &models.Step{
		Base: models.Base{
			ID: "step-1",
		},
		RunID:     "run-123",
		Name:      "policy_judge",
		Status:    models.StepStatusRunning,
		StartedAt: time.Now().UTC(),
	}

	createStepRes := runs.CreateStep(ctx, client, stepInput)
	if createStepRes.IsErr() {
		t.Fatalf("creating step failed: %v", createStepRes.Error())
	}
	createdStep := createStepRes.Unwrap()
	if createdStep.Name != "policy_judge" {
		t.Fatalf("expected step name policy_judge, got: %s", createdStep.Name)
	}

	// 3. Create Result
	resInput := &models.WfResult{
		Base: models.Base{
			ID: "res-1",
		},
		RunID:     "run-123",
		StepID:    "step-1",
		Judge:     "policy_judge",
		Outcome:   "APPROVED",
		Rationale: "All SLA conditions met",
		Attempt:   1,
	}

	createResRes := runs.CreateResult(ctx, client, resInput)
	if createResRes.IsErr() {
		t.Fatalf("creating result failed: %v", createResRes.Error())
	}
	createdResult := createResRes.Unwrap()
	if createdResult.Outcome != "APPROVED" {
		t.Fatalf("expected outcome APPROVED, got: %s", createdResult.Outcome)
	}

	// 4. Update Step & Run Status
	now := time.Now().UTC()
	updateStepRes := runs.UpdateStepStatus(ctx, client, "step-1", models.StepStatusCompleted, &now, map[string]any{"duration_ms": 450})
	if updateStepRes.IsErr() {
		t.Fatalf("updating step status failed: %v", updateStepRes.Error())
	}

	updateRunRes := runs.UpdateRunStatus(ctx, client, "run-123", models.RunStatusCompleted, &now, nil)
	if updateRunRes.IsErr() {
		t.Fatalf("updating run status failed: %v", updateRunRes.Error())
	}

	// 5. GetRun with loaded steps & results
	getRes := runs.GetRun(ctx, client, "run-123")
	if getRes.IsErr() {
		t.Fatalf("get run failed: %v", getRes.Error())
	}
	loadedRun := getRes.Unwrap()
	if loadedRun.Status != models.RunStatusCompleted {
		t.Errorf("expected run status COMPLETED, got: %s", loadedRun.Status)
	}
	if len(loadedRun.Steps) != 1 {
		t.Fatalf("expected 1 step in run, got: %d", len(loadedRun.Steps))
	}
	if len(loadedRun.Results) != 1 {
		t.Fatalf("expected 1 result in run, got: %d", len(loadedRun.Results))
	}

	// 6. ListRuns pagination
	listRes := runs.ListRuns(ctx, client, runs.ListFilter{TitleSlug: domainerrors.Some("eclipse")}, models.Pagination{Limit: 10, Page: 1})
	if listRes.IsErr() {
		t.Fatalf("list runs failed: %v", listRes.Error())
	}
	page := listRes.Unwrap()
	if page.TotalItems != 1 {
		t.Fatalf("expected 1 total item, got: %d", page.TotalItems)
	}
}
