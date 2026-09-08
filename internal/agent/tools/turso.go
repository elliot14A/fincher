package tools

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

const tursoQueryRowCap = 200

var writeVerbPattern = regexp.MustCompile(`(?i)\b(insert|update|delete|drop|alter|create|attach|detach|pragma|replace|vacuum|reindex|truncate|begin|commit|rollback)\b`)

type TursoQueryArgs struct {
	SQL string `json:"sql"`
}

func RunTursoQuery(ctx context.Context, db *sql.DB, rawSQL string) ([]map[string]any, error) {
	if db == nil {
		return nil, domainerrors.NewWithOp("tools.RunTursoQuery", domainerrors.CodeInvalidInput, "turso db is nil", nil)
	}

	trimmed := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(rawSQL), ";"))
	if trimmed == "" {
		return nil, domainerrors.NewWithOp("tools.RunTursoQuery", domainerrors.CodeInvalidInput, "sql cannot be empty", nil)
	}

	lower := strings.ToLower(trimmed)
	if !strings.HasPrefix(lower, "select") && !strings.HasPrefix(lower, "with") {
		return nil, domainerrors.NewWithOp("tools.RunTursoQuery", domainerrors.CodeInvalidInput, "only read-only SELECT/WITH queries are permitted", nil)
	}
	if strings.Contains(trimmed, ";") {
		return nil, domainerrors.NewWithOp("tools.RunTursoQuery", domainerrors.CodeInvalidInput, "multiple statements are not permitted", nil)
	}
	if writeVerbPattern.MatchString(trimmed) {
		return nil, domainerrors.NewWithOp("tools.RunTursoQuery", domainerrors.CodeInvalidInput, "write/DDL statements are not permitted", nil)
	}

	wrapped := fmt.Sprintf("SELECT * FROM (%s) LIMIT %d", trimmed, tursoQueryRowCap)

	logger.Info("tool.query_turso: executing", "sql", trimmed)

	rows, err := db.QueryContext(ctx, wrapped)
	if err != nil {
		logger.Warn("tool.query_turso: query failed", "sql", trimmed, "error", err)
		return nil, fmt.Errorf("query_turso failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	out := make([]map[string]any, 0)
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		m := make(map[string]any, len(cols))
		for i, c := range cols {
			v := vals[i]
			if b, ok := v.([]byte); ok {
				m[c] = string(b)
			} else {
				m[c] = v
			}
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	logger.Info("tool.query_turso: returned rows", "count", len(out))
	return out, nil
}

func NewTursoQueryTool(db *sql.DB) (tool.Tool, error) {
	if db == nil {
		return nil, domainerrors.NewWithOp("tools.NewTursoQueryTool", domainerrors.CodeInvalidInput, "turso db cannot be nil", nil)
	}
	return functiontool.New(
		functiontool.Config{
			Name:        "query_turso",
			Description: tursoQueryToolDescription,
		},
		func(ctx agent.Context, args TursoQueryArgs) ([]map[string]any, error) {
			return RunTursoQuery(ctx, db, args.SQL)
		},
	)
}

const tursoQueryToolDescription = `Run a read-only SQL SELECT against the Fincher operational SQLite database and get matching rows.
Only a single SELECT (or WITH ... SELECT) statement is allowed. Writes/DDL and multiple statements are rejected. Results are capped at 200 rows.

Schema:
- vendors(id TEXT, name TEXT, components TEXT_JSON, markets TEXT_JSON, hourly_rate_usd REAL, turnaround_hours INT)
- media_packages(id TEXT, title_id TEXT, component TEXT, language TEXT, version TEXT, vendor_id TEXT, derived_from_master_version TEXT, redelivery_count INT, status TEXT, market TEXT)
- deliveries(id TEXT, title_id TEXT, country TEXT, status TEXT, target_date DATETIME)
- titles(id TEXT, name TEXT, slug TEXT, type TEXT, premiere_date DATETIME, territories INT, current_master_version TEXT, overall_status TEXT)
- masters(id TEXT, title_id TEXT, version TEXT, supersedes_version TEXT, created_at DATETIME)

IMPORTANT: vendors.components and vendors.markets are JSON arrays stored as TEXT (e.g. ["AUDIO","SUBTITLE"], ["en-US","de-DE"]).
Filter them with json_each, NOT equality. Examples:
- AUDIO vendors serving de-DE:
  SELECT id, name, hourly_rate_usd, turnaround_hours FROM vendors
  WHERE EXISTS (SELECT 1 FROM json_each(components) WHERE value = 'AUDIO')
    AND EXISTS (SELECT 1 FROM json_each(markets) WHERE value = 'de-DE')
- Global VIDEO vendors (VIDEO is market-agnostic; ignore markets):
  SELECT id, name, hourly_rate_usd, turnaround_hours FROM vendors
  WHERE EXISTS (SELECT 1 FROM json_each(components) WHERE value = 'VIDEO')

Only query for the vendors/components/markets you actually need for the requirement at hand. Never fabricate vendor ids; use only ids returned by your queries.`
