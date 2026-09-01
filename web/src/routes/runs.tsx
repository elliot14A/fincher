import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import {
  AlertTriangle,
  Bot,
  CheckCircle2,
  Clock,
  Play,
  RotateCcw,
  Sparkles,
  Workflow,
  XCircle,
} from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { PaginationControls } from '#/components/ui/pagination'
import { runDetailQueryOptions, runsQueryOptions } from '#/features/runs'
import type { ModelsRun, ModelsStep, ModelsWfResult } from '#/lib/api'
import { useSelectableRow } from '#/lib/hooks'
import { formatDateTime } from '#/lib/utils'
import {
  attemptBadge,
  cardName,
  contentLayout,
  contextGrid,
  contextLabel,
  contextValue,
  emptyState,
  emptyText,
  emptyTitle,
  header,
  inspectorHeader,
  inspectorPanel,
  inspectorSubtitle,
  inspectorTitle,
  inspectorTopRow,
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
  rationaleBox,
  resultCard,
  resultHeader,
  resultJudge,
  resultOutcome,
  resultsList,
  row,
  rowActive,
  runIcon,
  sectionHeading,
  statusStack,
  stepCard,
  stepDurationText,
  stepHeader,
  stepHeaderLeft,
  stepMeta,
  stepName,
  stepsCount,
  stepsTimeline,
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
      return AlertTriangle
    case 'allocation':
      return Sparkles
    case 'resolution':
      return RotateCcw
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

  return (
    <div {...rowProps}>
      <div class={runIcon}>
        <Icon size={16} />
      </div>

      <div class={nameStack}>
        <span class={cardName}>
          {run.title_slug && run.title_slug !== 'SYSTEM' ? run.title_slug : run.id}
        </span>
        <div class={metaRow}>
          <span class={metaTrigger}>{run.trigger}</span>
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

function RunInspector({ runId }: { runId: string | null }) {
  const { data: run, isLoading } = useQuery({
    ...runDetailQueryOptions(runId ?? ''),
    enabled: Boolean(runId),
  })

  if (!runId) {
    return (
      <aside class={inspectorPanel}>
        <div class={emptyState}>
          <Bot size={36} />
          <div class={emptyTitle}>Select an Agent Run</div>
          <div class={emptyText}>
            Click any workflow execution trace on the left to inspect step transitions, LLM
            rationales, and policy decisions.
          </div>
        </div>
      </aside>
    )
  }

  if (isLoading || !run) {
    return (
      <aside class={inspectorPanel}>
        <div class={loadingState}>Loading execution trace...</div>
      </aside>
    )
  }

  const statusInfo = mapRunStatus(run.status)
  const duration = calculateDuration(run.started_at, run.ended_at)

  return (
    <aside class={inspectorPanel}>
      <div class={inspectorHeader}>
        <div class={inspectorTopRow}>
          <h3 class={inspectorTitle}>Run Execution Inspector</h3>
          <Badge variant={statusInfo.variant}>
            {statusInfo.isPulsing && <span class={pulseDot} />}
            {statusInfo.label}
          </Badge>
        </div>
        <span class={inspectorSubtitle}>ID: {run.id}</span>
      </div>

      <div>
        <h4 class={sectionHeading}>Workflow Context</h4>
        <div class={contextGrid}>
          <div>
            <span class={contextLabel}>Trigger:</span>{' '}
            <strong class={contextValue}>{run.trigger.toUpperCase()}</strong>
          </div>
          <div>
            <span class={contextLabel}>Duration:</span>{' '}
            <strong class={contextValue}>{duration}</strong>
          </div>
          <div>
            <span class={contextLabel}>Target Title:</span>{' '}
            <strong class={contextValue}>{run.title_slug ?? 'N/A'}</strong>
          </div>
          <div>
            <span class={contextLabel}>Started:</span>{' '}
            <span>{run.started_at ? formatDateTime(run.started_at) : '-'}</span>
          </div>
        </div>
      </div>

      {run.steps && run.steps.length > 0 && (
        <div>
          <h4 class={sectionHeading}>Execution Steps ({run.steps.length})</h4>
          <div class={stepsTimeline}>
            {run.steps.map((step: ModelsStep, index: number) => {
              const isCompleted = step.status === 'COMPLETED'
              const isFailed = step.status === 'FAILED'
              const stepDuration = calculateDuration(step.started_at, step.ended_at)

              return (
                <div key={step.id || index} class={stepCard}>
                  <div class={stepHeader}>
                    <div class={stepHeaderLeft}>
                      {isCompleted ? (
                        <CheckCircle2 size={14} />
                      ) : isFailed ? (
                        <XCircle size={14} />
                      ) : (
                        <Clock size={14} />
                      )}
                      <span class={stepName}>{step.name}</span>
                    </div>
                    <span class={stepDurationText}>{stepDuration}</span>
                  </div>

                  {step.metadata && Object.keys(step.metadata).length > 0 && (
                    <div class={stepMeta}>{JSON.stringify(step.metadata)}</div>
                  )}
                </div>
              )
            })}
          </div>
        </div>
      )}

      {run.results && run.results.length > 0 && (
        <div>
          <h4 class={sectionHeading}>AI Decision Rationales ({run.results.length})</h4>
          <div class={resultsList}>
            {run.results.map((res: ModelsWfResult, index: number) => (
              <div key={res.id || index} class={resultCard}>
                <div class={resultHeader}>
                  <span class={resultJudge}>{res.judge}</span>
                  {res.attempt && <span class={attemptBadge}>Attempt {res.attempt}</span>}
                </div>

                <div class={resultOutcome}>Outcome: {res.outcome}</div>

                {res.rationale && <div class={rationaleBox}>{res.rationale}</div>}
              </div>
            ))}
          </div>
        </div>
      )}
    </aside>
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

  return (
    <div class={pageClass}>
      <header class={header}>
        <div>
          <h1 class={pageTitle}>AI Agent Workflow Runs</h1>
          <span class={pageSubtitle}>
            Real-time execution traces, step transitions, and LLM policy verification rationales
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
              onClick={() => {
                setActiveTab(tab.id)
                setPage(1)
              }}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <div class={contentLayout}>
        <div class={mainListContainer}>
          {isLoading ? (
            <div class={loadingState}>Loading agent workflow runs...</div>
          ) : runs.length === 0 ? (
            <div class={emptyState}>
              <Play size={32} />
              <div class={emptyTitle}>No Agent Runs Found</div>
              <div class={emptyText}>
                No workflow runs matching trigger "{activeTab}". Runs trigger automatically on media
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

        <RunInspector runId={currentSelected} />
      </div>
    </div>
  )
}
