package titles

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	chevents "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// SendToQC handles POST /api/titles/:id/qc.
//
//	@Summary		Send Title for Master QC
//	@Description	Initiates master cut quality control inspection for a title. Rejects if the title is already in QC.
//	@Tags			titles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Title ID"
//	@Success		200	{object}	models.Title
//	@Failure		404	{object}	errors.DomainError
//	@Failure		409	{object}	errors.ErrorResponse
//	@Router			/titles/{id}/qc [post]
func SendToQC(client *ent.Client, chDB *sql.DB, _ func() model.LLM, sched *scheduler.Scheduler) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		ctx := c.Request().Context()

		getRes := tursotitles.Get(ctx, client, id)
		if getRes.IsErr() {
			return apierrors.Respond(c, getRes.Error())
		}
		title := getRes.Unwrap()

		// Guard: Reject if title is already undergoing QC
		if title.OverallStatus == models.StatusProcessing {
			return c.JSON(http.StatusConflict, apierrors.ErrorResponse{
				Code:    "ALREADY_IN_QC",
				Message: "Title is already in QC inspection",
			})
		}

		processingStatus := models.StatusProcessing
		updateRes := tursotitles.Update(ctx, client, id, &models.UpdateTitleInput{
			OverallStatus: &processingStatus,
		})
		if updateRes.IsErr() {
			return apierrors.Respond(c, updateRes.Error())
		}
		updatedTitle := updateRes.Unwrap()

		// 1. Emit Master QC Started Event to ClickHouse
		if chDB != nil {
			startEv := models.Event{
				ID:       "evt-" + uuid.NewString(),
				Type:     "fincher.master.qc.started",
				Source:   "fincher.titles.qc",
				Subject:  title.Slug,
				Time:     time.Now().UTC(),
				Severity: models.SeverityInfo,
				Data: map[string]any{
					"title_id":       title.ID,
					"title_slug":     title.Slug,
					"master_version": title.CurrentMasterVersion,
					"stage":          "MASTER_QC",
				},
			}
			if err := startEv.Validate(); err == nil {
				_ = chevents.InsertBatch(ctx, chDB, []models.Event{startEv})
			}
		}

		// 2. Schedule compressed-time Master QC task if scheduler is active
		if sched != nil {
			_, err := sched.ScheduleTask(
				scheduler.TaskKindMasterQC,
				title.ID,
				title.Slug,
				"vendor-deluxe-media",
				models.ComponentVideo,
				"PASSED",
				12.0, // 12h domain inspection -> 12s in compressed time
				func(t *scheduler.Task) {
					logger.Info("master qc: inspection task completed", "title_id", title.ID, "title_slug", title.Slug)
				},
			)
			if err != nil {
				logger.Warn("master qc: failed to schedule inspection task", "error", err)
			}
		}

		return c.JSON(http.StatusOK, updatedTitle)
	}
}
