package dependencies_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/dependencies"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/vendors"
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

	// 1. Seed Title
	_ = titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-eclipse"},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          40,
		CurrentMasterVersion: "V13",
		OverallStatus:        models.StatusAtRisk,
	})

	// 2. Seed Vendor
	_ = vendors.Create(ctx, client, &models.Vendor{
		Base:       models.Base{ID: "vendor_a"},
		Name:       "Vendor A",
		Components: []string{"AUDIO"},
		Markets:    []string{"en-US"},
	})

	// 3. Seed Packages: Video, Spanish Audio, Spanish Subtitles
	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-video-ov"},
		TitleID:                  "title-eclipse",
		Component:                models.ComponentVideo,
		Language:                 "ov",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V13",
		Status:                   models.PackageStatusValid,
	})

	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-audio-es"},
		TitleID:                  "title-eclipse",
		Component:                models.ComponentAudio,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V13",
		Status:                   models.PackageStatusValid,
	})

	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-sub-es"},
		TitleID:                  "title-eclipse",
		Component:                models.ComponentSubtitle,
		Language:                 "es",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V13",
		Status:                   models.PackageStatusValid,
	})

	return client
}

func TestDependencies_CRUD_And_LineageGraph(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create edge: Video -> Audio (MASTER_DERIVATION)
	dep1 := &models.Dependency{
		ID:             "dep-video-to-audio",
		ParentID:       "pkg-video-ov",
		ChildID:        "pkg-audio-es",
		DependencyType: models.DependencyMasterDerivation,
	}
	res1 := dependencies.Create(ctx, client, dep1)
	if res1.IsErr() {
		t.Fatalf("failed to create dep 1: %v", res1.Error())
	}

	// 2. Create edge: Audio -> Subtitles (AUDIO_SYNC)
	dep2 := &models.Dependency{
		ID:             "dep-audio-to-sub",
		ParentID:       "pkg-audio-es",
		ChildID:        "pkg-sub-es",
		DependencyType: models.DependencyAudioSync,
	}
	res2 := dependencies.Create(ctx, client, dep2)
	if res2.IsErr() {
		t.Fatalf("failed to create dep 2: %v", res2.Error())
	}

	// 3. Get
	getRes := dependencies.Get(ctx, client, "dep-video-to-audio")
	if getRes.IsErr() {
		t.Fatalf("failed to get dep: %v", getRes.Error())
	}
	if getRes.Unwrap().DependencyType != models.DependencyMasterDerivation {
		t.Errorf("unexpected dep type: %s", getRes.Unwrap().DependencyType)
	}

	// 4. List with Pagination
	p := models.NewPagination(1, 10, "asc", "")
	listRes := dependencies.List(ctx, client, dependencies.ListFilter{
		ParentID: domainerrors.Some("pkg-video-ov"),
		ChildID:  domainerrors.None[string](),
	}, p)
	if listRes.IsErr() {
		t.Fatalf("failed to list deps: %v", listRes.Error())
	}
	res := listRes.Unwrap()
	if len(res.Items) != 1 || res.TotalItems != 1 {
		t.Errorf("expected 1 child of video, got %d (total: %d)", len(res.Items), res.TotalItems)
	}

	// 5. Lineage Graph Traversal
	graphRes := dependencies.GetLineageGraph(ctx, client, "title-eclipse")
	if graphRes.IsErr() {
		t.Fatalf("failed to get lineage graph: %v", graphRes.Error())
	}
	graph := graphRes.Unwrap()
	if len(graph.Roots) != 1 {
		t.Fatalf("expected 1 root node (video), got %d", len(graph.Roots))
	}
	root := graph.Roots[0]
	if root.PackageID != "pkg-video-ov" {
		t.Errorf("expected root package to be pkg-video-ov, got %s", root.PackageID)
	}
	if len(root.Children) != 1 || root.Children[0].PackageID != "pkg-audio-es" {
		t.Fatalf("expected 1 child (audio-es), got %+v", root.Children)
	}
	if len(root.Children[0].Children) != 1 || root.Children[0].Children[0].PackageID != "pkg-sub-es" {
		t.Fatalf("expected 1 grandchild (sub-es), got %+v", root.Children[0].Children)
	}

	// 6. Delete
	delRes := dependencies.Delete(ctx, client, "dep-video-to-audio")
	if delRes.IsErr() {
		t.Fatalf("failed to delete dep: %v", delRes.Error())
	}
}

func TestDependencies_CyclePrevention(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// Video -> Audio -> Subtitles
	_ = dependencies.Create(ctx, client, &models.Dependency{
		ID:             "dep-1",
		ParentID:       "pkg-video-ov",
		ChildID:        "pkg-audio-es",
		DependencyType: models.DependencyMasterDerivation,
	})
	_ = dependencies.Create(ctx, client, &models.Dependency{
		ID:             "dep-2",
		ParentID:       "pkg-audio-es",
		ChildID:        "pkg-sub-es",
		DependencyType: models.DependencyAudioSync,
	})

	// Attempt to create circular dependency: Subtitles -> Video (cycle!)
	cycleDep := &models.Dependency{
		ID:             "dep-cycle",
		ParentID:       "pkg-sub-es",
		ChildID:        "pkg-video-ov",
		DependencyType: models.DependencySubtitleAlignment,
	}

	res := dependencies.Create(ctx, client, cycleDep)
	if res.IsOk() {
		t.Fatalf("expected circular dependency creation to fail, but succeeded")
	}

	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeConflict {
		t.Errorf("expected CodeConflict for cycle, got: %v", res.Error())
	}
}

func TestDependencies_CrossTitleRejection(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// Seed Title B and Package on Title B
	_ = titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-atlas"},
		Name:                 "Atlas",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(72 * time.Hour),
		Territories:          10,
		CurrentMasterVersion: "V1",
		OverallStatus:        models.StatusOnTrack,
	})

	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-atlas-video"},
		TitleID:                  "title-atlas",
		Component:                models.ComponentVideo,
		Language:                 "ov",
		Version:                  "v1",
		VendorID:                 "vendor_a",
		DerivedFromMasterVersion: "V1",
		Status:                   models.PackageStatusValid,
	})

	// Attempt cross-title dependency
	crossDep := &models.Dependency{
		ID:             "dep-cross-title",
		ParentID:       "pkg-video-ov",
		ChildID:        "pkg-atlas-video",
		DependencyType: models.DependencyMasterDerivation,
	}

	res := dependencies.Create(ctx, client, crossDep)
	if res.IsOk() {
		t.Fatalf("expected cross-title dependency creation to fail, but succeeded")
	}

	domErr, ok := res.Error().(*domainerrors.DomainError)
	if !ok || domErr.Code != domainerrors.CodeInvalidInput {
		t.Errorf("expected CodeInvalidInput for cross-title dependency, got: %v", res.Error())
	}
}
