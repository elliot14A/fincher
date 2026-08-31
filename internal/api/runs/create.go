package runs

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/agent/graph"
	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursoruns "github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// CreateRunRequest defines the parameters for triggering a workflow run.
type CreateRunRequest struct {
	Workflow           string         `json:"workflow"`
	TitleSlug          string         `json:"title_slug"`
	HoursUntilPremiere float64        `json:"hours_until_premiere"`
	Component          string         `json:"component,omitempty"`
	Event              *models.Event  `json:"event,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

// Create handles POST /api/runs.
func Create(client *ent.Client, chDB *sql.DB, modelProvider func() model.LLM) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateRunRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		var m model.LLM
		if modelProvider != nil {
			m = modelProvider()
		}
		if m == nil {
			return c.JSON(http.StatusServiceUnavailable, apierrors.ErrorResponse{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "AI model runtime is not initialized (GEMINI_API_KEY required)",
			})
		}

		wf := c.QueryParam("wf")
		if wf == "" {
			wf = req.Workflow
		}
		if wf == "" {
			if req.Component != "" && req.Event == nil {
				wf = "allocation"
			} else {
				wf = "incident"
			}
		}

		titleSlug := req.TitleSlug
		if titleSlug == "" && req.Event != nil {
			titleSlug = req.Event.Subject
		}
		if titleSlug == "" {
			titleSlug = models.DefaultTitleAgnosticSentinel
		}

		req.HoursUntilPremiere = tursotitles.ResolveHoursUntilPremiere(c.Request().Context(), client, titleSlug, req.HoursUntilPremiere)

		runID := "run-" + uuid.NewString()[:8]
		if req.Event != nil && req.Event.ID != "" {
			runID = "run-" + req.Event.ID
			existing := tursoruns.GetRun(c.Request().Context(), client, runID)
			if existing.IsOk() {
				return c.JSON(http.StatusOK, existing.Unwrap())
			}
		}

		switch wf {
		case "allocation":
			allocDeps := graph.AllocationGraphDeps{
				Model:       m,
				TursoClient: client,
				ClickHouse:  chDB,
			}
			runObj, _, err := graph.DispatchAllocation(c.Request().Context(), allocDeps, graph.AllocationInput{
				RunID:              runID,
				TitleSlug:          titleSlug,
				Component:          req.Component,
				HoursUntilPremiere: req.HoursUntilPremiere,
			})
			if err != nil {
				return apierrors.Respond(c, err)
			}
			return c.JSON(http.StatusCreated, runObj)

		case "resolution":
			resDeps := graph.ResolutionDeps{
				TursoClient: client,
				ClickHouse:  chDB,
			}
			runObj, _, err := graph.DispatchResolution(c.Request().Context(), resDeps, graph.ResolutionInput{
				RunID:     runID,
				Event:     req.Event,
				TitleSlug: titleSlug,
			})
			if err != nil {
				return apierrors.Respond(c, err)
			}
			return c.JSON(http.StatusCreated, runObj)

		default: // incident
			event := req.Event
			if event == nil {
				eventID := runID
				if len(eventID) > 4 {
					eventID = eventID[4:]
				}
				event = &models.Event{
					ID:       eventID,
					Type:     models.TypeOperatorForced,
					Source:   "fincher/api/runs",
					Subject:  titleSlug,
					Time:     time.Now().UTC(),
					Severity: models.SeverityCritical,
					Data:     map[string]any{"component": req.Component},
				}
			}
			incidentDeps := graph.IncidentGraphDeps{
				Model:       m,
				TursoClient: client,
				ClickHouse:  chDB,
				MaxAttempts: graph.DefaultMaxRemediationAttempts,
			}
			runObj, _, err := graph.DispatchIncident(c.Request().Context(), incidentDeps, graph.IncidentInput{
				RunID:              runID,
				Event:              event,
				HoursUntilPremiere: req.HoursUntilPremiere,
			})
			if err != nil {
				return apierrors.Respond(c, err)
			}
			return c.JSON(http.StatusCreated, runObj)
		}
	}
}
