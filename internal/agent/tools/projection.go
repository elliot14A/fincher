package tools

import (
	"context"
	"math"
	"time"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// TitleProjectionArgs defines the input schema for the title projection tool.
type TitleProjectionArgs struct {
	TitleSlug string `json:"title_slug" jsonschema:"description=The slug or ID of the title to project ready timeline for"`
}

// RepairProjection captures in-flight or estimated repair work for a package.
type RepairProjection struct {
	PackageID       string  `json:"package_id"`
	Component       string  `json:"component"`
	Market          string  `json:"market"`
	VendorID        string  `json:"vendor_id"`
	VendorName      string  `json:"vendor_name"`
	TurnaroundHours float64 `json:"turnaround_hours"`
	Status          string  `json:"status"`
}

// TitleProjection provides pre-computed timeline, critical path, buffer hours, and risk band.
type TitleProjection struct {
	TitleSlug              string             `json:"title_slug"`
	TitleName              string             `json:"title_name"`
	PremiereDate           string             `json:"premiere_date"`
	HoursUntilPremiere     float64            `json:"hours_until_premiere"`
	CriticalRemainingHours float64            `json:"critical_remaining_hours"`
	BufferHours            float64            `json:"buffer_hours"`
	RiskBand               string             `json:"risk_band"` // BREACH | TIGHT | WATCH | SAFE
	IsBreached             bool               `json:"is_breached"`
	IsUrgent               bool               `json:"is_urgent"`
	Repairs                []RepairProjection `json:"repairs"`
}

// GetTitleReadyProjection computes critical path and launch margin for a title.
func GetTitleReadyProjection(ctx context.Context, client *ent.Client, titleSlug string) (*TitleProjection, error) {
	if client == nil {
		return &TitleProjection{
			TitleSlug:  titleSlug,
			RiskBand:   "SAFE",
			IsBreached: false,
		}, nil
	}

	titleRes := tursotitles.FindByIDOrSlug(ctx, client, titleSlug)
	if titleRes.IsErr() {
		return nil, titleRes.Error()
	}
	titleObj := titleRes.Unwrap()

	now := time.Now().UTC()
	hoursUntilPremiere := math.Max(0, titleObj.PremiereDate.Sub(now).Hours())

	pkgListRes := tursopackages.List(ctx, client, tursopackages.ListFilter{
		TitleID: domainerrors.Some(titleObj.ID),
	}, models.Pagination{Limit: 200})

	var allPackages []*models.Package
	if pkgListRes.IsOk() {
		allPackages = pkgListRes.Unwrap().Items
	}

	var repairs []RepairProjection
	var maxTurnaround float64

	for _, pkg := range allPackages {
		// If package is not yet valid, it requires active work or QC
		if pkg.Status != models.PackageStatusValid {
			turnaround := 12.0 // domain default
			vendorName := "Unassigned"

			if pkg.VendorID != "" {
				vRes := tursovendors.Get(ctx, client, pkg.VendorID)
				if vRes.IsOk() {
					v := vRes.Unwrap()
					vendorName = v.Name
					if v.TurnaroundHours > 0 {
						turnaround = float64(v.TurnaroundHours)
					}
				}
			}

			repairs = append(repairs, RepairProjection{
				PackageID:       pkg.ID,
				Component:       string(pkg.Component),
				Market:          pkg.Market,
				VendorID:        pkg.VendorID,
				VendorName:      vendorName,
				TurnaroundHours: turnaround,
				Status:          string(pkg.Status),
			})

			if turnaround > maxTurnaround {
				maxTurnaround = turnaround
			}
		}
	}

	criticalRemainingHours := maxTurnaround
	bufferHours := hoursUntilPremiere - criticalRemainingHours

	var riskBand string
	isBreached := false
	isUrgent := false

	if bufferHours < 0 {
		riskBand = "BREACH"
		isBreached = true
		isUrgent = true
	} else if bufferHours < 6.0 {
		riskBand = "TIGHT"
		isUrgent = true
	} else if bufferHours < 24.0 {
		riskBand = "WATCH"
	} else {
		riskBand = "SAFE"
	}

	return &TitleProjection{
		TitleSlug:              titleObj.Slug,
		TitleName:              titleObj.Name,
		PremiereDate:           titleObj.PremiereDate.Format(time.RFC3339),
		HoursUntilPremiere:     math.Round(hoursUntilPremiere*10) / 10,
		CriticalRemainingHours: math.Round(criticalRemainingHours*10) / 10,
		BufferHours:            math.Round(bufferHours*10) / 10,
		RiskBand:               riskBand,
		IsBreached:             isBreached,
		IsUrgent:               isUrgent,
		Repairs:                repairs,
	}, nil
}

// NewProjectionTool constructs an ADK Go tool enabling agents to project title readiness.
func NewProjectionTool(client *ent.Client) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "get_title_ready_projection",
			Description: "Calculates hours until premiere, critical path remaining turnaround hours, buffer margin, and risk band for a title.",
		},
		func(ctx agent.Context, args TitleProjectionArgs) (*TitleProjection, error) {
			return GetTitleReadyProjection(ctx, client, args.TitleSlug)
		},
	)
}
