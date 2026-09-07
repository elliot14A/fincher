import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Play, RefreshCw, Route as RouteIcon, ShieldAlert, Workflow } from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { PaginationControls } from '#/components/ui/pagination'
import { RunInspector } from '#/features/runs/components/inspector'
import { runsQueryOptions } from '#/features/runs/queryOptions'
import type { ModelsRun } from '#/lib/api'
import { useSelectableRow } from '#/lib/hooks'
import { formatDateTime } from '#/lib/utils'
import {
  cardName,
  contentLayout,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  list,
  loadingState,
  mainListContainer,
  metaDate,
  metaDivider,
  metaRow,
  metaTrigger,
  nameStack,
  page as pageClass,
  pageSubtitle,
  pageTitle,
  pulseDot,
  row,
  rowActive,
  runIcon,
  statusStack,
  stepsCount,
  sublabelTag,
  timeStack,
  timeValue,
  toolbar,
  toolbarGroup,
  toolbarTab,
  toolbarTabActive,
} from '#/styles/routes/runs.css'

export const Route = createFileRoute('/runs')({
  component: RunsPage,
})

const TABS = [
  { id: 'ALL', label: 'All Workflows' },
  { id: 'INCIDENT', label: 'Incident Remediation' },
  { id: 'ALLOCATION', label: 'Vendor Allocation' },
  { id: 'RESOLUTION', label: 'Closed-Loop Resolution' },
] as const

type TabId = (typeof TABS)[number]['id']

function mapRunStatus(status: ModelsRun['status']): {
  label: string
  variant: BadgeProps['variant']
  isPulsing: boolean
} {
  switch (status) {
    case 'COMPLETED':
      return { label: 'Completed', variant: 'success', isPulsing: false }
    case 'RUNNING':
      return { label: 'Executing', variant: 'warning', isPulsing: true }
    case 'FAILED':
      return { label: 'Failed', variant: 'danger', isPulsing: false }
    case 'ESCALATED':
      return { label: 'Escalated', variant: 'danger', isPulsing: false }
    case 'PENDING':
      return { label: 'Pending', variant: 'neutral', isPulsing: false }
    default:
      return { label: status ?? 'Unknown', variant: 'neutral', isPulsing: false }
  }
}

function getTriggerIcon(trigger: string) {
  switch (trigger?.toLowerCase()) {
    case 'incident':
      return ShieldAlert
    case 'allocation':
      return RouteIcon
    case 'resolution':
      return RefreshCw
    default:
      return Workflow
  }
}

function calculateDuration(startedAt?: string, endedAt?: string): string {
  if (!startedAt) return '-'
  const start = new Date(startedAt).getTime()
  const end = endedAt ? new Date(endedAt).getTime() : Date.now()
  const diffMs = Math.max(0, end - start)
  if (diffMs < 1000) return `${diffMs}ms`
  const diffSec = (diffMs / 1000).toFixed(1)
  return `${diffSec}s`
}

