package graph_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/agent/graph"
	"github.com/elliot14A/fincher/internal/turso/runs"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestExecuteAllocation(t *testing.T) {
	ctx := context.Background()

	t.Run("Evaluates multi-requirement candidates and creates holistic staffing plan", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)

		vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "vendor-deluxe",
			},
			Name:            "Deluxe Audio",
			Components:      []string{"AUDIO"},
			Markets:         []string{"en-US", "de-DE"},
			HourlyRateUSD:   200.0,
			TurnaroundHours: 12,
		})

		vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "vendor-berlin",
			},
			Name:            "Berlin Synchron",
			Components:      []string{"AUDIO"},
			Markets:         []string{"en-US", "de-DE"},
			HourlyRateUSD:   120.0,
			TurnaroundHours: 24,
		})

		selectorResponse := `{
			"assignments": [
				{
					"component": "AUDIO",
					"market": "de-DE",
					"language": "de-DE",
					"winner_vendor_id": "vendor-berlin",
					"winner_vendor_name": "Berlin Synchron",
					"hourly_rate_usd": 120.0,
					"turnaround_hours": 24,
					"rationale": "Meets the 48h turnaround requirement comfortably with lowest rate ($120/hr)."
				}
			],
			"overall_summary": "German dubbing assigned to Berlin Synchron."
		}`

		llm := &mockLLM{
			responses: []string{selectorResponse},
		}

		deps := graph.AllocationGraphDeps{
			Model:       llm,
			TursoClient: client,
		}

		output, err := graph.ExecuteAllocation(ctx, deps, graph.AllocationInput{
			TitleSlug: "eclipse",
			Requirements: []models.AllocationRequirement{
				{
					Component: "AUDIO",
					Market:    "de-DE",
					Language:  "de-DE",
				},
			},
			HoursUntilPremiere: 48.0,
		})
		if err != nil {
			t.Fatalf("ExecuteAllocation returned error: %v", err)
		}

		if output.Plan == nil {
			t.Fatal("expected non-nil Plan")
		}
		if len(output.Plan.Assignments) != 1 {
			t.Fatalf("expected 1 assignment, got: %d", len(output.Plan.Assignments))
		}
		if output.Plan.Assignments[0].WinnerVendorID != "vendor-berlin" {
			t.Errorf("expected winner vendor-berlin, got: %s", output.Plan.Assignments[0].WinnerVendorID)
		}

		// Verify that candidate_gathering and vendor_selection steps were created
		runListRes := runs.ListRuns(ctx, client, runs.ListFilter{TitleSlug: domainerrors.Some("eclipse")}, models.Pagination{Limit: 10, Page: 1})
		if runListRes.IsErr() {
			t.Fatalf("failed to list runs: %v", runListRes.Error())
		}
		if len(runListRes.Unwrap().Items) < 1 {
			t.Fatal("expected at least 1 run in Turso")
		}
		loadedRun := runListRes.Unwrap().Items[0]
		if len(loadedRun.Steps) != 2 {
			t.Fatalf("expected 2 steps (candidates, selection), got: %d", len(loadedRun.Steps))
		}
		if loadedRun.Steps[0].Name != "candidate_gathering" {
			t.Errorf("step 0: expected candidate_gathering, got: %s", loadedRun.Steps[0].Name)
		}
		if loadedRun.Steps[1].Name != "vendor_selection" {
			t.Errorf("step 1: expected vendor_selection, got: %s", loadedRun.Steps[1].Name)
		}
		if len(loadedRun.Results) != 1 {
			t.Errorf("expected 1 selection WfResult, got: %d", len(loadedRun.Results))
		}
	})

	t.Run("Fails when title slug is empty", func(t *testing.T) {
		deps := graph.AllocationGraphDeps{
			Model: &mockLLM{},
		}
		_, err := graph.ExecuteAllocation(ctx, deps, graph.AllocationInput{
			TitleSlug: "",
			Requirements: []models.AllocationRequirement{
				{Component: "AUDIO", Market: "en-US"},
			},
		})
		if err == nil {
			t.Fatal("expected error for empty title slug, got nil")
		}
	})

	t.Run("Fails when requirements list is empty", func(t *testing.T) {
		deps := graph.AllocationGraphDeps{
			Model: &mockLLM{},
		}
		_, err := graph.ExecuteAllocation(ctx, deps, graph.AllocationInput{
			TitleSlug:    "avatar-fire-ash",
			Requirements: []models.AllocationRequirement{},
		})
		if err == nil {
			t.Fatal("expected error for empty requirements, got nil")
		}
	})
}
