import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { CheckCircle2, Film, Layers, Plus, Trash2 } from 'lucide-preact'
import { useEffect, useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { ActionMenu } from '#/components/ui/dropdown'
import { DeleteModal } from '#/components/ui/modal'
import { PaginationControls } from '#/components/ui/pagination'
import { deliveriesKeys } from '#/features/deliveries/queryKeys'
import { packagesKeys } from '#/features/packages/queryKeys'
import { runsKeys } from '#/features/runs/queryKeys'
import { getQcGating, getTitleStatusNote, mapTitleStatus } from '#/features/titles'
import { CreateTitleModal } from '#/features/titles/components/modals'
import { TitleSidebar } from '#/features/titles/components/sidebar'
import { titlesKeys } from '#/features/titles/queryKeys'
import { titlesQueryOptions } from '#/features/titles/queryOptions'
import { deleteTitlesById, type ModelsTitle, postTitlesByIdQc } from '#/lib/api'
import {
  syncSimulationAnchorWithTitles,
  useCountdown,
  useDisclosure,
  useSelectableRow,
  useTabbedQueryList,
} from '#/lib/hooks'
import {
  actions,
  cardName,
  contentLayout,
  countdownEmpty,
  countdownValue,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  mainListContainer,
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
  { id: 'DRAFT', label: 'Drafts' },
  { id: 'PROCESSING', label: 'In QC' },
  { id: 'ON_TRACK', label: 'Ready' },
  { id: 'HOLD', label: 'Holds' },
  { id: 'OVERDUE', label: 'Overdue' },
] as const

type TabId = (typeof TABS)[number]['id']

function TitleCountdown({
  premiereDate,
  status,
}: {
  premiereDate: string | undefined
  status?: string
}) {
  const schedule = useCountdown(premiereDate, 100, status)

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
  onSendToQC,
  isSendingQC,
  onDelete,
}: {
  titleItem: ModelsTitle
  isSelected: boolean
  onSelect: () => void
  onSendToQC: () => void
  isSendingQC: boolean
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
  const posterUrl = (titleItem.metadata as Record<string, string> | undefined)?.poster_url
  const { isInQC, canRunQC } = getQcGating(titleItem.overall_status)
  return (
    <div {...rowProps}>
      {posterUrl ? (
        <img src={posterUrl} alt={titleItem.name} class={posterThumb} />
      ) : (
        <div class={posterThumb}>
          <Film size={18} />
        </div>
      )}

      <div class={nameStack}>
        <span class={cardName}>{titleItem.name}</span>
        <span class={metaRow}>
          <span class={metaVersion}>{masterText}</span>
          <span class={metaDivider}>•</span>
          <span class={metaTerritories}>{territoriesText}</span>
        </span>
      </div>

      <div class={statusStack}>
        <Badge variant={statusInfo.variant} class={statusBadge}>
          {statusInfo.label}
        </Badge>
        <span class={statusNote}>{noteText}</span>
      </div>

      <TitleCountdown premiereDate={titleItem.premiere_date} status={titleItem.overall_status} />

      <div class={actions}>
        <ActionMenu
          ariaLabel={`Actions for ${titleItem.name}`}
          items={[
            {
              type: 'action',
              key: 'qc',
              label: isInQC ? 'In QC Inspection' : canRunQC ? 'Send for Master QC' : 'QC Complete',
              icon: CheckCircle2,
              disabled: !canRunQC || isSendingQC,
              onClick: onSendToQC,
            },
            {
              type: 'action',
              key: 'packages',
              label: 'View Packages',
              icon: Layers,
              onClick: () => navigate({ to: '/deliveries', search: { title: titleItem.id } }),
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
  const [selectedTitleId, setSelectedTitleId] = useState<string | null>(null)

  const {
    activeTab,
    onTabChange: onTabChangeInternal,
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

  const [hasUserClosedSidebar, setHasUserClosedSidebar] = useState(false)

  useEffect(() => {
    if (titles.length > 0) {
      syncSimulationAnchorWithTitles(titles)
      if (!selectedTitleId && !hasUserClosedSidebar) {
        setSelectedTitleId(titles[0].id)
      }
    }
  }, [titles, selectedTitleId, hasUserClosedSidebar])

  const sendToQcMutation = useMutation({
    mutationFn: async (id: string) => {
      const { error } = await postTitlesByIdQc({ path: { id } })
      if (error) {
        throw new Error((error as { message?: string }).message || 'Failed to send title to QC')
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: titlesKeys.all })
      queryClient.invalidateQueries({ queryKey: runsKeys.all })
      queryClient.invalidateQueries({ queryKey: deliveriesKeys.all })
      queryClient.invalidateQueries({ queryKey: packagesKeys.all })
      toast.success('Master QC initiated successfully')
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to initiate QC')
    },
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

  const onTabChange = (nextTab: TabId) => {
    onTabChangeInternal(nextTab)
    setSelectedTitleId(null)
  }

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

      <div class={contentLayout}>
        <div class={mainListContainer}>
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
                  isSelected={titleItem.id === selectedTitleId}
                  onSelect={() => {
                    if (titleItem.id === selectedTitleId) {
                      setSelectedTitleId(null)
                      setHasUserClosedSidebar(true)
                    } else {
                      setSelectedTitleId(titleItem.id)
                      setHasUserClosedSidebar(false)
                    }
                  }}
                  onSendToQC={() => sendToQcMutation.mutate(titleItem.id)}
                  isSendingQC={sendToQcMutation.isPending}
                  onDelete={() => {
                    setDeletingTitle(titleItem)
                    deleteModal.open()
                  }}
                />
              ))}
            </div>
          )}
        </div>

        {selectedTitleId ? (
          <TitleSidebar
            titleId={selectedTitleId}
            onClose={() => {
              setSelectedTitleId(null)
              setHasUserClosedSidebar(true)
            }}
            onSendToQC={(id) => sendToQcMutation.mutate(id)}
            isSendingQC={sendToQcMutation.isPending}
          />
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
