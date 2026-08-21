package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	tursotitles "github.com/elliot14A/fincher/pkg/turso/titles"
)

// Update handles PATCH /titles/:id.
func Update(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.UpdateTitleInput
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}

		res := tursotitles.Update(c.Request().Context(), client, id, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
