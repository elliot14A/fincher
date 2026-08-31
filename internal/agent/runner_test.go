package agent_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/agent/scheduler"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/runs"
	"github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestRunActionPlan(t *testing.T) {
	client := tursotest.NewMemoryClient(t)
	defer client.Close()

	ctx := context.Background()

	premiere := time.Now().Add(48 * time.Hour)
	titleRes := titles.Create(ctx, client, &models.Title{
		Base: models.Base{
			ID: "title-ecl",
		},
		Name:                 "Eclipse",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         premiere,
		Territories:          1,
		CurrentMasterVersion: "master-v1",
		OverallStatus:        models.StatusProcessing,
	})
	if titleRes.IsErr() {
		t.Fatalf("create title failed: %v", titleRes.Error())
	}
	title := titleRes.Unwrap()

	delRes := deliveries.Create(ctx, client, &models.Delivery{
		Base: models.Base{
			ID: "del-germany-1",
		},
		TitleID:    title.ID,
		Country:    "DE",
		Status:     models.DeliveryStatusReadyToShip,
		TargetDate: premiere,
	})
	if delRes.IsErr() {
		t.Fatalf("create delivery failed: %v", delRes.Error())
	}
	del := delRes.Unwrap()

	vBerlinRes := vendors.Create(ctx, client, &models.Vendor{
		Base: models.Base{
			ID: "vendor-berlin",
		},
		Name:            "Berlin Synchron",
		Components:      []string{"AUDIO"},
		Markets:         []string{"de-DE"},
		HourlyRateUSD:   120.0,
		TurnaroundHours: 24,
	})
	if vBerlinRes.IsErr() {
		t.Fatalf("create vendor failed: %v", vBerlinRes.Error())
	}
	vendorBerlin := vBerlinRes.Unwrap()

	vOldRes := vendors.Create(ctx, client, &models.Vendor{
		Base: models.Base{
			ID: "vendor-old",
		},
		Name:            "Old Vendor",
		Components:      []string{"AUDIO"},
		Markets:         []string{"de-DE"},
		HourlyRateUSD:   80.0,
		TurnaroundHours: 48,
	})
	if vOldRes.IsErr() {
		t.Fatalf("create old vendor failed: %v", vOldRes.Error())
	}

	pkgRes := packages.Create(ctx, client, &models.Package{
		Base: models.Base{
			ID: "pkg-dub-1",
		},
		TitleID:                  title.ID,
		VendorID:                 "vendor-old",
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Version:                  "v1.0",
		DerivedFromMasterVersion: "master-v1",
		Status:                   models.PackageStatusValid,
	})
	if pkgRes.IsErr() {
		t.Fatalf("create package failed: %v", pkgRes.Error())
	}
	pkg := pkgRes.Unwrap()

	runRes := runs.CreateRun(ctx, client, &models.Run{
		Base: models.Base{
			ID: "run-exec-1",
		},
		TitleSlug: "eclipse",
		Trigger:   "ANOMALY_SIGNAL",
		Status:    models.RunStatusRunning,
		StartedAt: time.Now().UTC(),
	})
	if runRes.IsErr() {
		t.Fatalf("create run failed: %v", runRes.Error())
	}

	stepRes := runs.CreateStep(ctx, client, &models.Step{
		Base: models.Base{
			ID: "step-exec-1",
		},
		RunID:     "run-exec-1",
		Name:      "runner",
		Status:    models.StepStatusRunning,
		StartedAt: time.Now().UTC(),
	})
	if stepRes.IsErr() {
		t.Fatalf("create step failed: %v", stepRes.Error())
	}

	plan := &models.ActionPlan{
		TitleSlug: "eclipse",
		Summary:   "Hold delivery and reassign to Berlin Synchron",
		Actions: []models.Action{
			{
				Type:     models.ActionHoldDelivery,
				TargetID: del.ID,
				Reason:   "Audio drift in German track",
			},
			{
				Type:     models.ActionReassignVendor,
				TargetID: vendorBerlin.ID,
				Reason:   "Reassigning to reliable vendor",
				Payload:  map[string]any{"package_id": pkg.ID},
			},
			{
				Type:     models.ActionEmailVendor,
				TargetID: vendorBerlin.ID,
				Reason:   "Expedited dub request",
				Payload:  map[string]any{"subject": "Urgent Reconform", "body": "Please deliver in 24h"},
			},
			{
				Type:     models.ActionNotifyStakeholders,
				TargetID: "slack-ops",
				Reason:   "Inform leadership of hold",
			},
		},
	}

	execRes := agent.RunActionPlan(ctx, client, nil, "run-exec-1", "step-exec-1", plan)
	if execRes.IsErr() {
		t.Fatalf("RunActionPlan failed: %v", execRes.Error())
	}
	result := execRes.Unwrap()

	if len(result.ExecutedActions) != 4 {
		t.Errorf("expected 4 executed actions, got: %d", len(result.ExecutedActions))
	}
	if len(result.Artifacts) != 2 {
		t.Errorf("expected 2 mock communication artifacts, got: %d", len(result.Artifacts))
	}

	updatedDelRes := deliveries.Get(ctx, client, del.ID)
	if updatedDelRes.IsErr() {
		t.Fatalf("get delivery failed: %v", updatedDelRes.Error())
	}
	if updatedDelRes.Unwrap().Status != models.DeliveryStatusHold {
		t.Errorf("expected delivery status HOLD, got: %s", updatedDelRes.Unwrap().Status)
	}

	updatedPkgRes := packages.Get(ctx, client, pkg.ID)
	if updatedPkgRes.IsErr() {
		t.Fatalf("get package failed: %v", updatedPkgRes.Error())
	}
	if updatedPkgRes.Unwrap().VendorID != vendorBerlin.ID {
		t.Errorf("expected vendor reassigned to %s, got: %s", vendorBerlin.ID, updatedPkgRes.Unwrap().VendorID)
	}

	updatedRunRes := runs.GetRun(ctx, client, "run-exec-1")
	if updatedRunRes.IsErr() {
		t.Fatalf("get run failed: %v", updatedRunRes.Error())
	}
	if updatedRunRes.Unwrap().Status != models.RunStatusCompleted {
		t.Errorf("expected run status COMPLETED, got: %s", updatedRunRes.Unwrap().Status)
	}
}

