import { useMutation, useQueryClient } from '@tanstack/react-query'
import { type ChatSendMessageResponse, type ChatTaggedObjectRef, postChat } from '#/lib/api'
import { chatKeys } from '../queryKeys'

export type SendMessageInput = {
  message: string
  sessionId?: string
  tagged?: ChatTaggedObjectRef[]
}

export function useSendMessage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      message,
      sessionId,
      tagged,
    }: SendMessageInput): Promise<ChatSendMessageResponse> => {
      const { data, error } = await postChat({
        body: {
          message,
          ...(sessionId ? { session_id: sessionId } : {}),
          ...(tagged && tagged.length > 0 ? { tagged_objects: tagged } : {}),
        },
      })
      if (error) {
        throw new Error((error as { message?: string }).message || 'Failed to reach the assistant')
      }
      if (!data) {
        throw new Error('The assistant returned an empty response')
      }
      return data
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: chatKeys.lists() })
      if (data.session_id) {
        queryClient.invalidateQueries({ queryKey: chatKeys.detail(data.session_id) })
      }
    },
  })
}
