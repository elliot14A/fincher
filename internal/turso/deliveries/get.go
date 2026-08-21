package deliveries

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Get fetches a single delivery by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Delivery] {
	d, err := client.Delivery.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Delivery](turso.MapEntError("deliveries.Get", "delivery", id, err))
	}
	return domainerrors.Ok(toDomain(d))
}
