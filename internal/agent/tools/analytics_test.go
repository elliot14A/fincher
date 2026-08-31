package tools_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestAnalyticsTool(t *testing.T) {
	ctx := context.Background()

	t.Run("Returns default summary when db is nil", func(t *testing.T) {
		summary, err := tools.FetchAnalytics(ctx, nil, tools.AnalyticsArgs{
			VendorID:  "vendor-1",
			TitleSlug: "eclipse",
			Component: "AUDIO",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if summary == nil {
			t.Fatal("expected non-nil summary")
		}
		if summary.VendorHistoricalAccuracy != models.UnmeasuredHistoricalAccuracy {
			t.Errorf("expected unmeasured accuracy %f, got: %f", models.UnmeasuredHistoricalAccuracy, summary.VendorHistoricalAccuracy)
		}
	})

	t.Run("Initializes ADK tool instance cleanly", func(t *testing.T) {
		adkTool, err := tools.NewAnalyticsTool(nil)
		if err != nil {
			t.Fatalf("failed to create analytics tool: %v", err)
		}
		if adkTool == nil {
			t.Fatal("expected non-nil ADK tool")
		}
		if adkTool.Name() != "query_analytics" {
			t.Errorf("expected tool name query_analytics, got: %s", adkTool.Name())
		}
	})
}
