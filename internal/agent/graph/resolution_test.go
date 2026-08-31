package graph_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/agent/graph"
	tursodeliveries "github.com/elliot14A/fincher/internal/turso/deliveries"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursoruns "github.com/elliot14A/fincher/internal/turso/runs"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestGraph_ExecuteResolution_EndToEnd(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	// Seed vendor
	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-deluxe"},
		Name:            "Deluxe Audio",
		Specialty:       "AUDIO_DUBBING",
		HourlyRateUSD:   200.0,
		TurnaroundHours: 12,
	})

	titleSlug := "avatar-fire-ash-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Avatar Fire and Ash",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusAtRisk,
	})
	if titleRes.IsErr() {
		t.Fatalf("failed to create title: %v", titleRes.Error())
	}
	titleObj := titleRes.Unwrap()

	// 1. Create media packages for US: Audio (RE_QC_PENDING) + Subtitle (VALID)
	pkgAudioID := "pkg-audio-" + uuid.NewString()[:8]
	pkgRes := tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgAudioID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})
	if pkgRes.IsErr() {
		t.Fatalf("failed to create audio package: %v", pkgRes.Error())
	}

	pkgSubID := "pkg-sub-" + uuid.NewString()[:8]
	pkgSubRes := tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgSubID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentSubtitle,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusValid,
	})
	if pkgSubRes.IsErr() {
		t.Fatalf("failed to create subtitle package: %v", pkgSubRes.Error())
	}

	// 2. Create delivery for US on HOLD
	delID := "del-us-" + uuid.NewString()[:8]
	delRes := tursodeliveries.Create(ctx, client, &models.Delivery{
		Base:       models.Base{ID: delID},
		TitleID:    titleObj.ID,
		Country:    "US",
		Status:     models.DeliveryStatusHold,
		TargetDate: time.Now().UTC().Add(24 * time.Hour),
	})
	if delRes.IsErr() {
		t.Fatalf("failed to create delivery: %v", delRes.Error())
	}

	// 3. Execute Resolution workflow on clean QC return for Audio package
	runID := "run-res-" + uuid.NewString()[:8]
	deps := graph.ResolutionDeps{
		TursoClient: client,
		ClickHouse:  nil,
	}

	out, err := graph.ExecuteResolution(ctx, deps, graph.ResolutionInput{
		RunID:     runID,
		TitleSlug: titleSlug,
		PackageID: pkgAudioID,
		Event: &models.Event{
			ID:      uuid.NewString(),
			Type:    models.TypeQCInspectionCompleted,
			Subject: titleSlug,
			Data: map[string]any{
				"package_id": pkgAudioID,
				"status":     "PASSED",
			},
		},
	})
	if err != nil {
		t.Fatalf("ExecuteResolution failed: %v", err)
	}

	if out.PackageID != pkgAudioID {
		t.Errorf("expected packageID %s, got %s", pkgAudioID, out.PackageID)
	}
	if len(out.ReleasedDeliveryIDs) != 1 || out.ReleasedDeliveryIDs[0] != delID {
		t.Errorf("expected released delivery [%s], got %v", delID, out.ReleasedDeliveryIDs)
	}

	// 4. Verify package is now VALID
	updatedPkg := tursopackages.Get(ctx, client, pkgAudioID).Unwrap()
	if updatedPkg.Status != models.PackageStatusValid {
		t.Errorf("expected package status VALID, got %s", updatedPkg.Status)
	}

	// 5. Verify delivery is now READY_TO_SHIP
	updatedDel := tursodeliveries.Get(ctx, client, delID).Unwrap()
	if updatedDel.Status != models.DeliveryStatusReadyToShip {
		t.Errorf("expected delivery status READY_TO_SHIP, got %s", updatedDel.Status)
	}

	// 6. Verify title is now ON_TRACK
	updatedTitle := tursotitles.Get(ctx, client, titleObj.ID).Unwrap()
	if updatedTitle.OverallStatus != models.StatusOnTrack {
		t.Errorf("expected title status ON_TRACK, got %s", updatedTitle.OverallStatus)
	}

	// 7. Verify Run and WfResult records
	runObj := tursoruns.GetRun(ctx, client, runID).Unwrap()
	if runObj.Status != models.RunStatusCompleted {
		t.Errorf("expected run status COMPLETED, got %s", runObj.Status)
	}
	if len(runObj.Results) == 0 || runObj.Results[0].Outcome != "RESOLVED" {
		t.Errorf("expected WfResult outcome RESOLVED, got %v", runObj.Results)
	}
}

