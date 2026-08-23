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
	"github.com/elliot14A/fincher/openapi"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/web"
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

// registerRoutes wires all API endpoints under /api via their respective entity route modules.
func (s *Server) registerRoutes() {
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Serve the code-first generated OpenAPI/Swagger specification
	s.echo.GET("/openapi.json", func(c echo.Context) error {
		return c.Blob(200, "application/json", openapi.SpecJSON)
	})

	apiGroup := s.echo.Group("/api")
	apiGroup.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	titles.RegisterRoutes(apiGroup.Group("/titles"), s.client)
	masters.RegisterRoutes(apiGroup.Group("/masters"), s.client)
	vendors.RegisterRoutes(apiGroup.Group("/vendors"), s.client)
	packages.RegisterRoutes(apiGroup.Group("/packages"), s.client)
	deliveries.RegisterRoutes(apiGroup.Group("/deliveries"), s.client)
	dependencies.RegisterRoutes(apiGroup.Group("/dependencies"), s.client)

	// Register embedded web MPA frontend and static assets
	web.RegisterRoutes(s.echo)
}
