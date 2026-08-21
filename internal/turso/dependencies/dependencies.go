package dependencies

import (
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func toDomain(d *ent.Dependency) *models.Dependency {
	if d == nil {
		return nil
	}
	return &models.Dependency{
		ID:             d.ID,
		ParentID:       d.ParentID,
		ChildID:        d.ChildID,
		DependencyType: models.DependencyType(d.DependencyType),
		CreatedAt:      d.CreatedAt,
	}
}

func toDomainList(items []*ent.Dependency) []*models.Dependency {
	res := make([]*models.Dependency, len(items))
	for i, item := range items {
		res[i] = toDomain(item)
	}
	return res
}
