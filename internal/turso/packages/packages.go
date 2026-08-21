package packages

import (
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// toDomain converts an Ent MediaPackage entity to a domain Package model.
func toDomain(p *ent.MediaPackage) *models.Package {
	if p == nil {
		return nil
	}
	return &models.Package{
		Base: models.Base{
			ID:        p.ID,
			Metadata:  p.Metadata,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
		TitleID:                  p.TitleID,
		Component:                models.ComponentType(p.Component),
		Language:                 p.Language,
		Version:                  p.Version,
		VendorID:                 p.VendorID,
		DerivedFromMasterVersion: p.DerivedFromMasterVersion,
		RedeliveryCount:          p.RedeliveryCount,
		Status:                   models.PackageStatus(p.Status),
	}
}

// toDomainList converts a slice of Ent MediaPackage entities to a slice of domain Package models.
func toDomainList(items []*ent.MediaPackage) []*models.Package {
	res := make([]*models.Package, len(items))
	for i, item := range items {
		res[i] = toDomain(item)
	}
	return res
}
