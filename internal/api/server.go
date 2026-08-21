package api

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/elliot14A/fincher/internal/api/titles"
	"github.com/elliot14A/fincher/pkg/ent"
)

// Server holds the HTTP router and database client.
type Server struct {
	echo   *echo.Echo
	client *ent.Client
}

// NewServer initializes the Echo HTTP server and registers routes.
func NewServer(client *ent.Client) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Simple structured HTTP request logging
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				slog.Error("http request error",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency_ms", v.Latency.Milliseconds(),
					"error", v.Error,
				)
			} else {
				slog.Info("http request",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency_ms", v.Latency.Milliseconds(),
				)
			}
			return nil
		},
	}))

	s := &Server{
		echo:   e,
		client: client,
	}

	s.registerRoutes()
	return s
}

// Router returns the underlying Echo instance.
func (s *Server) Router() *echo.Echo {
	return s.echo
}

// registerRoutes wires all API endpoints.
func (s *Server) registerRoutes() {
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Title REST routes
	t := s.echo.Group("/titles")
	t.POST("", titles.Create(s.client))
	t.GET("", titles.List(s.client))
	t.GET("/:id", titles.Get(s.client))
	t.PATCH("/:id", titles.Update(s.client))
	t.DELETE("/:id", titles.Delete(s.client))
}
