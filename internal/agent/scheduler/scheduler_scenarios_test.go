package scheduler_test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent/scheduler"
	"github.com/elliot14A/fincher/internal/config"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// B1. Concurrent schedule + cancel storm (the -race catcher)
func TestScheduler_B1_ConcurrentScheduleCancelStorm(t *testing.T) {
	timeScale := 500 * time.Microsecond // fast timers firing during storm
	sched := scheduler.NewScheduler(timeScale)

	const (
		numWorkers  = 50
		numPackages = 5
		iterations  = 30
	)

	var totalCompletions atomic.Int64
	var activeScheduleSucceeded atomic.Int64
	var firedTasks sync.Map
	var allScheduled sync.Map
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(int64(workerID*1000 + 42)))

			for i := 0; i < iterations; i++ {
				pkgID := fmt.Sprintf("pkg-storm-%d", rng.Intn(numPackages))
				op := rng.Intn(3)

				switch op {
				case 0, 1: // Schedule task with short duration (1-5h => 0.5-2.5ms)
					turnaround := float64(rng.Intn(5) + 1)
					task, err := sched.ScheduleTask(
						scheduler.TaskKindPackage,
						pkgID,
						"avatar-fire-ash",
						fmt.Sprintf("vendor-%d", workerID),
						models.ComponentAudio,
						"PASSED",
						turnaround,
						func(completed *scheduler.Task) {
							if completed.Status != scheduler.TaskStatusCompleted {
								t.Errorf("callback received non-completed status: %s", completed.Status)
							}
							if completed.CompletedAt == nil {
								t.Errorf("completed task has nil CompletedAt")
							}
							// Invariant: at most one completion per task ID
							if _, alreadyFired := firedTasks.LoadOrStore(completed.ID, completed); alreadyFired {
								t.Errorf("duplicate callback fired for task ID %s", completed.ID)
							}
							totalCompletions.Add(1)
						},
					)
					if err == nil {
						activeScheduleSucceeded.Add(1)
						allScheduled.Store(task.ID, task)
					}

				case 2: // Cancel in-flight tasks for package
					sched.CancelTasksForPackage(pkgID)
				}

				// Small jitter to interleave scheduling and cancellations
				time.Sleep(time.Duration(rng.Intn(200)) * time.Microsecond)
			}
		}(w)
	}

	wg.Wait()

	// Drain any remaining in-flight tasks
	time.Sleep(20 * time.Millisecond)

	// All tasks must have reached terminal state (none running)
	activeAfter := sched.GetActiveTasks()
	if len(activeAfter) != 0 {
		t.Errorf("expected 0 active tasks after drain, got: %d", len(activeAfter))
	}

	// Completions must not exceed total scheduled
	completions := totalCompletions.Load()
	schedules := activeScheduleSucceeded.Load()
	if completions > schedules {
		t.Errorf("completions (%d) exceeded successful schedules (%d)", completions, schedules)
	}

	// Invariant: For every task that fired onComplete, its final status in scheduler is COMPLETED (never CANCELLED)
	firedTasks.Range(func(key, value any) bool {
		taskID := key.(string)
		if snap, exists := sched.GetTask(taskID); exists {
			if snap.Status != scheduler.TaskStatusCompleted {
				t.Errorf("task %s fired callback but scheduler records status: %s", taskID, snap.Status)
			}
			if snap.CancelledAt != nil {
				t.Errorf("task %s fired callback but has non-nil CancelledAt", taskID)
			}
		}
		return true
	})
}

