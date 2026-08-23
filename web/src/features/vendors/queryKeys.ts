export const vendorsKeys = {
  all: ['vendors'] as const,
  lists: () => [...vendorsKeys.all, 'list'] as const,
  list: (specialty?: string) => [...vendorsKeys.lists(), { specialty }] as const,
  details: () => [...vendorsKeys.all, 'detail'] as const,
  detail: (id: string) => [...vendorsKeys.details(), id] as const,
}
