import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Globe, Layers, MessageSquare, Plus, Trash2 } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { ActionMenu } from '#/components/ui/dropdown'
import { DeleteModal } from '#/components/ui/modal'
import { PaginationControls } from '#/components/ui/pagination'
import { CreateDeliveryModal } from '#/features/deliveries/components/modals'
import { deliveriesKeys } from '#/features/deliveries/queryKeys'
import { deliveriesQueryOptions } from '#/features/deliveries/queryOptions'
import { deleteDeliveriesById, type ModelsDelivery } from '#/lib/api'
import { useDisclosure, useSelectableRow, useTabbedQueryList } from '#/lib/hooks'
import { formatDate } from '#/lib/utils'
import {
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
} from '#/styles/routes/deliveries.css'

export const Route = createFileRoute('/deliveries')({
  component: DeliveriesPage,
})

const TABS = [
  { id: 'ALL', label: 'All' },
  { id: 'READY_TO_SHIP', label: 'Ready to Ship' },
  { id: 'HOLD', label: 'Holds' },
  { id: 'SHIPPED', label: 'Shipped' },
] as const

type TabId = (typeof TABS)[number]['id']

function mapDeliveryStatus(status: ModelsDelivery['status']): {
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

function DeliveryRow({
  delivery,
  isSelected,
  onSelect,
  onDelete,
}: {
  delivery: ModelsDelivery
  isSelected: boolean
  onSelect: () => void
  onDelete: () => void
}) {
  const navigate = useNavigate()
  const { rowProps } = useSelectableRow({
    isSelected,
    onSelect,
    baseClassName: row,
    activeClassName: rowActive,
  })

  const statusInfo = mapDeliveryStatus(delivery.status)
  const formattedDate = formatDate(delivery.target_date, undefined, 'Unscheduled')

  return (
    <div {...rowProps}>
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

      <div class={scheduleStack}>
        <span class={scheduleLabel}>Target date</span>
        <span class={countdownValue}>{formattedDate}</span>
      </div>

      <div class={actions}>
        <ActionMenu
          ariaLabel={`Actions for ${delivery.id}`}
          items={[
            {
              type: 'action',
              key: 'chat',
              label: 'Ask Assistant',
              icon: MessageSquare,
              onClick: () => navigate({ to: '/' }),
            },
            {
              type: 'action',
              key: 'packages',
              label: 'View Packages',
              icon: Layers,
              onClick: () => navigate({ to: '/runs' }),
            },
            {
              type: 'divider',
              key: 'div-1',
            },
            {
              type: 'action',
              key: 'delete',
              label: 'Delete Delivery',
              icon: Trash2,
              danger: true,
              onClick: onDelete,
            },
          ]}
        />
      </div>
    </div>
  )
}

function DeliveriesPage() {
  const queryClient = useQueryClient()
  const createModal = useDisclosure()
  const deleteModal = useDisclosure()
  const [deletingDelivery, setDeletingDelivery] = useState<ModelsDelivery | null>(null)

  const {
    activeTab,
    onTabChange,
    setSelectedId,
    currentSelectedId,
    page,
    onPrevPage,
    onNextPage,
    items: deliveries,
    totalPages,
    hasNextPage,
    hasPrevPage,
    isLoading,
    isError,
    error,
  } = useTabbedQueryList<ModelsDelivery, TabId>({
    tabs: TABS,
    buildQueryOptions: ({ filter, page, limit, sort_order }) =>
      deliveriesQueryOptions({
        status: filter,
        page,
        limit,
        sort_order,
      }),
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      const { error } = await deleteDeliveriesById({ path: { id } })
      if (error) throw new Error(error.message || 'Failed to delete delivery')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deliveriesKeys.all })
      toast.success(`Delivery "${deletingDelivery?.id ?? ''}" deleted`)
      setDeletingDelivery(null)
      deleteModal.close()
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to delete delivery')
    },
  })

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Market Deliveries</h1>
          <span class={pageSubtitle}>
            Global release shipping targets, package readiness, and carrier dispatches.
          </span>
        </div>

        <Button variant="primary" size="sm" onClick={createModal.open}>
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
              onClick={() => onTabChange(tab.id)}
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
              ? 'No market deliveries scheduled yet.'
              : `No deliveries matching filter '${activeTab}'.`}
          </div>
        </div>
      ) : (
        <div class={list}>
          {deliveries.map((delivery) => (
            <DeliveryRow
              key={delivery.id}
              delivery={delivery}
              isSelected={delivery.id === currentSelectedId}
              onSelect={() => setSelectedId(delivery.id)}
              onDelete={() => {
                setDeletingDelivery(delivery)
                deleteModal.open()
              }}
            />
          ))}
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

      <CreateDeliveryModal isOpen={createModal.isOpen} onClose={createModal.close} />

      {deletingDelivery ? (
        <DeleteModal
          isOpen={deleteModal.isOpen}
          onClose={() => {
            deleteModal.close()
            setDeletingDelivery(null)
          }}
          onConfirm={() => deleteMutation.mutate(deletingDelivery.id)}
          entityType="Delivery"
          entityName={`Market ${deletingDelivery.country}`}
          entityId={deletingDelivery.id}
          warningMessage="Deleting this delivery will cancel this territory shipping commitment."
          isDeleting={deleteMutation.isPending}
        />
      ) : null}
    </div>
  )
}
