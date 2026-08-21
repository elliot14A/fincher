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

// List fetches dependencies matching filters.
func List(ctx context.Context, client *ent.Client, filter ListFilter) domainerrors.Result[[]*models.Dependency] {
	query := client.Dependency.Query().Order(ent.Asc(entdependency.FieldCreatedAt))

	if filter.ParentID.IsSome() {
		query = query.Where(entdependency.ParentIDEQ(filter.ParentID.Unwrap()))
	}
	if filter.ChildID.IsSome() {
		query = query.Where(entdependency.ChildIDEQ(filter.ChildID.Unwrap()))
	}

	deps, err := query.All(ctx)
	if err != nil {
		return domainerrors.Err[[]*models.Dependency](turso.NewError("dependencies.List", domainerrors.CodeInternal, "failed to query dependencies", err))
	}

	return domainerrors.Ok(toDomainList(deps))
}
