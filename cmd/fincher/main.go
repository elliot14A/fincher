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

	"github.com/elliot14A/fincher/internal/api"
	"github.com/elliot14A/fincher/internal/config"
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/pkg/logger"
)

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

	// Initialize Turso database connection
	dbClient, err := turso.Open(cfg.TursoURL, cfg.TursoToken)
	if err != nil {
		logger.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}
	defer dbClient.Close()

	// Run auto migrations on startup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := turso.AutoMigrate(ctx, dbClient); err != nil {
		logger.Error("failed to execute database schema migrations", "error", err)
		os.Exit(1)
	}

	// Initialize API server
	srv := api.NewServer(dbClient)

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	addr := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		logger.Info("api server listening", "address", addr)
		if err := srv.Router().Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	logger.Info("shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Router().Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	} else {
		logger.Info("fincher service stopped")
	}
}
