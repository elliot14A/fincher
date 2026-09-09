import { useQuery } from '@tanstack/react-query'
import { GitBranch, Workflow, X } from 'lucide-preact'
import { useEffect, useState } from 'preact/hooks'
import { runsQueryOptions } from '#/features/runs/queryOptions'
import type { ModelsRun, ModelsStep } from '#/lib/api'
import {
  backdrop,
  bar,
  barRight,
  barSubtitle,
  barTitle,
  barTitleGroup,
  closeBtn,
  loading,
  stage,
  toggleBtn,
  toggleBtnActive,
  toggleGroup,
} from './flowFullscreen.css'
import { RunFlow } from './runFlow'
import { TitleFlow } from './titleFlow'

export type FlowScope = 'run' | 'title'

export type FlowFullscreenProps = {
  run: ModelsRun
  titleName: string
  steps: ModelsStep[]
  selectedStepId: string | null
  onSelectStep: (stepId: string) => void
  onSelectRun: (runId: string) => void
  onClose: () => void
}

export function FlowFullscreen({
  run,
  titleName,
  steps,
  selectedStepId,
  onSelectStep,
  onSelectRun,
  onClose,
}: FlowFullscreenProps) {
  const [scope, setScope] = useState<FlowScope>('run')
  const slug = run.title_slug ?? ''
  const canGroup = Boolean(slug) && slug !== 'SYSTEM'

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [onClose])

  const { data: grouped, isLoading } = useQuery({
    ...runsQueryOptions({ title_slug: slug, limit: 100 }),
    enabled: scope === 'title' && canGroup,
  })

  const groupedRuns = grouped?.items ?? []

  return (
    <div class={backdrop} role="dialog" aria-modal="true" aria-label="Agent flow graph">
      <div class={bar}>
        <div class={barTitleGroup}>
          <span class={barTitle}>{titleName}</span>
          <span class={barSubtitle}>
            {scope === 'run' ? run.id : `${groupedRuns.length || '…'} runs · closed-loop lifecycle`}
          </span>
        </div>

        <div class={barRight}>
          {canGroup && (
            <div class={toggleGroup}>
              <button
                type="button"
                class={scope === 'run' ? `${toggleBtn} ${toggleBtnActive}` : toggleBtn}
                onClick={() => setScope('run')}
              >
                <Workflow size={13} /> This run
              </button>
              <button
                type="button"
                class={scope === 'title' ? `${toggleBtn} ${toggleBtnActive}` : toggleBtn}
                onClick={() => setScope('title')}
              >
                <GitBranch size={13} /> All runs for title
              </button>
            </div>
          )}
          <button type="button" class={closeBtn} onClick={onClose} aria-label="Close">
            <X size={16} />
          </button>
        </div>
      </div>

      <div class={stage}>
        {scope === 'run' ? (
          <RunFlow
            steps={steps}
            selectedStepId={selectedStepId}
            onSelectStep={onSelectStep}
            fullscreen
          />
        ) : isLoading ? (
          <div class={loading}>Loading title lifecycle…</div>
        ) : (
          <TitleFlow runs={groupedRuns} selectedRunId={run.id} onSelectRun={onSelectRun} />
        )}
      </div>
    </div>
  )
}
