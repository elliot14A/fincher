import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Film, Layers, MessageSquare, Plus, Trash2 } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { ActionMenu } from '#/components/ui/dropdown'
import { DeleteModal } from '#/components/ui/modal'
import { PaginationControls } from '#/components/ui/pagination'
import { CreateTitleModal } from '#/features/titles/components/modals'
import { titlesKeys } from '#/features/titles/queryKeys'
import { titlesQueryOptions } from '#/features/titles/queryOptions'
import { deleteTitlesById, type ModelsTitle } from '#/lib/api'
import { useCountdown, useDisclosure, useSelectableRow, useTabbedQueryList } from '#/lib/hooks'
import {
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
} from '#/styles/routes/titles.css'

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

function TitleCountdown({ premiereDate }: { premiereDate: string | undefined }) {
  const schedule = useCountdown(premiereDate)

  if (!schedule.scheduled) {
    return (
      <div class={scheduleStack}>
        <span class={scheduleLabelMuted}>Not scheduled</span>
        <span class={countdownEmpty}>-</span>
      </div>
    )
  }

  return (
    <div class={scheduleStack}>
      <span class={scheduleLabel}>{schedule.label}</span>
      <span class={countdownValue}>{schedule.timecode}</span>
    </div>
  )
}

function TitleRow({
  titleItem,
  isSelected,
  onSelect,
  onDelete,
}: {
  titleItem: ModelsTitle
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

  const statusInfo = mapTitleStatus(titleItem.overall_status)
  const territoriesText = `${titleItem.territories || 0} ${
    titleItem.territories === 1 ? 'market' : 'markets'
  }`
  const masterText = `Master ${titleItem.current_master_version || 'V01'}`
  const noteText = getTitleStatusNote(titleItem.overall_status)
  const avatarUrl = (titleItem.metadata as Record<string, string> | undefined)?.avatar_url

  return (
    <div {...rowProps}>
      {avatarUrl ? (
        <img src={avatarUrl} alt={titleItem.name} class={posterThumb} />
      ) : (
        <div class={posterThumb}>
          <Film size={18} />
        </div>
      )}

      <div class={nameStack}>
        <span class={cardName}>{titleItem.name}</span>
        <span class={metaRow}>
          <span class={metaVersion}>{masterText}</span>
          <span class={metaDivider}>·</span>
          <span class={metaTerritories}>{territoriesText}</span>
        </span>
      </div>

      <div class={statusStack}>
        <Badge variant={statusInfo.variant} class={statusBadge}>
          {statusInfo.label}
        </Badge>
        <span class={statusNote}>{noteText}</span>
      </div>

      <TitleCountdown premiereDate={titleItem.premiere_date} />

      <div class={actions}>
        <ActionMenu
          ariaLabel={`Actions for ${titleItem.name}`}
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
              label: 'Delete Title',
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

function TitlesPage() {
  const queryClient = useQueryClient()
  const createModal = useDisclosure()
  const deleteModal = useDisclosure()
  const [deletingTitle, setDeletingTitle] = useState<ModelsTitle | null>(null)

  const {
    activeTab,
    onTabChange,
    setSelectedId,
    currentSelectedId,
    page,
    onPrevPage,
    onNextPage,
    items: titles,
    totalPages,
    hasNextPage,
    hasPrevPage,
    isLoading,
    isError,
    error,
  } = useTabbedQueryList<ModelsTitle, TabId>({
    tabs: TABS,
    buildQueryOptions: ({ filter, page, limit, sort_order }) =>
      titlesQueryOptions({
        status: filter,
        page,
        limit,
        sort_order,
      }),
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      const { error } = await deleteTitlesById({ path: { id } })
      if (error) throw new Error(error.message || 'Failed to delete title')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: titlesKeys.all })
      toast.success(`Title "${deletingTitle?.name ?? ''}" deleted`)
      setDeletingTitle(null)
      deleteModal.close()
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to delete title')
    },
  })

  return (
    <div class={pageClass}>
      <div class={header}>
        <div>
          <h1 class={pageTitle}>Titles &amp; Releases</h1>
          <span class={pageSubtitle}>
            Catalog titles, premiere schedules, and active master cut revisions.
          </span>
        </div>

        <Button variant="primary" size="sm" onClick={createModal.open}>
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
              onClick={() => onTabChange(tab.id)}
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
          {titles.map((titleItem) => (
            <TitleRow
              key={titleItem.id}
              titleItem={titleItem}
              isSelected={titleItem.id === currentSelectedId}
              onSelect={() => setSelectedId(titleItem.id)}
              onDelete={() => {
                setDeletingTitle(titleItem)
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

      <CreateTitleModal isOpen={createModal.isOpen} onClose={createModal.close} />

      {deletingTitle ? (
        <DeleteModal
          isOpen={deleteModal.isOpen}
          onClose={() => {
            deleteModal.close()
            setDeletingTitle(null)
          }}
          onConfirm={() => deleteMutation.mutate(deletingTitle.id)}
          entityType="Title"
          entityName={deletingTitle.name}
          entityId={deletingTitle.id}
          warningMessage="Deleting this title will remove its master cut records, territory deliveries, and packages."
          isDeleting={deleteMutation.isPending}
        />
      ) : null}
    </div>
  )
}
