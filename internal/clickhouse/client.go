package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/elliot14A/fincher/pkg/logger"
)

// bootstrapDatabase connects briefly to the default database to ensure fincher exists.
func bootstrapDatabase(dsn string) error {
	rootDSN := dsn
	if !strings.HasPrefix(rootDSN, "clickhouse://") && !strings.HasPrefix(rootDSN, "tcp://") {
		rootDSN = fmt.Sprintf("clickhouse://%s/default", dsn)
	} else if strings.Contains(rootDSN, "/fincher") {
		rootDSN = strings.Replace(rootDSN, "/fincher", "/default", 1)
	}

	db, err := sql.Open("clickhouse", rootDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, "create database if not exists fincher;")
	return err
}

// Open establishes a connection pool to ClickHouse and guarantees the target database exists.
func Open(dsn string) (*sql.DB, error) {
	if err := bootstrapDatabase(dsn); err != nil {
		logger.Warn("clickhouse bootstrap warning (attempting direct connect)", "error", err)
	}

	connStr := dsn
	if !strings.HasPrefix(connStr, "clickhouse://") && !strings.HasPrefix(connStr, "tcp://") {
		connStr = fmt.Sprintf("clickhouse://%s/fincher?async_insert=1&wait_for_async_insert=1", dsn)
	} else if !strings.Contains(connStr, "async_insert") {
		if strings.Contains(connStr, "?") {
			connStr += "&async_insert=1&wait_for_async_insert=1"
		} else {
			connStr += "?async_insert=1&wait_for_async_insert=1"
		}
	}

	db, err := sql.Open("clickhouse", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening clickhouse database/sql pool: %w", err)
	}

	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(4)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("pinging clickhouse via database/sql: %w", err)
	}

	logger.Info("connected to clickhouse via database/sql", "dsn", dsn)
	return db, nil
}
