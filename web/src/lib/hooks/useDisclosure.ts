import { useCallback, useState } from 'preact/hooks'

export type UseDisclosureReturn = {
  isOpen: boolean
  open: () => void
  close: () => void
  toggle: () => void
  setIsOpen: (isOpen: boolean) => void
}

export function useDisclosure(defaultIsOpen = false): UseDisclosureReturn {
  const [isOpen, setIsOpen] = useState(defaultIsOpen)

  const open = useCallback(() => setIsOpen(true), [])
  const close = useCallback(() => setIsOpen(false), [])
  const toggle = useCallback(() => setIsOpen((prev) => !prev), [])

  return {
    isOpen,
    open,
    close,
    toggle,
    setIsOpen,
  }
}
