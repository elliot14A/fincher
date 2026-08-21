package turso

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/elliot14A/fincher/pkg/ent"
)

// Open creates a new database connection client.
func Open(dbURL, authToken string) (*ent.Client, error) {
	var connStr string
	var driverName string

	if strings.HasPrefix(dbURL, "libsql://") || strings.HasPrefix(dbURL, "https://") || strings.HasPrefix(dbURL, "http://") {
		driverName = "libsql"
		connStr = dbURL
		if authToken != "" {
			if strings.Contains(connStr, "?") {
				connStr += "&authToken=" + authToken
			} else {
				connStr += "?authToken=" + authToken
			}
		}
	} else {
		driverName = "sqlite3"
		if dbURL == ":memory:" {
			connStr = "file:ent?mode=memory&cache=shared&_fk=1"
		} else {
			connStr = fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_fk=1", dbURL)
		}
	}

	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	slog.Info("connected to database", "driver", driverName, "url", dbURL)

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))

	return client, nil
}

// AutoMigrate executes schema migrations.
func AutoMigrate(ctx context.Context, client *ent.Client) error {
	start := time.Now()
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("running automigration: %w", err)
	}
	slog.Info("schema migration completed", "duration_ms", time.Since(start).Milliseconds())
	return nil
}
