package agent

import (
	"context"
	"encoding/json"
	"fmt"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/prompts"
	"google.golang.org/adk/v2/model"
)

// PlanRemediation synthesizes an operational ActionPlan using Gemini Flash based on impact, analytics, and available vendors.
func PlanRemediation(
	ctx context.Context,
	m model.LLM,
	event *models.Event,
	impact *models.DeliveryImpact,
	analytics *models.AnalyticsSummary,
	candidates []models.VendorCandidate,
	feedback string,
) domainerrors.Result[*models.ActionPlan] {
	if m == nil {
		return domainerrors.Err[*models.ActionPlan](NewError("agent.PlanRemediation", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}
	if event == nil {
		return domainerrors.Err[*models.ActionPlan](NewError("agent.PlanRemediation", domainerrors.CodeInvalidInput, "event cannot be nil", nil))
	}

	impactJSON, _ := json.Marshal(impact)
	analyticsJSON, _ := json.Marshal(analytics)
	candidatesJSON, _ := json.Marshal(candidates)
	eventDataJSON, _ := json.Marshal(event.Data)

	feedbackContext := ""
	if feedback != "" {
		feedbackContext = fmt.Sprintf("\nPREVIOUS PLAN REJECTED BY POLICY ENGINE:\nReason: %s\nYou MUST revise the actions to satisfy this policy constraint.", feedback)
	}

	userPrompt := fmt.Sprintf(
		"Event: %s (Type: %s, Severity: %s, Subject: %s)\nEvent Data: %s\n\nDelivery Impact: %s\n\nHistorical Analytics: %s\n\nVendor Candidates: %s%s\n\nFormulate a compliant remediation action plan.",
		event.ID,
		event.Type,
		event.Severity,
		event.Subject,
		string(eventDataJSON),
		string(impactJSON),
		string(analyticsJSON),
		string(candidatesJSON),
		feedbackContext,
	)

	planRes := generateJSON[*models.ActionPlan](ctx, m, "agent.PlanRemediation", prompts.Planner, userPrompt)
	if planRes.IsErr() {
		return planRes
	}

	plan := planRes.Unwrap()
	if plan.TitleSlug == "" {
		plan.TitleSlug = event.Subject
	}

	return domainerrors.Ok(plan)
}
