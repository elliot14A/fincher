package masters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// List handles GET /api/masters.
//
//	@Summary		List all master cuts
//	@Description	Fetches master cut versions, optionally filtered by title_id.
//	@Tags			masters
//	@Produce		json
//	@Param			title_id	query		string	false	"Title ID filter"
//	@Success		200			{array}		models.Master
//	@Failure		500			{object}	errors.DomainError
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

		res := tursomasters.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
