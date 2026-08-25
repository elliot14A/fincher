package vendors

import (
	"context"
	"strings"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// Delete removes a vendor by ID and cascades avatar upload cleanup if present.
func Delete(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	v, err := client.Vendor.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[bool](turso.MapEntError("vendors.Delete", "vendor", id, err))
	}

	avatarURL := ""
	if v.Metadata != nil {
		if u, ok := v.Metadata["avatar_url"].(string); ok {
			avatarURL = u
		}
	}

	if err := client.Vendor.DeleteOneID(id).Exec(ctx); err != nil {
		return domainerrors.Err[bool](turso.MapEntError("vendors.Delete", "vendor", id, err))
	}

	// Clean up associated internal upload blob if referenced
	if strings.HasPrefix(avatarURL, "/api/uploads/") {
		uploadID := strings.TrimPrefix(avatarURL, "/api/uploads/")
		_ = client.Upload.DeleteOneID(uploadID).Exec(ctx)
	}

	return domainerrors.Ok(true)
}
