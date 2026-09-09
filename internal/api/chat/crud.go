package chat

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	tursochats "github.com/elliot14A/fincher/internal/turso/chats"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// List handles GET /api/chat.
//
//	@Summary		List chat sessions
//	@Description	Returns all persisted chat sessions newest-first (without messages).
//	@Tags			chat
//	@Produce		json
//	@Success		200	{array}		models.ChatSession
//	@Failure		500	{object}	errors.ErrorResponse
//	@Router			/chat [get]
func List(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := tursochats.ListSessions(c.Request().Context(), client)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}
		return c.JSON(http.StatusOK, res.Unwrap())
	}
}

// Get handles GET /api/chat/:id.
//
//	@Summary		Get a chat session with messages
//	@Description	Fetches a single chat session and its full message history (the audit log of what the assistant did, including SQL citations).
//	@Tags			chat
//	@Produce		json
//	@Param			id	path		string	true	"Session ID"
//	@Success		200	{object}	models.ChatSession
//	@Failure		404	{object}	errors.ErrorResponse
//	@Router			/chat/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := tursochats.GetSession(c.Request().Context(), client, c.Param("id"))
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}
		return c.JSON(http.StatusOK, res.Unwrap())
	}
}

// UpdateSessionRequest is the body for PATCH /api/chat/:id.
type UpdateSessionRequest struct {
	Title string `json:"title"`
}

// Update handles PATCH /api/chat/:id.
//
//	@Summary		Rename a chat session
//	@Description	Updates the title of a chat session.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Session ID"
//	@Param			request	body		chat.UpdateSessionRequest	true	"New title"
//	@Success		200		{object}	models.ChatSession
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		404		{object}	errors.ErrorResponse
//	@Router			/chat/{id} [patch]
func Update(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UpdateSessionRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}
		res := tursochats.UpdateTitle(c.Request().Context(), client, c.Param("id"), req.Title)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}
		return c.JSON(http.StatusOK, res.Unwrap())
	}
}

// Delete handles DELETE /api/chat/:id.
//
//	@Summary		Delete a chat session
//	@Description	Deletes a chat session and all of its messages.
//	@Tags			chat
//	@Produce		json
//	@Param			id	path	string	true	"Session ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	errors.ErrorResponse
//	@Router			/chat/{id} [delete]
func Delete(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := tursochats.DeleteSession(c.Request().Context(), client, c.Param("id"))
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}
		return c.NoContent(http.StatusNoContent)
	}
}
