package sink

import (
	"context"

	"github.com/elliot14A/fincher/internal/seed/entities"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// EntitySink defines the interface for persisting the assembled relational world graph.
type EntitySink interface {
	WriteWorld(ctx context.Context, world *entities.World) error
}

// EventSink defines the interface for inserting batches of generated events.
type EventSink interface {
	WriteEvents(ctx context.Context, events []models.Event) error
}
