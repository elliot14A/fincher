package masters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
)

// Delete handles DELETE /masters/:id.
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
