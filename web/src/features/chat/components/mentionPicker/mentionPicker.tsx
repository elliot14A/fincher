import { Film, LayoutGrid, Users } from 'lucide-preact'
import { useMemo } from 'preact/hooks'
import type { ModelsSearchResult, ModelsSearchResultKind } from '#/lib/api'
import {
  empty,
  groupLabel,
  iconBox,
  option,
  optionActive,
  optionBody,
  optionMeta,
  optionMetaActive,
  optionName,
  popover,
  poster,
} from './mentionPicker.css'

type MentionPickerProps = {
  isOpen: boolean
  query: string
  results: ModelsSearchResult[]
  isLoading: boolean
  activeIndex: number
  onActiveIndexChange: (i: number) => void
  onSelect: (i: number) => void
}

const KIND_META: Record<ModelsSearchResultKind, { group: string; icon: typeof Film }> = {
  title: { group: 'Titles', icon: Film },
  vendor: { group: 'Vendors', icon: Users },
  delivery: { group: 'Deliveries', icon: LayoutGrid },
}

export function MentionPicker({
  isOpen,
  query,
  results,
  isLoading,
  activeIndex,
  onActiveIndexChange,
  onSelect,
}: MentionPickerProps) {
  const groupOrder = useMemo<ModelsSearchResultKind[]>(() => {
    const seen = new Set<ModelsSearchResultKind>()
    const order: ModelsSearchResultKind[] = []
    for (const r of results) {
      if (!r.kind || seen.has(r.kind)) continue
      seen.add(r.kind)
      order.push(r.kind)
    }
    return order
  }, [results])

  if (!isOpen) return null

  const trimmed = query.trim()
  const tooShort = trimmed.length > 0 && trimmed.length < 2

  if (results.length === 0) {
    return (
      <div class={popover}>
        <div class={empty}>
          {isLoading
            ? 'Searching...'
            : trimmed.length === 0
              ? 'Type to reference a title, vendor, or delivery'
              : tooShort
                ? 'Type at least 2 characters'
                : `No matches for "${trimmed}"`}
        </div>
      </div>
    )
  }

  return (
    <div class={popover}>
      {groupOrder.map((kind) => {
        const meta = KIND_META[kind]
        const Icon = meta.icon
        const items = results.map((r, idx) => ({ r, idx })).filter(({ r }) => r.kind === kind)
        return (
          <div key={kind}>
            <div class={groupLabel}>{meta.group}</div>
            {items.map(({ r, idx }) => {
              const isActive = idx === activeIndex
              return (
                <button
                  key={`${r.kind}-${r.id}`}
                  type="button"
                  class={isActive ? `${option} ${optionActive}` : option}
                  onMouseDown={(e) => {
                    e.preventDefault()
                    onSelect(idx)
                  }}
                  onMouseEnter={() => onActiveIndexChange(idx)}
                >
                  {r.image_url ? (
                    <img class={poster} src={r.image_url} alt="" loading="lazy" />
                  ) : (
                    <span class={iconBox}>
                      <Icon size={14} />
                    </span>
                  )}
                  <span class={optionBody}>
                    <span class={optionName}>{r.display_name}</span>
                    <span class={isActive ? `${optionMeta} ${optionMetaActive}` : optionMeta}>
                      {r.subtitle ? `${r.subtitle} · ` : ''}@{r.identifier}
                    </span>
                  </span>
                </button>
              )
            })}
          </div>
        )
      })}
    </div>
  )
}
