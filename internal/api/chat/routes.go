package chat

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/mcp"
)

// RegisterRoutes registers the chat assistant endpoints on the given router group.
func RegisterRoutes(g *echo.Group, client *ent.Client, tursoDB *sql.DB, mcpClient *mcp.Client, modelProvider func() model.LLM) {
	g.POST("", Create(client, tursoDB, mcpClient, modelProvider))
	g.GET("", List(client))
	g.GET("/:id", Get(client))
	g.PATCH("/:id", Update(client))
	g.DELETE("/:id", Delete(client))
}
