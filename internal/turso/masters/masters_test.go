package masters_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/masters"
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
		t.Fatalf("failed to run schema automigrations: %v", err)
	}

	titleRes := titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-eclipse"},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V11",
		OverallStatus:        models.StatusAtRisk,
	})
	if titleRes.IsErr() {
		t.Fatalf("failed to seed title: %v", titleRes.Error())
	}

	return client
}

func TestMasters_CRUD(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create Master V12
	m1 := &models.Master{
		ID:                "master-eclipse-v12",
		TitleID:           "title-eclipse",
		Version:           "V12",
		SupersedesVersion: "V11",
		CreatedAt:         time.Now().UTC(),
	}
	res1 := masters.Create(ctx, client, m1)
	if res1.IsErr() {
		t.Fatalf("failed to create master V12: %v", res1.Error())
	}

	// 2. Create Master V13
	m2 := &models.Master{
		ID:                "master-eclipse-v13",
		TitleID:           "title-eclipse",
		Version:           "V13",
		SupersedesVersion: "V12",
		CreatedAt:         time.Now().UTC().Add(time.Hour),
	}
	res2 := masters.Create(ctx, client, m2)
	if res2.IsErr() {
		t.Fatalf("failed to create master V13: %v", res2.Error())
	}

	// 3. Get Master
	getRes := masters.Get(ctx, client, "master-eclipse-v13")
	if getRes.IsErr() {
		t.Fatalf("failed to get master: %v", getRes.Error())
	}

	// 4. List Masters with Pagination
	p := models.NewPagination(1, 10, "asc", "")
	listRes := masters.List(ctx, client, domainerrors.Some("title-eclipse"), p)
	if listRes.IsErr() {
		t.Fatalf("failed to list masters: %v", listRes.Error())
	}
	res := listRes.Unwrap()
	if len(res.Items) != 2 || res.TotalItems != 2 {
		t.Fatalf("expected 2 masters, got %d (total: %d)", len(res.Items), res.TotalItems)
	}

	// 5. Delete Master V13 -> Should reconcile Title.CurrentMasterVersion to V12
	delRes := masters.Delete(ctx, client, "master-eclipse-v13")
	if delRes.IsErr() {
		t.Fatalf("failed to delete master V13: %v", delRes.Error())
	}

	titleAfterV13Delete := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if titleAfterV13Delete.CurrentMasterVersion != "V12" {
		t.Fatalf("expected title master version to revert to V12, got '%s'", titleAfterV13Delete.CurrentMasterVersion)
	}

	// 6. Delete Master V12 (last remaining master) -> Should clear Title.CurrentMasterVersion to ""
	delV12Res := masters.Delete(ctx, client, "master-eclipse-v12")
	if delV12Res.IsErr() {
		t.Fatalf("failed to delete master V12: %v", delV12Res.Error())
	}

	titleAfterAllDelete := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if titleAfterAllDelete.CurrentMasterVersion != "" {
		t.Fatalf("expected title master version to be cleared to empty string, got '%s'", titleAfterAllDelete.CurrentMasterVersion)
	}

	// 7. Verify Delete
	notFoundRes := masters.Get(ctx, client, "master-eclipse-v13")
	if notFoundRes.IsOk() {
		t.Fatalf("expected deleted master to return error, got: %+v", notFoundRes.Unwrap())
	}
	domErr, ok := notFoundRes.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeNotFound {
		t.Fatalf("expected CodeNotFound error, got: %v", notFoundRes.Error())
	}
}

func TestMasters_FK_Constraint(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	orphanMaster := &models.Master{
		ID:        "master-orphan",
		TitleID:   "non-existent-title",
		Version:   "V1",
		CreatedAt: time.Now(),
	}

	res := masters.Create(ctx, client, orphanMaster)
	if res.IsOk() {
		t.Fatalf("expected foreign key constraint violation, but creation succeeded")
	}
	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeInvalidInput {
		t.Fatalf("expected CodeInvalidInput for FK constraint violation, got: %v", res.Error())
	}
}
