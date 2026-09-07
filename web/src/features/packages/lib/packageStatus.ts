import { FileCode, Film, Headphones, Subtitles } from 'lucide-preact'
import type { BadgeProps } from '#/components/ui/badge'
import type { ModelsPackage } from '#/lib/api'

export function mapPackageStatus(status?: ModelsPackage['status']): {
  label: string
  variant: BadgeProps['variant']
} {
  switch (status) {
    case 'VALID':
      return { label: 'Passed QC', variant: 'success' }
    case 'INVALIDATED':
      return { label: 'QC Failed', variant: 'danger' }
    case 'RE_QC_PENDING':
      return { label: 'Re-QC Pending', variant: 'warning' }
    case 'PENDING':
      return { label: 'Pending QC', variant: 'warning' }
    default:
      return { label: status ?? 'Pending', variant: 'neutral' }
  }
}

export function getComponentIcon(component?: ModelsPackage['component'] | string) {
  switch (component) {
    case 'VIDEO':
      return Film
    case 'AUDIO':
      return Headphones
    case 'SUBTITLE':
      return Subtitles
    default:
      return FileCode
  }
}
