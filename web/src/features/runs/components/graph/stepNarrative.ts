import type { LucideIcon } from 'lucide-preact'
import {
  Activity,
  Boxes,
  Database,
  GitFork,
  PackageCheck,
  Send,
  Users,
  Workflow,
} from 'lucide-preact'

export type IconType = LucideIcon

export type StepMeta = Record<string, unknown>

export type StepNarrative = {
  title: string
  agent: string
  icon: IconType
  describe: (meta: StepMeta) => string
}

export function str(v: unknown): string {
  return typeof v === 'string' ? v : ''
}

export function num(v: unknown): number {
  return typeof v === 'number' ? v : 0
}

export const STEP_NARRATIVES: Record<string, StepNarrative> = {
  triage_judge: {
    title: 'Triaged the incident',
    agent: 'Triage Judge',
    icon: Activity,
    describe: (m) => {
      const type = str(m.anomaly_type)
      const actionable = m.actionable === true
      if (type) {
        return actionable
          ? `Classified the event as ${type} and confirmed it is actionable.`
          : `Classified the event as ${type} — no action required.`
      }
      return 'Assessed the incoming event to decide whether it needs a response.'
    },
  },
  context_gathering: {
    title: 'Gathered live context',
    agent: 'Historian · ClickHouse MCP',
    icon: Database,
    describe: (m) => {
      const repairs = Array.isArray(m.repairs) ? m.repairs.length : 0
      const holds = num(m.deliveries_on_hold)
      return `Queried delivery, package and vendor state${
        repairs ? ` and mapped ${repairs} repair${repairs === 1 ? '' : 's'}` : ''
      }${holds ? ` with ${holds} deliveries on hold` : ''}.`
    },
  },
  remediation_loop: {
    title: 'Planned the remediation',
    agent: 'Action Planner + Policy Judge',
    icon: GitFork,
    describe: (m) => {
      const attempts = num(m.attempts)
      const decision = str(m.decision)
      if (decision) {
        return `Proposed a plan and ran the policy verification loop — ${decision}${
          attempts > 1 ? ` after ${attempts} attempts` : ''
        }.`
      }
      return 'Proposed an action plan and submitted it to the policy judge.'
    },
  },
  remediation_executor: {
    title: 'Executed the plan',
    agent: 'Executor',
    icon: Send,
    describe: (m) => {
      const count = num(m.artifacts_count)
      return count
        ? `Dispatched ${count} action${count === 1 ? '' : 's'} against production state.`
        : 'Applied the approved plan to production state.'
    },
  },
  candidate_gathering: {
    title: 'Gathered vendor candidates',
    agent: 'Historian · ClickHouse MCP',
    icon: Database,
    describe: (m) => {
      const reqs = num(m.requirements_count)
      return reqs
        ? `Resolved ${reqs} localization requirement${reqs === 1 ? '' : 's'} and pulled candidate vendors.`
        : 'Resolved localization requirements and pulled candidate vendors.'
    },
  },
  vendor_selection: {
    title: 'Selected vendors',
    agent: 'Vendor Selector',
    icon: Users,
    describe: (m) => {
      const count = num(m.assignments_count)
      return count
        ? `Matched ${count} task${count === 1 ? '' : 's'} to the best-fit vendor on turnaround and cost.`
        : 'Matched each task to the best-fit vendor on turnaround and cost.'
    },
  },
  provisioning: {
    title: 'Provisioned the work',
    agent: 'Executor',
    icon: Boxes,
    describe: (m) => {
      const d = num(m.deliveries_created)
      const p = num(m.packages_created)
      return `Created ${d} deliver${d === 1 ? 'y' : 'ies'} and ${p} package${
        p === 1 ? '' : 's'
      }, scheduling QC for each.`
    },
  },
  resolution_evaluation: {
    title: 'Evaluated resolution',
    agent: 'Resolution Engine',
    icon: PackageCheck,
    describe: (m) => {
      const status = str(m.title_status)
      const released = Array.isArray(m.released_deliveries) ? m.released_deliveries.length : 0
      if (released) {
        return `Verified clean QC and released ${released} deliver${released === 1 ? 'y' : 'ies'}${
          status ? ` — title now ${status}` : ''
        }.`
      }
      return status
        ? `Re-checked downstream state — title is ${status}.`
        : 'Re-checked downstream state after the last action.'
    },
  },
}

export function narrativeFor(name: string): StepNarrative {
  return (
    STEP_NARRATIVES[name] ?? {
      title: name.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase()),
      agent: 'Workflow node',
      icon: Workflow,
      describe: () => 'Executed a workflow node.',
    }
  )
}

export function stepLatencyMs(startedAt?: string, endedAt?: string): number {
  if (!startedAt) return 0
  const start = new Date(startedAt).getTime()
  const end = endedAt ? new Date(endedAt).getTime() : Date.now()
  return Math.max(0, end - start)
}

export function formatLatency(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

export function isApprovedOutcome(outcome: string): boolean {
  const o = outcome.toUpperCase()
  return (
    o.includes('APPROVE') ||
    o.includes('RESOLVE') ||
    o.includes('RELEASE') ||
    o.includes('PASS') ||
    o.includes('ACTIONABLE') ||
    o.startsWith('VND-')
  )
}

export function isRejectedOutcome(outcome: string): boolean {
  const o = outcome.toUpperCase()
  return o.includes('REJECT') || o.includes('FAIL') || o.includes('ESCALATE')
}
