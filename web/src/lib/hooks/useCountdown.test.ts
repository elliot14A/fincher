import { describe, expect, it } from 'bun:test'
import { calculateCountdown } from './useCountdown'

describe('calculateCountdown', () => {
  it('handles undefined or invalid dates safely', () => {
    const emptyResult = calculateCountdown(undefined)
    expect(emptyResult.scheduled).toBe(false)
    expect(emptyResult.timecode).toBe('-')
    expect(emptyResult.label).toBe('Not scheduled')

    const invalidResult = calculateCountdown('invalid-date')
    expect(invalidResult.scheduled).toBe(false)
    expect(invalidResult.timecode).toBe('-')
  })

  it('handles past/released dates', () => {
    const now = 1700000000000
    const pastDate = new Date(now - 10000).toISOString()
    const result = calculateCountdown(pastDate, now)

    expect(result.scheduled).toBe(true)
    expect(result.isPast).toBe(true)
    expect(result.timecode).toBe('00h 00m 00s')
    expect(result.label).toBe('Released')
  })

  it('calculates hours, minutes, seconds for future dates', () => {
    const now = 1700000000000
    // 2 hours, 15 minutes, 30 seconds
    const futureMs = now + (2 * 3600 + 15 * 60 + 30) * 1000
    const futureDate = new Date(futureMs).toISOString()

    const result = calculateCountdown(futureDate, now)
    expect(result.scheduled).toBe(true)
    expect(result.isPast).toBe(false)
    expect(result.hours).toBe(2)
    expect(result.minutes).toBe(15)
    expect(result.seconds).toBe(30)
    expect(result.timecode).toBe('2h 15m 30s')
    expect(result.label).toBe('2h 15m left')
  })

  it('formats multi-day range converting days into total hours', () => {
    const now = 1700000000000
    // 3 days (72h) + 4 hours = 76 hours, 12 minutes, 5 seconds
    const futureMs = now + (3 * 86400 + 4 * 3600 + 12 * 60 + 5) * 1000
    const futureDate = new Date(futureMs).toISOString()

    const result = calculateCountdown(futureDate, now)
    expect(result.scheduled).toBe(true)
    expect(result.hours).toBe(76)
    expect(result.minutes).toBe(12)
    expect(result.seconds).toBe(5)
    expect(result.label).toBe('76h 12m left')
    expect(result.timecode).toBe('76h 12m 05s')
  })
})
