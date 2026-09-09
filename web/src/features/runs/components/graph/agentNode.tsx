import { Handle, type NodeProps, Position } from '@xyflow/react'
import { CheckCircle2, Clock, XCircle } from 'lucide-preact'
import {
  agent,
  footerRow,
  iconBox,
  iconBoxDanger,
  iconBoxRunning,
  iconBoxSuccess,
  latency,
  node,
  nodeDanger,
  nodeHeader,
  nodeRunning,
  nodeSelected,
  nodeSuccess,
  stageIndex,
  title,
  titleGroup,
} from './agentNode.css'
import type { IconType } from './stepNarrative'

export type AgentNodeData = {
  title: string
  agent: string
  icon: IconType
  latency: string
  index: number
  total: number
  status: 'COMPLETED' | 'FAILED' | 'RUNNING' | 'PENDING' | 'SKIPPED'
}

export function AgentNode({ data, selected }: NodeProps) {
  const d = data as AgentNodeData
  const Icon = d.icon
  const isCompleted = d.status === 'COMPLETED'
  const isFailed = d.status === 'FAILED'
  const isRunning = d.status === 'RUNNING'

  const nodeClass = [
    node,
    isFailed ? nodeDanger : isRunning ? nodeRunning : isCompleted ? nodeSuccess : '',
    selected ? nodeSelected : '',
  ]
    .filter(Boolean)
    .join(' ')

  const iconClass = [
    iconBox,
    isFailed ? iconBoxDanger : isRunning ? iconBoxRunning : isCompleted ? iconBoxSuccess : '',
  ]
    .filter(Boolean)
    .join(' ')

  return (
    <div class={nodeClass}>
      <Handle type="target" position={Position.Left} isConnectable={false} />
      <div class={nodeHeader}>
        <div class={iconClass}>
          <Icon size={15} />
        </div>
        <div class={titleGroup}>
          <span class={title}>{d.title}</span>
          <span class={agent}>{d.agent}</span>
        </div>
      </div>
      <div class={footerRow}>
        <span class={stageIndex}>
          Stage {d.index + 1}/{d.total}
        </span>
        <span class={latency}>
          {isFailed ? (
            <XCircle size={11} />
          ) : isRunning ? (
            <Clock size={11} />
          ) : isCompleted ? (
            <CheckCircle2 size={11} />
          ) : null}{' '}
          {d.latency}
        </span>
      </div>
      <Handle type="source" position={Position.Right} isConnectable={false} />
    </div>
  )
}
