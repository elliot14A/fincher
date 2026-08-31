package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
)

// Scheduler manages in-memory compressed-time task execution and DAG duration chaining.
type Scheduler struct {
	mu        sync.RWMutex
	tasks     map[string]*Task
	timeScale time.Duration
}

// NewScheduler constructs a thread-safe scheduler instance with the specified time compression factor.
func NewScheduler(timeScale time.Duration) *Scheduler {
	if timeScale <= 0 {
		timeScale = time.Second
	}
	return &Scheduler{
		tasks:     make(map[string]*Task),
		timeScale: timeScale,
	}
}

// TimeScale returns the current time compression scale factor.
func (s *Scheduler) TimeScale() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.timeScale
}

// ScheduleTask registers and launches a compressed-time task.
// If a previous active task exists for the same package, it is cancelled to prevent duplicate completions.
func (s *Scheduler) ScheduleTask(
	kind TaskKind,
	targetID, titleSlug, vendorID string,
	turnaroundHours float64,
	onComplete func(t *Task),
) (*Task, error) {
	if targetID == "" {
		return nil, domainerrors.NewWithOp("scheduler.ScheduleTask", domainerrors.CodeInvalidInput, "target_id cannot be empty", nil)
	}
	if turnaroundHours < 0 {
		return nil, domainerrors.NewWithOp("scheduler.ScheduleTask", domainerrors.CodeInvalidInput, "turnaround_hours cannot be negative", nil)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Cancel existing in-flight task for the same package (idempotency guard)
	if kind == TaskKindPackage {
		for _, existing := range s.tasks {
			if existing.Kind == TaskKindPackage && existing.TargetID == targetID &&
				(existing.Status == TaskStatusRunning || existing.Status == TaskStatusScheduled) {
				if existing.timer != nil {
					existing.timer.Stop()
				}
				if existing.cancelFunc != nil {
					existing.cancelFunc()
				}
				existing.Status = TaskStatusCancelled
				logger.Info("scheduler: cancelled existing in-flight task on package re-assignment",
					"task_id", existing.ID,
					"package_id", targetID,
				)
			}
		}
	}

	taskID := fmt.Sprintf("task-%s-%s", kind, uuid.NewString()[:8])
	now := time.Now().UTC()
	realDuration := time.Duration(turnaroundHours * float64(s.timeScale))
	finishReal := now.Add(realDuration)

	ctx, cancel := context.WithCancel(context.Background())

	task := &Task{
		ID:              taskID,
		Kind:            kind,
		TargetID:        targetID,
		TitleSlug:       titleSlug,
		VendorID:        vendorID,
		TurnaroundHours: turnaroundHours,
		StartedAt:       now,
		FinishReal:      finishReal,
		Status:          TaskStatusRunning,
		cancelFunc:      cancel,
	}

	// Schedule execution callback
	task.timer = time.AfterFunc(realDuration, func() {
		s.mu.Lock()
		// Check if task was cancelled while waiting
		if task.Status != TaskStatusRunning {
			s.mu.Unlock()
			return
		}
		task.Status = TaskStatusCompleted
		s.mu.Unlock()

		select {
		case <-ctx.Done():
			return
		default:
			if onComplete != nil {
				onComplete(task)
			}
		}
	})

	s.tasks[taskID] = task

	logger.Info("scheduler: scheduled compressed-time task",
		"task_id", task.ID,
		"kind", task.Kind,
		"target_id", task.TargetID,
		"turnaround_hours", turnaroundHours,
		"real_duration", realDuration.String(),
		"finish_real", finishReal.Format(time.RFC3339),
	)

	return task, nil
}

// CancelTasksForPackage cancels all in-flight tasks for a specific package.
func (s *Scheduler) CancelTasksForPackage(packageID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	cancelled := 0
	for _, t := range s.tasks {
		if t.Kind == TaskKindPackage && t.TargetID == packageID &&
			(t.Status == TaskStatusRunning || t.Status == TaskStatusScheduled) {
			if t.timer != nil {
				t.timer.Stop()
			}
			if t.cancelFunc != nil {
				t.cancelFunc()
			}
			t.Status = TaskStatusCancelled
			cancelled++
		}
	}
	return cancelled
}

// GetTask returns a snapshot of a task by ID.
func (s *Scheduler) GetTask(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exists := s.tasks[id]
	if !exists {
		return nil, false
	}
	cp := *t
	return &cp, true
}

// GetActiveTasks returns all currently running tasks.
func (s *Scheduler) GetActiveTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*Task
	for _, t := range s.tasks {
		if t.Status == TaskStatusRunning {
			cp := *t
			active = append(active, &cp)
		}
	}
	return active
}

// GetTasksForTitle returns all tasks associated with a title slug.
func (s *Scheduler) GetTasksForTitle(titleSlug string) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*Task
	for _, t := range s.tasks {
		if t.TitleSlug == titleSlug {
			cp := *t
			list = append(list, &cp)
		}
	}
	return list
}
