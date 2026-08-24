package deliveries

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List handles GET /api/deliveries.
//
//	@Summary		List all territory deliveries
//	@Description	Fetches territory deliveries with pagination, optionally filtered by title, country, or status.
//	@Tags			deliveries
//	@Produce		json
//	@Param			title_id	query		string	false	"Title ID filter"
//	@Param			country		query		string	false	"Country code filter (e.g. US, ES, JP)"
//	@Param			status		query		string	false	"Status filter (PENDING, READY_TO_SHIP, HOLD, SHIPPED)"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			sort_order	query		string	false	"Sort order (asc, desc)"
//	@Param			search		query		string	false	"Search query"
//	@Success		200			{object}	models.DeliveryPaginationResult
//	@Failure		500			{object}	errors.ErrorResponse
//	@Router			/deliveries [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter deliveries.ListFilter

		if tID := c.QueryParam("title_id"); tID != "" {
			filter.TitleID = domainerrors.Some(tID)
		} else {
			filter.TitleID = domainerrors.None[string]()
		}

		if cCode := c.QueryParam("country"); cCode != "" {
			filter.Country = domainerrors.Some(cCode)
		} else {
			filter.Country = domainerrors.None[string]()
		}

		if stat := c.QueryParam("status"); stat != "" {
			filter.Status = domainerrors.Some(models.DeliveryStatus(stat))
		} else {
			filter.Status = domainerrors.None[models.DeliveryStatus]()
		}

		p := models.ParsePagination(
			c.QueryParam("page"),
			c.QueryParam("limit"),
			c.QueryParam("sort_order"),
			c.QueryParam("search"),
		)

		res := deliveries.List(c.Request().Context(), client, filter, p)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
