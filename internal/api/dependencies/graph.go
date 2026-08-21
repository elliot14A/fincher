package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// Graph handles GET /dependencies/graph/:title_id.
func Graph(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		titleID := c.Param("title_id")
		res := dependencies.GetLineageGraph(c.Request().Context(), client, titleID)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
