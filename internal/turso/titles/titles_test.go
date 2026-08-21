package titles_test

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

func setupTestEnt(t *testing.T) *ent.Client {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open memory ent client: %v", err)
	}

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to run ent automigration: %v", err)
	}

	return client
}

func TestTitles_EntCRUD(t *testing.T) {
	client := setupTestEnt(t)
	defer client.Close()

	ctx := context.Background()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	premiere := now.Add(48 * time.Hour)

	title := &models.Title{
		ID:                   "title-eclipse",
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         premiere,
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	}

	// 1. Create
	createRes := titles.Create(ctx, client, title)
	if createRes.IsErr() {
		t.Fatalf("failed to create title: %v", createRes.Error())
	}
	created := createRes.Unwrap()
	if created.ID != "title-eclipse" || created.Name != "Eclipse" {
		t.Errorf("unexpected created title: %+v", created)
	}

	// 2. Get
	getRes := titles.Get(ctx, client, "title-eclipse")
	if getRes.IsErr() {
		t.Fatalf("failed to get title: %v", getRes.Error())
	}
	fetched := getRes.Unwrap()
	if fetched.Name != "Eclipse" || fetched.Territories != 40 {
		t.Errorf("unexpected fetched title: %+v", fetched)
	}

	// 3. List (no filter)
	listRes := titles.List(ctx, client, domainerrors.None[models.TitleStatus]())
	if listRes.IsErr() {
		t.Fatalf("failed to list titles: %v", listRes.Error())
	}
	allTitles := listRes.Unwrap()
	if len(allTitles) != 1 {
		t.Errorf("expected 1 title, got %d", len(allTitles))
	}

	// 4. List (with filter)
	filterRes := titles.List(ctx, client, domainerrors.Some(models.StatusAtRisk))
	if filterRes.IsErr() {
		t.Fatalf("failed to list filtered titles: %v", filterRes.Error())
	}
	if len(filterRes.Unwrap()) != 1 {
		t.Errorf("expected 1 title matching AT_RISK, got %d", len(filterRes.Unwrap()))
	}

	noneRes := titles.List(ctx, client, domainerrors.Some(models.StatusOnTrack))
	if noneRes.IsErr() {
		t.Fatalf("failed to list filtered titles: %v", noneRes.Error())
	}
	if len(noneRes.Unwrap()) != 0 {
		t.Errorf("expected 0 titles matching ON_TRACK, got %d", len(noneRes.Unwrap()))
	}

	// 5. Partial Update (PATCH)
	newStatus := models.StatusHold
	newTerritories := 55
	updateRes := titles.Update(ctx, client, "title-eclipse", &models.UpdateTitleInput{
		OverallStatus: &newStatus,
		Territories:   &newTerritories,
	})
	if updateRes.IsErr() {
		t.Fatalf("failed to partially update title: %v", updateRes.Error())
	}

	updatedGet := titles.Get(ctx, client, "title-eclipse").Unwrap()
	if updatedGet.Territories != 55 || updatedGet.OverallStatus != models.StatusHold || updatedGet.Name != "Eclipse" {
		t.Errorf("partial update not persisted correctly: %+v", updatedGet)
	}

	// 6. Delete
	delRes := titles.Delete(ctx, client, "title-eclipse")
	if delRes.IsErr() {
		t.Fatalf("failed to delete title: %v", delRes.Error())
	}

	// 7. Get after Delete should return domainerrors.CodeNotFound
	afterDelRes := titles.Get(ctx, client, "title-eclipse")
	if afterDelRes.IsOk() {
		t.Errorf("expected error after delete, got title: %+v", afterDelRes.Unwrap())
	}
	var domErr *domainerrors.DomainError
	if domErr = afterDelRes.Error().(*domainerrors.DomainError); domErr.Code != domainerrors.CodeNotFound {
		t.Errorf("expected CodeNotFound, got %s", domErr.Code)
	}
}

func TestTitles_FK_DeleteBlockedByDependents(t *testing.T) {
	client := setupTestEnt(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create Title
	_ = titles.Create(ctx, client, &models.Title{
		ID:                   "title-active",
		Name:                 "Active Title",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(24 * time.Hour),
		Territories:          10,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusOnTrack,
	})

	// 2. Create Master referencing title-active
	_ = masters.Create(ctx, client, &models.Master{
		ID:      "master-active-v01",
		TitleID: "title-active",
		Version: "V01",
	})

	// 3. Deleting title-active must fail with CodeConflict (409)
	delRes := titles.Delete(ctx, client, "title-active")
	if delRes.IsOk() {
		t.Fatalf("expected title delete to be blocked by dependent master")
	}

	domErr, ok := delRes.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeConflict {
		t.Errorf("expected CodeConflict (409) for delete blocked by dependents, got: %v", delRes.Error())
	}
}
