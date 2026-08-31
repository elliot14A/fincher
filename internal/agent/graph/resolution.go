package graph

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	ch "github.com/elliot14A/fincher/internal/clickhouse/events"
	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// ResolutionDeps supplies dependencies to the resolution workflow.
type ResolutionDeps struct {
	TursoClient *ent.Client
	ClickHouse  *sql.DB
}

// ResolutionInput is the entry payload for clean-QC self-healing.
type ResolutionInput struct {
	RunID     string        `json:"run_id,omitempty"`
	Event     *models.Event `json:"event,omitempty"`
	PackageID string        `json:"package_id,omitempty"`
	TitleSlug string        `json:"title_slug,omitempty"`
}

// ResolutionOutput captures the entities unblocked by resolution.
type ResolutionOutput struct {
	PackageID           string   `json:"package_id"`
	ReleasedDeliveryIDs []string `json:"released_delivery_ids"`
	TitleStatus         string   `json:"title_status"`
}

// ensureRun guarantees the root Run record exists in Turso for direct invocations and dispatchers.
func ensureRun(ctx context.Context, client *ent.Client, runID, titleSlug, trigger string) (*models.Run, error) {
	if client == nil {
		return nil, nil
	}
	existing := runs.GetRun(ctx, client, runID)
	if existing.IsOk() {
		return existing.Unwrap(), nil
	}
	initialRun := &models.Run{
		Base:      models.Base{ID: runID},
		TitleSlug: titleSlug,
		Trigger:   trigger,
		Status:    models.RunStatusRunning,
		StartedAt: time.Now().UTC(),
	}
	createRes := runs.CreateRun(ctx, client, initialRun)
	if createRes.IsErr() {
		existingAfter := runs.GetRun(ctx, client, runID)
		if existingAfter.IsOk() {
			return existingAfter.Unwrap(), nil
		}
		return nil, createRes.Error()
	}
	return initialRun, nil
}

