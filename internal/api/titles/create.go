package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create handles POST /titles.
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
