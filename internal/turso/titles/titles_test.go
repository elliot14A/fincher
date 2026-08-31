package titles_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/masters"
	"github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/uploads"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func setupTestDB(t *testing.T) *ent.Client {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to run schema automigrations: %v", err)
	}

	return client
}

func TestTitles_EntCRUD(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create Title with Metadata
	now := time.Now().UTC().Truncate(time.Second)
	premiere := now.Add(48 * time.Hour)
	title := &models.Title{
		Base: models.Base{
			ID: "title-eclipse",
			Metadata: map[string]any{
				"poster_url": "https://cdn.lume.stream/posters/eclipse.jpg",
				"genre":      "Sci-Fi",
			},
		},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         premiere,
		Territories:          40,
		CurrentMasterVersion: "V12",
		OverallStatus:        models.StatusAtRisk,
	}

	createRes := titles.Create(ctx, client, title)
	if createRes.IsErr() {
		t.Fatalf("expected title creation to succeed, got error: %v", createRes.Error())
	}
	created := createRes.Unwrap()
	if created.ID != "title-eclipse" || created.Name != "Eclipse" {
		t.Fatalf("unexpected title data: %+v", created)
	}

	// 2. Get Title by ID, Slug, and FindByIDOrSlug
	getRes := titles.Get(ctx, client, "title-eclipse")
	if getRes.IsErr() {
		t.Fatalf("expected title fetch to succeed, got error: %v", getRes.Error())
	}
	fetched := getRes.Unwrap()
	if fetched.CurrentMasterVersion != "V12" {
		t.Errorf("expected CurrentMasterVersion 'V12', got '%s'", fetched.CurrentMasterVersion)
	}

	slugRes := titles.GetBySlug(ctx, client, "eclipse")
	if slugRes.IsErr() {
		t.Fatalf("expected title fetch by slug 'eclipse' to succeed, got error: %v", slugRes.Error())
	}
	if slugRes.Unwrap().ID != "title-eclipse" {
		t.Errorf("expected ID 'title-eclipse', got '%s'", slugRes.Unwrap().ID)
	}

	findRes := titles.FindByIDOrSlug(ctx, client, "Eclipse")
	if findRes.IsErr() {
		t.Fatalf("expected title fetch by name 'Eclipse' to succeed, got error: %v", findRes.Error())
	}
	if findRes.Unwrap().Slug != "eclipse" {
		t.Errorf("expected slug 'eclipse', got '%s'", findRes.Unwrap().Slug)
	}

	// 3. List Titles (Filter by status + Pagination)
	p := models.NewPagination(1, 10, "asc", "")
	listRes := titles.List(ctx, client, domainerrors.Some(models.StatusAtRisk), p)
	if listRes.IsErr() {
		t.Fatalf("expected title listing to succeed, got error: %v", listRes.Error())
	}
	res := listRes.Unwrap()
	if len(res.Items) != 1 || res.TotalItems != 1 {
		t.Fatalf("expected 1 title with status AT_RISK, got %d (total: %d)", len(res.Items), res.TotalItems)
	}

	// 4. Update Title
	newStatus := models.StatusHold
	newMaster := "V13"
	updateRes := titles.Update(ctx, client, "title-eclipse", &models.UpdateTitleInput{
		OverallStatus:        &newStatus,
		CurrentMasterVersion: &newMaster,
		Metadata: map[string]any{
			"poster_url": "https://cdn.lume.stream/posters/eclipse_v2.jpg",
			"genre":      "Sci-Fi",
			"priority":   "P0",
		},
	})
	if updateRes.IsErr() {
		t.Fatalf("expected title update to succeed, got error: %v", updateRes.Error())
	}
	updated := updateRes.Unwrap()
	if updated.OverallStatus != models.StatusHold || updated.CurrentMasterVersion != "V13" {
		t.Errorf("unexpected updated status or master version: %+v", updated)
	}

	// 5. Delete Title
	deleteRes := titles.Delete(ctx, client, "title-eclipse")
	if deleteRes.IsErr() {
		t.Fatalf("expected title deletion to succeed, got error: %v", deleteRes.Error())
	}

	// 6. Verify Delete
	notFoundRes := titles.Get(ctx, client, "title-eclipse")
	if notFoundRes.IsOk() {
		t.Fatalf("expected deleted title to return error, got title: %+v", notFoundRes.Unwrap())
	}
}

