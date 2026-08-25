import type { JSX } from 'preact'

export interface UseSelectableRowOptions {
  isSelected: boolean
  onSelect: () => void
  baseClassName?: string
  activeClassName?: string
}

export interface UseSelectableRowResult {
  isSelected: boolean
  rowProps: {
    role: 'button'
    tabIndex: 0
    class: string
    onClick: (event: JSX.TargetedMouseEvent<HTMLElement>) => void
    onKeyDown: (event: JSX.TargetedKeyboardEvent<HTMLElement>) => void
  }
}

export function useSelectableRow({
  isSelected,
  onSelect,
  baseClassName = '',
  activeClassName = '',
}: UseSelectableRowOptions): UseSelectableRowResult {
  const combinedClass =
    isSelected && activeClassName ? `${baseClassName} ${activeClassName}` : baseClassName

  const handleClick = (_event: JSX.TargetedMouseEvent<HTMLElement>) => {
    onSelect()
  }

  const handleKeyDown = (event: JSX.TargetedKeyboardEvent<HTMLElement>) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      onSelect()
    }
  }

  return {
    isSelected,
    rowProps: {
      role: 'button',
      tabIndex: 0,
      class: combinedClass,
      onClick: handleClick,
      onKeyDown: handleKeyDown,
    },
  }
}
