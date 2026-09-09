import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import {
  Building2,
  ChevronRight,
  Clock,
  DollarSign,
  Globe,
  Layers,
  ShieldCheck,
  X,
} from 'lucide-preact'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { getComponentIcon, mapPackageStatus, packagesQueryOptions } from '#/features/packages'
import { vendorDetailQueryOptions } from '#/features/vendors/queryOptions'
import { formatDate, formatDateTime } from '#/lib/utils/formatDate'
import {
  badgeGroup,
  capabilityTag,
  closeBtn,
  complianceBanner,
  complianceIcon,
  emptyNotice,
  entityAvatar,
  footer,
  header,
  headerMain,
  itemChevron,
  itemCode,
  itemDetail,
  itemLeft,
  itemList,
  itemRight,
  itemRow,
  marketTag,
  metadataCard,
  metadataKey,
  metadataRow,
  metadataVal,
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
  titleStack,
} from './vendorSidebar.css'

export type VendorSidebarProps = {
  vendorId: string | null
  onClose: () => void
  onSelectPackage?: (pkgId: string) => void
}

export function VendorSidebar({ vendorId, onClose, onSelectPackage }: VendorSidebarProps) {
  const navigate = useNavigate()

  const vendorQuery = useQuery({
    ...vendorDetailQueryOptions(vendorId ?? ''),
    enabled: Boolean(vendorId),
  })
  const vendor = vendorQuery.data

  const packagesQuery = useQuery({
    ...packagesQueryOptions({ vendor_id: vendorId ?? '', limit: 50 }),
    enabled: Boolean(vendorId),
  })
  const packages = packagesQuery.data?.items ?? []

  if (!vendorId) return null

  const meta = (vendor?.metadata as Record<string, unknown> | undefined) ?? {}
  const posterUrl = typeof meta.poster_url === 'string' ? meta.poster_url : undefined
  const componentsList = vendor?.components ?? []
  const marketsList = vendor?.markets ?? []

  const validCount = packages.filter((p) => p.status === 'VALID').length
  const passRate =
    packages.length > 0 ? `${Math.round((validCount / packages.length) * 100)}%` : '100%'

  return (
    <aside class={sidebarPanel({ size: 'standard' })} aria-label="Vendor Details">
      <header class={header}>
        <div class={headerMain}>
          {posterUrl ? (
            <img src={posterUrl} alt={vendor?.name ?? 'Vendor'} class={entityAvatar} />
          ) : (
            <div class={entityAvatar}>
              <Building2 size={22} />
            </div>
          )}
          <div class={titleStack}>
            <h2 class={sidebarTitle}>{vendor?.name ?? vendorId}</h2>
            <div class={sidebarMeta}>
              <span>{vendor?.id ?? vendorId}</span>
              <span>•</span>
              <Badge variant="success">Active Facility</Badge>
            </div>
          </div>
        </div>
        <button type="button" class={closeBtn} onClick={onClose} aria-label="Close sidebar">
          <X size={16} />
        </button>
      </header>

      <div class={scrollArea}>
        <div class={complianceBanner}>
          <ShieldCheck size={16} class={complianceIcon} />
          <div>
            <strong>MPAA / TPN Certified Facility:</strong> Compliant with studio security
            standards, digital watermarking, and air-gapped Dolby Atmos mixing environments.
          </div>
        </div>

        <div class={section}>
          <h3 class={sectionTitle}>
            <DollarSign size={12} /> Commercial &amp; SLA Metrics
          </h3>
          <div class={statGrid({ columns: 3 })}>
            <div class={statCard}>
              <span class={statValue}>${vendor?.hourly_rate_usd?.toFixed(2) ?? '0.00'}/hr</span>
              <span class={statLabel}>Billing Rate</span>
            </div>
            <div class={statCard}>
              <span class={statValue}>{vendor?.turnaround_hours ?? 24}h</span>
              <span class={statLabel}>Turnaround SLA</span>
            </div>
            <div class={statCard}>
              <span class={statValue}>{passRate}</span>
              <span class={statLabel}>QC Pass Rate</span>
            </div>
          </div>
        </div>

        <div class={section}>
          <h3 class={sectionTitle}>
            <Layers size={12} /> Supported Media Components
          </h3>
          <div class={badgeGroup}>
            {componentsList.length === 0 ? (
              <span class={capabilityTag}>General Post-Production</span>
            ) : (
              componentsList.map((comp) => (
                <span key={comp} class={capabilityTag}>
                  {comp === 'VIDEO'
                    ? 'Video Mastering & QC'
                    : comp === 'AUDIO'
                      ? 'Dolby Atmos & 5.1 Dubbing'
                      : comp === 'SUBTITLE'
                        ? 'Timed Subtitles & Closed Captions'
                        : comp}
                </span>
              ))
            )}
          </div>
        </div>

        <div class={section}>
          <h3 class={sectionTitle}>
            <Globe size={12} /> Localized Market Coverage
          </h3>
          <div class={badgeGroup}>
            {marketsList.length === 0 ? (
              <span class={marketTag}>GLOBAL</span>
            ) : (
              marketsList.map((m) => (
                <span key={m} class={marketTag}>
                  {m}
                </span>
              ))
            )}
          </div>
        </div>

        <div class={section}>
          <h3 class={sectionTitle}>
            <Layers size={12} /> Active Package Allocations ({packages.length})
          </h3>
          {packages.length === 0 ? (
            <div class={emptyNotice}>No media packages currently allocated to this facility.</div>
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
                    onClick={() => {
                      if (onSelectPackage) {
                        onSelectPackage(pkg.id)
                      } else {
                        navigate({ to: '/deliveries' })
                      }
                    }}
                    aria-label={`View package ${pkg.id}`}
                  >
                    <div class={itemLeft}>
                      <CompIcon size={14} />
                      <span class={itemCode}>{pkg.id}</span>
                      <span class={itemDetail}>
                        {pkg.language} • {pkg.component}
                      </span>
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
            <Clock size={12} /> Facility Registration
          </h3>
          <div class={metadataCard}>
            <div class={metadataRow}>
              <span class={metadataKey}>Partner ID</span>
              <span class={metadataVal}>{vendor?.id ?? vendorId}</span>
            </div>
            <div class={metadataRow}>
              <span class={metadataKey}>Onboarded</span>
              <span class={metadataVal}>
                {vendor?.created_at ? formatDate(vendor.created_at) : '—'}
              </span>
            </div>
            <div class={metadataRow}>
              <span class={metadataKey}>Last Verified</span>
              <span class={metadataVal}>
                {vendor?.updated_at ? formatDateTime(vendor.updated_at) : '—'}
              </span>
            </div>
          </div>
        </div>
      </div>

      <footer class={footer}>
        <Button variant="secondary" onClick={() => navigate({ to: '/deliveries' })}>
          <Layers size={14} />
          <span>Deliveries &amp; Packages</span>
        </Button>
      </footer>
    </aside>
  )
}
