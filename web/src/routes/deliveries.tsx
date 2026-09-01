import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import {
  ChevronDown,
  ChevronRight,
  FileCode,
  Film,
  Globe,
  Headphones,
  Layers,
  MessageSquare,
  Package as PackageIcon,
  Plus,
  Subtitles,
  Trash2,
} from 'lucide-preact'
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
import { packagesKeys } from '#/features/packages/queryKeys'
import { packagesQueryOptions } from '#/features/packages/queryOptions'
import {
  deleteDeliveriesById,
  deletePackagesById,
  type ModelsDelivery,
  type ModelsPackage,
} from '#/lib/api'
import { useDisclosure, useSelectableRow } from '#/lib/hooks'
import { formatDate, formatDateTime } from '#/lib/utils'
import {
  actions,
  cardName,
  cardTitleRow,
  componentIcon,
  countdownValue,
  countryBadge,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  metaDivider,
  metaRow,
  metaText,
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
  subItemComp,
  subItemIcon,
  subItemId,
  subItemLeft,
  subItemRow,
  subItemVendor,
  subListContainer,
  subListEmptyText,
  subListHeader,
  toolbar,
  toolbarGroup,
  toolbarTab,
  toolbarTabActive,
  viewModeBar,
  viewModeButton,
  viewModeButtonActive,
  viewModeGroup,
} from '#/styles/routes/deliveries.css'

export const Route = createFileRoute('/deliveries')({
  component: DeliveriesPage,
})

type ViewMode = 'DELIVERIES' | 'PACKAGES'

const DELIVERY_TABS = [
  { id: 'ALL', label: 'All Territories' },
  { id: 'READY_TO_SHIP', label: 'Ready to Ship' },
  { id: 'HOLD', label: 'Holds' },
  { id: 'SHIPPED', label: 'Shipped' },
] as const

const PACKAGE_TABS = [
  { id: 'ALL', label: 'All Components' },
  { id: 'VIDEO', label: 'Video' },
  { id: 'AUDIO', label: 'Audio Dubs' },
  { id: 'SUBTITLE', label: 'Subtitles' },
] as const

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

