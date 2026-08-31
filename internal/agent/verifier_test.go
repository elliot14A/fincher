package agent_test

import (
	"testing"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestVerifyPlan(t *testing.T) {
	impact := &models.DeliveryImpact{
		RootPackageID:      "pkg-german-dub",
		AffectedPackages:   []string{"pkg-german-dub"},
		AffectedDeliveries: []string{"del-germany", "del-austria"},
		AffectedMarkets:    []string{"DE", "AT"},
		HoursUntilPremiere: 48.0,
		IsPremiereUrgent:   true,
	}

	vendors := []models.VendorCandidate{
		{
			VendorID:           "vendor-berlin",
			VendorName:         "Berlin Synchron",
			Components:         []string{"AUDIO"},
			Markets:            []string{"de-DE"},
			HourlyRateUSD:      120.0,
			TurnaroundHours:    24,
			HistoricalAccuracy: 0.98,
		},
		{
			VendorID:           "vendor-slow",
			VendorName:         "Slow Dubs Inc",
			Components:         []string{"AUDIO"},
			Markets:            []string{"de-DE"},
			HourlyRateUSD:      60.0,
			TurnaroundHours:    72,
			HistoricalAccuracy: 0.95,
		},
		{
			VendorID:           "vendor-sloppy",
			VendorName:         "Sloppy Audio",
			Components:         []string{"AUDIO"},
			Markets:            []string{"de-DE"},
			HourlyRateUSD:      40.0,
			TurnaroundHours:    12,
			HistoricalAccuracy: 0.75,
		},
		{
			VendorID:           "vendor-new",
			VendorName:         "New Studio",
			Components:         []string{"AUDIO"},
			Markets:            []string{"de-DE"},
			HourlyRateUSD:      90.0,
			TurnaroundHours:    24,
			HistoricalAccuracy: models.UnmeasuredHistoricalAccuracy,
		},
	}

	t.Run("Approves compliant plan", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Hold German delivery and reassign to Berlin Synchron",
			Actions: []models.Action{
				{
					Type:     models.ActionHoldDelivery,
					TargetID: "del-germany",
					Reason:   "German audio track corrupt",
				},
				{
					Type:     models.ActionReassignVendor,
					TargetID: "vendor-berlin",
					Reason:   "Fast turnaround and high accuracy",
					Payload:  map[string]any{"package_id": "pkg-german-dub"},
				},
				{
					Type:     models.ActionEmailVendor,
					TargetID: "vendor-berlin",
					Reason:   "Dispatch expedited audio reconform",
				},
				{
					Type:     models.ActionNotifyStakeholders,
					TargetID: "slack-ops",
					Reason:   "Alert VP of Distribution to German hold",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionApproved {
			t.Fatalf("expected APPROVED, got: %s (rationale: %s)", verdict.Decision, verdict.Rationale)
		}
	})

	t.Run("Escalates to human operator when attempt reaches limit", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Attempt after multiple prior rejections",
			Actions: []models.Action{
				{
					Type:     models.ActionHoldDelivery,
					TargetID: "del-germany",
					Reason:   "Hold",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 3)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionEscalate {
			t.Fatalf("expected ESCALATE, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects empty action plan", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "No actions planned",
			Actions:   []models.Action{},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for empty actions, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects action with empty TargetID", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Malformed action missing target",
			Actions: []models.Action{
				{
					Type:     models.ActionHoldDelivery,
					TargetID: "",
					Reason:   "Missing target ID",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for empty target ID, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects contradictory plan targeting same delivery for HOLD and RELEASE", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Contradictory delivery instructions",
			Actions: []models.Action{
				{
					Type:     models.ActionHoldDelivery,
					TargetID: "del-germany",
					Reason:   "Hold defect",
				},
				{
					Type:     models.ActionReleaseDelivery,
					TargetID: "del-germany",
					Reason:   "Release defect",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for contradictory targets, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects market isolation violation", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Hold US delivery because of German audio issue",
			Actions: []models.Action{
				{
					Type:     models.ActionHoldDelivery,
					TargetID: "del-us",
					Reason:   "German audio issue",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for market isolation, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects releasing delivery when defective packages remain unresolved", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Premature release",
			Actions: []models.Action{
				{
					Type:     models.ActionReleaseDelivery,
					TargetID: "del-germany",
					Reason:   "Force release despite defective package",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for prerequisite failure, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects vendor reassignment missing package_id payload", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Reassign vendor without package context",
			Actions: []models.Action{
				{
					Type:     models.ActionReassignVendor,
					TargetID: "vendor-berlin",
					Reason:   "Missing package_id payload",
					Payload:  map[string]any{},
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for missing package_id, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects vendor reassignment for unrelated package not in blast radius", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Reassign unrelated Spanish package",
			Actions: []models.Action{
				{
					Type:     models.ActionReassignVendor,
					TargetID: "vendor-berlin",
					Reason:   "Reassign Spanish package when only German failed",
					Payload:  map[string]any{"package_id": "pkg-spanish-sub"},
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for unrelated package reassignment, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects vendor turnaround exceeding premiere deadline", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Reassign to slow vendor",
			Actions: []models.Action{
				{
					Type:     models.ActionReassignVendor,
					TargetID: "vendor-slow",
					Reason:   "Cheaper vendor",
					Payload:  map[string]any{"package_id": "pkg-german-dub"},
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for turnaround violation, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects vendor with low historical accuracy", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Reassign to sloppy vendor",
			Actions: []models.Action{
				{
					Type:     models.ActionReassignVendor,
					TargetID: "vendor-sloppy",
					Reason:   "Fastest vendor",
					Payload:  map[string]any{"package_id": "pkg-german-dub"},
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for quality floor violation, got: %s", verdict.Decision)
		}
	})

	t.Run("Approves vendor with unmeasured historical accuracy", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Reassign to new vendor without history",
			Actions: []models.Action{
				{
					Type:     models.ActionReassignVendor,
					TargetID: "vendor-new",
					Reason:   "New onboarding vendor with 24h turnaround",
					Payload:  map[string]any{"package_id": "pkg-german-dub"},
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionApproved {
			t.Fatalf("expected APPROVED for unmeasured vendor, got: %s (rationale: %s)", verdict.Decision, verdict.Rationale)
		}
	})

	t.Run("Rejects emailing unknown vendor", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Email unrecognized vendor",
			Actions: []models.Action{
				{
					Type:     models.ActionEmailVendor,
					TargetID: "vendor-ghost",
					Reason:   "Contact unknown third party",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for unknown vendor email, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects social update when premiere is distant", func(t *testing.T) {
		distantImpact := &models.DeliveryImpact{
			RootPackageID:      "pkg-german-dub",
			AffectedDeliveries: []string{"del-germany"},
			HoursUntilPremiere: 120.0,
			IsPremiereUrgent:   false,
		}

		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Post social delay notice 5 days early",
			Actions: []models.Action{
				{
					Type:     models.ActionHoldDelivery,
					TargetID: "del-germany",
					Reason:   "Hold delivery",
				},
				{
					Type:     models.ActionPostSocialUpdate,
					TargetID: "twitter",
					Reason:   "Tell fans German release is delayed",
				},
			},
		}

		res := agent.VerifyPlan(plan, distantImpact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for social guardrail violation, got: %s", verdict.Decision)
		}
	})

	t.Run("Rejects social update when no deliveries are on hold", func(t *testing.T) {
		plan := &models.ActionPlan{
			TitleSlug: "eclipse",
			Summary:   "Post public notice without any active holds",
			Actions: []models.Action{
				{
					Type:     models.ActionPostSocialUpdate,
					TargetID: "twitter",
					Reason:   "Post delay notice prematurely",
				},
			},
		}

		res := agent.VerifyPlan(plan, impact, vendors, 1)
		if res.IsErr() {
			t.Fatalf("VerifyPlan returned error: %v", res.Error())
		}
		verdict := res.Unwrap()
		if verdict.Decision != agent.DecisionRejected {
			t.Fatalf("expected REJECTED for unjustified social notice, got: %s", verdict.Decision)
		}
	})
}
