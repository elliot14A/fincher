package vendors_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/clickhouse"
	"github.com/elliot14A/fincher/internal/clickhouse/events"
	"github.com/elliot14A/fincher/internal/clickhouse/vendors"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestVendors_Validation(t *testing.T) {
	ctx := context.Background()

	res := vendors.RecencyWeightedAccuracy(ctx, nil, "", "AUDIO")
	if res.IsOk() {
		t.Fatal("expected empty vendorID to fail validation")
	}
	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok {
		t.Fatalf("expected *DomainError, got: %T", res.Error())
	}
	if domErr.Code != domainerrors.CodeInvalidInput {
		t.Fatalf("expected CodeInvalidInput, got: %s", domErr.Code)
	}

	metricsRes := vendors.GetMetrics(ctx, nil, "")
	if metricsRes.IsOk() {
		t.Fatal("expected empty vendorID in GetMetrics to fail validation")
	}
}

func TestVendors_Integration(t *testing.T) {
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

	vendorID := "vendor-acc-" + uuid.NewString()[:8]

	// Baseline accuracy for candidate with no historical data should be -1.0 (unmeasured)
	baselineRes := vendors.RecencyWeightedAccuracy(ctx, conn, vendorID, "AUDIO")
	if !baselineRes.IsOk() {
		t.Fatalf("calculating baseline accuracy: %v", baselineRes.Error())
	}
	if baselineRes.Unwrap() != -1.0 {
		t.Fatalf("expected baseline accuracy -1.0, got: %f", baselineRes.Unwrap())
	}

	// Insert 1 failed QC event using CloudEvents schema
	failedEvent := &models.Event{
		Type:     models.TypeQCInspectionCompleted,
		Source:   "qc.agent.audio",
		Subject:  "eclipse", // title slug
		Severity: models.SeverityWarn,
		Time:     time.Now().UTC(),
		Data: map[string]any{
			"package_id":    "pkg-acc-1",
			"vendor_id":     vendorID,
			"component":     "AUDIO",
			"status":        "FAILED",
			"sync_drift_ms": 130.0,
		},
	}

	if res := events.Insert(ctx, conn, failedEvent); !res.IsOk() {
		t.Fatalf("inserting test event: %v", res.Error())
	}

	// Accuracy should now be 0.0 (1 failed out of 1 measured)
	accRes := vendors.RecencyWeightedAccuracy(ctx, conn, vendorID, "AUDIO")
	if !accRes.IsOk() {
		t.Fatalf("calculating accuracy after failure: %v", accRes.Error())
	}
	if accRes.Unwrap() > 0.01 {
		t.Fatalf("expected accuracy ~0.0, got: %f", accRes.Unwrap())
	}

	// Verify GetMetrics fetches the rollup row
	rowsRes := vendors.GetMetrics(ctx, conn, vendorID)
	if !rowsRes.IsOk() {
		t.Fatalf("fetching metrics: %v", rowsRes.Error())
	}
	rows := rowsRes.Unwrap()
	if len(rows) == 0 {
		t.Fatal("expected at least 1 metrics row")
	}

	row := rows[0]
	if row.VendorID != vendorID || row.Component != "AUDIO" || row.FailedInspections != 1 {
		t.Fatalf("unexpected metric row: %+v", row)
	}

	// Cleanup test rows
	_, _ = conn.ExecContext(ctx, "alter table fincher.events delete where id = ?", failedEvent.ID)
}
