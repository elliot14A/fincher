package packages

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List handles GET /api/packages.
//
//	@Summary		List all media packages
//	@Description	Fetches media packages with optional filtering by title, vendor, component, or status.
//	@Tags			packages
//	@Produce		json
//	@Param			title_id	query		string	false	"Title ID filter"
//	@Param			vendor_id	query		string	false	"Vendor ID filter"
//	@Param			component	query		string	false	"Component filter (VIDEO, AUDIO, SUBTITLE, METADATA)"
//	@Param			status		query		string	false	"Status filter (PENDING, VALID, INVALIDATED, RE_QC_PENDING)"
//	@Success		200			{array}		models.Package
//	@Failure		500			{object}	errors.DomainError
//	@Router			/packages [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter tursopackages.ListFilter

		if tID := c.QueryParam("title_id"); tID != "" {
			filter.TitleID = domainerrors.Some(tID)
		} else {
			filter.TitleID = domainerrors.None[string]()
		}

		if vID := c.QueryParam("vendor_id"); vID != "" {
			filter.VendorID = domainerrors.Some(vID)
		} else {
			filter.VendorID = domainerrors.None[string]()
		}

		if comp := c.QueryParam("component"); comp != "" {
			filter.Component = domainerrors.Some(models.ComponentType(comp))
		} else {
			filter.Component = domainerrors.None[models.ComponentType]()
		}

		if stat := c.QueryParam("status"); stat != "" {
			filter.Status = domainerrors.Some(models.PackageStatus(stat))
		} else {
			filter.Status = domainerrors.None[models.PackageStatus]()
		}

		res := tursopackages.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
