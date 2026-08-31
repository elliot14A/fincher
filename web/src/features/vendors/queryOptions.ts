import { queryOptions } from '@tanstack/react-query'
import {
  getVendors,
  getVendorsById,
  type ModelsVendor,
  type ModelsVendorPaginationResult,
} from '#/lib/api'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'
import { vendorsKeys } from './queryKeys'

export type VendorsFilters = {
  component?: string
  specialty?: string
  page?: number
  limit?: number
  sort_order?: string
  search?: string
}

export const vendorsQueryOptions = (filters?: VendorsFilters | string) => {
  const normalized: VendorsFilters =
    typeof filters === 'string' ? { component: filters } : (filters ?? {})

  const page = normalized.page ?? DEFAULT_PAGE
  const limit = normalized.limit ?? DEFAULT_PAGE_LIMIT
  const sort_order = normalized.sort_order ?? DEFAULT_SORT_ORDER
  const comp = normalized.component ?? normalized.specialty
  const component = comp && comp !== 'ALL' ? comp : undefined
  const search = normalized.search

  return queryOptions({
    queryKey: vendorsKeys.list({
      specialty: component ?? 'ALL',
      page,
      limit,
      sort_order,
      search: search ?? '',
    }),
    queryFn: async (): Promise<ModelsVendorPaginationResult> => {
      const { data, error } = await getVendors({
        query: {
          page,
          limit,
          sort_order,
          ...(component ? { component } : {}),
          ...(search ? { search } : {}),
        },
      })
      if (error) {
        throw error
      }
      return (
        data ?? {
          items: [],
          page,
          limit,
          total_items: 0,
          total_pages: 1,
          has_next_page: false,
          has_prev_page: false,
        }
      )
    },
  })
}

export const vendorDetailQueryOptions = (id: string) =>
  queryOptions({
    queryKey: vendorsKeys.detail(id),
    queryFn: async (): Promise<ModelsVendor | undefined> => {
      const { data, error } = await getVendorsById({
        path: { id },
      })
      if (error) {
        throw error
      }
      return data
    },
  })
