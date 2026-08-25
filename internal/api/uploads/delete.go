package uploads

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursouploads "github.com/elliot14A/fincher/internal/turso/uploads"
)

// Delete handles DELETE /api/uploads/:id.
//
//	@Summary		Delete an upload
//	@Description	Removes an uploaded binary image from SQLite.
//	@Tags			uploads
//	@Produce		json
//	@Param			id	path		string	true	"Upload ID"
//	@Success		200	{object}	map[string]bool
//	@Failure		404	{object}	errors.DomainError
//	@Router			/uploads/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "missing upload id parameter",
			})
		}

		res := tursouploads.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, map[string]bool{"deleted": true})
	}
}
