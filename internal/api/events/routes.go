package events

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers the event batch ingestion endpoint.
func RegisterRoutes(g *echo.Group, db *sql.DB) {
	g.POST("", Create(db))
}
