import { describe, expect, it } from 'bun:test'
import {
  calculateCountdown,
  getSimulatedDomainTime,
  resetSimulationAnchor,
  syncSimulationAnchorWithTitles,
} from './useCountdown'

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

  it('handles overdue statuses when past deadline', () => {
    const now = 1700000000000
    const pastDate = new Date(now - 10000).toISOString()

    const overdueResult = calculateCountdown(pastDate, now, 'OVERDUE')
    expect(overdueResult.label).toBe('Overdue')
    expect(overdueResult.timecode).toBe('00h 00m 00s')

    const holdResult = calculateCountdown(pastDate, now, 'HOLD')
    expect(holdResult.label).toBe('Overdue')

    const qcResult = calculateCountdown(pastDate, now, 'PROCESSING')
    expect(qcResult.label).toBe('In QC (Overdue)')
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

describe('simulation anchor & compressed time', () => {
  it('advances simulated domain time at 3600x real time from anchor', () => {
    const base = 1700000000000
    resetSimulationAnchor(base)

    // 10 real seconds later = 10 domain hours (36,000,000 ms)
    const simulated = getSimulatedDomainTime(base + 10000)
    expect(simulated).toBe(base + 10 * 3600 * 1000)
  })

  it('syncSimulationAnchorWithTitles re-anchors when newer titles arrive', () => {
    const oldBase = 1700000000000
    resetSimulationAnchor(oldBase)

    const newerCreated = new Date(oldBase + 60000).toISOString()
    syncSimulationAnchorWithTitles([{ created_at: newerCreated }])

    // Should now be anchored to the newer title
    const simulated = getSimulatedDomainTime(oldBase + 60000)
    expect(simulated).toBe(oldBase + 60000)
  })
})
