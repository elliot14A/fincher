package deliveries

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// Get handles GET /api/deliveries/:id.
//
//	@Summary		Get delivery target by ID
//	@Description	Fetches territory release status, carrier details, and target shipping date.
//	@Tags			deliveries
//	@Produce		json
//	@Param			id	path		string	true	"Delivery ID"
//	@Success		200	{object}	models.Delivery
//	@Failure		404	{object}	errors.DomainError
//	@Router			/deliveries/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := deliveries.Get(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
