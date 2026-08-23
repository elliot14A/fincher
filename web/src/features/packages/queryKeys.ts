export const packagesKeys = {
  all: ['packages'] as const,
  lists: () => [...packagesKeys.all, 'list'] as const,
  list: (filters?: {
    title_id?: string
    vendor_id?: string
    component?: string
    status?: string
  }) => [...packagesKeys.lists(), filters] as const,
  details: () => [...packagesKeys.all, 'detail'] as const,
  detail: (id: string) => [...packagesKeys.details(), id] as const,
  lineage: (titleId: string) => [...packagesKeys.all, 'lineage', titleId] as const,
}
