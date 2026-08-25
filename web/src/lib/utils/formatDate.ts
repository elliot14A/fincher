export function formatDate(
  dateStr: string | undefined | null,
  options?: Intl.DateTimeFormatOptions,
  fallback = 'Unscheduled',
): string {
  if (!dateStr) return fallback
  const d = new Date(dateStr)
  if (Number.isNaN(d.getTime())) return fallback
  return d.toLocaleDateString(
    undefined,
    options ?? {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    },
  )
}

export function formatDateTime(
  dateStr: string | undefined | null,
  fallback = 'Registered',
): string {
  return formatDate(
    dateStr,
    {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    },
    fallback,
  )
}
