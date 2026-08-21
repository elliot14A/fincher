package packages

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	tursopackages "github.com/elliot14A/fincher/pkg/turso/packages"
)

// List handles GET /packages.
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		filter := tursopackages.ListFilter{
			TitleID:   domainerrors.None[string](),
			VendorID:  domainerrors.None[string](),
			Component: domainerrors.None[models.ComponentType](),
			Status:    domainerrors.None[models.PackageStatus](),
		}

		if tID := c.QueryParam("title_id"); tID != "" {
			filter.TitleID = domainerrors.Some(tID)
		}
		if vID := c.QueryParam("vendor_id"); vID != "" {
			filter.VendorID = domainerrors.Some(vID)
		}
		if comp := c.QueryParam("component"); comp != "" {
			filter.Component = domainerrors.Some(models.ComponentType(comp))
		}
		if st := c.QueryParam("status"); st != "" {
			filter.Status = domainerrors.Some(models.PackageStatus(st))
		}

		res := tursopackages.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
