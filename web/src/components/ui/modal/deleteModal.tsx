import { AlertTriangle, Trash2 } from 'lucide-preact'
import { Button } from '#/components/ui/button'
import {
  entityBadge,
  promptText,
  warningBox,
  warningContent,
  warningIcon,
  warningText as warningTextClass,
  warningTitle,
} from './deleteModal.css'
import { Modal } from './modal'

export type DeleteModalProps = {
  isOpen: boolean
  onClose: () => void
  onConfirm: () => void
  title?: string
  entityType: string
  entityName?: string
  entityId: string
  warningMessage?: string
  isDeleting?: boolean
}

export function DeleteModal({
  isOpen,
  onClose,
  onConfirm,
  title = 'Confirm Deletion',
  entityType,
  entityName,
  entityId,
  warningMessage = 'This action cannot be undone. Any dependent workflows or references may be affected.',
  isDeleting = false,
}: DeleteModalProps) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      description={`Permanently remove this ${entityType.toLowerCase()} from the operational database.`}
      footer={
        <>
          <Button variant="ghost" size="sm" onClick={onClose} disabled={isDeleting}>
            Cancel
          </Button>
          <Button variant="danger" size="sm" onClick={onConfirm} disabled={isDeleting}>
            <Trash2 size={14} />
            <span>{isDeleting ? 'Deleting...' : `Delete ${entityType}`}</span>
          </Button>
        </>
      }
    >
      <div class={warningBox}>
        <AlertTriangle size={18} class={warningIcon} />
        <div class={warningContent}>
          <div class={warningTitle}>Irreversible Action</div>
          <div class={warningTextClass}>{warningMessage}</div>
        </div>
      </div>

      <div class={promptText}>
        Are you sure you want to delete{' '}
        <strong>{entityName ? `"${entityName}"` : entityType}</strong>{' '}
        <span class={entityBadge}>({entityId})</span>?
      </div>
    </Modal>
  )
}
