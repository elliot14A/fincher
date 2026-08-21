package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
)

// Delete handles DELETE /api/vendors/:id.
//
//	@Summary		Delete a vendor
//	@Description	Removes a vendor if no dependent media packages exist.
//	@Tags			vendors
//	@Param			id	path	string	true	"Vendor ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	errors.DomainError
//	@Failure		409	{object}	errors.DomainError
//	@Router			/vendors/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursovendors.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
