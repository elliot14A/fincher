export const vendorsKeys = {
  all: ['vendors'] as const,
  lists: () => [...vendorsKeys.all, 'list'] as const,
  list: (filters?: unknown) => [...vendorsKeys.lists(), filters] as const,
  details: () => [...vendorsKeys.all, 'detail'] as const,
  detail: (id: string) => [...vendorsKeys.details(), id] as const,
}
