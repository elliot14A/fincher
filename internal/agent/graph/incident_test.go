package graph_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/agent/graph"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestExecuteIncident(t *testing.T) {
	ctx := context.Background()

	t.Run("Filters routine benign event without executing action plan", func(t *testing.T) {
		llm := &mockLLM{
			responses: []string{
				`{
					"actionable": false,
					"severity": "INFO",
					"anomaly_type": "NONE",
					"rationale": "Routine inspection completed within acceptable parameters."
				}`,
			},
		}

		event := &models.Event{
			ID:       "evt-benign-1",
			Type:     models.TypeQCInspectionCompleted,
			Severity: models.SeverityInfo,
			Subject:  "eclipse",
			Data:     map[string]any{"package_id": "pkg-1"},
		}

		client := tursotest.NewMemoryClient(t)
		defer client.Close()

		deps := graph.IncidentGraphDeps{
			Model:       llm,
			TursoClient: client,
		}

		output, err := graph.ExecuteIncident(ctx, deps, graph.IncidentInput{
			Event:              event,
			HoursUntilPremiere: 72.0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.Actionable {
			t.Errorf("expected actionable false, got true")
		}
		if output.Decision != "FILTERED" {
			t.Errorf("expected decision FILTERED, got: %s", output.Decision)
		}
		if output.RunnerResult != nil {
			t.Errorf("expected nil runner result for filtered event")
		}
	})

}
