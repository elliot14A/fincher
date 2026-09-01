import { queryOptions } from '@tanstack/react-query'
import { getRuns, getRunsById, type ModelsRun, type ModelsRunPaginationResult } from '#/lib/api'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT } from '#/lib/constants'
import { runsKeys } from './queryKeys'

export type RunsFilters = {
  trigger?: string
  status?: string
  title_slug?: string
  page?: number
  limit?: number
}

export const runsQueryOptions = (filters?: RunsFilters | string) => {
  const normalized: RunsFilters =
    typeof filters === 'string' ? { trigger: filters } : (filters ?? {})

  const page = normalized.page ?? DEFAULT_PAGE
  const limit = normalized.limit ?? DEFAULT_PAGE_LIMIT
  const trigger =
    normalized.trigger && normalized.trigger !== 'ALL'
      ? normalized.trigger.toLowerCase()
      : undefined
  const status =
    normalized.status && normalized.status !== 'ALL' ? normalized.status.toUpperCase() : undefined
  const title_slug = normalized.title_slug

  return queryOptions({
    queryKey: runsKeys.list({
      trigger: trigger ?? 'ALL',
      status: status ?? 'ALL',
      title_slug: title_slug ?? '',
      page,
      limit,
    }),
    queryFn: async (): Promise<ModelsRunPaginationResult> => {
      const { data, error } = await getRuns({
        query: {
          page,
          limit,
          ...(trigger ? { trigger } : {}),
          ...(status ? { status } : {}),
          ...(title_slug ? { title_slug } : {}),
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
    refetchInterval: (query) => {
      // Auto-poll actively if any run in current view is RUNNING
      const hasRunning = query.state.data?.items?.some(
        (r) => r.status === 'RUNNING' || r.status === 'PENDING',
      )
      return hasRunning ? 3000 : 10000
    },
  })
}

export const runDetailQueryOptions = (id: string) =>
  queryOptions({
    queryKey: runsKeys.detail(id),
    queryFn: async (): Promise<ModelsRun | undefined> => {
      const { data, error } = await getRunsById({
        path: { id },
      })
      if (error) {
        throw error
      }
      return data
    },
    refetchInterval: (query) => {
      const isRunning =
        query.state.data?.status === 'RUNNING' || query.state.data?.status === 'PENDING'
      return isRunning ? 2000 : false
    },
  })
