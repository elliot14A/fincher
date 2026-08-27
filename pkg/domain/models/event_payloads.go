package models

import (
	"encoding/json"
	"fmt"
)

// QCPayload represents the data payload for fincher.qc.completed events.
type QCPayload struct {
	PackageID            string  `json:"package_id"`
	VendorID             string  `json:"vendor_id"`
	Component            string  `json:"component"`
	Language             string  `json:"language"`
	Status               string  `json:"status"`
	SyncDriftMs          float64 `json:"sync_drift_ms,omitempty"`
	VideoCorruptionScore float64 `json:"video_corruption_score,omitempty"`
	DefectCategory       string  `json:"defect_category,omitempty"`
}

// MasterCutPayload represents the data payload for fincher.master.cut.revised events.
type MasterCutPayload struct {
	MasterID string `json:"master_id"`
	Version  int    `json:"version"`
}

// SLABreachPayload represents the data payload for fincher.vendor.sla_breach events.
type SLABreachPayload struct {
	VendorID  string `json:"vendor_id"`
	PackageID string `json:"package_id,omitempty"`
}

// PackageInvalidatedPayload represents the data payload for fincher.package.invalidated events.
type PackageInvalidatedPayload struct {
	PackageID string `json:"package_id"`
	Reason    string `json:"reason"`
}

// UnmarshalPayload decodes an Event's Data into a typed payload struct.
func UnmarshalPayload[T any](e *Event) (T, error) {
	var out T
	if e == nil || e.Data == nil {
		return out, nil
	}

	b, err := json.Marshal(e.Data)
	if err != nil {
		return out, fmt.Errorf("marshaling event %s data: %w", e.ID, err)
	}

	if err := json.Unmarshal(b, &out); err != nil {
		return out, fmt.Errorf("decoding event %s payload as %T: %w", e.ID, out, err)
	}

	return out, nil
}
