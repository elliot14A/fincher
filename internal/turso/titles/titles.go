package titles

import (
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// toDomain converts an Ent Title entity to a domain Title model.
func toDomain(t *ent.Title) *models.Title {
	if t == nil {
		return nil
	}
	return &models.Title{
		Base: models.Base{
			ID:        t.ID,
			Metadata:  t.Metadata,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
		Name:                 t.Name,
		Slug:                 t.Slug,
		Type:                 models.TitleType(t.Type),
		PremiereDate:         t.PremiereDate,
		Territories:          t.Territories,
		CurrentMasterVersion: t.CurrentMasterVersion,
		OverallStatus:        models.TitleStatus(t.OverallStatus),
	}
}

// toDomainList converts a slice of Ent Title entities to a slice of domain Title models.
func toDomainList(items []*ent.Title) []*models.Title {
	res := make([]*models.Title, len(items))
	for i, item := range items {
		res[i] = toDomain(item)
	}
	return res
}
