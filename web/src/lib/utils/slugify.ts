export function slugify(text: string, prefix = ''): string {
  const cleaned = text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')

  if (!cleaned) {
    return prefix ? `${prefix}-` : ''
  }

  return prefix ? `${prefix}-${cleaned}` : cleaned
}
