package turso

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/logger"
)

// Open creates a new ent.Client connected to Turso or local SQLite.
func Open(dbURL, authToken string) (*ent.Client, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	var (
		driverName string
		dsn        string
	)

	if strings.HasPrefix(dbURL, "libsql://") || strings.HasPrefix(dbURL, "https://") || strings.HasPrefix(dbURL, "http://") {
		driverName = "libsql"
		if authToken != "" {
			dsn = fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
		} else {
			dsn = dbURL
		}
	} else {
		driverName = "sqlite3"
		if dbURL == ":memory:" {
			dsn = "file::memory:?cache=shared&_fk=1"
		} else {
			dsn = fmt.Sprintf("file:%s?_fk=1&_journal=WAL", dbURL)
		}
	}

	logger.Debug("opening database connection", "driver", driverName, "url", dbURL)
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	logger.Info("connected to database", "driver", driverName, "url", dbURL)

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))

	return client, nil
}

// AutoMigrate executes schema migrations.
func AutoMigrate(ctx context.Context, client *ent.Client) error {
	start := time.Now()
	logger.Debug("starting database schema automigrations")
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("running automigration: %w", err)
	}
	logger.Info("schema migration completed", "duration_ms", time.Since(start).Milliseconds())
	return nil
}
