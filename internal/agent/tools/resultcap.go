package tools

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/elliot14A/fincher/pkg/logger"
)

// Hard limits on the tool output fed back to the LLM, to bound token usage.
const (
	maxResultRows      = 100
	maxResultBytes     = 24 * 1024
	clickhouseRowLimit = 100
)

var limitClausePattern = regexp.MustCompile(`(?i)\blimit\s+\d+`)

// capRows truncates a result set to a bounded number of rows AND a serialized
// byte budget, so a broad query can never flood the model context. It returns
// the (possibly trimmed) rows plus whether truncation occurred.
func capRows(rows []map[string]any) ([]map[string]any, bool) {
	truncated := false

	if len(rows) > maxResultRows {
		rows = rows[:maxResultRows]
		truncated = true
	}

	total := 0
	for i, row := range rows {
		b, err := json.Marshal(row)
		if err != nil {
			continue
		}
		total += len(b)
		if total > maxResultBytes {
			rows = rows[:i]
			truncated = true
			break
		}
	}

	if truncated {
		logger.Info("tool result capped", "returned_rows", len(rows), "byte_budget", maxResultBytes)
	}
	return rows, truncated
}

// ensureLimit appends a LIMIT clause to a raw SELECT that lacks one, so an
// unbounded ClickHouse scan cannot be issued.
func ensureLimit(sql string, limit int) string {
	if sql == "" {
		return sql
	}
	if limitClausePattern.MatchString(sql) {
		return sql
	}
	return sql + " LIMIT " + strconv.Itoa(limit)
}
