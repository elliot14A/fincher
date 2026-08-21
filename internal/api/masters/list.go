package masters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/ent"
	tursomasters "github.com/elliot14A/fincher/pkg/turso/masters"
)

// List handles GET /masters.
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		titleID := c.QueryParam("title_id")
		var filter domainerrors.Option[string]

		if titleID != "" {
			filter = domainerrors.Some(titleID)
		} else {
			filter = domainerrors.None[string]()
		}

		res := tursomasters.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
