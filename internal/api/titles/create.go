package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create handles POST /api/titles.
//
//	@Summary		Create a media title
//	@Description	Inserts a new release title into the launch calendar.
//	@Tags			titles
//	@Accept			json
//	@Produce		json
//	@Param			title	body		models.Title	true	"Title payload"
//	@Success		201		{object}	models.Title
//	@Failure		400		{object}	errors.DomainError
//	@Router			/titles [post]
func Create(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Title
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		res := tursotitles.Create(c.Request().Context(), client, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusCreated, res.Unwrap())
	}
}
