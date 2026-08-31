package tools_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/agent/tools"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestProjectionTool_RejectsNilClient(t *testing.T) {
	ctx := context.Background()
	_, err := tools.GetTitleReadyProjection(ctx, nil, "avatar-fire-ash")
	if err == nil {
		t.Fatal("expected error with nil client, got nil")
	}

	_, err = tools.NewProjectionTool(nil)
	if err == nil {
		t.Fatal("expected error with nil client in constructor, got nil")
	}
}

func TestProjectionTool_CalculatesSafeBuffer(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-deluxe"},
		Name:            "Deluxe Audio",
		Components:      []string{"AUDIO"},
		Markets:         []string{"de-DE"},
		TurnaroundHours: 12,
	})

	titleSlug := "proj-title-safe-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Avatar Safe Test",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(72 * time.Hour), // 72h until premiere
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusHold,
	})
	titleObj := titleRes.Unwrap()

	// Add 1 broken package with 12h repair
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-" + uuid.NewString()[:8]},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusInvalidated,
	})

	proj, err := tools.GetTitleReadyProjection(ctx, client, titleSlug)
	if err != nil {
		t.Fatalf("GetTitleReadyProjection failed: %v", err)
	}

	// Critical: 12h, Premiere: ~72h, Buffer: ~60h => SAFE
	if proj.RiskBand != "SAFE" {
		t.Errorf("expected RiskBand SAFE, got: %s", proj.RiskBand)
	}
	if proj.IsBreached {
		t.Errorf("expected IsBreached false")
	}
	if proj.CriticalRemainingHours != 12.0 {
		t.Errorf("expected CriticalRemainingHours 12.0, got: %v", proj.CriticalRemainingHours)
	}
}

func TestProjectionTool_CalculatesBreachBuffer(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-slow"},
		Name:            "Slow Dubbing",
		Components:      []string{"AUDIO"},
		Markets:         []string{"de-DE"},
		TurnaroundHours: 24,
	})

	titleSlug := "proj-title-breach-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Avatar Breach Test",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(8 * time.Hour), // Only 8h until premiere!
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusHold,
	})
	titleObj := titleRes.Unwrap()

	// 24h repair needed with only 8h remaining => -16h buffer => BREACH
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-" + uuid.NewString()[:8]},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "V01",
		VendorID:                 "vendor-slow",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusInvalidated,
	})

	proj, err := tools.GetTitleReadyProjection(ctx, client, titleSlug)
	if err != nil {
		t.Fatalf("GetTitleReadyProjection failed: %v", err)
	}

	if proj.RiskBand != "BREACH" {
		t.Errorf("expected RiskBand BREACH, got: %s", proj.RiskBand)
	}
	if !proj.IsBreached {
		t.Errorf("expected IsBreached true")
	}
	if proj.BufferHours >= 0 {
		t.Errorf("expected negative buffer hours, got: %v", proj.BufferHours)
	}
}

func TestProjectionTool_SequentialMasterPlusDubbingChaining(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-deluxe"},
		Name:            "Deluxe Audio",
		Components:      []string{"AUDIO"},
		Markets:         []string{"en-US"},
		TurnaroundHours: 12, // 12h dub turnaround
	})

	titleSlug := "proj-seq-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Avatar Sequential Test",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(48 * time.Hour), // 48h premiere
		Territories:          1,
		CurrentMasterVersion: "V02", // Master bumped to V02
		OverallStatus:        models.StatusHold,
	})
	titleObj := titleRes.Unwrap()

	// Localized package derived from V01 (needs 48h master reconform + 12h dub repair = 60h)
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-" + uuid.NewString()[:8]},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01", // Mismatch with title's V02
		Status:                   models.PackageStatusInvalidated,
	})

	proj, err := tools.GetTitleReadyProjection(ctx, client, titleSlug)
	if err != nil {
		t.Fatalf("GetTitleReadyProjection failed: %v", err)
	}

	// Critical remaining hours should sum 48h master reconform + 12h dubbing = 60h
	if proj.CriticalRemainingHours != 60.0 {
		t.Errorf("expected CriticalRemainingHours 60.0 (48h master + 12h dub), got: %v", proj.CriticalRemainingHours)
	}
}
