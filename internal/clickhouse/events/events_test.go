package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/clickhouse"
	"github.com/elliot14A/fincher/internal/clickhouse/events"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestInsert_Validation(t *testing.T) {
	ctx := context.Background()

	// Missing Type should fail validation
	invalidEvent := &models.Event{
		Source:   "qc.agent",
		Severity: models.SeverityInfo,
	}

	res := events.Insert(ctx, nil, invalidEvent)
	if res.IsOk() {
		t.Fatal("expected insert with missing type to fail")
	}

	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok {
		t.Fatalf("expected *DomainError, got: %T", res.Error())
	}
	if domErr.Code != domainerrors.CodeInvalidInput {
		t.Fatalf("expected CodeInvalidInput, got: %s", domErr.Code)
	}
}

func TestInsert_Integration(t *testing.T) {
	conn, err := clickhouse.Open("127.0.0.1:9000")
	if err != nil {
		t.Skip("skipping integration test: clickhouse not reachable:", err)
		return
	}
	defer conn.Close()

	ctx := context.Background()
	if err := clickhouse.AutoMigrate(ctx, conn); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	// 1. Title-scoped event with subject as title slug
	titleScoped := &models.Event{
		Type:     models.TypeQCInspectionCompleted,
		Source:   "qc.agent.audio",
		Subject:  "eclipse", // title slug
		Severity: models.SeverityWarn,
		Time:     time.Now().UTC(),
		Data: map[string]any{
			"package_id":    "pkg-test-audio",
			"vendor_id":     "vendor-test",
			"component":     "AUDIO",
			"status":        "FAILED",
			"sync_drift_ms": 140.0,
		},
	}

	res := events.Insert(ctx, conn, titleScoped)
	if !res.IsOk() {
		t.Fatalf("expected title-scoped insert to succeed, got: %v", res.Error())
	}

	// 2. Title-agnostic event (Subject is empty -> defaults to GLOBAL sentinel)
	titleAgnostic := &models.Event{
		Type:     models.TypeVendorHeartbeat,
		Source:   "gateway.munich",
		Severity: models.SeverityInfo,
		Time:     time.Now().UTC(),
		Data: map[string]any{
			"status": "ONLINE",
		},
	}

	res1 := events.Insert(ctx, conn, titleAgnostic)
	if !res1.IsOk() {
		t.Fatalf("expected title-agnostic insert to succeed, got: %v", res1.Error())
	}

	var recordedSubject string
	query := "select subject from fincher.events where id = ? limit 1"
	if err := conn.QueryRowContext(ctx, query, titleAgnostic.ID).Scan(&recordedSubject); err != nil {
		t.Fatalf("querying inserted title-agnostic event: %v", err)
	}

	if recordedSubject != models.DefaultTitleAgnosticSentinel {
		t.Fatalf("expected subject to be %q, got %q", models.DefaultTitleAgnosticSentinel, recordedSubject)
	}
}
