package titles

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	enttitle "github.com/elliot14A/fincher/internal/turso/ent/title"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create inserts a new title.
func Create(ctx context.Context, client *ent.Client, t *models.Title) domainerrors.Result[*models.Title] {
	if err := t.Validate(); err != nil {
		return domainerrors.Err[*models.Title](turso.NewError("titles.Create", domainerrors.CodeInvalidInput, "invalid title data", err))
	}

	created, err := client.Title.Create().
		SetID(t.ID).
		SetName(t.Name).
		SetType(enttitle.Type(t.Type)).
		SetPremiereDate(t.PremiereDate).
		SetTerritories(t.Territories).
		SetCurrentMasterVersion(t.CurrentMasterVersion).
		SetOverallStatus(enttitle.OverallStatus(t.OverallStatus)).
		Save(ctx)

	if err != nil {
		return domainerrors.Err[*models.Title](turso.MapEntError("titles.Create", "title", t.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}
