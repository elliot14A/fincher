package vendors

import (
	"context"
	"fmt"

	"github.com/elliot14A/fincher/internal/clickhouse"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/mcp"
)

const DecayHalfLifeDays = 120.0

func RecencyWeightedAccuracy(ctx context.Context, client *mcp.Client, vendorID, component string) domainerrors.Result[float64] {
	safeVendor, err := mcp.SafeIdentifier("vendor_id", vendorID)
	if err != nil {
		return domainerrors.Err[float64](clickhouse.NewError("vendors.RecencyWeightedAccuracy", domainerrors.CodeInvalidInput, err.Error(), nil))
	}
	safeComponent, err := mcp.SafeIdentifier("component", component)
	if err != nil {
		return domainerrors.Err[float64](clickhouse.NewError("vendors.RecencyWeightedAccuracy", domainerrors.CodeInvalidInput, err.Error(), nil))
	}

	query := fmt.Sprintf(`
		select
			sum(failed_inspections * exp(-dateDiff('day', recorded_date, today()) / %f)) as weighted_failed,
			sum(measured_status_count * exp(-dateDiff('day', recorded_date, today()) / %f)) as weighted_measured
		from fincher.vendor_metrics
		where vendor_id = '%s' and component = '%s'
	`, DecayHalfLifeDays, DecayHalfLifeDays, safeVendor, safeComponent)

	rows, err := mcp.RunQueryRows(ctx, client, query)
	if err != nil {
		return domainerrors.Err[float64](clickhouse.MapError("vendors.RecencyWeightedAccuracy", "vendor_metrics", vendorID, err))
	}
	if len(rows) == 0 {
		return domainerrors.Ok(-1.0)
	}

	weightedFailed := asFloat(rows[0]["weighted_failed"])
	weightedMeasured := asFloat(rows[0]["weighted_measured"])

	if weightedMeasured <= 0 {
		return domainerrors.Ok(-1.0)
	}

	accuracy := 1.0 - (weightedFailed / weightedMeasured)
	if accuracy < 0 {
		accuracy = 0
	}
	return domainerrors.Ok(accuracy)
}

func asFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int64:
		return float64(n)
	case int:
		return float64(n)
	default:
		return 0
	}
}
