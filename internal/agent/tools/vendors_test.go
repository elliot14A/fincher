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

	t.Run("Rejects nil client with error", func(t *testing.T) {
		_, err := tools.FetchVendorCandidates(ctx, nil, nil, tools.VendorCandidatesArgs{})
		if err == nil {
			t.Fatal("expected error for nil client, got nil")
		}
	})

	t.Run("Rejects nil client in tool constructor", func(t *testing.T) {
		_, err := tools.NewVendorCandidatesTool(nil, nil)
		if err == nil {
			t.Fatal("expected error for nil client constructor, got nil")
		}
	})

	t.Run("Filters vendor candidates by component and market", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)
		defer client.Close()

		// 1. Deluxe (Western + Indian, AUDIO + SUBTITLE)
		_ = vendors.Create(ctx, client, &models.Vendor{
			Base:            models.Base{ID: "vnd-deluxe"},
			Name:            "Deluxe Media",
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         []string{"en-US", "de-DE", "fr-FR", "hi-IN", "te-IN"},
			HourlyRateUSD:   200.0,
			TurnaroundHours: 12,
		})

		// 2. Iyuno (Western only, AUDIO + SUBTITLE)
		_ = vendors.Create(ctx, client, &models.Vendor{
			Base:            models.Base{ID: "vnd-iyuno"},
			Name:            "Iyuno SDI",
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         []string{"en-US", "de-DE", "fr-FR"},
			HourlyRateUSD:   70.0,
			TurnaroundHours: 36,
		})

		// 3. Sound & Vision India (Indian only, AUDIO + SUBTITLE)
		_ = vendors.Create(ctx, client, &models.Vendor{
			Base:            models.Base{ID: "vnd-sound-vision-india"},
			Name:            "Sound & Vision India",
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         []string{"hi-IN", "te-IN"},
			HourlyRateUSD:   90.0,
			TurnaroundHours: 20,
		})

		// 4. Technicolor (VIDEO only, global market)
		_ = vendors.Create(ctx, client, &models.Vendor{
			Base:            models.Base{ID: "vnd-technicolor"},
			Name:            "Technicolor",
			Components:      []string{"VIDEO"},
			Markets:         []string{},
			HourlyRateUSD:   185.0,
			TurnaroundHours: 16,
		})

		// Query 1: AUDIO for Telugu (te-IN) -> should match Deluxe and Sound & Vision India, EXCLUDE Iyuno and Technicolor
		teCandidates, err := tools.FetchVendorCandidates(ctx, client, nil, tools.VendorCandidatesArgs{
			Component: "AUDIO",
			Market:    "te-IN",
		})
		if err != nil {
			t.Fatalf("FetchVendorCandidates te-IN failed: %v", err)
		}
		if len(teCandidates) != 2 {
			t.Fatalf("expected 2 Telugu audio candidates, got: %d", len(teCandidates))
		}
		for _, c := range teCandidates {
			if c.VendorID == "vnd-iyuno" || c.VendorID == "vnd-technicolor" {
				t.Errorf("unexpected vendor %s for Telugu audio", c.VendorID)
			}
		}

		// Query 2: AUDIO for German (de-DE) -> should match Deluxe and Iyuno, EXCLUDE Sound & Vision India
		deCandidates, err := tools.FetchVendorCandidates(ctx, client, nil, tools.VendorCandidatesArgs{
			Component: "AUDIO",
			Market:    "de-DE",
		})
		if err != nil {
			t.Fatalf("FetchVendorCandidates de-DE failed: %v", err)
		}
		if len(deCandidates) != 2 {
			t.Fatalf("expected 2 German audio candidates, got: %d", len(deCandidates))
		}
		for _, c := range deCandidates {
			if c.VendorID == "vnd-sound-vision-india" || c.VendorID == "vnd-technicolor" {
				t.Errorf("unexpected vendor %s for German audio", c.VendorID)
			}
		}

		// Query 3: VIDEO (market == "") -> should match Technicolor, EXCLUDE Iyuno and Sound & Vision India
		videoCandidates, err := tools.FetchVendorCandidates(ctx, client, nil, tools.VendorCandidatesArgs{
			Component: "VIDEO",
		})
		if err != nil {
			t.Fatalf("FetchVendorCandidates VIDEO failed: %v", err)
		}
		if len(videoCandidates) != 1 || videoCandidates[0].VendorID != "vnd-technicolor" {
			t.Fatalf("expected 1 VIDEO candidate (vnd-technicolor), got: %v", videoCandidates)
		}
	})
}
