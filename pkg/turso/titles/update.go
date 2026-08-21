package titles

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
	"github.com/elliot14A/fincher/pkg/ent"
	enttitle "github.com/elliot14A/fincher/pkg/ent/title"
	"github.com/elliot14A/fincher/pkg/turso"
)

// Update modifies an existing title using partial input.
func Update(ctx context.Context, client *ent.Client, id string, input *models.UpdateTitleInput) domainerrors.Result[*models.Title] {
	if err := input.Validate(); err != nil {
		return domainerrors.Err[*models.Title](turso.NewError("titles.Update", domainerrors.CodeInvalidInput, "invalid title update input", err))
	}

	builder := client.Title.UpdateOneID(id)

	if input.Name != nil {
		builder.SetName(*input.Name)
	}
	if input.Type != nil {
		builder.SetType(enttitle.Type(*input.Type))
	}
	if input.PremiereDate != nil {
		builder.SetPremiereDate(*input.PremiereDate)
	}
	if input.Territories != nil {
		builder.SetTerritories(*input.Territories)
	}
	if input.CurrentMasterVersion != nil {
		builder.SetCurrentMasterVersion(*input.CurrentMasterVersion)
	}
	if input.OverallStatus != nil {
		builder.SetOverallStatus(enttitle.OverallStatus(*input.OverallStatus))
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Title](turso.MapEntError("titles.Update", "title", id, err))
	}

	return domainerrors.Ok(toDomain(updated))
}
