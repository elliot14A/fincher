package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// Delete handles DELETE /api/dependencies/:id.
//
//	@Summary		Delete a dependency edge
//	@Description	Removes a lineage edge between parent and child media packages.
//	@Tags			dependencies
//	@Param			id	path	string	true	"Dependency ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	errors.DomainError
//	@Router			/dependencies/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := dependencies.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
