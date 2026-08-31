package titles

import (
	"context"
	"time"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	enttitle "github.com/elliot14A/fincher/internal/turso/ent/title"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// DefaultPremiereUrgencyHours is the standard baseline operational urgency threshold.
const DefaultPremiereUrgencyHours = 48.0

// Get fetches a single title by ID.
func Get(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.Title] {
	t, err := client.Title.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[*models.Title](turso.MapEntError("titles.Get", "title", id, err))
	}
	return domainerrors.Ok(toDomain(t))
}

// GetBySlug queries a single title strictly by its slug.
func GetBySlug(ctx context.Context, client *ent.Client, slug string) domainerrors.Result[*models.Title] {
	if slug == "" {
		return domainerrors.Err[*models.Title](turso.NewError("titles.GetBySlug", domainerrors.CodeInvalidInput, "slug cannot be empty", nil))
	}

	t, err := client.Title.Query().Where(enttitle.SlugEQ(slug)).First(ctx)
	if err != nil {
		return domainerrors.Err[*models.Title](turso.MapEntError("titles.GetBySlug", "title", slug, err))
	}
	return domainerrors.Ok(toDomain(t))
}

// FindByIDOrSlug queries a title by Slug, ID, prefixed ID ("title-"+identifier), or Name (case-insensitive).
func FindByIDOrSlug(ctx context.Context, client *ent.Client, identifier string) domainerrors.Result[*models.Title] {
	if identifier == "" {
		return domainerrors.Err[*models.Title](turso.NewError("titles.FindByIDOrSlug", domainerrors.CodeInvalidInput, "title identifier cannot be empty", nil))
	}

	t, err := client.Title.Query().
		Where(enttitle.Or(
			enttitle.SlugEQ(identifier),
			enttitle.IDEQ(identifier),
			enttitle.IDEQ("title-"+identifier),
			enttitle.NameEqualFold(identifier),
		)).
		First(ctx)
	if err != nil {
		return domainerrors.Err[*models.Title](turso.MapEntError("titles.FindByIDOrSlug", "title", identifier, err))
	}
	return domainerrors.Ok(toDomain(t))
}

// ResolveHoursUntilPremiere calculates remaining hours until premiere for a title.
// If inputHours is > 0, it is returned directly (honoring explicit overrides).
// If identifier is empty or "GLOBAL", the baseline default urgency is returned without database lookup.
func ResolveHoursUntilPremiere(ctx context.Context, client *ent.Client, identifier string, inputHours float64) float64 {
	if inputHours > 0 {
		return inputHours
	}
	if identifier == "" || identifier == "GLOBAL" || client == nil {
		return DefaultPremiereUrgencyHours
	}

	res := FindByIDOrSlug(ctx, client, identifier)
	if res.IsOk() {
		t := res.Unwrap()
		if !t.PremiereDate.IsZero() {
			if h := time.Until(t.PremiereDate).Hours(); h > 0 {
				return h
			}
		}
	}

	return DefaultPremiereUrgencyHours
}
