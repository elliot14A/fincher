package chat

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/agent"
	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	tursochats "github.com/elliot14A/fincher/internal/turso/chats"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/mcp"
)

// TaggedObjectRef is an @-mentioned object the operator attached to a message.
type TaggedObjectRef struct {
	Kind        string `json:"kind"`
	Identifier  string `json:"identifier"`
	DisplayName string `json:"display_name"`
}

// SendMessageRequest is the body for POST /api/chat.
type SendMessageRequest struct {
	Message   string            `json:"message"`
	SessionID string            `json:"session_id,omitempty"`
	Tagged    []TaggedObjectRef `json:"tagged_objects,omitempty"`
}

// SendMessageResponse is returned after the assistant answers.
type SendMessageResponse struct {
	SessionID string             `json:"session_id"`
	Message   models.ChatMessage `json:"message"`
}

func deriveTitle(message string) string {
	title := strings.TrimSpace(message)
	if len(title) > 60 {
		title = title[:60] + "..."
	}
	return title
}

func taggedContext(tagged []TaggedObjectRef) string {
	if len(tagged) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n\nThe operator referenced these specific objects (resolve @mentions to them):")
	for _, t := range tagged {
		b.WriteString("\n- ")
		b.WriteString(t.Kind)
		b.WriteString(" \"")
		b.WriteString(t.DisplayName)
		b.WriteString("\" (identifier: ")
		b.WriteString(t.Identifier)
		b.WriteString(")")
	}
	return b.String()
}

// Create handles POST /api/chat.
//
//	@Summary		Send a message to the chat assistant
//	@Description	Sends an operator message to the read-only chat assistant, persists the turn, and returns the grounded assistant answer with SQL citations. Creates a new session when session_id is omitted.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Param			request	body		chat.SendMessageRequest	true	"Chat message"
//	@Success		200		{object}	chat.SendMessageResponse
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		503		{object}	errors.ErrorResponse
//	@Router			/chat [post]
func Create(client *ent.Client, tursoDB *sql.DB, mcpClient *mcp.Client, modelProvider func() model.LLM) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SendMessageRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "invalid request body",
			})
		}
		if strings.TrimSpace(req.Message) == "" {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "message cannot be empty",
			})
		}

		var m model.LLM
		if modelProvider != nil {
			m = modelProvider()
		}
		if m == nil || mcpClient == nil || tursoDB == nil {
			return c.JSON(http.StatusServiceUnavailable, apierrors.ErrorResponse{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "AI chat runtime is not initialized (GEMINI_API_KEY + ClickHouse MCP required)",
			})
		}

		ctx := c.Request().Context()

		sessionID := req.SessionID
		if sessionID == "" {
			created := tursochats.CreateSession(ctx, client, deriveTitle(req.Message))
			if created.IsErr() {
				return apierrors.Respond(c, created.Error())
			}
			sessionID = created.Unwrap().ID
		} else {
			if existing := tursochats.GetSession(ctx, client, sessionID); existing.IsErr() {
				return apierrors.Respond(c, existing.Error())
			}
		}

		history := loadHistory(ctx, client, sessionID)

		if userMsg := tursochats.AppendMessage(ctx, client, sessionID, models.ChatRoleUser, req.Message, nil, nil); userMsg.IsErr() {
			return apierrors.Respond(c, userMsg.Error())
		}

		agentMessage := req.Message + taggedContext(req.Tagged)

		answer := agent.AnswerQuery(ctx, m, client, tursoDB, mcpClient, history, agentMessage)
		if answer.IsErr() {
			return apierrors.Respond(c, answer.Error())
		}
		result := answer.Unwrap()

		citations := make([]models.ChatCitation, len(result.Citations))
		for i, cit := range result.Citations {
			citations[i] = models.ChatCitation{Label: cit.Label, Source: cit.Source, Tool: cit.Tool}
		}

		sources := make([]models.ChatSource, len(result.Sources))
		for i, s := range result.Sources {
			sources[i] = models.ChatSource{
				Kind:        s.Kind,
				ID:          s.ID,
				DisplayName: s.DisplayName,
				Identifier:  s.Identifier,
				Subtitle:    s.Subtitle,
				ImageURL:    s.ImageURL,
				Status:      s.Status,
			}
		}

		assistantMsg := tursochats.AppendMessage(ctx, client, sessionID, models.ChatRoleAssistant, result.Answer, citations, sources)
		if assistantMsg.IsErr() {
			return apierrors.Respond(c, assistantMsg.Error())
		}

		return c.JSON(http.StatusOK, SendMessageResponse{
			SessionID: sessionID,
			Message:   *assistantMsg.Unwrap(),
		})
	}
}

func loadHistory(ctx context.Context, client *ent.Client, sessionID string) []agent.ChatTurn {
	res := tursochats.GetSession(ctx, client, sessionID)
	if res.IsErr() {
		return nil
	}
	messages := res.Unwrap().Messages
	turns := make([]agent.ChatTurn, 0, len(messages))
	for _, msg := range messages {
		turns = append(turns, agent.ChatTurn{
			Role:    agent.ChatRole(msg.Role),
			Content: msg.Content,
		})
	}
	return turns
}
