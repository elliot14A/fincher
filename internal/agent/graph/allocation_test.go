package graph_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/agent/graph"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestExecuteAllocation(t *testing.T) {
	ctx := context.Background()

	t.Run("Fails when title slug is empty", func(t *testing.T) {
		deps := graph.AllocationGraphDeps{
			Model: &mockLLM{},
		}
		_, err := graph.ExecuteAllocation(ctx, deps, graph.AllocationInput{
			TitleSlug: "",
			Requirements: []models.AllocationRequirement{
				{Component: "AUDIO", Market: "en-US"},
			},
		})
		if err == nil {
			t.Fatal("expected error for empty title slug, got nil")
		}
	})

	t.Run("Fails when requirements list is empty", func(t *testing.T) {
		deps := graph.AllocationGraphDeps{
			Model: &mockLLM{},
		}
		_, err := graph.ExecuteAllocation(ctx, deps, graph.AllocationInput{
			TitleSlug:    "avatar-fire-ash",
			Requirements: []models.AllocationRequirement{},
		})
		if err == nil {
			t.Fatal("expected error for empty requirements, got nil")
		}
	})
}
