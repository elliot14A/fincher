import { queryOptions } from '@tanstack/react-query'
import {
  getTitles,
  getTitlesById,
  type ModelsTitle,
  type ModelsTitlePaginationResult,
} from '#/lib/api'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'
import { titlesKeys } from './queryKeys'

export type TitlesFilters = {
  status?: string
  page?: number
  limit?: number
  sort_order?: string
  search?: string
}

export const titlesQueryOptions = (filters?: TitlesFilters | string) => {
  const normalized: TitlesFilters =
    typeof filters === 'string' ? { status: filters } : (filters ?? {})

  const page = normalized.page ?? DEFAULT_PAGE
  const limit = normalized.limit ?? DEFAULT_PAGE_LIMIT
  const sort_order = normalized.sort_order ?? DEFAULT_SORT_ORDER
  const status = normalized.status && normalized.status !== 'ALL' ? normalized.status : undefined
  const search = normalized.search

  return queryOptions({
    queryKey: titlesKeys.list({
      status: status ?? 'ALL',
      page,
      limit,
      sort_order,
      search: search ?? '',
    }),
    queryFn: async (): Promise<ModelsTitlePaginationResult> => {
      const { data, error } = await getTitles({
        query: {
          page,
          limit,
          sort_order,
          ...(status ? { status } : {}),
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

export const titleDetailQueryOptions = (id: string) =>
  queryOptions({
    queryKey: titlesKeys.detail(id),
    queryFn: async (): Promise<ModelsTitle | undefined> => {
      const { data, error } = await getTitlesById({
        path: { id },
      })
      if (error) {
        throw error
      }
      return data
    },
  })
