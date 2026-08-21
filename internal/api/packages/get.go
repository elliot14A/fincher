package packages

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
)

// Get handles GET /api/packages/:id.
//
//	@Summary		Get media package by ID
//	@Description	Fetches package component, language, vendor, and derived master version.
//	@Tags			packages
//	@Produce		json
//	@Param			id	path		string	true	"Package ID"
//	@Success		200	{object}	models.Package
//	@Failure		404	{object}	errors.DomainError
//	@Router			/packages/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursopackages.Get(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
