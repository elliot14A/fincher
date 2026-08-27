package events

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	chEvents "github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create handles POST /api/events.
//
//	@Summary		Ingest event batch
//	@Description	Ingests an array of CloudEvents directly into ClickHouse.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			events	body		[]models.Event	true	"CloudEvents array"
//	@Success		201		{object}	models.EventBatchResponse
//	@Failure		400		{object}	errors.DomainError
//	@Router			/events [post]
func Create(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req []models.Event
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body: expected array of CloudEvents",
			})
		}

		if len(req) == 0 {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "event batch cannot be empty",
			})
		}

		// Validate all events upfront before writing to ClickHouse
		for i := range req {
			if err := req[i].Validate(); err != nil {
				return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
					Code:    "INVALID_INPUT",
					Message: fmt.Sprintf("event at index %d is invalid: %v", i, err),
				})
			}
		}

		ctx := c.Request().Context()
		for i := range req {
			res := chEvents.Insert(ctx, db, &req[i])
			if res.IsErr() {
				return apierrors.Respond(c, res.Error())
			}
		}

		return c.JSON(http.StatusCreated, models.EventBatchResponse{
			Status: "ingested",
			Count:  len(req),
		})
	}
}