func TestTitles_FK_DeleteBlockedByDependents(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	_ = titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-active"},
		Name:                 "Active Title",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(24 * time.Hour),
		Territories:          10,
		CurrentMasterVersion: "V1",
		OverallStatus:        models.StatusOnTrack,
	})

	_ = masters.Create(ctx, client, &models.Master{
		ID:        "master-active-v1",
		TitleID:   "title-active",
		Version:   "V1",
		CreatedAt: time.Now(),
	})

	delRes := titles.Delete(ctx, client, "title-active")
	if delRes.IsOk() {
		t.Fatalf("expected title delete to fail due to foreign key dependents")
	}
}

func TestTitles_CascadeUploadDelete(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create upload blob
	upRes := uploads.Create(ctx, client, &models.Upload{
		ID:        "upload-poster-123",
		Filename:  "poster.png",
		MimeType:  "image/png",
		SizeBytes: 4,
		Data:      []byte{0x89, 0x50, 0x4E, 0x47},
	})
	if upRes.IsErr() {
		t.Fatalf("failed to create upload: %v", upRes.Error())
	}

	// 2. Create title referencing this upload
	titleRes := titles.Create(ctx, client, &models.Title{
		Base: models.Base{
			ID: "title-with-avatar",
			Metadata: map[string]any{
				"avatar_url": "/api/uploads/upload-poster-123",
			},
		},
		Name:                 "Avatar Movie",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusOnTrack,
	})
	if titleRes.IsErr() {
		t.Fatalf("failed to create title: %v", titleRes.Error())
	}

	// 3. Delete title
	delRes := titles.Delete(ctx, client, "title-with-avatar")
	if delRes.IsErr() {
		t.Fatalf("failed to delete title: %v", delRes.Error())
	}

	// 4. Verify that upload was cascaded and deleted
	getUpRes := uploads.Get(ctx, client, "upload-poster-123")
	if getUpRes.IsOk() {
		t.Fatalf("expected upload blob to be deleted along with title, but it still exists")
	}
}

func TestTitles_ResolveHoursUntilPremiere_And_SlugDeduplication(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. GLOBAL sentinel returns standard 48h fallback without DB query
	globalHours := titles.ResolveHoursUntilPremiere(ctx, client, "GLOBAL", 0)
	if globalHours != 48.0 {
		t.Errorf("expected 48.0 for GLOBAL, got: %f", globalHours)
	}

	// 2. Custom explicit hours override is respected
	overrideHours := titles.ResolveHoursUntilPremiere(ctx, client, "GLOBAL", 12.5)
	if overrideHours != 12.5 {
		t.Errorf("expected 12.5 override, got: %f", overrideHours)
	}

	// 3. Create Title with Slug "avatar-fire-ash"
	premiere := time.Now().UTC().Add(24 * time.Hour)
	t1Res := titles.Create(ctx, client, &models.Title{
		Base: models.Base{
			ID: "title-avatar-1",
		},
		Name:                 "Avatar: Fire & Ash",
		Slug:                 "avatar-fire-ash",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         premiere,
		Territories:          10,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusProcessing,
	})
	if t1Res.IsErr() {
		t.Fatalf("failed to create first title: %v", t1Res.Error())
	}
	t1 := t1Res.Unwrap()
	if t1.Slug != "avatar-fire-ash" {
		t.Errorf("expected slug avatar-fire-ash, got: %s", t1.Slug)
	}

	// 4. Create second Title with identical Slug: should automatically deduplicate without 500 collision
	t2Res := titles.Create(ctx, client, &models.Title{
		Base: models.Base{
			ID: "title-avatar-2",
		},
		Name:                 "Avatar Fire and Ash Special Edition",
		Slug:                 "avatar-fire-ash",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         premiere,
		Territories:          10,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusProcessing,
	})
	if t2Res.IsErr() {
		t.Fatalf("expected duplicate slug creation to succeed via auto-dedup, got: %v", t2Res.Error())
	}
	t2 := t2Res.Unwrap()
	if t2.Slug == "avatar-fire-ash" {
		t.Errorf("expected slug to be deduplicated, but got identical: %s", t2.Slug)
	}
	if !strings.HasPrefix(t2.Slug, "avatar-fire-ash-") {
		t.Errorf("expected deduplicated slug prefix avatar-fire-ash-, got: %s", t2.Slug)
	}

	// 5. Resolve hours dynamically for avatar-fire-ash
	resolvedHours := titles.ResolveHoursUntilPremiere(ctx, client, "avatar-fire-ash", 0)
	if resolvedHours < 23.0 || resolvedHours > 25.0 {
		t.Errorf("expected ~24h for avatar-fire-ash, got: %f", resolvedHours)
	}
}
