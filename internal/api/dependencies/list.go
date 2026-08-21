package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// List handles GET /dependencies.
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		filter := dependencies.ListFilter{
			ParentID: domainerrors.None[string](),
			ChildID:  domainerrors.None[string](),
		}

		if parentID := c.QueryParam("parent_id"); parentID != "" {
			filter.ParentID = domainerrors.Some(parentID)
		}
		if childID := c.QueryParam("child_id"); childID != "" {
			filter.ChildID = domainerrors.Some(childID)
		}

		res := dependencies.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
