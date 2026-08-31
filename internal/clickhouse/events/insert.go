package events

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/clickhouse"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Insert writes an immutable event into fincher.events.
func Insert(ctx context.Context, db *sql.DB, event *models.Event) domainerrors.Result[*models.Event] {
	if err := event.Validate(); err != nil {
		return domainerrors.Err[*models.Event](clickhouse.NewError("events.Insert", domainerrors.CodeInvalidInput, "invalid event payload", err))
	}

	dataJSON, err := event.DataJSON()
	if err != nil {
		return domainerrors.Err[*models.Event](clickhouse.NewError("events.Insert", domainerrors.CodeInvalidInput, "failed to marshal event data", err))
	}

	subject := event.Subject
	if subject == "" {
		subject = models.DefaultTitleAgnosticSentinel
	}

	source := event.Source
	if source == "" {
		source = "fincher.system"
	}

	parsedID, err := uuid.Parse(event.ID)
	if err != nil {
		parsedID = uuid.New()
		event.ID = parsedID.String()
	}

	query := `
		insert into fincher.events (
			id, type, source, subject, time, data, severity, datacontenttype
		) values (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	if _, err := db.ExecContext(ctx, query,
		parsedID,
		event.Type,
		source,
		subject,
		event.Time,
		dataJSON,
		string(event.Severity),
		event.DataContentType,
	); err != nil {
		return domainerrors.Err[*models.Event](clickhouse.MapError("events.Insert", "event", event.ID, err))
	}

	return domainerrors.Ok(event)
}

// InsertBatch writes multiple immutable events into fincher.events in a single batch transaction.
func InsertBatch(ctx context.Context, db *sql.DB, events []models.Event) domainerrors.Result[int] {
	if len(events) == 0 {
		return domainerrors.Ok(0)
	}
	if db == nil {
		return domainerrors.Err[int](clickhouse.NewError("events.InsertBatch", domainerrors.CodeInvalidInput, "database connection cannot be nil", nil))
	}

	for i := range events {
		if err := events[i].Validate(); err != nil {
			return domainerrors.Err[int](clickhouse.NewError("events.InsertBatch", domainerrors.CodeInvalidInput, fmt.Sprintf("event at index %d is invalid", i), err))
		}
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return domainerrors.Err[int](clickhouse.MapError("events.InsertBatch", "events", "begin_tx", err))
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		insert into fincher.events (
			id, type, source, subject, time, data, severity, datacontenttype
		) values (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`)
	if err != nil {
		return domainerrors.Err[int](clickhouse.MapError("events.InsertBatch", "events", "prepare_stmt", err))
	}
	defer stmt.Close()

	for i := range events {
		dataJSON, err := events[i].DataJSON()
		if err != nil {
			return domainerrors.Err[int](clickhouse.NewError("events.InsertBatch", domainerrors.CodeInvalidInput, "failed to marshal event data", err))
		}

		subject := events[i].Subject
		if subject == "" {
			subject = models.DefaultTitleAgnosticSentinel
		}

		source := events[i].Source
		if source == "" {
			source = "fincher.system"
		}

		parsedID, err := uuid.Parse(events[i].ID)
		if err != nil {
			parsedID = uuid.New()
			events[i].ID = parsedID.String()
		}

		if _, err := stmt.ExecContext(ctx,
			parsedID,
			events[i].Type,
			source,
			subject,
			events[i].Time,
			dataJSON,
			string(events[i].Severity),
			events[i].DataContentType,
		); err != nil {
			return domainerrors.Err[int](clickhouse.MapError("events.InsertBatch", "events", events[i].ID, err))
		}
	}

	if err := tx.Commit(); err != nil {
		return domainerrors.Err[int](clickhouse.MapError("events.InsertBatch", "events", "commit", err))
	}

	return domainerrors.Ok(len(events))
}
