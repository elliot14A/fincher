package tools

import (
	"context"
	"strings"

	chvendors "github.com/elliot14A/fincher/internal/clickhouse/vendors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/mcp"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type VendorCandidatesArgs struct {
	Component string `json:"component,omitempty"`
	Market    string `json:"market,omitempty"`
}

func FetchVendorCandidates(ctx context.Context, client *ent.Client, mcpClient *mcp.Client, args VendorCandidatesArgs) ([]models.VendorCandidate, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.FetchVendorCandidates", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}

	filter := domainerrors.None[string]()
	targetComp := strings.ToUpper(strings.TrimSpace(args.Component))
	if targetComp != "" {
		filter = domainerrors.Some(targetComp)
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

	targetMarket := strings.TrimSpace(args.Market)

	for _, v := range vendorsList {
		if targetComp != "" {
			compFound := false
			for _, c := range v.Components {
				if strings.EqualFold(c, targetComp) {
					compFound = true
					break
				}
			}
			if !compFound {
				continue
			}
		}

		if targetComp != "VIDEO" && targetMarket != "" {
			marketFound := false
			for _, m := range v.Markets {
				if strings.EqualFold(m, targetMarket) {
					marketFound = true
					break
				}
			}
			if !marketFound {
				continue
			}
		}

		accuracy := models.UnmeasuredHistoricalAccuracy
		if mcpClient != nil && targetComp != "" {
			accRes := chvendors.RecencyWeightedAccuracy(ctx, mcpClient, v.ID, targetComp)
			if accRes.IsOk() {
				accuracy = accRes.Unwrap()
			}
		}

		candidates = append(candidates, models.VendorCandidate{
			VendorID:           v.ID,
			VendorName:         v.Name,
			Components:         v.Components,
			Markets:            v.Markets,
			HourlyRateUSD:      v.HourlyRateUSD,
			TurnaroundHours:    v.TurnaroundHours,
			HistoricalAccuracy: accuracy,
		})
	}

	return candidates, nil
}

func NewVendorCandidatesTool(client *ent.Client, mcpClient *mcp.Client) (tool.Tool, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.NewVendorCandidatesTool", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}

	return functiontool.New(
		functiontool.Config{
			Name:        "get_vendor_candidates",
			Description: "Retrieves qualified vendor partners for a given component and optional market with commercial terms and historical quality ratings.",
		},
		func(ctx agent.Context, args VendorCandidatesArgs) ([]models.VendorCandidate, error) {
			return FetchVendorCandidates(ctx, client, mcpClient, args)
		},
	)
}
