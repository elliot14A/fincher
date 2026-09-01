import { useEffect, useState } from 'preact/hooks'

// Compressed time factor: 1 real second = 1 domain hour (3600x compression).
// IMPORTANT: This constant is tightly coupled to Go backend config.DefaultTimeScale (time.Second).
// If the backend timescale changes, this constant must be synchronized accordingly.
export const TIME_COMPRESSION_FACTOR = 3600

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
    return {
      label: 'Released',
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
): CountdownResult {
  const [mountTime] = useState(() => Date.now())
  const [currentDomainTime, setCurrentDomainTime] = useState(() => Date.now())

  useEffect(() => {
    if (!targetDateStr) return

    const target = new Date(targetDateStr).getTime()
    if (Number.isNaN(target)) return

    const timer = setInterval(() => {
      const elapsedRealMs = Date.now() - mountTime
      const simulatedDomainMs = mountTime + elapsedRealMs * TIME_COMPRESSION_FACTOR
      setCurrentDomainTime(simulatedDomainMs)

      if (target - simulatedDomainMs <= 0) {
        clearInterval(timer)
      }
    }, intervalMs)

    return () => clearInterval(timer)
  }, [targetDateStr, intervalMs, mountTime])

  return calculateCountdown(targetDateStr, currentDomainTime)
}