func TestRunActionPlan_ForcedPass_EmitsQCCompleted(t *testing.T) {
	client := tursotest.NewMemoryClient(t)
	defer client.Close()
	ctx := context.Background()

	_ = titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-1"},
		Name:                 "Test Title",
		Slug:                 "avatar-fire-ash",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "master-v1",
		OverallStatus:        models.StatusProcessing,
	})

	_ = vendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-test"},
		Name:            "Test Vendor",
		Components:      []string{"AUDIO"},
		Markets:         []string{"de-DE"},
		TurnaroundHours: 1,
	})

	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-pass-test"},
		TitleID:                  "title-1",
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "v1.0",
		DerivedFromMasterVersion: "master-v1",
		Status:                   models.PackageStatusValid,
		VendorID:                 "vendor-test",
		RedeliveryCount:          0,
	})

	sched := scheduler.NewScheduler(time.Millisecond) // 1ms timescale
	var emittedEvent models.Event
	var wg sync.WaitGroup
	wg.Add(1)

	plan := &models.ActionPlan{
		TitleSlug: "avatar-fire-ash",
		Actions: []models.Action{
			{
				Type:     models.ActionReassignVendor,
				TargetID: "vendor-test",
				Payload: map[string]any{
					"package_id":    "pkg-pass-test",
					"force_outcome": "PASSED",
				},
			},
		},
	}

	execRes := agent.RunActionPlanWithDeps(ctx, agent.RunnerDeps{
		TursoClient: client,
		Scheduler:   sched,
		OnScheduleComplete: func(ev models.Event) {
			emittedEvent = ev
			wg.Done()
		},
	}, "run-pass-1", "step-pass-1", plan)

	if execRes.IsErr() {
		t.Fatalf("RunActionPlanWithDeps failed: %v", execRes.Error())
	}

	wg.Wait()

	if emittedEvent.Type != models.TypeQCInspectionCompleted {
		t.Fatalf("expected TypeQCInspectionCompleted, got: %s", emittedEvent.Type)
	}
	if emittedEvent.Data["status"] != "PASSED" {
		t.Errorf("expected data.status PASSED, got: %v", emittedEvent.Data["status"])
	}
}

