package api

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/api/deliveries"
	"github.com/elliot14A/fincher/internal/api/dependencies"
	"github.com/elliot14A/fincher/internal/api/events"
	"github.com/elliot14A/fincher/internal/api/masters"
	"github.com/elliot14A/fincher/internal/api/packages"
	"github.com/elliot14A/fincher/internal/api/runs"
	"github.com/elliot14A/fincher/internal/api/titles"
	"github.com/elliot14A/fincher/internal/api/uploads"
	"github.com/elliot14A/fincher/internal/api/vendors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/openapi"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/web"
)

type Server struct {
	echo   *echo.Echo
	client *ent.Client
	chDB   *sql.DB
	llm    model.LLM
}

func (s *Server) SetModel(m model.LLM) {
	s.llm = m
}

func NewServer(client *ent.Client, chDB ...*sql.DB) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CORS())

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
	if len(chDB) > 0 {
		s.chDB = chDB[0]
	}

	s.registerRoutes()
	return s
}

func (s *Server) Router() *echo.Echo {
	return s.echo
}

func (s *Server) registerRoutes() {
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

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
	uploads.RegisterRoutes(apiGroup.Group("/uploads"), s.client)
	runs.RegisterRoutes(apiGroup.Group("/runs"), s.client, s.chDB, func() model.LLM { return s.llm })

	if s.chDB != nil {
		events.RegisterRoutes(apiGroup.Group("/events"), s.chDB, s.client, func() model.LLM { return s.llm })
	}

	web.RegisterRoutes(s.echo)
}
