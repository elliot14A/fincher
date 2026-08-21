package packages

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Update handles PATCH /api/packages/:id.
//
//	@Summary		Partial update of a media package
//	@Description	Updates package status (e.g. RE_QC_PENDING, INVALIDATED), master version cut, or metadata.
//	@Tags			packages
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Package ID"
//	@Param			update	body		models.UpdatePackageInput	true	"Partial package update payload"
//	@Success		200		{object}	models.Package
//	@Failure		400		{object}	errors.DomainError
//	@Failure		404		{object}	errors.DomainError
//	@Router			/packages/{id} [patch]
func Update(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.UpdatePackageInput
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid update request body",
			})
		}

		res := tursopackages.Update(c.Request().Context(), client, id, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
