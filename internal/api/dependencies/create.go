package dependencies

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create handles POST /api/dependencies.
//
//	@Summary		Create an immutable dependency edge
//	@Description	Registers a lineage edge between two packages with cycle prevention and same-title checks.
//	@Tags			dependencies
//	@Accept			json
//	@Produce		json
//	@Param			dependency	body		models.Dependency	true	"Dependency payload"
//	@Success		201			{object}	models.Dependency
//	@Failure		400			{object}	errors.DomainError
//	@Failure		409			{object}	errors.DomainError
//	@Router			/dependencies [post]
func Create(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Dependency
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		res := dependencies.Create(c.Request().Context(), client, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusCreated, res.Unwrap())
	}
}
