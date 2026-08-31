package tools_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestImpactTool(t *testing.T) {
	ctx := context.Background()

	t.Run("Rejects nil client with error", func(t *testing.T) {
		_, err := tools.FetchDeliveryImpact(ctx, nil, tools.DeliveryImpactArgs{
			PackageID: "pkg-1",
		})
		if err == nil {
			t.Fatal("expected error for nil client, got nil")
		}
	})

	t.Run("Rejects nil client in tool constructor", func(t *testing.T) {
		_, err := tools.NewDeliveryImpactTool(nil)
		if err == nil {
			t.Fatal("expected error for nil client constructor, got nil")
		}
	})

	t.Run("Calculates market blast radius from SQLite", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)
		defer client.Close()

		premiere := time.Now().Add(36 * time.Hour)
		titleRes := titles.Create(ctx, client, &models.Title{
			Base: models.Base{
				ID: "title-matrix",
			},
			Name:                 "The Matrix",
			Type:                 models.TitleTypeFeature,
			PremiereDate:         premiere,
			Territories:          2,
			CurrentMasterVersion: "master-v1",
			OverallStatus:        models.StatusProcessing,
		})
		if titleRes.IsErr() {
			t.Fatalf("create title failed: %v", titleRes.Error())
		}
		title := titleRes.Unwrap()

		delDERes := deliveries.Create(ctx, client, &models.Delivery{
			Base: models.Base{
				ID: "del-de",
			},
			TitleID:    title.ID,
			Country:    "DE",
			Status:     models.DeliveryStatusReadyToShip,
			TargetDate: premiere,
		})
		if delDERes.IsErr() {
			t.Fatalf("create delivery DE failed: %v", delDERes.Error())
		}

		delUSRes := deliveries.Create(ctx, client, &models.Delivery{
			Base: models.Base{
				ID: "del-us",
			},
			TitleID:    title.ID,
			Country:    "US",
			Status:     models.DeliveryStatusReadyToShip,
			TargetDate: premiere,
		})
		if delUSRes.IsErr() {
			t.Fatalf("create delivery US failed: %v", delUSRes.Error())
		}

		vRes := vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "vendor-dub",
			},
			Name:            "Dub Studio",
			Components:      []string{"AUDIO"},
			Markets:         []string{"en-US"},
			HourlyRateUSD:   90.0,
			TurnaroundHours: 24,
		})
		if vRes.IsErr() {
			t.Fatalf("create vendor failed: %v", vRes.Error())
		}

		pkgRes := packages.Create(ctx, client, &models.Package{
			Base: models.Base{
				ID: "pkg-de-audio",
			},
			TitleID:                  title.ID,
			VendorID:                 "vendor-dub",
			Component:                models.ComponentAudio,
			Language:                 "de-DE",
			Version:                  "v1.0",
			DerivedFromMasterVersion: "master-v1",
			Status:                   models.PackageStatusValid,
		})
		if pkgRes.IsErr() {
			t.Fatalf("create package failed: %v", pkgRes.Error())
		}

		impact, err := tools.FetchDeliveryImpact(ctx, client, tools.DeliveryImpactArgs{
			PackageID: "pkg-de-audio",
		})
		if err != nil {
			t.Fatalf("FetchDeliveryImpact failed: %v", err)
		}

		if len(impact.AffectedDeliveries) != 1 || impact.AffectedDeliveries[0] != "del-de" {
			t.Errorf("expected only del-de affected, got: %v", impact.AffectedDeliveries)
		}
		if len(impact.AffectedMarkets) != 1 || impact.AffectedMarkets[0] != "DE" {
			t.Errorf("expected market DE, got: %v", impact.AffectedMarkets)
		}
		if !impact.IsPremiereUrgent {
			t.Errorf("expected IsPremiereUrgent to be true for 36h timeline, got false")
		}
	})
}
