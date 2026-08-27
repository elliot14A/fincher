package vendors

import (
	"context"
	"database/sql"
	"time"

	"github.com/elliot14A/fincher/internal/clickhouse"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// MetricRow represents a daily aggregate rollup row from fincher.vendor_metrics.
type MetricRow struct {
	RecordedDate        time.Time `ch:"recorded_date"`
	VendorID            string    `ch:"vendor_id"`
	Component           string    `ch:"component"`
	TotalInspections    uint64    `ch:"total_inspections"`
	FailedInspections   uint64    `ch:"failed_inspections"`
	WarningInspections  uint64    `ch:"warning_inspections"`
	MeasuredStatusCount uint64    `ch:"measured_status_count"`
	TotalSyncDriftMs    float64   `ch:"total_sync_drift_ms"`
	MeasuredDriftCount  uint64    `ch:"measured_drift_count"`
}

// GetMetrics fetches daily aggregated metrics for a vendor across all components.
func GetMetrics(ctx context.Context, db *sql.DB, vendorID string) domainerrors.Result[[]MetricRow] {
	if vendorID == "" {
		return domainerrors.Err[[]MetricRow](clickhouse.NewError("vendors.GetMetrics", domainerrors.CodeInvalidInput, "vendorID is required", nil))
	}

	query := `
		select
			recorded_date,
			vendor_id,
			component,
			sum(total_inspections) as total_inspections,
			sum(failed_inspections) as failed_inspections,
			sum(warning_inspections) as warning_inspections,
			sum(measured_status_count) as measured_status_count,
			sum(total_sync_drift_ms) as total_sync_drift_ms,
			sum(measured_drift_count) as measured_drift_count
		from fincher.vendor_metrics
		where vendor_id = ?
		group by recorded_date, vendor_id, component
		order by recorded_date desc
	`

	rows, err := db.QueryContext(ctx, query, vendorID)
	if err != nil {
		return domainerrors.Err[[]MetricRow](clickhouse.MapError("vendors.GetMetrics", "vendor_metrics", vendorID, err))
	}
	defer rows.Close()

	var result []MetricRow
	for rows.Next() {
		var row MetricRow
		if err := rows.Scan(
			&row.RecordedDate,
			&row.VendorID,
			&row.Component,
			&row.TotalInspections,
			&row.FailedInspections,
			&row.WarningInspections,
			&row.MeasuredStatusCount,
			&row.TotalSyncDriftMs,
			&row.MeasuredDriftCount,
		); err != nil {
			return domainerrors.Err[[]MetricRow](clickhouse.MapError("vendors.GetMetrics", "vendor_metrics", vendorID, err))
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return domainerrors.Err[[]MetricRow](clickhouse.MapError("vendors.GetMetrics", "vendor_metrics", vendorID, err))
	}

	return domainerrors.Ok(result)
}
