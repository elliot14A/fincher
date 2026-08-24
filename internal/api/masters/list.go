package masters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List handles GET /api/masters.
//
//	@Summary		List all master cuts
//	@Description	Fetches master cut versions with pagination, optionally filtered by title_id.
//	@Tags			masters
//	@Produce		json
//	@Param			title_id	query		string	false	"Title ID filter"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			sort_order	query		string	false	"Sort order (asc, desc)"
//	@Param			search		query		string	false	"Search query"
//	@Success		200			{object}	models.MasterPaginationResult
//	@Failure		500			{object}	errors.ErrorResponse
//	@Router			/masters [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		titleID := c.QueryParam("title_id")
		var filter domainerrors.Option[string]
		if titleID != "" {
			filter = domainerrors.Some(titleID)
		} else {
			filter = domainerrors.None[string]()
		}

		p := models.ParsePagination(
			c.QueryParam("page"),
			c.QueryParam("limit"),
			c.QueryParam("sort_order"),
			c.QueryParam("search"),
		)

		res := tursomasters.List(c.Request().Context(), client, filter, p)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
