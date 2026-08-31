package runs

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursoruns "github.com/elliot14A/fincher/internal/turso/runs"
)

// Get handles GET /api/runs/:id.
//
//	@Summary		Get workflow run by ID
//	@Description	Fetches a single workflow run by ID with all executed steps and policy results.
//	@Tags			runs
//	@Produce		json
//	@Param			id	path		string	true	"Run ID"
//	@Success		200	{object}	models.Run
//	@Failure		404	{object}	errors.ErrorResponse
//	@Failure		500	{object}	errors.ErrorResponse
//	@Router			/runs/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursoruns.GetRun(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
