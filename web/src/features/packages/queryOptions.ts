import { queryOptions } from '@tanstack/react-query'
import {
  getDependenciesGraphByTitleId,
  getPackages,
  getPackagesById,
  type ModelsLineageGraph,
  type ModelsPackage,
} from '#/lib/api'
import { packagesKeys } from './queryKeys'

export type PackagesFilters = {
  title_id?: string
  vendor_id?: string
  component?: 'VIDEO' | 'AUDIO' | 'SUBTITLE' | 'METADATA'
  status?: 'PENDING' | 'VALID' | 'INVALIDATED' | 'RE_QC_PENDING'
}

export const packagesQueryOptions = (filters?: PackagesFilters) =>
  queryOptions({
    queryKey: packagesKeys.list(filters),
    queryFn: async (): Promise<ModelsPackage[]> => {
      const { data, error } = await getPackages({
        query: filters,
      })
      if (error) {
        throw error
      }
      return data ?? []
    },
  })

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
