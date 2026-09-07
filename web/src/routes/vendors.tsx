import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Building2, Layers, MessageSquare, Plus, Trash2, Users } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { ActionMenu } from '#/components/ui/dropdown'
import { DeleteModal } from '#/components/ui/modal'
import { PaginationControls } from '#/components/ui/pagination'
import { CreateVendorModal } from '#/features/vendors/components/modals'
import { VendorSidebar } from '#/features/vendors/components/sidebar'
import { vendorsKeys } from '#/features/vendors/queryKeys'
import { vendorsQueryOptions } from '#/features/vendors/queryOptions'
import { deleteVendorsById, type ModelsVendor } from '#/lib/api'
import { useDisclosure, useSelectableRow, useTabbedQueryList } from '#/lib/hooks'
import {
  actions,
  cardName,
  componentBadge,
  contentLayout,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  mainListContainer,
  marketBadge,
  metaRow,
  nameHeader,
  nameStack,
  page as pageClass,
  pageSubtitle,
  pageTitle,
  rateLabel,
  rateStack,
  rateValue,
  row,
  rowActive,
  statusStack,
  tatLabel,
  tatStack,
  tatValue,
  toolbar,
  toolbarGroup,
  toolbarTab,
  toolbarTabActive,
  vendorAvatar,
  vendorCode,
} from '#/styles/routes/vendors.css'

export const Route = createFileRoute('/vendors')({
  component: VendorsPage,
})

const TABS = [
  { id: 'ALL', label: 'All' },
  { id: 'AUDIO', label: 'Audio Dubbing' },
  { id: 'SUBTITLE', label: 'Subtitles' },
  { id: 'VIDEO', label: 'QC Lab' },
] as const

type TabId = (typeof TABS)[number]['id']

function VendorRow({
  vendor,
  isSelected,
  onSelect,
  onDelete,
}: {
  vendor: ModelsVendor
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

  const posterUrl = (vendor.metadata as Record<string, string> | undefined)?.poster_url
  const componentsList = vendor.components ?? []
  const marketsList = vendor.markets ?? []
  const hourlyRate = vendor.hourly_rate_usd
    ? `$${vendor.hourly_rate_usd.toFixed(2)}/hr`
    : '$0.00/hr'
  const tat = vendor.turnaround_hours ? `${vendor.turnaround_hours}h SLA` : '24h SLA'

  return (
    <div {...rowProps}>
      {posterUrl ? (
        <img src={posterUrl} alt={vendor.name} class={vendorAvatar} />
      ) : (
        <div class={vendorAvatar}>
          <Building2 size={18} />
        </div>
      )}

      <div class={nameStack}>
        <div class={nameHeader}>
          <span class={cardName}>{vendor.name}</span>
          <span class={vendorCode}>{vendor.id}</span>
        </div>
        <div class={metaRow}>
          {componentsList.map((comp) => (
            <span key={comp} class={componentBadge}>
              {comp}
            </span>
          ))}
          {marketsList.length > 0 ? (
            <span class={marketBadge}>
              {marketsList.length === 5 ? 'All 5 Markets' : marketsList.join(', ')}
            </span>
          ) : (
            <span class={marketBadge}>Global</span>
          )}
        </div>
      </div>

      <div class={rateStack}>
        <span class={rateLabel}>Billing Rate</span>
        <span class={rateValue}>{hourlyRate}</span>
      </div>

      <div class={tatStack}>
        <span class={tatLabel}>Turnaround</span>
        <span class={tatValue}>{tat}</span>
      </div>

      <div class={statusStack}>
        <Badge variant="success">Active</Badge>
      </div>

      <div class={actions}>
        <ActionMenu
          ariaLabel={`Actions for ${vendor.name}`}
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
              label: 'Inspect Packages',
              icon: Layers,
              onClick: () => navigate({ to: '/deliveries' }),
            },
            {
              type: 'divider',
              key: 'div-1',
            },
            {
              type: 'action',
              key: 'delete',
              label: 'Delete Vendor',
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

function VendorsPage() {
  const queryClient = useQueryClient()
  const createModal = useDisclosure()
  const deleteModal = useDisclosure()
  const [deletingVendor, setDeletingVendor] = useState<ModelsVendor | null>(null)

  const {
    activeTab,
    onTabChange,
    setSelectedId,
    currentSelectedId,
    page,
    onPrevPage,
    onNextPage,
    items: vendors,
    totalPages,
    hasNextPage,
    hasPrevPage,
    isLoading,
    isError,
    error,
  } = useTabbedQueryList<ModelsVendor, TabId>({
    tabs: TABS,
    buildQueryOptions: ({ filter, page, limit, sort_order }) =>
      vendorsQueryOptions({
        component: filter,
        page,
        limit,
        sort_order,
      }),
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      const { error } = await deleteVendorsById({ path: { id } })
      if (error) throw new Error(error.message || 'Failed to delete vendor')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: vendorsKeys.all })
      toast.success(`Vendor "${deletingVendor?.name ?? ''}" deleted`)
      setDeletingVendor(null)
      deleteModal.close()
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to delete vendor')
    },
  })

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Vendor Facilities</h1>
          <span class={pageSubtitle}>
            Dubbing studios, subtitling vendors, and post-production facilities track record.
          </span>
        </div>

        <Button variant="primary" size="sm" onClick={createModal.open}>
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
              onClick={() => onTabChange(tab.id)}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <div class={contentLayout}>
        <div class={mainListContainer}>
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
              {vendors.map((vendor) => (
                <VendorRow
                  key={vendor.id}
                  vendor={vendor}
                  isSelected={vendor.id === currentSelectedId}
                  onSelect={() => setSelectedId(currentSelectedId === vendor.id ? null : vendor.id)}
                  onDelete={() => {
                    setDeletingVendor(vendor)
                    deleteModal.open()
                  }}
                />
              ))}
            </div>
          )}
        </div>

        {currentSelectedId ? (
          <VendorSidebar vendorId={currentSelectedId} onClose={() => setSelectedId(null)} />
        ) : null}
      </div>

      <PaginationControls
        page={page}
        totalPages={totalPages}
        hasNextPage={hasNextPage}
        hasPrevPage={hasPrevPage}
        onPrevPage={onPrevPage}
        onNextPage={onNextPage}
      />

      <CreateVendorModal isOpen={createModal.isOpen} onClose={createModal.close} />

      {deletingVendor ? (
        <DeleteModal
          isOpen={deleteModal.isOpen}
          onClose={() => {
            deleteModal.close()
            setDeletingVendor(null)
          }}
          onConfirm={() => deleteMutation.mutate(deletingVendor.id)}
          entityType="Vendor"
          entityName={deletingVendor.name}
          entityId={deletingVendor.id}
          warningMessage="Deleting this vendor will prevent packages from referencing this facility."
          isDeleting={deleteMutation.isPending}
        />
      ) : null}
    </div>
  )
}
