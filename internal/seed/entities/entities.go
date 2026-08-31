package entities

import (
	"time"

	"github.com/elliot14A/fincher/internal/seed/catalog"
	"github.com/elliot14A/fincher/internal/seed/types"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// BuildWorld constructs the starter relational graph of domain entities (Vendors + Titles) for SQLite.
func BuildWorld(cfg *types.SeedConfig, rng *types.RNG, now time.Time) (*World, error) {
	if cfg == nil {
		cfg = types.DefaultConfig()
	}

	world := &World{
		Vendors: make([]*models.Vendor, 0),
		Titles:  make([]*models.Title, 0),
	}

	// 1. Build Vendors (Curated + Optional Filler)
	curatedVendors := catalog.CuratedVendors()
	for _, vSpec := range curatedVendors {
		world.Vendors = append(world.Vendors, &models.Vendor{
			Base: models.Base{
				ID: vSpec.ID,
				Metadata: map[string]any{
					"poster_url": vSpec.PosterURL,
				},
			},
			Name:            vSpec.Name,
			Components:      vSpec.Components,
			Markets:         vSpec.Markets,
			HourlyRateUSD:   vSpec.HourlyRateUSD,
			TurnaroundHours: int(vSpec.TurnaroundHours),
		})
	}

	fillerCount := cfg.FillerVendors
	if fillerCount > 0 {
		fillers := catalog.GenerateFillerVendors(fillerCount, rng)
		for _, vSpec := range fillers {
			world.Vendors = append(world.Vendors, &models.Vendor{
				Base: models.Base{
					ID: vSpec.ID,
					Metadata: map[string]any{
						"poster_url": vSpec.PosterURL,
					},
				},
				Name:            vSpec.Name,
				Components:      vSpec.Components,
				Markets:         vSpec.Markets,
				HourlyRateUSD:   vSpec.HourlyRateUSD,
				TurnaroundHours: int(vSpec.TurnaroundHours),
			})
		}
	}

	// 2. Build Titles
	curatedTitles := catalog.CuratedTitles()
	titleLimit := cfg.Titles
	if titleLimit <= 0 || titleLimit > len(curatedTitles) {
		titleLimit = len(curatedTitles)
	}

	for i := 0; i < titleLimit; i++ {
		tSpec := curatedTitles[i]
		premiereDate := now.Add(time.Duration(tSpec.PremiereOffsetHours) * time.Hour)

		title := &models.Title{
			Base: models.Base{
				ID: tSpec.ID,
				Metadata: map[string]any{
					"poster_url": tSpec.PosterURL,
					"genre":      tSpec.Genre,
					"synopsis":   tSpec.Synopsis,
					"markets":    tSpec.Markets,
				},
			},
			Name:                 tSpec.Name,
			Slug:                 tSpec.Slug,
			Type:                 tSpec.Type,
			PremiereDate:         premiereDate,
			Territories:          len(tSpec.Markets),
			CurrentMasterVersion: "V01",
			OverallStatus:        tSpec.OverallStatus,
		}
		world.Titles = append(world.Titles, title)
	}

	return world, nil
}
