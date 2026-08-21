package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
)

// Delete handles DELETE /api/titles/:id.
//
//	@Summary		Delete a media title
//	@Description	Removes a title if it has no dependent packages or masters.
//	@Tags			titles
//	@Param			id	path	string	true	"Title ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	errors.DomainError
//	@Failure		409	{object}	errors.DomainError
//	@Router			/titles/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursotitles.Delete(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
