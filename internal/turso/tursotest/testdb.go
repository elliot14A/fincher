package tursotest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
)

// NewMemoryClient initializes an isolated in-memory Turso SQLite client with all migrations applied.
func NewMemoryClient(t *testing.T) *ent.Client {
	t.Helper()

	testName := strings.ReplaceAll(t.Name(), "/", "_")
	dbURL := fmt.Sprintf("file:mem_%s_%d?mode=memory&cache=shared&_fk=1&_busy_timeout=5000", testName, time.Now().UnixNano())

	client, err := turso.Open(dbURL, "")
	if err != nil {
		t.Fatalf("tursotest: failed to open in-memory database: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("tursotest: failed to run automigrations: %v", err)
	}

	return client
}
