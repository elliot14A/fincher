import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { FileCode, Film, Headphones, MessageSquare, Play, Subtitles, Trash2 } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { ActionMenu } from '#/components/ui/dropdown'
import { DeleteModal } from '#/components/ui/modal'
import { PaginationControls } from '#/components/ui/pagination'
import { packagesKeys } from '#/features/packages/queryKeys'
import { packagesQueryOptions } from '#/features/packages/queryOptions'
import { deletePackagesById, type ModelsPackage } from '#/lib/api'
import { useDisclosure, useSelectableRow, useTabbedQueryList } from '#/lib/hooks'
import { formatDateTime } from '#/lib/utils'
import {
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
} from '#/styles/routes/runs.css'

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

  const statusInfo = mapPackageStatus(pkg.status)
  const Icon = getComponentIcon(pkg.component)
  const formattedDate = formatDateTime(pkg.updated_at, 'Registered')

  return (
    <div {...rowProps}>
      <div class={componentIcon}>
        <Icon size={18} />
      </div>

      <div class={nameStack}>
        <span class={cardName}>{pkg.id}</span>
        <div class={metaRow}>
          <span class={metaVersion}>Master {pkg.derived_from_master_version || 'V01'}</span>
          <span class={metaDivider}>·</span>
          <span class={metaVendor}>Lang: {pkg.language || 'en'}</span>
          <span class={metaDivider}>·</span>
          <span class={metaVendor}>Vendor: {pkg.vendor_id}</span>
        </div>
      </div>

      <div class={statusStack}>
        <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
      </div>

      <div class={scheduleStack}>
        <span class={scheduleLabel}>Last evaluated</span>
        <span class={countdownValue}>{formattedDate}</span>
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
              type: 'action',
              key: 'titles',
              label: 'View Title',
              icon: Film,
              onClick: () => navigate({ to: '/titles' }),
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

function RunsPage() {
  const queryClient = useQueryClient()
  const deleteModal = useDisclosure()
  const [deletingPackage, setDeletingPackage] = useState<ModelsPackage | null>(null)

  const {
    activeTab,
    onTabChange,
    setSelectedId,
    currentSelectedId,
    page,
    onPrevPage,
    onNextPage,
    items: packages,
    totalPages,
    hasNextPage,
    hasPrevPage,
    isLoading,
    isError,
    error,
  } = useTabbedQueryList<ModelsPackage, TabId>({
    tabs: TABS,
    buildQueryOptions: ({ filter, page, limit, sort_order }) =>
      packagesQueryOptions({
        component: filter as ModelsPackage['component'],
        page,
        limit,
        sort_order,
      }),
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      const { error } = await deletePackagesById({ path: { id } })
      if (error) throw new Error(error.message || 'Failed to delete package')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: packagesKeys.all })
      toast.success(`Package "${deletingPackage?.id ?? ''}" deleted`)
      setDeletingPackage(null)
      deleteModal.close()
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to delete package')
    },
  })

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Media Packages &amp; Component Runs</h1>
          <span class={pageSubtitle}>
            Derived video, dubbing, and subtitle package lineage states across master cuts.
          </span>
        </div>
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
          {packages.map((pkg) => (
            <PackageRow
              key={pkg.id}
              pkg={pkg}
              isSelected={pkg.id === currentSelectedId}
              onSelect={() => setSelectedId(pkg.id)}
              onDelete={() => {
                setDeletingPackage(pkg)
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

      {deletingPackage ? (
        <DeleteModal
          isOpen={deleteModal.isOpen}
          onClose={() => {
            deleteModal.close()
            setDeletingPackage(null)
          }}
          onConfirm={() => deleteMutation.mutate(deletingPackage.id)}
          entityType="Package"
          entityName={deletingPackage.id}
          entityId={deletingPackage.id}
          warningMessage="Deleting this package will remove its QC validation history and lineage links."
          isDeleting={deleteMutation.isPending}
        />
      ) : null}
    </div>
  )
}
