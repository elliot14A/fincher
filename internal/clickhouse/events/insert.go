package events

import (
	"context"
	"database/sql"

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
