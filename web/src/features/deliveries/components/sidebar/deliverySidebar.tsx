import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { ArrowRight, ChevronRight, Clock, Film, Globe, Layers, ShieldAlert, X } from 'lucide-preact'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { deliveryDetailQueryOptions } from '#/features/deliveries/queryOptions'
import { getComponentIcon, mapPackageStatus, packagesQueryOptions } from '#/features/packages'
import { titleDetailQueryOptions } from '#/features/titles/queryOptions'
import type { ModelsDelivery } from '#/lib/api'
import { formatDate, formatDateTime } from '#/lib/utils/formatDate'
import {
  closeBtn,
  emptyNotice,
  footer,
  header,
  headerMain,
  holdNotice,
  holdNoticeIcon,
  itemChevron,
  itemCode,
  itemDetail,
  itemLeft,
  itemList,
  itemRight,
  itemRow,
  metadataCard,
  metadataKey,
  metadataRow,
  metadataVal,
  parentTitleCard,
  parentTitleChevron,
  parentTitleInfo,
  parentTitleLeft,
  parentTitleMeta,
  parentTitleName,
  scrollArea,
  section,
  sectionTitle,
  sidebarMeta,
  sidebarPanel,
  sidebarTitle,
  statCard,
  statGrid,
  statLabel,
  statValue,
  territoryAvatar,
  titleStack,
} from './deliverySidebar.css'

export type DeliverySidebarProps = {
  deliveryId: string | null
  onClose: () => void
  onSelectPackage?: (pkgId: string) => void
}

function mapDeliveryStatus(status: ModelsDelivery['status'] | undefined): {
  label: string
  variant: BadgeProps['variant']
} {
  switch (status) {
    case 'HOLD':
      return { label: 'Hold Active', variant: 'danger' }
    case 'READY_TO_SHIP':
      return { label: 'Ready to Ship', variant: 'success' }
    case 'SHIPPED':
      return { label: 'Shipped', variant: 'neutral' }
    case 'PENDING':
      return { label: 'Pending QC', variant: 'warning' }
    default:
      return { label: status ?? 'Pending', variant: 'neutral' }
  }
}

