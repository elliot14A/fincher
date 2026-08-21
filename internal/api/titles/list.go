package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	tursotitles "github.com/elliot14A/fincher/pkg/turso/titles"
)

// List handles GET /titles.
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		statusParam := c.QueryParam("status")
		var filter domainerrors.Option[models.TitleStatus]

		if statusParam != "" {
			filter = domainerrors.Some(models.TitleStatus(statusParam))
		} else {
			filter = domainerrors.None[models.TitleStatus]()
		}

		res := tursotitles.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
