package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// Delete handles DELETE /dependencies/:id.
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
