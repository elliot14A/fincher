package uploads

import (
	"github.com/labstack/echo/v4"

	"github.com/elliot14A/fincher/internal/turso/ent"
)

// RegisterRoutes mounts all upload endpoints on the provided Echo group.
func RegisterRoutes(g *echo.Group, client *ent.Client) {
	g.POST("", Upload(client))
	g.GET("/:id", Get(client))
	g.DELETE("/:id", Delete(client))
}
