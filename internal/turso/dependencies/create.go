package dependencies

import (
	"context"
	"fmt"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entdependency "github.com/elliot14A/fincher/internal/turso/ent/dependency"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// wouldCreateCycle checks if adding parent -> child creates a cycle in the DAG.
func wouldCreateCycle(ctx context.Context, client *ent.Client, parentID, childID string) (bool, error) {
	if parentID == childID {
		return true, nil
	}

	visited := make(map[string]bool)
	queue := []string{childID}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == parentID {
			return true, nil
		}

		if visited[curr] {
			continue
		}
		visited[curr] = true

		deps, err := client.Dependency.Query().
			Where(entdependency.ParentIDEQ(curr)).
			All(ctx)
		if err != nil {
			return false, err
		}

		for _, dep := range deps {
			if !visited[dep.ChildID] {
				queue = append(queue, dep.ChildID)
			}
		}
	}

	return false, nil
}

// Create inserts a dependency edge after validating same-title ownership and DAG cycle prevention.
func Create(ctx context.Context, client *ent.Client, d *models.Dependency) domainerrors.Result[*models.Dependency] {
	if err := d.Validate(); err != nil {
		return domainerrors.Err[*models.Dependency](turso.NewError("dependencies.Create", domainerrors.CodeInvalidInput, "invalid dependency data", err))
	}

	// Verify both parent and child exist and belong to the same Title
	parentPkg, err := client.MediaPackage.Get(ctx, d.ParentID)
	if err != nil {
		return domainerrors.Err[*models.Dependency](turso.MapEntError("dependencies.Create", "parent_package", d.ParentID, err))
	}

	childPkg, err := client.MediaPackage.Get(ctx, d.ChildID)
	if err != nil {
		return domainerrors.Err[*models.Dependency](turso.MapEntError("dependencies.Create", "child_package", d.ChildID, err))
	}

	if parentPkg.TitleID != childPkg.TitleID {
		return domainerrors.Err[*models.Dependency](turso.NewError(
			"dependencies.Create",
			domainerrors.CodeInvalidInput,
			fmt.Sprintf("cross-title dependency rejected: parent '%s' (title '%s') and child '%s' (title '%s') belong to different titles",
				d.ParentID, parentPkg.TitleID, d.ChildID, childPkg.TitleID),
			nil,
		))
	}

	isCycle, err := wouldCreateCycle(ctx, client, d.ParentID, d.ChildID)
	if err != nil {
		return domainerrors.Err[*models.Dependency](turso.NewError("dependencies.Create", domainerrors.CodeInternal, "failed cycle detection check", err))
	}
	if isCycle {
		return domainerrors.Err[*models.Dependency](turso.NewError(
			"dependencies.Create",
			domainerrors.CodeConflict,
			fmt.Sprintf("circular dependency detected: %s cannot depend on %s", d.ChildID, d.ParentID),
			nil,
		))
	}

	created, err := client.Dependency.Create().
		SetID(d.ID).
		SetParentID(d.ParentID).
		SetChildID(d.ChildID).
		SetDependencyType(entdependency.DependencyType(d.DependencyType)).
		Save(ctx)

	if err != nil {
		return domainerrors.Err[*models.Dependency](turso.MapEntError("dependencies.Create", "dependency", d.ID, err))
	}

	return domainerrors.Ok(toDomain(created))
}
