package graph_test

import (
	"context"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/agent/graph"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/deliveries"
	"github.com/elliot14A/fincher/internal/turso/packages"
	"github.com/elliot14A/fincher/internal/turso/runs"
	"github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/internal/turso/tursotest"
	"github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestExecuteIncident(t *testing.T) {
	ctx := context.Background()

	t.Run("Filters routine benign event without executing action plan", func(t *testing.T) {
		llm := &mockLLM{
			responses: []string{
				`{
					"actionable": false,
					"severity": "INFO",
					"anomaly_type": "NONE",
					"rationale": "Routine inspection completed within acceptable parameters."
				}`,
			},
		}

		event := &models.Event{
			ID:       "evt-benign-1",
			Type:     models.TypeQCInspectionCompleted,
			Severity: models.SeverityInfo,
			Subject:  "eclipse",
			Data:     map[string]any{"package_id": "pkg-1"},
		}

		client := tursotest.NewMemoryClient(t)
		defer client.Close()

		deps := graph.IncidentGraphDeps{
			Model:       llm,
			TursoClient: client,
		}

		output, err := graph.ExecuteIncident(ctx, deps, graph.IncidentInput{
			Event:              event,
			HoursUntilPremiere: 72.0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.Actionable {
			t.Errorf("expected actionable false, got true")
		}
		if output.Decision != "FILTERED" {
			t.Errorf("expected decision FILTERED, got: %s", output.Decision)
		}
		if output.RunnerResult != nil {
			t.Errorf("expected nil runner result for filtered event")
		}
	})

	t.Run("Processes actionable defect end-to-end and mutates operational state", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)

		premiere := time.Now().Add(48 * time.Hour)
		tRes := titles.Create(ctx, client, &models.Title{
			Base: models.Base{
				ID: "title-eclipse",
			},
			Name:                 "Eclipse",
			Type:                 models.TitleTypeFeature,
			PremiereDate:         premiere,
			Territories:          1,
			CurrentMasterVersion: "master-v1",
			OverallStatus:        models.StatusProcessing,
		})
		if tRes.IsErr() {
			t.Fatalf("failed to create title: %v", tRes.Error())
		}
		title := tRes.Unwrap()

		vRes := vendors.Create(ctx, client, &models.Vendor{
			Base: models.Base{
				ID: "vendor-berlin",
			},
			Name:            "Berlin Synchron",
			Components:      []string{"AUDIO"},
			Markets:         []string{"de-DE"},
			HourlyRateUSD:   120.0,
			TurnaroundHours: 24,
		})
		if vRes.IsErr() {
			t.Fatalf("failed to create vendor: %v", vRes.Error())
		}

		pRes := packages.Create(ctx, client, &models.Package{
			Base: models.Base{
				ID: "pkg-german-dub",
			},
			TitleID:                  title.ID,
			VendorID:                 "vendor-berlin",
			Component:                models.ComponentAudio,
			Language:                 "de-DE",
			Version:                  "v1.0",
			DerivedFromMasterVersion: "master-v1",
			Status:                   models.PackageStatusValid,
		})
		if pRes.IsErr() {
			t.Fatalf("failed to create package: %v", pRes.Error())
		}

		dRes := deliveries.Create(ctx, client, &models.Delivery{
			Base: models.Base{
				ID: "del-germany",
			},
			TitleID:    title.ID,
			Country:    "DE",
			Status:     models.DeliveryStatusReadyToShip,
			TargetDate: premiere,
		})
		if dRes.IsErr() {
			t.Fatalf("failed to create delivery: %v", dRes.Error())
		}

		filterOutput := `{
			"actionable": true,
			"severity": "CRITICAL",
			"anomaly_type": "AUDIO_SYNC_DRIFT",
			"rationale": "Audio sync drift of +380ms detected on German localized dub track."
		}`

		planOutput := `{
			"title_slug": "eclipse",
			"summary": "Hold German broadcast delivery and reassign to Berlin Synchron for expedited conform.",
			"actions": [
				{
					"type": "HOLD_DELIVERY",
					"target_id": "del-germany",
					"reason": "Hold broadcast package delivery while sync drift is corrected.",
					"payload": {}
				},
				{
					"type": "REASSIGN_VENDOR",
					"target_id": "vendor-berlin",
					"reason": "Dispatch expedited audio reconform.",
					"payload": {
						"package_id": "pkg-german-dub"
					}
				}
			]
		}`

		llm := &mockLLM{
			responses: []string{filterOutput, planOutput},
		}

		event := &models.Event{
			ID:       "evt-defect-sync-1",
			Type:     models.TypeAudioSyncDriftDetected,
			Severity: models.SeverityCritical,
			Subject:  "eclipse",
			Data: map[string]any{
				"package_id": "pkg-german-dub",
				"vendor_id":  "vendor-berlin",
				"component":  "AUDIO",
			},
		}

		deps := graph.IncidentGraphDeps{
			Model:       llm,
			TursoClient: client,
			MaxAttempts: 3,
		}

		output, err := graph.ExecuteIncident(ctx, deps, graph.IncidentInput{
			Event:              event,
			HoursUntilPremiere: 48.0,
		})
		if err != nil {
			t.Fatalf("ExecuteIncident returned error: %v", err)
		}

		if !output.Actionable {
			t.Errorf("expected actionable true, got false")
		}
		if output.Decision != agent.DecisionApproved {
			t.Errorf("expected APPROVED decision, got: %s", output.Decision)
		}
		if output.RunnerResult == nil {
			t.Fatalf("expected non-nil RunnerResult")
		}
		if len(output.RunnerResult.ExecutedActions) != 2 {
			t.Errorf("expected 2 executed actions, got: %d", len(output.RunnerResult.ExecutedActions))
		}

		// Verify state mutation
		delCheck := deliveries.Get(ctx, client, "del-germany")
		if delCheck.IsErr() {
			t.Fatalf("failed to fetch delivery del-germany: %v", delCheck.Error())
		}
		if delCheck.Unwrap().Status != models.DeliveryStatusHold {
			t.Errorf("expected delivery status HOLD, got: %s", delCheck.Unwrap().Status)
		}

		// Verify that all 4 lifecycle stages were persisted as Steps in Turso
		runCheck := runs.GetRun(ctx, client, "run-"+event.ID)
		if runCheck.IsErr() {
			t.Fatalf("failed to fetch run: %v", runCheck.Error())
		}
		loadedRun := runCheck.Unwrap()
		if len(loadedRun.Steps) != 4 {
			t.Fatalf("expected 4 lifecycle steps, got: %d", len(loadedRun.Steps))
		}
		stepNames := []string{"triage_judge", "context_gathering", "remediation_loop", "remediation_executor"}
		for i, s := range loadedRun.Steps {
			if s.Name != stepNames[i] {
				t.Errorf("step %d: expected name %s, got %s", i, stepNames[i], s.Name)
			}
			if s.Status != models.StepStatusCompleted {
				t.Errorf("step %d: expected status COMPLETED, got %s", i, s.Status)
			}
		}
		if len(loadedRun.Results) < 2 {
			t.Errorf("expected at least 2 WfResults, got: %d", len(loadedRun.Results))
		}
	})

	t.Run("Escalates to human operator when plan is repeatedly rejected", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)

		filterOutput := `{
			"actionable": true,
			"severity": "CRITICAL",
			"anomaly_type": "AUDIO_SYNC_DRIFT",
			"rationale": "Audio sync drift detected."
		}`

		// Contradictory action plan that will be rejected by policy
		invalidPlanOutput := `{
			"title_slug": "eclipse",
			"summary": "Contradictory plan holding and releasing same target.",
			"actions": [
				{
					"type": "HOLD_DELIVERY",
					"target_id": "del-germany",
					"reason": "Hold delivery",
					"payload": {}
				},
				{
					"type": "RELEASE_DELIVERY",
					"target_id": "del-germany",
					"reason": "Release delivery",
					"payload": {}
				}
			]
		}`

		llm := &mockLLM{
			responses: []string{filterOutput, invalidPlanOutput, invalidPlanOutput, invalidPlanOutput},
		}

		event := &models.Event{
			ID:       "evt-defect-contradictory",
			Type:     models.TypeAudioSyncDriftDetected,
			Severity: models.SeverityCritical,
			Subject:  "eclipse",
			Data: map[string]any{
				"package_id": "pkg-german-dub",
			},
		}

		deps := graph.IncidentGraphDeps{
			Model:       llm,
			TursoClient: client,
			MaxAttempts: 3,
		}

		output, err := graph.ExecuteIncident(ctx, deps, graph.IncidentInput{
			Event:              event,
			HoursUntilPremiere: 48.0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Decision != agent.DecisionEscalate {
			t.Errorf("expected ESCALATE decision, got: %s", output.Decision)
		}
		if output.Attempts != 3 {
			t.Errorf("expected 3 attempts, got: %d", output.Attempts)
		}

		// Verify run status in Turso is ESCALATED
		runCheck := runs.GetRun(ctx, client, "run-"+event.ID)
		if runCheck.IsErr() {
			t.Fatalf("failed to fetch run: %v", runCheck.Error())
		}
		if runCheck.Unwrap().Status != models.RunStatusEscalated {
			t.Errorf("expected run status ESCALATED, got: %s", runCheck.Unwrap().Status)
		}
	})

	t.Run("Cancels in-flight tasks for title when master cut revision event arrives", func(t *testing.T) {
		client := tursotest.NewMemoryClient(t)
		defer client.Close()

		premiere := time.Now().Add(48 * time.Hour)
		tRes := titles.Create(ctx, client, &models.Title{
			Base:                 models.Base{ID: "title-master-rev"},
			Name:                 "Master Rev Title",
			Slug:                 "master-rev-title",
			Type:                 models.TitleTypeFeature,
			PremiereDate:         premiere,
			Territories:          1,
			CurrentMasterVersion: "V01",
			OverallStatus:        models.StatusProcessing,
		})
		if tRes.IsErr() {
			t.Fatalf("failed to create title: %v", tRes.Error())
		}

		sched := scheduler.NewScheduler(time.Hour)
		defer sched.Stop()

		// Schedule 2 in-flight repair tasks for this title
		t1, _ := sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-1", "master-rev-title", "vendor-1", models.ComponentAudio, "", 10.0, nil)
		t2, _ := sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-2", "master-rev-title", "vendor-2", models.ComponentVideo, "", 10.0, nil)
		// Schedule 1 task for a DIFFERENT title (must not be cancelled)
		t3, _ := sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-3", "other-title", "vendor-3", models.ComponentSubtitle, "", 10.0, nil)

		filterOutput := `{
			"actionable": true,
			"severity": "CRITICAL",
			"anomaly_type": "MASTER_REVISED",
			"rationale": "Master cut revised from V01 to V02."
		}`
		planOutput := `{
			"title_slug": "master-rev-title",
			"summary": "Hold title while master is reconformed.",
			"actions": [
				{
					"type": "HOLD_TITLE",
					"target_id": "title-master-rev",
					"reason": "Master cut revised."
				}
			]
		}`

		llm := &mockLLM{
			responses: []string{filterOutput, planOutput},
		}

		event := &models.Event{
			ID:       "evt-master-rev-1",
			Type:     models.TypeMasterCutRevised,
			Severity: models.SeverityCritical,
			Subject:  "master-rev-title",
			Data: map[string]any{
				"new_version": "V02",
				"old_version": "V01",
			},
		}

		deps := graph.IncidentGraphDeps{
			Model:       llm,
			TursoClient: client,
			Scheduler:   sched,
			MaxAttempts: 3,
		}

		output, err := graph.ExecuteIncident(ctx, deps, graph.IncidentInput{
			Event:              event,
			HoursUntilPremiere: 48.0,
		})
		if err != nil {
			t.Fatalf("ExecuteIncident returned error: %v", err)
		}
		if !output.Actionable || output.Decision != agent.DecisionApproved {
			t.Fatalf("expected actionable approved incident, got: %s", output.Decision)
		}

		// Verify tasks for master-rev-title were moved to CANCELLED
		s1, _ := sched.GetTask(t1.ID)
		s2, _ := sched.GetTask(t2.ID)
		s3, _ := sched.GetTask(t3.ID)

		if s1.Status != scheduler.TaskStatusCancelled {
			t.Errorf("expected t1 to be CANCELLED, got: %s", s1.Status)
		}
		if s2.Status != scheduler.TaskStatusCancelled {
			t.Errorf("expected t2 to be CANCELLED, got: %s", s2.Status)
		}
		if s3.Status != scheduler.TaskStatusRunning {
			t.Errorf("expected t3 for other-title to remain RUNNING, got: %s", s3.Status)
		}
	})
}