func TestGraph_ExecuteResolution_MultiTerritory_1toN_Isolation(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	// Seed vendor
	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-deluxe"},
		Name:            "Deluxe Audio",
		Specialty:       "AUDIO_DUBBING",
		HourlyRateUSD:   200.0,
		TurnaroundHours: 12,
	})

	titleSlug := "avatar-multi-territory-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Avatar Fire and Ash Multi",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(48 * time.Hour),
		Territories:          2,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusHold,
	})
	if titleRes.IsErr() {
		t.Fatalf("failed to create title: %v", titleRes.Error())
	}
	titleObj := titleRes.Unwrap()

	// 1. Create US packages: Audio (RE_QC_PENDING) + Subtitle (VALID)
	pkgUSAudioID := "pkg-us-audio-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgUSAudioID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})

	pkgUSSubID := "pkg-us-sub-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgUSSubID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentSubtitle,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusValid,
	})

	// 2. Create DE packages: Audio (RE_QC_PENDING) + Subtitle (VALID)
	pkgDEAudioID := "pkg-de-audio-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgDEAudioID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})

	pkgDESubID := "pkg-de-sub-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgDESubID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentSubtitle,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusValid,
	})

	// 3. Create two deliveries on HOLD: US and DE
	delUSID := "del-us-" + uuid.NewString()[:8]
	_ = tursodeliveries.Create(ctx, client, &models.Delivery{
		Base:       models.Base{ID: delUSID},
		TitleID:    titleObj.ID,
		Country:    "US",
		Status:     models.DeliveryStatusHold,
		TargetDate: time.Now().UTC().Add(24 * time.Hour),
	})

	delDEID := "del-de-" + uuid.NewString()[:8]
	_ = tursodeliveries.Create(ctx, client, &models.Delivery{
		Base:       models.Base{ID: delDEID},
		TitleID:    titleObj.ID,
		Country:    "DE",
		Status:     models.DeliveryStatusHold,
		TargetDate: time.Now().UTC().Add(24 * time.Hour),
	})

	deps := graph.ResolutionDeps{
		TursoClient: client,
		ClickHouse:  nil,
	}

	// 4. Clean QC return arrives ONLY for US audio package
	outUS, err := graph.ExecuteResolution(ctx, deps, graph.ResolutionInput{
		RunID:     "run-res-us",
		TitleSlug: titleSlug,
		PackageID: pkgUSAudioID,
	})
	if err != nil {
		t.Fatalf("US ExecuteResolution failed: %v", err)
	}

	// Assert ONLY US delivery is released, DE delivery is NOT released
	if len(outUS.ReleasedDeliveryIDs) != 1 || outUS.ReleasedDeliveryIDs[0] != delUSID {
		t.Fatalf("expected ONLY US delivery [%s] to be released, got: %v", delUSID, outUS.ReleasedDeliveryIDs)
	}

	// Assert DE delivery is STILL on HOLD in database
	delDEObj := tursodeliveries.Get(ctx, client, delDEID).Unwrap()
	if delDEObj.Status != models.DeliveryStatusHold {
		t.Errorf("expected DE delivery to REMAIN on HOLD, got: %s", delDEObj.Status)
	}

	// Assert Title OverallStatus is STILL on HOLD because DE is still broken
	titleMid := tursotitles.Get(ctx, client, titleObj.ID).Unwrap()
	if titleMid.OverallStatus != models.StatusHold {
		t.Errorf("expected title to REMAIN on HOLD while DE delivery is held, got: %s", titleMid.OverallStatus)
	}

	// 5. Clean QC return arrives for DE audio package
	outDE, err := graph.ExecuteResolution(ctx, deps, graph.ResolutionInput{
		RunID:     "run-res-de",
		TitleSlug: titleSlug,
		PackageID: pkgDEAudioID,
	})
	if err != nil {
		t.Fatalf("DE ExecuteResolution failed: %v", err)
	}

	if len(outDE.ReleasedDeliveryIDs) != 1 || outDE.ReleasedDeliveryIDs[0] != delDEID {
		t.Fatalf("expected DE delivery [%s] to be released, got: %v", delDEID, outDE.ReleasedDeliveryIDs)
	}

	// Assert Title OverallStatus is NOW self-healed to ON_TRACK
	titleFinal := tursotitles.Get(ctx, client, titleObj.ID).Unwrap()
	if titleFinal.OverallStatus != models.StatusOnTrack {
		t.Errorf("expected title to self-heal to ON_TRACK when all territories are valid, got: %s", titleFinal.OverallStatus)
	}
}

