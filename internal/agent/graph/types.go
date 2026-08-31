package graph

import (
	"database/sql"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/agent/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"google.golang.org/adk/v2/model"
)

// IncidentGraphDeps supplies runtime dependencies to the incident workflow.
type IncidentGraphDeps struct {
	Model              model.LLM
	TursoClient        *ent.Client
	ClickHouse         *sql.DB
	MaxAttempts        int
	Scheduler          *scheduler.Scheduler
	OnScheduleComplete func(event models.Event)
}

// IncidentInput is the entry payload passed to the incident workflow.
type IncidentInput struct {
	RunID              string        `json:"run_id,omitempty"`
	Event              *models.Event `json:"event"`
	HoursUntilPremiere float64       `json:"hours_until_premiere"`
}

// IncidentOutput captures the complete end-to-end outcome of an incident run.
type IncidentOutput struct {
	Actionable   bool                       `json:"actionable"`
	Decision     agent.VerificationDecision `json:"decision"`
	Rationale    string                     `json:"rationale"`
	ActionPlan   *models.ActionPlan         `json:"action_plan,omitempty"`
	RunnerResult *agent.RunnerResult        `json:"runner_result,omitempty"`
	Attempts     int                        `json:"attempts"`
}

// AllocationGraphDeps supplies dependencies to the vendor allocation workflow.
type AllocationGraphDeps struct {
	Model       model.LLM
	TursoClient *ent.Client
	ClickHouse  *sql.DB
}

// AllocationInput is the entry payload for holistic vendor allocation.
type AllocationInput struct {
	RunID              string                         `json:"run_id,omitempty"`
	TitleSlug          string                         `json:"title_slug"`
	Requirements       []models.AllocationRequirement `json:"requirements"`
	HoursUntilPremiere float64                        `json:"hours_until_premiere"`
}

// AllocationOutput captures the holistic staffing plan.
type AllocationOutput struct {
	Plan     *models.AllocationPlan   `json:"plan,omitempty"`
	Decision *agent.SelectionDecision `json:"decision,omitempty"`
}
