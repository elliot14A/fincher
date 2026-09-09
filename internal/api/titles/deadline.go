package titles

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/api/events"
	chevents "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/scheduler"
	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/mcp"
)

// ArmTitleDeadline schedules a compressed-time deadline timer for a title.
func ArmTitleDeadline(
	client *ent.Client,
	chDB *sql.DB,
	mcpClient *mcp.Client,
	tursoDB *sql.DB,
	modelProvider func() model.LLM,
	sched *scheduler.Scheduler,
	title *models.Title,
) {
	if sched == nil || title == nil {
		return
	}

	_, _ = sched.ArmTitleDeadline(title.ID, title.Slug, title.PremiereDate, func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tRes := tursotitles.Get(bgCtx, client, title.ID)
		if tRes.IsErr() {
			return
		}
		currentTitle := tRes.Unwrap()
		if currentTitle.OverallStatus == models.StatusShipped || currentTitle.OverallStatus == models.StatusOverdue {
			return
		}

		pRes := tursopackages.List(bgCtx, client, tursopackages.ListFilter{
			TitleID: domainerrors.Some(title.ID),
		}, models.Pagination{Limit: 100})

		allValid := false
		if pRes.IsOk() {
			pkgs := pRes.Unwrap().Items
			if len(pkgs) > 0 {
				allValid = true
				for _, p := range pkgs {
					if p.Status != models.PackageStatusValid {
						allValid = false
						break
					}
				}
			}
		}

		if allValid {
			logger.Info("scheduler: title is fully ready at premiere deadline, marking SHIPPED", "title_slug", title.Slug)
			shippedStatus := models.StatusShipped
			_ = tursotitles.Update(bgCtx, client, title.ID, &models.UpdateTitleInput{
				OverallStatus: &shippedStatus,
			})

			dRes := tursodeliveries.List(bgCtx, client, tursodeliveries.ListFilter{
				TitleID: domainerrors.Some(title.ID),
			}, models.Pagination{Limit: 100})
			if dRes.IsOk() {
				shippedDelivery := models.DeliveryStatusShipped
				for _, d := range dRes.Unwrap().Items {
					_ = tursodeliveries.Update(bgCtx, client, d.ID, &models.UpdateDeliveryInput{
						Status: &shippedDelivery,
					})
				}
			}

			if chDB != nil {
				shipEv := models.Event{
					ID:       "evt-shipped-" + uuid.NewString()[:8],
					Source:   "fincher.scheduler",
					Type:     "fincher.title.shipped",
					Subject:  title.Slug,
					Time:     time.Now().UTC(),
					Severity: models.SeverityInfo,
					Data: map[string]any{
						"title_id":   title.ID,
						"title_slug": title.Slug,
						"status":     "SHIPPED",
					},
				}
				_ = chevents.InsertBatch(bgCtx, chDB, []models.Event{shipEv})
			}
			return
		}

		overdueStatus := models.StatusOverdue
		_ = tursotitles.Update(bgCtx, client, title.ID, &models.UpdateTitleInput{
			OverallStatus: &overdueStatus,
		})

		breachEvent := models.Event{
			ID:              "evt-deadline-" + uuid.NewString()[:8],
			Source:          "fincher.scheduler",
			Type:            models.TypeTitleDeadlineReached,
			Subject:         title.Slug,
			Time:            time.Now().UTC(),
			Severity:        models.SeverityCritical,
			DataContentType: "application/json",
			Data: map[string]any{
				"title_id":   title.ID,
				"title_slug": title.Slug,
				"stage":      "DEADLINE_BREACH",
			},
		}

		logger.Warn("scheduler: premiere deadline reached for unready title, emitting breach event",
			"title_id", title.ID,
			"title_slug", title.Slug,
		)

		if chDB != nil && client != nil {
			_, _ = events.IngestAndRoute(bgCtx, chDB, mcpClient, tursoDB, client, modelProvider, []models.Event{breachEvent}, sched)
		}
	})
}
