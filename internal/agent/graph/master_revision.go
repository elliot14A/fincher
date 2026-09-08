package graph

import (
	"context"
	"strings"

	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

type masterRevisionResult struct {
	NewVersion            string
	InvalidatedPackages   int
	InvalidatedPackageIDs []string
}

func applyMasterRevision(ctx context.Context, deps IncidentGraphDeps, event *models.Event) (masterRevisionResult, error) {
	var out masterRevisionResult
	if deps.TursoClient == nil || event == nil {
		return out, nil
	}

	newVersion := ""
	if event.Data != nil {
		if v, ok := event.Data["new_master_version"].(string); ok {
			newVersion = strings.TrimSpace(v)
		}
	}
	if newVersion == "" {
		return out, nil
	}

	tRes := tursotitles.FindByIDOrSlug(ctx, deps.TursoClient, event.Subject)
	if tRes.IsErr() {
		return out, tRes.Error()
	}
	title := tRes.Unwrap()

	if strings.EqualFold(title.CurrentMasterVersion, newVersion) {
		return out, nil
	}

	supersedes := title.CurrentMasterVersion
	out.NewVersion = newVersion

	masterID := models.MasterIDFor(title.Slug, newVersion)
	mRes := tursomasters.Create(ctx, deps.TursoClient, &models.Master{
		ID:                masterID,
		TitleID:           title.ID,
		Version:           newVersion,
		SupersedesVersion: supersedes,
	})
	if mRes.IsErr() {
		logger.Warn("master-revision: failed to create master row",
			"title_slug", title.Slug, "master_id", masterID, "error", mRes.Error())
	}

	pkgRes := tursopackages.List(ctx, deps.TursoClient, tursopackages.ListFilter{
		TitleID: domainerrors.Some(title.ID),
	}, models.Pagination{Page: 1, Limit: 500})
	if pkgRes.IsErr() {
		logger.Warn("master-revision: failed to list packages, bumping version without invalidation",
			"title_slug", title.Slug, "error", pkgRes.Error())
	} else {
		invalidated := models.PackageStatusInvalidated
		for _, p := range pkgRes.Unwrap().Items {
			if p.IsStaleAgainst(newVersion) && p.Status != models.PackageStatusInvalidated {
				if r := tursopackages.Update(ctx, deps.TursoClient, p.ID, &models.UpdatePackageInput{
					Status: &invalidated,
				}); r.IsErr() {
					logger.Warn("master-revision: failed to invalidate stale package",
						"package_id", p.ID, "error", r.Error())
					continue
				}
				out.InvalidatedPackages++
				out.InvalidatedPackageIDs = append(out.InvalidatedPackageIDs, p.ID)
			}
		}
	}

	titleUpd := &models.UpdateTitleInput{CurrentMasterVersion: &newVersion}
	if out.InvalidatedPackages > 0 {
		processing := models.StatusProcessing
		titleUpd.OverallStatus = &processing
	}
	uRes := tursotitles.Update(ctx, deps.TursoClient, title.ID, titleUpd)
	if uRes.IsErr() {
		return out, uRes.Error()
	}

	if out.InvalidatedPackages > 0 {
		delRes := tursodeliveries.List(ctx, deps.TursoClient, tursodeliveries.ListFilter{
			TitleID: domainerrors.Some(title.ID),
		}, models.Pagination{Page: 1, Limit: 500})
		if delRes.IsOk() {
			hold := models.DeliveryStatusHold
			for _, d := range delRes.Unwrap().Items {
				if d.Status == models.DeliveryStatusReadyToShip {
					if r := tursodeliveries.Update(ctx, deps.TursoClient, d.ID, &models.UpdateDeliveryInput{
						Status: &hold,
					}); r.IsErr() {
						logger.Warn("master-revision: failed to downgrade delivery to HOLD",
							"delivery_id", d.ID, "error", r.Error())
					}
				}
			}
		}
	}

	logger.Info("master-revision: applied new master version",
		"title_slug", title.Slug,
		"new_version", newVersion,
		"supersedes", supersedes,
		"invalidated_packages", out.InvalidatedPackages,
	)
	return out, nil
}
