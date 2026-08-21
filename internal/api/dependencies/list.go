package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// List handles GET /api/dependencies.
//
//	@Summary		List all dependency edges
//	@Description	Fetches dependency edges, optionally filtered by parent or child package ID.
//	@Tags			dependencies
//	@Produce		json
//	@Param			parent_id	query		string	false	"Parent package ID"
//	@Param			child_id	query		string	false	"Child package ID"
//	@Success		200			{array}		models.Dependency
//	@Failure		500			{object}	errors.DomainError
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

		res := dependencies.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
