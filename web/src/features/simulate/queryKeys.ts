export const simulateKeys = {
  all: ['simulate'] as const,
  titles: () => [...simulateKeys.all, 'titles'] as const,
  packages: (titleId: string) => [...simulateKeys.all, 'packages', titleId] as const,
}
