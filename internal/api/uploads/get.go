package uploads

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursouploads "github.com/elliot14A/fincher/internal/turso/uploads"
)

// Get handles GET /api/uploads/:id.
//
//	@Summary		Stream uploaded image
//	@Description	Streams the raw binary image BLOB with caching and security headers.
//	@Tags			uploads
//	@Produce		image/png
//	@Produce		image/jpeg
//	@Produce		image/webp
//	@Produce		image/gif
//	@Produce		octet-stream
//	@Param			id	path		string	true	"Upload ID"
//	@Success		200	{file}		binary
//	@Failure		404	{object}	errors.DomainError
//	@Router			/uploads/{id} [get]
func Get(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, apierrors.ErrorResponse{
				Code:    "INVALID_INPUT",
				Message: "missing upload id parameter",
			})
		}

		res := tursouploads.Get(c.Request().Context(), client, id)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}

		upload := res.Unwrap()

		// Set security and immutable cache headers
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("Content-Security-Policy", "default-src 'none'")
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", upload.SizeBytes))

		return c.Blob(http.StatusOK, upload.MimeType, upload.Data)
	}
}
