package models_test

import (
	"testing"
	"time"

	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestEvent_ValidateAndDefaults(t *testing.T) {
	e := &models.Event{
		Type: models.TypeVendorHeartbeat,
	}

	if err := e.Validate(); err != nil {
		t.Fatalf("expected validation to succeed, got: %v", err)
	}

	if e.ID == "" {
		t.Error("expected ID to be defaulted")
	}
	if e.Source != "fincher.system" {
		t.Errorf("expected source to default to fincher.system, got %q", e.Source)
	}
	if e.Subject != models.DefaultTitleAgnosticSentinel {
		t.Errorf("expected subject to default to %q, got %q", models.DefaultTitleAgnosticSentinel, e.Subject)
	}
	if e.DataContentType != "application/json" {
		t.Errorf("expected datacontenttype to default to application/json, got %q", e.DataContentType)
	}
	if e.Severity != models.SeverityInfo {
		t.Errorf("expected severity to default to INFO, got %q", e.Severity)
	}
	if e.Time.IsZero() {
		t.Error("expected time to default to now")
	}
}

func TestEvent_DataJSON(t *testing.T) {
	e := &models.Event{
		Type: models.TypeQCInspectionCompleted,
		Data: map[string]any{
			"status": "PASSED",
			"score":  99.5,
		},
	}

	jsonStr, err := e.DataJSON()
	if err != nil {
		t.Fatalf("expected DataJSON to succeed, got: %v", err)
	}
	if jsonStr == "" || jsonStr == "{}" {
		t.Fatalf("expected non-empty JSON, got %q", jsonStr)
	}
}

func TestEvent_Classify(t *testing.T) {
	tests := []struct {
		name     string
		event    *models.Event
		expected models.EventCategory
	}{
		{
			name: "vendor heartbeat is telemetry",
			event: &models.Event{
				Type: models.TypeVendorHeartbeat,
			},
			expected: models.CategoryTelemetry,
		},
		{
			name: "package download progress is telemetry",
			event: &models.Event{
				Type: models.TypePackageDownloadProgress,
			},
			expected: models.CategoryTelemetry,
		},
		{
			name: "qc inspection passed is routine outcome",
			event: &models.Event{
				Type: models.TypeQCInspectionCompleted,
				Data: map[string]any{"status": "PASSED"},
			},
			expected: models.CategoryRoutineOutcome,
		},
		{
			name: "qc inspection failed is incident",
			event: &models.Event{
				Type: models.TypeQCInspectionCompleted,
				Data: map[string]any{"status": "FAILED"},
			},
			expected: models.CategoryIncident,
		},
		{
			name: "audio sync drift detected is incident",
			event: &models.Event{
				Type: models.TypeAudioSyncDriftDetected,
			},
			expected: models.CategoryIncident,
		},
		{
			name: "master cut revised is incident",
			event: &models.Event{
				Type: models.TypeMasterCutRevised,
			},
			expected: models.CategoryIncident,
		},
		{
			name: "title created is allocation",
			event: &models.Event{
				Type: models.TypeTitleCreated,
			},
			expected: models.CategoryAllocation,
		},
		{
			name: "operator forced is operator forced",
			event: &models.Event{
				Type: models.TypeOperatorForced,
			},
			expected: models.CategoryOperatorForced,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.event.Time = time.Now().UTC()
			cat := tt.event.Classify()
			if cat != tt.expected {
				t.Fatalf("expected category %s, got %s", tt.expected, cat)
			}
		})
	}
}
