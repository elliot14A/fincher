package titles

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/api/events"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// Update handles PATCH /api/titles/:id.
//
//	@Summary		Partial update of a media title
//	@Description	Updates specific fields (overall_status, master version, metadata, premiere_date) on a title. When premiere_date is updated, re-arms the deadline timer and re-drives the workflow.
//	@Tags			titles
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Title ID"
//	@Param			update	body		models.UpdateTitleInput	true	"Partial title updates"
//	@Success		200		{object}	models.Title
//	@Failure		400		{object}	errors.DomainError
//	@Failure		404		{object}	errors.DomainError
//	@Router			/titles/{id} [patch]
func Update(client *ent.Client, chDB *sql.DB, modelProvider func() model.LLM, sched *scheduler.Scheduler) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.UpdateTitleInput
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid update request body",
			})
		}

		ctx := c.Request().Context()
		res := tursotitles.Update(ctx, client, id, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}
		updated := res.Unwrap()

		if req.PremiereDate != nil {
			if updated.OverallStatus == models.StatusOverdue && req.PremiereDate.After(time.Now().UTC()) {
				procStatus := models.StatusProcessing
				updStatusRes := tursotitles.Update(ctx, client, id, &models.UpdateTitleInput{
					OverallStatus: &procStatus,
				})
				if updStatusRes.IsOk() {
					updated = updStatusRes.Unwrap()
				}
			}

			ArmTitleDeadline(client, chDB, modelProvider, sched, updated)

			if chDB != nil {
				resumeEv := models.Event{
					ID:              "evt-rescheduled-" + uuid.NewString()[:8],
					Source:          "fincher.titles.update",
					Type:            models.TypeInvestigationTriggered,
					Subject:         updated.Slug,
					Time:            time.Now().UTC(),
					Severity:        models.SeverityInfo,
					DataContentType: "application/json",
					Data: map[string]any{
						"title_id":      updated.ID,
						"title_slug":    updated.Slug,
						"premiere_date": updated.PremiereDate.Format(time.RFC3339),
						"reason":        "OPERATOR_RESCHEDULED_PREMIERE",
					},
				}
				bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				_, err := events.IngestAndRoute(bgCtx, chDB, client, modelProvider, []models.Event{resumeEv}, sched)
				if err != nil {
					logger.Warn("titles.Update: failed to route resume event on premiere date change", "error", err)
				}
			}
		}

		return c.JSON(http.StatusOK, updated)
	}
}
