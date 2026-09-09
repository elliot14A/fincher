import { useQuery } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'
import { ArrowRight, CheckCircle2, Send } from 'lucide-preact'
import { useEffect, useMemo, useState } from 'preact/hooks'
import { toast } from 'sonner'
import { Button } from '#/components/ui/button'
import { Select, type SelectOption } from '#/components/ui/select'
import {
  buildEvent,
  DEFECT_CATEGORIES,
  SCENARIOS,
  type ScenarioId,
  type ScenarioInputs,
  scenarioById,
  simulatePackagesQueryOptions,
  simulateTitlesQueryOptions,
  useEmitEvent,
} from '#/features/simulate'
import type { ModelsEvent, ModelsEventBatchResponse } from '#/lib/api'
import {
  body,
  column,
  emitRow,
  field,
  fieldLabel,
  fieldRow,
  header,
  hintText,
  input,
  noRunNote,
  page,
  pageSubtitle,
  pageTitle,
  previewCard,
  previewPre,
  resultCard,
  resultMeta,
  resultTitle,
  runLink,
  runLinkList,
  scenarioCard,
  scenarioCardActive,
  scenarioDesc,
  scenarioGrid,
  scenarioHead,
  scenarioIcon,
  scenarioName,
  sectionLabel,
  severityBadge,
  severityBadgeCritical,
  severityBadgeInfo,
  severityBadgeWarn,
  textarea,
  triggerTag,
} from '#/styles/routes/simulate.css'

export const Route = createFileRoute('/simulate')({
  component: SimulatePage,
})

function severityBadgeClass(severity: ModelsEvent['severity']): string {
  if (severity === 'CRITICAL') return severityBadgeCritical
  if (severity === 'WARN') return severityBadgeWarn
  return severityBadgeInfo
}

const DEFAULT_CUSTOM = JSON.stringify(
  [
    {
      type: 'fincher.package.invalidated',
      source: 'fincher.ops.console',
      subject: 'GLOBAL',
      severity: 'CRITICAL',
      data: { package_id: '', reason: 'manual test' },
    },
  ],
  null,
  2,
)

