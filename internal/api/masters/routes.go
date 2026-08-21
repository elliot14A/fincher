package masters

import (
	"github.com/labstack/echo/v4"

	"github.com/elliot14A/fincher/pkg/ent"
)

// RegisterRoutes registers all master endpoints on the given router group.
func RegisterRoutes(g *echo.Group, client *ent.Client) {
	g.POST("", Create(client))
	g.GET("", List(client))
	g.GET("/:id", Get(client))
	g.DELETE("/:id", Delete(client))
}
