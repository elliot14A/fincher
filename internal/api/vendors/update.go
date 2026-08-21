package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Update handles PATCH /api/vendors/:id.
//
//	@Summary		Partial update of a vendor
//	@Description	Updates vendor name, specialty, or metadata tags.
//	@Tags			vendors
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Vendor ID"
//	@Param			update	body		models.UpdateVendorInput	true	"Partial vendor update payload"
//	@Success		200		{object}	models.Vendor
//	@Failure		400		{object}	errors.DomainError
//	@Failure		404		{object}	errors.DomainError
//	@Router			/vendors/{id} [patch]
func Update(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.UpdateVendorInput
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid update request body",
			})
		}

		res := tursovendors.Update(c.Request().Context(), client, id, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
