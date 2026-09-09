package search

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entdelivery "github.com/elliot14A/fincher/internal/turso/ent/delivery"
	enttitle "github.com/elliot14A/fincher/internal/turso/ent/title"
	entvendor "github.com/elliot14A/fincher/internal/turso/ent/vendor"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

const defaultSearchLimit = 8

func metaString(metadata map[string]any, key string) string {
	if metadata == nil {
		return ""
	}
	if v, ok := metadata[key].(string); ok {
		return v
	}
	return ""
}

// Search performs a case-insensitive lookup across titles, vendors, and deliveries.
func Search(ctx context.Context, client *ent.Client, query string, limit int) domainerrors.Result[[]models.SearchResult] {
	if limit <= 0 {
		limit = defaultSearchLimit
	}

	out := make([]models.SearchResult, 0, limit*3)

	titles, err := client.Title.Query().
		Where(enttitle.Or(enttitle.NameContainsFold(query), enttitle.SlugContainsFold(query))).
		Limit(limit).
		All(ctx)
	if err != nil {
		return domainerrors.Err[[]models.SearchResult](turso.MapEntError("search.Search", "title", query, err))
	}
	for _, t := range titles {
		out = append(out, models.SearchResult{
			Kind:        models.SearchKindTitle,
			ID:          t.ID,
			DisplayName: t.Name,
			Identifier:  t.Slug,
			Subtitle:    metaString(t.Metadata, "genre"),
			ImageURL:    metaString(t.Metadata, "poster_url"),
		})
	}

	vendors, err := client.Vendor.Query().
		Where(entvendor.NameContainsFold(query)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return domainerrors.Err[[]models.SearchResult](turso.MapEntError("search.Search", "vendor", query, err))
	}
	for _, v := range vendors {
		out = append(out, models.SearchResult{
			Kind:        models.SearchKindVendor,
			ID:          v.ID,
			DisplayName: v.Name,
			Identifier:  v.ID,
		})
	}

	deliveries, err := client.Delivery.Query().
		Where(entdelivery.CountryContainsFold(query)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return domainerrors.Err[[]models.SearchResult](turso.MapEntError("search.Search", "delivery", query, err))
	}
	for _, d := range deliveries {
		out = append(out, models.SearchResult{
			Kind:        models.SearchKindDelivery,
			ID:          d.ID,
			DisplayName: d.Country,
			Identifier:  d.ID,
		})
	}

	return domainerrors.Ok(out)
}
