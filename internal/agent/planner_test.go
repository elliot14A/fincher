package agent_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestPlanRemediation(t *testing.T) {
	ctx := context.Background()

	event := &models.Event{
		ID:       "evt-defect-1",
		Type:     models.TypeQCInspectionCompleted,
		Severity: models.SeverityCritical,
		Subject:  "eclipse",
		Time:     time.Now().UTC(),
		Data: map[string]any{
			"package_id":    "pkg-de-dub",
			"status":        "FAILED",
			"sync_drift_ms": 145.0,
		},
	}

	impact := &models.DeliveryImpact{
		RootPackageID:      "pkg-de-dub",
		AffectedPackages:   []string{"pkg-de-dub"},
		AffectedDeliveries: []string{"del-germany"},
		AffectedMarkets:    []string{"DE"},
		HoursUntilPremiere: 48.0,
		IsPremiereUrgent:   true,
	}

	analytics := &models.AnalyticsSummary{
		VendorHistoricalAccuracy: 0.82,
		SimilarDefectOccurrences: 2,
		PriorIncidentsForVendor:  3,
	}

	candidates := []models.VendorCandidate{
		{
			VendorID:           "vendor-berlin",
			VendorName:         "Berlin Synchron",
			Components:         []string{"AUDIO"},
			Markets:            []string{"en-US"},
			HourlyRateUSD:      120.0,
			TurnaroundHours:    24,
			HistoricalAccuracy: 0.98,
		},
	}

	t.Run("Fails when model is nil", func(t *testing.T) {
		res := agent.PlanRemediation(ctx, nil, event, impact, analytics, candidates, nil, "")
		if res.IsOk() {
			t.Fatal("expected error for nil model, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeInvalidInput {
			t.Errorf("expected INVALID_INPUT code, got: %s", domErr.Code)
		}
	})

	t.Run("Fails when event is nil", func(t *testing.T) {
		llm := &mockLLM{response: `{}`}
		res := agent.PlanRemediation(ctx, llm, nil, impact, analytics, candidates, nil, "")
		if res.IsOk() {
			t.Fatal("expected error for nil event, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeInvalidInput {
			t.Errorf("expected INVALID_INPUT code, got: %s", domErr.Code)
		}
	})

	t.Run("Synthesizes compliant remediation action plan", func(t *testing.T) {
		jsonOutput := `{
			"title_slug": "eclipse",
			"summary": "Hold German delivery and reassign to Berlin Synchron due to 145ms audio drift.",
			"actions": [
				{
					"type": "HOLD_DELIVERY",
					"target_id": "del-germany",
					"reason": "Audio sync drift of 145ms requires reconform before release."
				},
				{
					"type": "REASSIGN_VENDOR",
					"target_id": "vendor-berlin",
					"reason": "Berlin Synchron has 24h turnaround and 98% accuracy.",
					"payload": {
						"package_id": "pkg-de-dub"
					}
				},
				{
					"type": "EMAIL_VENDOR",
					"target_id": "vendor-berlin",
					"reason": "Dispatch expedited German dub reconform."
				},
				{
					"type": "NOTIFY_STAKEHOLDERS",
					"target_id": "slack-ops",
					"reason": "Alert distribution lead of German hold."
				}
			]
		}`
		llm := &mockLLM{response: jsonOutput}

		res := agent.PlanRemediation(ctx, llm, event, impact, analytics, candidates, nil, "")
		if res.IsErr() {
			t.Fatalf("PlanRemediation returned error: %v", res.Error())
		}
		plan := res.Unwrap()

		if plan.TitleSlug != "eclipse" {
			t.Errorf("expected title slug eclipse, got: %s", plan.TitleSlug)
		}
		if len(plan.Actions) != 4 {
			t.Fatalf("expected 4 actions, got: %d", len(plan.Actions))
		}
		if plan.Actions[0].Type != models.ActionHoldDelivery {
			t.Errorf("expected first action HOLD_DELIVERY, got: %s", plan.Actions[0].Type)
		}
		if plan.Actions[1].Type != models.ActionReassignVendor {
			t.Errorf("expected second action REASSIGN_VENDOR, got: %s", plan.Actions[1].Type)
		}

		verifRes := agent.VerifyPlan(plan, impact, candidates, nil, 1)
		if verifRes.IsErr() {
			t.Fatalf("VerifyPlan error: %v", verifRes.Error())
		}
		if verifRes.Unwrap().Decision != agent.DecisionApproved {
			t.Fatalf("expected plan to pass verification, got: %s (rationale: %s)", verifRes.Unwrap().Decision, verifRes.Unwrap().Rationale)
		}
	})

	t.Run("Fails when model returns malformed JSON", func(t *testing.T) {
		llm := &mockLLM{response: "malformed non-json response"}

		res := agent.PlanRemediation(ctx, llm, event, impact, analytics, candidates, nil, "")
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

		res := agent.PlanRemediation(ctx, llm, event, impact, analytics, candidates, nil, "")
		if res.IsOk() {
			t.Fatal("expected error when model fails, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeBudgetExceeded {
			t.Errorf("expected BUDGET_EXCEEDED for quota error, got: %s", domErr.Code)
		}
	})
}
