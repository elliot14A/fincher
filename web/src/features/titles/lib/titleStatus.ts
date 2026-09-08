import type { BadgeProps } from '#/components/ui/badge'
import type { ModelsTitle } from '#/lib/api'

type TitleStatus = ModelsTitle['overall_status'] | undefined

export function mapTitleStatus(status: TitleStatus): {
  label: string
  variant: BadgeProps['variant']
} {
  switch (status) {
    case 'DRAFT':
      return { label: 'Draft', variant: 'neutral' }
    case 'HOLD':
      return { label: 'Hold', variant: 'danger' }
    case 'OVERDUE':
      return { label: 'Overdue', variant: 'danger' }
    case 'AT_RISK':
      return { label: 'At Risk', variant: 'warning' }
    case 'PROCESSING':
      return { label: 'In QC', variant: 'warning' }
    case 'ON_TRACK':
      return { label: 'Ready', variant: 'success' }
    case 'SHIPPED':
      return { label: 'Shipped', variant: 'success' }
    default:
      return { label: status ?? 'Draft', variant: 'neutral' }
  }
}

export function getTitleStatusNote(status: TitleStatus): string {
  switch (status) {
    case 'DRAFT':
      return 'Awaiting Master QC'
    case 'HOLD':
      return 'Delivery hold active'
    case 'OVERDUE':
      return 'Premiere window breached'
    case 'AT_RISK':
      return 'High risk package drift'
    case 'PROCESSING':
      return 'QC inspection in progress'
    case 'ON_TRACK':
      return 'All packages verified'
    case 'SHIPPED':
      return 'Worldwide delivery completed'
    default:
      return 'Awaiting cut confirmation'
  }
}

export function getQcGating(status: TitleStatus): {
  isInQC: boolean
  canRunQC: boolean
} {
  return {
    isInQC: status === 'PROCESSING',
    canRunQC: status === 'DRAFT' || status === 'AT_RISK',
  }
}
