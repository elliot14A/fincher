package dependencies

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entdependency "github.com/elliot14A/fincher/internal/turso/ent/dependency"
	entmediapackage "github.com/elliot14A/fincher/internal/turso/ent/mediapackage"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// GetLineageGraph constructs the resolved dependency tree for a title.
func GetLineageGraph(ctx context.Context, client *ent.Client, titleID string) domainerrors.Result[*models.LineageGraph] {
	packages, err := client.MediaPackage.Query().
		Where(entmediapackage.TitleIDEQ(titleID)).
		All(ctx)
	if err != nil {
		return domainerrors.Err[*models.LineageGraph](turso.NewError("dependencies.GetLineageGraph", domainerrors.CodeInternal, "failed to query packages", err))
	}

	if len(packages) == 0 {
		return domainerrors.Ok(&models.LineageGraph{
			TitleID: titleID,
			Roots:   []*models.LineageNode{},
		})
	}

	pkgIDs := make([]string, len(packages))
	pkgMap := make(map[string]*ent.MediaPackage, len(packages))
	for i, p := range packages {
		pkgIDs[i] = p.ID
		pkgMap[p.ID] = p
	}

	allDeps, err := client.Dependency.Query().
		Where(entdependency.Or(
			entdependency.ParentIDIn(pkgIDs...),
			entdependency.ChildIDIn(pkgIDs...),
		)).
		All(ctx)
	if err != nil {
		return domainerrors.Err[*models.LineageGraph](turso.NewError("dependencies.GetLineageGraph", domainerrors.CodeInternal, "failed to query dependencies", err))
	}

	parentToChildren := make(map[string][]struct {
		ChildID string
		DepType models.DependencyType
	})
	childHasParent := make(map[string]bool)

	for _, dep := range allDeps {
		if _, ok := pkgMap[dep.ParentID]; ok {
			if _, okChild := pkgMap[dep.ChildID]; okChild {
				parentToChildren[dep.ParentID] = append(parentToChildren[dep.ParentID], struct {
					ChildID string
					DepType models.DependencyType
				}{
					ChildID: dep.ChildID,
					DepType: models.DependencyType(dep.DependencyType),
				})
				childHasParent[dep.ChildID] = true
			}
		}
	}

	var buildNode func(pkgID string, depType models.DependencyType, visited map[string]bool) *models.LineageNode
	buildNode = func(pkgID string, depType models.DependencyType, visited map[string]bool) *models.LineageNode {
		if visited[pkgID] {
			return nil
		}
		visited[pkgID] = true

		p, ok := pkgMap[pkgID]
		if !ok {
			return nil
		}

		node := &models.LineageNode{
			PackageID:      p.ID,
			TitleID:        p.TitleID,
			Component:      models.ComponentType(p.Component),
			Language:       p.Language,
			Status:         models.PackageStatus(p.Status),
			DependencyType: depType,
			Children:       make([]*models.LineageNode, 0),
		}

		for _, edge := range parentToChildren[pkgID] {
			childVisited := make(map[string]bool, len(visited))
			for k, v := range visited {
				childVisited[k] = v
			}
			if childNode := buildNode(edge.ChildID, edge.DepType, childVisited); childNode != nil {
				node.Children = append(node.Children, childNode)
			}
		}
		return node
	}

	roots := make([]*models.LineageNode, 0)
	for _, p := range packages {
		if !childHasParent[p.ID] {
			if rootNode := buildNode(p.ID, "", make(map[string]bool)); rootNode != nil {
				roots = append(roots, rootNode)
			}
		}
	}

	return domainerrors.Ok(&models.LineageGraph{
		TitleID: titleID,
		Roots:   roots,
	})
}
