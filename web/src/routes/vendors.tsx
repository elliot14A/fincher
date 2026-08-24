import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Building2, Plus, Users } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { PaginationControls } from '#/components/ui/pagination'
import { vendorsQueryOptions } from '#/features/vendors'
import { DEFAULT_PAGE, DEFAULT_PAGE_LIMIT, DEFAULT_SORT_ORDER } from '#/lib/constants'
import {
  actionLink,
  actions,
  cardName,
  countdownValue,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  metaRow,
  metaSpecialty,
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
  vendorAvatar,
} from './-vendors.css'

export const Route = createFileRoute('/vendors')({
  component: VendorsPage,
})

const TABS = [
  { id: 'ALL', label: 'All' },
  { id: 'AUDIO_DUBBING', label: 'Audio Dubbing' },
  { id: 'SUBTITLES', label: 'Subtitles' },
  { id: 'QC_LAB', label: 'QC Lab' },
] as const

type TabId = (typeof TABS)[number]['id']

function VendorsPage() {
  const [activeTab, setActiveTab] = useState<TabId>('ALL')
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [page, setPage] = useState(DEFAULT_PAGE)

  const {
    data: vendorsResult,
    isLoading,
    isError,
    error,
  } = useQuery(
    vendorsQueryOptions({
      specialty: activeTab === 'ALL' ? undefined : activeTab,
      page,
      limit: DEFAULT_PAGE_LIMIT,
      sort_order: DEFAULT_SORT_ORDER,
    }),
  )

  const vendors = vendorsResult?.items ?? []
  const totalPages = vendorsResult?.total_pages ?? 1
  const hasNextPage = vendorsResult?.has_next_page ?? false
  const hasPrevPage = vendorsResult?.has_prev_page ?? false

  const isSelectedPresent = vendors.some((v) => v.id === selectedId)
  const currentSelectedId = isSelectedPresent ? selectedId : (vendors[0]?.id ?? null)

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Vendor Facilities</h1>
          <span class={pageSubtitle}>
            Dubbing studios, subtitling vendors, and post-production facilities track record.
          </span>
        </div>

        <Button variant="primary" size="sm">
          <Plus size={14} />
          <span>New Vendor</span>
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
        <div class={loadingState}>Loading vendors directory from database...</div>
      ) : isError ? (
        <div class={emptyState}>
          <div class={emptyTitle}>Failed to load vendors</div>
          <div class={emptyText}>
            {error instanceof Error ? error.message : 'An unexpected error occurred.'}
          </div>
        </div>
      ) : vendors.length === 0 ? (
        <div class={emptyState}>
          <Users size={24} />
          <div class={emptyTitle}>No vendors found</div>
          <div class={emptyText}>
            {activeTab === 'ALL'
              ? 'No post-production facilities registered yet.'
              : `No vendors found for specialty '${activeTab}'.`}
          </div>
        </div>
      ) : (
        <div class={list}>
          {vendors.map((vendor) => {
            const isSelected = vendor.id === currentSelectedId
            const formattedDate = vendor.created_at
              ? new Date(vendor.created_at).toLocaleDateString(undefined, {
                  month: 'short',
                  day: 'numeric',
                  year: 'numeric',
                })
              : 'Registered'

            return (
              // biome-ignore lint/a11y/useSemanticElements: row contains nested actions
              <div
                key={vendor.id}
                role="button"
                tabIndex={0}
                class={isSelected ? `${row} ${rowActive}` : row}
                onClick={() => setSelectedId(vendor.id)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault()
                    setSelectedId(vendor.id)
                  }
                }}
              >
                <div class={vendorAvatar}>
                  <Building2 size={18} />
                </div>

                <div class={nameStack}>
                  <span class={cardName}>{vendor.name}</span>
                  <div class={metaRow}>
                    <span class={metaSpecialty}>{vendor.specialty}</span>
                  </div>
                </div>

                <div class={statusStack}>
                  <Badge variant="neutral">Active</Badge>
                </div>

                <span />

                <div class={scheduleStack}>
                  <span class={scheduleLabel}>Onboarded</span>
                  <span class={countdownValue}>{formattedDate}</span>
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
                    Audit
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}

      <PaginationControls
        page={vendorsResult?.page ?? page}
        totalPages={totalPages}
        hasNextPage={hasNextPage}
        hasPrevPage={hasPrevPage}
        onPrevPage={() => setPage((p) => Math.max(1, p - 1))}
        onNextPage={() => setPage((p) => Math.min(totalPages, p + 1))}
      />
    </div>
  )
}
