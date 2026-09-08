import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import {
  ArrowRight,
  Bot,
  Box,
  CheckCircle2,
  Clock,
  Film,
  Globe,
  Layers,
  TrendingUp,
  X,
} from 'lucide-preact'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { deliveriesQueryOptions } from '#/features/deliveries/queryOptions'
import { mapPackageStatus, packagesQueryOptions } from '#/features/packages'
import { runsQueryOptions } from '#/features/runs/queryOptions'
import { getQcGating, mapTitleStatus } from '#/features/titles/lib'
import { titleDetailQueryOptions } from '#/features/titles/queryOptions'
import { useCountdown } from '#/lib/hooks'
import { formatDateTime } from '#/lib/utils/formatDate'
import {
  badgeGroup,
  closeBtn,
  emptyNotice,
  footer,
  header,
  headerMain,
  itemCode,
  itemDetail,
  itemLeft,
  itemList,
  itemRow,
  posterPlaceholder,
  posterThumb,
  progressBar,
  progressCard,
  progressFill,
  progressHeader,
  progressLabel,
  progressPct,
  scrollArea,
  section,
  sectionHeaderRow,
  sectionTitle,
  sidebarPanel,
  skeletonBlock,
  skeletonLine,
  skeletonLineMedium,
  skeletonLineShort,
  skeletonLineThin,
  statCard,
  statGrid,
  statLabel,
  statValue,
  statValueDanger,
  synopsisBox,
  titleInfo,
  titleMeta,
  titleName,
} from './titleSidebar.css'

export type TitleSidebarProps = {
  titleId: string | null
  onClose: () => void
  onSendToQC: (id: string) => void
  isSendingQC: boolean
}

