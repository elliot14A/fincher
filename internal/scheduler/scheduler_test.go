package scheduler_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestScheduler_ScheduleTask_CompletesAfterCompressedDuration(t *testing.T) {
	// 10ms real per domain hour
	timeScale := 10 * time.Millisecond
	sched := scheduler.NewScheduler(timeScale)

	var wg sync.WaitGroup
	wg.Add(1)

	var completedTask *scheduler.Task
	task, err := sched.ScheduleTask(
		scheduler.TaskKindPackage,
		"pkg-audio-de",
		"avatar-fire-ash",
		"vendor-deluxe",
		models.ComponentAudio,
		"PASSED",
		5.0, // 5 hours => 50ms real
		func(completed *scheduler.Task) {
			completedTask = completed
			wg.Done()
		},
	)
	if err != nil {
		t.Fatalf("unexpected error scheduling task: %v", err)
	}

	if task.Status != scheduler.TaskStatusRunning {
		t.Errorf("expected status RUNNING, got: %s", task.Status)
	}

	// Verify it shows in GetActiveTasks
	active := sched.GetActiveTasks()
	if len(active) != 1 || active[0].ID != task.ID {
		t.Errorf("expected 1 active task, got: %v", active)
	}

	wg.Wait()

	if completedTask == nil {
		t.Fatalf("expected callback to execute")
	}
	if completedTask.Status != scheduler.TaskStatusCompleted {
		t.Errorf("expected completed task status COMPLETED, got: %s", completedTask.Status)
	}

	// Verify no longer active
	activeAfter := sched.GetActiveTasks()
	if len(activeAfter) != 0 {
		t.Errorf("expected 0 active tasks after completion, got: %d", len(activeAfter))
	}
}

func TestScheduler_ScheduleTask_CancelsOnPackageReassignment(t *testing.T) {
	timeScale := 50 * time.Millisecond
	sched := scheduler.NewScheduler(timeScale)

	task1Completed := false
	task1, err := sched.ScheduleTask(
		scheduler.TaskKindPackage,
		"pkg-audio-de",
		"avatar-fire-ash",
		"vendor-slow",
		models.ComponentAudio,
		"",
		10.0, // 500ms
		func(completed *scheduler.Task) {
			task1Completed = true
		},
	)
	if err != nil {
		t.Fatalf("scheduling task 1: %v", err)
	}

	// Re-assign package to fast vendor before task1 completes
	time.Sleep(10 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	task2Completed := false
	task2, err := sched.ScheduleTask(
		scheduler.TaskKindPackage,
		"pkg-audio-de",
		"avatar-fire-ash",
		"vendor-fast",
		models.ComponentAudio,
		"",
		2.0, // 100ms
		func(completed *scheduler.Task) {
			task2Completed = true
			wg.Done()
		},
	)
	if err != nil {
		t.Fatalf("scheduling task 2: %v", err)
	}

	// Check task1 status was cancelled immediately
	t1Snapshot, _ := sched.GetTask(task1.ID)
	if t1Snapshot.Status != scheduler.TaskStatusCancelled {
		t.Errorf("expected task 1 status CANCELLED, got: %s", t1Snapshot.Status)
	}

	wg.Wait()

	// Task 2 completed
	if !task2Completed {
		t.Errorf("expected task 2 to complete")
	}

	// Wait past task 1's original timer to ensure callback never fired
	time.Sleep(550 * time.Millisecond)

	if task1Completed {
		t.Errorf("task 1 callback fired even though it was cancelled on reassignment")
	}

	t2Snapshot, _ := sched.GetTask(task2.ID)
	if t2Snapshot.Status != scheduler.TaskStatusCompleted {
		t.Errorf("expected task 2 status COMPLETED, got: %s", t2Snapshot.Status)
	}
}

func TestScheduler_CancelTasksForPackage(t *testing.T) {
	timeScale := 100 * time.Millisecond
	sched := scheduler.NewScheduler(timeScale)

	task, _ := sched.ScheduleTask(
		scheduler.TaskKindPackage,
		"pkg-cancel-target",
		"avatar-fire-ash",
		"vendor-deluxe",
		models.ComponentAudio,
		"",
		10.0,
		nil,
	)

	cancelledCount := sched.CancelTasksForPackage("pkg-cancel-target")
	if cancelledCount != 1 {
		t.Errorf("expected 1 task cancelled, got: %d", cancelledCount)
	}

	tSnapshot, _ := sched.GetTask(task.ID)
	if tSnapshot.Status != scheduler.TaskStatusCancelled {
		t.Errorf("expected status CANCELLED, got: %s", tSnapshot.Status)
	}
}

func TestScheduler_SequentialDAG_DurationAddition(t *testing.T) {
	// Validates sequential master reconform (6h) -> dubbing repair (12h)
	// Callback chaining preserves exact duration without polling jitter.
	timeScale := 5 * time.Millisecond // 5ms per hour
	sched := scheduler.NewScheduler(timeScale)

	var wg sync.WaitGroup
	wg.Add(1)

	startTime := time.Now()
	var masterFinishTime time.Time
	var dubFinishTime time.Time

	// Step 1: Master Reconform (6 hours => 30ms)
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

			// Step 2: Dependent Dubbing repair starts upon master completion (12 hours => 60ms)
			_, errDub := sched.ScheduleTask(
				scheduler.TaskKindPackage,
				"pkg-audio-de",
				"avatar-fire-ash",
				"vendor-deluxe",
				models.ComponentAudio,
				"PASSED",
				12.0,
				func(dubTask *scheduler.Task) {
					dubFinishTime = time.Now()
					wg.Done()
				},
			)
			if errDub != nil {
				t.Errorf("scheduling dependent dub task: %v", errDub)
			}
		},
	)
	if err != nil {
		t.Fatalf("scheduling master task: %v", err)
	}

	wg.Wait()

	masterDuration := masterFinishTime.Sub(startTime)
	dubDuration := dubFinishTime.Sub(masterFinishTime)
	totalDuration := dubFinishTime.Sub(startTime)

	// Master ~30ms, Dub ~60ms, Total ~90ms (within tight test tolerances checking upper bounds to prevent §3 inflation)
	if masterDuration < 25*time.Millisecond || masterDuration > 60*time.Millisecond {
		t.Errorf("expected master duration ~30ms (range [25ms, 60ms]), got: %v", masterDuration)
	}
	if dubDuration < 50*time.Millisecond || dubDuration > 95*time.Millisecond {
		t.Errorf("expected dub duration ~60ms (range [50ms, 95ms]), got: %v", dubDuration)
	}
	if totalDuration < 80*time.Millisecond || totalDuration > 150*time.Millisecond {
		t.Errorf("expected total sequential duration ~90ms (range [80ms, 150ms]), got: %v", totalDuration)
	}
}

