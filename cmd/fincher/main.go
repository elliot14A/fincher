package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"

	"github.com/elliot14A/fincher/internal/api"
	"github.com/elliot14A/fincher/pkg/domain/config"
	"github.com/elliot14A/fincher/pkg/turso"
)

var CLI struct {
	config.Config
}

func main() {
	kong.Parse(&CLI,
		kong.Name("fincher"),
		kong.Description("Fincher: Autonomous Delivery-Integrity Engine for LUME"),
		kong.UsageOnError(),
	)

	cfg := &CLI.Config
	if err := cfg.Validate(); err != nil {
		slog.Error("configuration validation failed", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting fincher server",
		"port", cfg.Port,
		"env", cfg.Environment,
		"turso_url", cfg.TursoURL,
		"mcp_url", cfg.MCPURL,
	)

	// Initialize database connection and run schema migrations
	client, err := turso.Open(cfg.TursoURL, cfg.TursoToken)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := turso.AutoMigrate(ctx, client); err != nil {
		slog.Error("failed to run database automigration", "error", err)
		os.Exit(1)
	}
	slog.Info("database automigration completed successfully")

	// Initialize HTTP server
	server := api.NewServer(client)
	addr := fmt.Sprintf(":%d", cfg.Port)

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Router().Start(addr); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("fincher ready to serve traffic", "addr", addr)

	<-stop
	slog.Info("shutting down fincher gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Router().Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
		os.Exit(1)
	}

	slog.Info("server shutdown complete")
}