// B2. Fail->reassign->fail chain reaching the cap
func TestScheduler_B2_FailReassignChainReachingCap(t *testing.T) {
	timeScale := time.Millisecond // 1ms per hour
	// Seed chosen so Audio component failure occurs
	sched := scheduler.NewScheduler(timeScale, 20260828)

	pkgID := "pkg-audio-chain"
	var completedCount atomic.Int64
	var wg sync.WaitGroup

	maxAttempts := config.MaxRedeliveryAttempts // 3
	wg.Add(maxAttempts)

	var runChain func(attempt int)
	runChain = func(attempt int) {
		if attempt > maxAttempts {
			return
		}

		_, err := sched.ScheduleTask(
			scheduler.TaskKindPackage,
			pkgID,
			"avatar-fire-ash",
			fmt.Sprintf("vendor-attempt-%d", attempt),
			models.ComponentAudio,
			"FAILED", // force failure to model defect loop
			1.0,      // 1ms real duration
			func(completed *scheduler.Task) {
				completedCount.Add(1)
				wg.Done()

				if attempt < maxAttempts {
					// Trigger next reassignment in chain
					runChain(attempt + 1)
				}
			},
		)
		if err != nil {
			t.Errorf("attempt %d schedule error: %v", attempt, err)
		}
	}

	runChain(1)
	wg.Wait()

	if completedCount.Load() != int64(maxAttempts) {
		t.Errorf("expected %d completed iterations, got: %d", maxAttempts, completedCount.Load())
	}
}

// B3. Sequential DAG with a mid-chain cancellation (master->dub, dub reassigned mid-flight)
func TestScheduler_B3_SequentialDAG_MidChainCancellation(t *testing.T) {
	timeScale := 5 * time.Millisecond // 5ms per hour
	sched := scheduler.NewScheduler(timeScale)

	var wg sync.WaitGroup
	wg.Add(1)

	startTime := time.Now()
	var masterFinishTime time.Time
	var originalDubFired atomic.Bool
	var fastDubFinishTime time.Time

	// Step 1: Schedule master reconform (6 hours => 30ms)
	_, err := sched.ScheduleTask(
		scheduler.TaskKindMaster,
		"master-V02",
		"avatar-fire-ash",
		"vendor-editorial",
		models.ComponentVideo,
		"PASSED",
		6.0,
		func(masterTask *scheduler.Task) {
			masterFinishTime = time.Now()

			// Step 2: Schedule slow dub (20 hours => 100ms)
			slowTask, _ := sched.ScheduleTask(
				scheduler.TaskKindPackage,
				"pkg-audio-de",
				"avatar-fire-ash",
				"vendor-slow",
				models.ComponentAudio,
				"PASSED",
				20.0,
				func(t *scheduler.Task) {
					originalDubFired.Store(true)
				},
			)

			// Step 3: Mid-flight cancellation after 15ms
			time.Sleep(15 * time.Millisecond)
			sched.CancelTasksForPackage("pkg-audio-de")

			snap, _ := sched.GetTask(slowTask.ID)
			if snap.Status != scheduler.TaskStatusCancelled {
				t.Errorf("expected slow task status CANCELLED, got: %s", snap.Status)
			}
			if snap.CancelledAt == nil {
				t.Errorf("expected non-nil CancelledAt on cancelled task")
			}

			// Step 4: Reassign to fast vendor (6 hours => 30ms)
			_, _ = sched.ScheduleTask(
				scheduler.TaskKindPackage,
				"pkg-audio-de",
				"avatar-fire-ash",
				"vendor-fast",
				models.ComponentAudio,
				"PASSED",
				6.0,
				func(fastTask *scheduler.Task) {
					fastDubFinishTime = time.Now()
					wg.Done()
				},
			)
		},
	)
	if err != nil {
		t.Fatalf("scheduling master task: %v", err)
	}

	wg.Wait()

	// Wait past slow task's original duration to ensure it never fires
	time.Sleep(120 * time.Millisecond)

	if originalDubFired.Load() {
		t.Fatalf("cancelled original dub task callback fired!")
	}

	masterDur := masterFinishTime.Sub(startTime)
	totalDur := fastDubFinishTime.Sub(startTime)

	// Master ~30ms [25ms, 60ms]
	if masterDur < 25*time.Millisecond || masterDur > 60*time.Millisecond {
		t.Errorf("master duration ~30ms out of bounds: %v", masterDur)
	}

	// Total ~30ms master + 15ms mid-cancel + 30ms fast dub = ~75ms [60ms, 140ms]
	if totalDur < 60*time.Millisecond || totalDur > 140*time.Millisecond {
		t.Errorf("total sequential wall time ~75ms out of bounds: %v", totalDur)
	}
}

