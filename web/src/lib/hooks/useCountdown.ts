import { useEffect, useState } from 'preact/hooks'

export const TIME_COMPRESSION_FACTOR = 3600
export const SIMULATION_ANCHOR_KEY = 'fincher_sim_anchor_v1'
export const MAX_ANCHOR_AGE_MS = 24 * 60 * 60 * 1000

export interface SimulationAnchor {
  realStart: number
  domainStart: number
}

let memoryAnchor: SimulationAnchor | null = null

/**
 * Retrieves the persistent simulation anchor from browser storage or memory fallback.
 * If none exists or if it has expired, initializes a new anchor.
 */
export function getSimulationAnchor(): SimulationAnchor {
  if (typeof window !== 'undefined' && window.localStorage) {
    try {
      const raw = localStorage.getItem(SIMULATION_ANCHOR_KEY)
      if (raw) {
        const parsed = JSON.parse(raw)
        if (
          typeof parsed.realStart === 'number' &&
          typeof parsed.domainStart === 'number' &&
          !Number.isNaN(parsed.realStart) &&
          !Number.isNaN(parsed.domainStart) &&
          Math.abs(Date.now() - parsed.realStart) < MAX_ANCHOR_AGE_MS
        ) {
          return parsed
        }
      }
    } catch {}
  }

  if (memoryAnchor) {
    return memoryAnchor
  }

  const now = Date.now()
  const anchor: SimulationAnchor = { realStart: now, domainStart: now }
  memoryAnchor = anchor
  if (typeof window !== 'undefined' && window.localStorage) {
    try {
      localStorage.setItem(SIMULATION_ANCHOR_KEY, JSON.stringify(anchor))
    } catch {}
  }
  return anchor
}

export const ANCHOR_CHANGED_EVENT = 'fincher:sim_anchor_changed'

/**
 * Resets the persistent simulation anchor (e.g. when database is re-seeded).
 */
export function resetSimulationAnchor(baseTimeMs: number = Date.now()): SimulationAnchor {
  const anchor: SimulationAnchor = { realStart: baseTimeMs, domainStart: baseTimeMs }
  memoryAnchor = anchor
  if (typeof window !== 'undefined') {
    if (window.localStorage) {
      try {
        localStorage.setItem(SIMULATION_ANCHOR_KEY, JSON.stringify(anchor))
      } catch {}
    }
    window.dispatchEvent(new CustomEvent(ANCHOR_CHANGED_EVENT, { detail: anchor }))
  }
  return anchor
}

/**
 * Syncs the simulation anchor with fetched titles.
 * If titles were seeded or created after the current anchor was set,
 * or if the anchor has diverged from the current dataset's creation time,
 * automatically resets the anchor to the latest title creation time.
 */
export function syncSimulationAnchorWithTitles(titles: Array<{ created_at?: string }>): void {
  if (!titles || titles.length === 0) return

  let latestCreatedAtMs = 0
  for (const t of titles) {
    if (t.created_at) {
      const ms = new Date(t.created_at).getTime()
      if (!Number.isNaN(ms) && ms > latestCreatedAtMs) {
        latestCreatedAtMs = ms
      }
    }
  }

  if (latestCreatedAtMs > 0) {
    const currentAnchor = getSimulationAnchor()
    // If titles were created at a different time than current anchor (e.g. re-seed or newly created),
    // re-anchor to the latest title creation time.
    if (Math.abs(latestCreatedAtMs - currentAnchor.realStart) > 5000) {
      resetSimulationAnchor(latestCreatedAtMs)
    }
  }
}

/**
 * Calculates current simulated domain time based on real elapsed time from anchor.
 */
export function getSimulatedDomainTime(now: number = Date.now()): number {
  const anchor = getSimulationAnchor()
  const elapsedRealMs = now - anchor.realStart
  return anchor.domainStart + elapsedRealMs * TIME_COMPRESSION_FACTOR
}

export interface CountdownResult {
  label: string
  timecode: string
  scheduled: boolean
  isPast: boolean
  diffMs: number
  hours: number
  minutes: number
  seconds: number
}

export function calculateCountdown(
  targetDateStr: string | undefined,
  nowMs: number = Date.now(),
  status?: string,
): CountdownResult {
  if (!targetDateStr) {
    return {
      label: 'Not scheduled',
      timecode: '-',
      scheduled: false,
      isPast: false,
      diffMs: 0,
      hours: 0,
      minutes: 0,
      seconds: 0,
    }
  }

  const target = new Date(targetDateStr).getTime()
  if (Number.isNaN(target)) {
    return {
      label: 'Not scheduled',
      timecode: '-',
      scheduled: false,
      isPast: false,
      diffMs: 0,
      hours: 0,
      minutes: 0,
      seconds: 0,
    }
  }

  const diffMs = target - nowMs
  if (diffMs <= 0) {
    let label = 'Released'
    if (status === 'HOLD' || status === 'OVERDUE') {
      label = 'Overdue'
    } else if (status === 'PROCESSING') {
      label = 'In QC (Overdue)'
    } else if (status === 'DRAFT') {
      label = 'Draft (Overdue)'
    }

    return {
      label,
      timecode: '00h 00m 00s',
      scheduled: true,
      isPast: true,
      diffMs,
      hours: 0,
      minutes: 0,
      seconds: 0,
    }
  }

  const totalSeconds = Math.floor(diffMs / 1000)
  const totalHours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60

  const pad = (n: number) => n.toString().padStart(2, '0')

  let label: string
  let timecode: string

  if (totalHours > 0) {
    label = minutes > 0 ? `${totalHours}h ${minutes}m left` : `${totalHours}h left`
    timecode = `${totalHours}h ${pad(minutes)}m ${pad(seconds)}s`
  } else if (minutes > 0) {
    label = `${minutes}m ${seconds}s left`
    timecode = `${pad(minutes)}m ${pad(seconds)}s`
  } else {
    label = `${seconds}s left`
    timecode = `${pad(seconds)}s`
  }

  return {
    label,
    timecode,
    scheduled: true,
    isPast: false,
    diffMs,
    hours: totalHours,
    minutes,
    seconds,
  }
}

export function useCountdown(
  targetDateStr: string | undefined,
  intervalMs: number = 100,
  status?: string,
): CountdownResult {
  const [currentDomainTime, setCurrentDomainTime] = useState(() => getSimulatedDomainTime())

  useEffect(() => {
    if (!targetDateStr) return

    const target = new Date(targetDateStr).getTime()
    if (Number.isNaN(target)) return

    let timer: ReturnType<typeof setInterval> | null = null

    const update = () => {
      const simTime = getSimulatedDomainTime()
      setCurrentDomainTime(simTime)
      return simTime
    }

    const startTimer = () => {
      if (timer) {
        clearInterval(timer)
        timer = null
      }
      const current = update()
      if (target - current <= 0) {
        return
      }
      timer = setInterval(() => {
        const simTime = update()
        if (target - simTime <= 0) {
          if (timer) clearInterval(timer)
          timer = null
        }
      }, intervalMs)
    }

    startTimer()

    const onAnchorChanged = () => {
      startTimer()
    }

    if (typeof window !== 'undefined') {
      window.addEventListener(ANCHOR_CHANGED_EVENT, onAnchorChanged)
    }

    return () => {
      if (timer) clearInterval(timer)
      if (typeof window !== 'undefined') {
        window.removeEventListener(ANCHOR_CHANGED_EVENT, onAnchorChanged)
      }
    }
  }, [targetDateStr, intervalMs])

  return calculateCountdown(targetDateStr, currentDomainTime, status)
}
