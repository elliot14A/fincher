package tools

import (
	"context"
	"strings"

	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

// DeliveryImpactArgs defines the target package and title for blast radius analysis.
type DeliveryImpactArgs struct {
	PackageID          string  `json:"package_id"`
	TitleID            string  `json:"title_id,omitempty"`
	HoursUntilPremiere float64 `json:"hours_until_premiere,omitempty"`
}

// FetchDeliveryImpact computes affected packages, deliveries, and premiere timeline from Turso.
func FetchDeliveryImpact(ctx context.Context, client *ent.Client, args DeliveryImpactArgs) (*models.DeliveryImpact, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.FetchDeliveryImpact", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}

	titleID := args.TitleID
	pkgLang := ""

	if args.PackageID != "" {
		pkgRes := tursopackages.Get(ctx, client, args.PackageID)
		if pkgRes.IsOk() {
			pkg := pkgRes.Unwrap()
			if titleID == "" {
				titleID = pkg.TitleID
			}
			pkgLang = pkg.Language
		}
	}

	hoursUntilPremiere := args.HoursUntilPremiere
	if hoursUntilPremiere <= 0 {
		hoursUntilPremiere = tursotitles.ResolveHoursUntilPremiere(ctx, client, titleID, 72.0)
	}

	affectedPackages := []string{}
	if args.PackageID != "" {
		affectedPackages = append(affectedPackages, args.PackageID)
	}

	affectedDeliveries := []string{}
	affectedMarkets := make(map[string]bool)

	if titleID != "" {
		delRes := tursodeliveries.List(ctx, client, tursodeliveries.ListFilter{
			TitleID: domainerrors.Some(titleID),
		}, models.Pagination{
			Page:  1,
			Limit: 100,
		})
		if delRes.IsOk() {
			for _, d := range delRes.Unwrap().Items {
				if pkgLang != "" {
					parts := strings.Split(pkgLang, "-")
					country := parts[len(parts)-1]
					if strings.EqualFold(d.Country, country) {
						affectedDeliveries = append(affectedDeliveries, d.ID)
						affectedMarkets[d.Country] = true
					}
				} else {
					affectedDeliveries = append(affectedDeliveries, d.ID)
					affectedMarkets[d.Country] = true
				}
			}
		}
	}

	marketList := make([]string, 0, len(affectedMarkets))
	for m := range affectedMarkets {
		marketList = append(marketList, m)
	}

	isUrgent := hoursUntilPremiere <= 72.0 && hoursUntilPremiere > 0

	return &models.DeliveryImpact{
		RootPackageID:      args.PackageID,
		AffectedPackages:   affectedPackages,
		AffectedDeliveries: affectedDeliveries,
		AffectedMarkets:    marketList,
		HoursUntilPremiere: hoursUntilPremiere,
		IsPremiereUrgent:   isUrgent,
	}, nil
}

// NewDeliveryImpactTool creates an ADK tool wrapping FetchDeliveryImpact.
func NewDeliveryImpactTool(client *ent.Client) (tool.Tool, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.NewDeliveryImpactTool", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	return functiontool.New(
		functiontool.Config{
			Name:        "get_delivery_impact",
			Description: "Calculates the blast radius of a defective package: affected downstream deliveries, markets, and premiere deadline.",
		},
		func(ctx agent.Context, args DeliveryImpactArgs) (*models.DeliveryImpact, error) {
			return FetchDeliveryImpact(ctx, client, args)
		},
	)
}
