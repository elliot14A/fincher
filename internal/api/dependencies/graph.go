package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// Graph handles GET /api/dependencies/graph/:title_id.
//
//	@Summary		Fetch lineage DAG tree for a title
//	@Description	Traverses packages and dependency edges to build a recursive hierarchical tree from root video masters.
//	@Tags			dependencies
//	@Produce		json
//	@Param			title_id	path		string	true	"Title ID"
//	@Success		200			{object}	models.LineageGraph
//	@Failure		404			{object}	errors.DomainError
//	@Router			/dependencies/graph/{title_id} [get]
func Graph(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		titleID := c.Param("title_id")
		res := dependencies.GetLineageGraph(c.Request().Context(), client, titleID)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
