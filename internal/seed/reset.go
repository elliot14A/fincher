package seed

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/logger"
)

// PreflightCheck inspects existing record counts across SQLite and ClickHouse tables and fails fast without --reset.
func PreflightCheck(ctx context.Context, client *ent.Client, db *sql.DB, reset bool) error {
	var tursoTitles int
	var tursoVendors int
	var chEvents int
	var chQC int
	var chMetrics int

	if client != nil {
		tursoTitles, _ = client.Title.Query().Count(ctx)
		tursoVendors, _ = client.Vendor.Query().Count(ctx)
	}

	if db != nil {
		_ = db.QueryRowContext(ctx, "SELECT count() FROM fincher.events").Scan(&chEvents)
		_ = db.QueryRowContext(ctx, "SELECT count() FROM fincher.qc").Scan(&chQC)
		_ = db.QueryRowContext(ctx, "SELECT count() FROM fincher.vendor_metrics").Scan(&chMetrics)
	}

	hasData := tursoTitles > 0 || tursoVendors > 0 || chEvents > 0 || chQC > 0 || chMetrics > 0
	if hasData && !reset {
		return fmt.Errorf("existing data found (Turso titles: %d, vendors: %d, ClickHouse events: %d, qc: %d, vendor_metrics: %d); pass --reset to wipe and re-seed",
			tursoTitles, tursoVendors, chEvents, chQC, chMetrics)
	}

	return nil
}

// ResetDatabases clears all Turso relational tables (in reverse FK order) and truncates ClickHouse event/MV tables.
func ResetDatabases(ctx context.Context, client *ent.Client, db *sql.DB) error {
	if client != nil {
		logger.Info("clearing Turso / SQLite tables in reverse FK order")

		// Reverse FK deletion order
		if _, err := client.Dependency.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear dependencies: %w", err)
		}
		if _, err := client.MediaPackage.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear media packages: %w", err)
		}
		if _, err := client.Delivery.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear deliveries: %w", err)
		}
		if _, err := client.Master.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear masters: %w", err)
		}
		if _, err := client.Title.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear titles: %w", err)
		}
		if _, err := client.Vendor.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear vendors: %w", err)
		}
		if _, err := client.WfResult.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear wf_results: %w", err)
		}
		if _, err := client.Step.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear steps: %w", err)
		}
		if _, err := client.Run.Delete().Exec(ctx); err != nil {
			return fmt.Errorf("failed to clear runs: %w", err)
		}
		logger.Info("successfully cleared Turso tables")
	}

	if db != nil {
		logger.Info("truncating ClickHouse events and materialized views")

		tables := []string{
			"fincher.events",
			"fincher.qc",
			"fincher.vendor_metrics",
		}

		for _, table := range tables {
			query := fmt.Sprintf("TRUNCATE TABLE %s", table)
			if _, err := db.ExecContext(ctx, query); err != nil {
				return fmt.Errorf("failed to truncate %s: %w", table, err)
			}
		}
		logger.Info("successfully truncated ClickHouse tables")
	}

	return nil
}
