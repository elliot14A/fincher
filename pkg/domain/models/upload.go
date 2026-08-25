package models

import (
	"fmt"
	"time"
)

// MaxUploadSizeBytes defines the strict 1MB size limit for file uploads.
const MaxUploadSizeBytes = 1 * 1024 * 1024

// Upload represents an uploaded binary file stored in SQLite.
type Upload struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	MimeType  string    `json:"mime_type"`
	Data      []byte    `json:"-"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadResponse represents the public metadata returned after uploading a file.
type UploadResponse struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Filename  string    `json:"filename"`
	MimeType  string    `json:"mime_type"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate checks the upload's integrity and size constraints.
func (u *Upload) Validate() error {
	if u.ID == "" {
		return fmt.Errorf("upload id is required")
	}
	if u.Filename == "" {
		return fmt.Errorf("filename is required")
	}
	if u.MimeType == "" {
		return fmt.Errorf("mime_type is required")
	}
	if len(u.Data) == 0 {
		return fmt.Errorf("upload data cannot be empty")
	}
	if u.SizeBytes <= 0 {
		u.SizeBytes = int64(len(u.Data))
	}
	if u.SizeBytes > MaxUploadSizeBytes {
		return fmt.Errorf("file size (%d bytes) exceeds maximum limit of %d bytes", u.SizeBytes, MaxUploadSizeBytes)
	}
	return nil
}
