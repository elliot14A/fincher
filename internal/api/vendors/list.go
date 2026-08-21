package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// List handles GET /api/vendors.
//
//	@Summary		List all vendors
//	@Description	Fetches registered vendors, optionally filtered by specialty.
//	@Tags			vendors
//	@Produce		json
//	@Param			specialty	query		string	false	"Specialty filter (e.g. AUDIO_DUBBING, SUBTITLES)"
//	@Success		200			{array}		models.Vendor
//	@Failure		500			{object}	errors.DomainError
//	@Router			/vendors [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spec := c.QueryParam("specialty")
		var filter domainerrors.Option[string]
		if spec != "" {
			filter = domainerrors.Some(spec)
		} else {
			filter = domainerrors.None[string]()
		}

		res := tursovendors.List(c.Request().Context(), client, filter)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
