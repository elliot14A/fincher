export const deliveriesKeys = {
  all: ['deliveries'] as const,
  lists: () => [...deliveriesKeys.all, 'list'] as const,
  list: (filters?: { title_id?: string; country?: string; status?: string }) =>
    [...deliveriesKeys.lists(), filters] as const,
  details: () => [...deliveriesKeys.all, 'detail'] as const,
  detail: (id: string) => [...deliveriesKeys.details(), id] as const,
}
