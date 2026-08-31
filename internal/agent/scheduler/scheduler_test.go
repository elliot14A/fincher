package scheduler_test

import (
	"sync"
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/agent/scheduler"
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
		6.0,
		func(masterTask *scheduler.Task) {
			masterFinishTime = time.Now()

			// Step 2: Dependent Dubbing repair starts upon master completion (12 hours => 60ms)
			_, errDub := sched.ScheduleTask(
				scheduler.TaskKindPackage,
				"pkg-audio-de",
				"avatar-fire-ash",
				"vendor-deluxe",
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

	// Master ~30ms, Dub ~60ms, Total ~90ms (within reasonable test tolerances)
	if masterDuration < 25*time.Millisecond {
		t.Errorf("expected master duration ~30ms, got: %v", masterDuration)
	}
	if dubDuration < 50*time.Millisecond {
		t.Errorf("expected dub duration ~60ms, got: %v", dubDuration)
	}
	if totalDuration < 80*time.Millisecond {
		t.Errorf("expected total sequential duration ~90ms, got: %v", totalDuration)
	}
}
