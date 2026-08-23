import { useMemo, useState } from 'preact/hooks'

export type UsePaginatedListOptions = {
  defaultPage?: number
  limit?: number
}

export function usePaginatedList<T>(items: T[], options: UsePaginatedListOptions = {}) {
  const limit = options.limit ?? 10
  const [page, setPage] = useState(options.defaultPage ?? 1)

  const totalPages = Math.max(1, Math.ceil(items.length / limit))
  const safePage = Math.min(Math.max(1, page), totalPages)

  const paginatedItems = useMemo(() => {
    const start = (safePage - 1) * limit
    return items.slice(start, start + limit)
  }, [items, safePage, limit])

  const hasNextPage = safePage < totalPages
  const hasPrevPage = safePage > 1

  return {
    page: safePage,
    limit,
    totalPages,
    paginatedItems,
    hasNextPage,
    hasPrevPage,
    setPage,
    onPrevPage: () => setPage((p) => Math.max(1, p - 1)),
    onNextPage: () => setPage((p) => Math.min(totalPages, p + 1)),
  }
}
