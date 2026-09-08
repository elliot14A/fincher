package graph

import (
	"database/sql"

	"github.com/elliot14A/fincher/internal/agent"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/mcp"
	"google.golang.org/adk/v2/model"
)

type IncidentGraphDeps struct {
	Model              model.LLM
	TursoClient        *ent.Client
	TursoDB            *sql.DB
	ClickHouse         *sql.DB
	MCP                *mcp.Client
	MaxAttempts        int
	Scheduler          *scheduler.Scheduler
	OnScheduleComplete func(event models.Event)
}

type IncidentInput struct {
	RunID              string        `json:"run_id,omitempty"`
	Event              *models.Event `json:"event"`
	HoursUntilPremiere float64       `json:"hours_until_premiere"`
}

type IncidentOutput struct {
	Actionable   bool                       `json:"actionable"`
	Decision     agent.VerificationDecision `json:"decision"`
	Rationale    string                     `json:"rationale"`
	ActionPlan   *models.ActionPlan         `json:"action_plan,omitempty"`
	RunnerResult *agent.RunnerResult        `json:"runner_result,omitempty"`
	Attempts     int                        `json:"attempts"`
}

type AllocationGraphDeps struct {
	Model              model.LLM
	TursoClient        *ent.Client
	TursoDB            *sql.DB
	ClickHouse         *sql.DB
	MCP                *mcp.Client
	Scheduler          *scheduler.Scheduler
	OnScheduleComplete func(event models.Event)
}

type AllocationInput struct {
	RunID              string                         `json:"run_id,omitempty"`
	TitleSlug          string                         `json:"title_slug"`
	Requirements       []models.AllocationRequirement `json:"requirements"`
	HoursUntilPremiere float64                        `json:"hours_until_premiere"`
}

type AllocationOutput struct {
	Plan     *models.AllocationPlan   `json:"plan,omitempty"`
	Decision *agent.SelectionDecision `json:"decision,omitempty"`
}
