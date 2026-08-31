package runs

import (
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// toDomainResult converts an Ent WfResult entity to a domain WfResult model.
func toDomainResult(r *ent.WfResult) models.WfResult {
	if r == nil {
		return models.WfResult{}
	}
	return models.WfResult{
		Base: models.Base{
			ID:        r.ID,
			Metadata:  r.Metadata,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		RunID:     r.RunID,
		StepID:    r.StepID,
		Judge:     r.Judge,
		Outcome:   r.Outcome,
		Rationale: r.Rationale,
		Attempt:   r.Attempt,
	}
}

// toDomainStep converts an Ent Step entity to a domain Step model.
func toDomainStep(s *ent.Step) models.Step {
	if s == nil {
		return models.Step{}
	}
	step := models.Step{
		Base: models.Base{
			ID:        s.ID,
			Metadata:  s.Metadata,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		RunID:     s.RunID,
		Name:      s.Name,
		Status:    models.StepStatus(s.Status),
		StartedAt: s.StartedAt,
		EndedAt:   s.EndedAt,
	}

	if s.Edges.Results != nil {
		step.Results = make([]models.WfResult, len(s.Edges.Results))
		for i, res := range s.Edges.Results {
			step.Results[i] = toDomainResult(res)
		}
	}
	return step
}

// toDomain converts an Ent Run entity to a domain Run model.
func toDomain(r *ent.Run) *models.Run {
	if r == nil {
		return nil
	}
	run := &models.Run{
		Base: models.Base{
			ID:        r.ID,
			Metadata:  r.Metadata,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		TitleSlug: r.TitleSlug,
		Trigger:   r.Trigger,
		Status:    models.RunStatus(r.Status),
		StartedAt: r.StartedAt,
		EndedAt:   r.EndedAt,
	}

	if r.Edges.Steps != nil {
		run.Steps = make([]models.Step, len(r.Edges.Steps))
		for i, s := range r.Edges.Steps {
			run.Steps[i] = toDomainStep(s)
		}
	}

	if r.Edges.Results != nil {
		run.Results = make([]models.WfResult, len(r.Edges.Results))
		for i, res := range r.Edges.Results {
			run.Results[i] = toDomainResult(res)
		}
	}

	return run
}

// toDomainList converts a slice of Ent Run entities to domain Run models.
func toDomainList(items []*ent.Run) []*models.Run {
	res := make([]*models.Run, len(items))
	for i, item := range items {
		res[i] = toDomain(item)
	}
	return res
}