// ExecuteResolution handles clean QC returns: marks package VALID,
// releases dependent storefront deliveries whose required components (AUDIO + SUBTITLE) are all VALID,
// and applies three-way Title self-healing (HOLD / PROCESSING / ON_TRACK).
func ExecuteResolution(ctx context.Context, deps ResolutionDeps, input ResolutionInput) (*ResolutionOutput, error) {
	// 1. Resolve Package ID and Title Slug
	packageID := input.PackageID
	titleSlug := input.TitleSlug

	if packageID == "" && input.Event != nil && input.Event.Data != nil {
		if pid, ok := input.Event.Data["package_id"].(string); ok {
			packageID = pid
		}
	}
	if titleSlug == "" && input.Event != nil {
		titleSlug = input.Event.Subject
	}
	if titleSlug == "" {
		titleSlug = models.DefaultTitleAgnosticSentinel
	}

	if packageID == "" {
		return nil, domainerrors.NewWithOp("graph.ExecuteResolution", domainerrors.CodeInvalidInput, "target package_id is required", nil)
	}

	// 2. Resolve Run ID and ensure root Run record exists
	runID := input.RunID
	if runID == "" {
		if input.Event != nil && input.Event.ID != "" {
			runID = fmt.Sprintf("run-%s", input.Event.ID)
		} else {
			runID = "run-" + uuid.NewString()[:8]
		}
	}

	if deps.TursoClient != nil {
		if _, err := ensureRun(ctx, deps.TursoClient, runID, titleSlug, "resolution"); err != nil {
			return nil, fmt.Errorf("failed to ensure run record: %w", err)
		}
	}

	// Create Step row for live telemetry and streaming
	now := time.Now().UTC()
	stepID := fmt.Sprintf("step-%s-resolution", runID)
	if deps.TursoClient != nil {
		_ = runs.CreateStep(ctx, deps.TursoClient, &models.Step{
			Base:      models.Base{ID: stepID},
			RunID:     runID,
			Name:      "resolution_evaluation",
			Status:    models.StepStatusRunning,
			StartedAt: now,
		})
	}

	var releasedDeliveries []string
	var heldReasons []map[string]any
	finalTitleStatus := "UNKNOWN"

	if deps.TursoClient != nil {
		// ORDERING INVARIANT: The package VALID update MUST execute and persist in the database
		// BEFORE the title's package list is queried below, guaranteeing that the 1:N readiness check
		// evaluates against the newly valid package state.
		validStatus := models.PackageStatusValid
		pkgUpdRes := tursopackages.Update(ctx, deps.TursoClient, packageID, &models.UpdatePackageInput{
			Status: &validStatus,
		})
		if pkgUpdRes.IsErr() {
			return nil, fmt.Errorf("resolution aborted: failed to update package %s to VALID: %w", packageID, pkgUpdRes.Error())
		}

		// Find title object
		titleRes := tursotitles.FindByIDOrSlug(ctx, deps.TursoClient, titleSlug)
		if titleRes.IsErr() {
			logger.Warn("resolution: could not resolve title for self-healing", "title_slug", titleSlug, "error", titleRes.Error())
		} else {
			titleObj := titleRes.Unwrap()

			// Fetch all packages and deliveries for this title
			pkgListRes := tursopackages.List(ctx, deps.TursoClient, tursopackages.ListFilter{
				TitleID: domainerrors.Some(titleObj.ID),
			}, models.Pagination{Limit: 200})

			delListRes := tursodeliveries.List(ctx, deps.TursoClient, tursodeliveries.ListFilter{
				TitleID: domainerrors.Some(titleObj.ID),
			}, models.Pagination{Limit: 200})

			var allPackages []*models.Package
			if pkgListRes.IsOk() {
				allPackages = pkgListRes.Unwrap().Items
			}

			var allDeliveries []*models.Delivery
			if delListRes.IsOk() {
				allDeliveries = delListRes.Unwrap().Items
			}

			// Count held deliveries upfront for diagnostics
			heldDeliveryCount := 0
			for _, del := range allDeliveries {
				if del.Status == models.DeliveryStatusHold {
					heldDeliveryCount++
				}
			}

			// Task 2: Warn on localized packages with empty market if title has held deliveries
			if heldDeliveryCount > 0 {
				for _, pkg := range allPackages {
					if (pkg.Component == models.ComponentAudio || pkg.Component == models.ComponentSubtitle) && pkg.Market == "" {
						logger.Warn("resolution: localized package has empty market, cannot match any delivery",
							"title_id", titleObj.ID,
							"package_id", pkg.ID,
							"component", pkg.Component,
							"language", pkg.Language,
						)
					}
				}
			}

			// 4. Multi-Package Delivery Resolution (Strict 1:N Deterministic Market & Component Matching)
			// A delivery on HOLD is releasable iff:
			// - It has at least one AUDIO package and at least one SUBTITLE package matching del.Country / Market
			// - ALL relevant packages (Audio, Subtitle, global Video) are in VALID status.
			for _, del := range allDeliveries {
				if del.Status == models.DeliveryStatusHold {
					var relevantPackages []*models.Package
					hasAudio := false
					hasSubtitle := false

					for _, pkg := range allPackages {
						// Match market/country or global video
						isMarketMatch := pkg.Market != "" && strings.EqualFold(pkg.Market, del.Country)
						isVideo := pkg.Component == models.ComponentVideo

						if isMarketMatch || isVideo {
							relevantPackages = append(relevantPackages, pkg)
							if pkg.Component == models.ComponentAudio {
								hasAudio = true
							}
							if pkg.Component == models.ComponentSubtitle {
								hasSubtitle = true
							}
						}
					}

					// Required components: must have both AUDIO and SUBTITLE for the territory
					hasRequiredComponents := hasAudio && hasSubtitle
					allValid := hasRequiredComponents && len(relevantPackages) > 0

					if allValid {
						for _, reqPkg := range relevantPackages {
							if reqPkg.Status != models.PackageStatusValid {
								allValid = false
								break
							}
						}
					}

					if allValid {
						readyStatus := models.DeliveryStatusReadyToShip
						updRes := tursodeliveries.Update(ctx, deps.TursoClient, del.ID, &models.UpdateDeliveryInput{
							Status: &readyStatus,
						})
						if updRes.IsOk() {
							releasedDeliveries = append(releasedDeliveries, del.ID)
							del.Status = models.DeliveryStatusReadyToShip
						}
					} else {
						// Task 1: Record structured diagnostics for held deliveries not released
						reason := ""
						if len(relevantPackages) == 0 {
							reason = "no_relevant_packages"
						} else if !hasAudio || !hasSubtitle {
							reason = "missing_required_component"
						} else {
							reason = "package_not_valid"
						}

						logger.Warn("resolution: held delivery not released",
							"delivery_id", del.ID,
							"country", del.Country,
							"reason", reason,
							"has_audio", hasAudio,
							"has_subtitle", hasSubtitle,
							"relevant_package_count", len(relevantPackages),
						)

						heldReasons = append(heldReasons, map[string]any{
							"delivery_id":            del.ID,
							"country":                del.Country,
							"reason":                 reason,
							"has_audio":              hasAudio,
							"has_subtitle":           hasSubtitle,
							"relevant_package_count": len(relevantPackages),
						})
					}
				}
			}

			// Task 4: Log summary tripwire if title has held deliveries but none were released
			if heldDeliveryCount > 0 && len(releasedDeliveries) == 0 {
				logger.Warn("resolution: no deliveries released for title with held deliveries",
					"title_id", titleObj.ID,
					"held_count", heldDeliveryCount,
					"resolved_package_id", packageID,
				)
			}

			// 5. Three-Way Title Self-Healing Logic:
			// - Stay HOLD if any delivery is still HOLD or any package is INVALIDATED/RE_QC_PENDING
			// - Set PROCESSING if any package is PENDING
			// - Set ON_TRACK only when all packages are VALID and all deliveries are clear of HOLD
			hasHoldDelivery := false
			for _, del := range allDeliveries {
				if del.Status == models.DeliveryStatusHold {
					hasHoldDelivery = true
					break
				}
			}

			hasBrokenPackage := false
			hasPendingPackage := false
			for _, pkg := range allPackages {
				if pkg.Status == models.PackageStatusInvalidated || pkg.Status == models.PackageStatusReQCPending {
					hasBrokenPackage = true
				} else if pkg.Status == models.PackageStatusPending {
					hasPendingPackage = true
				}
			}

			var targetStatus models.TitleStatus
			if hasHoldDelivery || hasBrokenPackage {
				targetStatus = models.StatusHold
			} else if hasPendingPackage {
				targetStatus = models.StatusProcessing
			} else {
				targetStatus = models.StatusOnTrack
			}

			_ = tursotitles.Update(ctx, deps.TursoClient, titleObj.ID, &models.UpdateTitleInput{
				OverallStatus: &targetStatus,
			})
			finalTitleStatus = string(targetStatus)
		}
	}

	// 6. Emit downstream fincher.delivery.released event to ClickHouse
	if deps.ClickHouse != nil {
		downstream := []models.Event{
			{
				ID:              "evt-" + uuid.NewString()[:8],
				Source:          "fincher/resolution",
				Type:            models.TypeDeliveryReleased,
				Subject:         titleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityInfo,
				DataContentType: "application/json",
				Data: map[string]any{
					"package_id":          packageID,
					"released_deliveries": releasedDeliveries,
					"title_status":        finalTitleStatus,
					"status":              "RESOLVED",
				},
			},
		}
		_ = ch.InsertBatch(ctx, deps.ClickHouse, downstream)
	}

	// 7. Persist completion Step and WfResult in Turso
	completedAt := time.Now().UTC()
	if deps.TursoClient != nil {
		stepMeta := map[string]any{
			"package_id":          packageID,
			"released_deliveries": releasedDeliveries,
			"title_status":        finalTitleStatus,
		}
		if len(heldReasons) > 0 {
			stepMeta["held_deliveries"] = heldReasons
		}

		_ = runs.UpdateStepStatus(ctx, deps.TursoClient, stepID, models.StepStatusCompleted, &completedAt, stepMeta)

		_ = runs.CreateResult(ctx, deps.TursoClient, &models.WfResult{
			Base:      models.Base{ID: fmt.Sprintf("res-%s-resolution", runID)},
			RunID:     runID,
			StepID:    stepID,
			Judge:     "resolution_engine",
			Outcome:   "RESOLVED",
			Rationale: fmt.Sprintf("Clean QC on package %s verified. Deliveries released: %v. Title %s -> %s.", packageID, releasedDeliveries, titleSlug, finalTitleStatus),
			Attempt:   1,
		})

		_ = runs.UpdateRunStatus(ctx, deps.TursoClient, runID, models.RunStatusCompleted, &completedAt, nil)
	}

	return &ResolutionOutput{
		PackageID:           packageID,
		ReleasedDeliveryIDs: releasedDeliveries,
		TitleStatus:         finalTitleStatus,
	}, nil
}
