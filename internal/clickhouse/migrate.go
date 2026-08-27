package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/elliot14A/fincher/migrations"
	"github.com/elliot14A/fincher/pkg/logger"
)

// AutoMigrate executes all embedded ClickHouse migration scripts in lexicographical order.
func AutoMigrate(ctx context.Context, db *sql.DB) error {
	entries, err := fs.ReadDir(migrations.FS, "clickhouse")
	if err != nil {
		return fmt.Errorf("reading embedded migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		path := "clickhouse/" + file
		content, err := migrations.FS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading migration file %s: %w", file, err)
		}

		query := strings.TrimSpace(string(content))
		if query == "" {
			continue
		}

		logger.Debug("applying clickhouse migration", "file", file)
		if _, err := db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("executing migration %s: %w", file, err)
		}
	}

	logger.Info("clickhouse migrations completed successfully", "count", len(files))
	return nil
}
