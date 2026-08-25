import { describe, expect, it } from 'bun:test'
import { formatDate, formatDateTime } from './formatDate'

describe('formatDate', () => {
  it('returns fallback for undefined, null, or invalid dates', () => {
    expect(formatDate(undefined)).toBe('Unscheduled')
    expect(formatDate(null)).toBe('Unscheduled')
    expect(formatDate('not-a-date', undefined, 'Fallback')).toBe('Fallback')
  })

  it('formats valid ISO date string properly', () => {
    const formatted = formatDate('2026-08-24T12:00:00Z')
    expect(formatted).toContain('2026')
    expect(formatted).toContain('Aug')
  })
})

describe('formatDateTime', () => {
  it('returns fallback for undefined or invalid dates', () => {
    expect(formatDateTime(undefined)).toBe('Registered')
    expect(formatDateTime('invalid', 'CustomFallback')).toBe('CustomFallback')
  })

  it('formats valid ISO date-time string properly', () => {
    const formatted = formatDateTime('2026-08-24T12:30:00Z')
    expect(formatted).toContain('Aug')
  })
})
