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

// FilterDecision records the event classification outcome.
type FilterDecision struct {
	Actionable  bool                 `json:"actionable"`
	Severity    models.EventSeverity `json:"severity"`
	AnomalyType string               `json:"anomaly_type"`
	Rationale   string               `json:"rationale"`
}

// IsCandidateTrigger determines if an incoming event type warrants LLM screening.
func IsCandidateTrigger(eventType string) bool {
	switch eventType {
	case models.TypeQCInspectionCompleted,
		models.TypeAudioSyncDriftDetected,
		models.TypeMasterCutRevised,
		models.TypeVendorSLABreach,
		models.TypePackageInvalidated,
		models.TypeOperatorForced,
		models.TypeInvestigationTriggered,
		models.TypeTitleDeadlineReached:
		return true
	default:
		return false
	}
}

// FilterEvent screens an incoming event using Gemini Flash to decide if an investigation is required.
func FilterEvent(ctx context.Context, m model.LLM, event *models.Event, hoursUntilPremiere float64) domainerrors.Result[*FilterDecision] {
	if m == nil {
		return domainerrors.Err[*FilterDecision](NewError("agent.FilterEvent", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}
	if event == nil {
		return domainerrors.Err[*FilterDecision](NewError("agent.FilterEvent", domainerrors.CodeInvalidInput, "event cannot be nil", nil))
	}

	eventDataJSON, _ := json.Marshal(event.Data)
	userPrompt := fmt.Sprintf(
		"Event ID: %s\nType: %s\nSeverity: %s\nSubject: %s\nHours Until Premiere: %.1f\nData: %s\n\nEvaluate whether this event requires investigation.",
		event.ID,
		event.Type,
		event.Severity,
		event.Subject,
		hoursUntilPremiere,
		string(eventDataJSON),
	)

	return generateJSON[*FilterDecision](ctx, m, "agent.FilterEvent", prompts.Filter, userPrompt)
}