func TestRunActionPlan_ForcedFail_UnderCap_IncrementsRedeliveryAndEmitsDefect(t *testing.T) {
	client := tursotest.NewMemoryClient(t)
	defer client.Close()
	ctx := context.Background()

	_ = titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-1"},
		Name:                 "Test Title",
		Slug:                 "avatar-fire-ash",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "master-v1",
		OverallStatus:        models.StatusProcessing,
	})

	_ = vendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-test"},
		Name:            "Test Vendor",
		Components:      []string{"AUDIO"},
		Markets:         []string{"de-DE"},
		TurnaroundHours: 1,
	})

	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-fail-test"},
		TitleID:                  "title-1",
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "v1.0",
		DerivedFromMasterVersion: "master-v1",
		Status:                   models.PackageStatusValid,
		VendorID:                 "vendor-test",
		RedeliveryCount:          0,
	})

	sched := scheduler.NewScheduler(time.Millisecond)
	var emittedEvent models.Event
	var wg sync.WaitGroup
	wg.Add(1)

	plan := &models.ActionPlan{
		TitleSlug: "avatar-fire-ash",
		Actions: []models.Action{
			{
				Type:     models.ActionReassignVendor,
				TargetID: "vendor-test",
				Payload: map[string]any{
					"package_id":    "pkg-fail-test",
					"force_outcome": "FAILED",
				},
			},
		},
	}

	execRes := agent.RunActionPlanWithDeps(ctx, agent.RunnerDeps{
		TursoClient: client,
		Scheduler:   sched,
		OnScheduleComplete: func(ev models.Event) {
			emittedEvent = ev
			wg.Done()
		},
	}, "run-fail-1", "step-fail-1", plan)

	if execRes.IsErr() {
		t.Fatalf("RunActionPlanWithDeps failed: %v", execRes.Error())
	}

	wg.Wait()

	// Audio defect emitted
	if emittedEvent.Type != models.TypeAudioSyncDriftDetected {
		t.Fatalf("expected TypeAudioSyncDriftDetected on audio failure, got: %s", emittedEvent.Type)
	}
	if emittedEvent.Severity != models.SeverityWarn {
		t.Errorf("expected SeverityWarn, got: %s", emittedEvent.Severity)
	}

	// Redelivery count incremented in Turso
	pRes := packages.Get(ctx, client, "pkg-fail-test")
	if pRes.IsErr() || pRes.Unwrap().RedeliveryCount != 1 {
		t.Fatalf("expected RedeliveryCount 1, got: %v (err: %v)", pRes.Unwrap().RedeliveryCount, pRes.Error())
	}
}

func TestRunActionPlan_ForcedFail_AtCap_EmitsSLABreach(t *testing.T) {
	client := tursotest.NewMemoryClient(t)
	defer client.Close()
	ctx := context.Background()

	_ = titles.Create(ctx, client, &models.Title{
		Base:                 models.Base{ID: "title-1"},
		Name:                 "Test Title",
		Slug:                 "avatar-fire-ash",
		Type:                 models.TitleTypeFeature,
		PremiereDate:         time.Now().Add(48 * time.Hour),
		Territories:          1,
		CurrentMasterVersion: "master-v1",
		OverallStatus:        models.StatusProcessing,
	})

	_ = vendors.Create(ctx, client, &models.Vendor{
		Base:            models.Base{ID: "vendor-test"},
		Name:            "Test Vendor",
		Components:      []string{"AUDIO"},
		Markets:         []string{"de-DE"},
		TurnaroundHours: 1,
	})

	// Package already at MaxRedeliveryAttempts (3)
	_ = packages.Create(ctx, client, &models.Package{
		Base:                     models.Base{ID: "pkg-cap-test"},
		TitleID:                  "title-1",
		Component:                models.ComponentAudio,
		Language:                 "de-DE",
		Market:                   "DE",
		Version:                  "v1.0",
		DerivedFromMasterVersion: "master-v1",
		Status:                   models.PackageStatusValid,
		VendorID:                 "vendor-test",
		RedeliveryCount:          3,
	})

	sched := scheduler.NewScheduler(time.Millisecond)
	var emittedEvent models.Event
	var wg sync.WaitGroup
	wg.Add(1)

	plan := &models.ActionPlan{
		TitleSlug: "avatar-fire-ash",
		Actions: []models.Action{
			{
				Type:     models.ActionReassignVendor,
				TargetID: "vendor-test",
				Payload: map[string]any{
					"package_id":    "pkg-cap-test",
					"force_outcome": "FAILED",
				},
			},
		},
	}

	execRes := agent.RunActionPlanWithDeps(ctx, agent.RunnerDeps{
		TursoClient: client,
		Scheduler:   sched,
		OnScheduleComplete: func(ev models.Event) {
			emittedEvent = ev
			wg.Done()
		},
	}, "run-cap-1", "step-cap-1", plan)

	if execRes.IsErr() {
		t.Fatalf("RunActionPlanWithDeps failed: %v", execRes.Error())
	}

	wg.Wait()

	if emittedEvent.Type != models.TypeVendorSLABreach {
		t.Fatalf("expected TypeVendorSLABreach on cap exceeded, got: %s", emittedEvent.Type)
	}
	if emittedEvent.Severity != models.SeverityCritical {
		t.Errorf("expected SeverityCritical on SLA breach, got: %s", emittedEvent.Severity)
	}
	if emittedEvent.Data["reason"] != "redelivery_cap_exceeded" {
		t.Errorf("expected reason redelivery_cap_exceeded, got: %v", emittedEvent.Data["reason"])
	}
}
