package titles

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Update handles PATCH /api/titles/:id.
//
//	@Summary		Partial update of a media title
//	@Description	Updates specific fields (overall_status, master version, metadata) on a title.
//	@Tags			titles
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Title ID"
//	@Param			update	body		models.UpdateTitleInput	true	"Partial title updates"
//	@Success		200		{object}	models.Title
//	@Failure		400		{object}	errors.DomainError
//	@Failure		404		{object}	errors.DomainError
//	@Router			/titles/{id} [patch]
func Update(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.UpdateTitleInput
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid update request body",
			})
		}

		res := tursotitles.Update(c.Request().Context(), client, id, &req)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		return c.JSON(http.StatusOK, res.Unwrap())
	}
}
