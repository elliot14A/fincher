package tools

import (
	"context"
	"database/sql"

	chvendors "github.com/elliot14A/fincher/internal/clickhouse/vendors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

// VendorCandidatesArgs defines filtering parameters for vendor candidates.
type VendorCandidatesArgs struct {
	Specialty string `json:"specialty,omitempty"`
	Component string `json:"component,omitempty"`
}

// FetchVendorCandidates queries candidate vendors from Turso SQLite and enriches with ClickHouse accuracy.
func FetchVendorCandidates(ctx context.Context, client *ent.Client, chDB *sql.DB, args VendorCandidatesArgs) ([]models.VendorCandidate, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.FetchVendorCandidates", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}

	filter := domainerrors.None[string]()
	if args.Specialty != "" {
		filter = domainerrors.Some(args.Specialty)
	}

	p := models.Pagination{
		Page:  1,
		Limit: 50,
	}

	listRes := tursovendors.List(ctx, client, filter, p)
	if listRes.IsErr() {
		return nil, listRes.Error()
	}

	vendorsList := listRes.Unwrap().Items
	candidates := make([]models.VendorCandidate, 0, len(vendorsList))

	for _, v := range vendorsList {
		accuracy := models.UnmeasuredHistoricalAccuracy
		if chDB != nil && args.Component != "" {
			accRes := chvendors.RecencyWeightedAccuracy(ctx, chDB, v.ID, args.Component)
			if accRes.IsOk() {
				accuracy = accRes.Unwrap()
			}
		}

		candidates = append(candidates, models.VendorCandidate{
			VendorID:           v.ID,
			VendorName:         v.Name,
			Specialty:          v.Specialty,
			HourlyRateUSD:      v.HourlyRateUSD,
			TurnaroundHours:    v.TurnaroundHours,
			HistoricalAccuracy: accuracy,
		})
	}

	return candidates, nil
}

// NewVendorCandidatesTool creates an ADK tool wrapping FetchVendorCandidates.
func NewVendorCandidatesTool(client *ent.Client, chDB *sql.DB) (tool.Tool, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.NewVendorCandidatesTool", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	return functiontool.New(
		functiontool.Config{
			Name:        "list_vendor_candidates",
			Description: "Lists available vendor candidates from the commercial rate card with turnaround hours and historical accuracy.",
		},
		func(ctx agent.Context, args VendorCandidatesArgs) ([]models.VendorCandidate, error) {
			return FetchVendorCandidates(ctx, client, chDB, args)
		},
	)
}
