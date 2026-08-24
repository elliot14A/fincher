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
//	@Description	Fetches releases in the launch calendar with pagination, optionally filtered by status.
//	@Tags			titles
//	@Produce		json
//	@Param			status		query		string	false	"Status filter (ON_TRACK, AT_RISK, HOLD, PROCESSING, SHIPPED)"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			sort_order	query		string	false	"Sort order (asc, desc)"
//	@Param			search		query		string	false	"Search query"
//	@Success		200			{object}	models.TitlePaginationResult
//	@Failure		500			{object}	errors.ErrorResponse
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

		p := models.ParsePagination(
			c.QueryParam("page"),
			c.QueryParam("limit"),
			c.QueryParam("sort_order"),
			c.QueryParam("search"),
		)

		res := tursotitles.List(c.Request().Context(), client, filter, p)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
