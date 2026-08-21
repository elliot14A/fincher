package titles

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	enttitle "github.com/elliot14A/fincher/pkg/ent/title"
	"github.com/elliot14A/fincher/pkg/turso"
)

// List fetches all titles, optionally filtered by overall status.
func List(ctx context.Context, client *ent.Client, statusFilter domainerrors.Option[models.TitleStatus]) domainerrors.Result[[]*models.Title] {
	query := client.Title.Query().Order(ent.Asc(enttitle.FieldPremiereDate))

	if statusFilter.IsSome() {
		query = query.Where(enttitle.OverallStatusEQ(enttitle.OverallStatus(statusFilter.Unwrap())))
	}

	titlesList, err := query.All(ctx)
	if err != nil {
		return domainerrors.Err[[]*models.Title](turso.NewError("titles.List", domainerrors.CodeInternal, "failed to query titles", err))
	}

	return domainerrors.Ok(toDomainList(titlesList))
}
