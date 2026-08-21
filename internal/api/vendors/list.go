package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// List handles GET /vendors.
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		specialty := c.QueryParam("specialty")
		var filter domainerrors.Option[string]

		if specialty != "" {
			filter = domainerrors.Some(specialty)
		} else {
			filter = domainerrors.None[string]()
		}

		res := tursovendors.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
