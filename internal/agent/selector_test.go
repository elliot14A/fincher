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
			Specialty:          "AUDIO",
			HourlyRateUSD:      150.0,
			TurnaroundHours:    24,
			HistoricalAccuracy: 0.98,
		},
		{
			VendorID:           "vendor-efficient",
			VendorName:         "Efficient Dubs",
			Specialty:          "AUDIO",
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
