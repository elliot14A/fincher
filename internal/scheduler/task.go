package scheduler

import (
	"time"

	"github.com/elliot14A/fincher/pkg/domain/models"
)

// TaskKind represents the domain entity being repaired or conformed.
type TaskKind string

const (
	TaskKindPackage       TaskKind = "package"
	TaskKindMaster        TaskKind = "master"
	TaskKindMasterQC      TaskKind = "master_qc"
	TaskKindTitleDeadline TaskKind = "title_deadline"
)

// TaskStatus represents the lifecycle of a compressed-time simulation task.
type TaskStatus string

const (
	TaskStatusRunning   TaskStatus = "RUNNING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusCancelled TaskStatus = "CANCELLED"
)

// Task represents an active or completed background repair job.
type Task struct {
	ID              string               `json:"id"`
	Kind            TaskKind             `json:"kind"`
	TargetID        string               `json:"target_id"` // package ID or master version
	TitleSlug       string               `json:"title_slug"`
	VendorID        string               `json:"vendor_id"`
	Component       models.ComponentType `json:"component,omitempty"`
	ForceOutcome    string               `json:"force_outcome,omitempty"`
	TurnaroundHours float64              `json:"turnaround_hours"`
	StartedAt       time.Time            `json:"started_at"`
	FinishReal      time.Time            `json:"finish_real"`
	CompletedAt     *time.Time           `json:"completed_at,omitempty"`
	CancelledAt     *time.Time           `json:"cancelled_at,omitempty"`
	Status          TaskStatus           `json:"status"`

	timer *time.Timer
}

// Snapshot returns a clean snapshot of the task with unexported concurrency primitives omitted.
func (t *Task) Snapshot() *Task {
	if t == nil {
		return nil
	}
	var compAt *time.Time
	if t.CompletedAt != nil {
		c := *t.CompletedAt
		compAt = &c
	}
	var cancAt *time.Time
	if t.CancelledAt != nil {
		c := *t.CancelledAt
		cancAt = &c
	}

	return &Task{
		ID:              t.ID,
		Kind:            t.Kind,
		TargetID:        t.TargetID,
		TitleSlug:       t.TitleSlug,
		VendorID:        t.VendorID,
		Component:       t.Component,
		ForceOutcome:    t.ForceOutcome,
		TurnaroundHours: t.TurnaroundHours,
		StartedAt:       t.StartedAt,
		FinishReal:      t.FinishReal,
		CompletedAt:     compAt,
		CancelledAt:     cancAt,
		Status:          t.Status,
	}
}
