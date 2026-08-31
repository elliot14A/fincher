package events

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/seed/catalog"
	"github.com/elliot14A/fincher/internal/seed/entities"
	"github.com/elliot14A/fincher/internal/seed/types"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

type coveredPair struct {
	Component string
	MarketTag string
}

// GenerateVendorHistory creates historical QC inspection events for all vendors, emitting events ONLY for covered (component, market) pairs.
func GenerateVendorHistory(world *entities.World, cfg *types.SeedConfig, rng *types.RNG, now time.Time) []models.Event {
	if cfg == nil {
		cfg = types.DefaultConfig()
	}

	curatedSpecs := make(map[string]catalog.VendorSpec)
	for _, spec := range catalog.CuratedVendors() {
		curatedSpecs[spec.ID] = spec
	}

	var titleSlugs []string
	for _, t := range world.Titles {
		titleSlugs = append(titleSlugs, t.Slug)
	}
	if len(titleSlugs) == 0 {
		titleSlugs = []string{"avatar-fire-ash", "pushpa-the-rule", "dune-part-two", "shogun"}
	}

	totalCapacity := len(world.Vendors) * cfg.EventsPerVendor
	allEvents := make([]models.Event, 0, totalCapacity)

	for _, v := range world.Vendors {
		spec, isCurated := curatedSpecs[v.ID]
		if !isCurated {
			spec = catalog.VendorSpec{
				ID:              v.ID,
				Name:            v.Name,
				Components:      v.Components,
				Markets:         v.Markets,
				HourlyRateUSD:   v.HourlyRateUSD,
				TurnaroundHours: float64(v.TurnaroundHours),
				TargetAccuracy:  0.90,
				DriftMeanMS:     40.0,
				DriftStdDevMS:   20.0,
			}
		}

		// Enumerate all covered (component, market) pairs for this vendor
		var pairs []coveredPair
		for _, comp := range spec.Components {
			if comp == "VIDEO" {
				pairs = append(pairs, coveredPair{Component: "VIDEO", MarketTag: ""})
			} else {
				for _, m := range spec.Markets {
					pairs = append(pairs, coveredPair{Component: comp, MarketTag: m})
				}
			}
		}

		if len(pairs) == 0 {
			// Fallback for safety
			pairs = append(pairs, coveredPair{Component: "AUDIO", MarketTag: "en-US"})
		}

		eventsCount := cfg.EventsPerVendor
		if eventsCount <= 0 {
			eventsCount = 10000
		}

		historyDays := float64(cfg.HistoryDays)
		if historyDays <= 0 {
			historyDays = 120.0
		}

		for i := 0; i < eventsCount; i++ {
			pair := pairs[i%len(pairs)]

			// Exponential recency-weighted day offset
			u := rng.Float64()
			decayConst := historyDays / 2.0
			dayOffset := -decayConst * math.Log(1.0-u*(1.0-math.Exp(-historyDays/decayConst)))
			if dayOffset < 0 {
				dayOffset = 0
			} else if dayOffset > historyDays {
				dayOffset = historyDays
			}

			intraDayHours := rng.FloatInRange(0.0, 24.0)
			eventTime := now.Add(-time.Duration((dayOffset*24.0)+intraDayHours) * time.Hour)

			// Draw QC outcome: P(FAILED) = 1.0 - TargetAccuracy
			roll := rng.Float64()
			failThreshold := 1.0 - spec.TargetAccuracy
			warnThreshold := failThreshold + (spec.TargetAccuracy * 0.02)

			status := "PASSED"
			defectCat := "NONE"
			severity := models.SeverityInfo

			var syncDrift *float64
			var videoScore *float64

			if roll < failThreshold {
				status = "FAILED"
				severity = models.SeverityWarn

				switch pair.Component {
				case "AUDIO":
					defectCat = "AUDIO_SYNC_DRIFT"
					driftVal := spec.DriftMeanMS + rng.NormFloat64(0, spec.DriftStdDevMS)
					if driftVal < 60.0 {
						driftVal = 60.0 + rng.FloatInRange(5.0, 50.0)
					}
					syncDrift = &driftVal
				case "VIDEO":
					defectCat = "CORRUPT_FRAME"
					scoreVal := rng.FloatInRange(0.35, 0.95)
					videoScore = &scoreVal
				case "SUBTITLE":
					defectCat = "SUBTITLE_OVERLAP"
				default:
					defectCat = "OTHER"
				}
			} else if roll < warnThreshold {
				status = "WARNING"
				severity = models.SeverityWarn

				if pair.Component == "AUDIO" {
					defectCat = "AUDIO_SYNC_DRIFT"
					driftVal := rng.FloatInRange(35.0, 55.0)
					syncDrift = &driftVal
				} else if pair.Component == "VIDEO" {
					defectCat = "CORRUPT_FRAME"
					scoreVal := rng.FloatInRange(0.10, 0.30)
					videoScore = &scoreVal
				} else {
					defectCat = "SUBTITLE_OVERLAP"
				}
			} else {
				if pair.Component == "AUDIO" {
					driftVal := rng.FloatInRange(2.0, 25.0)
					syncDrift = &driftVal
				} else if pair.Component == "VIDEO" {
					scoreVal := rng.FloatInRange(0.0, 0.05)
					videoScore = &scoreVal
				}
			}

			titleSlug := types.Choice(rng, titleSlugs)
			pkgID := fmt.Sprintf("pkg-hist-%s-%06d", spec.ID, i)

			// Deterministic UUIDv5
			seedKey := fmt.Sprintf("fincher-seed-%d-%s-%06d", cfg.Seed, spec.ID, i)
			eventUUID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(seedKey)).String()
			inspectorSource := fmt.Sprintf("qc.%s", spec.ID)

			allEvents = append(allEvents, BuildQCEvent(QCEventParams{
				ID:                   eventUUID,
				Time:                 eventTime,
				TitleSlug:            titleSlug,
				InspectorAgent:       inspectorSource,
				PackageID:            pkgID,
				VendorID:             spec.ID,
				Component:            pair.Component,
				Language:             pair.MarketTag,
				Status:               status,
				SyncDriftMS:          syncDrift,
				VideoCorruptionScore: videoScore,
				DefectCategory:       defectCat,
				Severity:             severity,
			}))
		}
	}

	return allEvents
}
