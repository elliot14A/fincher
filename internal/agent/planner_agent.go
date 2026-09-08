package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elliot14A/fincher/internal/agent/tools"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/prompts"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
)

func PlanRemediationViaTools(
	ctx context.Context,
	m model.LLM,
	tursoDB *sql.DB,
	event *models.Event,
	impact *models.DeliveryImpact,
	analytics *models.AnalyticsSummary,
	projection *tools.TitleProjection,
	feedback string,
) domainerrors.Result[*models.ActionPlan] {
	if m == nil {
		return domainerrors.Err[*models.ActionPlan](NewError("agent.PlanRemediationViaTools", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}
	if event == nil {
		return domainerrors.Err[*models.ActionPlan](NewError("agent.PlanRemediationViaTools", domainerrors.CodeInvalidInput, "event cannot be nil", nil))
	}
	if tursoDB == nil {
		return domainerrors.Err[*models.ActionPlan](NewError("agent.PlanRemediationViaTools", domainerrors.CodeInvalidInput, "turso db cannot be nil", nil))
	}

	queryTool, err := tools.NewTursoQueryTool(tursoDB)
	if err != nil {
		return domainerrors.Err[*models.ActionPlan](MapError("agent.PlanRemediationViaTools", "query_turso", err))
	}

	impactJSON, _ := json.Marshal(impact)
	analyticsJSON, _ := json.Marshal(analytics)
	projectionJSON, _ := json.Marshal(projection)
	eventDataJSON, _ := json.Marshal(event.Data)

	feedbackContext := ""
	if feedback != "" {
		feedbackContext = fmt.Sprintf("\nPREVIOUS PLAN REJECTED BY POLICY ENGINE:\nReason: %s\nYou MUST revise the actions to satisfy this policy constraint.", feedback)
	}

	affectedContext := ""
	if impact != nil && len(impact.AffectedPackages) > 0 {
		affectedContext = fmt.Sprintf(
			"\n\nAFFECTED PACKAGES REQUIRING REMEDIATION (%d): %s\nIf you proceed with repair (deadline permits), you MUST emit exactly one REASSIGN_VENDOR action per package id in this list — %d packages means %d REASSIGN_VENDOR actions, each with payload.package_id set to that package. Query the database for an eligible vendor for each package's component/market. If the deadline can NOT absorb re-deriving all of them, emit HOLD_TITLE instead.",
			len(impact.AffectedPackages), strings.Join(impact.AffectedPackages, ", "),
			len(impact.AffectedPackages), len(impact.AffectedPackages),
		)
	}

	userPrompt := fmt.Sprintf(
		"Event: %s (Type: %s, Severity: %s, Subject: %s)\nEvent Data: %s\n\nDelivery Impact: %s\n\nTitle Launch Projection: %s\n\nHistorical Analytics: %s%s%s\n\nUse the query_turso tool to find eligible vendors before choosing one. Then formulate a compliant remediation action plan.",
		event.ID, event.Type, event.Severity, event.Subject,
		string(eventDataJSON), string(impactJSON), string(projectionJSON), string(analyticsJSON), affectedContext, feedbackContext,
	)

	logger.Info("planner-agent: starting two-phase remediation planning",
		"event_type", event.Type, "subject", event.Subject, "has_feedback", feedback != "")

	emitInstruction := prompts.Planner + "\n\nYou are now emitting the final structured remediation plan. Do not call any tools. Return the plan as JSON conforming to the schema, including one action per required item based on the findings provided."

	clean, toolCalls, err := runToolThenSchema(
		ctx, m, "remediation", "plan-"+event.Subject,
		prompts.Planner, userPrompt, []tool.Tool{queryTool},
		emitInstruction, actionPlanSchema(),
	)
	if err != nil {
		return domainerrors.Err[*models.ActionPlan](MapError("agent.PlanRemediationViaTools", "run", err))
	}

	var plan models.ActionPlan
	if err := json.Unmarshal([]byte(clean), &plan); err != nil {
		return domainerrors.Err[*models.ActionPlan](NewError("agent.PlanRemediationViaTools", domainerrors.CodeInternal, fmt.Sprintf("failed to parse remediation plan: %s", clean), err))
	}
	if plan.TitleSlug == "" {
		plan.TitleSlug = event.Subject
	}

	logger.Info("planner-agent: completed",
		"event_type", event.Type, "subject", event.Subject, "actions", len(plan.Actions), "tool_calls", toolCalls)

	return domainerrors.Ok(&plan)
}

func actionPlanSchema() *genai.Schema {
	action := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"type":      {Type: genai.TypeString},
			"target_id": {Type: genai.TypeString},
			"reason":    {Type: genai.TypeString},
			"payload":   {Type: genai.TypeObject},
		},
		Required: []string{"type", "target_id", "reason"},
	}
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title_slug": {Type: genai.TypeString},
			"summary":    {Type: genai.TypeString},
			"actions":    {Type: genai.TypeArray, Items: action},
		},
		Required: []string{"summary", "actions"},
	}
}
