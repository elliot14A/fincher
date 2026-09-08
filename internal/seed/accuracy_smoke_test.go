package seed_test

import (
	"context"
	"math"
	"testing"

	chvendors "github.com/elliot14A/fincher/internal/clickhouse/vendors"
	"github.com/elliot14A/fincher/pkg/mcp"
)

func TestLiveClickHouse_AccuracyRollups(t *testing.T) {
	ctx := context.Background()
	mcpClient, err := mcp.NewClient("http://127.0.0.1:8000/mcp")
	if err != nil {
		t.Skipf("skipping live mcp test: %v", err)
	}
	if err := mcpClient.Ping(ctx); err != nil {
		t.Skipf("skipping live mcp test: mcp server not reachable: %v", err)
	}

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
		res := chvendors.RecencyWeightedAccuracy(ctx, mcpClient, tc.vendorID, tc.component)
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
