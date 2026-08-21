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

	// Seed parent title
	titleRes := titles.Create(ctx, client, &models.Title{
		ID:                   "title-eclipse",
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

	m1 := &models.Master{
		ID:                "master-eclipse-v12",
		TitleID:           "title-eclipse",
		Version:           "V12",
		SupersedesVersion: "V11",
	}

	m2 := &models.Master{
		ID:                "master-eclipse-v13",
		TitleID:           "title-eclipse",
		Version:           "V13",
		SupersedesVersion: "V12",
	}

	// 1. Create m1 and verify title current_master_version is updated
	res1 := masters.Create(ctx, client, m1)
	if res1.IsErr() {
		t.Fatalf("failed to create master 1: %v", res1.Error())
	}
	title1 := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if title1.CurrentMasterVersion != "V12" {
		t.Errorf("expected title CurrentMasterVersion to be V12, got %s", title1.CurrentMasterVersion)
	}

	// 2. Create m2 and verify title is updated to V13
	res2 := masters.Create(ctx, client, m2)
	if res2.IsErr() {
		t.Fatalf("failed to create master 2: %v", res2.Error())
	}
	title2 := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if title2.CurrentMasterVersion != "V13" {
		t.Errorf("expected title CurrentMasterVersion to be V13, got %s", title2.CurrentMasterVersion)
	}

	// 3. Get
	getRes := masters.Get(ctx, client, "master-eclipse-v13")
	if getRes.IsErr() {
		t.Fatalf("failed to get master: %v", getRes.Error())
	}
	if getRes.Unwrap().Version != "V13" || getRes.Unwrap().SupersedesVersion != "V12" {
		t.Errorf("unexpected master fields: %+v", getRes.Unwrap())
	}

	// 4. List
	listRes := masters.List(ctx, client, domainerrors.Some("title-eclipse"))
	if listRes.IsErr() {
		t.Fatalf("failed to list masters: %v", listRes.Error())
	}
	if len(listRes.Unwrap()) != 2 {
		t.Errorf("expected 2 masters, got %d", len(listRes.Unwrap()))
	}

	// 5. Delete
	delRes := masters.Delete(ctx, client, "master-eclipse-v12")
	if delRes.IsErr() {
		t.Fatalf("failed to delete master: %v", delRes.Error())
	}

	// 6. Get after Delete
	afterDel := masters.Get(ctx, client, "master-eclipse-v12")
	if afterDel.IsOk() {
		t.Errorf("expected error after delete, got master")
	}
}

func TestMasters_FK_Constraint(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// Negative Test: Creating Master for non-existent Title should fail with CodeInvalidInput (400)
	orphanMaster := &models.Master{
		ID:      "master-orphan",
		TitleID: "non-existent-title",
		Version: "V01",
	}

	res := masters.Create(ctx, client, orphanMaster)
	if res.IsOk() {
		t.Fatalf("expected orphan master creation to fail, but succeeded")
	}

	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeInvalidInput {
		t.Errorf("expected CodeInvalidInput for orphan FK, got: %v", res.Error())
	}
}
