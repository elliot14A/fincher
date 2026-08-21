package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
)

// Get handles GET /api/vendors/:id.
//
//	@Summary		Get vendor by ID
//	@Description	Fetches post facility metadata, specialty, and contact tags.
//	@Tags			vendors
//	@Produce		json
//	@Param			id	path		string	true	"Vendor ID"
//	@Success		200	{object}	models.Vendor
//	@Failure		404	{object}	errors.DomainError
//	@Router			/vendors/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursovendors.Get(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
