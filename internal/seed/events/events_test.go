package events_test

import (
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/seed"
	"github.com/elliot14A/fincher/internal/seed/entities"
	"github.com/elliot14A/fincher/internal/seed/events"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestBuildQCEvent_FullFidelity(t *testing.T) {
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	drift := 115.5

	event := events.BuildQCEvent(events.QCEventParams{
		ID:             "evt-test-1",
		Time:           now,
		TitleSlug:      "avatar-fire-ash",
		InspectorAgent: "qc.deluxe",
		PackageID:      "pkg-audio-de",
		VendorID:       "vnd-deluxe",
		Component:      "AUDIO",
		Language:       "de-DE",
		Status:         "FAILED",
		SyncDriftMS:    &drift,
		DefectCategory: "AUDIO_SYNC_DRIFT",
		Severity:       models.SeverityWarn,
	})

	if event.ID != "evt-test-1" {
		t.Errorf("expected ID 'evt-test-1', got '%s'", event.ID)
	}
	if event.Type != models.TypeQCInspectionCompleted {
		t.Errorf("expected Type '%s', got '%s'", models.TypeQCInspectionCompleted, event.Type)
	}
	if event.Subject != "avatar-fire-ash" {
		t.Errorf("expected Subject 'avatar-fire-ash', got '%s'", event.Subject)
	}

	data := event.Data
	if data["vendor_id"] != "vnd-deluxe" {
		t.Errorf("expected vendor_id 'vnd-deluxe', got %v", data["vendor_id"])
	}
	if data["package_id"] != "pkg-audio-de" {
		t.Errorf("expected package_id 'pkg-audio-de', got %v", data["package_id"])
	}
	if data["component"] != "AUDIO" {
		t.Errorf("expected component 'AUDIO', got %v", data["component"])
	}
	if data["status"] != "FAILED" {
		t.Errorf("expected status 'FAILED', got %v", data["status"])
	}
	if data["defect_category"] != "AUDIO_SYNC_DRIFT" {
		t.Errorf("expected defect_category 'AUDIO_SYNC_DRIFT', got %v", data["defect_category"])
	}
	if data["sync_drift_ms"] != 115.5 {
		t.Errorf("expected sync_drift_ms 115.5, got %v", data["sync_drift_ms"])
	}
}

func TestGenerateVendorHistory_CoverageAndDistribution(t *testing.T) {
	cfg := seed.DefaultConfig()
	cfg.Titles = 2
	cfg.FillerVendors = 0
	cfg.EventsPerVendor = 2000
	cfg.HistoryDays = 120

	rng := seed.NewRNG(cfg.Seed)
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)

	world, err := entities.BuildWorld(cfg, rng, now)
	if err != nil {
		t.Fatalf("BuildWorld failed: %v", err)
	}

	histEvents := events.GenerateVendorHistory(world, cfg, rng, now)

	// 8 curated vendors * 2000 events = 16000 events
	expectedCount := 8 * cfg.EventsPerVendor
	if len(histEvents) != expectedCount {
		t.Fatalf("expected %d events, got %d", expectedCount, len(histEvents))
	}

	vendorFails := make(map[string]int)
	vendorTotals := make(map[string]int)

	for _, ev := range histEvents {
		// Verify event ID is a valid UUIDv5
		if _, err := uuid.Parse(ev.ID); err != nil {
			t.Errorf("event ID '%s' is not a valid UUID: %v", ev.ID, err)
		}

		vID, _ := ev.Data["vendor_id"].(string)
		status, _ := ev.Data["status"].(string)
		comp, _ := ev.Data["component"].(string)
		lang, _ := ev.Data["language"].(string)

		vendorTotals[vID]++
		if status == "FAILED" {
			vendorFails[vID]++
		}

		// Verify Technicolor (VIDEO only) never has AUDIO/SUBTITLE or non-empty language
		if vID == "vnd-technicolor" {
			if comp != "VIDEO" {
				t.Errorf("Technicolor should only have VIDEO events, got %s", comp)
			}
			if lang != "" {
				t.Errorf("Technicolor VIDEO should have empty market language, got %s", lang)
			}
		}

		// Verify Sound & Vision India only has hi-IN and te-IN
		if vID == "vnd-sound-vision-india" {
			if lang != "hi-IN" && lang != "te-IN" {
				t.Errorf("Sound & Vision India should only have hi-IN or te-IN, got %s", lang)
			}
		}

		// Verify timestamp is within 125 days trailing window
		diffDays := now.Sub(ev.Time).Hours() / 24.0
		if diffDays < 0 || diffDays > 125 {
			t.Errorf("event time out of range: diffDays = %f", diffDays)
		}
	}

	// Verify fail rates roughly match expected (1 - TargetAccuracy)
	targets := map[string]float64{
		"vnd-deluxe":             0.01,
		"vnd-testronic":          0.15,
		"vnd-iyuno":              0.07,
		"vnd-pixelogic":          0.04,
		"vnd-sound-vision-india": 0.05,
		"vnd-prasad":             0.08,
		"vnd-technicolor":        0.02,
		"vnd-prime-focus":        0.11,
	}

	for vID, expectedFailRate := range targets {
		total := vendorTotals[vID]
		fails := vendorFails[vID]
		actualFailRate := float64(fails) / float64(total)

		if math.Abs(actualFailRate-expectedFailRate) > 0.03 {
			t.Errorf("vendor %s fail rate expected ~%f, got %f (%d/%d)", vID, expectedFailRate, actualFailRate, fails, total)
		}
	}
}