func TestGraph_ExecuteResolution_MissingSubtitle_StaysHeld(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-deluxe"},
		Name:            "Deluxe Audio",
		Specialty:       "AUDIO_DUBBING",
		HourlyRateUSD:   200.0,
		TurnaroundHours: 12,
	})

	titleSlug := "title-missing-sub-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Missing Subtitle Test",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusHold,
	})
	titleObj := titleRes.Unwrap()

	// Audio is RE_QC_PENDING
	pkgAudioID := "pkg-audio-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgAudioID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})

	// Subtitle package is ALSO RE_QC_PENDING (not valid yet!)
	pkgSubID := "pkg-sub-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgSubID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentSubtitle,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})

	delID := "del-us-" + uuid.NewString()[:8]
	_ = tursodeliveries.Create(ctx, client, &models.Delivery{
		Base:       models.Base{ID: delID},
		TitleID:    titleObj.ID,
		Country:    "US",
		Status:     models.DeliveryStatusHold,
		TargetDate: time.Now().UTC().Add(24 * time.Hour),
	})

	deps := graph.ResolutionDeps{
		TursoClient: client,
		ClickHouse:  nil,
	}

	// Audio returns clean QC, but Subtitle is still RE_QC_PENDING
	out, err := graph.ExecuteResolution(ctx, deps, graph.ResolutionInput{
		RunID:     "run-res-sub-held",
		TitleSlug: titleSlug,
		PackageID: pkgAudioID,
	})
	if err != nil {
		t.Fatalf("ExecuteResolution failed: %v", err)
	}

	// Delivery must NOT be released because subtitle is still RE_QC_PENDING
	if len(out.ReleasedDeliveryIDs) != 0 {
		t.Errorf("expected 0 deliveries released when subtitle is not valid, got: %v", out.ReleasedDeliveryIDs)
	}

	delObj := tursodeliveries.Get(ctx, client, delID).Unwrap()
	if delObj.Status != models.DeliveryStatusHold {
		t.Errorf("expected delivery to remain on HOLD, got: %s", delObj.Status)
	}

	titleObjAfter := tursotitles.Get(ctx, client, titleObj.ID).Unwrap()
	if titleObjAfter.OverallStatus != models.StatusHold {
		t.Errorf("expected title to remain on HOLD, got: %s", titleObjAfter.OverallStatus)
	}
}

func TestGraph_ExecuteResolution_Processing_Branch(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-deluxe"},
		Name:            "Deluxe Audio",
		Specialty:       "AUDIO_DUBBING",
		HourlyRateUSD:   200.0,
		TurnaroundHours: 12,
	})

	titleSlug := "title-processing-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Processing Branch Test",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusHold,
	})
	titleObj := titleRes.Unwrap()

	// Audio returning clean QC
	pkgAudioID := "pkg-audio-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgAudioID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})

	// Subtitle is VALID
	pkgSubID := "pkg-sub-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgSubID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentSubtitle,
		Language:                 "en-US",
		Market:                   "US",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusValid,
	})

	// A 3rd package (e.g. metadata for JP) is PENDING (in progress, not broken)
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-meta-" + uuid.NewString()[:8]},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentMetadata,
		Language:                 "ja-JP",
		Market:                   "JP",
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusPending,
	})

	delID := "del-us-" + uuid.NewString()[:8]
	_ = tursodeliveries.Create(ctx, client, &models.Delivery{
		Base:       models.Base{ID: delID},
		TitleID:    titleObj.ID,
		Country:    "US",
		Status:     models.DeliveryStatusHold,
		TargetDate: time.Now().UTC().Add(24 * time.Hour),
	})

	deps := graph.ResolutionDeps{
		TursoClient: client,
		ClickHouse:  nil,
	}

	out, err := graph.ExecuteResolution(ctx, deps, graph.ResolutionInput{
		RunID:     "run-res-proc",
		TitleSlug: titleSlug,
		PackageID: pkgAudioID,
	})
	if err != nil {
		t.Fatalf("ExecuteResolution failed: %v", err)
	}

	// US delivery is released (Audio + Subtitle are now valid)
	if len(out.ReleasedDeliveryIDs) != 1 || out.ReleasedDeliveryIDs[0] != delID {
		t.Fatalf("expected US delivery to be released, got: %v", out.ReleasedDeliveryIDs)
	}

	// Title status becomes PROCESSING (not ON_TRACK) because metadata package is PENDING
	if out.TitleStatus != string(models.StatusProcessing) {
		t.Errorf("expected TitleStatus PROCESSING, got: %s", out.TitleStatus)
	}

	titleObjAfter := tursotitles.Get(ctx, client, titleObj.ID).Unwrap()
	if titleObjAfter.OverallStatus != models.StatusProcessing {
		t.Errorf("expected title OverallStatus to be PROCESSING, got: %s", titleObjAfter.OverallStatus)
	}
}

