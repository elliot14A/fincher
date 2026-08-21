package vendors

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	tursovendors "github.com/elliot14A/fincher/pkg/turso/vendors"
)

// Create handles POST /vendors.
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
