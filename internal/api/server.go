package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/elliot14A/fincher/internal/api/deliveries"
	"github.com/elliot14A/fincher/internal/api/dependencies"
	"github.com/elliot14A/fincher/internal/api/masters"
	"github.com/elliot14A/fincher/internal/api/packages"
	"github.com/elliot14A/fincher/internal/api/titles"
	"github.com/elliot14A/fincher/internal/api/vendors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/logger"
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

	// Structured HTTP request logging using pkg/logger
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				logger.Error("http request error",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency_ms", v.Latency.Milliseconds(),
					"error", v.Error,
				)
			} else {
				logger.Info("http request",
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

// registerRoutes wires all API endpoints via their respective entity route modules.
func (s *Server) registerRoutes() {
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	titles.RegisterRoutes(s.echo.Group("/titles"), s.client)
	masters.RegisterRoutes(s.echo.Group("/masters"), s.client)
	vendors.RegisterRoutes(s.echo.Group("/vendors"), s.client)
	packages.RegisterRoutes(s.echo.Group("/packages"), s.client)
	deliveries.RegisterRoutes(s.echo.Group("/deliveries"), s.client)
	dependencies.RegisterRoutes(s.echo.Group("/dependencies"), s.client)

	logger.Debug("registered all api routes", "routes_count", len(s.echo.Routes()))
}
