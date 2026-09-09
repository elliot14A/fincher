import { useMutation, useQueryClient } from '@tanstack/react-query'
import { deleteChatById, patchChatById } from '#/lib/api'
import { chatKeys } from '../queryKeys'

export function useRenameSession() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, title }: { id: string; title: string }) => {
      const { error } = await patchChatById({ path: { id }, body: { title } })
      if (error) {
        throw new Error((error as { message?: string }).message || 'Failed to rename chat')
      }
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: chatKeys.lists() })
      queryClient.invalidateQueries({ queryKey: chatKeys.detail(variables.id) })
    },
  })
}

export function useDeleteSession() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const { error } = await deleteChatById({ path: { id } })
      if (error) {
        throw new Error((error as { message?: string }).message || 'Failed to delete chat')
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatKeys.lists() })
    },
  })
}
