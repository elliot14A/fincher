package scheduler

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/recovery"
)

// DefaultSchedulerSeed is the fixed constant used for reproducible demo runs.
const DefaultSchedulerSeed = int64(20260828)

// Retention constants defining the bounded task history with hysteresis band.
const (
	maxRetainedTasks = 500
	reapBatch        = 128
)

// Scheduler manages in-memory compressed-time task execution and DAG duration chaining.
type Scheduler struct {
	mu          sync.RWMutex
	tasks       map[string]*Task
	timeScale   time.Duration
	rng         *rand.Rand
	rngMu       sync.Mutex
	runVariance float64
}

// NewScheduler constructs a thread-safe scheduler instance with the specified time compression factor and optional seed.
func NewScheduler(timeScale time.Duration, seed ...int64) *Scheduler {
	if timeScale <= 0 {
		timeScale = time.Second
	}
	seedVal := DefaultSchedulerSeed
	if len(seed) > 0 && seed[0] != 0 {
		seedVal = seed[0]
	}

	r := rand.New(rand.NewSource(seedVal))
	// runVariance in range [0.5, 1.5), drawn once per scheduler instance from seeded RNG
	variance := 0.5 + r.Float64()

	return &Scheduler{
		tasks:       make(map[string]*Task),
		timeScale:   timeScale,
		rng:         r,
		runVariance: variance,
	}
}

// TimeScale returns the current time compression scale factor.
func (s *Scheduler) TimeScale() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.timeScale
}

// RunVariance returns the per-instance variance multiplier drawn at initialization.
func (s *Scheduler) RunVariance() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.runVariance
}

// TaskCount returns the total number of retained tasks in the scheduler map.
func (s *Scheduler) TaskCount() int {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks)
}

// reapLocked prunes oldest terminal tasks (Completed/Cancelled) when total count exceeds maxRetainedTasks + reapBatch.
// Must be called with s.mu write lock held. RUNNING tasks are never reaped.
func (s *Scheduler) reapLocked() {
	if len(s.tasks) < maxRetainedTasks+reapBatch {
		return
	}

	var terminal []*Task
	for _, t := range s.tasks {
		if t.Status != TaskStatusRunning {
			terminal = append(terminal, t)
		}
	}
	if len(terminal) == 0 {
		return
	}

	// Sort terminal tasks by StartedAt ascending (oldest first)
	sort.Slice(terminal, func(i, j int) bool {
		return terminal[i].StartedAt.Before(terminal[j].StartedAt)
	})

	excess := len(s.tasks) - maxRetainedTasks
	for i := 0; i < excess && i < len(terminal); i++ {
		delete(s.tasks, terminal[i].ID)
	}
}

// ScheduleTask registers and launches a compressed-time task.
// If a previous active task exists for the same package, it is cancelled to prevent duplicate completions.
func (s *Scheduler) ScheduleTask(
	kind TaskKind,
	targetID, titleSlug, vendorID string,
	component models.ComponentType,
	forceOutcome string,
	turnaroundHours float64,
	onComplete func(t *Task),
) (*Task, error) {
	if s == nil {
		return nil, nil
	}
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
			if existing.Kind == TaskKindPackage && existing.TargetID == targetID && existing.Status == TaskStatusRunning {
				if existing.timer != nil {
					existing.timer.Stop()
				}
				nowCanc := time.Now().UTC()
				existing.Status = TaskStatusCancelled
				existing.CancelledAt = &nowCanc
				logger.Debug("scheduler: cancelled existing in-flight task on package re-assignment",
					"task_id", existing.ID,
					"package_id", targetID,
				)
			}
		}
	}

	taskID := fmt.Sprintf("task-%s-%s", kind, uuid.NewString()[:8])
	now := time.Now().UTC()

	// Safe duration calculation: guarded against float overflow and rounded
	raw := turnaroundHours * float64(s.timeScale)
	if raw > float64(math.MaxInt64) {
		raw = float64(math.MaxInt64)
	}
	realDuration := time.Duration(math.Round(raw))
	if realDuration < 0 {
		realDuration = 0
	}
	finishReal := now.Add(realDuration)

	task := &Task{
		ID:              taskID,
		Kind:            kind,
		TargetID:        targetID,
		TitleSlug:       titleSlug,
		VendorID:        vendorID,
		Component:       component,
		ForceOutcome:    forceOutcome,
		TurnaroundHours: turnaroundHours,
		StartedAt:       now,
		FinishReal:      finishReal,
		Status:          TaskStatusRunning,
	}

	// Schedule execution callback with atomic status transition and immutable snapshot delivery
	task.timer = time.AfterFunc(realDuration, func() {
		recovery.WrapPanic(fmt.Sprintf("scheduler.callback.task=%s.target=%s", taskID, targetID), func() {
			s.mu.Lock()
			shouldFire := (task.Status == TaskStatusRunning)
			var snap *Task
			if shouldFire {
				nowComp := time.Now().UTC()
				task.Status = TaskStatusCompleted
				task.CompletedAt = &nowComp
				snap = task.Snapshot()
			}
			s.mu.Unlock()

			if !shouldFire {
				return
			}

			if onComplete != nil {
				onComplete(snap)
			}
		}, nil)
	})

	s.tasks[taskID] = task
	s.reapLocked()

	logger.Debug("scheduler: scheduled compressed-time task",
		"task_id", task.ID,
		"kind", task.Kind,
		"target_id", task.TargetID,
		"turnaround_hours", turnaroundHours,
		"real_duration", realDuration.String(),
		"finish_real", finishReal.Format(time.RFC3339),
	)

	return task.Snapshot(), nil
}

