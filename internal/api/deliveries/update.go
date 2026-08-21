package deliveries

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Update handles PATCH /api/deliveries/:id.
//
//	@Summary		Partial update of a delivery target
//	@Description	Updates territory release shipping state (e.g. HOLD, READY_TO_SHIP) or metadata.
//	@Tags			deliveries
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Delivery ID"
//	@Param			update	body		models.UpdateDeliveryInput	true	"Partial delivery update payload"
//	@Success		200		{object}	models.Delivery
//	@Failure		400		{object}	errors.DomainError
//	@Failure		404		{object}	errors.DomainError
//	@Router			/deliveries/{id} [patch]
func Update(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.UpdateDeliveryInput
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid update request body",
			})
		}

		res := deliveries.Update(c.Request().Context(), client, id, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
