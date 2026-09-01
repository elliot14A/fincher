package events

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/agent/graph"
	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	chEvents "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

// IngestAndRoute validates, stores, and routes a batch of CloudEvents.
// It inserts all events directly into ClickHouse and asynchronously dispatches
// AI workflows (Incident investigation, Vendor allocation) based on event classification.
func IngestAndRoute(
	ctx context.Context,
	db *sql.DB,
	tursoClient *ent.Client,
	modelProvider func() model.LLM,
	events []models.Event,
	sched *scheduler.Scheduler,
) (*models.EventBatchResponse, error) {
	if db == nil {
		return nil, domainerrors.NewWithOp("events.IngestAndRoute", domainerrors.CodeInvalidInput, "clickhouse db connection cannot be nil", nil)
	}
	if tursoClient == nil {
		return nil, domainerrors.NewWithOp("events.IngestAndRoute", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	if len(events) == 0 {
		return nil, domainerrors.NewWithOp("events.IngestAndRoute", domainerrors.CodeInvalidInput, "event batch cannot be empty", nil)
	}

	// 1. Validate all events upfront before writing to ClickHouse
	for i := range events {
		if err := events[i].Validate(); err != nil {
			return nil, domainerrors.NewWithOp("events.IngestAndRoute", domainerrors.CodeInvalidInput, fmt.Sprintf("event at index %d is invalid: %v", i, err), err)
		}
	}

	// 2. Insert into ClickHouse
	res := chEvents.InsertBatch(ctx, db, events)
	if res.IsErr() {
		return nil, res.Error()
	}

	// 3. Inspect events and trigger appropriate workflows
	var runIDs []string
	var m model.LLM
	if modelProvider != nil {
		m = modelProvider()
	}

	for _, ev := range events {
		category := ev.Classify()

		switch category {
		case models.CategoryIncident:
			if m == nil {
				logger.Warn("incident event received but AI model runtime is not initialized, skipping dispatch", "event_id", ev.ID, "type", ev.Type)
				continue
			}
			incidentDeps := graph.IncidentGraphDeps{
				Model:       m,
				TursoClient: tursoClient,
				ClickHouse:  db,
				MaxAttempts: graph.DefaultMaxRemediationAttempts,
				Scheduler:   sched,
				OnScheduleComplete: func(qcEvent models.Event) {
					bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
					defer cancel()
					_, err := IngestAndRoute(bgCtx, db, tursoClient, modelProvider, []models.Event{qcEvent}, sched)
					if err != nil {
						logger.Error("events: failed to re-ingest scheduled QC completion event",
							"event_id", qcEvent.ID,
							"event_type", qcEvent.Type,
							"subject", qcEvent.Subject,
							"error", err,
						)
					}
				},
			}
			runObj, _, err := graph.DispatchIncident(ctx, incidentDeps, graph.IncidentInput{
				Event: &ev,
			})
			if err != nil {
				logger.Error("failed to dispatch incident from ingestion", "event_id", ev.ID, "error", err)
			} else if runObj != nil {
				runIDs = append(runIDs, runObj.ID)
			}

		case models.CategoryAllocation:
			if m == nil {
				logger.Warn("allocation event received but AI model runtime is not initialized, skipping dispatch", "event_id", ev.ID, "type", ev.Type)
				continue
			}
			component := "AUDIO"
			market := ""
			language := "en-US"
			if ev.Data != nil {
				if compStr, ok := ev.Data["component"].(string); ok && compStr != "" {
					component = compStr
				}
				if mStr, ok := ev.Data["market"].(string); ok && mStr != "" {
					market = mStr
				}
				if lStr, ok := ev.Data["language"].(string); ok && lStr != "" {
					language = lStr
				}
			}
			allocDeps := graph.AllocationGraphDeps{
				Model:       m,
				TursoClient: tursoClient,
				ClickHouse:  db,
			}
			runObj, _, err := graph.DispatchAllocation(ctx, allocDeps, graph.AllocationInput{
				RunID:     "run-" + ev.ID,
				TitleSlug: ev.Subject,
				Requirements: []models.AllocationRequirement{
					{
						Component: component,
						Market:    market,
						Language:  language,
					},
				},
			})
			if err != nil {
				logger.Error("failed to dispatch allocation from ingestion", "event_id", ev.ID, "error", err)
			} else if runObj != nil {
				runIDs = append(runIDs, runObj.ID)
			}

		case models.CategoryRoutineOutcome:
			// Clean QC inspection return triggers closed-loop resolution & self-healing
			if ev.Type == models.TypeQCInspectionCompleted && ev.Data != nil {
				if status, ok := ev.Data["status"].(string); ok && strings.ToUpper(status) == "PASSED" {
					resDeps := graph.ResolutionDeps{
						TursoClient: tursoClient,
						ClickHouse:  db,
					}
					runObj, _, err := graph.DispatchResolution(ctx, resDeps, graph.ResolutionInput{
						Event: &ev,
					})
					if err != nil {
						logger.Error("failed to dispatch resolution from ingestion", "event_id", ev.ID, "error", err)
					} else if runObj != nil {
						runIDs = append(runIDs, runObj.ID)
					}
				}
			}
		}
	}

	return &models.EventBatchResponse{
		Status: "ingested",
		Count:  len(events),
		RunIDs: runIDs,
	}, nil
}

// Create handles POST /api/events.
//
//	@Summary		Ingest event batch
//	@Description	Ingests an array of CloudEvents directly into ClickHouse and routes actionable events to agent workflows.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			events	body		[]models.Event	true	"CloudEvents array"
//	@Success		201		{object}	models.EventBatchResponse
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/events [post]
func Create(db *sql.DB, tursoClient *ent.Client, modelProvider func() model.LLM, sched *scheduler.Scheduler) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req []models.Event
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body: expected array of CloudEvents",
			})
		}

		ctx := c.Request().Context()
		resp, err := IngestAndRoute(ctx, db, tursoClient, modelProvider, req, sched)
		if err != nil {
			return apierrors.Respond(c, err)
		}

		return c.JSON(http.StatusCreated, resp)
	}
}
