package packages

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create handles POST /api/packages.
//
//	@Summary		Create a media package
//	@Description	Registers a video, audio dub, or subtitle package derived from a master cut.
//	@Tags			packages
//	@Accept			json
//	@Produce		json
//	@Param			package	body		models.Package	true	"Package payload"
//	@Success		201		{object}	models.Package
//	@Failure		400		{object}	errors.DomainError
//	@Router			/packages [post]
func Create(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Package
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		res := tursopackages.Create(c.Request().Context(), client, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusCreated, res.Unwrap())
	}
}
