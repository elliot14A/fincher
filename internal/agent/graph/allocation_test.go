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

	t.Run("Evaluates candidates and selects optimal vendor partner", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)

		vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "vendor-deluxe",
			},
			Name:            "Deluxe Audio",
			Specialty:       "AUDIO_DUBBING",
			HourlyRateUSD:   200.0,
			TurnaroundHours: 12,
		})

		vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "vendor-berlin",
			},
			Name:            "Berlin Synchron",
			Specialty:       "AUDIO_DUBBING",
			HourlyRateUSD:   120.0,
			TurnaroundHours: 24,
		})

		selectorResponse := `{
			"winner_vendor_id": "vendor-berlin",
			"winner_vendor_name": "Berlin Synchron",
			"hourly_rate_usd": 120.0,
			"turnaround_hours": 24,
			"rationale": "Meets the 48h turnaround requirement comfortably with the lowest commercial hourly rate ($120/hr)."
		}`

		llm := &mockLLM{
			responses: []string{selectorResponse},
		}

		deps := graph.AllocationGraphDeps{
			Model:       llm,
			TursoClient: client,
		}

		output, err := graph.ExecuteAllocation(ctx, deps, graph.AllocationInput{
			TitleSlug:          "eclipse",
			Component:          "AUDIO_DUBBING",
			HoursUntilPremiere: 48.0,
		})
		if err != nil {
			t.Fatalf("ExecuteAllocation returned error: %v", err)
		}

		if output.Decision == nil {
			t.Fatal("expected non-nil Decision")
		}
		if output.Decision.WinnerVendorID != "vendor-berlin" {
			t.Errorf("expected winner vendor-berlin, got: %s", output.Decision.WinnerVendorID)
		}
		if output.Decision.HourlyRateUSD != 120.0 {
			t.Errorf("expected hourly rate 120.0, got: %f", output.Decision.HourlyRateUSD)
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
			Component: "AUDIO_DUBBING",
		})
		if err == nil {
			t.Fatal("expected error for empty title slug, got nil")
		}
	})
}
