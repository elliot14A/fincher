package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/elliot14A/fincher/internal/seed/entities"
	"github.com/elliot14A/fincher/internal/seed/events"
	"github.com/elliot14A/fincher/internal/seed/sink"
	"github.com/elliot14A/fincher/internal/seed/types"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/logger"
)

// SeedSummary captures the record counts and duration of a completed seed operation.
type SeedSummary struct {
	Vendors  int
	Titles   int
	Events   int
	Duration time.Duration
}

// Seeder coordinates the generation and persistence of baseline domain entities and historical ClickHouse events.
type Seeder struct {
	cfg        *types.SeedConfig
	turso      *ent.Client
	chDB       *sql.DB
	entitySink sink.EntitySink
	eventSink  sink.EventSink
	nowFunc    func() time.Time
}

// NewSeeder creates a new Seeder instance.
func NewSeeder(cfg *types.SeedConfig, turso *ent.Client, chDB *sql.DB) *Seeder {
	if cfg == nil {
		cfg = types.DefaultConfig()
	}

	var entitySink sink.EntitySink
	if turso != nil {
		entitySink = sink.NewTursoEntitySink(turso)
	}

	var eventSink sink.EventSink
	if chDB != nil {
		eventSink = sink.NewClickHouseEventSink(chDB, cfg.BatchSize)
	}

	return &Seeder{
		cfg:        cfg,
		turso:      turso,
		chDB:       chDB,
		entitySink: entitySink,
		eventSink:  eventSink,
		nowFunc:    time.Now,
	}
}

// WithSinks overrides default database sinks (useful for isolated testing).
func (s *Seeder) WithSinks(entitySink sink.EntitySink, eventSink sink.EventSink) *Seeder {
	s.entitySink = entitySink
	s.eventSink = eventSink
	return s
}

// WithNow overrides the current timestamp provider.
func (s *Seeder) WithNow(nowFunc func() time.Time) *Seeder {
	s.nowFunc = nowFunc
	return s
}

// Run executes the complete data-generation and persistence pipeline.
func (s *Seeder) Run(ctx context.Context) (*SeedSummary, error) {
	startTime := time.Now()
	now := s.nowFunc().UTC()

	logger.Info("starting fincher starter seeder",
		"reset", s.cfg.Reset,
		"seed", s.cfg.Seed,
		"events_per_vendor", s.cfg.EventsPerVendor,
	)

	// 1. Preflight guard & reset
	if err := PreflightCheck(ctx, s.turso, s.chDB, s.cfg.Reset); err != nil {
		return nil, fmt.Errorf("preflight check failed: %w", err)
	}

	if s.cfg.Reset {
		if err := ResetDatabases(ctx, s.turso, s.chDB); err != nil {
			return nil, fmt.Errorf("database reset failed: %w", err)
		}
	}

	// 2. Initialize deterministic RNG
	rng := types.NewRNG(s.cfg.Seed)

	// 3. Assemble domain entities
	world, err := entities.BuildWorld(s.cfg, rng, now)
	if err != nil {
		return nil, fmt.Errorf("failed to build starter entities: %w", err)
	}

	// 4. Persist domain entities into Turso / SQLite
	if s.entitySink != nil {
		if err := s.entitySink.WriteWorld(ctx, world); err != nil {
			return nil, fmt.Errorf("failed to write entities: %w", err)
		}
	}

	// 5. Generate historical ClickHouse events
	histEvents := events.GenerateVendorHistory(world, s.cfg, rng, now)

	// 6. Persist historical events into ClickHouse
	if s.eventSink != nil && len(histEvents) > 0 {
		if err := s.eventSink.WriteEvents(ctx, histEvents); err != nil {
			return nil, fmt.Errorf("failed to write historical events: %w", err)
		}
	}

	summary := &SeedSummary{
		Vendors:  len(world.Vendors),
		Titles:   len(world.Titles),
		Events:   len(histEvents),
		Duration: time.Since(startTime),
	}

	logger.Info("seeding completed successfully",
		"vendors", summary.Vendors,
		"titles", summary.Titles,
		"events", summary.Events,
		"duration_ms", summary.Duration.Milliseconds(),
	)

	return summary, nil
}