// CancelTasksForPackage cancels all in-flight tasks for a specific package.
func (s *Scheduler) CancelTasksForPackage(packageID string) int {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	cancelled := 0
	for _, t := range s.tasks {
		if t.Kind == TaskKindPackage && t.TargetID == packageID && t.Status == TaskStatusRunning {
			if t.timer != nil {
				t.timer.Stop()
			}
			nowCanc := time.Now().UTC()
			t.Status = TaskStatusCancelled
			t.CancelledAt = &nowCanc
			cancelled++
		}
	}
	return cancelled
}

// CancelTasksForTitle cancels all in-flight tasks for a specific title slug.
func (s *Scheduler) CancelTasksForTitle(titleSlug string) int {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	cancelled := 0
	for _, t := range s.tasks {
		if t.TitleSlug == titleSlug && t.Status == TaskStatusRunning {
			if t.timer != nil {
				t.timer.Stop()
			}
			nowCanc := time.Now().UTC()
			t.Status = TaskStatusCancelled
			t.CancelledAt = &nowCanc
			cancelled++
		}
	}
	return cancelled
}

// ArmTitleDeadline arms or replaces a one-shot deadline timer for a title.
// Note: In-memory timers are not rehydrated on process restart (acceptable for hackathon runtime simulation).
func (s *Scheduler) ArmTitleDeadline(
	titleID, titleSlug string,
	premiereDate time.Time,
	onBreach func(),
) (*Task, error) {
	if s == nil {
		return nil, nil
	}
	if titleID == "" {
		return nil, domainerrors.NewWithOp("scheduler.ArmTitleDeadline", domainerrors.CodeInvalidInput, "title_id cannot be empty", nil)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Cancel existing deadline timer for this titleID (idempotency guard)
	for _, existing := range s.tasks {
		if existing.Kind == TaskKindTitleDeadline && existing.TargetID == titleID && existing.Status == TaskStatusRunning {
			if existing.timer != nil {
				existing.timer.Stop()
			}
			nowCanc := time.Now().UTC()
			existing.Status = TaskStatusCancelled
			existing.CancelledAt = &nowCanc
			logger.Debug("scheduler: cancelled existing deadline timer on title re-arm",
				"task_id", existing.ID,
				"title_id", titleID,
			)
		}
	}

	now := time.Now().UTC()
	hoursUntil := premiereDate.Sub(now).Hours()
	if hoursUntil < 0 {
		hoursUntil = 0
	}

	raw := hoursUntil * float64(s.timeScale)
	if raw > float64(math.MaxInt64) {
		raw = float64(math.MaxInt64)
	}
	realDuration := time.Duration(math.Round(raw))
	if realDuration < 0 {
		realDuration = 0
	}
	finishReal := now.Add(realDuration)

	taskID := fmt.Sprintf("task-deadline-%s", uuid.NewString()[:8])
	task := &Task{
		ID:              taskID,
		Kind:            TaskKindTitleDeadline,
		TargetID:        titleID,
		TitleSlug:       titleSlug,
		TurnaroundHours: hoursUntil,
		StartedAt:       now,
		FinishReal:      finishReal,
		Status:          TaskStatusRunning,
	}

	task.timer = time.AfterFunc(realDuration, func() {
		recovery.WrapPanic(fmt.Sprintf("scheduler.deadline.task=%s.title=%s", taskID, titleSlug), func() {
			s.mu.Lock()
			shouldFire := (task.Status == TaskStatusRunning)
			if shouldFire {
				nowComp := time.Now().UTC()
				task.Status = TaskStatusCompleted
				task.CompletedAt = &nowComp
			}
			s.mu.Unlock()

			if !shouldFire {
				return
			}

			if onBreach != nil {
				onBreach()
			}
		}, nil)
	})

	s.tasks[taskID] = task
	s.reapLocked()

	logger.Debug("scheduler: armed title deadline timer",
		"task_id", task.ID,
		"title_id", titleID,
		"title_slug", titleSlug,
		"premiere_date", premiereDate.Format(time.RFC3339),
		"hours_until", hoursUntil,
		"real_duration", realDuration.String(),
	)

	return task.Snapshot(), nil
}

// CancelDeadlineForTitle cancels active deadline timer for a title.
func (s *Scheduler) CancelDeadlineForTitle(titleID string) int {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	cancelled := 0
	for _, t := range s.tasks {
		if t.Kind == TaskKindTitleDeadline && t.TargetID == titleID && t.Status == TaskStatusRunning {
			if t.timer != nil {
				t.timer.Stop()
			}
			nowCanc := time.Now().UTC()
			t.Status = TaskStatusCancelled
			t.CancelledAt = &nowCanc
			cancelled++
		}
	}
	return cancelled
}

// Stop cancels all running timers and shuts down the scheduler.
func (s *Scheduler) Stop() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	for _, t := range s.tasks {
		if t.Status == TaskStatusRunning {
			if t.timer != nil {
				t.timer.Stop()
			}
			t.Status = TaskStatusCancelled
			t.CancelledAt = &now
		}
	}
}

// GetTask returns a clean snapshot of a task by ID.
func (s *Scheduler) GetTask(id string) (*Task, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exists := s.tasks[id]
	if !exists {
		return nil, false
	}
	return t.Snapshot(), true
}

// GetActiveTasks returns clean snapshots of all currently running tasks.
func (s *Scheduler) GetActiveTasks() []*Task {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*Task
	for _, t := range s.tasks {
		if t.Status == TaskStatusRunning {
			active = append(active, t.Snapshot())
		}
	}
	return active
}

// GetTasksForTitle returns clean snapshots of all tasks associated with a title slug.
func (s *Scheduler) GetTasksForTitle(titleSlug string) []*Task {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*Task
	for _, t := range s.tasks {
		if t.TitleSlug == titleSlug {
			list = append(list, t.Snapshot())
		}
	}
	return list
}
