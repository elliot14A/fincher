import { queryOptions } from '@tanstack/react-query'
import { getPackages, getTitles, type ModelsPackage, type ModelsTitle } from '#/lib/api'
import { simulateKeys } from './queryKeys'

export const simulateTitlesQueryOptions = () =>
  queryOptions({
    queryKey: simulateKeys.titles(),
    queryFn: async (): Promise<ModelsTitle[]> => {
      const { data, error } = await getTitles({ query: { limit: 100, sort_order: 'asc' } })
      if (error) {
        throw error
      }
      return data?.items ?? []
    },
  })

export const simulatePackagesQueryOptions = (titleId: string) =>
  queryOptions({
    queryKey: simulateKeys.packages(titleId),
    enabled: Boolean(titleId),
    queryFn: async (): Promise<ModelsPackage[]> => {
      const { data, error } = await getPackages({ query: { title_id: titleId, limit: 100 } })
      if (error) {
        throw error
      }
      return data?.items ?? []
    },
  })
