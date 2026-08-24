import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Film, Plus } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { PaginationControls } from '#/components/ui/pagination'
import { titlesQueryOptions } from '#/features/titles'
import type { ModelsTitle } from '#/lib/api'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'
import {
  actionLink,
  actions,
  cardName,
  countdownEmpty,
  countdownValue,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  metaDivider,
  metaRow,
  metaTerritories,
  metaVersion,
  nameStack,
  page as pageClass,
  pageSubtitle,
  pageTitle,
  posterThumb,
  row,
  rowActive,
  scheduleLabel,
  scheduleLabelMuted,
  scheduleStack,
  statusBadge,
  statusNote,
  statusStack,
  toolbar,
  toolbarGroup,
  toolbarTab,
  toolbarTabActive,
} from './-titles.css'

export const Route = createFileRoute('/titles')({
  component: TitlesPage,
})

const TABS = [
  { id: 'ALL', label: 'All' },
  { id: 'HOLD', label: 'Holds' },
  { id: 'PROCESSING', label: 'QC' },
  { id: 'ON_TRACK', label: 'Ready' },
] as const

type TabId = (typeof TABS)[number]['id']

function mapTitleStatus(status: ModelsTitle['overall_status'] | undefined): {
  label: string
  variant: BadgeProps['variant']
} {
  switch (status) {
    case 'HOLD':
      return { label: 'Hold', variant: 'danger' }
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

function getTitleStatusNote(status: ModelsTitle['overall_status'] | undefined): string {
  switch (status) {
    case 'HOLD':
      return 'Delivery hold active'
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

function formatPremiereCountdown(premiereDateStr: string | undefined): {
  label: string
  timecode: string
  scheduled: boolean
} {
  if (!premiereDateStr) {
    return { label: 'Not scheduled', timecode: '-', scheduled: false }
  }
  const target = new Date(premiereDateStr).getTime()
  if (Number.isNaN(target)) {
    return { label: 'Not scheduled', timecode: '-', scheduled: false }
  }
  const now = Date.now()
  const diffMs = target - now
  if (diffMs <= 0) {
    return { label: 'Released', timecode: '00:00:00:00', scheduled: true }
  }
  const totalSeconds = Math.floor(diffMs / 1000)
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60
  const frames = Math.floor((diffMs % 1000) / 40)
  const pad = (n: number) => n.toString().padStart(2, '0')
  return {
    label: `${hours}h left`,
    timecode: `${pad(hours)}:${pad(minutes)}:${pad(seconds)}:${pad(frames)}`,
    scheduled: true,
  }
}

function TitlesPage() {
  const [activeTab, setActiveTab] = useState<TabId>('ALL')
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [page, setPage] = useState(DEFAULT_PAGE)

  const {
    data: titlesResult,
    isLoading,
    isError,
    error,
  } = useQuery(
    titlesQueryOptions({
      status: activeTab === 'ALL' ? undefined : activeTab,
      page,
      limit: DEFAULT_PAGE_LIMIT,
      sort_order: DEFAULT_SORT_ORDER,
    }),
  )

  const titles = titlesResult?.items ?? []
  const totalPages = titlesResult?.total_pages ?? 1
  const hasNextPage = titlesResult?.has_next_page ?? false
  const hasPrevPage = titlesResult?.has_prev_page ?? false

  const isSelectedPresent = titles.some((t) => t.id === selectedId)
  const currentSelectedId = isSelectedPresent ? selectedId : (titles[0]?.id ?? null)

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Titles &amp; Releases</h1>
          <span class={pageSubtitle}>
            Catalog titles, premiere schedules, and active master cut revisions.
          </span>
        </div>

        <Button variant="primary" size="sm">
          <Plus size={14} />
          <span>New Title</span>
        </Button>
      </div>

      <div class={toolbar}>
        <div class={toolbarGroup}>
          {TABS.map((tab) => (
            <button
              key={tab.id}
              type="button"
              class={activeTab === tab.id ? `${toolbarTab} ${toolbarTabActive}` : toolbarTab}
              onClick={() => {
                setActiveTab(tab.id)
                setSelectedId(null)
                setPage(DEFAULT_PAGE)
              }}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {isLoading ? (
        <div class={loadingState}>Loading titles from database...</div>
      ) : isError ? (
        <div class={emptyState}>
          <div class={emptyTitle}>Failed to load titles</div>
          <div class={emptyText}>
            {error instanceof Error ? error.message : 'An unexpected error occurred.'}
          </div>
        </div>
      ) : titles.length === 0 ? (
        <div class={emptyState}>
          <Film size={24} />
          <div class={emptyTitle}>No titles found</div>
          <div class={emptyText}>
            {activeTab === 'ALL'
              ? 'No media titles registered yet. Create your first title to begin.'
              : `No titles found matching status '${activeTab}'.`}
          </div>
        </div>
      ) : (
        <div class={list}>
          {titles.map((titleItem) => {
            const statusInfo = mapTitleStatus(titleItem.overall_status)
            const schedule = formatPremiereCountdown(titleItem.premiere_date)
            const territoriesText = `${titleItem.territories || 0} ${
              titleItem.territories === 1 ? 'territory' : 'territories'
            }`
            const masterText = `Master ${titleItem.current_master_version || 'V01'}`
            const noteText = getTitleStatusNote(titleItem.overall_status)
            const isSelected = titleItem.id === currentSelectedId

            return (
              // biome-ignore lint/a11y/useSemanticElements: row contains nested action buttons, so it cannot itself be a <button>
              <div
                key={titleItem.id}
                role="button"
                tabIndex={0}
                class={isSelected ? `${row} ${rowActive}` : row}
                onClick={() => setSelectedId(titleItem.id)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault()
                    setSelectedId(titleItem.id)
                  }
                }}
              >
                <div class={posterThumb}>
                  <Film size={18} />
                </div>

                <div class={nameStack}>
                  <span class={cardName}>{titleItem.name}</span>
                  <span class={metaRow}>
                    <span class={metaVersion}>{masterText}</span>
                    <span class={metaDivider}>·</span>
                    <span class={metaTerritories}>{territoriesText}</span>
                  </span>
                </div>

                <span />

                <div class={statusStack}>
                  <Badge variant={statusInfo.variant} class={statusBadge}>
                    {statusInfo.label}
                  </Badge>
                  <span class={statusNote}>{noteText}</span>
                </div>

                <div class={scheduleStack}>
                  {schedule.scheduled ? (
                    <>
                      <span class={scheduleLabel}>{schedule.label}</span>
                      <span class={countdownValue}>{schedule.timecode}</span>
                    </>
                  ) : (
                    <>
                      <span class={scheduleLabelMuted}>Not scheduled</span>
                      <span class={countdownEmpty}>-</span>
                    </>
                  )}
                </div>

                <div class={actions}>
                  <button
                    type="button"
                    class={actionLink}
                    onClick={(event) => event.stopPropagation()}
                  >
                    Inspect
                  </button>
                  <button
                    type="button"
                    class={actionLink}
                    onClick={(event) => event.stopPropagation()}
                  >
                    More
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}

      <PaginationControls
        page={titlesResult?.page ?? page}
        totalPages={totalPages}
        hasNextPage={hasNextPage}
        hasPrevPage={hasPrevPage}
        onPrevPage={() => setPage((p) => Math.max(1, p - 1))}
        onNextPage={() => setPage((p) => Math.min(totalPages, p + 1))}
      />
    </div>
  )
}
