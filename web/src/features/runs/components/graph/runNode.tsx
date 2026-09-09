import { Handle, type NodeProps, Position } from '@xyflow/react'
import { ChevronDown, ChevronRight, Layers } from 'lucide-preact'
import {
  accentAllocation,
  accentDefault,
  accentIncident,
  accentResolution,
  expandBtn,
  headerRow,
  metaRow,
  metaText,
  node,
  nodeSelected,
  statusDotDanger,
  statusDotNeutral,
  statusDotSuccess,
  timeText,
  triggerLabel,
} from './runNode.css'

export type RunNodeData = {
  trigger: string
  status: string
  stageCount: number
  latency: string
  clock: string
  expanded: boolean
  onToggle: () => void
}

function accentFor(trigger: string): string {
  switch (trigger.toLowerCase()) {
    case 'incident':
      return accentIncident
    case 'allocation':
      return accentAllocation
    case 'resolution':
      return accentResolution
    default:
      return accentDefault
  }
}

export function RunNode({ data, selected }: NodeProps) {
  const d = data as RunNodeData
  const isFailed = d.status === 'FAILED' || d.status === 'ESCALATED'
  const isDone = d.status === 'COMPLETED'
  const dotClass = isFailed ? statusDotDanger : isDone ? statusDotSuccess : statusDotNeutral

  const nodeClass = [node, accentFor(d.trigger), selected ? nodeSelected : '']
    .filter(Boolean)
    .join(' ')

  return (
    <div class={nodeClass}>
      <Handle type="target" position={Position.Left} isConnectable={false} />
      <div class={headerRow}>
        <span class={triggerLabel}>
          <span class={dotClass} />
          {d.trigger}
        </span>
        <span class={timeText}>{d.clock}</span>
      </div>
      <div class={metaRow}>
        <span class={metaText}>
          {d.stageCount} stage{d.stageCount === 1 ? '' : 's'} · {d.latency}
        </span>
        <button
          type="button"
          class={expandBtn}
          onClick={(e) => {
            e.stopPropagation()
            d.onToggle()
          }}
          aria-label={d.expanded ? 'Collapse stages' : 'Expand stages'}
        >
          {d.expanded ? <ChevronDown size={11} /> : <ChevronRight size={11} />}
          <Layers size={11} />
        </button>
      </div>
      <Handle type="source" position={Position.Right} isConnectable={false} />
    </div>
  )
}