function SimulatePage() {
  const [scenarioId, setScenarioId] = useState<ScenarioId>('master_revision')
  const [titleSlug, setTitleSlug] = useState('')
  const [titleId, setTitleId] = useState('')
  const [packageId, setPackageId] = useState('')
  const [newMasterVersion, setNewMasterVersion] = useState('v2.0')
  const [invalidationReason, setInvalidationReason] = useState('corrupt source master detected')
  const [driftMs, setDriftMs] = useState(140)
  const [defectCategory, setDefectCategory] = useState<string>(DEFECT_CATEGORIES[0])
  const [customJson, setCustomJson] = useState(DEFAULT_CUSTOM)
  const [result, setResult] = useState<ModelsEventBatchResponse | null>(null)

  const scenario = scenarioById(scenarioId)
  const severity = scenario.defaultSeverity
  const emit = useEmitEvent()

  const { data: titles = [] } = useQuery(simulateTitlesQueryOptions())
  const { data: packages = [] } = useQuery(simulatePackagesQueryOptions(titleId))

  useEffect(() => {
    setResult(null)
  }, [scenarioId])

  useEffect(() => {
    if (!titleSlug && titles.length > 0) {
      setTitleSlug(titles[0].slug)
      setTitleId(titles[0].id)
    }
  }, [titles, titleSlug])

  useEffect(() => {
    setPackageId('')
  }, [titleId])

  const selectedPackage = packages.find((p) => p.id === packageId) ?? packages[0]

  const titleOptions: SelectOption[] = useMemo(
    () =>
      titles.map((t) => ({
        value: t.slug,
        label: t.name,
        meta: t.overall_status,
        keywords: t.slug,
      })),
    [titles],
  )

  const packageOptions: SelectOption[] = useMemo(
    () =>
      packages.map((p) => ({
        value: p.id,
        label: `${p.component} · ${p.market || p.language} · ${p.vendor_id}`,
        meta: p.status,
        keywords: `${p.id} ${p.vendor_id} ${p.component} ${p.market ?? ''} ${p.language}`,
      })),
    [packages],
  )

  const defectOptions: SelectOption[] = DEFECT_CATEGORIES.map((c) => ({ value: c, label: c }))

  const inputs: ScenarioInputs = {
    subject: titleSlug,
    pkg: selectedPackage,
    newMasterVersion,
    invalidationReason,
    driftMs,
    defectCategory,
  }

  const event = useMemo(
    () => (scenario.id === 'custom' ? null : buildEvent(scenario, severity, inputs)),
    [
      scenario,
      severity,
      titleSlug,
      selectedPackage,
      newMasterVersion,
      invalidationReason,
      driftMs,
      defectCategory,
    ],
  )

  const previewText = scenario.id === 'custom' ? customJson : JSON.stringify([event], null, 2)

  const handleTitleChange = (slug: string) => {
    const t = titles.find((x) => x.slug === slug)
    setTitleSlug(slug)
    setTitleId(t?.id ?? '')
  }

  const handleEmit = async () => {
    let payload: ModelsEvent[]
    if (scenario.id === 'custom') {
      try {
        const parsed = JSON.parse(customJson)
        payload = Array.isArray(parsed) ? parsed : [parsed]
      } catch {
        toast.error('Custom payload is not valid JSON')
        return
      }
    } else if (event) {
      payload = [event]
    } else {
      return
    }

    try {
      const res = await emit.mutateAsync(payload)
      setResult(res)
      const runCount = res.run_ids?.length ?? 0
      toast.success(
        runCount > 0
          ? `Emitted — triggered ${runCount} agent run${runCount === 1 ? '' : 's'}`
          : 'Event ingested (no workflow triggered)',
      )
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to emit event')
    }
  }

  return (
    <div class={page}>
      <header class={header}>
        <h1 class={pageTitle}>Event Simulator</h1>
        <span class={pageSubtitle}>
          Emit production events to drive the autonomous agent pipeline live.
        </span>
      </header>

      <div class={body}>
        <div class={column}>
          <div>
            <div class={sectionLabel}>Scenario</div>
            <div class={scenarioGrid}>
              {SCENARIOS.map((s) => {
                const Icon = s.icon
                return (
                  <button
                    key={s.id}
                    type="button"
                    class={
                      s.id === scenarioId ? `${scenarioCard} ${scenarioCardActive}` : scenarioCard
                    }
                    onClick={() => setScenarioId(s.id)}
                  >
                    <div class={scenarioHead}>
                      <Icon size={16} class={scenarioIcon} />
                      <span class={scenarioName}>{s.label}</span>
                    </div>
                    <span class={scenarioDesc}>{s.description}</span>
                    <span class={triggerTag}>{s.trigger}</span>
                  </button>
                )
              })}
            </div>
          </div>

          {scenario.id === 'custom' ? (
            <div class={field}>
              <span class={fieldLabel}>Raw event payload (JSON array)</span>
              <textarea
                class={textarea}
                value={customJson}
                onInput={(e) => setCustomJson((e.target as HTMLTextAreaElement).value)}
                spellcheck={false}
              />
              <span class={hintText}>
                Sent as-is to POST /api/events. Must be a JSON array of events.
              </span>
            </div>
          ) : (
            <div class={column}>
              <div class={fieldRow}>
                <div class={field}>
                  <span class={fieldLabel}>Title</span>
                  <Select
                    options={titleOptions}
                    value={titleSlug}
                    onChange={handleTitleChange}
                    placeholder="Select a title"
                    ariaLabel="Select title"
                  />
                </div>
                <div class={field}>
                  <span class={fieldLabel}>Severity</span>
                  <span class={`${severityBadge} ${severityBadgeClass(severity)}`}>
                    {severity}
                    <span class={hintText}>· auto</span>
                  </span>
                </div>
              </div>

              {scenario.id === 'master_revision' && (
                <div class={field}>
                  <span class={fieldLabel}>New master version</span>
                  <input
                    class={input}
                    value={newMasterVersion}
                    onInput={(e) => setNewMasterVersion((e.target as HTMLInputElement).value)}
                  />
                </div>
              )}

              {scenario.needsPackage && (
                <div class={field}>
                  <span class={fieldLabel}>Package</span>
                  <Select
                    options={packageOptions}
                    value={selectedPackage?.id ?? ''}
                    onChange={setPackageId}
                    placeholder="No packages for this title"
                    ariaLabel="Select package"
                  />
                </div>
              )}

              {scenario.id === 'package_invalidation' && (
                <div class={field}>
                  <span class={fieldLabel}>Reason</span>
                  <input
                    class={input}
                    value={invalidationReason}
                    onInput={(e) => setInvalidationReason((e.target as HTMLInputElement).value)}
                  />
                </div>
              )}

              {scenario.id === 'qc_failed' && (
                <div class={fieldRow}>
                  <div class={field}>
                    <span class={fieldLabel}>Sync drift (ms)</span>
                    <input
                      class={input}
                      type="number"
                      value={String(driftMs)}
                      onInput={(e) => setDriftMs(Number((e.target as HTMLInputElement).value) || 0)}
                    />
                  </div>
                  <div class={field}>
                    <span class={fieldLabel}>Defect category</span>
                    <Select
                      options={defectOptions}
                      value={defectCategory}
                      onChange={setDefectCategory}
                      searchable={false}
                      ariaLabel="Select defect category"
                    />
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        <div class={previewCard}>
          <div class={sectionLabel}>Payload preview</div>
          <pre class={previewPre}>{previewText}</pre>

          <div class={emitRow}>
            <Button variant="primary" onClick={handleEmit} disabled={emit.isPending}>
              <Send size={14} />
              <span>{emit.isPending ? 'Emitting…' : 'Emit event'}</span>
            </Button>
          </div>

          {result && (
            <div class={resultCard}>
              <div class={resultTitle}>
                <CheckCircle2 size={15} />
                <span>Ingested</span>
              </div>
              <span class={resultMeta}>
                {result.count ?? 0} event{result.count === 1 ? '' : 's'} · status {result.status}
              </span>
              {result.run_ids && result.run_ids.length > 0 ? (
                <div class={runLinkList}>
                  {result.run_ids.map((id) => (
                    <Link key={id} to="/runs" class={runLink}>
                      {id} <ArrowRight size={12} />
                    </Link>
                  ))}
                </div>
              ) : (
                <span class={noRunNote}>No agent workflow was triggered by this event.</span>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
