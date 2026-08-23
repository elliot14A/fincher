# Pagination Engineering Rules & Standards

This document establishes the universal pagination invariant and contract for both the Go backend and Preact frontend in Fincher, matching the unified standard from **Gaur** (`gaur-server` and `gaur-web`).

---

## 1. Backend Contract (`pkg/domain/models/pagination.go`)

All paginated list endpoints, repository queries, and analytical listing APIs must adhere to the standard `Pagination` input and `PaginationResult[T]` envelope.

### 1.1 Pagination Parameters
* **`page`** (`int`, query: `page`): 1-indexed current page number.
  - Default: `1`
  - Minimum: `1`
* **`limit`** (`int`, query: `limit`): Items per page.
  - Default: `10`
  - Bounds: `1 <= limit <= 100` (Max: `100`)
* **`sort_order`** (`string`, query: `sort_order`): Sort direction.
  - Allowed: `"asc"`, `"desc"`
  - Default: `"asc"`
* **`search`** (`string`, query: `search`, optional): Filter string.
  - Max length: `120` characters (trimmed)
* **`Offset()` calculation**:
  $$\text{offset} = (\text{page} - 1) \times \text{limit}$$

### 1.2 `PaginationResult[T]` Response Structure
Every paginated API response must return the structured envelope:
```json
{
  "items": [...],
  "total_items": 42,
  "page": 1,
  "limit": 10,
  "total_pages": 5,
  "has_next_page": true,
  "has_prev_page": false
}
```

* **`total_pages`**: $\lceil \text{total\_items} / \text{limit} \rceil$ (if $\text{total\_items} = 0 \implies 1$).
* **`has_next_page`**: $\text{page} < \text{total\_pages}$
* **`has_prev_page`**: $\text{page} > 1$

---

## 2. Frontend Contract (`web/src/components/ui/pagination/`)

The web operations console implements the exact same pagination contract and UI component established in Gaur.

### 2.1 TypeScript Contract
```ts
export type PaginationResult<T> = {
  items: T[]
  total_items: number
  page: number
  limit: number
  total_pages: number
  has_next_page: boolean
  has_prev_page: boolean
}

export type PaginationParams = {
  page?: number
  limit?: number
  sort_order?: 'asc' | 'desc'
  search?: string
}
```

### 2.2 `<PaginationControls />` Component
Co-located under `src/components/ui/pagination/` with Zero-Runtime Vanilla Extract styling:

* **Component Props**:
  ```ts
  export type PaginationControlsProps = {
    page: number
    totalPages: number
    hasNextPage: boolean
    hasPrevPage: boolean
    onPrevPage: () => void
    onNextPage: () => void
  }
  ```
* **UI Invariants**:
  - Full-width bottom bar with top subtle border (`vars.color.borderSubtle`).
  - Left: **Previous** button (disabled state when `!hasPrevPage`).
  - Center: **`Page {page} of {totalPages}`** in secondary/muted typography.
  - Right: **Next** button (disabled state when `!hasNextPage`).
  - Accessible focus rings, hover transitions, and cursor state (`cursor: not-allowed`, `opacity: 0.5` when disabled).

---

## 3. Engineering Invariants

1. **No Ad-Hoc Pagination**: Endpoints and UI grids must not invent custom pagination shapes (e.g. `pageSize`, `currentPage`, `skip`, `take`).
2. **Deterministic Query Keys**: TanStack Query keys for paginated lists must include the pagination object (e.g. `titlesKeys.list({ page, limit, status })`).
3. **Route Search Sync**: Page and filter parameters in the UI synchronize with route search parameters (`useSearch` in TanStack Router) where deep-linking is required.
