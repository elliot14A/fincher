package models_test

import (
	"testing"

	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestUnmarshalPayload_QC(t *testing.T) {
	ev := &models.Event{
		ID:   "evt-1",
		Type: models.TypeQCInspectionCompleted,
		Data: map[string]any{
			"package_id":      "pkg-audio-de",
			"vendor_id":       "vendor-berlin",
			"component":       "AUDIO",
			"language":        "de",
			"status":          "FAILED",
			"sync_drift_ms":   145.5,
			"defect_category": "SYNC_DRIFT",
		},
	}

	payload, err := models.UnmarshalPayload[models.QCPayload](ev)
	if err != nil {
		t.Fatalf("failed to unmarshal QC payload: %v", err)
	}

	if payload.PackageID != "pkg-audio-de" {
		t.Errorf("expected package_id pkg-audio-de, got: %s", payload.PackageID)
	}
	if payload.VendorID != "vendor-berlin" {
		t.Errorf("expected vendor_id vendor-berlin, got: %s", payload.VendorID)
	}
	if payload.Status != "FAILED" {
		t.Errorf("expected status FAILED, got: %s", payload.Status)
	}
	if payload.SyncDriftMs != 145.5 {
		t.Errorf("expected drift 145.5, got: %f", payload.SyncDriftMs)
	}
}

func TestUnmarshalPayload_MasterCut(t *testing.T) {
	ev := &models.Event{
		ID:   "evt-2",
		Type: models.TypeMasterCutRevised,
		Data: map[string]any{
			"master_id": "master-v2",
			"version":   2,
		},
	}

	payload, err := models.UnmarshalPayload[models.MasterCutPayload](ev)
	if err != nil {
		t.Fatalf("failed to unmarshal MasterCut payload: %v", err)
	}

	if payload.MasterID != "master-v2" {
		t.Errorf("expected master_id master-v2, got: %s", payload.MasterID)
	}
	if payload.Version != 2 {
		t.Errorf("expected version 2, got: %d", payload.Version)
	}
}
