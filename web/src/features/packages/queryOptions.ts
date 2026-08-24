import { queryOptions } from '@tanstack/react-query'
import {
  getDependenciesGraphByTitleId,
  getPackages,
  getPackagesById,
  type ModelsLineageGraph,
  type ModelsPackage,
  type ModelsPackagePaginationResult,
} from '#/lib/api'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'
import { packagesKeys } from './queryKeys'

export type PackagesFilters = {
  title_id?: string
  vendor_id?: string
  component?: 'VIDEO' | 'AUDIO' | 'SUBTITLE' | 'METADATA'
  status?: 'PENDING' | 'VALID' | 'INVALIDATED' | 'RE_QC_PENDING'
  page?: number
  limit?: number
  sort_order?: string
  search?: string
}

export const packagesQueryOptions = (filters?: PackagesFilters) => {
  const page = filters?.page ?? DEFAULT_PAGE
  const limit = filters?.limit ?? DEFAULT_PAGE_LIMIT
  const sort_order = filters?.sort_order ?? DEFAULT_SORT_ORDER

  return queryOptions({
    queryKey: packagesKeys.list({
      page,
      limit,
      sort_order,
      ...filters,
    }),
    queryFn: async (): Promise<ModelsPackagePaginationResult> => {
      const { data, error } = await getPackages({
        query: {
          page,
          limit,
          sort_order,
          ...filters,
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

export const packageDetailQueryOptions = (id: string) =>
  queryOptions({
    queryKey: packagesKeys.detail(id),
    queryFn: async (): Promise<ModelsPackage | undefined> => {
      const { data, error } = await getPackagesById({
        path: { id },
      })
      if (error) {
        throw error
      }
      return data
    },
  })

export const lineageGraphQueryOptions = (titleId: string) =>
  queryOptions({
    queryKey: packagesKeys.lineage(titleId),
    queryFn: async (): Promise<ModelsLineageGraph | undefined> => {
      const { data, error } = await getDependenciesGraphByTitleId({
        path: { title_id: titleId },
      })
      if (error) {
        throw error
      }
      return data
    },
    enabled: Boolean(titleId),
  })
