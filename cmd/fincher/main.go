package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/api"
	"github.com/elliot14A/fincher/internal/clickhouse"
	"github.com/elliot14A/fincher/internal/config"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/mcp"
)

//	@title			Fincher Media Delivery Operations API
//	@version		1.0.0
//	@description	Autonomous delivery-integrity workflow engine for LUME streaming operations.
//	@BasePath		/api
//	@produce		json
//	@consume		json

var CLI struct {
	Config config.Config `kong:"embed"`
}

func main() {
	kctx := kong.Parse(&CLI,
		kong.Name("fincher"),
		kong.Description("Fincher — Autonomous Post-Production Incident Orchestrator"),
		kong.UsageOnError(),
	)

	cfg := &CLI.Config
	if err := cfg.Validate(); err != nil {
		kctx.FatalIfErrorf(fmt.Errorf("configuration validation failed: %w", err))
	}

	logger.Init(cfg.Environment, os.Stdout)
	logger.Info("starting fincher service", "environment", cfg.Environment, "port", cfg.Port)

	dbClient, dbSQL, err := turso.Open(cfg.TursoURL, cfg.TursoToken)
	if err != nil {
		logger.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}
	defer dbClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := turso.AutoMigrate(ctx, dbClient); err != nil {
		logger.Error("failed to execute database schema migrations", "error", err)
		return
	}

	chDB, err := clickhouse.Open(cfg.ClickHouseDSN)
	if err != nil {
		logger.Warn("failed to connect to clickhouse", "error", err)
	} else {
		defer chDB.Close()
		if err := clickhouse.AutoMigrate(ctx, chDB); err != nil {
			logger.Warn("failed to execute clickhouse schema migrations", "error", err)
		}
	}

	srv := api.NewServer(dbClient, chDB)
	srv.SetTursoDB(dbSQL)
	srv.SetScheduler(scheduler.NewScheduler(config.DefaultTimeScale))

	mcpClient, err := mcp.NewClient(cfg.MCPURL)
	if err != nil {
		logger.Error("failed to initialize clickhouse mcp client", "endpoint", cfg.MCPURL, "error", err)
		return
	}
	if err := mcpClient.Ping(ctx); err != nil {
		logger.Error("clickhouse mcp server unreachable; agent analytics require it", "endpoint", cfg.MCPURL, "error", err)
		return
	}
	srv.SetMCP(mcpClient)
	logger.Info("connected to clickhouse mcp server", "endpoint", cfg.MCPURL)
	if cfg.GeminiAPIKey != "" {
		modelRes := agent.NewModel(ctx, cfg.GeminiAPIKey, cfg.GeminiModel, cfg.GeminiOptions)
		if modelRes.IsErr() {
			logger.Warn("failed to initialize gemini model", "error", modelRes.Error())
		} else {
			srv.SetModel(modelRes.Unwrap())
			logger.Info("initialized gemini model runtime", "model", cfg.GeminiModel)
		}
	}
	e := srv.Router()

	addr := fmt.Sprintf(":%d", cfg.Port)
	e.Server = &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("api server listening", "address", addr)
		if err := e.StartServer(e.Server); err != nil && err != http.ErrServerClosed {
			logger.Error("server stopped unexpectedly", "error", err)
			stop <- syscall.SIGTERM
		}
	}()

	<-stop
	logger.Info("shutting down gracefully...")

	if sched := srv.Scheduler(); sched != nil {
		sched.Stop()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	} else {
		logger.Info("fincher service stopped")
	}
}
