package models

import (
	"fmt"
	"strings"
	"time"
)

// Master represents an immutable editorial cut or master video version of a Title.
// Master cuts are append-only; newer cuts supersede older cuts via the supersedes_version field.
type Master struct {
	ID                string    `json:"id" validate:"required"`
	TitleID           string    `json:"title_id" validate:"required"`
	Version           string    `json:"version" validate:"required"`
	SupersedesVersion string    `json:"supersedes_version,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

// Validate verifies master attributes.
func (m *Master) Validate() error {
	return validate.Struct(m)
}

// MasterIDFor generates a deterministic master entity ID from a title slug and version string.
func MasterIDFor(titleSlug, version string) string {
	if version == "" {
		version = "v01"
	}
	return fmt.Sprintf("mst-%s-%s", titleSlug, strings.ToLower(version))
}
