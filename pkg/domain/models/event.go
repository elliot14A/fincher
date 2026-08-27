package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// EventSeverity denotes the urgency level of an ingested event.
type EventSeverity string

const (
	SeverityInfo     EventSeverity = "INFO"
	SeverityWarn     EventSeverity = "WARN"
	SeverityCritical EventSeverity = "CRITICAL"
)

// Canonical event types across Fincher.
const (
	TypeVendorHeartbeat           = "fincher.vendor.heartbeat"
	TypePackageDownloadStarted    = "fincher.package.download.started"
	TypePackageDownloadProgress   = "fincher.package.download.progress"
	TypeQCInspectionCompleted     = "fincher.qc.completed"
	TypeAudioSyncDriftDetected    = "fincher.audio.sync_drift"
	TypeMasterCutRevised          = "fincher.master.cut.revised"
	TypeVendorSLABreach           = "fincher.vendor.sla_breach"
	TypePackageInvalidated        = "fincher.package.invalidated"
	TypeTitleCreated              = "fincher.title.created"
	TypePackageRequired           = "fincher.package.required"
	TypeVendorReconformDispatched = "fincher.vendor.reconform.dispatched"
	TypeOperatorForced            = "fincher.operator.forced"
	TypeInvestigationTriggered    = "fincher.investigation.triggered"
)

// EventCategory classifies events for downstream routing.
type EventCategory string

const (
	CategoryTelemetry         EventCategory = "TELEMETRY"
	CategoryRoutineOutcome    EventCategory = "ROUTINE_OUTCOME"
	CategoryAnomalySignal     EventCategory = "ANOMALY_SIGNAL"
	CategoryAllocationRequest EventCategory = "ALLOCATION_REQUEST"
	CategoryOperatorForced    EventCategory = "OPERATOR_FORCED"
)

// DefaultTitleAgnosticSentinel is used when an event is not scoped to a specific title.
const DefaultTitleAgnosticSentinel = "GLOBAL"

// Event represents an immutable CloudEvent.
type Event struct {
	ID              string         `json:"id"`
	Type            string         `json:"type" validate:"required,min=1"`
	Source          string         `json:"source" validate:"required,min=1"`
	Subject         string         `json:"subject"`
	Time            time.Time      `json:"time"`
	Data            map[string]any `json:"data"`
	DataContentType string         `json:"datacontenttype"`
	Severity        EventSeverity  `json:"severity" validate:"required,oneof=INFO WARN CRITICAL"`
	CreatedAt       time.Time      `json:"created_at"`
}

// EventBatchResponse represents the response returned after event batch ingestion.
type EventBatchResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// Validate ensures required event attributes are present and defaults are populated.
func (e *Event) Validate() error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	if e.Source == "" {
		e.Source = "fincher.system"
	}
	if e.Subject == "" {
		e.Subject = DefaultTitleAgnosticSentinel
	}
	if e.DataContentType == "" {
		e.DataContentType = "application/json"
	}
	if e.Severity == "" {
		e.Severity = SeverityInfo
	}
	return validate.Struct(e)
}

// DataJSON serializes the data payload into a JSON string.
func (e *Event) DataJSON() (string, error) {
	if e.Data == nil {
		return "{}", nil
	}
	b, err := json.Marshal(e.Data)
	if err != nil {
		return "", fmt.Errorf("marshaling event data: %w", err)
	}
	return string(b), nil
}

// Classify determines the EventCategory for downstream routing.
func (e *Event) Classify() EventCategory {
	switch strings.ToLower(e.Type) {
	case TypeVendorHeartbeat, TypePackageDownloadStarted, TypePackageDownloadProgress:
		return CategoryTelemetry

	case TypeQCInspectionCompleted:
		if e.Data != nil {
			if status, ok := e.Data["status"].(string); ok {
				upperStatus := strings.ToUpper(status)
				if upperStatus == "FAILED" || upperStatus == "WARNING" {
					return CategoryAnomalySignal
				}
				if upperStatus == "PASSED" {
					return CategoryRoutineOutcome
				}
			}
		}
		if e.Severity == SeverityCritical || e.Severity == SeverityWarn {
			return CategoryAnomalySignal
		}
		return CategoryRoutineOutcome

	case TypeAudioSyncDriftDetected, TypeMasterCutRevised, TypeVendorSLABreach, TypePackageInvalidated:
		return CategoryAnomalySignal

	case TypeTitleCreated, TypePackageRequired, TypeVendorReconformDispatched:
		return CategoryAllocationRequest

	case TypeOperatorForced, TypeInvestigationTriggered:
		return CategoryOperatorForced

	default:
		if e.Severity == SeverityCritical || e.Severity == SeverityWarn {
			return CategoryAnomalySignal
		}
		return CategoryTelemetry
	}
}
