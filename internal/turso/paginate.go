package turso

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// OrderBy selects the appropriate ascending or descending ordering option based on pagination.
func OrderBy[T any](p models.Pagination, ascOpt, descOpt T) T {
	if p.IsAsc() {
		return ascOpt
	}
	return descOpt
}

// Paginate executes a count query, fetches a bounded slice with limit and offset,
// maps raw database entities to domain models, and wraps them in a standard PaginationResult.
func Paginate[E any, D any](
	ctx context.Context,
	op string,
	p models.Pagination,
	countFn func(ctx context.Context) (int, error),
	fetchFn func(ctx context.Context, limit, offset int) ([]E, error),
	mapFn func([]E) []D,
) domainerrors.Result[models.PaginationResult[D]] {
	totalItems, err := countFn(ctx)
	if err != nil {
		return domainerrors.Err[models.PaginationResult[D]](NewError(op, domainerrors.CodeInternal, "failed to count records", err))
	}

	rawItems, err := fetchFn(ctx, p.Limit, p.Offset())
	if err != nil {
		return domainerrors.Err[models.PaginationResult[D]](NewError(op, domainerrors.CodeInternal, "failed to query records", err))
	}

	return domainerrors.Ok(models.NewPaginationResult(mapFn(rawItems), totalItems, p))
}