function RunRow({
  run,
  isSelected,
  onSelect,
}: {
  run: ModelsRun
  isSelected: boolean
  onSelect: () => void
}) {
  const { rowProps } = useSelectableRow({
    isSelected,
    onSelect,
    baseClassName: row,
    activeClassName: rowActive,
  })
  const Icon = getTriggerIcon(run.trigger)
  const statusInfo = mapRunStatus(run.status)
  const duration = calculateDuration(run.started_at, run.ended_at)

  const meta = (run.metadata as Record<string, unknown> | undefined) ?? {}
  const titleName =
    typeof meta.title_name === 'string'
      ? meta.title_name
      : run.title_slug && run.title_slug !== 'SYSTEM'
        ? run.title_slug
        : run.id
  const incidentType = typeof meta.incident_type === 'string' ? meta.incident_type : undefined
  const allocationTarget =
    typeof meta.allocation_target === 'string'
      ? `${meta.allocation_target}`
      : typeof meta.allocation_type === 'string'
        ? meta.allocation_type
        : undefined
  const resolutionLoop = typeof meta.resolution_loop === 'string' ? meta.resolution_loop : undefined
  const sublabel = incidentType || allocationTarget || resolutionLoop

  return (
    <div {...rowProps}>
      <div class={runIcon}>
        <Icon size={16} />
      </div>

      <div class={nameStack}>
        <span class={cardName}>{titleName}</span>
        <div class={metaRow}>
          <span class={metaTrigger}>{run.trigger}</span>
          {sublabel && <span class={sublabelTag}>{sublabel}</span>}
          <span class={metaDivider}>•</span>
          <span class={metaDate}>
            {run.started_at ? formatDateTime(run.started_at) : 'Just now'}
          </span>
        </div>
      </div>

      <div class={statusStack}>
        <Badge variant={statusInfo.variant}>
          {statusInfo.isPulsing && <span class={pulseDot} />}
          {statusInfo.label}
        </Badge>
      </div>

      <div class={timeStack}>
        <span class={timeValue}>{duration}</span>
        <span class={stepsCount}>
          {run.steps?.length ?? 0} {run.steps?.length === 1 ? 'step' : 'steps'}
        </span>
      </div>
    </div>
  )
}

function RunsPage() {
  const [activeTab, setActiveTab] = useState<TabId>('ALL')
  const [page, setPage] = useState(1)
  const [selectedRunId, setSelectedRunId] = useState<string | null>(null)

  const { data: runsResult, isLoading } = useQuery(
    runsQueryOptions({
      trigger: activeTab,
      page,
      limit: 15,
    }),
  )

  const runs = runsResult?.items ?? []
  const currentSelected = selectedRunId ?? (runs.length > 0 ? runs[0].id : null)

  const handleTabChange = (tabId: TabId) => {
    setActiveTab(tabId)
    setPage(1)
    setSelectedRunId(null)
  }

  return (
    <div class={pageClass}>
      <header class={header}>
        <div>
          <h1 class={pageTitle}>AI Agent Workflow Runs</h1>
          <span class={pageSubtitle}>
            Real-time execution traces, step latency waterfall, ClickHouse MCP queries, and LLM
            policy verification rationales
          </span>
        </div>
      </header>

      <div class={toolbar}>
        <div class={toolbarGroup}>
          {TABS.map((tab) => (
            <button
              key={tab.id}
              type="button"
              class={activeTab === tab.id ? `${toolbarTab} ${toolbarTabActive}` : toolbarTab}
              onClick={() => handleTabChange(tab.id)}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <div class={contentLayout}>
        <div class={mainListContainer}>
          {isLoading ? (
            <div class={loadingState}>Loading agent workflow runs from database...</div>
          ) : runs.length === 0 ? (
            <div class={emptyState}>
              <Play size={32} />
              <div class={emptyTitle}>No Agent Runs Found</div>
              <div class={emptyText}>
                No workflow runs matching filter "{activeTab}". Runs trigger automatically on media
                ingestion events, title onboarding, or incident remediations.
              </div>
            </div>
          ) : (
            <>
              <div class={list}>
                {runs.map((run: ModelsRun) => (
                  <RunRow
                    key={run.id}
                    run={run}
                    isSelected={currentSelected === run.id}
                    onSelect={() => setSelectedRunId(run.id)}
                  />
                ))}
              </div>

              <PaginationControls
                page={runsResult?.page ?? page}
                totalPages={runsResult?.total_pages ?? 1}
                hasNextPage={runsResult?.has_next_page ?? false}
                hasPrevPage={runsResult?.has_prev_page ?? false}
                onPrevPage={() => setPage((p) => Math.max(1, p - 1))}
                onNextPage={() => setPage((p) => p + 1)}
              />
            </>
          )}
        </div>

        <RunInspector runId={currentSelected} onClose={() => setSelectedRunId(null)} />
      </div>
    </div>
  )
}
