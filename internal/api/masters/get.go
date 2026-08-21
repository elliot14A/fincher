package masters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
)

// Get handles GET /api/masters/:id.
//
//	@Summary		Get master cut by ID
//	@Description	Fetches master cut version and supersedes metadata.
//	@Tags			masters
//	@Produce		json
//	@Param			id	path		string	true	"Master ID"
//	@Success		200	{object}	models.Master
//	@Failure		404	{object}	errors.DomainError
//	@Router			/masters/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursomasters.Get(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
