package models_test

import (
	"testing"

	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestPagination_DefaultsAndBounds(t *testing.T) {
	p := models.NewPagination(0, 0, "invalid", "")
	if p.Page != models.DefaultPage {
		t.Errorf("expected page %d, got %d", models.DefaultPage, p.Page)
	}
	if p.Limit != models.DefaultLimit {
		t.Errorf("expected limit %d, got %d", models.DefaultLimit, p.Limit)
	}
	if p.SortOrder != models.DefaultSortOrder {
		t.Errorf("expected sort_order %s, got %s", models.DefaultSortOrder, p.SortOrder)
	}
	if p.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", p.Offset())
	}

	pOver := models.NewPagination(3, 500, "desc", "query")
	if pOver.Page != 3 {
		t.Errorf("expected page 3, got %d", pOver.Page)
	}
	if pOver.Limit != models.MaxLimit {
		t.Errorf("expected limit %d, got %d", models.MaxLimit, pOver.Limit)
	}
	if pOver.SortOrder != "desc" {
		t.Errorf("expected sort_order desc, got %s", pOver.SortOrder)
	}
	if pOver.Offset() != 200 {
		t.Errorf("expected offset 200, got %d", pOver.Offset())
	}
}

func TestPagination_Calculations(t *testing.T) {
	p := models.NewPagination(1, 10, "asc", "")
	if p.TotalPages(0) != 1 {
		t.Errorf("expected total pages 1 for 0 items, got %d", p.TotalPages(0))
	}
	if p.HasNextPage(0) {
		t.Errorf("expected has_next_page false for 0 items")
	}
	if p.HasPrevPage() {
		t.Errorf("expected has_prev_page false for page 1")
	}

	if p.TotalPages(25) != 3 {
		t.Errorf("expected total pages 3 for 25 items, got %d", p.TotalPages(25))
	}
	if !p.HasNextPage(25) {
		t.Errorf("expected has_next_page true on page 1 of 3")
	}
	if p.HasPrevPage() {
		t.Errorf("expected has_prev_page false on page 1")
	}

	p2 := models.NewPagination(2, 10, "asc", "")
	if !p2.HasNextPage(25) {
		t.Errorf("expected has_next_page true on page 2 of 3")
	}
	if !p2.HasPrevPage() {
		t.Errorf("expected has_prev_page true on page 2")
	}

	p3 := models.NewPagination(3, 10, "asc", "")
	if p3.HasNextPage(25) {
		t.Errorf("expected has_next_page false on page 3 of 3")
	}
	if !p3.HasPrevPage() {
		t.Errorf("expected has_prev_page true on page 3")
	}
}

func TestPaginationResult_Construct(t *testing.T) {
	items := []string{"item1", "item2"}
	p := models.NewPagination(1, 10, "asc", "")
	res := models.NewPaginationResult(items, 25, p)

	if len(res.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(res.Items))
	}
	if res.TotalItems != 25 {
		t.Errorf("expected total_items 25, got %d", res.TotalItems)
	}
	if res.Page != 1 {
		t.Errorf("expected page 1, got %d", res.Page)
	}
	if res.Limit != 10 {
		t.Errorf("expected limit 10, got %d", res.Limit)
	}
	if res.TotalPages != 3 {
		t.Errorf("expected total_pages 3, got %d", res.TotalPages)
	}
	if !res.HasNextPage {
		t.Errorf("expected has_next_page true")
	}
	if res.HasPrevPage {
		t.Errorf("expected has_prev_page false")
	}
}
