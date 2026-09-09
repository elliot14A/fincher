import {
  Background,
  BackgroundVariant,
  Controls,
  type Edge,
  type Node,
  ReactFlow,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { Maximize2 } from 'lucide-preact'
import { useMemo } from 'preact/hooks'
import type { ModelsStep } from '#/lib/api'
import { AgentNode, type AgentNodeData } from './agentNode'
import { canvas, flowShell, flowShellFull, maximizeBtn } from './runFlow.css'
import { formatLatency, narrativeFor, stepLatencyMs } from './stepNarrative'

const nodeTypes = { agent: AgentNode }

const NODE_WIDTH = 188
const X_GAP = 72
const Y_POS = 60

export type RunFlowProps = {
  steps: ModelsStep[]
  selectedStepId: string | null
  onSelectStep: (stepId: string) => void
  fullscreen?: boolean
  onMaximize?: () => void
}

export function RunFlow({
  steps,
  selectedStepId,
  onSelectStep,
  fullscreen = false,
  onMaximize,
}: RunFlowProps) {
  const { nodes, edges } = useMemo(() => {
    const n: Node[] = steps.map((step, idx) => {
      const narrative = narrativeFor(step.name)
      const data: AgentNodeData = {
        title: narrative.title,
        agent: narrative.agent,
        icon: narrative.icon,
        latency: formatLatency(stepLatencyMs(step.started_at, step.ended_at)),
        index: idx,
        total: steps.length,
        status: step.status,
      }
      return {
        id: step.id || `step-${idx}`,
        type: 'agent',
        position: { x: idx * (NODE_WIDTH + X_GAP), y: Y_POS },
        data: data as unknown as Record<string, unknown>,
        selected: (step.id || `step-${idx}`) === selectedStepId,
        draggable: false,
      }
    })

    const e: Edge[] = []
    for (let i = 0; i < n.length - 1; i++) {
      e.push({
        id: `e-${n[i].id}-${n[i + 1].id}`,
        source: n[i].id,
        target: n[i + 1].id,
        animated: steps[i + 1]?.status === 'RUNNING',
      })
    }
    return { nodes: n, edges: e }
  }, [steps, selectedStepId])

  return (
    <div class={fullscreen ? `${flowShell} ${flowShellFull}` : flowShell}>
      {onMaximize && (
        <button type="button" class={maximizeBtn} onClick={onMaximize} aria-label="Maximize graph">
          <Maximize2 size={12} />
          <span>Expand</span>
        </button>
      )}
      <div class={canvas}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          onNodeClick={(_e, node) => onSelectStep(node.id)}
          fitView
          fitViewOptions={{ padding: 0.2 }}
          nodesConnectable={false}
          nodesDraggable={false}
          elementsSelectable
          panOnScroll={fullscreen}
          zoomOnScroll={fullscreen}
          preventScrolling={fullscreen}
          proOptions={{ hideAttribution: true }}
          minZoom={0.3}
          maxZoom={1.6}
        >
          <Background variant={BackgroundVariant.Dots} gap={16} size={1} />
          <Controls showInteractive={false} />
        </ReactFlow>
      </div>
    </div>
  )
}
