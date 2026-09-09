import { useQuery } from '@tanstack/react-query'
import type { RefObject } from 'preact'
import { useCallback, useEffect, useMemo, useRef, useState } from 'preact/hooks'
import type { ModelsSearchResult } from '#/lib/api'
import { searchQueryOptions } from '../queryOptions'

const MIN_QUERY_CHARS = 2
const SEARCH_LIMIT = 8
const MENTION_REGEX = /(?:^|\s)@(\S*)$/

export type UseObjectMentionsOptions = {
  value: string
  setValue: (next: string) => void
  textareaRef: RefObject<HTMLTextAreaElement>
  onSelectObject: (obj: ModelsSearchResult) => void
}

export type UseObjectMentionsResult = {
  isOpen: boolean
  query: string
  results: ModelsSearchResult[]
  isLoading: boolean
  activeIndex: number
  setActiveIndex: (i: number | ((prev: number) => number)) => void
  selectActive: () => void
  close: () => void
  handleKeyDown: (event: KeyboardEvent) => void
}

type Trigger = {
  start: number
  end: number
  query: string
}

export function useObjectMentions({
  value,
  setValue,
  textareaRef,
  onSelectObject,
}: UseObjectMentionsOptions): UseObjectMentionsResult {
  const [trigger, setTrigger] = useState<Trigger | null>(null)
  const [activeIndex, setActiveIndex] = useState(0)
  const [debouncedQuery, setDebouncedQuery] = useState('')
  const dismissedRef = useRef<{ start: number; query: string } | null>(null)

  useEffect(() => {
    const cursor = textareaRef.current?.selectionStart ?? value.length
    const before = value.slice(0, cursor)
    const match = MENTION_REGEX.exec(before)
    if (!match) {
      setTrigger(null)
      return
    }
    const queryPart = match[1] ?? ''
    const atIndex = match.index + (match[0].startsWith('@') ? 0 : 1)
    const start = atIndex
    const dismissed = dismissedRef.current
    if (dismissed && dismissed.start === start && dismissed.query === queryPart) {
      return
    }
    dismissedRef.current = null
    setTrigger({ start, end: cursor, query: queryPart })
  }, [value, textareaRef])

  useEffect(() => {
    if (!trigger) {
      setDebouncedQuery('')
      return
    }
    const handle = window.setTimeout(() => setDebouncedQuery(trigger.query), 150)
    return () => window.clearTimeout(handle)
  }, [trigger])

  const isOpen = trigger !== null
  const queryEnabled = isOpen && debouncedQuery.trim().length >= MIN_QUERY_CHARS

  const searchResults = useQuery({
    ...searchQueryOptions(debouncedQuery.trim(), SEARCH_LIMIT),
    enabled: queryEnabled,
  })

  const results = useMemo<ModelsSearchResult[]>(
    () => searchResults.data ?? [],
    [searchResults.data],
  )

  const effectiveActiveIndex = results.length === 0 ? 0 : Math.min(activeIndex, results.length - 1)

  const close = useCallback(() => {
    if (trigger) dismissedRef.current = { start: trigger.start, query: trigger.query }
    setTrigger(null)
  }, [trigger])

  const selectActive = useCallback(() => {
    if (!trigger) return
    const candidate = results[effectiveActiveIndex]
    if (!candidate) return
    const replacement = `@${candidate.identifier ?? candidate.display_name} `
    const next = value.slice(0, trigger.start) + replacement + value.slice(trigger.end)
    const newCursor = trigger.start + replacement.length
    setValue(next)
    onSelectObject(candidate)
    dismissedRef.current = null
    setTrigger(null)
    queueMicrotask(() => {
      const el = textareaRef.current
      if (el) {
        el.focus()
        el.setSelectionRange(newCursor, newCursor)
      }
    })
  }, [trigger, results, effectiveActiveIndex, value, setValue, onSelectObject, textareaRef])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (!isOpen) return
      switch (event.key) {
        case 'ArrowDown':
          event.preventDefault()
          setActiveIndex((prev) => (results.length ? (prev + 1) % results.length : 0))
          break
        case 'ArrowUp':
          event.preventDefault()
          setActiveIndex((prev) =>
            results.length ? (prev - 1 + results.length) % results.length : 0,
          )
          break
        case 'Enter':
          if (results.length > 0) {
            event.preventDefault()
            selectActive()
          }
          break
        case 'Escape':
          event.preventDefault()
          close()
          break
        case 'Tab':
          if (results.length > 0) {
            event.preventDefault()
            selectActive()
          }
          break
        default:
          break
      }
    },
    [isOpen, results.length, selectActive, close],
  )

  return {
    isOpen,
    query: trigger?.query ?? '',
    results,
    isLoading: searchResults.isFetching,
    activeIndex: effectiveActiveIndex,
    setActiveIndex,
    selectActive,
    close,
    handleKeyDown,
  }
}
