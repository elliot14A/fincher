package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List handles GET /api/vendors.
//
//	@Summary		List all vendors
//	@Description	Fetches registered vendors with pagination, optionally filtered by specialty.
//	@Tags			vendors
//	@Produce		json
//	@Param			component	query		string	false	"Component filter (e.g. AUDIO, SUBTITLE, VIDEO)"
//	@Param			specialty	query		string	false	"Legacy specialty filter"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			sort_order	query		string	false	"Sort order (asc, desc)"
//	@Param			search		query		string	false	"Search query"
//	@Success		200			{object}	models.VendorPaginationResult
//	@Failure		500			{object}	errors.ErrorResponse
//	@Router			/vendors [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		comp := c.QueryParam("component")
		if comp == "" {
			comp = c.QueryParam("specialty")
		}
		var filter domainerrors.Option[string]
		if comp != "" {
			filter = domainerrors.Some(comp)
		} else {
			filter = domainerrors.None[string]()
		}

		p := models.ParsePagination(
			c.QueryParam("page"),
			c.QueryParam("limit"),
			c.QueryParam("sort_order"),
			c.QueryParam("search"),
		)

		res := tursovendors.List(c.Request().Context(), client, filter, p)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