func TestScheduler_ArmTitleDeadline_FiresCallback(t *testing.T) {
	timeScale := 5 * time.Millisecond // 5ms per hour
	sched := scheduler.NewScheduler(timeScale)

	var wg sync.WaitGroup
	wg.Add(1)

	fired := false
	premiere := time.Now().UTC().Add(4 * time.Hour) // 4h => 20ms real

	_, err := sched.ArmTitleDeadline("title-dead-1", "dead-title", premiere, func() {
		fired = true
		wg.Done()
	})
	if err != nil {
		t.Fatalf("ArmTitleDeadline failed: %v", err)
	}

	wg.Wait()

	if !fired {
		t.Error("expected deadline callback to fire")
	}
}

func TestScheduler_ArmTitleDeadline_ReArmsAndCancelsOld(t *testing.T) {
	timeScale := 20 * time.Millisecond
	sched := scheduler.NewScheduler(timeScale)

	fired1 := false
	fired2 := false
	var wg sync.WaitGroup
	wg.Add(1)

	// Arm title for 10h (200ms)
	_, _ = sched.ArmTitleDeadline("title-rearm", "rearm-title", time.Now().UTC().Add(10*time.Hour), func() {
		fired1 = true
	})

	time.Sleep(10 * time.Millisecond)

	// Re-arm for 2h (40ms)
	_, _ = sched.ArmTitleDeadline("title-rearm", "rearm-title", time.Now().UTC().Add(2*time.Hour), func() {
		fired2 = true
		wg.Done()
	})

	wg.Wait()

	if !fired2 {
		t.Error("expected second deadline callback to fire")
	}

	// Sleep past first timer
	time.Sleep(250 * time.Millisecond)
	if fired1 {
		t.Error("first deadline callback fired even though it was re-armed and cancelled")
	}
}

