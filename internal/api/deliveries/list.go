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

// List handles GET /deliveries.
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		filter := deliveries.ListFilter{
			TitleID: domainerrors.None[string](),
			Country: domainerrors.None[string](),
			Status:  domainerrors.None[models.DeliveryStatus](),
		}

		if tID := c.QueryParam("title_id"); tID != "" {
			filter.TitleID = domainerrors.Some(tID)
		}
		if country := c.QueryParam("country"); country != "" {
			filter.Country = domainerrors.Some(country)
		}
		if status := c.QueryParam("status"); status != "" {
			filter.Status = domainerrors.Some(models.DeliveryStatus(status))
		}

		res := deliveries.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
