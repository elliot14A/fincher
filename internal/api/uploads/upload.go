package uploads

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursouploads "github.com/elliot14A/fincher/internal/turso/uploads"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

var allowedMimeTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/webp": true,
	"image/gif":  true,
}

func generateUploadID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("img_%s", hex.EncodeToString(b))
}

// Upload handles POST /api/uploads.
//
//	@Summary		Upload an image asset
//	@Description	Uploads a raster image (PNG, JPEG, WebP, GIF <= 1MB) stored directly in SQLite as binary BLOB.
//	@Tags			uploads
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"Image binary file"
//	@Success		201		{object}	models.UploadResponse
//	@Failure		400		{object}	errors.DomainError
//	@Router			/uploads [post]
func Upload(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "missing 'file' field in multipart request",
			})
		}

		if fileHeader.Size > models.MaxUploadSizeBytes {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: fmt.Sprintf("file size (%d bytes) exceeds maximum allowed limit of %d bytes (1MB)", fileHeader.Size, models.MaxUploadSizeBytes),
			})
		}

		src, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "failed to read uploaded file",
			})
		}
		defer src.Close()

		// Read up to MaxUploadSizeBytes + 1 byte to strictly enforce boundary
		data, err := io.ReadAll(io.LimitReader(src, models.MaxUploadSizeBytes+1))
		if err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "failed to buffer uploaded file bytes",
			})
		}

		if int64(len(data)) > models.MaxUploadSizeBytes {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: fmt.Sprintf("file size (%d bytes) exceeds maximum allowed limit of %d bytes (1MB)", len(data), models.MaxUploadSizeBytes),
			})
		}

		// Sniff real content type from bytes (first 512 bytes)
		sniffLen := 512
		if len(data) < sniffLen {
			sniffLen = len(data)
		}
		detectedMime := http.DetectContentType(data[:sniffLen])

		// Normalize mime type (strip parameters like charset)
		baseMime := strings.TrimSpace(strings.Split(detectedMime, ";")[0])
		if !allowedMimeTypes[baseMime] {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: fmt.Sprintf("unsupported image format '%s'. Only PNG, JPEG, WebP, and GIF raster images are accepted", baseMime),
			})
		}

		uploadID := generateUploadID()
		u := &models.Upload{
			ID:        uploadID,
			Filename:  fileHeader.Filename,
			MimeType:  baseMime,
			Data:      data,
			SizeBytes: int64(len(data)),
			CreatedAt: time.Now().UTC(),
		}

		res := tursouploads.Create(c.Request().Context(), client, u)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		created := res.Unwrap()
		return c.JSON(http.StatusCreated, models.UploadResponse{
			ID:        created.ID,
			URL:       fmt.Sprintf("/api/uploads/%s", created.ID),
			Filename:  created.Filename,
			MimeType:  created.MimeType,
			SizeBytes: created.SizeBytes,
			CreatedAt: created.CreatedAt,
		})
	}
}
