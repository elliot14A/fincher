package deliveries

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// Delete handles DELETE /api/deliveries/:id.
//
//	@Summary		Delete a territory delivery target
//	@Description	Removes a territory delivery entry.
//	@Tags			deliveries
//	@Param			id	path	string	true	"Delivery ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	errors.DomainError
//	@Router			/deliveries/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := deliveries.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
