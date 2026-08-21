package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
)

// Get handles GET /vendors/:id.
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursovendors.Get(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
