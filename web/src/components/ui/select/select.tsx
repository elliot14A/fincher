import { Check, ChevronDown, Search } from 'lucide-preact'
import { useEffect, useRef, useState } from 'preact/hooks'
import {
  chevron,
  container,
  empty,
  list,
  option,
  optionActive,
  optionMain,
  optionMeta,
  panel,
  placeholder as placeholderClass,
  searchInput,
  searchRow,
  trigger,
  triggerLabel,
  triggerOpen,
} from './select.css'

export type SelectOption = {
  value: string
  label: string
  meta?: string
  keywords?: string
}

export type SelectProps = {
  options: SelectOption[]
  value: string
  onChange: (value: string) => void
  placeholder?: string
  searchable?: boolean
  ariaLabel?: string
}

export function Select({
  options,
  value,
  onChange,
  placeholder = 'Select…',
  searchable = true,
  ariaLabel = 'Select option',
}: SelectProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [query, setQuery] = useState('')
  const containerRef = useRef<HTMLDivElement>(null)
  const searchRef = useRef<HTMLInputElement>(null)

  const selected = options.find((o) => o.value === value)

  useEffect(() => {
    if (!isOpen) return

    const onClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false)
      }
    }
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setIsOpen(false)
    }
    document.addEventListener('mousedown', onClickOutside)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('mousedown', onClickOutside)
      document.removeEventListener('keydown', onKey)
    }
  }, [isOpen])

  useEffect(() => {
    if (isOpen) {
      setQuery('')
      searchRef.current?.focus()
    }
  }, [isOpen])

  const q = query.trim().toLowerCase()
  const filtered = q
    ? options.filter(
        (o) =>
          o.label.toLowerCase().includes(q) ||
          o.meta?.toLowerCase().includes(q) ||
          o.keywords?.toLowerCase().includes(q),
      )
    : options

  return (
    <div ref={containerRef} class={container}>
      <button
        type="button"
        class={isOpen ? `${trigger} ${triggerOpen}` : trigger}
        onClick={() => setIsOpen((p) => !p)}
        aria-haspopup="listbox"
        aria-expanded={isOpen}
        aria-label={ariaLabel}
      >
        <span class={selected ? triggerLabel : `${triggerLabel} ${placeholderClass}`}>
          {selected ? selected.label : placeholder}
        </span>
        <ChevronDown size={14} class={chevron} />
      </button>

      {isOpen && (
        <div class={panel} role="listbox">
          {searchable && (
            <div class={searchRow}>
              <Search size={13} />
              <input
                ref={searchRef}
                class={searchInput}
                value={query}
                placeholder="Search…"
                onInput={(e) => setQuery((e.target as HTMLInputElement).value)}
              />
            </div>
          )}
          <div class={list}>
            {filtered.length === 0 ? (
              <div class={empty}>No matches</div>
            ) : (
              filtered.map((o) => (
                <button
                  key={o.value}
                  type="button"
                  role="option"
                  aria-selected={o.value === value}
                  class={o.value === value ? `${option} ${optionActive}` : option}
                  onClick={() => {
                    onChange(o.value)
                    setIsOpen(false)
                  }}
                >
                  <span class={optionMain}>{o.label}</span>
                  {o.value === value ? (
                    <Check size={13} />
                  ) : o.meta ? (
                    <span class={optionMeta}>{o.meta}</span>
                  ) : null}
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
