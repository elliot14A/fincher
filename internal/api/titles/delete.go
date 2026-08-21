package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/pkg/ent"
	tursotitles "github.com/elliot14A/fincher/pkg/turso/titles"
)

// Delete handles DELETE /titles/:id.
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursotitles.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
