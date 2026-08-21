package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// MaxMetadataSizeBytes caps the serialized metadata size to 64KB.
const MaxMetadataSizeBytes = 64 * 1024

// Base contains common entity fields: ID, metadata, and timestamps.
type Base struct {
	ID        string         `json:"id" validate:"required"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ValidateMetadata checks that the metadata does not exceed MaxMetadataSizeBytes when serialized.
func (b *Base) ValidateMetadata() error {
	if len(b.Metadata) == 0 {
		return nil
	}
	bytes, err := json.Marshal(b.Metadata)
	if err != nil {
		return fmt.Errorf("invalid metadata json: %w", err)
	}
	if len(bytes) > MaxMetadataSizeBytes {
		return fmt.Errorf("metadata size (%d bytes) exceeds limit of %d bytes", len(bytes), MaxMetadataSizeBytes)
	}
	return nil
}

// ValidateMetadataMap validates any standalone metadata map.
func ValidateMetadataMap(m map[string]any) error {
	if len(m) == 0 {
		return nil
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("invalid metadata json: %w", err)
	}
	if len(bytes) > MaxMetadataSizeBytes {
		return fmt.Errorf("metadata size (%d bytes) exceeds limit of %d bytes", len(bytes), MaxMetadataSizeBytes)
	}
	return nil
}
