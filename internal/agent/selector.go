package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/prompts"
	"google.golang.org/adk/v2/model"
)

// bestCandidate returns the preferred fallback vendor from a pool when the model's
// assignment must be substituted (dropped requirement or ineligible/hallucinated pick).
// Preference order: highest measured historical accuracy, then lowest hourly rate,
// then lowest turnaround. Unmeasured accuracy (sentinel -1.0) is treated as neutral,
// so an unmeasured vendor never outranks a proven one purely on the sentinel.
func bestCandidate(cands []models.VendorCandidate) models.VendorCandidate {
	best := cands[0]
	accOf := func(c models.VendorCandidate) float64 {
		if c.HistoricalAccuracy == models.UnmeasuredHistoricalAccuracy {
			return 0 // neutral: below any measured pass rate, above nothing
		}
		return c.HistoricalAccuracy
	}
	for _, c := range cands[1:] {
		ca, ba := accOf(c), accOf(best)
		switch {
		case ca != ba:
			if ca > ba {
				best = c
			}
		case c.HourlyRateUSD != best.HourlyRateUSD:
			if c.HourlyRateUSD < best.HourlyRateUSD {
				best = c
			}
		default:
			if c.TurnaroundHours < best.TurnaroundHours {
				best = c
			}
		}
	}
	return best
}

// SelectionDecision records the vendor allocation choice and rationale for single component selection.
type SelectionDecision struct {
	WinnerVendorID   string  `json:"winner_vendor_id"`
	WinnerVendorName string  `json:"winner_vendor_name"`
	HourlyRateUSD    float64 `json:"hourly_rate_usd"`
	TurnaroundHours  int     `json:"turnaround_hours"`
	Rationale        string  `json:"rationale"`
}

// SelectVendor evaluates candidate vendors for a single component and selects the optimal partner.
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

