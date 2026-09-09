import {
  AlertTriangle,
  FileWarning,
  GitCommitVertical,
  type LucideIcon,
  Wrench,
} from 'lucide-preact'
import type { ModelsEvent, ModelsPackage } from '#/lib/api'

export type ScenarioId = 'master_revision' | 'package_invalidation' | 'qc_failed' | 'custom'

export type ScenarioTrigger = 'incident' | 'allocation' | 'resolution' | 'custom'

export type ScenarioDef = {
  id: ScenarioId
  label: string
  icon: LucideIcon
  trigger: ScenarioTrigger
  eventType: string
  defaultSeverity: ModelsEvent['severity']
  description: string
  needsTitle: boolean
  needsPackage: boolean
}

export const SCENARIOS: ScenarioDef[] = [
  {
    id: 'master_revision',
    label: 'New master release',
    icon: GitCommitVertical,
    trigger: 'incident',
    eventType: 'fincher.master.cut.revised',
    defaultSeverity: 'CRITICAL',
    description:
      'Editorial ships a revised master. The agent invalidates stale packages, holds affected deliveries, and re-plans QC.',
    needsTitle: true,
    needsPackage: false,
  },
  {
    id: 'package_invalidation',
    label: 'Package invalidated',
    icon: FileWarning,
    trigger: 'incident',
    eventType: 'fincher.package.invalidated',
    defaultSeverity: 'CRITICAL',
    description:
      'A localization package is flagged bad. The agent triages the blast radius and plans remediation under policy review.',
    needsTitle: true,
    needsPackage: true,
  },
  {
    id: 'qc_failed',
    label: 'QC inspection failed',
    icon: AlertTriangle,
    trigger: 'incident',
    eventType: 'fincher.qc.completed',
    defaultSeverity: 'WARN',
    description:
      'An automated QC pass fails on a package. The agent investigates the vendor and decides on reconform or reassignment.',
    needsTitle: true,
    needsPackage: true,
  },
  {
    id: 'custom',
    label: 'Custom event',
    icon: Wrench,
    trigger: 'custom',
    eventType: 'fincher.package.invalidated',
    defaultSeverity: 'INFO',
    description: 'Hand-craft any event payload and emit it directly to the ingestion API.',
    needsTitle: false,
    needsPackage: false,
  },
]

export function scenarioById(id: ScenarioId): ScenarioDef {
  return SCENARIOS.find((s) => s.id === id) ?? SCENARIOS[0]
}

export type ScenarioInputs = {
  subject: string
  pkg?: ModelsPackage
  newMasterVersion: string
  invalidationReason: string
  driftMs: number
  defectCategory: string
}

export function buildEventData(
  scenario: ScenarioDef,
  inputs: ScenarioInputs,
): Record<string, unknown> {
  switch (scenario.id) {
    case 'master_revision':
      return {
        new_master_version: inputs.newMasterVersion,
      }
    case 'package_invalidation':
      return {
        package_id: inputs.pkg?.id ?? '',
        reason: inputs.invalidationReason,
      }
    case 'qc_failed':
      return {
        package_id: inputs.pkg?.id ?? '',
        vendor_id: inputs.pkg?.vendor_id ?? '',
        component: inputs.pkg?.component ?? 'AUDIO',
        language: inputs.pkg?.language ?? '',
        status: 'FAILED',
        sync_drift_ms: inputs.driftMs,
        defect_category: inputs.defectCategory,
      }
    default:
      return {}
  }
}

export function buildEvent(
  scenario: ScenarioDef,
  severity: ModelsEvent['severity'],
  inputs: ScenarioInputs,
): ModelsEvent {
  return {
    type: scenario.eventType,
    source: 'fincher.ops.console',
    subject: inputs.subject || 'GLOBAL',
    severity,
    data: buildEventData(scenario, inputs),
  }
}

export const DEFECT_CATEGORIES = [
  'AUDIO_SYNC_DRIFT',
  'CORRUPT_FRAME',
  'SUBTITLE_OVERLAP',
  'OTHER',
] as const
