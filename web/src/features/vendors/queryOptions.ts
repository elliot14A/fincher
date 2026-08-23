import { queryOptions } from '@tanstack/react-query'
import { getVendors, getVendorsById, type ModelsVendor } from '#/lib/api'
import { vendorsKeys } from './queryKeys'

export const vendorsQueryOptions = (specialty?: string) =>
  queryOptions({
    queryKey: vendorsKeys.list(specialty),
    queryFn: async (): Promise<ModelsVendor[]> => {
      const { data, error } = await getVendors({
        query: specialty && specialty !== 'ALL' ? { specialty } : undefined,
      })
      if (error) {
        throw error
      }
      return data ?? []
    },
  })

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
