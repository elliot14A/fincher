package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type queryEnvelope struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
}

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9_.\-]+$`)

func SafeIdentifier(kind, v string) (string, error) {
	if v == "" {
		return "", fmt.Errorf("mcp: %s cannot be empty", kind)
	}
	if !identifierPattern.MatchString(v) {
		return "", fmt.Errorf("mcp: %s %q contains illegal characters", kind, v)
	}
	return v, nil
}

func RunQueryRows(ctx context.Context, c *Client, sql string) ([]map[string]any, error) {
	if c == nil {
		return nil, fmt.Errorf("mcp: client is nil")
	}

	raw, err := c.RunQuery(ctx, sql)
	if err != nil {
		return nil, err
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []map[string]any{}, nil
	}

	sanitized := nanInfPattern.ReplaceAllString(raw, "null")

	var env queryEnvelope
	if err := json.Unmarshal([]byte(sanitized), &env); err != nil {
		return nil, fmt.Errorf("mcp: failed to parse run_query response %q: %w", truncate(raw, 200), err)
	}

	out := make([]map[string]any, 0, len(env.Rows))
	for _, row := range env.Rows {
		m := make(map[string]any, len(env.Columns))
		for i, col := range env.Columns {
			if i < len(row) {
				m[col] = row[i]
			}
		}
		out = append(out, m)
	}
	return out, nil
}

func RunQueryScalar(ctx context.Context, c *Client, sql string) (any, error) {
	rows, err := RunQueryRows(ctx, c, sql)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	for _, v := range rows[0] {
		return v, nil
	}
	return nil, nil
}

var nanInfPattern = regexp.MustCompile(`\b(-?NaN|-?Inf(inity)?)\b`)

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
