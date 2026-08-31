package entities_test

import (
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/seed/entities"
	"github.com/elliot14A/fincher/internal/seed/types"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestBuildWorld(t *testing.T) {
	cfg := types.DefaultConfig()
	cfg.Titles = 7
	cfg.FillerVendors = 2

	rng := types.NewRNG(cfg.Seed)
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)

	world, err := entities.BuildWorld(cfg, rng, now)
	if err != nil {
		t.Fatalf("BuildWorld failed: %v", err)
	}

	// 8 curated + 2 filler = 10 vendors
	if len(world.Vendors) != 10 {
		t.Errorf("expected 10 vendors, got %d", len(world.Vendors))
	}

	// 7 titles
	if len(world.Titles) != 7 {
		t.Errorf("expected 7 titles, got %d", len(world.Titles))
	}

	// Find hero title
	var heroTitle *models.Title
	for _, title := range world.Titles {
		if title.Slug == "avatar-fire-ash" {
			heroTitle = title
			break
		}
	}
	if heroTitle == nil {
		t.Fatal("expected hero title 'avatar-fire-ash' in world")
	}

	if heroTitle.OverallStatus != models.StatusOnTrack {
		t.Errorf("expected hero title status ON_TRACK, got %s", heroTitle.OverallStatus)
	}

	// Check poster_url and markets in title metadata
	if heroTitle.Metadata == nil || heroTitle.Metadata["poster_url"] == nil || heroTitle.Metadata["poster_url"] == "" {
		t.Errorf("expected non-empty poster_url in hero title metadata, got %v", heroTitle.Metadata)
	}
	marketsRaw, ok := heroTitle.Metadata["markets"].([]string)
	if !ok || len(marketsRaw) != 5 {
		t.Errorf("expected 5 markets in hero title metadata, got %v", heroTitle.Metadata["markets"])
	}

	// Check components and markets on vendors
	for _, v := range world.Vendors {
		if len(v.Components) == 0 {
			t.Errorf("vendor %s has no components", v.ID)
		}
		if v.Metadata == nil || v.Metadata["poster_url"] == nil || v.Metadata["poster_url"] == "" {
			t.Errorf("expected non-empty poster_url in vendor metadata, got %v", v.Metadata)
		}
	}

	// Verify scarce te-IN coverage
	teINCount := 0
	for _, v := range world.Vendors {
		for _, m := range v.Markets {
			if m == "te-IN" {
				teINCount++
				break
			}
		}
	}
	if teINCount < 4 { // Deluxe, Pixelogic, Sound & Vision, Prasad (+ Prime Focus)
		t.Errorf("expected at least 4 vendors supporting te-IN, got %d", teINCount)
	}
}
