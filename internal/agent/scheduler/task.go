package scheduler

import (
	"context"
	"time"
)

// TaskKind represents the domain entity being repaired or conformed.
type TaskKind string

const (
	TaskKindPackage TaskKind = "package"
	TaskKindMaster  TaskKind = "master"
)

// TaskStatus represents the lifecycle of a compressed-time simulation task.
type TaskStatus string

const (
	TaskStatusScheduled TaskStatus = "SCHEDULED"
	TaskStatusRunning   TaskStatus = "RUNNING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusCancelled TaskStatus = "CANCELLED"
)

// Task represents an active or completed background repair job.
type Task struct {
	ID              string     `json:"id"`
	Kind            TaskKind   `json:"kind"`
	TargetID        string     `json:"target_id"` // package ID or master version
	TitleSlug       string     `json:"title_slug"`
	VendorID        string     `json:"vendor_id"`
	TurnaroundHours float64    `json:"turnaround_hours"`
	StartedAt       time.Time  `json:"started_at"`
	FinishReal      time.Time  `json:"finish_real"`
	Status          TaskStatus `json:"status"`

	cancelFunc context.CancelFunc
	timer      *time.Timer
}
