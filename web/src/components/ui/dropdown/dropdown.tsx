import { MoreVertical } from 'lucide-preact'
import type { ComponentType, JSX } from 'preact'
import { useEffect, useRef, useState } from 'preact/hooks'
import {
  menuContainer,
  menuDivider,
  menuDropdown,
  menuItem,
  menuItemDanger,
  menuItemIcon,
  menuTrigger,
  menuTriggerActive,
} from './dropdown.css'

export interface DropdownActionItem {
  type?: 'action'
  key: string
  label: string
  // biome-ignore lint/suspicious/noExplicitAny: supports Lucide icons and custom svg components
  icon?: ComponentType<any>
  danger?: boolean
  disabled?: boolean
  onClick: () => void
}

export interface DropdownDividerItem {
  type: 'divider'
  key: string
}

export type DropdownMenuItem = DropdownActionItem | DropdownDividerItem

export interface ActionMenuProps {
  items: DropdownMenuItem[]
  ariaLabel?: string
}

export function ActionMenu({ items, ariaLabel = 'Row actions' }: ActionMenuProps) {
  const [isOpen, setIsOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!isOpen) return

    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen])

  const handleToggle = (e: JSX.TargetedMouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()
    setIsOpen((prev) => !prev)
  }

  return (
    <div ref={containerRef} class={menuContainer}>
      <button
        type="button"
        class={isOpen ? `${menuTrigger} ${menuTriggerActive}` : menuTrigger}
        onClick={handleToggle}
        aria-label={ariaLabel}
        aria-expanded={isOpen}
      >
        <MoreVertical size={16} />
      </button>

      {isOpen ? (
        <div class={menuDropdown} role="menu">
          {items.map((item) => {
            if (item.type === 'divider') {
              return <div key={item.key} class={menuDivider} />
            }

            const Icon = item.icon
            const itemClass = item.danger ? `${menuItem} ${menuItemDanger}` : menuItem

            return (
              <button
                key={item.key}
                type="button"
                role="menuitem"
                disabled={item.disabled}
                class={itemClass}
                onClick={(e) => {
                  e.stopPropagation()
                  if (item.disabled) return
                  setIsOpen(false)
                  item.onClick()
                }}
              >
                {Icon ? <Icon size={14} class={menuItemIcon} /> : null}
                <span>{item.label}</span>
              </button>
            )
          })}
        </div>
      ) : null}
    </div>
  )
}
