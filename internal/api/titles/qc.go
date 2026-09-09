package titles

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/agent/graph"
	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/api/events"
	chevents "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/mcp"
)

// SendToQC handles POST /api/titles/:id/qc.
//
//	@Summary		Send Title for Master QC
//	@Description	Initiates master cut quality control inspection and dispatches localization allocation for a title.
//	@Tags			titles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Title ID"
//	@Success		200	{object}	models.Title
//	@Failure		404	{object}	errors.DomainError
//	@Failure		409	{object}	errors.ErrorResponse
//	@Router			/titles/{id}/qc [post]
func SendToQC(client *ent.Client, chDB *sql.DB, mcpClient *mcp.Client, tursoDB *sql.DB, modelProvider func() model.LLM, sched *scheduler.Scheduler) echo.HandlerFunc {
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

		// 2. Schedule compressed-time Master QC task and arm premiere deadline if scheduler is active
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
			ArmTitleDeadline(client, chDB, mcpClient, tursoDB, modelProvider, sched, updatedTitle)
		}

		// 3. Dispatch Autonomous Multi-Agent Allocation Workflow
		if modelProvider != nil && modelProvider() != nil {
			var requirements []models.AllocationRequirement
			markets := updatedTitle.Metadata["markets"]
			if marketList, ok := markets.([]any); ok {
				for _, m := range marketList {
					mStr := fmt.Sprintf("%v", m)
					requirements = append(requirements,
						models.AllocationRequirement{Component: "AUDIO", Market: mStr, Language: mStr},
						models.AllocationRequirement{Component: "SUBTITLE", Market: mStr, Language: mStr},
					)
				}
			}
			if len(requirements) == 0 {
				requirements = append(requirements,
					models.AllocationRequirement{Component: "AUDIO", Market: "en-US", Language: "en-US"},
					models.AllocationRequirement{Component: "SUBTITLE", Market: "en-US", Language: "en-US"},
				)
			}

			hoursUntil := tursotitles.ResolveHoursUntilPremiere(ctx, client, updatedTitle.Slug, 0)
			allocDeps := graph.AllocationGraphDeps{
				Model:       modelProvider(),
				TursoClient: client,
				ClickHouse:  chDB,
				MCP:         mcpClient,
				TursoDB:     tursoDB,
				Scheduler:   sched,
				OnScheduleComplete: func(qcEvent models.Event) {
					bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
					defer cancel()
					_, err := events.IngestAndRoute(bgCtx, chDB, mcpClient, tursoDB, client, modelProvider, []models.Event{qcEvent}, sched)
					if err != nil {
						logger.Error("titles/qc: failed to re-ingest allocation-scheduled QC event",
							"title_slug", updatedTitle.Slug,
							"event_id", qcEvent.ID,
							"event_type", qcEvent.Type,
							"error", err,
						)
					}
				},
			}

			go func() {
				bgCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
				defer cancel()
				runID := "run-" + uuid.NewString()[:8]
				_, _, err := graph.DispatchAllocation(bgCtx, allocDeps, graph.AllocationInput{
					RunID:              runID,
					TitleSlug:          updatedTitle.Slug,
					Requirements:       requirements,
					HoursUntilPremiere: hoursUntil,
				})
				if err != nil {
					logger.Warn("titles/qc: failed to dispatch allocation workflow for title",
						"title_slug", updatedTitle.Slug,
						"error", err,
					)
				} else {
					logger.Info("titles/qc: dispatched allocation workflow for title",
						"title_slug", updatedTitle.Slug,
						"run_id", runID,
					)
				}
			}()
		}

		return c.JSON(http.StatusOK, updatedTitle)
	}
}
