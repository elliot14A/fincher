package deliveries_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/titles"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func setupTestDB(t *testing.T) *ent.Client {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}

	titleRes := titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-eclipse"},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})
	if titleRes.IsErr() {
		t.Fatalf("failed to seed title: %v", titleRes.Error())
	}

	return client
}

func TestDeliveries_CRUD(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	d1 := &models.Delivery{
		Base: models.Base{
			ID: "del-eclipse-us",
			Metadata: map[string]any{
				"carrier": "DirectLine",
				"qc_gate": "passed",
			},
		},
		TitleID:    "title-eclipse",
		Country:    "US",
		Status:     models.DeliveryStatusPending,
		TargetDate: time.Now().Add(24 * time.Hour),
	}

	// 1. Create
	createRes := deliveries.Create(ctx, client, d1)
	if createRes.IsErr() {
		t.Fatalf("failed to create delivery: %v", createRes.Error())
	}
	if createRes.Unwrap().Country != "US" {
		t.Errorf("unexpected delivery country: %s", createRes.Unwrap().Country)
	}
	if createRes.Unwrap().Metadata["carrier"] != "DirectLine" {
		t.Errorf("expected metadata carrier, got: %v", createRes.Unwrap().Metadata["carrier"])
	}

	// 2. Get
	getRes := deliveries.Get(ctx, client, "del-eclipse-us")
	if getRes.IsErr() {
		t.Fatalf("failed to get delivery: %v", getRes.Error())
	}

	// 3. List
	p := models.NewPagination(1, 10, "asc", "")
	listRes := deliveries.List(ctx, client, deliveries.ListFilter{
		TitleID: domainerrors.Some("title-eclipse"),
		Country: domainerrors.None[string](),
		Status:  domainerrors.None[models.DeliveryStatus](),
	}, p)
	if listRes.IsErr() {
		t.Fatalf("failed to list deliveries: %v", listRes.Error())
	}
	res := listRes.Unwrap()
	if len(res.Items) != 1 || res.TotalItems != 1 {
		t.Errorf("expected 1 delivery, got %d (total: %d)", len(res.Items), res.TotalItems)
	}

	// 4. Update
	newStatus := models.DeliveryStatusReadyToShip
	upRes := deliveries.Update(ctx, client, "del-eclipse-us", &models.UpdateDeliveryInput{
		Status: &newStatus,
		Metadata: map[string]any{
			"carrier": "DirectLineExpress",
			"qc_gate": "passed",
		},
	})
	if upRes.IsErr() {
		t.Fatalf("failed to update delivery: %v", upRes.Error())
	}
	if upRes.Unwrap().Status != models.DeliveryStatusReadyToShip {
		t.Errorf("expected status READY_TO_SHIP, got %s", upRes.Unwrap().Status)
	}
	if upRes.Unwrap().Metadata["carrier"] != "DirectLineExpress" {
		t.Errorf("expected updated carrier, got: %v", upRes.Unwrap().Metadata["carrier"])
	}

	// 5. Delete
	delRes := deliveries.Delete(ctx, client, "del-eclipse-us")
	if delRes.IsErr() {
		t.Fatalf("failed to delete delivery: %v", delRes.Error())
	}
}

func TestDeliveries_FK_Constraint(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	orphanDel := &models.Delivery{
		Base:       models.Base{ID: "del-orphan"},
		TitleID:    "non-existent-title",
		Country:    "US",
		Status:     models.DeliveryStatusPending,
		TargetDate: time.Now().Add(24 * time.Hour),
	}

	res := deliveries.Create(ctx, client, orphanDel)
	if res.IsOk() {
		t.Fatalf("expected orphan delivery creation to fail")
	}

	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeInvalidInput {
		t.Errorf("expected CodeInvalidInput for orphan FK, got: %v", res.Error())
	}
}
