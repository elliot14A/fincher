import { useMemo, useState } from 'preact/hooks'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT } from '#/lib/constants'

export type PaginationCalculation = {
  page: number
  limit: number
  totalPages: number
  startIndex: number
  endIndex: number
  hasNextPage: boolean
  hasPrevPage: boolean
}

export function calculatePagination(
  totalItems: number,
  page = DEFAULT_PAGE,
  limit = DEFAULT_PAGE_LIMIT,
): PaginationCalculation {
  const safeLimit = Math.max(1, limit)
  const totalPages = Math.max(1, Math.ceil(totalItems / safeLimit))
  const safePage = Math.min(Math.max(1, page), totalPages)
  const startIndex = (safePage - 1) * safeLimit
  const endIndex = Math.min(totalItems, startIndex + safeLimit)

  return {
    page: safePage,
    limit: safeLimit,
    totalPages,
    startIndex,
    endIndex,
    hasNextPage: safePage < totalPages,
    hasPrevPage: safePage > 1,
  }
}

export function paginateArray<T>(
  items: T[],
  page = DEFAULT_PAGE,
  limit = DEFAULT_PAGE_LIMIT,
): { items: T[]; calculation: PaginationCalculation } {
  const calculation = calculatePagination(items.length, page, limit)
  return {
    items: items.slice(calculation.startIndex, calculation.endIndex),
    calculation,
  }
}

export type UsePaginatedListOptions = {
  defaultPage?: number
  limit?: number
}

export function usePaginatedList<T>(items: T[], options: UsePaginatedListOptions = {}) {
  const limit = options.limit ?? DEFAULT_PAGE_LIMIT
  const [page, setPage] = useState(options.defaultPage ?? DEFAULT_PAGE)

  const calculation = calculatePagination(items.length, page, limit)

  const paginatedItems = useMemo(() => {
    return items.slice(calculation.startIndex, calculation.endIndex)
  }, [items, calculation.startIndex, calculation.endIndex])

  return {
    page: calculation.page,
    limit: calculation.limit,
    totalPages: calculation.totalPages,
    paginatedItems,
    hasNextPage: calculation.hasNextPage,
    hasPrevPage: calculation.hasPrevPage,
    setPage,
    onPrevPage: () => setPage((p) => Math.max(1, p - 1)),
    onNextPage: () => setPage((p) => Math.min(calculation.totalPages, p + 1)),
  }
}
