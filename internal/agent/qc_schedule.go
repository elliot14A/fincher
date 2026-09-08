package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/elliot14A/fincher/internal/config"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursopackages "github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
)

type QCScheduleDeps struct {
	TursoClient        *ent.Client
	Scheduler          SchedulerInterface
	OnScheduleComplete func(event models.Event)
	LogContext         string
}

func BuildQCCompletionCallback(ctx context.Context, deps QCScheduleDeps) func(t *scheduler.Task) {
	return func(t *scheduler.Task) {
		var outcome scheduler.QCOutcome
		if deps.Scheduler != nil {
			outcome = deps.Scheduler.DecideOutcome(t.ForceOutcome, t.Component)
		} else {
			outcome = scheduler.QCOutcomePass
		}

		if outcome == scheduler.QCOutcomePass {
			pRes := tursopackages.Get(ctx, deps.TursoClient, t.TargetID)
			if pRes.IsOk() {
				pkg := pRes.Unwrap()
				tRes := tursotitles.Get(ctx, deps.TursoClient, pkg.TitleID)
				if tRes.IsOk() {
					title := tRes.Unwrap()
					if pkg.IsStaleAgainst(title.CurrentMasterVersion) {
						logger.Warn("qc: discarding QC pass for stale package against revised master",
							"package_id", pkg.ID,
							"derived_from", pkg.DerivedFromMasterVersion,
							"active_master", title.CurrentMasterVersion,
						)
						return
					}
				}
			}

			validStatus := models.PackageStatusValid
			updPkgRes := tursopackages.Update(ctx, deps.TursoClient, t.TargetID, &models.UpdatePackageInput{
				Status: &validStatus,
			})
			if updPkgRes.IsErr() {
				logger.Error("qc: failed to update package status to VALID on QC pass",
					"log_context", deps.LogContext,
					"package_id", t.TargetID,
					"task_id", t.ID,
					"error", updPkgRes.Error(),
				)
			}

			qcEvent := models.Event{
				ID:              fmt.Sprintf("evt-qc-%s", t.ID),
				Source:          "fincher/qc.agent",
				Type:            models.TypeQCInspectionCompleted,
				Subject:         t.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityInfo,
				DataContentType: "application/json",
				Data: map[string]any{
					"package_id": t.TargetID,
					"status":     "PASSED",
					"vendor_id":  t.VendorID,
				},
			}
			if deps.OnScheduleComplete != nil {
				deps.OnScheduleComplete(qcEvent)
			}
			return
		}

		currentRedelivery := 0
		pRes := tursopackages.Get(ctx, deps.TursoClient, t.TargetID)
		if pRes.IsOk() {
			currentRedelivery = pRes.Unwrap().RedeliveryCount
		}

		if currentRedelivery >= config.MaxRedeliveryAttempts {
			slaEvent := models.Event{
				ID:              fmt.Sprintf("evt-sla-breach-%s", t.ID),
				Source:          "fincher/qc.agent",
				Type:            models.TypeVendorSLABreach,
				Subject:         t.TitleSlug,
				Time:            time.Now().UTC(),
				Severity:        models.SeverityCritical,
				DataContentType: "application/json",
				Data: map[string]any{
					"package_id":       t.TargetID,
					"vendor_id":        t.VendorID,
					"reason":           "redelivery_cap_exceeded",
					"redelivery_count": currentRedelivery,
				},
			}
			if deps.OnScheduleComplete != nil {
				deps.OnScheduleComplete(slaEvent)
			}
			return
		}

		newRedelivery := currentRedelivery + 1
		updRedelivRes := tursopackages.Update(ctx, deps.TursoClient, t.TargetID, &models.UpdatePackageInput{
			RedeliveryCount: &newRedelivery,
		})
		if updRedelivRes.IsErr() {
			logger.Error("qc: failed to update package redelivery count",
				"log_context", deps.LogContext,
				"package_id", t.TargetID,
				"task_id", t.ID,
				"new_count", newRedelivery,
				"error", updRedelivRes.Error(),
			)
		}

		defectEventType, defectSeverity := scheduler.DefectEventTypeFor(t.Component)
		defectData := map[string]any{
			"package_id":       t.TargetID,
			"vendor_id":        t.VendorID,
			"defect_type":      "REPAIR_INSPECTION_FAILED",
			"redelivery_count": newRedelivery,
		}
		if t.Component == models.ComponentAudio {
			defectData["drift_ms"] = 110.0
			defectData["defect_type"] = "AUDIO_SYNC_DRIFT"
		}

		defectEvent := models.Event{
			ID:              fmt.Sprintf("evt-defect-%s-%d", t.ID, newRedelivery),
			Source:          "fincher/qc.agent",
			Type:            defectEventType,
			Subject:         t.TitleSlug,
			Time:            time.Now().UTC(),
			Severity:        defectSeverity,
			DataContentType: "application/json",
			Data:            defectData,
		}
		if deps.OnScheduleComplete != nil {
			deps.OnScheduleComplete(defectEvent)
		}
	}
}
