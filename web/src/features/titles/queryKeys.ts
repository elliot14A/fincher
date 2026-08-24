export const titlesKeys = {
  all: ['titles'] as const,
  lists: () => [...titlesKeys.all, 'list'] as const,
  list: (filters?: unknown) => [...titlesKeys.lists(), filters] as const,
  details: () => [...titlesKeys.all, 'detail'] as const,
  detail: (id: string) => [...titlesKeys.details(), id] as const,
}
