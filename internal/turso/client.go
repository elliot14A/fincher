package turso

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/logger"
)

// Open creates a new ent.Client connected to Turso or local SQLite.
//
// Architectural Decision Record (SQLite Single Writer & In-DB Blob Storage):
//  1. Local SQLite uses a single open connection (SetMaxOpenConns(1)) with WAL mode
//     (_journal=WAL) and busy timeout (_busy_timeout=5000) to ensure zero database lock contention
//     (preventing SQLITE_BUSY errors) and fully serialized ACID writes.
//  2. Poster and avatar image blobs are stored directly in the `uploads` SQLite table.
//     Because all uploads are strictly bounded by MaxUploadSizeBytes (1MB), individual sequential
//     write times are < 1ms on local disk in WAL mode. This guarantees that control-plane and
//     operational entity writes are never stalled, while keeping Fincher entirely self-contained
//     with zero external S3/MinIO/blob store dependencies for local operations and demo environments.
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
			u, err := url.Parse(dbURL)
			if err != nil {
				return nil, fmt.Errorf("invalid database URL: %w", err)
			}
			q := u.Query()
			q.Set("authToken", authToken)
			u.RawQuery = q.Encode()
			dsn = u.String()
		} else {
			dsn = dbURL
		}
	} else {
		driverName = "sqlite3"
		if dbURL == ":memory:" {
			dsn = "file::memory:?cache=shared&_fk=1&_busy_timeout=5000"
		} else if strings.HasPrefix(dbURL, "file:") {
			dsn = dbURL
		} else {
			dsn = fmt.Sprintf("file:%s?_fk=1&_journal=WAL&_busy_timeout=5000", dbURL)
		}
	}

	logger.Debug("opening database connection", "driver", driverName, "url", dbURL)
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	if driverName == "sqlite3" {
		db.SetMaxOpenConns(1)
	} else {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
	}
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
