package models

import "time"

// DependencyType indicates the coupling relation between packages.
type DependencyType string

const (
	DependencyAudioSync         DependencyType = "AUDIO_SYNC"
	DependencySubtitleAlignment DependencyType = "SUBTITLE_ALIGNMENT"
	DependencyMasterDerivation  DependencyType = "MASTER_DERIVATION"
)

// Dependency represents a directional link from parent to child package.
// Dependency edges are immutable (no partial update path); changes must be made by deleting and recreating the edge.
type Dependency struct {
	ID             string         `json:"id" validate:"required"`
	ParentID       string         `json:"parent_id" validate:"required"`
	ChildID        string         `json:"child_id" validate:"required"`
	DependencyType DependencyType `json:"dependency_type" validate:"required,oneof=AUDIO_SYNC SUBTITLE_ALIGNMENT MASTER_DERIVATION"`
	CreatedAt      time.Time      `json:"created_at"`
}

// Validate verifies dependency attributes.
func (d *Dependency) Validate() error {
	return validate.Struct(d)
}

// LineageNode represents a node in the resolved dependency tree.
type LineageNode struct {
	PackageID      string         `json:"package_id"`
	TitleID        string         `json:"title_id"`
	Component      ComponentType  `json:"component"`
	Language       string         `json:"language"`
	Status         PackageStatus  `json:"status"`
	DependencyType DependencyType `json:"dependency_type,omitempty"`
	Children       []*LineageNode `json:"children,omitempty"`
}

// LineageGraph represents the full dependency graph for a title.
type LineageGraph struct {
	TitleID string         `json:"title_id"`
	Roots   []*LineageNode `json:"roots"`
}
