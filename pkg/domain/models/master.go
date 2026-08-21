package models

import "time"

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
