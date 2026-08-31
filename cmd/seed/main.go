package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"

	"github.com/elliot14A/fincher/internal/seed"
	"github.com/elliot14A/fincher/internal/seed/types"
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/pkg/logger"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var cfg types.SeedConfig
	kong.Parse(&cfg,
		kong.Name("fincher-seed"),
		kong.Description("Fincher baseline data-generation and historical seed tool"),
		kong.UsageOnError(),
	)

	ctx := context.Background()

	// Connect Turso SQLite
	tursoClient, err := turso.Open(cfg.TursoURL, cfg.TursoToken)
	if err != nil {
		logger.Error("failed to open Turso / SQLite client", "error", err)
		os.Exit(1)
	}
	defer tursoClient.Close()

	if err := turso.AutoMigrate(ctx, tursoClient); err != nil {
		logger.Error("failed to run Turso automigrations", "error", err)
		os.Exit(1)
	}

	// Connect ClickHouse
	chDB, err := sql.Open("clickhouse", fmt.Sprintf("clickhouse://%s?database=fincher&debug=false", cfg.ClickHouseDSN))
	if err != nil {
		logger.Error("failed to open ClickHouse connection", "error", err)
		os.Exit(1)
	}
	defer chDB.Close()

	if err := chDB.PingContext(ctx); err != nil {
		logger.Error("failed to ping ClickHouse", "error", err)
		os.Exit(1)
	}

	// Execute seed pipeline
	seeder := seed.NewSeeder(&cfg, tursoClient, chDB)
	summary, err := seeder.Run(ctx)
	if err != nil {
		logger.Error("seed pipeline failed", "error", err)
		os.Exit(1)
	}

	fmt.Printf("\n=== Fincher Seed Completed Successfully ===\n")
	fmt.Printf("  Vendors:      %d\n", summary.Vendors)
	fmt.Printf("  Titles:       %d\n", summary.Titles)
	fmt.Printf("  Events:       %d\n", summary.Events)
	fmt.Printf("  Duration:     %s\n", summary.Duration.Round(time.Millisecond))
	fmt.Printf("===========================================\n\n")
}
