package runs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/elliot14A/fincher/internal/turso/ent"
	tursoruns "github.com/elliot14A/fincher/internal/turso/runs"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Stream handles GET /api/runs/:id/stream.
//
//	@Summary		Stream live run updates via SSE
//	@Description	Establishes a Server-Sent Events (SSE) connection streaming real-time step progressions and status transitions for a run.
//	@Tags			runs
//	@Produce		text/event-stream
//	@Param			id	path	string	true	"Run ID"
//	@Success		200	{string}	string	"Stream of events: event: update\\ndata: {...}\\n\\n"
//	@Failure		404	{object}	errors.ErrorResponse
//	@Router			/runs/{id}/stream [get]
func Stream(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		ctx := c.Request().Context()

		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
		c.Response().Header().Set("X-Accel-Buffering", "no")
		c.Response().WriteHeader(http.StatusOK)

		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()

		var lastStepCount int = -1
		var lastStatus models.RunStatus = ""

		for {
			select {
			case <-ctx.Done():
				return nil

			case <-ticker.C:
				res := tursoruns.GetRun(ctx, client, id)
				if res.IsErr() {
					continue
				}
				run := res.Unwrap()

				hasChanged := len(run.Steps) != lastStepCount || run.Status != lastStatus
				if hasChanged {
					lastStepCount = len(run.Steps)
					lastStatus = run.Status

					data, err := json.Marshal(run)
					if err == nil {
						fmt.Fprintf(c.Response(), "event: update\ndata: %s\n\n", data)
						c.Response().Flush()
					}
				}

				if run.Status == models.RunStatusCompleted || run.Status == models.RunStatusFailed || run.Status == models.RunStatusEscalated {
					fmt.Fprintf(c.Response(), "event: done\ndata: {}\n\n")
					c.Response().Flush()
					return nil
				}
			}
		}
	}
}