export function TitleSidebar({ titleId, onClose, onSendToQC, isSendingQC }: TitleSidebarProps) {
  const navigate = useNavigate()

  const titleQuery = useQuery({
    ...titleDetailQueryOptions(titleId ?? ''),
    enabled: Boolean(titleId),
  })

  const title = titleQuery.data

  const packagesQuery = useQuery({
    ...packagesQueryOptions({ title_id: titleId ?? '', limit: 100 }),
    enabled: Boolean(titleId),
  })

  const deliveriesQuery = useQuery({
    ...deliveriesQueryOptions({ title_id: titleId ?? '', limit: 50 }),
    enabled: Boolean(titleId),
  })

  const runsQuery = useQuery({
    ...runsQueryOptions({ title_slug: title?.slug, limit: 10 }),
    enabled: Boolean(title?.slug),
  })

  const countdown = useCountdown(title?.premiere_date, 100, title?.overall_status)

  if (!titleId) return null

  const statusInfo = mapTitleStatus(title?.overall_status)
  const meta = (title?.metadata as Record<string, unknown> | undefined) ?? {}
  const posterUrl = typeof meta.poster_url === 'string' ? meta.poster_url : undefined
  const synopsis = typeof meta.synopsis === 'string' ? meta.synopsis : undefined
  const genre = typeof meta.genre === 'string' ? meta.genre : undefined
  const markets = Array.isArray(meta.markets) ? (meta.markets as string[]) : []

  const packages = packagesQuery.data?.items ?? []
  const totalPackages = packages.length
  const validPackages = packages.filter((p) => p.status === 'VALID').length
  const failedPackages = packages.filter((p) => p.status === 'INVALIDATED').length
  const progressPctNum = totalPackages > 0 ? Math.round((validPackages / totalPackages) * 100) : 0

  const deliveries = deliveriesQuery.data?.items ?? []
  const runs = runsQuery.data?.items ?? []

  const { isInQC, canRunQC } = getQcGating(title?.overall_status)

  return (
    <aside
      class={sidebarPanel({ size: 'compact' })}
      aria-label={`Details for ${title?.name ?? 'Title'}`}
    >
      <header class={header}>
        <div class={headerMain}>
          {posterUrl ? (
            <img src={posterUrl} alt={title?.name} class={posterThumb} />
          ) : (
            <div class={posterPlaceholder}>
              <Film size={20} />
            </div>
          )}
          <div class={titleInfo}>
            <h2 class={titleName}>{title?.name ?? 'Loading...'}</h2>
            <div class={titleMeta}>
              <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
              <span>•</span>
              <span>Master {title?.current_master_version || 'V01'}</span>
              {genre ? (
                <>
                  <span>•</span>
                  <span>{genre}</span>
                </>
              ) : null}
            </div>
          </div>
        </div>
        <button
          type="button"
          class={closeBtn}
          onClick={(e) => {
            e.stopPropagation()
            onClose()
          }}
          aria-label="Close sidebar"
        >
          <X size={16} />
        </button>
      </header>

      <div class={scrollArea}>
        <section class={section}>
          <h3 class={sectionTitle}>
            <Clock size={12} />
            <span>Premiere Timeline &amp; Countdown</span>
          </h3>
          <div class={progressCard}>
            <div class={progressHeader}>
              <span class={progressLabel}>
                {countdown.scheduled
                  ? countdown.isPast
                    ? 'Premiere Window'
                    : 'Countdown to Premiere'
                  : 'Unscheduled'}
              </span>
              <span class={progressPct}>{countdown.timecode}</span>
            </div>
            <div class={itemDetail}>
              Target Premiere: {formatDateTime(title?.premiere_date)} ({countdown.label})
            </div>
          </div>
        </section>

        <section class={section}>
          <h3 class={sectionTitle}>
            <TrendingUp size={12} />
            <span>Package Completion Progress</span>
          </h3>
          {packagesQuery.isLoading ? (
            <div class={skeletonBlock}>
              <div class={`${skeletonLine} ${skeletonLineShort}`} />
              <div class={`${skeletonLine} ${skeletonLineThin}`} />
            </div>
          ) : (
            <div class={progressCard}>
              <div class={progressHeader}>
                <span class={progressLabel}>
                  {validPackages} of {totalPackages} Assets Validated
                </span>
                <span class={progressPct}>{progressPctNum}%</span>
              </div>
              <div class={progressBar}>
                <div class={progressFill} style={{ width: `${progressPctNum}%` }} />
              </div>
            </div>
          )}

          <div class={statGrid({ columns: 3 })}>
            <div class={statCard}>
              <span class={statValue}>
                {deliveriesQuery.isLoading ? '...' : (title?.territories ?? deliveries.length ?? 0)}
              </span>
              <span class={statLabel}>Territories</span>
            </div>
            <div class={statCard}>
              <span class={statValue}>{packagesQuery.isLoading ? '...' : totalPackages}</span>
              <span class={statLabel}>Packages</span>
            </div>
            <div class={statCard}>
              <span class={failedPackages > 0 ? `${statValue} ${statValueDanger}` : statValue}>
                {packagesQuery.isLoading ? '...' : failedPackages}
              </span>
              <span class={statLabel}>Defects</span>
            </div>
          </div>
        </section>

        {markets.length > 0 ? (
          <section class={section}>
            <h3 class={sectionTitle}>
              <Globe size={12} />
              <span>Targeted Language Markets</span>
            </h3>
            <div class={badgeGroup}>
              {markets.map((m) => (
                <Badge key={m} variant="neutral">
                  {m}
                </Badge>
              ))}
            </div>
          </section>
        ) : null}

        {synopsis ? (
          <section class={section}>
            <h3 class={sectionTitle}>
              <Film size={12} />
              <span>Overview &amp; Synopsis</span>
            </h3>
            <p class={synopsisBox}>{synopsis}</p>
          </section>
        ) : null}

        <section class={section}>
          <div class={sectionHeaderRow}>
            <h3 class={sectionTitle}>
              <Layers size={12} />
              <span>Media Packages ({totalPackages})</span>
            </h3>
            {totalPackages > 0 ? (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => navigate({ to: '/deliveries', search: { title: title?.id } })}
              >
                <span>View All</span>
                <ArrowRight size={12} />
              </Button>
            ) : null}
          </div>

          {packagesQuery.isLoading ? (
            <div class={skeletonBlock}>
              <div class={`${skeletonLine} ${skeletonLineMedium}`} />
              <div class={`${skeletonLine} ${skeletonLineShort}`} />
            </div>
          ) : packages.length === 0 ? (
            <div class={emptyNotice}>
              No component packages generated yet. Initiate Master QC to bundle audio, video, and
              subtitles.
            </div>
          ) : (
            <div class={itemList}>
              {packages.slice(0, 4).map((pkg) => {
                const pkgStatus = mapPackageStatus(pkg.status)
                return (
                  <div key={pkg.id} class={itemRow}>
                    <div class={itemLeft}>
                      <Badge variant={pkgStatus.variant}>{pkg.component}</Badge>
                      <span class={itemCode}>{pkg.market || pkg.language}</span>
                      <span class={itemDetail}>Master {pkg.derived_from_master_version}</span>
                    </div>
                    <Badge variant={pkgStatus.variant}>{pkg.status}</Badge>
                  </div>
                )
              })}
            </div>
          )}
        </section>

        <section class={section}>
          <div class={sectionHeaderRow}>
            <h3 class={sectionTitle}>
              <Bot size={12} />
              <span>Autonomous Agent Runs ({runs.length})</span>
            </h3>
            {runs.length > 0 ? (
              <Button variant="ghost" size="sm" onClick={() => navigate({ to: '/runs' })}>
                <span>History</span>
                <ArrowRight size={12} />
              </Button>
            ) : null}
          </div>

          {runs.length === 0 ? (
            <div class={emptyNotice}>
              No autonomous agent workflows triggered for this title yet.
            </div>
          ) : (
            <div class={itemList}>
              {runs.slice(0, 3).map((run) => (
                <div key={run.id} class={itemRow}>
                  <div class={itemLeft}>
                    <span class={itemCode}>Run #{run.id.slice(0, 8)}</span>
                    <span class={itemDetail}>{run.trigger}</span>
                  </div>
                  <Badge
                    variant={
                      run.status === 'COMPLETED'
                        ? 'success'
                        : run.status === 'FAILED'
                          ? 'danger'
                          : 'warning'
                    }
                  >
                    {run.status}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </section>
      </div>

      <footer class={footer}>
        <Button
          variant="secondary"
          size="sm"
          onClick={() => navigate({ to: '/deliveries', search: { title: title?.id } })}
        >
          <Box size={14} />
          <span>Deliveries &amp; Packages</span>
        </Button>

        <Button
          variant="primary"
          size="sm"
          disabled={!canRunQC || isSendingQC}
          onClick={() => canRunQC && title?.id && onSendToQC(title.id)}
        >
          <CheckCircle2 size={14} />
          <span>{isInQC ? 'In QC' : canRunQC ? 'Initiate Master QC' : 'QC Complete'}</span>
        </Button>
      </footer>
    </aside>
  )
}
