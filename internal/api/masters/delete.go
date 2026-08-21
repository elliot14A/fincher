package masters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
)

// Delete handles DELETE /api/masters/:id.
//
//	@Summary		Delete a master cut
//	@Description	Removes an unused master cut entry.
//	@Tags			masters
//	@Param			id	path	string	true	"Master ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	errors.DomainError
//	@Router			/masters/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursomasters.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