func TestGraph_ExecuteResolution_EmptyMarket_IsObservable(t *testing.T) {
	ctx := context.Background()
	client := tursotest.NewMemoryClient(t)

	_ = tursovendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-deluxe"},
		Name:            "Deluxe Audio",
		Specialty:       "AUDIO_DUBBING",
		HourlyRateUSD:   200.0,
		TurnaroundHours: 12,
	})

	titleSlug := "title-empty-market-" + uuid.NewString()[:8]
	titleRes := tursotitles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-" + uuid.NewString()[:8]},
		Name:                 "Empty Market Observability Test",
		Slug:                 titleSlug,
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().UTC().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "V01",
		OverallStatus:        models.StatusHold,
	})
	titleObj := titleRes.Unwrap()

	// 1. Audio and Subtitle packages with EMPTY Market ("")
	pkgAudioID := "pkg-audio-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgAudioID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentAudio,
		Language:                 "en-US",
		Market:                   "", // empty market
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})

	pkgSubID := "pkg-sub-" + uuid.NewString()[:8]
	_ = tursopackages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: pkgSubID},
		TitleID:                  titleObj.ID,
		Component:                models.ComponentSubtitle,
		Language:                 "en-US",
		Market:                   "", // empty market
		Version:                  "V01",
		VendorID:                 "vendor-deluxe",
		DerivedFromMasterVersion: "V01",
		Status:                   models.PackageStatusReQCPending,
	})

	// 2. Delivery for US on HOLD
	delID := "del-us-" + uuid.NewString()[:8]
	_ = tursodeliveries.Create(ctx, client, &models.Delivery{
		Base:       models.Base{ID: delID},
		TitleID:    titleObj.ID,
		Country:    "US",
		Status:     models.DeliveryStatusHold,
		TargetDate: time.Now().UTC().Add(24 * time.Hour),
	})

	deps := graph.ResolutionDeps{
		TursoClient: client,
		ClickHouse:  nil,
	}

	runID := "run-res-empty-market"
	out, err := graph.ExecuteResolution(ctx, deps, graph.ResolutionInput{
		RunID:     runID,
		TitleSlug: titleSlug,
		PackageID: pkgAudioID,
	})
	if err != nil {
		t.Fatalf("ExecuteResolution failed: %v", err)
	}

	// Unchanged outcome: 0 deliveries released, delivery stays HOLD, title stays HOLD
	if len(out.ReleasedDeliveryIDs) != 0 {
		t.Errorf("expected 0 deliveries released with empty market, got: %v", out.ReleasedDeliveryIDs)
	}

	delObj := tursodeliveries.Get(ctx, client, delID).Unwrap()
	if delObj.Status != models.DeliveryStatusHold {
		t.Errorf("expected delivery to remain on HOLD, got: %s", delObj.Status)
	}

	titleObjAfter := tursotitles.Get(ctx, client, titleObj.ID).Unwrap()
	if titleObjAfter.OverallStatus != models.StatusHold {
		t.Errorf("expected title to remain on HOLD, got: %s", titleObjAfter.OverallStatus)
	}

	// Verify Step metadata contains held_deliveries with reason "no_relevant_packages"
	runObj := tursoruns.GetRun(ctx, client, runID).Unwrap()
	if len(runObj.Steps) == 0 {
		t.Fatalf("expected run to contain at least 1 step")
	}

	step := runObj.Steps[0]
	heldEntries, ok := step.Metadata["held_deliveries"]
	if !ok || heldEntries == nil {
		t.Fatalf("expected step metadata to contain 'held_deliveries', got: %v", step.Metadata)
	}

	// Unmarshal or assert heldEntries
	heldSlice, ok := heldEntries.([]any)
	if !ok {
		// Or []map[string]any
		if directSlice, okDirect := heldEntries.([]map[string]any); okDirect {
			if len(directSlice) == 0 || directSlice[0]["reason"] != "no_relevant_packages" {
				t.Errorf("expected reason 'no_relevant_packages', got: %v", directSlice)
			}
			return
		}
		t.Fatalf("unexpected type for held_deliveries: %T (%v)", heldEntries, heldEntries)
	}

	if len(heldSlice) == 0 {
		t.Fatalf("expected non-empty held_deliveries slice")
	}

	firstEntry, ok := heldSlice[0].(map[string]any)
	if !ok {
		t.Fatalf("expected held entry map[string]any, got: %T", heldSlice[0])
	}

	if firstEntry["reason"] != "no_relevant_packages" {
		t.Errorf("expected reason 'no_relevant_packages', got: %v", firstEntry["reason"])
	}
	if firstEntry["country"] != "US" {
		t.Errorf("expected country 'US', got: %v", firstEntry["country"])
	}
}
