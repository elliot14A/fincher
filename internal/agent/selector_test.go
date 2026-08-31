package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elliot14A/fincher/internal/agent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestSelectVendor(t *testing.T) {
	ctx := context.Background()

	candidates := []models.VendorCandidate{
		{
			VendorID:           "vendor-fast",
			VendorName:         "Fast Dubs",
			Components:         []string{"AUDIO"},
			Markets:            []string{"en-US"},
			HourlyRateUSD:      150.0,
			TurnaroundHours:    24,
			HistoricalAccuracy: 0.98,
		},
		{
			VendorID:           "vendor-efficient",
			VendorName:         "Efficient Dubs",
			Components:         []string{"AUDIO"},
			Markets:            []string{"en-US"},
			HourlyRateUSD:      75.0,
			TurnaroundHours:    48,
			HistoricalAccuracy: 0.96,
		},
	}

	t.Run("Fails when model is nil", func(t *testing.T) {
		res := agent.SelectVendor(ctx, nil, "eclipse", "AUDIO", candidates, 72.0)
		if res.IsOk() {
			t.Fatal("expected error for nil model, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeInvalidInput {
			t.Errorf("expected INVALID_INPUT code, got: %s", domErr.Code)
		}
	})

	t.Run("Fails when candidates list is empty", func(t *testing.T) {
		llm := &mockLLM{response: `{}`}
		res := agent.SelectVendor(ctx, llm, "eclipse", "AUDIO", []models.VendorCandidate{}, 72.0)
		if res.IsOk() {
			t.Fatal("expected error for empty candidates, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeInvalidInput {
			t.Errorf("expected INVALID_INPUT code, got: %s", domErr.Code)
		}
	})

	t.Run("Selects optimal vendor partner based on cost and feasibility", func(t *testing.T) {
		jsonOutput := `{
			"winner_vendor_id": "vendor-efficient",
			"winner_vendor_name": "Efficient Dubs",
			"hourly_rate_usd": 75.0,
			"turnaround_hours": 48,
			"rationale": "Efficient Dubs meets the 72h premiere timeline with 48h turnaround and offers the lowest rate ($75/h) while maintaining 96% accuracy."
		}`
		llm := &mockLLM{response: jsonOutput}

		res := agent.SelectVendor(ctx, llm, "eclipse", "AUDIO", candidates, 72.0)
		if res.IsErr() {
			t.Fatalf("SelectVendor returned error: %v", res.Error())
		}
		decision := res.Unwrap()

		if decision.WinnerVendorID != "vendor-efficient" {
			t.Errorf("expected winner vendor-efficient, got: %s", decision.WinnerVendorID)
		}
		if decision.HourlyRateUSD != 75.0 {
			t.Errorf("expected $75.0 hourly rate, got: %f", decision.HourlyRateUSD)
		}
		if decision.TurnaroundHours != 48 {
			t.Errorf("expected 48h turnaround, got: %d", decision.TurnaroundHours)
		}
	})

	t.Run("Fails when model returns malformed JSON", func(t *testing.T) {
		llm := &mockLLM{response: "invalid non json response"}

		res := agent.SelectVendor(ctx, llm, "eclipse", "AUDIO", candidates, 72.0)
		if res.IsOk() {
			t.Fatal("expected error for malformed JSON, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeInternal {
			t.Errorf("expected INTERNAL code, got: %s", domErr.Code)
		}
	})

	t.Run("Propagates model generation error", func(t *testing.T) {
		llm := &mockLLM{err: errors.New("upstream quota exceeded")}

		res := agent.SelectVendor(ctx, llm, "eclipse", "AUDIO", candidates, 72.0)
		if res.IsOk() {
			t.Fatal("expected error when model fails, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeBudgetExceeded {
			t.Errorf("expected BUDGET_EXCEEDED for quota error, got: %s", domErr.Code)
		}
	})
}

func TestSelectVendorsForPlan(t *testing.T) {
	ctx := context.Background()

	reqs := []models.AllocationRequirement{
		{Component: "VIDEO", Market: "", Language: "en-US"},
		{Component: "AUDIO", Market: "de-DE", Language: "de-DE"},
		{Component: "SUBTITLE", Market: "de-DE", Language: "de-DE"},
	}

	candPools := map[string][]models.VendorCandidate{
		"VIDEO|": {
			{VendorID: "vnd-technicolor", VendorName: "Technicolor", HourlyRateUSD: 185.0, TurnaroundHours: 16, HistoricalAccuracy: 0.98},
		},
		"AUDIO|de-DE": {
			{VendorID: "vnd-iyuno", VendorName: "Iyuno SDI", HourlyRateUSD: 70.0, TurnaroundHours: 36, HistoricalAccuracy: 0.93},
			{VendorID: "vnd-deluxe", VendorName: "Deluxe Media", HourlyRateUSD: 200.0, TurnaroundHours: 12, HistoricalAccuracy: 0.99},
		},
		"SUBTITLE|de-DE": {
			{VendorID: "vnd-pixelogic", VendorName: "Pixelogic", HourlyRateUSD: 80.0, TurnaroundHours: 8, HistoricalAccuracy: 0.96},
		},
	}

	t.Run("Selects holistic multi-requirement plan with single LLM call", func(t *testing.T) {
		jsonOutput := `{
			"assignments": [
				{
					"component": "VIDEO",
					"market": "",
					"language": "en-US",
					"winner_vendor_id": "vnd-technicolor",
					"winner_vendor_name": "Technicolor",
					"hourly_rate_usd": 185.0,
					"turnaround_hours": 16,
					"rationale": "Sole qualified video/QC facility."
				},
				{
					"component": "AUDIO",
					"market": "de-DE",
					"language": "de-DE",
					"winner_vendor_id": "vnd-iyuno",
					"winner_vendor_name": "Iyuno SDI",
					"hourly_rate_usd": 70.0,
					"turnaround_hours": 36,
					"rationale": "Lowest cost qualified German dubbing partner."
				},
				{
					"component": "SUBTITLE",
					"market": "de-DE",
					"language": "de-DE",
					"winner_vendor_id": "vnd-pixelogic",
					"winner_vendor_name": "Pixelogic",
					"hourly_rate_usd": 80.0,
					"turnaround_hours": 8,
					"rationale": "High accuracy subtitle specialist with rapid 8h turnaround."
				}
			],
			"overall_summary": "All 3 requirements staffed with top-tier vendors within 48h premiere buffer."
		}`
		llm := &mockLLM{response: jsonOutput}

		res := agent.SelectVendorsForPlan(ctx, llm, "avatar-fire-ash", reqs, candPools, 72.0)
		if res.IsErr() {
			t.Fatalf("SelectVendorsForPlan returned error: %v", res.Error())
		}
		plan := res.Unwrap()

		if len(plan.Assignments) != 3 {
			t.Fatalf("expected 3 assignments, got: %d", len(plan.Assignments))
		}
		if plan.Assignments[0].WinnerVendorID != "vnd-technicolor" {
			t.Errorf("expected video winner vnd-technicolor, got: %s", plan.Assignments[0].WinnerVendorID)
		}
		if plan.Assignments[1].WinnerVendorID != "vnd-iyuno" {
			t.Errorf("expected audio winner vnd-iyuno, got: %s", plan.Assignments[1].WinnerVendorID)
		}
		if plan.Assignments[2].WinnerVendorID != "vnd-pixelogic" {
			t.Errorf("expected sub winner vnd-pixelogic, got: %s", plan.Assignments[2].WinnerVendorID)
		}
	})

	t.Run("Enforces no_eligible_vendor and completes dropped requirements defensively in Go", func(t *testing.T) {
		reqsWithEmpty := []models.AllocationRequirement{
			{Component: "VIDEO", Market: "", Language: "en-US"},
			{Component: "AUDIO", Market: "te-IN", Language: "te-IN"}, // empty candidate pool
			{Component: "SUBTITLE", Market: "te-IN", Language: "te-IN"},
		}

		candPoolsWithEmpty := map[string][]models.VendorCandidate{
			"VIDEO|": {
				{VendorID: "vnd-technicolor", VendorName: "Technicolor", HourlyRateUSD: 185.0, TurnaroundHours: 16, HistoricalAccuracy: 0.98},
			},
			"AUDIO|te-IN": {}, // empty pool
			"SUBTITLE|te-IN": {
				{VendorID: "vnd-pixelogic", VendorName: "Pixelogic", HourlyRateUSD: 80.0, TurnaroundHours: 8, HistoricalAccuracy: 0.96},
			},
		}

		// Model only returns assignment for VIDEO, dropping AUDIO and SUBTITLE
		partialJSON := `{
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
				}
			],
			"overall_summary": "Partial model response"
		}`
		llm := &mockLLM{response: partialJSON}

		res := agent.SelectVendorsForPlan(ctx, llm, "avatar-fire-ash", reqsWithEmpty, candPoolsWithEmpty, 72.0)
		if res.IsErr() {
			t.Fatalf("SelectVendorsForPlan returned error: %v", res.Error())
		}
		plan := res.Unwrap()

		if len(plan.Assignments) != 3 {
			t.Fatalf("expected 3 verified assignments, got: %d", len(plan.Assignments))
		}

		// 1. VIDEO matched model assignment
		if plan.Assignments[0].WinnerVendorID != "vnd-technicolor" {
			t.Errorf("expected video winner vnd-technicolor, got: %s", plan.Assignments[0].WinnerVendorID)
		}

		// 2. AUDIO (empty pool) enforced no_eligible_vendor
		if plan.Assignments[1].WinnerVendorID != "no_eligible_vendor" {
			t.Errorf("expected audio winner no_eligible_vendor, got: %s", plan.Assignments[1].WinnerVendorID)
		}

		// 3. SUBTITLE (dropped by model, pool available) recovered top candidate
		if plan.Assignments[2].WinnerVendorID != "vnd-pixelogic" {
			t.Errorf("expected subtitle recovered winner vnd-pixelogic, got: %s", plan.Assignments[2].WinnerVendorID)
		}
	})

	t.Run("Overrides ineligible model pick with preferred candidate (accuracy > rate > turnaround)", func(t *testing.T) {
		reqsOne := []models.AllocationRequirement{
			{Component: "AUDIO", Market: "de-DE", Language: "de-DE"},
		}
		// Pool ordered so the FIRST entry is NOT the best: a low-accuracy cheap vendor first,
		// a high-accuracy vendor second. bestCandidate must prefer accuracy over rate.
		poolsMulti := map[string][]models.VendorCandidate{
			"AUDIO|de-DE": {
				{VendorID: "vnd-cheap-weak", VendorName: "Cheap Weak", HourlyRateUSD: 50.0, TurnaroundHours: 40, HistoricalAccuracy: 0.80},
				{VendorID: "vnd-deluxe", VendorName: "Deluxe Media", HourlyRateUSD: 200.0, TurnaroundHours: 12, HistoricalAccuracy: 0.99},
			},
		}
		// Model hallucinates a vendor not in the pool.
		hallucinatedJSON := `{
			"assignments": [
				{
					"component": "AUDIO",
					"market": "de-DE",
					"language": "de-DE",
					"winner_vendor_id": "vnd-does-not-exist",
					"winner_vendor_name": "Ghost Vendor",
					"hourly_rate_usd": 10.0,
					"turnaround_hours": 5,
					"rationale": "Hallucinated."
				}
			],
			"overall_summary": "Invalid pick"
		}`
		llm := &mockLLM{response: hallucinatedJSON}

		res := agent.SelectVendorsForPlan(ctx, llm, "avatar-fire-ash", reqsOne, poolsMulti, 72.0)
		if res.IsErr() {
			t.Fatalf("SelectVendorsForPlan returned error: %v", res.Error())
		}
		plan := res.Unwrap()

		if len(plan.Assignments) != 1 {
			t.Fatalf("expected 1 assignment, got: %d", len(plan.Assignments))
		}
		// Must override the hallucinated pick with the HIGHEST-ACCURACY eligible vendor,
		// not simply the first in the pool (which is cheaper but weaker).
		if plan.Assignments[0].WinnerVendorID != "vnd-deluxe" {
			t.Errorf("expected override to preferred vnd-deluxe (highest accuracy), got: %s", plan.Assignments[0].WinnerVendorID)
		}
	})
}
