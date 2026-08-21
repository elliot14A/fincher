package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
)

// Get handles GET /api/titles/:id.
//
//	@Summary		Get media title by ID
//	@Description	Fetches detailed launch metadata for a specific title.
//	@Tags			titles
//	@Produce		json
//	@Param			id	path		string	true	"Title ID"
//	@Success		200	{object}	models.Title
//	@Failure		404	{object}	errors.DomainError
//	@Router			/titles/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		res := tursotitles.Get(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
