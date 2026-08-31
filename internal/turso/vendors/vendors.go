package vendors

import (
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// toDomain converts an Ent Vendor entity to a domain Vendor model.
func toDomain(v *ent.Vendor) *models.Vendor {
	if v == nil {
		return nil
	}
	return &models.Vendor{
		Base: models.Base{
			ID:        v.ID,
			Metadata:  v.Metadata,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		},
		Name:            v.Name,
		Components:      v.Components,
		Markets:         v.Markets,
		HourlyRateUSD:   v.HourlyRateUsd,
		TurnaroundHours: v.TurnaroundHours,
	}
}

// toDomainList converts a slice of Ent Vendor entities to a slice of domain Vendor models.
func toDomainList(items []*ent.Vendor) []*models.Vendor {
	res := make([]*models.Vendor, len(items))
	for i, item := range items {
		res[i] = toDomain(item)
	}
	return res
}
