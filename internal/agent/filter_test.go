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

func TestFilterEvent(t *testing.T) {
	ctx := context.Background()

	t.Run("Fails when model is nil", func(t *testing.T) {
		event := &models.Event{
			ID:       "evt-1",
			Type:     models.TypeQCInspectionCompleted,
			Severity: models.SeverityCritical,
			Time:     time.Now().UTC(),
		}
		res := agent.FilterEvent(ctx, nil, event, 48.0)
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
		res := agent.FilterEvent(ctx, llm, nil, 48.0)
		if res.IsOk() {
			t.Fatal("expected error for nil event, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeInvalidInput {
			t.Errorf("expected INVALID_INPUT code, got: %s", domErr.Code)
		}
	})

	t.Run("Screens actionable audio drift anomaly", func(t *testing.T) {
		jsonOutput := `{
			"actionable": true,
			"severity": "CRITICAL",
			"anomaly_type": "AUDIO_SYNC_DRIFT",
			"rationale": "Audio sync drift of 145ms exceeds 50ms broadcast tolerance."
		}`
		llm := &mockLLM{response: jsonOutput}

		event := &models.Event{
			ID:       "evt-defect-1",
			Type:     models.TypeQCInspectionCompleted,
			Severity: models.SeverityCritical,
			Subject:  "eclipse",
			Time:     time.Now().UTC(),
			Data: map[string]any{
				"package_id":    "pkg-german-dub",
				"status":        "FAILED",
				"sync_drift_ms": 145.0,
			},
		}

		res := agent.FilterEvent(ctx, llm, event, 24.0)
		if res.IsErr() {
			t.Fatalf("FilterEvent returned error: %v", res.Error())
		}
		decision := res.Unwrap()

		if !decision.Actionable {
			t.Errorf("expected actionable true, got false")
		}
		if decision.Severity != models.SeverityCritical {
			t.Errorf("expected CRITICAL severity, got: %s", decision.Severity)
		}
		if decision.AnomalyType != "AUDIO_SYNC_DRIFT" {
			t.Errorf("expected AUDIO_SYNC_DRIFT, got: %s", decision.AnomalyType)
		}

		if llm.lastRequest == nil || llm.lastRequest.Config == nil {
			t.Fatal("expected LLMRequest with non-nil Config")
		}
		if llm.lastRequest.Config.ResponseMIMEType != "application/json" {
			t.Errorf("expected application/json MIME type, got: %s", llm.lastRequest.Config.ResponseMIMEType)
		}
		if llm.lastRequest.Config.ResponseJsonSchema == nil {
			t.Error("expected enforced ResponseJsonSchema to be set on LLM request")
		}
	})

	t.Run("Screens benign routine QC pass as non-actionable", func(t *testing.T) {
		jsonOutput := `{
			"actionable": false,
			"severity": "INFO",
			"anomaly_type": "NONE",
			"rationale": "Routine inspection completed within acceptable parameters."
		}`
		llm := &mockLLM{response: jsonOutput}

		event := &models.Event{
			ID:       "evt-pass-1",
			Type:     models.TypeQCInspectionCompleted,
			Severity: models.SeverityInfo,
			Subject:  "eclipse",
			Time:     time.Now().UTC(),
			Data: map[string]any{
				"package_id":    "pkg-english-master",
				"status":        "PASSED",
				"sync_drift_ms": 1.2,
			},
		}

		res := agent.FilterEvent(ctx, llm, event, 72.0)
		if res.IsErr() {
			t.Fatalf("FilterEvent returned error: %v", res.Error())
		}
		decision := res.Unwrap()

		if decision.Actionable {
			t.Errorf("expected actionable false for benign event, got true")
		}
		if decision.Severity != models.SeverityInfo {
			t.Errorf("expected INFO severity, got: %s", decision.Severity)
		}
	})

	t.Run("Fails when model returns malformed JSON", func(t *testing.T) {
		llm := &mockLLM{response: "not a valid json response"}
		event := &models.Event{
			ID:       "evt-1",
			Type:     models.TypeQCInspectionCompleted,
			Severity: models.SeverityCritical,
			Time:     time.Now().UTC(),
		}

		res := agent.FilterEvent(ctx, llm, event, 48.0)
		if res.IsOk() {
			t.Fatal("expected error for malformed model JSON, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeInternal {
			t.Errorf("expected INTERNAL code, got: %s", domErr.Code)
		}
	})

	t.Run("Propagates model generation error", func(t *testing.T) {
		llm := &mockLLM{err: errors.New("upstream quota exceeded")}
		event := &models.Event{
			ID:       "evt-1",
			Type:     models.TypeQCInspectionCompleted,
			Severity: models.SeverityCritical,
			Time:     time.Now().UTC(),
		}

		res := agent.FilterEvent(ctx, llm, event, 48.0)
		if res.IsOk() {
			t.Fatal("expected error when model fails, got Ok")
		}
		domErr := res.Error().(*domainerrors.DomainError)
		if domErr.Code != domainerrors.CodeBudgetExceeded {
			t.Errorf("expected BUDGET_EXCEEDED for quota error, got: %s", domErr.Code)
		}
	})
}

func TestIsCandidateTrigger(t *testing.T) {
	candidates := []string{
		models.TypeQCInspectionCompleted,
		models.TypeAudioSyncDriftDetected,
		models.TypeMasterCutRevised,
		models.TypeVendorSLABreach,
		models.TypePackageInvalidated,
		models.TypeOperatorForced,
		models.TypeInvestigationTriggered,
		models.TypeTitleDeadlineReached,
	}

	for _, c := range candidates {
		if !agent.IsCandidateTrigger(c) {
			t.Errorf("expected %s to be a candidate trigger, got false", c)
		}
	}

	routine := []string{
		models.TypePackageDownloadStarted,
		models.TypePackageDownloadProgress,
		models.TypeVendorHeartbeat,
		"fincher.delivery.held",
		"fincher.vendor.assigned",
	}

	for _, r := range routine {
		if agent.IsCandidateTrigger(r) {
			t.Errorf("expected %s to NOT be a candidate trigger, got true", r)
		}
	}
}
