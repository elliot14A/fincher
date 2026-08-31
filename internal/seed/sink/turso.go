package sink

import (
	"context"
	"fmt"

	"github.com/elliot14A/fincher/internal/seed/entities"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/logger"
)

type tursoEntitySink struct {
	client *ent.Client
}

// NewTursoEntitySink creates an EntitySink that persists starter domain models (Vendors + Titles) into Turso / SQLite.
func NewTursoEntitySink(client *ent.Client) EntitySink {
	return &tursoEntitySink{client: client}
}

func (s *tursoEntitySink) WriteWorld(ctx context.Context, world *entities.World) error {
	if s.client == nil {
		return fmt.Errorf("tursoEntitySink: ent client is nil")
	}

	// 1. Vendors
	for _, v := range world.Vendors {
		res := vendors.Create(ctx, s.client, v)
		if res.IsErr() {
			return fmt.Errorf("failed to persist vendor %s: %w", v.ID, res.Error())
		}
	}
	logger.Info("seeded vendors", "count", len(world.Vendors))

	// 2. Titles
	for _, t := range world.Titles {
		res := titles.Create(ctx, s.client, t)
		if res.IsErr() {
			return fmt.Errorf("failed to persist title %s: %w", t.ID, res.Error())
		}
	}
	logger.Info("seeded titles", "count", len(world.Titles))

	return nil
}
