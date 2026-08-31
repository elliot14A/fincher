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

// SelectionDecision records the vendor allocation choice and rationale.
type SelectionDecision struct {
	WinnerVendorID   string  `json:"winner_vendor_id"`
	WinnerVendorName string  `json:"winner_vendor_name"`
	HourlyRateUSD    float64 `json:"hourly_rate_usd"`
	TurnaroundHours  int     `json:"turnaround_hours"`
	Rationale        string  `json:"rationale"`
}

// SelectVendor evaluates candidate vendors using Gemini Flash and selects the optimal partner.
func SelectVendor(
	ctx context.Context,
	m model.LLM,
	titleSlug string,
	component string,
	candidates []models.VendorCandidate,
	hoursUntilPremiere float64,
) domainerrors.Result[*SelectionDecision] {
	if m == nil {
		return domainerrors.Err[*SelectionDecision](NewError("agent.SelectVendor", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}
	if len(candidates) == 0 {
		return domainerrors.Err[*SelectionDecision](NewError("agent.SelectVendor", domainerrors.CodeInvalidInput, "candidates list cannot be empty", nil))
	}

	candidatesJSON, _ := json.Marshal(candidates)

	userPrompt := fmt.Sprintf(
		"Title: %s\nComponent: %s\nHours Until Premiere: %.1f\n\nCandidate Vendors:\n%s\n\nSelect the best vendor partner according to turnaround feasibility, quality floor (>= 90%%), and commercial rate.",
		titleSlug,
		component,
		hoursUntilPremiere,
		string(candidatesJSON),
	)

	return generateJSON[*SelectionDecision](ctx, m, "agent.SelectVendor", prompts.Selector, userPrompt)
}