export function DeliverySidebar({ deliveryId, onClose, onSelectPackage }: DeliverySidebarProps) {
  const navigate = useNavigate()

  const deliveryQuery = useQuery({
    ...deliveryDetailQueryOptions(deliveryId ?? ''),
    enabled: Boolean(deliveryId),
  })
  const delivery = deliveryQuery.data

  const titleQuery = useQuery({
    ...titleDetailQueryOptions(delivery?.title_id ?? ''),
    enabled: Boolean(delivery?.title_id),
  })
  const title = titleQuery.data

  const packagesQuery = useQuery({
    ...packagesQueryOptions({ title_id: delivery?.title_id ?? '', limit: 50 }),
    enabled: Boolean(delivery?.title_id),
  })

  if (!deliveryId) return null

  const statusInfo = mapDeliveryStatus(delivery?.status)
  const meta = (delivery?.metadata as Record<string, unknown> | undefined) ?? {}
  const platform = typeof meta.platform === 'string' ? meta.platform : undefined
  const priority = typeof meta.priority === 'string' ? meta.priority : undefined
  const holdReason = typeof meta.hold_reason === 'string' ? meta.hold_reason : undefined

  // Filter packages for this delivery's market + root video
  const packages = (packagesQuery.data?.items ?? []).filter(
    (p) => p.market === delivery?.country || p.component === 'VIDEO' || !p.market,
  )

  const isHold = delivery?.status === 'HOLD'

  return (
    <aside class={sidebarPanel({ size: 'standard' })} aria-label="Delivery Details">
      <header class={header}>
        <div class={headerMain}>
          <div class={territoryAvatar}>{delivery?.country ?? '??'}</div>
          <div class={titleStack}>
            <h2 class={sidebarTitle}>{delivery?.id ?? deliveryId}</h2>
            <div class={sidebarMeta}>
              <span>Territory: {delivery?.country ?? 'Global'}</span>
              <span>•</span>
              <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
            </div>
          </div>
        </div>
        <button type="button" class={closeBtn} onClick={onClose} aria-label="Close sidebar">
          <X size={16} />
        </button>
      </header>

      <div class={scrollArea}>
        {isHold && (
          <div class={holdNotice}>
            <ShieldAlert size={16} class={holdNoticeIcon} />
            <div>
              <strong>Fulfillment Blocked:</strong>{' '}
              {holdReason ?? 'Active hold prevents media package handoff.'}
            </div>
          </div>
        )}

        <div class={section}>
          <h3 class={sectionTitle}>
            <Clock size={12} /> Target Release Milestone
          </h3>
          <div class={statGrid({ columns: 2 })}>
            <div class={statCard}>
              <span class={statValue}>
                {delivery?.target_date ? formatDate(delivery.target_date) : 'Unscheduled'}
              </span>
              <span class={statLabel}>Target Date</span>
            </div>
            <div class={statCard}>
              <span class={statValue}>{priority ?? 'STANDARD'}</span>
              <span class={statLabel}>Priority Tier</span>
            </div>
          </div>
        </div>

        {title && (
          <div class={section}>
            <h3 class={sectionTitle}>
              <Film size={12} /> Parent Media Title
            </h3>
            <button
              type="button"
              class={parentTitleCard}
              onClick={() => navigate({ to: '/titles' })}
            >
              <div class={parentTitleLeft}>
                <div class={parentTitleInfo}>
                  <span class={parentTitleName}>{title.name}</span>
                  <span class={parentTitleMeta}>
                    Master: {title.current_master_version || 'V01'} • Status: {title.overall_status}
                  </span>
                </div>
              </div>
              <ArrowRight size={14} class={parentTitleChevron} />
            </button>
          </div>
        )}

        <div class={section}>
          <h3 class={sectionTitle}>
            <Layers size={12} /> Constituent Media Packages ({packages.length})
          </h3>
          {packages.length === 0 ? (
            <div class={emptyNotice}>No constituent packages found for this territory yet.</div>
          ) : (
            <div class={itemList}>
              {packages.map((pkg) => {
                const pkgStatus = mapPackageStatus(pkg.status)
                const CompIcon = getComponentIcon(pkg.component)
                return (
                  <button
                    key={pkg.id}
                    type="button"
                    class={itemRow}
                    onClick={() => onSelectPackage?.(pkg.id)}
                    aria-label={`Inspect package ${pkg.id}`}
                  >
                    <div class={itemLeft}>
                      <CompIcon size={14} />
                      <span class={itemCode}>{pkg.id}</span>
                      <span class={itemDetail}>({pkg.component})</span>
                    </div>
                    <div class={itemRight}>
                      <Badge variant={pkgStatus.variant}>{pkgStatus.label}</Badge>
                      <ChevronRight size={13} class={itemChevron} />
                    </div>
                  </button>
                )
              })}
            </div>
          )}
        </div>

        <div class={section}>
          <h3 class={sectionTitle}>
            <Globe size={12} /> Platform &amp; Fulfillment Details
          </h3>
          <div class={metadataCard}>
            <div class={metadataRow}>
              <span class={metadataKey}>Platform Target</span>
              <span class={metadataVal}>{platform ?? 'General OTT / Theatrical'}</span>
            </div>
            <div class={metadataRow}>
              <span class={metadataKey}>Created</span>
              <span class={metadataVal}>
                {delivery?.created_at ? formatDateTime(delivery.created_at) : '—'}
              </span>
            </div>
            <div class={metadataRow}>
              <span class={metadataKey}>Last Updated</span>
              <span class={metadataVal}>
                {delivery?.updated_at ? formatDateTime(delivery.updated_at) : '—'}
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