// SelectVendorsForPlan reasons holistically across multiple title requirements in a single LLM call.
func SelectVendorsForPlan(
	ctx context.Context,
	m model.LLM,
	titleSlug string,
	requirements []models.AllocationRequirement,
	candidatesByRequirement map[string][]models.VendorCandidate,
	hoursUntilPremiere float64,
) domainerrors.Result[*models.AllocationPlan] {
	if m == nil {
		return domainerrors.Err[*models.AllocationPlan](NewError("agent.SelectVendorsForPlan", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}
	if len(requirements) == 0 {
		return domainerrors.Err[*models.AllocationPlan](NewError("agent.SelectVendorsForPlan", domainerrors.CodeInvalidInput, "requirements list cannot be empty", nil))
	}

	type reqCandidatePool struct {
		Component  string                   `json:"component"`
		Market     string                   `json:"market"`
		Language   string                   `json:"language"`
		Candidates []models.VendorCandidate `json:"candidates"`
	}

	pools := make([]reqCandidatePool, 0, len(requirements))
	for _, req := range requirements {
		key := fmt.Sprintf("%s|%s", strings.ToUpper(req.Component), req.Market)
		cands := candidatesByRequirement[key]
		if cands == nil {
			cands = []models.VendorCandidate{}
		}
		pools = append(pools, reqCandidatePool{
			Component:  req.Component,
			Market:     req.Market,
			Language:   req.Language,
			Candidates: cands,
		})
	}

	poolsJSON, _ := json.MarshalIndent(pools, "", "  ")

	userPrompt := fmt.Sprintf(
		"Title: %s\nHours Until Premiere: %.1f\n\nRequirements and Verified Candidate Pools:\n%s\n\nCreate a complete, balanced staffing allocation plan across all requirements.",
		titleSlug,
		hoursUntilPremiere,
		string(poolsJSON),
	)

	planRes := generateJSON[*models.AllocationPlan](ctx, m, "agent.SelectVendorsForPlan", prompts.PlanSelector, userPrompt)
	if planRes.IsErr() {
		return planRes
	}

	plan := planRes.Unwrap()
	if plan == nil {
		plan = &models.AllocationPlan{Assignments: []models.RequirementAssignment{}}
	}

	// Index model assignments by (Component, Market)
	assignmentMap := make(map[string]models.RequirementAssignment)
	for _, a := range plan.Assignments {
		key := fmt.Sprintf("%s|%s", strings.ToUpper(strings.TrimSpace(a.Component)), strings.TrimSpace(a.Market))
		assignmentMap[key] = a
	}

	verifiedAssignments := make([]models.RequirementAssignment, 0, len(requirements))
	for _, req := range requirements {
		reqComp := strings.ToUpper(strings.TrimSpace(req.Component))
		reqMarket := strings.TrimSpace(req.Market)
		key := fmt.Sprintf("%s|%s", reqComp, reqMarket)

		cands := candidatesByRequirement[key]
		candMap := make(map[string]models.VendorCandidate)
		for _, c := range cands {
			candMap[c.VendorID] = c
		}

		assignment, found := assignmentMap[key]
		if !found || assignment.WinnerVendorID == "" {
			if len(cands) == 0 {
				logger.Warn("selector: requirement missing from model plan with no eligible vendor",
					"title_slug", titleSlug, "component", req.Component, "market", req.Market)
				assignment = models.RequirementAssignment{
					Component:        req.Component,
					Market:           req.Market,
					Language:         req.Language,
					WinnerVendorID:   "no_eligible_vendor",
					WinnerVendorName: "None",
					HourlyRateUSD:    0,
					TurnaroundHours:  0,
					Rationale:        "no eligible vendor in pool for component and market",
				}
			} else {
				top := bestCandidate(cands)
				logger.Warn("selector: requirement dropped by model, recovered top candidate defensively",
					"title_slug", titleSlug, "component", req.Component, "market", req.Market, "recovered_vendor", top.VendorID)
				assignment = models.RequirementAssignment{
					Component:        req.Component,
					Market:           req.Market,
					Language:         req.Language,
					WinnerVendorID:   top.VendorID,
					WinnerVendorName: top.VendorName,
					HourlyRateUSD:    top.HourlyRateUSD,
					TurnaroundHours:  top.TurnaroundHours,
					Rationale:        "Recovered top eligible candidate (requirement omitted from model plan)",
				}
			}
		} else if len(cands) == 0 {
			logger.Warn("selector: model assigned a vendor for a requirement with no eligible pool, overriding to no_eligible_vendor",
				"title_slug", titleSlug, "component", req.Component, "market", req.Market, "model_vendor", assignment.WinnerVendorID)
			assignment.WinnerVendorID = "no_eligible_vendor"
			assignment.WinnerVendorName = "None"
			assignment.HourlyRateUSD = 0
			assignment.TurnaroundHours = 0
			if assignment.Rationale == "" {
				assignment.Rationale = "no eligible vendor in pool for component and market"
			}
		} else if _, inPool := candMap[assignment.WinnerVendorID]; !inPool && assignment.WinnerVendorID != "no_eligible_vendor" {
			top := bestCandidate(cands)
			logger.Warn("selector: model selected an ineligible vendor, overriding with top eligible candidate",
				"title_slug", titleSlug, "component", req.Component, "market", req.Market,
				"model_vendor", assignment.WinnerVendorID, "override_vendor", top.VendorID)
			assignment.WinnerVendorID = top.VendorID
			assignment.WinnerVendorName = top.VendorName
			assignment.HourlyRateUSD = top.HourlyRateUSD
			assignment.TurnaroundHours = top.TurnaroundHours
		}

		assignment.Component = req.Component
		assignment.Market = req.Market
		if assignment.Language == "" {
			assignment.Language = req.Language
		}
		verifiedAssignments = append(verifiedAssignments, assignment)
	}

	plan.Assignments = verifiedAssignments
	if plan.OverallSummary == "" {
		plan.OverallSummary = fmt.Sprintf("Completed allocation plan for %d requirements.", len(verifiedAssignments))
	}

	return domainerrors.Ok(plan)
}
