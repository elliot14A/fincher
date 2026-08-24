package dependencies

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entdependency "github.com/elliot14A/fincher/internal/turso/ent/dependency"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// ListFilter contains optional search filters for dependencies.
type ListFilter struct {
	ParentID domainerrors.Option[string]
	ChildID  domainerrors.Option[string]
}

// List fetches paginated dependencies matching filters.
func List(ctx context.Context, client *ent.Client, filter ListFilter, p models.Pagination) domainerrors.Result[models.PaginationResult[*models.Dependency]] {
	query := client.Dependency.Query()

	if filter.ParentID.IsSome() {
		query = query.Where(entdependency.ParentIDEQ(filter.ParentID.Unwrap()))
	}
	if filter.ChildID.IsSome() {
		query = query.Where(entdependency.ChildIDEQ(filter.ChildID.Unwrap()))
	}

	query = query.Order(turso.OrderBy(p, ent.Asc(entdependency.FieldCreatedAt), ent.Desc(entdependency.FieldCreatedAt)))

	return turso.Paginate(
		ctx,
		"dependencies.List",
		p,
		query.Count,
		func(ctx context.Context, limit, offset int) ([]*ent.Dependency, error) {
			return query.Limit(limit).Offset(offset).All(ctx)
		},
		toDomainList,
	)
}
