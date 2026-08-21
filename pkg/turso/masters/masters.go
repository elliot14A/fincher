package masters

import (
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
)

// toDomain converts an Ent Master entity to a domain Master model.
func toDomain(m *ent.Master) *models.Master {
	if m == nil {
		return nil
	}
	return &models.Master{
		ID:                m.ID,
		TitleID:           m.TitleID,
		Version:           m.Version,
		SupersedesVersion: m.SupersedesVersion,
		CreatedAt:         m.CreatedAt,
	}
}

// toDomainList converts a slice of Ent Master entities to domain Master models.
func toDomainList(items []*ent.Master) []*models.Master {
	res := make([]*models.Master, len(items))
	for i, item := range items {
		res[i] = toDomain(item)
	}
	return res
}
