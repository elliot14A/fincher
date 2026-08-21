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
		t.Fatalf("failed to automigrate: %v", err)
	}

	// Seed Title
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

	// 1. Create Master V12 superseding V11
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
	if res1.Unwrap().Version != "V12" {
		t.Fatalf("expected version V12, got %s", res1.Unwrap().Version)
	}

	// Verify Title.CurrentMasterVersion updated to V12
	t1 := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if t1.CurrentMasterVersion != "V12" {
		t.Fatalf("expected title master version V12, got %s", t1.CurrentMasterVersion)
	}

	// 2. Create Master V13 superseding V12
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

	// Verify Title.CurrentMasterVersion updated sequentially to V13
	t2 := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if t2.CurrentMasterVersion != "V13" {
		t.Fatalf("expected title master version V13, got %s", t2.CurrentMasterVersion)
	}

	// 3. Get Master
	getRes := masters.Get(ctx, client, "master-eclipse-v13")
	if getRes.IsErr() {
		t.Fatalf("failed to get master: %v", getRes.Error())
	}

	// 4. List Masters (both V12 and V13)
	listRes := masters.List(ctx, client, domainerrors.Some("title-eclipse"))
	if listRes.IsErr() {
		t.Fatalf("failed to list masters: %v", listRes.Error())
	}
	if len(listRes.Unwrap()) != 2 {
		t.Fatalf("expected 2 masters, got %d", len(listRes.Unwrap()))
	}

	// 5. Delete Master
	delRes := masters.Delete(ctx, client, "master-eclipse-v13")
	if delRes.IsErr() {
		t.Fatalf("failed to delete master: %v", delRes.Error())
	}

	// 6. Verify Delete (Get returns NotFound)
	notFoundRes := masters.Get(ctx, client, "master-eclipse-v13")
	if notFoundRes.IsOk() {
		t.Fatalf("expected deleted master to return error, got: %+v", notFoundRes.Unwrap())
	}
	domErr, ok := notFoundRes.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeNotFound {
		t.Fatalf("expected CodeNotFound error after deletion, got: %v", notFoundRes.Error())
	}
}

func TestMasters_FK_Constraint(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// Attempt to create master for non-existent title
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
		t.Fatalf("expected CodeInvalidInput for orphan FK, got: %v", res.Error())
	}
}
