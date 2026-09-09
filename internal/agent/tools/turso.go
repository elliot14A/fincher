package tools

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

const tursoQueryRowCap = 100

var writeVerbPattern = regexp.MustCompile(`(?i)\b(insert|update|delete|drop|alter|create|attach|detach|pragma|replace|vacuum|reindex|truncate|begin|commit|rollback)\b`)

// ExtractSQLArg pulls the SQL string out of a loosely-typed tool argument map.
// The model occasionally emits the value as an array of strings or under an
// alternate key; this coerces all of those shapes into a single statement.
func ExtractSQLArg(args map[string]any) string {
	if args == nil {
		return ""
	}
	if v, ok := args["sql"]; ok {
		if s := coerceArgToString(v); s != "" {
			return s
		}
	}
	for _, v := range args {
		if s := coerceArgToString(v); s != "" {
			return s
		}
	}
	return ""
}

func coerceArgToString(v any) string {
	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val)
	case []any:
		parts := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				parts = append(parts, s)
			}
		}
		return strings.TrimSpace(strings.Join(parts, " "))
	case []string:
		return strings.TrimSpace(strings.Join(val, " "))
	default:
		return ""
	}
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
			switch v := vals[i].(type) {
			case []byte:
				m[c] = string(v)
			case time.Time:
				m[c] = v.UTC().Format(time.RFC3339)
			case nil, bool, int64, float64, string:
				m[c] = v
			default:
				m[c] = fmt.Sprintf("%v", v)
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
		func(ctx agent.Context, args map[string]any) ([]map[string]any, error) {
			rows, err := RunTursoQuery(ctx, db, ExtractSQLArg(args))
			if err != nil {
				return nil, err
			}
			capped, _ := capRows(rows)
			return capped, nil
		},
	)
}

const tursoQueryToolDescription = `Run a read-only SQL SELECT against the Fincher operational SQLite database and get matching rows.
Only a single SELECT (or WITH ... SELECT) statement is allowed. Writes/DDL and multiple statements are rejected. Results are capped at 100 rows (and a byte budget), so use aggregates (COUNT, GROUP BY) or tight WHERE/LIMIT clauses rather than selecting whole tables.

Schema:
- vendors(id TEXT, name TEXT, components TEXT_JSON, markets TEXT_JSON, hourly_rate_usd REAL, turnaround_hours INT)
- media_packages(id TEXT, title_id TEXT, component TEXT, language TEXT, version TEXT, vendor_id TEXT, derived_from_master_version TEXT, redelivery_count INT, status TEXT, market TEXT)
- deliveries(id TEXT, title_id TEXT, country TEXT, status TEXT, target_date DATETIME)
- titles(id TEXT, name TEXT, slug TEXT, type TEXT, premiere_date DATETIME, territories INT, current_master_version TEXT, overall_status TEXT)
- masters(id TEXT, title_id TEXT, version TEXT, supersedes_version TEXT, created_at DATETIME)

Agent audit trail (records WHAT the agent did and WHY — use these to answer "why"/"when" questions):
- runs(id TEXT, title_slug TEXT, trigger TEXT, status TEXT, started_at DATETIME, ended_at DATETIME, created_at DATETIME) — one row per agent workflow execution. status in (PENDING,RUNNING,COMPLETED,FAILED,ESCALATED).
- steps(id TEXT, run_id TEXT, name TEXT, status TEXT, started_at DATETIME, ended_at DATETIME, metadata TEXT_JSON, created_at DATETIME) — individual nodes within a run.
- wf_results(id TEXT, run_id TEXT, step_id TEXT, judge TEXT, outcome TEXT, rationale TEXT, attempt INT, created_at DATETIME) — each judge/verdict produced during a run. rationale is the recorded natural-language reasoning = the literal "why the agent decided this". Join back to runs via run_id (and optionally steps via step_id).

Examples:
- Most recent run for a title and its verdicts (the "why"):
  SELECT r.id, r.trigger, r.status, r.started_at FROM runs r WHERE r.title_slug = '<slug>' ORDER BY r.started_at DESC LIMIT 1
  SELECT judge, outcome, rationale, attempt FROM wf_results WHERE run_id = '<run_id>' ORDER BY created_at
- Titles premiering in a date window:
  SELECT name, slug, premiere_date, overall_status FROM titles WHERE premiere_date BETWEEN '2026-09-09' AND '2026-09-10' ORDER BY premiere_date

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
