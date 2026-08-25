import { useEffect, useState } from 'preact/hooks'

export interface CountdownResult {
  label: string
  timecode: string
  scheduled: boolean
  isPast: boolean
  diffMs: number
  days: number
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
      days: 0,
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
      days: 0,
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
      days: 0,
      hours: 0,
      minutes: 0,
      seconds: 0,
    }
  }

  const totalSeconds = Math.floor(diffMs / 1000)
  const days = Math.floor(totalSeconds / 86400)
  const hours = Math.floor((totalSeconds % 86400) / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60

  const pad = (n: number) => n.toString().padStart(2, '0')

  let label: string
  let timecode: string

  if (days > 0) {
    label = hours > 0 ? `${days}d ${hours}h left` : `${days}d left`
    timecode = `${days}d ${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`
  } else if (hours > 0) {
    label = `${hours}h ${minutes}m left`
    timecode = `${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`
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
    days,
    hours,
    minutes,
    seconds,
  }
}

export function useCountdown(
  targetDateStr: string | undefined,
  intervalMs: number = 1000,
): CountdownResult {
  const [now, setNow] = useState(() => Date.now())

  useEffect(() => {
    if (!targetDateStr) return

    const target = new Date(targetDateStr).getTime()
    if (Number.isNaN(target)) return

    const timer = setInterval(() => {
      const current = Date.now()
      setNow(current)
      if (target - current <= 0) {
        clearInterval(timer)
      }
    }, intervalMs)

    return () => clearInterval(timer)
  }, [targetDateStr, intervalMs])

  return calculateCountdown(targetDateStr, now)
}
