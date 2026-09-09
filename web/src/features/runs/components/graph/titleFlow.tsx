import {
  Background,
  BackgroundVariant,
  Controls,
  type Edge,
  type Node,
  ReactFlow,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { useMemo, useState } from 'preact/hooks'
import type { ModelsRun } from '#/lib/api'
import { AgentNode, type AgentNodeData } from './agentNode'
import { canvas } from './runFlow.css'
import { RunNode, type RunNodeData } from './runNode'
import { formatLatency, narrativeFor, stepLatencyMs } from './stepNarrative'
import {
  canvasHost,
  legend,
  legendItem,
  legendSwatchAllocation,
  legendSwatchIncident,
  legendSwatchResolution,
  wrap,
} from './titleFlow.css'

const nodeTypes = { run: RunNode, agent: AgentNode }

const RUN_GAP = 300
const RUN_Y = 0
const STAGE_Y = 150
const STAGE_WIDTH = 188
const STAGE_GAP = 40

export type TitleFlowProps = {
  runs: ModelsRun[]
  selectedRunId: string | null
  onSelectRun: (runId: string) => void
}

function clock(ts?: string): string {
  if (!ts) return '—'
  const d = new Date(ts)
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

export function TitleFlow({ runs, selectedRunId, onSelectRun }: TitleFlowProps) {
  const [expanded, setExpanded] = useState<Record<string, boolean>>({})

  const sorted = useMemo(
    () =>
      [...runs].sort(
        (a, b) =>
          new Date(a.started_at ?? a.created_at ?? 0).getTime() -
          new Date(b.started_at ?? b.created_at ?? 0).getTime(),
      ),
    [runs],
  )

  const { nodes, edges } = useMemo(() => {
    const n: Node[] = []
    const e: Edge[] = []

    sorted.forEach((run, idx) => {
      const steps = run.steps ?? []
      const isExpanded = Boolean(expanded[run.id])
      const runX = idx * RUN_GAP
      const runData: RunNodeData = {
        trigger: run.trigger,
        status: run.status,
        stageCount: steps.length,
        latency: formatLatency(stepLatencyMs(run.started_at, run.ended_at)),
        clock: clock(run.started_at),
        expanded: isExpanded,
        onToggle: () => setExpanded((prev) => ({ ...prev, [run.id]: !prev[run.id] })),
      }

      n.push({
        id: run.id,
        type: 'run',
        position: { x: runX, y: RUN_Y },
        data: runData as unknown as Record<string, unknown>,
        selected: run.id === selectedRunId,
        draggable: false,
      })

      if (idx > 0) {
        e.push({
          id: `loop-${sorted[idx - 1].id}-${run.id}`,
          source: sorted[idx - 1].id,
          target: run.id,
          animated: run.status === 'RUNNING',
        })
      }

      if (isExpanded && steps.length > 0) {
        steps.forEach((step, sIdx) => {
          const narrative = narrativeFor(step.name)
          const stageId = `${run.id}::${step.id || sIdx}`
          const stageData: AgentNodeData = {
            title: narrative.title,
            agent: narrative.agent,
            icon: narrative.icon,
            latency: formatLatency(stepLatencyMs(step.started_at, step.ended_at)),
            index: sIdx,
            total: steps.length,
            status: step.status,
          }
          n.push({
            id: stageId,
            type: 'agent',
            position: { x: runX + sIdx * (STAGE_WIDTH + STAGE_GAP), y: STAGE_Y },
            data: stageData as unknown as Record<string, unknown>,
            draggable: false,
            selectable: false,
          })
          if (sIdx === 0) {
            e.push({ id: `s-${run.id}-${stageId}`, source: run.id, target: stageId })
          } else {
            const prevId = `${run.id}::${steps[sIdx - 1].id || sIdx - 1}`
            e.push({ id: `s-${prevId}-${stageId}`, source: prevId, target: stageId })
          }
        })
      }
    })

    return { nodes: n, edges: e }
  }, [sorted, expanded, selectedRunId])

  return (
    <div class={wrap}>
      <div class={legend}>
        <span class={legendItem}>
          <span class={legendSwatchAllocation} /> Allocation
        </span>
        <span class={legendItem}>
          <span class={legendSwatchIncident} /> Incident
        </span>
        <span class={legendItem}>
          <span class={legendSwatchResolution} /> Resolution
        </span>
      </div>
      <div class={canvasHost}>
        <div class={canvas}>
          <ReactFlow
            nodes={nodes}
            edges={edges}
            nodeTypes={nodeTypes}
            onNodeClick={(_e, node) => {
              if (node.type === 'run') onSelectRun(node.id)
            }}
            fitView
            fitViewOptions={{ padding: 0.15 }}
            nodesConnectable={false}
            nodesDraggable={false}
            elementsSelectable
            proOptions={{ hideAttribution: true }}
            minZoom={0.2}
            maxZoom={1.6}
          >
            <Background variant={BackgroundVariant.Dots} gap={16} size={1} />
            <Controls showInteractive={false} />
          </ReactFlow>
        </div>
      </div>
    </div>
  )
}
