package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List handles GET /api/dependencies.
//
//	@Summary		List all dependency edges
//	@Description	Fetches dependency edges with pagination, optionally filtered by parent or child package ID.
//	@Tags			dependencies
//	@Produce		json
//	@Param			parent_id	query		string	false	"Parent package ID"
//	@Param			child_id	query		string	false	"Child package ID"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			sort_order	query		string	false	"Sort order (asc, desc)"
//	@Param			search		query		string	false	"Search query"
//	@Success		200			{object}	models.DependencyPaginationResult
//	@Failure		500			{object}	errors.ErrorResponse
//	@Router			/dependencies [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter dependencies.ListFilter

		if pID := c.QueryParam("parent_id"); pID != "" {
			filter.ParentID = domainerrors.Some(pID)
		} else {
			filter.ParentID = domainerrors.None[string]()
		}

		if cID := c.QueryParam("child_id"); cID != "" {
			filter.ChildID = domainerrors.Some(cID)
		} else {
			filter.ChildID = domainerrors.None[string]()
		}

		p := models.ParsePagination(
			c.QueryParam("page"),
			c.QueryParam("limit"),
			c.QueryParam("sort_order"),
			c.QueryParam("search"),
		)

		res := dependencies.List(c.Request().Context(), client, filter, p)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
