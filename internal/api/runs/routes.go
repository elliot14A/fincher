package runs

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/turso/ent"
)

// RegisterRoutes registers all workflow run endpoints on the given router group.
func RegisterRoutes(g *echo.Group, client *ent.Client, chDB *sql.DB, modelProvider func() model.LLM) {
	g.POST("", Create(client, chDB, modelProvider))
	g.GET("", List(client))
	g.GET("/:id", Get(client))
	g.GET("/:id/stream", Stream(client))
}
