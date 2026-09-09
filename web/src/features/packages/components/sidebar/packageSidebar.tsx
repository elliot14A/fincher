import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import {
  ArrowRight,
  Building2,
  ChevronLeft,
  FileCode,
  Film,
  Layers,
  ShieldAlert,
  X,
} from 'lucide-preact'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { getComponentIcon, mapPackageStatus } from '#/features/packages/lib'
import { packageDetailQueryOptions } from '#/features/packages/queryOptions'
import { titleDetailQueryOptions } from '#/features/titles/queryOptions'
import { vendorDetailQueryOptions } from '#/features/vendors/queryOptions'
import { formatDateTime } from '#/lib/utils/formatDate'
import {
  backBtn,
  closeBtn,
  componentAvatar,
  entityCard,
  entityCardChevron,
  entityCardInfo,
  entityCardLeft,
  entityCardMeta,
  entityCardName,
  failedNotice,
  failedNoticeIcon,
  footer,
  header,
  headerActions,
  headerMain,
  metadataCard,
  metadataKey,
  metadataRow,
  metadataVal,
  packageIdClass,
  packageInfo,
  packageMeta,
  scrollArea,
  section,
  sectionTitle,
  sidebarPanel,
  statCard,
  statGrid,
  statLabel,
  statValue,
} from './packageSidebar.css'

export type PackageSidebarProps = {
  packageId: string | null
  onClose: () => void
  onBack?: () => void
  backLabel?: string
}

export function PackageSidebar({
  packageId,
  onClose,
  onBack,
  backLabel = 'Back',
}: PackageSidebarProps) {
  const navigate = useNavigate()

  const packageQuery = useQuery({
    ...packageDetailQueryOptions(packageId ?? ''),
    enabled: Boolean(packageId),
  })
  const pkg = packageQuery.data

  const titleQuery = useQuery({
    ...titleDetailQueryOptions(pkg?.title_id ?? ''),
    enabled: Boolean(pkg?.title_id),
  })
  const title = titleQuery.data

  const vendorQuery = useQuery({
    ...vendorDetailQueryOptions(pkg?.vendor_id ?? ''),
    enabled: Boolean(pkg?.vendor_id),
  })
  const vendor = vendorQuery.data

  if (!packageId) return null

  const statusInfo = mapPackageStatus(pkg?.status)
  const CompIcon = getComponentIcon(pkg?.component ?? 'VIDEO')
  const meta = (pkg?.metadata as Record<string, unknown> | undefined) ?? {}
  const isFailed = pkg?.status === 'INVALIDATED' || pkg?.status === 'RE_QC_PENDING'

  return (
    <aside class={sidebarPanel({ size: 'standard' })} aria-label="Package Details">
      <header class={header}>
        <div class={headerMain}>
          <div class={componentAvatar}>
            <CompIcon size={20} />
          </div>
          <div class={packageInfo}>
            <h2 class={packageIdClass}>{pkg?.id ?? packageId}</h2>
            <div class={packageMeta}>
              <span>{pkg?.component ?? 'PACKAGE'}</span>
              <span>•</span>
              <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
            </div>
          </div>
        </div>
        <div class={headerActions}>
          {onBack && (
            <button type="button" class={backBtn} onClick={onBack} aria-label={backLabel}>
              <ChevronLeft size={13} />
              <span>{backLabel}</span>
            </button>
          )}
          <button type="button" class={closeBtn} onClick={onClose} aria-label="Close sidebar">
            <X size={16} />
          </button>
        </div>
      </header>

      <div class={scrollArea}>
        {isFailed && (
          <div class={failedNotice}>
            <ShieldAlert size={16} class={failedNoticeIcon} />
            <div>
              <strong>Quality Check Exception:</strong> Package failed validation rules or was
              invalidated due to upstream master drift. Requires vendor re-delivery.
            </div>
          </div>
        )}

        <div class={section}>
          <h3 class={sectionTitle}>
            <Layers size={12} /> Asset Specifications
          </h3>
          <div class={statGrid({ columns: 3 })}>
            <div class={statCard}>
              <span class={statValue}>{pkg?.language || 'en-US'}</span>
              <span class={statLabel}>Language</span>
            </div>
            <div class={statCard}>
              <span class={statValue}>{pkg?.market || 'GLOBAL'}</span>
              <span class={statLabel}>Market</span>
            </div>
            <div class={statCard}>
              <span class={statValue}>{pkg?.version || 'V01'}</span>
              <span class={statLabel}>Version</span>
            </div>
          </div>
        </div>

        {title && (
          <div class={section}>
            <h3 class={sectionTitle}>
              <Film size={12} /> Parent Media Title
            </h3>
            <button type="button" class={entityCard} onClick={() => navigate({ to: '/titles' })}>
              <div class={entityCardLeft}>
                <div class={entityCardInfo}>
                  <span class={entityCardName}>{title.name}</span>
                  <span class={entityCardMeta}>
                    Master: {title.current_master_version || 'V01'} • Status: {title.overall_status}
                  </span>
                </div>
              </div>
              <ArrowRight size={14} class={entityCardChevron} />
            </button>
          </div>
        )}

        {vendor && (
          <div class={section}>
            <h3 class={sectionTitle}>
              <Building2 size={12} /> Assigned Vendor Facility
            </h3>
            <button type="button" class={entityCard} onClick={() => navigate({ to: '/vendors' })}>
              <div class={entityCardLeft}>
                <div class={entityCardInfo}>
                  <span class={entityCardName}>{vendor.name}</span>
                  <span class={entityCardMeta}>
                    Rate: ${vendor.hourly_rate_usd}/hr • TAT: {vendor.turnaround_hours}h
                  </span>
                </div>
              </div>
              <ArrowRight size={14} class={entityCardChevron} />
            </button>
          </div>
        )}

        <div class={section}>
          <h3 class={sectionTitle}>
            <FileCode size={12} /> Technical Metadata & Lineage
          </h3>
          <div class={metadataCard}>
            <div class={metadataRow}>
              <span class={metadataKey}>Derived From Master</span>
              <span class={metadataVal}>{pkg?.derived_from_master_version ?? 'V01'}</span>
            </div>
            <div class={metadataRow}>
              <span class={metadataKey}>Redelivery Count</span>
              <span class={metadataVal}>{pkg?.redelivery_count ?? 0}</span>
            </div>
            {Object.entries(meta).map(([k, v]) => (
              <div key={k} class={metadataRow}>
                <span class={metadataKey}>{k}</span>
                <span class={metadataVal}>{String(v)}</span>
              </div>
            ))}
            <div class={metadataRow}>
              <span class={metadataKey}>Created</span>
              <span class={metadataVal}>
                {pkg?.created_at ? formatDateTime(pkg.created_at) : '—'}
              </span>
            </div>
          </div>
        </div>
      </div>

      <footer class={footer}>
        <Button variant="secondary" onClick={() => navigate({ to: '/titles' })}>
          <Film size={14} />
          <span>View Title</span>
        </Button>
      </footer>
    </aside>
  )
}
