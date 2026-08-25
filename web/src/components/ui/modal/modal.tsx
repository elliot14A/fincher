import { X } from 'lucide-preact'
import type { ComponentChildren } from 'preact'
import { useEffect } from 'preact/hooks'
import {
  backdrop,
  body,
  closeButton,
  description as descriptionClass,
  footer as footerClass,
  header,
  modalContainer,
  title as titleClass,
  titleGroup,
} from './modal.css'

export type ModalProps = {
  isOpen: boolean
  onClose: () => void
  title: string
  description?: string
  children: ComponentChildren
  footer?: ComponentChildren
}

export function Modal({ isOpen, onClose, title, description, children, footer }: ModalProps) {
  useEffect(() => {
    if (!isOpen) return

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose()
      }
    }

    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [isOpen, onClose])

  if (!isOpen) return null

  return (
    // biome-ignore lint/a11y/noStaticElementInteractions: backdrop click-to-close overlay
    <div
      role="presentation"
      class={backdrop}
      onClick={(e) => {
        if (e.target === e.currentTarget) {
          onClose()
        }
      }}
    >
      <div role="dialog" aria-modal="true" aria-labelledby="modal-title" class={modalContainer}>
        <div class={header}>
          <div class={titleGroup}>
            <h2 id="modal-title" class={titleClass}>
              {title}
            </h2>
            {description ? <p class={descriptionClass}>{description}</p> : null}
          </div>
          <button type="button" class={closeButton} onClick={onClose} aria-label="Close dialog">
            <X size={16} />
          </button>
        </div>

        <div class={body}>{children}</div>

        {footer ? <div class={footerClass}>{footer}</div> : null}
      </div>
    </div>
  )
}