function mapPackageStatus(status: ModelsPackage['status']): {
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
  const [isExpanded, setIsExpanded] = useState(false)
  const { rowProps } = useSelectableRow({
    isSelected,
    onSelect: () => {
      onSelect()
      setIsExpanded((prev) => !prev)
    },
    baseClassName: row,
    activeClassName: rowActive,
  })

  const { data: packagesResult } = useQuery({
    ...packagesQueryOptions({ title_id: delivery.title_id, limit: 10 }),
    enabled: isExpanded,
  })
  const constituentPackages = (packagesResult?.items ?? []).filter(
    (p) => p.market === delivery.country || p.component === 'VIDEO' || !p.market,
  )

  const statusInfo = mapDeliveryStatus(delivery.status)
  const formattedDate = formatDate(delivery.target_date, undefined, 'Unscheduled')

  return (
    <>
      <div {...rowProps}>
        <div class={countryBadge}>{delivery.country}</div>

        <div class={nameStack}>
          <div class={cardTitleRow}>
            <span class={cardName}>{delivery.id}</span>
            {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
          </div>
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

      {isExpanded && (
        <div class={subListContainer}>
          <div class={subListHeader}>
            Constituent Package Dependencies ({constituentPackages.length})
          </div>
          {constituentPackages.length === 0 ? (
            <div class={subListEmptyText}>No dependent packages found for this territory yet.</div>
          ) : (
            constituentPackages.map((pkg: ModelsPackage) => {
              const pkgStatus = mapPackageStatus(pkg.status)
              const CompIcon = getComponentIcon(pkg.component)
              return (
                <div key={pkg.id} class={subItemRow}>
                  <div class={subItemLeft}>
                    <CompIcon size={13} class={subItemIcon} />
                    <span class={subItemId}>{pkg.id}</span>
                    <span class={subItemComp}>({pkg.component})</span>
                    {pkg.vendor_id && <span class={subItemVendor}>• Vendor: {pkg.vendor_id}</span>}
                  </div>
                  <Badge variant={pkgStatus.variant}>{pkgStatus.label}</Badge>
                </div>
              )
            })
          )}
        </div>
      )}
    </>
  )
}

function PackageRow({
  pkg,
  isSelected,
  onSelect,
  onDelete,
}: {
  pkg: ModelsPackage
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

  const CompIcon = getComponentIcon(pkg.component)
  const statusInfo = mapPackageStatus(pkg.status)

  return (
    <div {...rowProps}>
      <div class={componentIcon}>
        <CompIcon size={16} />
      </div>

      <div class={nameStack}>
        <span class={cardName}>{pkg.id}</span>
        <div class={metaRow}>
          <span class={metaVersion}>{pkg.version ?? 'V01'}</span>
          <span class={metaDivider}>•</span>
          <span class={metaText}>Title: {pkg.title_id}</span>
          {pkg.market && (
            <>
              <span class={metaDivider}>•</span>
              <span class={metaText}>Market: {pkg.market}</span>
            </>
          )}
          {pkg.vendor_id && (
            <>
              <span class={metaDivider}>•</span>
              <span class={metaText}>Vendor: {pkg.vendor_id}</span>
            </>
          )}
        </div>
      </div>

      <div class={statusStack}>
        <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
      </div>

      <div class={scheduleStack}>
        <span class={scheduleLabel}>Created</span>
        <span class={countdownValue}>
          {pkg.created_at ? formatDateTime(pkg.created_at) : 'Just now'}
        </span>
      </div>

      <div class={actions}>
        <ActionMenu
          ariaLabel={`Actions for ${pkg.id}`}
          items={[
            {
              type: 'action',
              key: 'chat',
              label: 'Ask Assistant',
              icon: MessageSquare,
              onClick: () => navigate({ to: '/' }),
            },
            {
              type: 'divider',
              key: 'div-1',
            },
            {
              type: 'action',
              key: 'delete',
              label: 'Delete Package',
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
  const [viewMode, setViewMode] = useState<ViewMode>('DELIVERIES')
  const [deliveryTab, setDeliveryTab] = useState<string>('ALL')
  const [packageTab, setPackageTab] = useState<string>('ALL')
  const [page, setPage] = useState(1)
  const [selectedId, setSelectedId] = useState<string | null>(null)

  const createModal = useDisclosure()
  const deleteModal = useDisclosure()
  const [deletingItem, setDeletingItem] = useState<{
    id: string
    type: 'delivery' | 'package'
  } | null>(null)

  const { data: deliveriesResult, isLoading: isDeliveriesLoading } = useQuery(
    deliveriesQueryOptions({
      status: deliveryTab !== 'ALL' ? deliveryTab : undefined,
      page,
      limit: 15,
    }),
  )

  const { data: packagesResult, isLoading: isPackagesLoading } = useQuery(
    packagesQueryOptions({
      component: packageTab !== 'ALL' ? (packageTab as ModelsPackage['component']) : undefined,
      page,
      limit: 15,
    }),
  )

  const deleteDeliveryMutation = useMutation({
    mutationFn: async (id: string) => {
      const { error } = await deleteDeliveriesById({ path: { id } })
      if (error) throw error
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deliveriesKeys.all })
      toast.success('Delivery deleted')
      deleteModal.close()
    },
    onError: (err: unknown) => {
      const msg = err instanceof Error ? err.message : 'Failed to delete delivery'
      toast.error(msg)
    },
  })

  const deletePackageMutation = useMutation({
    mutationFn: async (id: string) => {
      const { error } = await deletePackagesById({ path: { id } })
      if (error) throw error
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: packagesKeys.all })
      toast.success('Package deleted')
      deleteModal.close()
    },
    onError: (err: unknown) => {
      const msg = err instanceof Error ? err.message : 'Failed to delete package'
      toast.error(msg)
    },
  })

  const handleDeleteConfirm = () => {
    if (!deletingItem) return
    if (deletingItem.type === 'delivery') {
      deleteDeliveryMutation.mutate(deletingItem.id)
    } else {
      deletePackageMutation.mutate(deletingItem.id)
    }
  }

  const deliveries = deliveriesResult?.items ?? []
  const packages = packagesResult?.items ?? []

  return (
    <div class={pageClass}>
      <header class={header}>
        <div>
          <h1 class={pageTitle}>Deliveries & Packages</h1>
          <span class={pageSubtitle}>
            Territory fulfillment matrix, constituent localized media packages, and release
            integrity
          </span>
        </div>
        {viewMode === 'DELIVERIES' && (
          <Button variant="primary" onClick={createModal.open}>
            <Plus size={14} />
            <span>Create delivery</span>
          </Button>
        )}
      </header>

      <div class={viewModeBar}>
        <div class={viewModeGroup}>
          <button
            type="button"
            class={
              viewMode === 'DELIVERIES'
                ? `${viewModeButton} ${viewModeButtonActive}`
                : viewModeButton
            }
            onClick={() => {
              setViewMode('DELIVERIES')
              setPage(1)
            }}
          >
            <Globe size={14} />
            <span>Territory Deliveries</span>
          </button>
          <button
            type="button"
            class={
              viewMode === 'PACKAGES' ? `${viewModeButton} ${viewModeButtonActive}` : viewModeButton
            }
            onClick={() => {
              setViewMode('PACKAGES')
              setPage(1)
            }}
          >
            <Layers size={14} />
            <span>Media Packages</span>
          </button>
        </div>
      </div>

      <div class={toolbar}>
        <div class={toolbarGroup}>
          {viewMode === 'DELIVERIES'
            ? DELIVERY_TABS.map((tab) => (
                <button
                  key={tab.id}
                  type="button"
                  class={deliveryTab === tab.id ? `${toolbarTab} ${toolbarTabActive}` : toolbarTab}
                  onClick={() => {
                    setDeliveryTab(tab.id)
                    setPage(1)
                  }}
                >
                  {tab.label}
                </button>
              ))
            : PACKAGE_TABS.map((tab) => (
                <button
                  key={tab.id}
                  type="button"
                  class={packageTab === tab.id ? `${toolbarTab} ${toolbarTabActive}` : toolbarTab}
                  onClick={() => {
                    setPackageTab(tab.id)
                    setPage(1)
                  }}
                >
                  {tab.label}
                </button>
              ))}
        </div>
      </div>

      {viewMode === 'DELIVERIES' ? (
        isDeliveriesLoading ? (
          <div class={loadingState}>Loading territory deliveries...</div>
        ) : deliveries.length === 0 ? (
          <div class={emptyState}>
            <Globe size={32} />
            <div class={emptyTitle}>No Deliveries Found</div>
            <div class={emptyText}>
              No territory deliveries match the selected filter. Create a delivery to schedule
              localized platform releases.
            </div>
          </div>
        ) : (
          <>
            <div class={list}>
              {deliveries.map((d: ModelsDelivery) => (
                <DeliveryRow
                  key={d.id}
                  delivery={d}
                  isSelected={selectedId === d.id}
                  onSelect={() => setSelectedId(d.id)}
                  onDelete={() => {
                    setDeletingItem({ id: d.id, type: 'delivery' })
                    deleteModal.open()
                  }}
                />
              ))}
            </div>

            <PaginationControls
              page={deliveriesResult?.page ?? page}
              totalPages={deliveriesResult?.total_pages ?? 1}
              hasNextPage={deliveriesResult?.has_next_page ?? false}
              hasPrevPage={deliveriesResult?.has_prev_page ?? false}
              onPrevPage={() => setPage((p) => Math.max(1, p - 1))}
              onNextPage={() => setPage((p) => p + 1)}
            />
          </>
        )
      ) : isPackagesLoading ? (
        <div class={loadingState}>Loading media packages...</div>
      ) : packages.length === 0 ? (
        <div class={emptyState}>
          <PackageIcon size={32} />
          <div class={emptyTitle}>No Packages Found</div>
          <div class={emptyText}>
            No media asset packages match the selected component filter. Packages are created during
            title onboarding and vendor allocations.
          </div>
        </div>
      ) : (
        <>
          <div class={list}>
            {packages.map((pkg: ModelsPackage) => (
              <PackageRow
                key={pkg.id}
                pkg={pkg}
                isSelected={selectedId === pkg.id}
                onSelect={() => setSelectedId(pkg.id)}
                onDelete={() => {
                  setDeletingItem({ id: pkg.id, type: 'package' })
                  deleteModal.open()
                }}
              />
            ))}
          </div>

          <PaginationControls
            page={packagesResult?.page ?? page}
            totalPages={packagesResult?.total_pages ?? 1}
            hasNextPage={packagesResult?.has_next_page ?? false}
            hasPrevPage={packagesResult?.has_prev_page ?? false}
            onPrevPage={() => setPage((p) => Math.max(1, p - 1))}
            onNextPage={() => setPage((p) => p + 1)}
          />
        </>
      )}

      <CreateDeliveryModal isOpen={createModal.isOpen} onClose={createModal.close} />

      <DeleteModal
        isOpen={deleteModal.isOpen}
        onClose={deleteModal.close}
        onConfirm={handleDeleteConfirm}
        entityType={deletingItem?.type === 'delivery' ? 'Delivery' : 'Package'}
        entityId={deletingItem?.id ?? ''}
        isDeleting={deleteDeliveryMutation.isPending || deletePackageMutation.isPending}
      />
    </div>
  )
}
