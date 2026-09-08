package graph

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/config"
	"github.com/elliot14A/fincher/internal/scheduler"
	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

type provisionResult struct {
	PackagesCreated   int
	DeliveriesCreated int
	QCScheduled       int
	SkippedInfeasible int
}

func provisionAllocationAssets(
	ctx context.Context,
	deps AllocationGraphDeps,
	titleSlug string,
	plan *models.AllocationPlan,
) (provisionResult, error) {
	var out provisionResult

	if deps.TursoClient == nil {
		return out, domainerrors.NewWithOp("graph.provisionAllocationAssets", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	if plan == nil || len(plan.Assignments) == 0 {
		return out, nil
	}

	tRes := tursotitles.FindByIDOrSlug(ctx, deps.TursoClient, titleSlug)
	if tRes.IsErr() {
		return out, domainerrors.NewWithOp("graph.provisionAllocationAssets", domainerrors.CodeInternal,
			fmt.Sprintf("failed to resolve title %q", titleSlug), tRes.Error())
	}
	title := tRes.Unwrap()

	masterVersion := title.CurrentMasterVersion
	if masterVersion == "" {
		masterVersion = "V01"
	}

	existingPkgKey := make(map[string]bool)
	if pl := tursopackages.List(ctx, deps.TursoClient, tursopackages.ListFilter{
		TitleID: domainerrors.Some(title.ID),
	}, models.Pagination{Page: 1, Limit: 500}); pl.IsOk() {
		for _, p := range pl.Unwrap().Items {
			existingPkgKey[string(p.Component)+"|"+p.Market] = true
		}
	}
	existingDeliveryCountry := make(map[string]bool)
	if dl := tursodeliveries.List(ctx, deps.TursoClient, tursodeliveries.ListFilter{
		TitleID: domainerrors.Some(title.ID),
	}, models.Pagination{Page: 1, Limit: 500}); dl.IsOk() {
		for _, d := range dl.Unwrap().Items {
			existingDeliveryCountry[d.Country] = true
		}
	}

	marketsSeen := make(map[string]bool)
	for i, a := range plan.Assignments {
		component := strings.ToUpper(strings.TrimSpace(a.Component))
		market := strings.TrimSpace(a.Market)
		language := strings.TrimSpace(a.Language)
		if language == "" {
			language = "en-US"
		}

		if market != "" {
			marketsSeen[market] = true
		}

		if strings.TrimSpace(a.WinnerVendorID) == "" {
			out.SkippedInfeasible++
			continue
		}

		pkgKey := component + "|" + market
		if existingPkgKey[pkgKey] {
			continue
		}

		pkgID := fmt.Sprintf("pkg-%s-%s-%d", title.Slug, strings.ToLower(component), i+1)
		if market != "" {
			pkgID = fmt.Sprintf("pkg-%s-%s-%s", title.Slug, strings.ToLower(component), strings.ToLower(market))
		}

		pkg := &models.Package{
			Base:                     models.Base{ID: pkgID},
			TitleID:                  title.ID,
			Component:                models.ComponentType(component),
			Language:                 language,
			Version:                  "v01",
			VendorID:                 a.WinnerVendorID,
			DerivedFromMasterVersion: masterVersion,
			RedeliveryCount:          0,
			Status:                   models.PackageStatusPending,
			Market:                   market,
		}
		cRes := tursopackages.Create(ctx, deps.TursoClient, pkg)
		if cRes.IsErr() {
			logger.Warn("allocation-executor: failed to create package",
				"title_slug", title.Slug, "package_id", pkgID, "error", cRes.Error())
			continue
		}
		existingPkgKey[pkgKey] = true
		out.PackagesCreated++

		if deps.Scheduler != nil && deps.OnScheduleComplete != nil {
			turnaround := config.DefaultTurnaroundHours
			if a.TurnaroundHours > 0 {
				turnaround = float64(a.TurnaroundHours)
			}
			_, schedErr := deps.Scheduler.ScheduleTask(
				scheduler.TaskKindPackage,
				pkgID,
				title.Slug,
				a.WinnerVendorID,
				models.ComponentType(component),
				"",
				turnaround,
				agent.BuildQCCompletionCallback(ctx, agent.QCScheduleDeps{
					TursoClient:        deps.TursoClient,
					Scheduler:          deps.Scheduler,
					OnScheduleComplete: deps.OnScheduleComplete,
					LogContext:         "allocation:" + title.Slug,
				}),
			)
			if schedErr != nil {
				logger.Warn("allocation-executor: failed to schedule QC for package",
					"title_slug", title.Slug, "package_id", pkgID, "error", schedErr)
			} else {
				out.QCScheduled++
			}
		}
	}

	targetDate := title.PremiereDate
	if targetDate.IsZero() {
		targetDate = time.Now().UTC().Add(72 * time.Hour)
	}
	for market := range marketsSeen {
		if existingDeliveryCountry[market] {
			continue
		}
		delID := fmt.Sprintf("del-%s-%s", title.Slug, strings.ToLower(market))
		del := &models.Delivery{
			Base:       models.Base{ID: delID},
			TitleID:    title.ID,
			Country:    market,
			Status:     models.DeliveryStatusHold,
			TargetDate: targetDate,
		}
		dRes := tursodeliveries.Create(ctx, deps.TursoClient, del)
		if dRes.IsErr() {
			logger.Warn("allocation-executor: failed to create delivery",
				"title_slug", title.Slug, "delivery_id", delID, "error", dRes.Error())
			continue
		}
		existingDeliveryCountry[market] = true
		out.DeliveriesCreated++
	}

	processing := models.StatusProcessing
	if uRes := tursotitles.Update(ctx, deps.TursoClient, title.ID, &models.UpdateTitleInput{
		OverallStatus: &processing,
	}); uRes.IsErr() {
		logger.Warn("allocation-executor: failed to set title PROCESSING",
			"title_slug", title.Slug, "error", uRes.Error())
	}

	return out, nil
}
