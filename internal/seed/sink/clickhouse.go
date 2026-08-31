package sink

import (
	"context"
	"database/sql"
	"fmt"

	chevents "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

type clickHouseEventSink struct {
	db        *sql.DB
	batchSize int
}

// NewClickHouseEventSink creates an EventSink that bulk-inserts CloudEvents into ClickHouse fincher.events.
func NewClickHouseEventSink(db *sql.DB, batchSize int) EventSink {
	if batchSize <= 0 {
		batchSize = 2000
	}
	return &clickHouseEventSink{
		db:        db,
		batchSize: batchSize,
	}
}

func (s *clickHouseEventSink) WriteEvents(ctx context.Context, events []models.Event) error {
	if s.db == nil {
		return fmt.Errorf("clickHouseEventSink: database connection is nil")
	}

	total := len(events)
	if total == 0 {
		return nil
	}

	inserted := 0
	for i := 0; i < total; i += s.batchSize {
		end := i + s.batchSize
		if end > total {
			end = total
		}

		chunk := events[i:end]
		res := chevents.InsertBatch(ctx, s.db, chunk)
		if res.IsErr() {
			return fmt.Errorf("failed to insert batch at offset %d: %w", i, res.Error())
		}

		inserted += res.Unwrap()
		if inserted%10000 == 0 || end == total {
			logger.Info("inserted historical QC events into ClickHouse", "inserted", inserted, "total", total)
		}
	}

	return nil
}
