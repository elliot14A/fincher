package events

import (
	"time"

	"github.com/elliot14A/fincher/pkg/domain/models"
)

// QCEventParams contains all input parameters required to construct a full-fidelity QC CloudEvent.
type QCEventParams struct {
	ID                   string
	Time                 time.Time
	TitleSlug            string
	InspectorAgent       string
	PackageID            string
	VendorID             string
	Component            string
	Language             string
	Status               string
	SyncDriftMS          *float64
	VideoCorruptionScore *float64
	DefectCategory       string
	Severity             models.EventSeverity
}

// BuildQCEvent constructs a models.Event conforming to the fincher.qc and fincher.vendor_metrics schema.
func BuildQCEvent(params QCEventParams) models.Event {
	data := map[string]any{
		"package_id":      params.PackageID,
		"vendor_id":       params.VendorID,
		"component":       params.Component,
		"language":        params.Language,
		"status":          params.Status,
		"defect_category": params.DefectCategory,
	}

	if params.SyncDriftMS != nil {
		data["sync_drift_ms"] = *params.SyncDriftMS
	}
	if params.VideoCorruptionScore != nil {
		data["video_corruption_score"] = *params.VideoCorruptionScore
	}

	subject := params.TitleSlug
	if subject == "" {
		subject = "GLOBAL"
	}

	source := params.InspectorAgent
	if source == "" {
		source = "qc.automated-inspector"
	}

	severity := params.Severity
	if severity == "" {
		switch params.Status {
		case "FAILED", "WARNING":
			severity = models.SeverityWarn
		default:
			severity = models.SeverityInfo
		}
	}

	return models.Event{
		ID:              params.ID,
		Type:            models.TypeQCInspectionCompleted,
		Source:          source,
		Subject:         subject,
		Time:            params.Time,
		Severity:        severity,
		DataContentType: "application/json",
		Data:            data,
	}
}
