package vendors

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create inserts a new vendor.
func Create(ctx context.Context, client *ent.Client, v *models.Vendor) domainerrors.Result[*models.Vendor] {
	if err := v.Validate(); err != nil {
		return domainerrors.Err[*models.Vendor](turso.NewError("vendors.Create", domainerrors.CodeInvalidInput, "invalid vendor data", err))
	}

	builder := client.Vendor.Create().
		SetID(v.ID).
		SetName(v.Name).
		SetComponents(v.Components).
		SetMarkets(v.Markets).
		SetHourlyRateUsd(v.HourlyRateUSD).
		SetTurnaroundHours(v.TurnaroundHours)

	if v.Metadata != nil {
		builder.SetMetadata(v.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Vendor](turso.MapEntError("vendors.Create", "vendor", v.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}
