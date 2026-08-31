package tools_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestVendorCandidatesTool(t *testing.T) {
	ctx := context.Background()

	t.Run("Returns empty list when client is nil", func(t *testing.T) {
		candidates, err := tools.FetchVendorCandidates(ctx, nil, nil, tools.VendorCandidatesArgs{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(candidates) != 0 {
			t.Errorf("expected 0 candidates, got: %d", len(candidates))
		}
	})

	t.Run("Filters vendor candidates by specialty from SQLite", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)
		defer client.Close()

		v1Res := vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "v-audio",
			},
			Name:            "Audio Labs",
			Specialty:       "AUDIO",
			HourlyRateUSD:   110.0,
			TurnaroundHours: 24,
		})
		if v1Res.IsErr() {
			t.Fatalf("create v-audio failed: %v", v1Res.Error())
		}

		v2Res := vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "v-subs",
			},
			Name:            "Sub World",
			Specialty:       "SUBTITLE",
			HourlyRateUSD:   45.0,
			TurnaroundHours: 12,
		})
		if v2Res.IsErr() {
			t.Fatalf("create v-subs failed: %v", v2Res.Error())
		}

		candidates, err := tools.FetchVendorCandidates(ctx, client, nil, tools.VendorCandidatesArgs{
			Specialty: "AUDIO",
		})
		if err != nil {
			t.Fatalf("FetchVendorCandidates failed: %v", err)
		}

		if len(candidates) != 1 {
			t.Fatalf("expected 1 audio candidate, got: %d", len(candidates))
		}
		if candidates[0].VendorID != "v-audio" {
			t.Errorf("expected v-audio, got: %s", candidates[0].VendorID)
		}
		if candidates[0].HourlyRateUSD != 110.0 {
			t.Errorf("expected $110.0 rate, got: %f", candidates[0].HourlyRateUSD)
		}
		if candidates[0].TurnaroundHours != 24 {
			t.Errorf("expected 24h turnaround, got: %d", candidates[0].TurnaroundHours)
		}
	})
}
