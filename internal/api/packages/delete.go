package packages

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
)

// Delete handles DELETE /api/packages/:id.
//
//	@Summary		Delete a media package
//	@Description	Removes a media package if no child dependencies exist.
//	@Tags			packages
//	@Param			id	path	string	true	"Package ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	errors.DomainError
//	@Failure		409	{object}	errors.DomainError
//	@Router			/packages/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursopackages.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
