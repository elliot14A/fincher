package tools

import (
	"context"
	"database/sql"

	"github.com/elliot14A/fincher/internal/clickhouse/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

// AnalyticsArgs defines parameters for ClickHouse analytics inspection.
type AnalyticsArgs struct {
	VendorID  string `json:"vendor_id"`
	TitleSlug string `json:"title_slug"`
	Component string `json:"component"`
}

// FetchAnalytics compiles historical metrics and defect logs from ClickHouse.
func FetchAnalytics(ctx context.Context, db *sql.DB, args AnalyticsArgs) (*models.AnalyticsSummary, error) {
	if db == nil {
		return &models.AnalyticsSummary{
			VendorHistoricalAccuracy: models.UnmeasuredHistoricalAccuracy,
			SimilarDefectOccurrences: 0,
			PriorIncidentsForVendor:  0,
			RelevantHistoricalLogs:   []string{},
		}, nil
	}

	accuracy := models.UnmeasuredHistoricalAccuracy
	if args.VendorID != "" && args.Component != "" {
		accRes := vendors.RecencyWeightedAccuracy(ctx, db, args.VendorID, args.Component)
		if accRes.IsOk() {
			accuracy = accRes.Unwrap()
		}
	}

	similarCount := 0
	if args.TitleSlug != "" {
		row := db.QueryRowContext(ctx, `
			select count() from fincher.events
			where subject = ? and severity in ('WARN', 'CRITICAL')
		`, args.TitleSlug)
		_ = row.Scan(&similarCount)
	}

	vendorIncidents := 0
	if args.VendorID != "" {
		row := db.QueryRowContext(ctx, `
			select count() from fincher.events
			where JSONExtractString(data, 'vendor_id') = ? and severity in ('WARN', 'CRITICAL')
		`, args.VendorID)
		_ = row.Scan(&vendorIncidents)
	}

	var logs []string
	if args.TitleSlug != "" || args.VendorID != "" {
		rows, err := db.QueryContext(ctx, `
			select concat(toString(time), ' | ', type, ' | ', severity)
			from fincher.events
			where (subject = ? or JSONExtractString(data, 'vendor_id') = ?) and severity in ('WARN', 'CRITICAL')
			order by time desc
			limit 5
		`, args.TitleSlug, args.VendorID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var l string
				if err := rows.Scan(&l); err == nil {
					logs = append(logs, l)
				}
			}
			if err := rows.Err(); err != nil {
				return nil, err
			}
		}
	}

	return &models.AnalyticsSummary{
		VendorHistoricalAccuracy: accuracy,
		SimilarDefectOccurrences: similarCount,
		PriorIncidentsForVendor:  vendorIncidents,
		RelevantHistoricalLogs:   logs,
	}, nil
}

// NewAnalyticsTool creates an ADK tool wrapping FetchAnalytics.
func NewAnalyticsTool(db *sql.DB) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "query_analytics",
			Description: "Queries ClickHouse for historical vendor accuracy, defect recurrence, and past incident logs.",
		},
		func(ctx agent.Context, args AnalyticsArgs) (*models.AnalyticsSummary, error) {
			return FetchAnalytics(ctx, db, args)
		},
	)
}
