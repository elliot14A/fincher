package titles

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// RegisterRoutes registers all title endpoints on the given router group.
func RegisterRoutes(g *echo.Group, client *ent.Client, chDB *sql.DB, modelProvider func() model.LLM, sched *scheduler.Scheduler) {
	g.POST("", Create(client, chDB, modelProvider, sched))
	g.GET("", List(client))
	g.GET("/:id", Get(client))
	g.PATCH("/:id", Update(client, chDB, modelProvider, sched))
	g.DELETE("/:id", Delete(client))
	g.POST("/:id/qc", SendToQC(client, chDB, modelProvider, sched))
}
