package agent

import (
	"fmt"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// VerificationDecision represents the gate evaluation outcome.
type VerificationDecision string

const (
	DecisionApproved VerificationDecision = "APPROVED"
	DecisionRejected VerificationDecision = "REJECTED"
	DecisionEscalate VerificationDecision = "ESCALATE"
)

const (
	MaxRemediationAttempts     = 3
	VendorAccuracyFloor        = 0.90
	SocialNoticeThresholdHours = 72.0
)

// VerificationResult captures the policy evaluation outcome and explanation.
type VerificationResult struct {
	Decision  VerificationDecision `json:"decision"`
	Rationale string               `json:"rationale"`
	Attempt   int                  `json:"attempt"`
}

// VerifyPlan rigorously validates a proposed ActionPlan against operational policies, SLA bounds, and market isolation gates.
func VerifyPlan(plan *models.ActionPlan, impact *models.DeliveryImpact, vendors []models.VendorCandidate, projection *models.TitleProjection, attempt int) domainerrors.Result[*VerificationResult] {
	if attempt >= MaxRemediationAttempts {
		return domainerrors.Ok(&VerificationResult{
			Decision:  DecisionEscalate,
			Rationale: fmt.Sprintf("Automated plan failed to pass verification after %d attempts; escalating to operations lead.", attempt),
			Attempt:   attempt,
		})
	}

	reject := func(format string, args ...any) domainerrors.Result[*VerificationResult] {
		return domainerrors.Ok(&VerificationResult{
			Decision:  DecisionRejected,
			Rationale: fmt.Sprintf(format, args...),
			Attempt:   attempt,
		})
	}

	if plan == nil || len(plan.Actions) == 0 {
		return reject("Action plan is empty; no operational remediation proposed.")
	}

	// Feasibility gate: If title readiness projection is breached (critical path + reconform exceeds remaining premiere window)
	// and the plan does not propose holding or alerting, reject it.
	if projection != nil && projection.IsBreached {
		hasMitigation := false
		for _, action := range plan.Actions {
			if action.Type == models.ActionHoldTitle || action.Type == models.ActionHoldDelivery || action.Type == models.ActionNotifyStakeholders {
				hasMitigation = true
				break
			}
		}
		if !hasMitigation {
			return reject("Feasibility breach: title critical path (%.1fh) exceeds remaining window to premiere (buffer: %.1fh); remediation plan must include HOLD_TITLE, HOLD_DELIVERY, or NOTIFY_STAKEHOLDERS.", projection.CriticalRemainingHours, projection.BufferHours)
		}
	}

	affectedDeliveryMap := make(map[string]bool)
	affectedPackageMap := make(map[string]bool)
	if impact != nil {
		for _, d := range impact.AffectedDeliveries {
			affectedDeliveryMap[d] = true
		}
		for _, p := range impact.AffectedPackages {
			affectedPackageMap[p] = true
		}
	}

	vendorMap := make(map[string]models.VendorCandidate)
	for _, v := range vendors {
		vendorMap[v.VendorID] = v
	}

	holdTargets := make(map[string]bool)
	releaseTargets := make(map[string]bool)

	for _, action := range plan.Actions {
		if action.TargetID == "" {
			return reject("Malformed action: action type %s has an empty target identifier.", action.Type)
		}

		switch action.Type {
		case models.ActionHoldDelivery:
			if releaseTargets[action.TargetID] {
				return reject("Contradictory plan: delivery %s is targeted for both HOLD and RELEASE within the same action plan.", action.TargetID)
			}
			holdTargets[action.TargetID] = true

			if impact != nil && len(affectedDeliveryMap) > 0 && !affectedDeliveryMap[action.TargetID] {
				return reject("Market isolation violation: delivery %s is not affected by failing asset (affected deliveries: %v).", action.TargetID, impact.AffectedDeliveries)
			}

		case models.ActionReleaseDelivery:
			if holdTargets[action.TargetID] {
				return reject("Contradictory plan: delivery %s is targeted for both HOLD and RELEASE within the same action plan.", action.TargetID)
			}
			releaseTargets[action.TargetID] = true

			if impact != nil && len(impact.AffectedPackages) > 0 {
				return reject("Prerequisite check failed: cannot release delivery %s while defective packages %v remain unresolved.", action.TargetID, impact.AffectedPackages)
			}

		case models.ActionReassignVendor:
			candidate, exists := vendorMap[action.TargetID]
			if !exists {
				return reject("Vendor validation failure: candidate vendor %s does not exist or is inactive.", action.TargetID)
			}

			pkgIDVal, hasPkg := action.Payload["package_id"]
			pkgID, isStr := pkgIDVal.(string)
			if !hasPkg || !isStr || pkgID == "" {
				return reject("Payload validation failure: reassigning vendor %s requires a valid 'package_id' in action payload.", action.TargetID)
			}

			if impact != nil && len(affectedPackageMap) > 0 && !affectedPackageMap[pkgID] {
				return reject("Blast radius violation: target package %s is not among affected defective packages (%v).", pkgID, impact.AffectedPackages)
			}

			if impact != nil && impact.HoursUntilPremiere > 0 && float64(candidate.TurnaroundHours) > impact.HoursUntilPremiere {
				return reject("Turnaround violation: vendor %s turnaround (%dh) exceeds hours until premiere (%.1fh).", candidate.VendorName, candidate.TurnaroundHours, impact.HoursUntilPremiere)
			}

			if candidate.HistoricalAccuracy >= 0 && candidate.HistoricalAccuracy < VendorAccuracyFloor {
				return reject("Quality floor violation: vendor %s historical accuracy (%.1f%%) is below required %.0f%% floor.", candidate.VendorName, candidate.HistoricalAccuracy*100, VendorAccuracyFloor*100)
			}

		case models.ActionEmailVendor:
			if _, exists := vendorMap[action.TargetID]; !exists {
				return reject("Vendor validation failure: cannot dispatch email to unknown vendor %s.", action.TargetID)
			}

		case models.ActionPostSocialUpdate:
			if impact != nil && !impact.IsPremiereUrgent && impact.HoursUntilPremiere > SocialNoticeThresholdHours {
				return reject("Social guardrail violation: premiere is %.1fh away (> %.0fh threshold); public delay notice prohibited while internal remediation is possible.", impact.HoursUntilPremiere, SocialNoticeThresholdHours)
			}
		}
	}

	for _, action := range plan.Actions {
		if action.Type == models.ActionPostSocialUpdate {
			if len(holdTargets) == 0 {
				return reject("Social notice unjustified: cannot post public delay announcement when no delivery holds are active.")
			}
		}
	}

	return domainerrors.Ok(&VerificationResult{
		Decision:  DecisionApproved,
		Rationale: "All proposed actions comply with market isolation, turnaround deadlines, prerequisite gates, and operational boundaries.",
		Attempt:   attempt,
	})
}