func TestScheduler_CancelTasksForTitle(t *testing.T) {
	sched := scheduler.NewScheduler(time.Hour)

	_, _ = sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-1", "title-slug-a", "vendor-1", models.ComponentAudio, "", 10.0, nil)
	_, _ = sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-2", "title-slug-a", "vendor-2", models.ComponentVideo, "", 10.0, nil)
	_, _ = sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-3", "title-slug-b", "vendor-3", models.ComponentSubtitle, "", 10.0, nil)

	cancelled := sched.CancelTasksForTitle("title-slug-a")
	if cancelled != 2 {
		t.Errorf("expected 2 tasks cancelled for title-slug-a, got: %d", cancelled)
	}

	tasksA := sched.GetTasksForTitle("title-slug-a")
	for _, task := range tasksA {
		if task.Status != scheduler.TaskStatusCancelled {
			t.Errorf("expected task %s to be CANCELLED, got: %s", task.ID, task.Status)
		}
	}
}

func TestScheduler_Stop_CancelsAllRunning(t *testing.T) {
	sched := scheduler.NewScheduler(time.Hour)

	_, _ = sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-stop-1", "title-a", "vendor-1", models.ComponentAudio, "", 10.0, nil)
	_, _ = sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-stop-2", "title-b", "vendor-2", models.ComponentVideo, "", 10.0, nil)

	if len(sched.GetActiveTasks()) != 2 {
		t.Fatalf("expected 2 active tasks before stop")
	}

	sched.Stop()

	if len(sched.GetActiveTasks()) != 0 {
		t.Errorf("expected 0 active tasks after Stop(), got: %d", len(sched.GetActiveTasks()))
	}
}

func TestScheduler_Deadline_ConcurrencyStorm(t *testing.T) {
	// Stresses shared mutex map across deadlines, tasks, cancellations, and queries under -race
	sched := scheduler.NewScheduler(10 * time.Millisecond)
	defer sched.Stop()

	const numGoroutines = 30
	const numIterations = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		go func(gid int) {
			defer wg.Done()
			titleSlug := fmt.Sprintf("title-storm-%d", gid%5)
			titleID := fmt.Sprintf("id-storm-%d", gid%5)

			for i := 0; i < numIterations; i++ {
				switch (gid + i) % 6 {
				case 0:
					_, _ = sched.ArmTitleDeadline(titleID, titleSlug, time.Now().UTC().Add(time.Duration(i)*time.Hour), func() {})
				case 1:
					_ = sched.CancelDeadlineForTitle(titleID)
				case 2:
					_, _ = sched.ScheduleTask(scheduler.TaskKindPackage, fmt.Sprintf("pkg-%d-%d", gid, i), titleSlug, "vendor-1", models.ComponentAudio, "PASSED", 5.0, nil)
				case 3:
					_ = sched.CancelTasksForTitle(titleSlug)
				case 4:
					_ = sched.GetActiveTasks()
				case 5:
					_ = sched.GetTasksForTitle(titleSlug)
				}
			}
		}(g)
	}

	wg.Wait()
}

func TestScheduler_Stop_IdempotencyAndPostStopSafety(t *testing.T) {
	sched := scheduler.NewScheduler(time.Hour)

	_, _ = sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-post-1", "title-1", "vendor-1", models.ComponentAudio, "", 10.0, nil)
	_, _ = sched.ArmTitleDeadline("id-post-1", "title-1", time.Now().UTC().Add(time.Hour), func() {})

	// First Stop()
	sched.Stop()
	if len(sched.GetActiveTasks()) != 0 {
		t.Errorf("expected 0 active tasks after first Stop(), got: %d", len(sched.GetActiveTasks()))
	}

	// Calling Stop() a second time must be completely safe (idempotent)
	sched.Stop()

	// Post-Stop operations must not panic
	_, errTask := sched.ScheduleTask(scheduler.TaskKindPackage, "pkg-post-2", "title-1", "vendor-1", models.ComponentAudio, "", 10.0, nil)
	if errTask != nil {
		t.Logf("post-stop schedule task safely handled: %v", errTask)
	}

	_, errDeadline := sched.ArmTitleDeadline("id-post-2", "title-1", time.Now().UTC().Add(time.Hour), func() {})
	if errDeadline != nil {
		t.Logf("post-stop arm deadline safely handled: %v", errDeadline)
	}

	_ = sched.CancelDeadlineForTitle("id-post-1")
	_ = sched.CancelTasksForTitle("title-1")
	_ = sched.GetActiveTasks()
}
