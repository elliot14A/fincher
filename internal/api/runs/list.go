package runs

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursoruns "github.com/elliot14A/fincher/internal/turso/runs"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List handles GET /api/runs.
//
//	@Summary		List workflow runs
//	@Description	Fetches paginated agent workflow runs, optionally filtered by workflow trigger, status, or title slug.
//	@Tags			runs
//	@Produce		json
//	@Param			wf			query		string	false	"Workflow trigger filter (e.g. incident, allocation)"
//	@Param			status		query		string	false	"Status filter (PENDING, RUNNING, COMPLETED, FAILED, ESCALATED)"
//	@Param			title_slug	query		string	false	"Title slug filter"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			sort_order	query		string	false	"Sort order (asc, desc)"
//	@Param			search		query		string	false	"Search query"
//	@Success		200			{object}	models.RunPaginationResult
//	@Failure		500			{object}	errors.ErrorResponse
//	@Router			/runs [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter tursoruns.ListFilter

		if wf := c.QueryParam("wf"); wf != "" {
			filter.Workflow = domainerrors.Some(wf)
		} else {
			filter.Workflow = domainerrors.None[string]()
		}

		if stat := c.QueryParam("status"); stat != "" {
			filter.Status = domainerrors.Some(models.RunStatus(stat))
		} else {
			filter.Status = domainerrors.None[models.RunStatus]()
		}

		if tSlug := c.QueryParam("title_slug"); tSlug != "" {
			filter.TitleSlug = domainerrors.Some(tSlug)
		} else {
			filter.TitleSlug = domainerrors.None[string]()
		}

		p := models.ParsePagination(
			c.QueryParam("page"),
			c.QueryParam("limit"),
			c.QueryParam("sort_order"),
			c.QueryParam("search"),
		)

		res := tursoruns.ListRuns(c.Request().Context(), client, filter, p)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
