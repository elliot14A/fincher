package seed_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/seed"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
)

func TestSeeder_PipelineWithTurso(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)
	defer client.Close()

	cfg := seed.DefaultConfig()
	cfg.Titles = 7
	cfg.FillerVendors = 2
	cfg.EventsPerVendor = 10 // small volume for unit test
	cfg.Reset = false

	seeder := seed.NewSeeder(cfg, client, nil)
	summary, err := seeder.Run(ctx)
	if err != nil {
		t.Fatalf("seeder.Run failed: %v", err)
	}

	// 8 curated + 2 filler = 10 vendors
	if summary.Vendors != 10 {
		t.Errorf("expected 10 vendors, got %d", summary.Vendors)
	}
	if summary.Titles != 7 {
		t.Errorf("expected 7 titles, got %d", summary.Titles)
	}

	// 2. Test preflight guard: second run without reset must fail
	_, err = seeder.Run(ctx)
	if err == nil {
		t.Fatal("expected preflight guard error on second run without --reset, got nil")
	}

	// 3. Third run with reset must succeed
	cfg.Reset = true
	summary2, err := seeder.Run(ctx)
	if err != nil {
		t.Fatalf("expected seeder with reset to succeed, got: %v", err)
	}
	if summary2.Vendors != 10 {
		t.Errorf("expected 10 vendors after reset, got %d", summary2.Vendors)
	}
}
