package tools

import (
	"context"
	"fmt"

	"github.com/elliot14A/fincher/internal/clickhouse/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/mcp"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type AnalyticsArgs struct {
	VendorID  string `json:"vendor_id"`
	TitleSlug string `json:"title_slug"`
	Component string `json:"component"`
}

func FetchAnalytics(ctx context.Context, client *mcp.Client, args AnalyticsArgs) (*models.AnalyticsSummary, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.FetchAnalytics", domainerrors.CodeInvalidInput, "mcp client cannot be nil", nil)
	}

	accuracy := models.UnmeasuredHistoricalAccuracy
	if args.VendorID != "" && args.Component != "" {
		accRes := vendors.RecencyWeightedAccuracy(ctx, client, args.VendorID, args.Component)
		if accRes.IsOk() {
			accuracy = accRes.Unwrap()
		}
	}

	similarCount := 0
	if args.TitleSlug != "" {
		slug, err := mcp.SafeIdentifier("title_slug", args.TitleSlug)
		if err != nil {
			return nil, err
		}
		rows, err := mcp.RunQueryRows(ctx, client, fmt.Sprintf(
			`select count() as c from fincher.events where subject = '%s' and severity in ('WARN','CRITICAL')`, slug))
		if err != nil {
			return nil, err
		}
		if len(rows) > 0 {
			similarCount = asInt(rows[0]["c"])
		}
	}

	vendorIncidents := 0
	if args.VendorID != "" {
		vendorID, err := mcp.SafeIdentifier("vendor_id", args.VendorID)
		if err != nil {
			return nil, err
		}
		rows, err := mcp.RunQueryRows(ctx, client, fmt.Sprintf(
			`select count() as c from fincher.events where JSONExtractString(data,'vendor_id') = '%s' and severity in ('WARN','CRITICAL')`, vendorID))
		if err != nil {
			return nil, err
		}
		if len(rows) > 0 {
			vendorIncidents = asInt(rows[0]["c"])
		}
	}

	var logs []string
	if args.TitleSlug != "" || args.VendorID != "" {
		slug := args.TitleSlug
		vendorID := args.VendorID
		if slug != "" {
			if _, err := mcp.SafeIdentifier("title_slug", slug); err != nil {
				return nil, err
			}
		}
		if vendorID != "" {
			if _, err := mcp.SafeIdentifier("vendor_id", vendorID); err != nil {
				return nil, err
			}
		}
		rows, err := mcp.RunQueryRows(ctx, client, fmt.Sprintf(
			`select concat(toString(time),' | ',type,' | ',severity) as log
			 from fincher.events
			 where (subject = '%s' or JSONExtractString(data,'vendor_id') = '%s') and severity in ('WARN','CRITICAL')
			 order by time desc limit 5`, slug, vendorID))
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			if s, ok := r["log"].(string); ok {
				logs = append(logs, s)
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

func NewAnalyticsTool(client *mcp.Client) (tool.Tool, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.NewAnalyticsTool", domainerrors.CodeInvalidInput, "mcp client cannot be nil", nil)
	}
	return functiontool.New(
		functiontool.Config{
			Name:        "query_analytics",
			Description: "Queries ClickHouse via MCP for historical vendor accuracy, defect recurrence, and past incident logs.",
		},
		func(ctx agent.Context, args AnalyticsArgs) (*models.AnalyticsSummary, error) {
			return FetchAnalytics(ctx, client, args)
		},
	)
}

func asInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}
