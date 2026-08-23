import { queryOptions } from '@tanstack/react-query'
import { getTitles, getTitlesById, type ModelsTitle } from '#/lib/api'
import { titlesKeys } from './queryKeys'

export const titlesQueryOptions = (status?: string) =>
  queryOptions({
    queryKey: titlesKeys.list(status),
    queryFn: async (): Promise<ModelsTitle[]> => {
      const { data, error } = await getTitles({
        query: status && status !== 'ALL' ? { status } : undefined,
      })
      if (error) {
        throw error
      }
      return data ?? []
    },
  })

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
