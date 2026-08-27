package vendors

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elliot14A/fincher/internal/clickhouse"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

const DecayHalfLifeDays = 120.0

// RecencyWeightedAccuracy computes the 120-day exponential decay pass rate for a vendor component.
// If no measurements exist for this vendor component, it returns -1.0 (unmeasured sentinel).
func RecencyWeightedAccuracy(ctx context.Context, db *sql.DB, vendorID, component string) domainerrors.Result[float64] {
	if vendorID == "" || component == "" {
		return domainerrors.Err[float64](clickhouse.NewError("vendors.RecencyWeightedAccuracy", domainerrors.CodeInvalidInput, "vendorID and component are required", nil))
	}

	query := fmt.Sprintf(`
		select
			sum(failed_inspections * exp(-dateDiff('day', recorded_date, today()) / %f)) as weighted_failed,
			sum(measured_status_count * exp(-dateDiff('day', recorded_date, today()) / %f)) as weighted_measured
		from fincher.vendor_metrics
		where vendor_id = ? and component = ?
	`, DecayHalfLifeDays, DecayHalfLifeDays)

	var (
		weightedFailed   float64
		weightedMeasured float64
	)

	row := db.QueryRowContext(ctx, query, vendorID, component)
	if err := row.Scan(&weightedFailed, &weightedMeasured); err != nil {
		return domainerrors.Err[float64](clickhouse.MapError("vendors.RecencyWeightedAccuracy", "vendor_metrics", vendorID, err))
	}

	if weightedMeasured <= 0 {
		return domainerrors.Ok(-1.0)
	}

	accuracy := 1.0 - (weightedFailed / weightedMeasured)
	if accuracy < 0 {
		accuracy = 0
	}
	return domainerrors.Ok(accuracy)
}
