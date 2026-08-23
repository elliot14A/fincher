package models

const (
	DefaultPage      = 1
	DefaultLimit     = 10
	MaxLimit         = 100
	DefaultSortOrder = "asc"
	MaxSearchLength  = 120
)

type Pagination struct {
	Page      int               `json:"page" query:"page"`
	Limit     int               `json:"limit" query:"limit"`
	SortOrder string            `json:"sort_order" query:"sort_order"`
	Search    string            `json:"search,omitempty" query:"search"`
	Params    map[string]string `json:"-"`
}

func NewPagination(page, limit int, sortOrder, search string) Pagination {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	} else if limit > MaxLimit {
		limit = MaxLimit
	}
	if sortOrder != "desc" {
		sortOrder = DefaultSortOrder
	}
	if len(search) > MaxSearchLength {
		search = search[:MaxSearchLength]
	}
	return Pagination{
		Page:      page,
		Limit:     limit,
		SortOrder: sortOrder,
		Search:    search,
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func (p Pagination) TotalPages(totalItems int) int {
	if totalItems <= 0 {
		return 1
	}
	return (totalItems + p.Limit - 1) / p.Limit
}

func (p Pagination) HasNextPage(totalItems int) bool {
	return p.Page < p.TotalPages(totalItems)
}

func (p Pagination) HasPrevPage() bool {
	return p.Page > 1
}

type PaginationResult[T any] struct {
	Items       []T  `json:"items"`
	TotalItems  int  `json:"total_items"`
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalPages  int  `json:"total_pages"`
	HasNextPage bool `json:"has_next_page"`
	HasPrevPage bool `json:"has_prev_page"`
}

func NewPaginationResult[T any](items []T, totalItems int, p Pagination) PaginationResult[T] {
	if items == nil {
		items = []T{}
	}
	return PaginationResult[T]{
		Items:       items,
		TotalItems:  totalItems,
		Page:        p.Page,
		Limit:       p.Limit,
		TotalPages:  p.TotalPages(totalItems),
		HasNextPage: p.HasNextPage(totalItems),
		HasPrevPage: p.HasPrevPage(),
	}
}
