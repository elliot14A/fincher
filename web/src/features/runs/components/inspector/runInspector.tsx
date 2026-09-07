import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import {
  Activity,
  Bot,
  Brain,
  Check,
  CheckCircle2,
  ChevronDown,
  ChevronRight,
  Clock,
  Copy,
  Database,
  Film,
  GitFork,
  Scale,
  Send,
  ShieldCheck,
  Terminal,
  Workflow,
  X,
  XCircle,
} from 'lucide-preact'
import { Fragment } from 'preact'
import { useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { runDetailQueryOptions } from '#/features/runs/queryOptions'
import type { ModelsRun, ModelsStep, ModelsWfResult } from '#/lib/api'
import { formatDateTime } from '#/lib/utils/formatDate'
import {
  attemptPill,
  closeBtn,
  contextCard,
  copyBtn,
  decisionCard,
  decisionHeader,
  decisionList,
  decisionMeta,
  emptyState,
  emptyText,
  emptyTitle,
  footer,
  header,
  headerTopRight,
  headerTopRow,
  inspectorPanel,
  judgeTitle,
  keyValGrid,
  kvKey,
  kvVal,
  latencyBarContainer,
  latencyBarFill,
  latencyText,
  outcomeLabel,
  outcomeRow,
  pillTag,
  pulseDot,
  rationaleCard,
  rationaleText,
  rawJsonHeading,
  rawJsonPre,
  runIdRow,
  runIdText,
  runTitle,
  scrollArea,
  section,
  sectionHeading,
  sqlBlock,
  statBar,
  statCard,
  statLabel,
  statusIconClock,
  statusIconDanger,
  statusIconSuccess,
  statValue,
  stepBody,
  stepCard,
  stepCategoryBadge,
  stepHeaderBtn,
  stepHeaderLeft,
  stepHeaderRight,
  stepIconBox,
  stepIndex,
  stepName,
  tabBtn,
  tabBtnActive,
  tabCountBadge,
  tabsNav,
  tagList,
  titleStack,
  waterfallList,
} from './runInspector.css'

export type RunInspectorProps = {
  runId: string | null
  onClose?: () => void
}

type TabType = 'steps' | 'policy' | 'telemetry'

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

function mapOutcomeVariant(outcome: string): BadgeProps['variant'] {
  const o = outcome.toUpperCase()
  if (o.includes('APPROVE') || o.includes('RELEASE') || o.includes('VALID') || o.includes('PASS')) {
    return 'success'
  }
  if (o.includes('HOLD') || o.includes('REVISE') || o.includes('RETRY')) {
    return 'warning'
  }
  if (o.includes('ESCALATE') || o.includes('REJECT') || o.includes('FAIL')) {
    return 'danger'
  }
  return 'neutral'
}

function getStepCategory(name: string): {
  label: string
  icon: typeof Activity
} {
  switch (name.toUpperCase()) {
    case 'ANOMALY_TRIAGE':
      return { label: 'Triage Engine', icon: Activity }
    case 'HISTORIAN_CH_QUERY':
      return { label: 'ClickHouse MCP', icon: Database }
    case 'DEPENDENCY_GRAPH_WALK':
      return { label: 'Graph Traversal', icon: GitFork }
    case 'ACTION_PLANNING':
      return { label: 'ADK Planner', icon: Brain }
    case 'POLICY_VERIFICATION':
      return { label: 'Policy Judge', icon: ShieldCheck }
    case 'EXECUTOR_DISPATCH':
      return { label: 'Executor', icon: Send }
    default:
      return { label: 'Workflow Node', icon: Terminal }
  }
}

function calculateStepLatencyMs(startedAt?: string, endedAt?: string): number {
  if (!startedAt) return 0
  const start = new Date(startedAt).getTime()
  const end = endedAt ? new Date(endedAt).getTime() : Date.now()
  return Math.max(0, end - start)
}

function formatLatency(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

function StepRow({
  step,
  index,
  totalRunMs,
  isInitiallyOpen,
}: {
  step: ModelsStep
  index: number
  totalRunMs: number
  isInitiallyOpen: boolean
}) {
  const [isOpen, setIsOpen] = useState(isInitiallyOpen)
  const category = getStepCategory(step.name)
  const Icon = category.icon
  const isCompleted = step.status === 'COMPLETED'
  const isFailed = step.status === 'FAILED'
  const isRunning = step.status === 'RUNNING'
  const stepMs = calculateStepLatencyMs(step.started_at, step.ended_at)
  const percentOfTotal =
    totalRunMs > 0 ? Math.min(100, Math.max(10, Math.round((stepMs / totalRunMs) * 100))) : 50

  const meta = (step.metadata as Record<string, unknown> | undefined) ?? {}
  const hasMeta = Object.keys(meta).length > 0

  return (
    <div class={stepCard}>
      <button
        type="button"
        class={stepHeaderBtn}
        onClick={() => setIsOpen(!isOpen)}
        aria-expanded={isOpen}
      >
        <div class={stepHeaderLeft}>
          <span class={stepIndex}>{index + 1}.</span>
          <div class={stepIconBox}>
            <Icon size={13} />
          </div>
          <span class={stepName}>{step.name}</span>
          <span class={stepCategoryBadge}>{category.label}</span>
        </div>

        <div class={stepHeaderRight}>
          <div class={latencyBarContainer} title={`${percentOfTotal}% of total latency`}>
            <div class={latencyBarFill} style={{ width: `${percentOfTotal}%` }} />
          </div>
          <span class={latencyText}>{formatLatency(stepMs)}</span>

          {isCompleted ? (
            <CheckCircle2 size={14} class={statusIconSuccess} />
          ) : isFailed ? (
            <XCircle size={14} class={statusIconDanger} />
          ) : isRunning ? (
            <span class={pulseDot} />
          ) : (
            <Clock size={14} class={statusIconClock} />
          )}

          {hasMeta ? isOpen ? <ChevronDown size={14} /> : <ChevronRight size={14} /> : null}
        </div>
      </button>

      {isOpen && hasMeta && (
        <div class={stepBody}>
          {typeof meta.query === 'string' && (
            <div>
              <div class={sectionHeading}>
                <Database size={11} /> ClickHouse SQL Analytical Query
              </div>
              <div class={sqlBlock}>{meta.query}</div>
            </div>
          )}

          <div class={keyValGrid}>
            {Object.entries(meta).map(([k, v]) => {
              if (k === 'query') return null

              const formattedKey = k.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())

              if (Array.isArray(v)) {
                return (
                  <Fragment key={`arr-${k}`}>
                    <span class={kvKey}>{formattedKey}:</span>
                    <div class={tagList}>
                      {v.map((item) => (
                        <span key={String(item)} class={pillTag}>
                          {String(item)}
                        </span>
                      ))}
                    </div>
                  </Fragment>
                )
              }

              if (typeof v === 'object' && v !== null) {
                return (
                  <Fragment key={`obj-${k}`}>
                    <span class={kvKey}>{formattedKey}:</span>
                    <div class={tagList}>
                      {Object.entries(v).map(([subK, subV]) => (
                        <span key={subK} class={pillTag}>
                          {subK}: {String(subV)}
                        </span>
                      ))}
                    </div>
                  </Fragment>
                )
              }

              return (
                <Fragment key={`val-${k}`}>
                  <span class={kvKey}>{formattedKey}:</span>
                  <span class={kvVal}>{String(v)}</span>
                </Fragment>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}

export function RunInspector({ runId, onClose }: RunInspectorProps) {
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<TabType>('steps')
  const [copied, setCopied] = useState(false)

  const { data: run, isLoading } = useQuery({
    ...runDetailQueryOptions(runId ?? ''),
    enabled: Boolean(runId),
  })

  if (!runId) {
    return (
      <aside class={inspectorPanel}>
        <div class={emptyState}>
          <Bot size={36} />
          <div class={emptyTitle}>Select an Agent Workflow Run</div>
          <div class={emptyText}>
            Click any execution trace on the left to inspect step latency spans, ClickHouse MCP
            queries, and policy verification scorecards.
          </div>
        </div>
      </aside>
    )
  }

  if (isLoading || !run) {
    return (
      <aside class={inspectorPanel}>
        <div class={emptyState}>
          <Clock size={32} />
          <div class={emptyTitle}>Loading Execution Trace...</div>
        </div>
      </aside>
    )
  }

  const statusInfo = mapRunStatus(run.status)
  const meta = (run.metadata as Record<string, unknown> | undefined) ?? {}
  const titleName =
    typeof meta.title_name === 'string'
      ? meta.title_name
      : run.title_slug && run.title_slug !== 'SYSTEM'
        ? run.title_slug
        : run.id

  const totalRunMs = calculateStepLatencyMs(run.started_at, run.ended_at)
  const steps = run.steps ?? []
  const results = run.results ?? []

  const handleCopyId = () => {
    if (run.id) {
      navigator.clipboard.writeText(run.id)
      setCopied(true)
      toast.success('Run ID copied to clipboard')
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <aside class={inspectorPanel} aria-label="Run Execution Details">
      <header class={header}>
        <div class={headerTopRow}>
          <div class={titleStack}>
            <h3 class={runTitle}>{titleName}</h3>
            <div class={runIdRow}>
              <span class={runIdText}>{run.id}</span>
              <button type="button" class={copyBtn} onClick={handleCopyId} aria-label="Copy run ID">
                {copied ? <Check size={11} /> : <Copy size={11} />}
                <span>{copied ? 'Copied' : 'Copy'}</span>
              </button>
            </div>
          </div>

          <div class={headerTopRight}>
            <Badge variant={statusInfo.variant}>
              {statusInfo.isPulsing && <span class={pulseDot} />}
              {statusInfo.label}
            </Badge>
            {onClose && (
              <button type="button" class={closeBtn} onClick={onClose} aria-label="Close inspector">
                <X size={16} />
              </button>
            )}
          </div>
        </div>

        <div class={statBar}>
          <div class={statCard}>
            <span class={statLabel}>Duration</span>
            <span class={statValue}>{formatLatency(totalRunMs)}</span>
          </div>
          <div class={statCard}>
            <span class={statLabel}>Trigger</span>
            <span class={statValue}>{run.trigger.toUpperCase()}</span>
          </div>
          <div class={statCard}>
            <span class={statLabel}>Steps</span>
            <span class={statValue}>{steps.length} total</span>
          </div>
          <div class={statCard}>
            <span class={statLabel}>Model</span>
            <span class={statValue}>gemini-2.5</span>
          </div>
        </div>
      </header>

      <nav class={tabsNav} aria-label="Inspector tabs">
        <button
          type="button"
          class={activeTab === 'steps' ? `${tabBtn} ${tabBtnActive}` : tabBtn}
          onClick={() => setActiveTab('steps')}
        >
          <Workflow size={13} />
          <span>Waterfall &amp; Steps</span>
          <span class={tabCountBadge}>{steps.length}</span>
        </button>

        <button
          type="button"
          class={activeTab === 'policy' ? `${tabBtn} ${tabBtnActive}` : tabBtn}
          onClick={() => setActiveTab('policy')}
        >
          <ShieldCheck size={13} />
          <span>Policy Judges</span>
          <span class={tabCountBadge}>{results.length}</span>
        </button>

        <button
          type="button"
          class={activeTab === 'telemetry' ? `${tabBtn} ${tabBtnActive}` : tabBtn}
          onClick={() => setActiveTab('telemetry')}
        >
          <Terminal size={13} />
          <span>Context &amp; Telemetry</span>
        </button>
      </nav>

      <div class={scrollArea}>
        {activeTab === 'steps' && (
          <div class={section}>
            <div class={sectionHeading}>
              <Workflow size={12} /> Execution Latency Waterfall ({steps.length} spans)
            </div>
            <div class={waterfallList}>
              {steps.map((step, idx) => (
                <StepRow
                  key={step.id || idx}
                  step={step}
                  index={idx}
                  totalRunMs={totalRunMs}
                  isInitiallyOpen={idx === 0 || idx === steps.length - 1}
                />
              ))}
            </div>
          </div>
        )}

        {activeTab === 'policy' && (
          <div class={section}>
            <div class={sectionHeading}>
              <Scale size={12} /> Policy Verification Scorecards ({results.length} evaluations)
            </div>

            {results.length === 0 ? (
              <div class={emptyState}>
                <ShieldCheck size={28} />
                <div class={emptyTitle}>No Policy Evaluations Recorded</div>
                <div class={emptyText}>
                  This workflow execution did not trigger bounded policy loop verification.
                </div>
              </div>
            ) : (
              <div class={decisionList}>
                {results.map((res: ModelsWfResult, idx: number) => {
                  const outcomeVariant = mapOutcomeVariant(res.outcome)
                  return (
                    <div key={res.id || idx} class={decisionCard}>
                      <div class={decisionHeader}>
                        <div class={judgeTitle}>
                          <ShieldCheck size={14} />
                          <span>{res.judge}</span>
                        </div>
                        <span class={attemptPill}>Attempt {res.attempt ?? 1} of 3</span>
                      </div>

                      <div class={outcomeRow}>
                        <span class={outcomeLabel}>Policy Verdict:</span>
                        <Badge variant={outcomeVariant}>{res.outcome}</Badge>
                      </div>

                      <div class={rationaleCard}>
                        <p class={rationaleText}>{res.rationale}</p>
                      </div>

                      <div class={decisionMeta}>
                        <span>Target: {run.title_slug}</span>
                        <span>•</span>
                        <span>Evaluation Model: gemini-2.5-pro</span>
                        <span>•</span>
                        <span>Temp: 0.1</span>
                      </div>
                    </div>
                  )
                })}
              </div>
            )}
          </div>
        )}

        {activeTab === 'telemetry' && (
          <div class={section}>
            <div class={sectionHeading}>
              <Terminal size={12} /> Workflow Context &amp; Metadata
            </div>

            <div class={contextCard}>
              <div class={keyValGrid}>
                <span class={kvKey}>Target Title:</span>
                <span class={kvVal}>
                  {titleName} ({run.title_slug})
                </span>

                <span class={kvKey}>Trigger Reason:</span>
                <span class={kvVal}>{run.trigger.toUpperCase()}</span>

                <span class={kvKey}>Started At:</span>
                <span class={kvVal}>{run.started_at ? formatDateTime(run.started_at) : '—'}</span>

                <span class={kvKey}>Completed At:</span>
                <span class={kvVal}>
                  {run.ended_at ? formatDateTime(run.ended_at) : 'In progress'}
                </span>

                {Object.entries(meta).map(([k, v]) => {
                  if (k === 'title_name') return null
                  const formattedKey = k.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
                  return (
                    <Fragment key={`meta-${k}`}>
                      <span class={kvKey}>{formattedKey}:</span>
                      <span class={kvVal}>
                        {typeof v === 'object' ? JSON.stringify(v) : String(v)}
                      </span>
                    </Fragment>
                  )
                })}
              </div>
            </div>

            <div class={rawJsonHeading}>
              <Terminal size={12} /> Raw JSON Payload
            </div>
            <pre class={rawJsonPre}>{JSON.stringify(run, null, 2)}</pre>
          </div>
        )}
      </div>

      <footer class={footer}>
        <Button variant="secondary" size="sm" onClick={() => navigate({ to: '/titles' })}>
          <Film size={13} />
          <span>View Titles</span>
        </Button>
        <Button variant="primary" size="sm" onClick={() => navigate({ to: '/' })}>
          <Bot size={13} />
          <span>Ask Assistant</span>
        </Button>
      </footer>
    </aside>
  )
}
