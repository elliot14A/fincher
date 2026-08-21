package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create handles POST /api/vendors.
//
//	@Summary		Create a post-production vendor
//	@Description	Registers a localization, dubbing, or subtitling facility.
//	@Tags			vendors
//	@Accept			json
//	@Produce		json
//	@Param			vendor	body		models.Vendor	true	"Vendor payload"
//	@Success		201		{object}	models.Vendor
//	@Failure		400		{object}	errors.DomainError
//	@Router			/vendors [post]
func Create(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Vendor
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		res := tursovendors.Create(c.Request().Context(), client, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusCreated, res.Unwrap())
	}
}
