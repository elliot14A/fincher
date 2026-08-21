package deliveries

import (
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func toDomain(d *ent.Delivery) *models.Delivery {
	if d == nil {
		return nil
	}
	return &models.Delivery{
		ID:         d.ID,
		TitleID:    d.TitleID,
		Country:    d.Country,
		Status:     models.DeliveryStatus(d.Status),
		TargetDate: d.TargetDate,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

func toDomainList(items []*ent.Delivery) []*models.Delivery {
	res := make([]*models.Delivery, len(items))
	for i, item := range items {
		res[i] = toDomain(item)
	}
	return res
}
