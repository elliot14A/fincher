package titles

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	enttitle "github.com/elliot14A/fincher/internal/turso/ent/title"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create inserts a new title.
func Create(ctx context.Context, client *ent.Client, t *models.Title) domainerrors.Result[*models.Title] {
	if client == nil {
		return domainerrors.Err[*models.Title](turso.NewError("titles.Create", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil))
	}
	if t == nil {
		return domainerrors.Err[*models.Title](turso.NewError("titles.Create", domainerrors.CodeInvalidInput, "title cannot be nil", nil))
	}

	if t.Slug == "" {
		if t.ID != "" && strings.HasPrefix(t.ID, "title-") {
			t.Slug = strings.TrimPrefix(t.ID, "title-")
		} else if t.Name != "" {
			t.Slug = strings.ToLower(strings.ReplaceAll(t.Name, " ", "-"))
		} else {
			t.Slug = t.ID
		}
	}

	// Ensure slug uniqueness: if slug already exists in Turso, append a random entropy suffix
	exists, err := client.Title.Query().Where(enttitle.SlugEQ(t.Slug)).Exist(ctx)
	if err == nil && exists {
		t.Slug = fmt.Sprintf("%s-%s", t.Slug, uuid.NewString()[:6])
	}

	if err := t.Validate(); err != nil {
		return domainerrors.Err[*models.Title](turso.NewError("titles.Create", domainerrors.CodeInvalidInput, "invalid title data", err))
	}

	builder := client.Title.Create().
		SetID(t.ID).
		SetName(t.Name).
		SetSlug(t.Slug).
		SetType(enttitle.Type(t.Type)).
		SetPremiereDate(t.PremiereDate).
		SetTerritories(t.Territories).
		SetCurrentMasterVersion(t.CurrentMasterVersion).
		SetOverallStatus(enttitle.OverallStatus(t.OverallStatus))

	if t.Metadata != nil {
		builder.SetMetadata(t.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.Title](turso.MapEntError("titles.Create", "title", t.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}
