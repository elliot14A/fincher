package events

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// RegisterRoutes registers the event batch ingestion and routing endpoint.
func RegisterRoutes(g *echo.Group, db *sql.DB, tursoClient *ent.Client, modelProvider func() model.LLM, sched *scheduler.Scheduler) {
	g.POST("", Create(db, tursoClient, modelProvider, sched))
}
