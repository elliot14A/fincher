import { queryOptions } from '@tanstack/react-query'
import { getDeliveries, getDeliveriesById, type ModelsDelivery } from '#/lib/api'
import { deliveriesKeys } from './queryKeys'

export type DeliveriesFilters = {
  title_id?: string
  country?: string
  status?: string
}

export const deliveriesQueryOptions = (filters?: DeliveriesFilters) =>
  queryOptions({
    queryKey: deliveriesKeys.list(filters),
    queryFn: async (): Promise<ModelsDelivery[]> => {
      const { data, error } = await getDeliveries({
        query: filters,
      })
      if (error) {
        throw error
      }
      return data ?? []
    },
  })

export const deliveryDetailQueryOptions = (id: string) =>
  queryOptions({
    queryKey: deliveriesKeys.detail(id),
    queryFn: async (): Promise<ModelsDelivery | undefined> => {
      const { data, error } = await getDeliveriesById({
        path: { id },
      })
      if (error) {
        throw error
      }
      return data
    },
  })
