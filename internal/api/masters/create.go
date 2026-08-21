package masters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursomasters "github.com/elliot14A/fincher/internal/turso/masters"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create handles POST /api/masters.
//
//	@Summary		Create an immutable master cut
//	@Description	Publishes a new master cut and updates Title.current_master_version.
//	@Tags			masters
//	@Accept			json
//	@Produce		json
//	@Param			master	body		models.Master	true	"Master payload"
//	@Success		201		{object}	models.Master
//	@Failure		400		{object}	errors.DomainError
//	@Router			/masters [post]
func Create(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Master
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		res := tursomasters.Create(c.Request().Context(), client, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusCreated, res.Unwrap())
	}
}
