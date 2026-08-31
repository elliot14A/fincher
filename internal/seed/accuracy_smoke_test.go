package seed_test

import (
	"context"
	"math"
	"testing"

	"github.com/elliot14A/fincher/internal/clickhouse"
	chvendors "github.com/elliot14A/fincher/internal/clickhouse/vendors"
)

func TestLiveClickHouse_AccuracyRollups(t *testing.T) {
	db, err := clickhouse.Open("127.0.0.1:9000")
	if err != nil {
		t.Skipf("skipping live clickhouse test: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	targets := []struct {
		vendorID       string
		component      string
		targetAccuracy float64
	}{
		{"vnd-deluxe", "AUDIO", 0.99},
		{"vnd-iyuno", "AUDIO", 0.93},
		{"vnd-testronic", "AUDIO", 0.85},
		{"vnd-pixelogic", "SUBTITLE", 0.96},
		{"vnd-sound-vision-india", "AUDIO", 0.95},
		{"vnd-prasad", "AUDIO", 0.92},
		{"vnd-technicolor", "VIDEO", 0.98},
		{"vnd-prime-focus", "VIDEO", 0.89},
	}

	for _, tc := range targets {
		res := chvendors.RecencyWeightedAccuracy(ctx, db, tc.vendorID, tc.component)
		if res.IsErr() {
			t.Fatalf("failed to calculate accuracy for %s: %v", tc.vendorID, res.Error())
		}
		actualAcc := res.Unwrap()
		t.Logf("Vendor: %s, Component: %s, Target: %.3f, Actual: %.3f", tc.vendorID, tc.component, tc.targetAccuracy, actualAcc)

		if math.Abs(actualAcc-tc.targetAccuracy) > 0.03 {
			t.Errorf("vendor %s accuracy expected ~%.2f, got %.3f", tc.vendorID, tc.targetAccuracy, actualAcc)
		}
	}
}
