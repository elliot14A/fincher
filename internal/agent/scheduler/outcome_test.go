package scheduler_test

import (
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent/scheduler"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestDecideOutcome_ForceOverridesRNG(t *testing.T) {
	sched := scheduler.NewScheduler(time.Second)

	// Forced PASSED must always return PASS for every component type
	for _, comp := range []models.ComponentType{models.ComponentAudio, models.ComponentSubtitle, models.ComponentVideo, models.ComponentMetadata} {
		for i := 0; i < 20; i++ {
			outcome := sched.DecideOutcome("PASSED", comp)
			if outcome != scheduler.QCOutcomePass {
				t.Fatalf("expected forced PASSED to yield PASS, got: %s", outcome)
			}
			outcomeLower := sched.DecideOutcome("passed", comp)
			if outcomeLower != scheduler.QCOutcomePass {
				t.Fatalf("expected case-insensitive forced passed to yield PASS, got: %s", outcomeLower)
			}
		}
	}

	// Forced FAILED must always return FAIL for every component type
	for _, comp := range []models.ComponentType{models.ComponentAudio, models.ComponentSubtitle, models.ComponentVideo, models.ComponentMetadata} {
		for i := 0; i < 20; i++ {
			outcome := sched.DecideOutcome("FAILED", comp)
			if outcome != scheduler.QCOutcomeFail {
				t.Fatalf("expected forced FAILED to yield FAIL, got: %s", outcome)
			}
			outcomeLower := sched.DecideOutcome("failed", comp)
			if outcomeLower != scheduler.QCOutcomeFail {
				t.Fatalf("expected case-insensitive forced failed to yield FAIL, got: %s", outcomeLower)
			}
		}
	}
}

func TestDecideOutcome_SeededDeterminism(t *testing.T) {
	seed := int64(99887766)
	sched1 := scheduler.NewScheduler(time.Second, seed)
	sched2 := scheduler.NewScheduler(time.Second, seed)

	if sched1.RunVariance() != sched2.RunVariance() {
		t.Fatalf("expected identical runVariance for identical seed, got: %v vs %v", sched1.RunVariance(), sched2.RunVariance())
	}

	// Over 50 rolls, both schedulers with the same seed must produce identical outcomes
	for i := 0; i < 50; i++ {
		o1 := sched1.DecideOutcome("", models.ComponentAudio)
		o2 := sched2.DecideOutcome("", models.ComponentAudio)
		if o1 != o2 {
			t.Fatalf("roll %d mismatch between identical-seed schedulers: %s vs %s", i, o1, o2)
		}
	}

	// Different seed produces different variance
	sched3 := scheduler.NewScheduler(time.Second, int64(11223344))
	if sched1.RunVariance() == sched3.RunVariance() {
		t.Fatalf("expected different runVariance for different seed, got equal: %v", sched1.RunVariance())
	}
}

func TestDecideOutcome_ProbabilityBoundsAndVarianceClamping(t *testing.T) {
	sched := scheduler.NewScheduler(time.Second, 12345)
	variance := sched.RunVariance()

	if variance < 0.5 || variance >= 1.5 {
		t.Fatalf("expected runVariance in [0.5, 1.5), got: %v", variance)
	}

	// Base rates verified in code
	if scheduler.BaseFailureRate[models.ComponentAudio] != 0.40 {
		t.Errorf("expected Audio base failure rate 0.40, got: %v", scheduler.BaseFailureRate[models.ComponentAudio])
	}
	if scheduler.BaseFailureRate[models.ComponentSubtitle] != 0.25 {
		t.Errorf("expected Subtitle base failure rate 0.25, got: %v", scheduler.BaseFailureRate[models.ComponentSubtitle])
	}
	if scheduler.BaseFailureRate[models.ComponentVideo] != 0.15 {
		t.Errorf("expected Video base failure rate 0.15, got: %v", scheduler.BaseFailureRate[models.ComponentVideo])
	}
	if scheduler.BaseFailureRate[models.ComponentMetadata] != 0.10 {
		t.Errorf("expected Metadata base failure rate 0.10, got: %v", scheduler.BaseFailureRate[models.ComponentMetadata])
	}
}

func TestDefectEventTypeFor(t *testing.T) {
	// Audio maps to audio sync drift
	evType, sev := scheduler.DefectEventTypeFor(models.ComponentAudio)
	if evType != models.TypeAudioSyncDriftDetected || sev != models.SeverityWarn {
		t.Errorf("Audio defect mismatch: got (%s, %s)", evType, sev)
	}

	// Subtitle maps to package invalidated
	evType, sev = scheduler.DefectEventTypeFor(models.ComponentSubtitle)
	if evType != models.TypePackageInvalidated || sev != models.SeverityWarn {
		t.Errorf("Subtitle defect mismatch: got (%s, %s)", evType, sev)
	}

	// Video maps to package invalidated
	evType, sev = scheduler.DefectEventTypeFor(models.ComponentVideo)
	if evType != models.TypePackageInvalidated || sev != models.SeverityWarn {
		t.Errorf("Video defect mismatch: got (%s, %s)", evType, sev)
	}

	// Metadata and unknown map to package invalidated
	evType, sev = scheduler.DefectEventTypeFor(models.ComponentMetadata)
	if evType != models.TypePackageInvalidated || sev != models.SeverityWarn {
		t.Errorf("Metadata defect mismatch: got (%s, %s)", evType, sev)
	}
}
