import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { FileCode, Film, Headphones, Play, Plus, Subtitles } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { PaginationControls } from '#/components/ui/pagination'
import { packagesQueryOptions } from '#/features/packages'
import type { ModelsPackage } from '#/lib/api'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'
import {
  actionLink,
  actions,
  cardName,
  componentIcon,
  countdownValue,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  metaDivider,
  metaRow,
  metaVendor,
  metaVersion,
  nameStack,
  page as pageClass,
  pageSubtitle,
  pageTitle,
  row,
  rowActive,
  scheduleLabel,
  scheduleStack,
  statusStack,
  toolbar,
  toolbarGroup,
  toolbarTab,
  toolbarTabActive,
} from './-runs.css'

export const Route = createFileRoute('/runs')({
  component: RunsPage,
})

const TABS = [
  { id: 'ALL', label: 'All Packages' },
  { id: 'VIDEO', label: 'Video' },
  { id: 'AUDIO', label: 'Audio' },
  { id: 'SUBTITLE', label: 'Subtitles' },
] as const

type TabId = (typeof TABS)[number]['id']

function mapPackageStatus(status: ModelsPackage['status']): {
  label: string
  variant: BadgeProps['variant']
} {
  switch (status) {
    case 'VALID':
      return { label: 'Valid', variant: 'success' }
    case 'INVALIDATED':
      return { label: 'Invalidated', variant: 'danger' }
    case 'RE_QC_PENDING':
      return { label: 'Re-QC Pending', variant: 'warning' }
    case 'PENDING':
      return { label: 'Pending QC', variant: 'warning' }
    default:
      return { label: status ?? 'Pending', variant: 'neutral' }
  }
}

function getComponentIcon(component: ModelsPackage['component']) {
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

function RunsPage() {
  const [activeTab, setActiveTab] = useState<TabId>('ALL')
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [page, setPage] = useState(DEFAULT_PAGE)

  const {
    data: packagesResult,
    isLoading,
    isError,
    error,
  } = useQuery(
    packagesQueryOptions({
      component: activeTab === 'ALL' ? undefined : (activeTab as ModelsPackage['component']),
      page,
      limit: DEFAULT_PAGE_LIMIT,
      sort_order: DEFAULT_SORT_ORDER,
    }),
  )

  const packages = packagesResult?.items ?? []
  const totalPages = packagesResult?.total_pages ?? 1
  const hasNextPage = packagesResult?.has_next_page ?? false
  const hasPrevPage = packagesResult?.has_prev_page ?? false

  const isSelectedPresent = packages.some((pkg) => pkg.id === selectedId)
  const currentSelectedId = isSelectedPresent ? selectedId : (packages[0]?.id ?? null)

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Autonomous DAG Runs &amp; Media Packages</h1>
          <span class={pageSubtitle}>
            Derived video, dubbing, and subtitle package lineage states across master cuts.
          </span>
        </div>

        <Button variant="primary" size="sm">
          <Plus size={14} />
          <span>New Package</span>
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
        <div class={loadingState}>Loading media package runs from database...</div>
      ) : isError ? (
        <div class={emptyState}>
          <div class={emptyTitle}>Failed to load packages</div>
          <div class={emptyText}>
            {error instanceof Error ? error.message : 'An unexpected error occurred.'}
          </div>
        </div>
      ) : packages.length === 0 ? (
        <div class={emptyState}>
          <Play size={24} />
          <div class={emptyTitle}>No media packages found</div>
          <div class={emptyText}>
            {activeTab === 'ALL'
              ? 'No derived media packages registered yet.'
              : `No packages matching component '${activeTab}'.`}
          </div>
        </div>
      ) : (
        <div class={list}>
          {packages.map((pkg) => {
            const statusInfo = mapPackageStatus(pkg.status)
            const Icon = getComponentIcon(pkg.component)
            const isSelected = pkg.id === currentSelectedId
            const formattedDate = pkg.updated_at
              ? new Date(pkg.updated_at).toLocaleDateString(undefined, {
                  month: 'short',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })
              : 'Registered'

            return (
              // biome-ignore lint/a11y/useSemanticElements: row contains nested actions
              <div
                key={pkg.id}
                role="button"
                tabIndex={0}
                class={isSelected ? `${row} ${rowActive}` : row}
                onClick={() => setSelectedId(pkg.id)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault()
                    setSelectedId(pkg.id)
                  }
                }}
              >
                <div class={componentIcon}>
                  <Icon size={18} />
                </div>

                <div class={nameStack}>
                  <span class={cardName}>{pkg.id}</span>
                  <div class={metaRow}>
                    <span class={metaVersion}>
                      Master {pkg.derived_from_master_version || 'V01'}
                    </span>
                    <span class={metaDivider}>·</span>
                    <span class={metaVendor}>Lang: {pkg.language || 'en'}</span>
                    <span class={metaDivider}>·</span>
                    <span class={metaVendor}>Vendor: {pkg.vendor_id}</span>
                  </div>
                </div>

                <div class={statusStack}>
                  <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
                </div>

                <span />

                <div class={scheduleStack}>
                  <span class={scheduleLabel}>Last evaluated</span>
                  <span class={countdownValue}>{formattedDate}</span>
                </div>

                <div class={actions}>
                  <button
                    type="button"
                    class={actionLink}
                    onClick={(event) => event.stopPropagation()}
                  >
                    Inspect DAG
                  </button>
                  <button
                    type="button"
                    class={actionLink}
                    onClick={(event) => event.stopPropagation()}
                  >
                    Re-QC
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}

      <PaginationControls
        page={packagesResult?.page ?? page}
        totalPages={totalPages}
        hasNextPage={hasNextPage}
        hasPrevPage={hasPrevPage}
        onPrevPage={() => setPage((p) => Math.max(1, p - 1))}
        onNextPage={() => setPage((p) => Math.min(totalPages, p + 1))}
      />
    </div>
  )
}
