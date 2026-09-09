package search

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	apierrors "github.com/elliot14A/fincher/internal/api/errors"
	"github.com/elliot14A/fincher/internal/turso/ent"
	tursosearch "github.com/elliot14A/fincher/internal/turso/search"
)

// Search handles GET /api/search.
//
//	@Summary		Search taggable objects
//	@Description	Case-insensitive search across titles, vendors, and deliveries for @-mention tagging and reference resolution.
//	@Tags			search
//	@Produce		json
//	@Param			q		query		string	true	"Search query"
//	@Param			limit	query		int		false	"Max results per kind"
//	@Success		200		{array}		models.SearchResult
//	@Failure		400		{object}	errors.ErrorResponse
//	@Router			/search [get]
func Search(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := strings.TrimSpace(c.QueryParam("q"))
		if query == "" {
			return c.JSON(http.StatusOK, []any{})
		}

		limit := 0
		if raw := c.QueryParam("limit"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil {
				limit = parsed
			}
		}

		res := tursosearch.Search(c.Request().Context(), client, query, limit)
		if res.IsErr() {
			return apierrors.Respond(c, res.Error())
		}
		return c.JSON(http.StatusOK, res.Unwrap())
	}
}

// RegisterRoutes registers the search endpoint on the given router group.
func RegisterRoutes(g *echo.Group, client *ent.Client) {
	g.GET("", Search(client))
}
