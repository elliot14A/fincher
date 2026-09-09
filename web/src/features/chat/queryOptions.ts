import { queryOptions } from '@tanstack/react-query'
import {
  getChat,
  getChatById,
  getSearch,
  type ModelsChatSession,
  type ModelsSearchResult,
} from '#/lib/api'
import { chatKeys } from './queryKeys'

export const chatListQueryOptions = () =>
  queryOptions({
    queryKey: chatKeys.list(),
    queryFn: async (): Promise<ModelsChatSession[]> => {
      const { data, error } = await getChat()
      if (error) {
        throw error
      }
      return data ?? []
    },
  })

export const chatDetailQueryOptions = (id: string) =>
  queryOptions({
    queryKey: chatKeys.detail(id),
    queryFn: async (): Promise<ModelsChatSession | undefined> => {
      const { data, error } = await getChatById({ path: { id } })
      if (error) {
        throw error
      }
      return data
    },
  })

export const searchQueryOptions = (query: string, limit = 8) =>
  queryOptions({
    queryKey: [...chatKeys.all, 'search', query, limit] as const,
    queryFn: async (): Promise<ModelsSearchResult[]> => {
      const { data, error } = await getSearch({ query: { q: query, limit } })
      if (error) {
        throw error
      }
      return data ?? []
    },
    staleTime: 30_000,
  })
