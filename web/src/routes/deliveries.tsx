import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Globe, Plus } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { PaginationControls } from '#/components/ui/pagination'
import { deliveriesQueryOptions } from '#/features/deliveries'
import type { ModelsDelivery } from '#/lib/api'
import { usePaginatedList } from '#/lib/hooks'
import {
  actionLink,
  actions,
  cardName,
  countdownValue,
  countryBadge,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  metaRow,
  metaText,
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
} from './-deliveries.css'

export const Route = createFileRoute('/deliveries')({
  component: DeliveriesPage,
})

const TABS = [
  { id: 'ALL', label: 'All' },
  { id: 'HOLD', label: 'Holds' },
  { id: 'READY_TO_SHIP', label: 'Ready' },
  { id: 'SHIPPED', label: 'Shipped' },
] as const

type TabId = (typeof TABS)[number]['id']

function mapDeliveryStatus(status: ModelsDelivery['status'] | undefined): {
  label: string
  variant: BadgeProps['variant']
} {
  switch (status) {
    case 'HOLD':
      return { label: 'Hold', variant: 'danger' }
    case 'READY_TO_SHIP':
      return { label: 'Ready to ship', variant: 'success' }
    case 'SHIPPED':
      return { label: 'Shipped', variant: 'success' }
    case 'PENDING':
      return { label: 'Pending QC', variant: 'warning' }
    default:
      return { label: status ?? 'Pending', variant: 'neutral' }
  }
}

function DeliveriesPage() {
  const [activeTab, setActiveTab] = useState<TabId>('ALL')
  const [selectedId, setSelectedId] = useState<string | null>(null)

  const {
    data: deliveries = [],
    isLoading,
    isError,
    error,
  } = useQuery(deliveriesQueryOptions(activeTab === 'ALL' ? undefined : { status: activeTab }))

  const {
    page,
    totalPages,
    paginatedItems: paginatedDeliveries,
    hasNextPage,
    hasPrevPage,
    onPrevPage,
    onNextPage,
  } = usePaginatedList(deliveries, { limit: 10 })

  const isSelectedPresent = deliveries.some((d) => d.id === selectedId)
  const currentSelectedId = isSelectedPresent ? selectedId : (deliveries[0]?.id ?? null)

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Territory Deliveries</h1>
          <span class={pageSubtitle}>
            Global territory shipping targets, package readiness, and carrier dispatches.
          </span>
        </div>

        <Button variant="primary" size="sm">
          <Plus size={14} />
          <span>New Delivery</span>
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
              }}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {isLoading ? (
        <div class={loadingState}>Loading delivery matrix from database...</div>
      ) : isError ? (
        <div class={emptyState}>
          <div class={emptyTitle}>Failed to load deliveries</div>
          <div class={emptyText}>
            {error instanceof Error ? error.message : 'An unexpected error occurred.'}
          </div>
        </div>
      ) : deliveries.length === 0 ? (
        <div class={emptyState}>
          <Globe size={24} />
          <div class={emptyTitle}>No deliveries found</div>
          <div class={emptyText}>
            {activeTab === 'ALL'
              ? 'No territory deliveries scheduled yet.'
              : `No deliveries matching filter '${activeTab}'.`}
          </div>
        </div>
      ) : (
        <div class={list}>
          {paginatedDeliveries.map((delivery) => {
            const statusInfo = mapDeliveryStatus(delivery.status)
            const isSelected = delivery.id === currentSelectedId
            const formattedDate = delivery.target_date
              ? new Date(delivery.target_date).toLocaleDateString(undefined, {
                  month: 'short',
                  day: 'numeric',
                  year: 'numeric',
                })
              : 'Unscheduled'

            return (
              // biome-ignore lint/a11y/useSemanticElements: row contains nested actions
              <div
                key={delivery.id}
                role="button"
                tabIndex={0}
                class={isSelected ? `${row} ${rowActive}` : row}
                onClick={() => setSelectedId(delivery.id)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault()
                    setSelectedId(delivery.id)
                  }
                }}
              >
                <div class={countryBadge}>{delivery.country}</div>

                <div class={nameStack}>
                  <span class={cardName}>{delivery.id}</span>
                  <div class={metaRow}>
                    <span class={metaText}>Title: {delivery.title_id}</span>
                  </div>
                </div>

                <div class={statusStack}>
                  <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
                </div>

                <span />

                <div class={scheduleStack}>
                  <span class={scheduleLabel}>Target date</span>
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
                    Hold
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}

      <PaginationControls
        page={page}
        totalPages={totalPages}
        hasNextPage={hasNextPage}
        hasPrevPage={hasPrevPage}
        onPrevPage={onPrevPage}
        onNextPage={onNextPage}
      />
    </div>
  )
}
