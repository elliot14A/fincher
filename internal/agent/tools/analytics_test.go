package tools_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/agent/tools"
)

func TestAnalyticsTool(t *testing.T) {
	ctx := context.Background()

	t.Run("Rejects nil db with error", func(t *testing.T) {
		_, err := tools.FetchAnalytics(ctx, nil, tools.AnalyticsArgs{
			VendorID:  "vendor-1",
			TitleSlug: "eclipse",
			Component: "AUDIO",
		})
		if err == nil {
			t.Fatal("expected error for nil db, got nil")
		}
	})

	t.Run("Rejects nil db in tool constructor", func(t *testing.T) {
		_, err := tools.NewAnalyticsTool(nil)
		if err == nil {
			t.Fatal("expected error for nil db constructor, got nil")
		}
	})
}
