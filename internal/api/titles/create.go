package titles

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/agent/graph"
	"github.com/elliot14A/fincher/internal/agent/scheduler"
	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	chevents "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// Create handles POST /api/titles.
//
//	@Summary		Create a media title
//	@Description	Inserts a new release title into the launch calendar, bundles Master V01, and triggers holistic allocation.
//	@Tags			titles
//	@Accept			json
//	@Produce		json
//	@Param			title	body		models.Title	true	"Title payload"
//	@Success		201		{object}	models.Title
//	@Failure		400		{object}	errors.DomainError
//	@Router			/titles [post]
func Create(client *ent.Client, chDB *sql.DB, modelProvider func() model.LLM, _ *scheduler.Scheduler) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Title
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		ctx := c.Request().Context()

		res := tursotitles.Create(ctx, client, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}
		created := res.Unwrap()

		masterVer := created.CurrentMasterVersion
		if masterVer == "" {
			masterVer = "V01"
		}
		masterID := fmt.Sprintf("mst-%s-%s", created.Slug, strings.ToLower(masterVer))
		master := &models.Master{
			ID:      masterID,
			TitleID: created.ID,
			Version: masterVer,
		}
		mRes := tursomasters.Create(ctx, client, master)
		if mRes.IsErr() {
			logger.Error("title onboarding: failed to bundle initial master",
				"title_id", created.ID,
				"title_slug", created.Slug,
				"error", mRes.Error(),
			)
		}

		requirements := []models.AllocationRequirement{
			{
				Component: "VIDEO",
				Market:    "",
				Language:  "en-US",
			},
		}

		var targetMarkets []string
		if created.Metadata != nil {
			if mList, ok := created.Metadata["markets"].([]string); ok && len(mList) > 0 {
				targetMarkets = mList
			} else if mAny, ok := created.Metadata["markets"].([]any); ok && len(mAny) > 0 {
				for _, item := range mAny {
					if str, ok := item.(string); ok && str != "" {
						targetMarkets = append(targetMarkets, str)
					}
				}
			}
		}
		if len(targetMarkets) == 0 {
			targetMarkets = []string{"en-US", "de-DE", "fr-FR", "hi-IN", "te-IN"}
		}

		for _, m := range targetMarkets {
			requirements = append(requirements, models.AllocationRequirement{
				Component: "AUDIO",
				Market:    m,
				Language:  m,
			})
			requirements = append(requirements, models.AllocationRequirement{
				Component: "SUBTITLE",
				Market:    m,
				Language:  m,
			})
		}

		if chDB != nil {
			reqJSON, _ := json.Marshal(requirements)
			ev := models.Event{
				ID:       "evt-" + uuid.NewString(),
				Type:     models.TypePackageRequired,
				Source:   "fincher.titles.onboarding",
				Subject:  created.Slug,
				Time:     time.Now().UTC(),
				Severity: models.SeverityInfo,
				Data: map[string]any{
					"title_id":     created.ID,
					"requirements": string(reqJSON),
				},
			}
			if insRes := chevents.InsertBatch(ctx, chDB, []models.Event{ev}); insRes.IsErr() {
				logger.Warn("title onboarding: failed to insert package.required audit event into clickhouse",
					"title_slug", created.Slug,
					"error", insRes.Error(),
				)
			}
		}

		if modelProvider != nil && modelProvider() != nil {
			deps := graph.AllocationGraphDeps{
				Model:       modelProvider(),
				TursoClient: client,
				ClickHouse:  chDB,
			}
			_, _, err := graph.DispatchAllocation(ctx, deps, graph.AllocationInput{
				TitleSlug:    created.Slug,
				Requirements: requirements,
			})
			if err != nil {
				logger.Warn("title onboarding: failed to dispatch allocation workflow",
					"title_slug", created.Slug,
					"error", err,
				)
			}
		} else {
			logger.Warn("title onboarding: allocation skipped (model provider unavailable)",
				"title_slug", created.Slug,
			)
		}

		return c.JSON(http.StatusCreated, created)
	}
}