// B4. Reaper bounds the map under sustained ambient load
func TestScheduler_B4_ReaperBoundsMapUnderAmbientLoad(t *testing.T) {
	timeScale := 10 * time.Microsecond
	sched := scheduler.NewScheduler(timeScale)

	const taskBatch = 1500
	var wg sync.WaitGroup
	wg.Add(taskBatch)

	for i := 0; i < taskBatch; i++ {
		_, err := sched.ScheduleTask(
			scheduler.TaskKindPackage,
			fmt.Sprintf("pkg-ambient-%d", i),
			"avatar-fire-ash",
			"vendor-auto",
			models.ComponentAudio,
			"PASSED",
			0.1, // 1 microsecond
			func(t *scheduler.Task) {
				wg.Done()
			},
		)
		if err != nil {
			t.Fatalf("scheduling task %d: %v", i, err)
		}
	}

	wg.Wait()

	// Verify task map is bounded under hysteresis band (maxRetainedTasks + reapBatch = 628) and did not grow to 1500
	count := sched.TaskCount()
	if count > 628 {
		t.Errorf("expected reaper to bound task count <= 628, got: %d", count)
	}
}

// B5. Force override under concurrency
func TestScheduler_B5_ForceOverrideUnderConcurrency(t *testing.T) {
	timeScale := 500 * time.Microsecond
	sched := scheduler.NewScheduler(timeScale)

	const workers = 20
	const iterations = 20
	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				force := "PASSED"
				if (workerID+i)%2 == 1 {
					force = "FAILED"
				}

				outcome := sched.DecideOutcome(force, models.ComponentAudio)
				if force == "PASSED" && outcome != scheduler.QCOutcomePass {
					t.Errorf("worker %d: forced PASSED yielded: %s", workerID, outcome)
				}
				if force == "FAILED" && outcome != scheduler.QCOutcomeFail {
					t.Errorf("worker %d: forced FAILED yielded: %s", workerID, outcome)
				}
			}
		}(w)
	}

	wg.Wait()
}

// B6. Real-mode duration sanity (overflow protection & hour-scale)
func TestScheduler_B6_RealModeDurationSanity(t *testing.T) {
	// Real mode scale: 1 hour timescale
	sched := scheduler.NewScheduler(time.Hour)

	// 36h turnaround
	task36, err := sched.ScheduleTask(
		scheduler.TaskKindPackage,
		"pkg-real-36",
		"avatar-fire-ash",
		"vendor-deluxe",
		models.ComponentAudio,
		"PASSED",
		36.0,
		nil,
	)
	if err != nil {
		t.Fatalf("scheduling 36h task: %v", err)
	}

	expected36Duration := 36 * time.Hour
	actual36Duration := task36.FinishReal.Sub(task36.StartedAt)
	if actual36Duration != expected36Duration {
		t.Errorf("expected 36h duration %v, got: %v", expected36Duration, actual36Duration)
	}

	// Pathological turnaround hours check (must clamp without overflow or panic)
	pathologicalHours := 1e12
	taskPathological, err := sched.ScheduleTask(
		scheduler.TaskKindPackage,
		"pkg-pathological",
		"avatar-fire-ash",
		"vendor-deluxe",
		models.ComponentAudio,
		"PASSED",
		pathologicalHours,
		nil,
	)
	if err != nil {
		t.Fatalf("scheduling pathological task: %v", err)
	}

	if taskPathological.FinishReal.Before(taskPathological.StartedAt) {
		t.Errorf("pathological duration overflowed into the past: %v < %v", taskPathological.FinishReal, taskPathological.StartedAt)
	}
}
