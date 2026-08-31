package models

import (
	"time"
)

// RunStatus represents the lifecycle state of an agent workflow run.
type RunStatus string

const (
	RunStatusPending   RunStatus = "PENDING"
	RunStatusRunning   RunStatus = "RUNNING"
	RunStatusCompleted RunStatus = "COMPLETED"
	RunStatusFailed    RunStatus = "FAILED"
	RunStatusEscalated RunStatus = "ESCALATED"
)

// Run represents an execution of an agent workflow.
type Run struct {
	Base
	TitleSlug string     `json:"title_slug,omitempty"`
	Trigger   string     `json:"trigger" validate:"required"`
	Status    RunStatus  `json:"status" validate:"required"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Steps     []Step     `json:"steps,omitempty"`
	Results   []WfResult `json:"results,omitempty"`
}

// Validate verifies Run constraints.
func (r *Run) Validate() error {
	if err := r.ValidateMetadata(); err != nil {
		return err
	}
	if r.TitleSlug == "" {
		r.TitleSlug = DefaultTitleAgnosticSentinel
	}
	return validate.Struct(r)
}

// StepStatus represents the execution state of an individual workflow step.
type StepStatus string

const (
	StepStatusPending   StepStatus = "PENDING"
	StepStatusRunning   StepStatus = "RUNNING"
	StepStatusCompleted StepStatus = "COMPLETED"
	StepStatusFailed    StepStatus = "FAILED"
	StepStatusSkipped   StepStatus = "SKIPPED"
)

// Step represents a node execution within a workflow Run.
type Step struct {
	Base
	RunID     string     `json:"run_id" validate:"required"`
	Name      string     `json:"name" validate:"required"`
	Status    StepStatus `json:"status" validate:"required"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Results   []WfResult `json:"results,omitempty"`
}

// Validate verifies Step constraints.
func (s *Step) Validate() error {
	if err := s.ValidateMetadata(); err != nil {
		return err
	}
	return validate.Struct(s)
}

// WfResult records a judge decision, verdict, or node evaluation outcome.
type WfResult struct {
	Base
	RunID     string `json:"run_id" validate:"required"`
	StepID    string `json:"step_id,omitempty"`
	Judge     string `json:"judge" validate:"required"`
	Outcome   string `json:"outcome" validate:"required"`
	Rationale string `json:"rationale"`
	Attempt   int    `json:"attempt"`
}

// Validate verifies WfResult constraints.
func (w *WfResult) Validate() error {
	if err := w.ValidateMetadata(); err != nil {
		return err
	}
	return validate.Struct(w)
}
