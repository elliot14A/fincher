package vendors

import (
	"context"
	"strings"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entvendor "github.com/elliot14A/fincher/internal/turso/ent/vendor"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// List fetches paginated vendors, optionally filtered by component and search term.
func List(ctx context.Context, client *ent.Client, componentFilter domainerrors.Option[string], p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Vendor]] {
	query := client.Vendor.Query()

	if p.Search != "" {
		query = query.Where(entvendor.NameContainsFold(p.Search))
	}

	query = query.Order(turso.OrderBy(p, ent.Asc(entvendor.FieldName), ent.Desc(entvendor.FieldName)))

	if componentFilter.IsNone() || componentFilter.Unwrap() == "" || strings.EqualFold(componentFilter.Unwrap(), "ALL") {
		return turso.Paginate(
			ctx,
			"vendors.List",
			p,
			query.Count,
			func(ctx context.Context, limit, offset int) ([]*ent.Vendor, error) {
				return query.Limit(limit).Offset(offset).All(ctx)
			},
			toDomainList,
		)
	}

	// In-Go filtering by component
	targetComp := strings.ToUpper(strings.TrimSpace(componentFilter.Unwrap()))
	all, err := query.All(ctx)
	if err != nil {
		return domainerrors.Err[models.PaginationResult[*models.Vendor]](turso.MapEntError("vendors.List", "vendor", "", err))
	}

	var filtered []*ent.Vendor
	for _, v := range all {
		for _, c := range v.Components {
			if strings.EqualFold(c, targetComp) {
				filtered = append(filtered, v)
				break
			}
		}
	}

	totalItems := len(filtered)
	page := p.Page
	if page < 1 {
		page = 1
	}
	limit := p.Limit
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var paged []*ent.Vendor
	if offset < totalItems {
		end := offset + limit
		if end > totalItems {
			end = totalItems
		}
		paged = filtered[offset:end]
	}

	totalPages := (totalItems + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return domainerrors.Ok(models.PaginationResult[*models.Vendor]{
		Items:       toDomainList(paged),
		TotalItems:  totalItems,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	})
}
