import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import {
  Bot,
  Braces,
  CheckCircle2,
  Check as CheckIcon,
  ChevronDown,
  ChevronRight,
  Clock,
  Copy,
  Film,
  ShieldAlert,
  ShieldCheck,
  Workflow,
  X,
  XCircle,
} from 'lucide-preact'
import { useEffect, useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Badge, type BadgeProps } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import {
  FlowFullscreen,
  formatLatency,
  isApprovedOutcome,
  isRejectedOutcome,
  narrativeFor,
  num,
  RunFlow,
  type StepMeta,
  stepLatencyMs,
  str,
} from '#/features/runs/components/graph'
import { runDetailQueryOptions } from '#/features/runs/queryOptions'
import type { ModelsRun, ModelsStep, ModelsWfResult } from '#/lib/api'
import {
  actionIcon,
  actionList,
  actionReason,
  actionRow,
  actionText,
  actionType,
  attemptPill,
  chip,
  chipLabel,
  chipRow,
  chipValue,
  closeBtn,
  copyBtn,
  dataCellPrimary,
  dataCellRight,
  dataCellSecondary,
  dataTable,
  dataTableHeadCell,
  dataTableRow,
  emptyState,
  emptyText,
  emptyTitle,
  footer,
  header,
  headerTopRight,
  headerTopRow,
  inspectorPanel,
  judgeCard,
  judgeCardApproved,
  judgeCardRejected,
  judgeHeader,
  judgeHeaderLeft,
  judgeLoop,
  judgeName,
  judgeRationale,
  pulseDot,
  rawJsonPre,
  rawToggle,
  riskBandLabel,
  riskBanner,
  riskBannerBreach,
  riskBannerUrgent,
  riskHeaderRow,
  riskMetric,
  riskMetricLabel,
  riskMetrics,
  riskMetricValue,
  riskMetricValueDanger,
  runIdRow,
  runIdText,
  runTitle,
  scrollArea,
  sectionHeading,
  sqlBlock,
  statBar,
  statCard,
  statLabel,
  statValue,
  stepAgent,
  stepDescription,
  stepTitle,
  stepTitleGroup,
  stepTitleRow,
  summaryLede,
  summaryStrong,
  titleStack,
  verdictBanner,
  verdictBannerApproved,
  verdictBannerRejected,
} from './runInspector.css'

export type RunInspectorProps = {
  runId: string | null
  onClose?: () => void
  onSelectRun?: (runId: string) => void
}

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

function triggerHeadline(trigger: string, titleName: string): string {
  switch (trigger.toLowerCase()) {
    case 'incident':
      return `An incident was raised on ${titleName}. The agent triaged it, gathered live context, planned a fix under policy review, and executed.`
    case 'allocation':
      return `A localization order came in for ${titleName}. The agent selected vendors and provisioned the delivery pipeline.`
    case 'resolution':
      return `A downstream event fired for ${titleName}. The agent re-evaluated the workflow to drive it toward resolution.`
    default:
      return `The agent ran a ${trigger} workflow for ${titleName}.`
  }
}

