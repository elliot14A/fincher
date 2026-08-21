package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List handles GET /api/titles.
//
//	@Summary		List all media titles
//	@Description	Fetches all releases in the launch calendar, optionally filtered by status.
//	@Tags			titles
//	@Produce		json
//	@Param			status	query		string	false	"Status filter (ON_TRACK, AT_RISK, HOLD, PROCESSING, SHIPPED)"
//	@Success		200		{array}		models.Title
//	@Failure		500		{object}	errors.DomainError
//	@Router			/titles [get]
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
