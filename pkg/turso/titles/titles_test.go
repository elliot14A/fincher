package titles_test

import (
	"context"
	"testing"
	"time"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	"github.com/elliot14A/fincher/pkg/turso"
	"github.com/elliot14A/fincher/pkg/turso/titles"
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
	if updatedGet.Territories != 55 || updatedGet.OverallStatus != "HOLD" || updatedGet.Name != "Eclipse" {
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
