package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/elliot14A/fincher/internal/agent/tools"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/prompts"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
)

func SelectVendorsViaTools(
	ctx context.Context,
	m model.LLM,
	tursoDB *sql.DB,
	titleSlug string,
	requirements []models.AllocationRequirement,
	hoursUntilPremiere float64,
) domainerrors.Result[*models.AllocationPlan] {
	if m == nil {
		return domainerrors.Err[*models.AllocationPlan](NewError("agent.SelectVendorsViaTools", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}
	if tursoDB == nil {
		return domainerrors.Err[*models.AllocationPlan](NewError("agent.SelectVendorsViaTools", domainerrors.CodeInvalidInput, "turso db cannot be nil", nil))
	}
	if len(requirements) == 0 {
		return domainerrors.Err[*models.AllocationPlan](NewError("agent.SelectVendorsViaTools", domainerrors.CodeInvalidInput, "requirements cannot be empty", nil))
	}

	queryTool, err := tools.NewTursoQueryTool(tursoDB)
	if err != nil {
		return domainerrors.Err[*models.AllocationPlan](MapError("agent.SelectVendorsViaTools", "query_turso", err))
	}

	reqJSON, _ := json.MarshalIndent(requirements, "", "  ")
	userPrompt := fmt.Sprintf(
		"Title: %s\nHours Until Premiere: %.1f\n\nRequirements needing a vendor (component/market/language):\n%s\n\nFor each requirement, query the vendor database for vendors that actually handle that component and market, then choose the best one (quality first, then cost, then turnaround). Only choose vendors your queries returned.",
		titleSlug, hoursUntilPremiere, string(reqJSON),
	)

	logger.Info("allocation-agent: starting two-phase vendor selection",
		"title_slug", titleSlug, "requirements", len(requirements), "hours_until_premiere", hoursUntilPremiere)

	emitInstruction := prompts.PlanSelector + "\n\nYou are now emitting the final structured allocation plan. Do not call any tools. Return one assignment per requirement based strictly on the vendor findings provided."

	clean, toolCalls, err := runToolThenSchema(
		ctx, m, "allocation", "alloc-"+titleSlug,
		prompts.PlanSelector, userPrompt, []tool.Tool{queryTool},
		emitInstruction, allocationPlanSchema(),
	)
	if err != nil {
		return domainerrors.Err[*models.AllocationPlan](MapError("agent.SelectVendorsViaTools", "run", err))
	}

	var plan models.AllocationPlan
	if err := json.Unmarshal([]byte(clean), &plan); err != nil {
		return domainerrors.Err[*models.AllocationPlan](NewError("agent.SelectVendorsViaTools", domainerrors.CodeInternal, fmt.Sprintf("failed to parse allocation plan: %s", clean), err))
	}

	logger.Info("allocation-agent: completed",
		"title_slug", titleSlug, "assignments", len(plan.Assignments), "tool_calls", toolCalls)

	return domainerrors.Ok(&plan)
}

func allocationPlanSchema() *genai.Schema {
	assignment := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"component":          {Type: genai.TypeString},
			"market":             {Type: genai.TypeString},
			"language":           {Type: genai.TypeString},
			"winner_vendor_id":   {Type: genai.TypeString},
			"winner_vendor_name": {Type: genai.TypeString},
			"rationale":          {Type: genai.TypeString},
			"hourly_rate_usd":    {Type: genai.TypeNumber},
			"turnaround_hours":   {Type: genai.TypeInteger},
		},
		Required: []string{"component", "market", "winner_vendor_id", "rationale"},
	}
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"assignments":     {Type: genai.TypeArray, Items: assignment},
			"overall_summary": {Type: genai.TypeString},
		},
		Required: []string{"assignments", "overall_summary"},
	}
}
