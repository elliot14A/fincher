package dependencies

import (
	"github.com/labstack/echo/v4"

	"github.com/elliot14A/fincher/internal/turso/ent"
)

// RegisterRoutes registers all dependency endpoints on the given router group.
func RegisterRoutes(g *echo.Group, client *ent.Client) {
	g.POST("", Create(client))
	g.GET("", List(client))
	g.GET("/graph/:title_id", Graph(client))
	g.DELETE("/:id", Delete(client))
}