function Projection({ projection }: { projection: Record<string, unknown> }) {
  const band = str(projection.risk_band).toUpperCase()
  const isBreach = projection.is_breached === true || band === 'BREACH'
  const isUrgent = !isBreach && (projection.is_urgent === true || band === 'URGENT')
  const bannerClass = isBreach
    ? `${riskBanner} ${riskBannerBreach}`
    : isUrgent
      ? `${riskBanner} ${riskBannerUrgent}`
      : riskBanner
  const hours = num(projection.hours_until_premiere)
  const buffer = num(projection.buffer_hours)
  const repairs = Array.isArray(projection.repairs)
    ? (projection.repairs as Record<string, unknown>[])
    : []

  return (
    <div class={bannerClass}>
      <div class={riskHeaderRow}>
        {isBreach ? <ShieldAlert size={13} /> : <Clock size={13} />}
        <span class={riskBandLabel}>{band || 'ON TRACK'}</span>
      </div>
      <div class={riskMetrics}>
        <div class={riskMetric}>
          <span class={riskMetricValue}>{hours}h</span>
          <span class={riskMetricLabel}>Until premiere</span>
        </div>
        <div class={riskMetric}>
          <span
            class={buffer < 0 ? `${riskMetricValue} ${riskMetricValueDanger}` : riskMetricValue}
          >
            {buffer > 0 ? `+${buffer}` : buffer}h
          </span>
          <span class={riskMetricLabel}>Buffer</span>
        </div>
        <div class={riskMetric}>
          <span class={riskMetricValue}>{repairs.length}</span>
          <span class={riskMetricLabel}>Repairs</span>
        </div>
      </div>

      {repairs.length > 0 && (
        <div class={dataTable}>
          <div class={dataTableRow}>
            <span class={dataTableHeadCell}>Component / Market</span>
            <span class={dataTableHeadCell}>Vendor</span>
            <span class={dataTableHeadCell}>ETA</span>
          </div>
          {repairs.map((r, i) => (
            <div class={dataTableRow} key={str(r.package_id) || i}>
              <span class={dataCellPrimary}>
                {str(r.component)}
                {str(r.market) ? ` · ${str(r.market)}` : ''}
              </span>
              <span class={dataCellSecondary}>{str(r.vendor_name) || str(r.vendor_id)}</span>
              <span class={dataCellRight}>{num(r.turnaround_hours)}h</span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

function ExecutedActions({ meta }: { meta: StepMeta }) {
  let artifacts: Record<string, unknown>[] = []
  const raw = meta.artifacts_json
  if (typeof raw === 'string') {
    try {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed)) artifacts = parsed
    } catch {
      artifacts = []
    }
  } else if (Array.isArray(meta.artifacts)) {
    artifacts = meta.artifacts as Record<string, unknown>[]
  }

  if (artifacts.length === 0) return null

  return (
    <div class={actionList}>
      {artifacts.map((a, i) => (
        <div class={actionRow} key={str(a.dispatch_id) || str(a.target_id) || i}>
          <CheckCircle2 size={14} class={actionIcon} />
          <div class={actionText}>
            <span class={actionType}>{str(a.type) || 'ACTION'}</span>
            {str(a.reason) && <span class={actionReason}>{str(a.reason)}</span>}
          </div>
        </div>
      ))}
    </div>
  )
}

function StepDetail({ step }: { step: ModelsStep }) {
  const narrative = narrativeFor(step.name)
  const meta = (step.metadata as StepMeta | undefined) ?? {}
  const extras: preact.JSX.Element[] = []

  if (step.name === 'context_gathering' && meta.projection && typeof meta.projection === 'object') {
    extras.push(<Projection key="proj" projection={meta.projection as Record<string, unknown>} />)
  }
  if (step.name === 'remediation_executor' && (meta.artifacts_json || meta.artifacts)) {
    extras.push(<ExecutedActions key="actions" meta={meta} />)
  }
  if (step.name === 'vendor_selection' && str(meta.overall_summary)) {
    extras.push(
      <p key="vsum" class={stepDescription}>
        {str(meta.overall_summary)}
      </p>,
    )
  }
  if (step.name === 'provisioning') {
    const chips = [
      { value: num(meta.deliveries_created), label: 'deliveries' },
      { value: num(meta.packages_created), label: 'packages' },
      { value: num(meta.qc_scheduled), label: 'QC jobs' },
    ]
    extras.push(
      <div key="chips" class={chipRow}>
        {chips.map((c) => (
          <span class={chip} key={c.label}>
            <span class={chipValue}>{c.value}</span>
            <span class={chipLabel}>{c.label}</span>
          </span>
        ))}
      </div>,
    )
  }
  if (typeof meta.query === 'string') {
    extras.push(
      <div key="sql" class={sqlBlock}>
        {meta.query}
      </div>,
    )
  }

  return (
    <div>
      <div class={stepTitleRow}>
        <div class={stepTitleGroup}>
          <span class={stepTitle}>{narrative.title}</span>
          <span class={stepAgent}>{narrative.agent}</span>
        </div>
        <span class={stepAgent}>
          {formatLatency(stepLatencyMs(step.started_at, step.ended_at))}
        </span>
      </div>
      <p class={stepDescription}>{narrative.describe(meta)}</p>
      {extras}
    </div>
  )
}

function JudgeLoop({ results }: { results: ModelsWfResult[] }) {
  return (
    <div class={judgeLoop}>
      {results.map((res, idx) => {
        const approved = isApprovedOutcome(res.outcome)
        const rejected = isRejectedOutcome(res.outcome)
        const cardClass = approved
          ? `${judgeCard} ${judgeCardApproved}`
          : rejected
            ? `${judgeCard} ${judgeCardRejected}`
            : judgeCard
        const verdictClass = approved
          ? `${verdictBanner} ${verdictBannerApproved}`
          : rejected
            ? `${verdictBanner} ${verdictBannerRejected}`
            : verdictBanner
        return (
          <div class={cardClass} key={res.id || idx}>
            <div class={judgeHeader}>
              <div class={judgeHeaderLeft}>
                <ShieldCheck size={14} />
                <span class={judgeName}>{res.judge}</span>
              </div>
              <span class={attemptPill}>Attempt {res.attempt ?? 1}</span>
            </div>
            <div class={verdictClass}>
              {approved ? <CheckCircle2 size={13} /> : rejected ? <XCircle size={13} /> : null}
              <span>{res.outcome}</span>
            </div>
            {res.rationale && <p class={judgeRationale}>{res.rationale}</p>}
          </div>
        )
      })}
    </div>
  )
}

export function RunInspector({ runId, onClose, onSelectRun }: RunInspectorProps) {
  const navigate = useNavigate()
  const [copied, setCopied] = useState(false)
  const [showRaw, setShowRaw] = useState(false)
  const [selectedStepId, setSelectedStepId] = useState<string | null>(null)
  const [maximized, setMaximized] = useState(false)

  const { data: run, isLoading } = useQuery({
    ...runDetailQueryOptions(runId ?? ''),
    enabled: Boolean(runId),
  })

  const steps = run?.steps ?? []

  useEffect(() => {
    setSelectedStepId(null)
    setMaximized(false)
  }, [runId])

  if (!runId) {
    return (
      <aside class={inspectorPanel}>
        <div class={emptyState}>
          <Bot size={36} />
          <div class={emptyTitle}>Select an Agent Workflow Run</div>
          <div class={emptyText}>
            Pick a run on the left to see the agent flow — what triggered it, what each agent did,
            and how the policy judge ruled.
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

  const totalRunMs = stepLatencyMs(run.started_at, run.ended_at)
  const results = run.results ?? []

  const activeStepId = selectedStepId ?? steps[0]?.id ?? null
  const activeStep = steps.find((s) => (s.id ?? '') === activeStepId) ?? steps[0]

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
                {copied ? <CheckIcon size={11} /> : <Copy size={11} />}
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
            <span class={statLabel}>Trigger</span>
            <span class={statValue}>{run.trigger.toUpperCase()}</span>
          </div>
          <div class={statCard}>
            <span class={statLabel}>Stages</span>
            <span class={statValue}>{steps.length}</span>
          </div>
          <div class={statCard}>
            <span class={statLabel}>Duration</span>
            <span class={statValue}>{formatLatency(totalRunMs)}</span>
          </div>
        </div>
      </header>

      <div class={scrollArea}>
        <p class={summaryLede}>
          <span class={summaryStrong}>{run.trigger.toUpperCase()}. </span>
          {triggerHeadline(run.trigger, titleName)}
        </p>

        {steps.length > 0 && (
          <div>
            <div class={sectionHeading}>
              <Workflow size={12} /> Agent flow — tap a stage
            </div>
            <RunFlow
              steps={steps}
              selectedStepId={activeStepId}
              onSelectStep={setSelectedStepId}
              onMaximize={() => setMaximized(true)}
            />
          </div>
        )}

        {activeStep && <StepDetail step={activeStep} />}

        {results.length > 0 && (
          <div>
            <div class={sectionHeading}>
              <ShieldCheck size={12} /> Policy verification
            </div>
            <JudgeLoop results={results} />
          </div>
        )}

        <div>
          <button type="button" class={rawToggle} onClick={() => setShowRaw(!showRaw)}>
            {showRaw ? <ChevronDown size={12} /> : <ChevronRight size={12} />}
            <Braces size={12} /> Raw JSON payload
          </button>
          {showRaw && <pre class={rawJsonPre}>{JSON.stringify(run, null, 2)}</pre>}
        </div>
      </div>

      <footer class={footer}>
        <Button variant="secondary" size="sm" onClick={() => navigate({ to: '/titles' })}>
          <Film size={13} />
          <span>View Titles</span>
        </Button>
      </footer>

      {maximized && (
        <FlowFullscreen
          run={run}
          titleName={titleName}
          steps={steps}
          selectedStepId={activeStepId}
          onSelectStep={setSelectedStepId}
          onSelectRun={(id) => {
            onSelectRun?.(id)
            setSelectedStepId(null)
          }}
          onClose={() => setMaximized(false)}
        />
      )}
    </aside>
  )
}
