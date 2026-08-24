import { queryOptions } from '@tanstack/react-query'
import {
  getDeliveries,
  getDeliveriesById,
  type ModelsDelivery,
  type ModelsDeliveryPaginationResult,
} from '#/lib/api'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'
import { deliveriesKeys } from './queryKeys'

export type DeliveriesFilters = {
  title_id?: string
  country?: string
  status?: string
  page?: number
  limit?: number
  sort_order?: string
  search?: string
}

export const deliveriesQueryOptions = (filters?: DeliveriesFilters) => {
  const page = filters?.page ?? DEFAULT_PAGE
  const limit = filters?.limit ?? DEFAULT_PAGE_LIMIT
  const sort_order = filters?.sort_order ?? DEFAULT_SORT_ORDER

  return queryOptions({
    queryKey: deliveriesKeys.list({
      page,
      limit,
      sort_order,
      ...filters,
    }),
    queryFn: async (): Promise<ModelsDeliveryPaginationResult> => {
      const { data, error } = await getDeliveries({
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
