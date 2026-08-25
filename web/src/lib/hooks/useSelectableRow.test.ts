import { describe, expect, it } from 'bun:test'
import type { JSX } from 'preact'
import { useSelectableRow } from './useSelectableRow'

describe('useSelectableRow', () => {
  it('returns appropriate base class and a11y attributes when not selected', () => {
    let selected = false
    const { rowProps, isSelected } = useSelectableRow({
      isSelected: false,
      onSelect: () => {
        selected = true
      },
      baseClassName: 'row-base',
      activeClassName: 'row-active',
    })

    expect(isSelected).toBe(false)
    expect(rowProps.role).toBe('button')
    expect(rowProps.tabIndex).toBe(0)
    expect(rowProps.class).toBe('row-base')

    const dummyClickEvent = {} as JSX.TargetedMouseEvent<HTMLElement>
    rowProps.onClick(dummyClickEvent)
    expect(selected).toBe(true)
  })

  it('combines base and active class names when isSelected is true', () => {
    const { rowProps, isSelected } = useSelectableRow({
      isSelected: true,
      onSelect: () => {},
      baseClassName: 'row-base',
      activeClassName: 'row-active',
    })

    expect(isSelected).toBe(true)
    expect(rowProps.class).toBe('row-base row-active')
  })

  it('triggers onSelect when Enter or Space key is pressed', () => {
    let callCount = 0
    const onSelect = () => {
      callCount++
    }

    const { rowProps } = useSelectableRow({
      isSelected: false,
      onSelect,
    })

    let defaultPrevented = false
    const enterEvent = {
      key: 'Enter',
      preventDefault: () => {
        defaultPrevented = true
      },
    } as unknown as JSX.TargetedKeyboardEvent<HTMLElement>

    rowProps.onKeyDown(enterEvent)
    expect(callCount).toBe(1)
    expect(defaultPrevented).toBe(true)

    const spaceEvent = {
      key: ' ',
      preventDefault: () => {},
    } as unknown as JSX.TargetedKeyboardEvent<HTMLElement>
    rowProps.onKeyDown(spaceEvent)
    expect(callCount).toBe(2)

    const tabEvent = {
      key: 'Tab',
      preventDefault: () => {},
    } as unknown as JSX.TargetedKeyboardEvent<HTMLElement>
    rowProps.onKeyDown(tabEvent)
    expect(callCount).toBe(2)
  })
})
