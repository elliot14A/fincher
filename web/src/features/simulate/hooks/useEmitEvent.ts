import { useMutation, useQueryClient } from '@tanstack/react-query'
import { runsKeys } from '#/features/runs/queryKeys'
import { type ModelsEvent, type ModelsEventBatchResponse, postEvents } from '#/lib/api'

export function useEmitEvent() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (events: ModelsEvent[]): Promise<ModelsEventBatchResponse> => {
      const { data, error } = await postEvents({ body: events })
      if (error) {
        throw new Error((error as { message?: string }).message || 'Failed to emit event')
      }
      if (!data) {
        throw new Error('The ingestion API returned an empty response')
      }
      return data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: runsKeys.lists() })
    },
  })
}
