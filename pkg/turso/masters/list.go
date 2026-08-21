package masters

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	entmaster "github.com/elliot14A/fincher/pkg/ent/master"
	"github.com/elliot14A/fincher/pkg/turso"
)

// List fetches masters, optionally filtered by title_id.
func List(ctx context.Context, client *ent.Client, titleIDFilter domainerrors.Option[string]) domainerrors.Result[[]*models.Master] {
	query := client.Master.Query().Order(ent.Desc(entmaster.FieldCreatedAt))

	if titleIDFilter.IsSome() {
		query = query.Where(entmaster.TitleIDEQ(titleIDFilter.Unwrap()))
	}

	mastersList, err := query.All(ctx)
	if err != nil {
		return domainerrors.Err[[]*models.Master](turso.NewError("masters.List", domainerrors.CodeInternal, "failed to query masters", err))
	}

	return domainerrors.Ok(toDomainList(mastersList))
}
