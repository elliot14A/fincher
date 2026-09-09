import { createFileRoute, Link } from '@tanstack/react-router'
import {
  Activity,
  AlertOctagon,
  ArrowRight,
  Bot,
  Check,
  CheckCircle2,
  ChevronRight,
  Cpu,
  Database,
  Film,
  GitFork,
  HelpCircle,
  Layers,
  Lock,
  Play,
  RefreshCw,
  Scale,
  ShieldAlert,
  ShieldCheck,
  Sparkles,
  Terminal,
  Users,
  Workflow,
  Zap,
} from 'lucide-preact'
import { useState } from 'preact/hooks'
import { Badge } from '#/components/ui/badge'
import { Logo } from '#/components/ui/logo'
import {
  archCodePre,
  archConnectorCol,
  archConnectorText,
  archDot,
  archEventBadge,
  archEventBox,
  archGrid,
  archHeaderLeft,
  architectureCard,
  architectureCardHeader,
  archLivePill,
  archNodeCard,
  archNodeIcon,
  archNodeLeft,
  archNodeList,
  archNodeSubtext,
  archNodeText,
  archTitle,
  chatAiAnswerText,
  chatAiBubble,
  chatCitationList,
  chatCitationPill,
  chatCitationsLabel,
  chatCitationsWrapper,
  chatFooterActionRow,
  chatFooterNote,
  chatGoToChatLink,
  chatPromptButton,
  chatPromptButtonActive,
  chatPromptIcon,
  chatPromptLeft,
  chatPromptsCol,
  chatPromptsHeader,
  chatPromptsSubtitle,
  chatPromptsTitle,
  chatPromptTag,
  chatPromptText,
  chatShowcaseContainer,
  chatUserBubble,
  chatWindowAvatar,
  chatWindowBody,
  chatWindowCard,
  chatWindowHeader,
  chatWindowHeaderLeft,
  chatWindowStatus,
  chatWindowTitle,
  demoBody,
  demoContainer,
  demoNav,
  demoScorecardCol,
  demoScorecardText,
  demoScorecardTitle,
  demoScorecardVerdict,
  demoStepBadge,
  demoStepLatency,
  demoStepLeft,
  demoStepName,
  demoStepRow,
  demoTabBtn,
  demoTabBtnActive,
  demoWaterfallCol,
  demoWaterfallHeader,
  dotGreen,
  dotRed,
  dotYellow,
  executedActionIcon,
  executedActionPill,
  executedActionsTitle,
  executedActionsWrapper,
  featureCard,
  featureGrid,
  featureIconBox,
  featurePill,
  featurePillList,
  featureText,
  featureTitle,
  feedbackBadge,
  finalCtaButtons,
  finalCtaCard,
  finalCtaSubtitle,
  finalCtaTitle,
  heroActions,
  heroBadge,
  heroBadgeDot,
  heroBadgePulse,
  heroContainer,
  heroEyebrow,
  heroGlow,
  heroGridTexture,
  heroSection,
  heroStepArrow,
  heroStepBody,
  heroStepIcon,
  heroStepItem,
  heroStepLabel,
  heroSteps,
  heroStepText,
  heroSubtitle,
  heroTitle,
  heroTitleGradient,
  landingContainer,
  loopCard,
  loopCardActive,
  loopGrid,
  loopHeader,
  loopSection,
  loopStepDesc,
  loopStepNum,
  loopStepTitle,
  metricCard,
  metricLabel,
  metricSubtext,
  metricsGrid,
  metricValue,
  navCtaBtn,
  navHeader,
  navInner,
  navLinkItem,
  navLinks,
  navLogoArea,
  navRight,
  primaryCtaBtn,
  problemCard,
  problemGrid,
  problemIconBox,
  problemText,
  problemTitle,
  safetyCard,
  safetyDesc,
  safetyGrid,
  safetyIconBox,
  safetyTitle,
  secondaryCtaBtn,
  sectionBadge,
  sectionContainer,
  sectionHeader,
  sectionHeading,
  sectionSubheading,
  stepCheckIcon,
  stepIndexText,
  stepRightBox,
  techPill,
  techPillIcon,
  techPillList,
  techStrip,
  techStripLabel,
} from '#/styles/routes/landing.css'

export const Route = createFileRoute('/')({
  component: LandingPage,
})

type DemoScenario = 'drift' | 'allocation' | 'master'

type ChatPromptKey = 'hold_reason' | 'vendor_score' | 'german_status' | 'master_revision'

interface ChatSample {
  key: ChatPromptKey
  prompt: string
  tag: string
  answer: string
  citations: Array<{ label: string; type: 'sqlite' | 'clickhouse' }>
}

const CHAT_SAMPLES: Record<ChatPromptKey, ChatSample> = {
  hold_reason: {
    key: 'hold_reason',
    prompt: 'Why is Avatar: Fire and Ash on HOLD, and what is the blast radius?',
    tag: 'Incident Triage',
    answer:
      'Avatar: Fire and Ash was placed on HOLD due to a critical audio synchronization drift (+42.5ms) detected in package pkg-de-audio-v01 during ingest QC. Lineage analysis determined that the blast radius is strictly confined to the German (DE) and synchronized European theatrical deliveries. North American and Asian releases remain unaffected and on schedule. Deluxe Digital has been assigned to re-conform the audio stems.',
    citations: [
      { label: 'SQLite: titles.status = "HOLD" · blast_radius = 2', type: 'sqlite' },
      {
        label: 'ClickHouse: qc_inspections (drift_ms = 42.5, severity = "CRITICAL")',
        type: 'clickhouse',
      },
      { label: 'ADK Lineage: 2 dependent packages contained, 6 safe', type: 'sqlite' },
    ],
  },
  vendor_score: {
    key: 'vendor_score',
    prompt: 'Which vendor has the highest accuracy for Japanese dubbing and QC?',
    tag: 'Vendor Intelligence',
    answer:
      'According to ClickHouse historical QC telemetry across 2,418 completed runs over the last 120 days (with exponential half-life decay), Iyuno holds the highest performance rating at 98.8% first-pass accuracy with an average turnaround time of 14.2 hours. Pixelogic follows at 96.4% accuracy (18.6h TAT), while Deluxe Digital scores 94.1% on Japanese localization.',
    citations: [
      {
        label: 'ClickHouse: vendor_metrics (facility = "Iyuno", accuracy = 0.988, runs = 2418)',
        type: 'clickhouse',
      },
      { label: 'ClickHouse: 120-day half-life decay model applied', type: 'clickhouse' },
      { label: 'SQLite: active_contracts (Iyuno SLA = 16h)', type: 'sqlite' },
    ],
  },
  german_status: {
    key: 'german_status',
    prompt: 'Are German dubs on track for the Dune: Part Two premiere window?',
    tag: 'Release Readiness',
    answer:
      'Yes. The re-mastered German 7.1 Dolby Atmos stem (pkg-de-dub-v02) successfully passed automated audio sync validation with zero drift (0.2ms vs 15ms threshold). The Policy Verification Judge cleared all safety checks, and the delivery status has been restored to READY_TO_SHIP with a +52 hour buffer prior to the global storefront premiere.',
    citations: [
      { label: 'SQLite: deliveries.status = "READY_TO_SHIP" · buffer = "+52h"', type: 'sqlite' },
      {
        label: 'ClickHouse: qc_inspections (sync_drift_ms = 0.2, verdict = "PASS")',
        type: 'clickhouse',
      },
      { label: 'Policy Judge: verification attempt 1/3 PASSED', type: 'sqlite' },
    ],
  },
  master_revision: {
    key: 'master_revision',
    prompt: 'What happened when the Wicked master cut bumped from V01 to V02?',
    tag: 'Lineage & Containment',
    answer:
      'When event fincher.master.bumped was received, Fincher’s stale-pass guard triggered immediately. The lineage engine traversed the component graph and identified 4 downstream localized packages (Spanish Subtitles, French Dub, German Dub, Italian Subtitles) rendered against V01. These 4 packages were automatically invalidated and re-queued for conform without halting the entire title catalog.',
    citations: [
      { label: 'SQLite: lineage_edges (invalidated_count = 4, unaffected = 12)', type: 'sqlite' },
      { label: 'ClickHouse: system_events (type = "fincher.master.bumped")', type: 'clickhouse' },
      { label: 'Software Executor: cancelled 2 stale jobs, scheduled 4 re-QCs', type: 'sqlite' },
    ],
  },
}

const DEMO_SCENARIOS = {
  drift: {
    id: 'drift',
    title: 'Audio Sync Drift Remediation',
    trigger: 'Anomaly Signal (fincher.audio.sync_drift)',
    steps: [
      { name: 'ANOMALY_TRIAGE', label: 'Triage Engine', latency: '42ms' },
      { name: 'HISTORIAN_CH_QUERY', label: 'ClickHouse MCP', latency: '18ms' },
      { name: 'DEPENDENCY_GRAPH_WALK', label: 'Lineage Traversal', latency: '6ms' },
      { name: 'POLICY_VERIFICATION', label: 'Policy Judge', latency: '840ms' },
      { name: 'EXECUTOR_DISPATCH', label: 'Software Executor', latency: '12ms' },
    ],
    judge: 'Policy Verification Judge (Attempt 1 of 3)',
    verdict: 'APPROVED · SAFE_TO_EXECUTE',
    rationale:
      'Remediation plan places German dub delivery on temporary HOLD and reassigns to Deluxe Digital. Deluxe historical accuracy is 99.2% with 12h turnaround. Projected buffer remains +48h ahead of premiere window.',
    actions: [
      'HOLD_DELIVERY: US & DE Releases',
      'REASSIGN_VENDOR: Deluxe Digital',
      'DISPATCH: Vendor SLA Notice',
    ],
  },
  allocation: {
    id: 'allocation',
    title: 'Multi-Market Vendor Allocation',
    trigger: 'Allocation Request (fincher.title.onboarded)',
    steps: [
      { name: 'REQUIREMENTS_PARSER', label: 'ADK Router', latency: '35ms' },
      { name: 'HISTORIAN_CH_QUERY', label: 'ClickHouse MCP', latency: '22ms' },
      { name: 'VENDOR_SELECTION', label: 'Allocation Judge', latency: '790ms' },
      { name: 'SCHEDULE_DISPATCH', label: 'Software Executor', latency: '9ms' },
    ],
    judge: 'Vendor Allocation Judge',
    verdict: 'ALLOCATED · OPTIMAL_COST_AND_TIME',
    rationale:
      'Evaluated 4 candidate QC facilities across 5 language markets. Iyuno selected for JP/ES dubbing ($120/hr, 16h TAT), Pixelogic selected for Subtitles ($85/hr, 8h TAT). All delivery SLA targets satisfied.',
    actions: [
      'ASSIGN_PACKAGE: PKG-AUDIO-JA',
      'ASSIGN_PACKAGE: PKG-SUB-ES',
      'DISPATCH: Work Orders Emitted',
    ],
  },
  master: {
    id: 'master',
    title: 'Upstream Master Revision Handoff',
    trigger: 'Master Cut Bump (fincher.master.bumped)',
    steps: [
      { name: 'LINEAGE_INVALIDATION', label: 'Lineage Engine', latency: '8ms' },
      { name: 'BLAST_RADIUS_EVAL', label: 'Graph Traversal', latency: '14ms' },
      { name: 'POLICY_VERIFICATION', label: 'Policy Judge', latency: '620ms' },
      { name: 'RE_QC_DISPATCH', label: 'Software Executor', latency: '11ms' },
    ],
    judge: 'Master Revision Safety Judge',
    verdict: 'CONTAINED · RE_QC_DISPATCHED',
    rationale:
      'Editorial Master cut revised from V01 to V02. Stale-pass guard activated; invalidated 6 downstream localized audio/subtitle packages without full catalog halting. Automated re-QC scheduled.',
    actions: [
      'INVALIDATE: 6 Dependent Packages',
      'CANCEL_TASKS: Stale In-Flight Runs',
      'SCHEDULE: Master Re-Conform',
    ],
  },
}

function LandingPage() {
  const [activeScenario, setActiveScenario] = useState<DemoScenario>('drift')
  const [activeChatKey, setActiveChatKey] = useState<ChatPromptKey>('hold_reason')
  const currentDemo = DEMO_SCENARIOS[activeScenario]
  const currentChat = CHAT_SAMPLES[activeChatKey]

  return (
    <div class={landingContainer}>
      {/* Sticky Top Navigation */}
      <header class={navHeader}>
        <div class={navInner}>
          <Link to="/" class={navLogoArea}>
            <Logo size="md" />
          </Link>

          <nav class={navLinks} aria-label="Main Navigation">
            <a href="#features" class={navLinkItem}>
              Capabilities
            </a>
            <a href="#feedback-loop" class={navLinkItem}>
              Closed Loop
            </a>
            <a href="#intelligence" class={navLinkItem}>
              AI Copilot
            </a>
            <a href="#architecture" class={navLinkItem}>
              Architecture
            </a>
            <a href="#playground" class={navLinkItem}>
              Workflows
            </a>
          </nav>

          <div class={navRight}>
            <Link to="/titles" class={navCtaBtn}>
              <span>Try Fincher Console</span>
              <ArrowRight size={14} />
            </Link>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section class={heroSection}>
        <div class={heroGridTexture} />
        <div class={heroGlow} />

        <div class={heroContainer}>
          <div class={heroBadge}>
            <span class={heroBadgeDot}>
              <span class={heroBadgePulse} />
            </span>
            <span>Autonomous Media Supply Chain Engine</span>
          </div>

          <p class={heroEyebrow}>An AI operations engineer for global film &amp; TV delivery</p>

          <h1 class={heroTitle}>
            Ship every localized title <br />
            <span class={heroTitleGradient}>on time, with zero defects.</span>
          </h1>

          <p class={heroSubtitle}>
            Fincher watches every audio, video, and subtitle package as it moves toward premiere.
            The moment QC drift or a bad master appears, its AI agents diagnose the blast radius,
            re-assign the best vendor, and heal the release &mdash; without halting the catalog or
            waiting on a human.
          </p>

          <div class={heroActions}>
            <Link to="/titles" class={primaryCtaBtn}>
              <span>Open Fincher Console</span>
              <ArrowRight size={15} />
            </Link>
            <Link to="/runs" class={secondaryCtaBtn}>
              <Play size={14} />
              <span>Watch Execution Runs</span>
            </Link>
          </div>

          {/* How it works, at a glance */}
          <div class={heroSteps}>
            <div class={heroStepItem}>
              <span class={heroStepIcon}>
                <Activity size={16} />
              </span>
              <span class={heroStepBody}>
                <span class={heroStepLabel}>1 · Detect</span>
                <span class={heroStepText}>
                  Ingests live QC events &amp; spots defects the instant they land
                </span>
              </span>
            </div>
            <span class={heroStepArrow}>
              <ArrowRight size={16} />
            </span>
            <div class={heroStepItem}>
              <span class={heroStepIcon}>
                <Bot size={16} />
              </span>
              <span class={heroStepBody}>
                <span class={heroStepLabel}>2 · Reason</span>
                <span class={heroStepText}>
                  ADK agents query ClickHouse history &amp; propose a verified fix
                </span>
              </span>
            </div>
            <span class={heroStepArrow}>
              <ArrowRight size={16} />
            </span>
            <div class={heroStepItem}>
              <span class={heroStepIcon}>
                <RefreshCw size={16} />
              </span>
              <span class={heroStepBody}>
                <span class={heroStepLabel}>3 · Self-heal</span>
                <span class={heroStepText}>
                  Executes safely, then re-checks itself until the title ships
                </span>
              </span>
            </div>
          </div>

          {/* Architecture Visual Graphic */}
          <div id="architecture" class={architectureCard}>
            <div class={architectureCardHeader}>
              <div class={archHeaderLeft}>
                <span class={`${archDot} ${dotRed}`} />
                <span class={`${archDot} ${dotYellow}`} />
                <span class={`${archDot} ${dotGreen}`} />
                <span class={archTitle}>event_pipeline_orchestrator.json</span>
              </div>
              <span class={archLivePill}>
                <Activity size={10} />
                <span>Live Event Stream</span>
              </span>
            </div>

            <div class={archGrid}>
              <div class={archEventBox}>
                <div class={archEventBadge}>
                  <ShieldAlert size={12} />
                  <span>Incoming CloudEvent (Ingested)</span>
                </div>
                <pre class={archCodePre}>
                  {`{
  "specversion": "1.0",
  "type": "fincher.audio.sync_drift",
  "subject": "avatar-fire-ash",
  "data": {
    "package_id": "pkg-de-audio-v01",
    "language": "de-DE",
    "drift_ms": 42.5,
    "severity": "CRITICAL_DRIFT"
  }
}`}
                </pre>
              </div>

              <div class={archConnectorCol}>
                <Workflow size={20} />
                <span class={archConnectorText}>ADK Graph</span>
                <ArrowRight size={18} />
              </div>

              <div class={archNodeList}>
                <div class={archNodeCard}>
                  <div class={archNodeLeft}>
                    <ShieldAlert size={14} class={archNodeIcon} />
                    <div>
                      <div class={archNodeText}>1. Delivery Halt &amp; Containment</div>
                      <div class={archNodeSubtext}>
                        Market US/DE blocked · 0 uncontained releases
                      </div>
                    </div>
                  </div>
                  <Badge variant="danger">HOLD</Badge>
                </div>

                <div class={archNodeCard}>
                  <div class={archNodeLeft}>
                    <Database size={14} class={archNodeIcon} />
                    <div>
                      <div class={archNodeText}>2. ClickHouse Vendor Scoring</div>
                      <div class={archNodeSubtext}>
                        Deluxe Digital: 99.2% accuracy · 12h turnaround
                      </div>
                    </div>
                  </div>
                  <Badge variant="warning">REASSIGN</Badge>
                </div>

                <div class={archNodeCard}>
                  <div class={archNodeLeft}>
                    <Scale size={14} class={archNodeIcon} />
                    <div>
                      <div class={archNodeText}>3. Policy Verification Gate</div>
                      <div class={archNodeSubtext}>
                        Attempt 1 of 3 · Premiere buffer verified (+48h)
                      </div>
                    </div>
                  </div>
                  <Badge variant="success">APPROVED</Badge>
                </div>

                <div class={archNodeCard}>
                  <div class={archNodeLeft}>
                    <CheckCircle2 size={14} class={archNodeIcon} />
                    <div>
                      <div class={archNodeText}>4. Closed-Loop Title Self-Healing</div>
                      <div class={archNodeSubtext}>
                        Clean QC return unholds delivery · Ready to ship
                      </div>
                    </div>
                  </div>
                  <Badge variant="success">ON_TRACK</Badge>
                </div>
              </div>
            </div>
          </div>

          {/* Tech Ecosystem Strip */}
          <div class={techStrip}>
            <span class={techStripLabel}>Engineered on First Principles With</span>
            <div class={techPillList}>
              <span class={techPill}>
                <Bot size={13} class={techPillIcon} />
                <span>Google ADK Go v2</span>
              </span>
              <span class={techPill}>
                <Sparkles size={13} class={techPillIcon} />
                <span>Gemini Structured Reasoning</span>
              </span>
              <span class={techPill}>
                <Database size={13} class={techPillIcon} />
                <span>ClickHouse Official MCP</span>
              </span>
              <span class={techPill}>
                <Cpu size={13} class={techPillIcon} />
                <span>Turso / SQLite WAL</span>
              </span>
              <span class={techPill}>
                <Layers size={13} class={techPillIcon} />
                <span>Preact + Vanilla Extract</span>
              </span>
            </div>
          </div>
        </div>
      </section>

      {/* The Problem Section */}
      <section class={sectionContainer}>
        <div class={sectionHeader}>
          <div class={sectionBadge}>
            <AlertOctagon size={13} />
            <span>The Media Supply Chain Dilemma</span>
          </div>
          <h2 class={sectionHeading}>Traditional media operations fail in silence.</h2>
          <p class={sectionSubheading}>
            Global synchronized releases across 50+ territories require zero-defect packages. When
            issues arise, manual coordination causes missed premiere windows and broadcast halts.
          </p>
        </div>

        <div class={problemGrid}>
          <div class={problemCard}>
            <div class={problemIconBox}>
              <ShieldAlert size={20} />
            </div>
            <h3 class={problemTitle}>Silent Multi-Market Drift</h3>
            <p class={problemText}>
              Out-of-sync audio dubs and corrupted subtitle cues are often discovered hours before
              global storefront premieres, risking multi-million dollar release delays.
            </p>
          </div>

          <div class={problemCard}>
            <div class={problemIconBox}>
              <Cpu size={20} />
            </div>
            <h3 class={problemTitle}>Manual Spreadsheet Chaos</h3>
            <p class={problemText}>
              Operations coordinators juggle dozens of localization facilities across static sheets
              without analytical historical track records or real-time performance telemetry.
            </p>
          </div>

          <div class={problemCard}>
            <div class={problemIconBox}>
              <AlertOctagon size={20} />
            </div>
            <h3 class={problemTitle}>Uncontained Blast Radiuses</h3>
            <p class={problemText}>
              A single defective master version triggers manual panic, halting entire catalog
              releases rather than surgically isolating the specific affected language packages.
            </p>
          </div>
        </div>
      </section>

      {/* Core Capabilities Section */}
      <section id="features" class={sectionContainer}>
        <div class={sectionHeader}>
          <div class={sectionBadge}>
            <Sparkles size={13} />
            <span>Engineered for Studio Precision</span>
          </div>
          <h2 class={sectionHeading}>Autonomous quality control from ingest to premiere.</h2>
          <p class={sectionSubheading}>
            Fincher connects ClickHouse analytical evidence with Google ADK multi-agent reasoning to
            remediate quality incidents in seconds.
          </p>
        </div>

        <div class={featureGrid}>
          <div class={featureCard}>
            <div class={featureIconBox}>
              <Database size={22} />
            </div>
            <h3 class={featureTitle}>ClickHouse Analytical Intelligence</h3>
            <p class={featureText}>
              Query millions of historical QC inspections in sub-milliseconds via ClickHouse MCP.
              Fincher applies a 120-day recency decay model to calculate true facility accuracy and
              turnaround reliability.
            </p>
            <div class={featurePillList}>
              <span class={featurePill}>ClickHouse MCP</span>
              <span class={featurePill}>120-Day Half-Life Decay</span>
              <span class={featurePill}>Sub-Millisecond Querying</span>
            </div>
          </div>

          <div class={featureCard}>
            <div class={featureIconBox}>
              <GitFork size={22} />
            </div>
            <h3 class={featureTitle}>Multi-Market Lineage Traversal</h3>
            <p class={featureText}>
              Parallel Historian and Dependency agents map deep component graphs across audio,
              video, and subtitle tracks to compute the exact blast radius without halting
              unaffected storefronts.
            </p>
            <div class={featurePillList}>
              <span class={featurePill}>Parallel ADK Agents</span>
              <span class={featurePill}>Surgical Blast Radius</span>
              <span class={featurePill}>1:N Delivery Resolution</span>
            </div>
          </div>

          <div class={featureCard}>
            <div class={featureIconBox}>
              <Scale size={22} />
            </div>
            <h3 class={featureTitle}>Bounded Policy Verification Gate</h3>
            <p class={featureText}>
              AI agents never mutate production state directly. Structured action proposals pass
              through a multi-judge verification loop evaluating premiere countdown urgency,
              contractual SLAs, and budget feasibility.
            </p>
            <div class={featurePillList}>
              <span class={featurePill}>Zero Unchecked AI</span>
              <span class={featurePill}>Capped 3-Attempt Loop</span>
              <span class={featurePill}>Transactional Go Executor</span>
            </div>
          </div>

          <div class={featureCard}>
            <div class={featureIconBox}>
              <RefreshCw size={22} />
            </div>
            <h3 class={featureTitle}>Closed-Loop Title Self-Healing</h3>
            <p class={featureText}>
              Fincher listens to its own downstream execution events. Clean vendor re-deliveries
              automatically unhold dependent packages, release storefront handoffs, and restore
              Title status from HOLD back to ON_TRACK.
            </p>
            <div class={featurePillList}>
              <span class={featurePill}>Zero-Touch Unholding</span>
              <span class={featurePill}>Simulated Communications</span>
              <span class={featurePill}>Real-Time SSE Sync</span>
            </div>
          </div>
        </div>

        {/* Deterministic Safety Guarantees */}
        <div class={safetyGrid}>
          <div class={safetyCard}>
            <div class={safetyIconBox}>
              <Lock size={18} />
            </div>
            <h3 class={safetyTitle}>Zero Unchecked AI Authority</h3>
            <p class={safetyDesc}>
              AI agents are strictly read-only analytical intelligence. Structured action proposals
              must pass multi-judge policy evaluation before the Go application layer executes
              transactional state mutations.
            </p>
          </div>

          <div class={safetyCard}>
            <div class={safetyIconBox}>
              <Database size={18} />
            </div>
            <h3 class={safetyTitle}>ClickHouse Half-Life Decay</h3>
            <p class={safetyDesc}>
              Facility accuracy, turnaround speed, and defect rates are calculated over millions of
              historical QC inspections using an exponential 120-day recency decay model via
              ClickHouse MCP.
            </p>
          </div>

          <div class={safetyCard}>
            <div class={safetyIconBox}>
              <RefreshCw size={18} />
            </div>
            <h3 class={safetyTitle}>Bounded Self-Correction</h3>
            <p class={safetyDesc}>
              Policy verification retries are strictly capped at 3 attempts. When safety margins or
              premiere countdown buffers cannot be satisfied, Fincher safely halts the affected
              delivery and alerts operations.
            </p>
          </div>
        </div>
      </section>

      {/* 100% Event-Driven Closed-Loop Feedback Section */}
      <section id="feedback-loop" class={sectionContainer}>
        <div class={loopSection}>
          <div class={loopHeader}>
            <div class={feedbackBadge}>
              <RefreshCw size={11} />
              <span>100% Event-Driven Closed Loop</span>
            </div>
            <h2 class={sectionHeading}>
              Every executed action feeds ClickHouse &amp; trains the system.
            </h2>
            <p class={sectionSubheading}>
              Fincher never operates in an open loop. When a remediation plan executes, the
              resulting delivery states, vendor re-conforms, and downstream QC results emit new
              CloudEvents directly into ClickHouse, continuously enriching historical intelligence
              in real-time.
            </p>
          </div>

          <div class={loopGrid}>
            <div class={loopCard}>
              <span class={loopStepNum}>STAGE 01</span>
              <h3 class={loopStepTitle}>CloudEvent Ingest</h3>
              <p class={loopStepDesc}>
                Real-time audio drift or corrupted asset signal ingested via HTTP webhook.
              </p>
            </div>

            <div class={loopCard}>
              <span class={loopStepNum}>STAGE 02</span>
              <h3 class={loopStepTitle}>ClickHouse Evidence</h3>
              <p class={loopStepDesc}>
                Historian agent queries past vendor reliability with 120-day half-life decay.
              </p>
            </div>

            <div class={loopCard}>
              <span class={loopStepNum}>STAGE 03</span>
              <h3 class={loopStepTitle}>Policy Verification</h3>
              <p class={loopStepDesc}>
                Multi-judge gate verifies blast radius, premiere countdown, and budget limits.
              </p>
            </div>

            <div class={loopCard}>
              <span class={loopStepNum}>STAGE 04</span>
              <h3 class={loopStepTitle}>Software Execution</h3>
              <p class={loopStepDesc}>
                Go transactional engine halts deliveries &amp; re-routes stems to top-ranked vendor.
              </p>
            </div>

            <div class={`${loopCard} ${loopCardActive}`}>
              <span class={loopStepNum}>STAGE 05 · FEEDBACK</span>
              <h3 class={loopStepTitle}>ClickHouse Continuous Learning</h3>
              <p class={loopStepDesc}>
                Executed actions emit downstream events back into ClickHouse, sharpening future AI
                judgments.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Conversational AI Intelligence Showcase */}
      <section id="intelligence" class={sectionContainer}>
        <div class={sectionHeader}>
          <div class={sectionBadge}>
            <Bot size={13} />
            <span>Fincher Intelligence Engine</span>
          </div>
          <h2 class={sectionHeading}>Fincher AI answers every supply chain query.</h2>
          <p class={sectionSubheading}>
            Ask anything about live release blockers, historical facility error rates, blast radius
            calculations, or localized audio dub statuses. Fincher synthesizes real-time SQLite
            records with ClickHouse analytical queries to deliver verifiable, cited answers.
          </p>
        </div>

        <div class={chatShowcaseContainer}>
          <div class={chatPromptsCol}>
            <div class={chatPromptsHeader}>
              <h3 class={chatPromptsTitle}>Sample Operator Inquiries</h3>
              <p class={chatPromptsSubtitle}>
                Click any prompt to inspect Fincher AI's multi-source synthesis.
              </p>
            </div>

            {Object.values(CHAT_SAMPLES).map((sample) => (
              <button
                key={sample.key}
                type="button"
                class={
                  activeChatKey === sample.key
                    ? `${chatPromptButton} ${chatPromptButtonActive}`
                    : chatPromptButton
                }
                onClick={() => setActiveChatKey(sample.key)}
              >
                <div class={chatPromptLeft}>
                  <HelpCircle size={15} class={chatPromptIcon} />
                  <span class={chatPromptText}>{sample.prompt}</span>
                </div>
                <span class={chatPromptTag}>{sample.tag}</span>
              </button>
            ))}
          </div>

          <div class={chatWindowCard}>
            <div class={chatWindowHeader}>
              <div class={chatWindowHeaderLeft}>
                <div class={chatWindowAvatar}>
                  <Bot size={14} />
                </div>
                <span class={chatWindowTitle}>Fincher Supply Chain Copilot</span>
              </div>
              <span class={chatWindowStatus}>
                <Sparkles size={11} />
                <span>Active Intelligence</span>
              </span>
            </div>

            <div class={chatWindowBody}>
              <div class={chatUserBubble}>{currentChat.prompt}</div>

              <div class={chatAiBubble}>
                <p class={chatAiAnswerText}>{currentChat.answer}</p>

                <div class={chatCitationsWrapper}>
                  <span class={chatCitationsLabel}>Live Audit &amp; Data Citations</span>
                  <div class={chatCitationList}>
                    {currentChat.citations.map((cit) => (
                      <span key={cit.label} class={chatCitationPill}>
                        {cit.type === 'clickhouse' ? (
                          <Database size={10} />
                        ) : (
                          <Terminal size={10} />
                        )}
                        <span>{cit.label}</span>
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            <div class={chatFooterActionRow}>
              <span class={chatFooterNote}>
                Synthesizing live SQLite state + ClickHouse MCP analytical store
              </span>
              <Link to="/chat" class={chatGoToChatLink}>
                <span>Open Full Operator Copilot</span>
                <ChevronRight size={13} />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Interactive Workflow Playground */}
      <section id="playground" class={sectionContainer}>
        <div class={sectionHeader}>
          <div class={sectionBadge}>
            <Zap size={13} />
            <span>Interactive Operational Workflows</span>
          </div>
          <h2 class={sectionHeading}>Inspect autonomous agent execution live.</h2>
          <p class={sectionSubheading}>
            Select an operational scenario below to explore how Fincher triages anomalies, queries
            ClickHouse, evaluates policy verifiers, and executes state transitions.
          </p>
        </div>

        <div class={demoContainer}>
          <div class={demoNav}>
            <button
              type="button"
              class={activeScenario === 'drift' ? `${demoTabBtn} ${demoTabBtnActive}` : demoTabBtn}
              onClick={() => setActiveScenario('drift')}
            >
              <ShieldAlert size={14} />
              <span>Audio Sync Drift</span>
            </button>
            <button
              type="button"
              class={
                activeScenario === 'allocation' ? `${demoTabBtn} ${demoTabBtnActive}` : demoTabBtn
              }
              onClick={() => setActiveScenario('allocation')}
            >
              <Users size={14} />
              <span>Multi-Market Allocation</span>
            </button>
            <button
              type="button"
              class={activeScenario === 'master' ? `${demoTabBtn} ${demoTabBtnActive}` : demoTabBtn}
              onClick={() => setActiveScenario('master')}
            >
              <Film size={14} />
              <span>Master Revision Handoff</span>
            </button>
          </div>

          <div class={demoBody}>
            <div class={demoWaterfallCol}>
              <div class={demoWaterfallHeader}>
                Execution Latency Waterfall · {currentDemo.trigger}
              </div>
              {currentDemo.steps.map((s, idx) => (
                <div key={s.name} class={demoStepRow}>
                  <div class={demoStepLeft}>
                    <span class={stepIndexText}>{idx + 1}.</span>
                    <span class={demoStepName}>{s.name}</span>
                    <span class={demoStepBadge}>{s.label}</span>
                  </div>
                  <div class={stepRightBox}>
                    <span class={demoStepLatency}>{s.latency}</span>
                    <Check size={13} class={stepCheckIcon} />
                  </div>
                </div>
              ))}
            </div>

            <div class={demoScorecardCol}>
              <div class={demoScorecardTitle}>
                <ShieldCheck size={14} />
                <span>{currentDemo.judge}</span>
              </div>

              <div class={demoScorecardVerdict}>
                <Badge variant="success">{currentDemo.verdict}</Badge>
              </div>

              <p class={demoScorecardText}>{currentDemo.rationale}</p>

              <div class={executedActionsWrapper}>
                <span class={executedActionsTitle}>Executed Software Actions</span>
                {currentDemo.actions.map((act) => (
                  <div key={act} class={executedActionPill}>
                    <ArrowRight size={11} class={executedActionIcon} />
                    <span>{act}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Metrics Strip */}
      <section id="metrics" class={sectionContainer}>
        <div class={metricsGrid}>
          <div class={metricCard}>
            <span class={metricValue}>&lt; 12s</span>
            <span class={metricLabel}>Remediation Velocity</span>
            <span class={metricSubtext}>Full triage to dispatch execution</span>
          </div>

          <div class={metricCard}>
            <span class={metricValue}>100%</span>
            <span class={metricLabel}>Policy Provenance</span>
            <span class={metricSubtext}>Every action audited in SQLite &amp; ClickHouse</span>
          </div>

          <div class={metricCard}>
            <span class={metricValue}>5+</span>
            <span class={metricLabel}>Simultaneous Markets</span>
            <span class={metricSubtext}>Synchronized multi-territory delivery DAGs</span>
          </div>

          <div class={metricCard}>
            <span class={metricValue}>0</span>
            <span class={metricLabel}>Manual Spreadsheets</span>
            <span class={metricSubtext}>Continuous self-healing release pipelines</span>
          </div>
        </div>
      </section>

      {/* Final Call to Action */}
      <section class={sectionContainer}>
        <div class={finalCtaCard}>
          <h2 class={finalCtaTitle}>Experience autonomous media operations.</h2>
          <p class={finalCtaSubtitle}>
            Launch Fincher to inspect live title countdowns, multi-market package lineage, and
            real-time agent workflow traces.
          </p>
          <div class={finalCtaButtons}>
            <Link to="/titles" class={primaryCtaBtn}>
              <span>Open Fincher Console</span>
              <ArrowRight size={15} />
            </Link>
            <Link to="/chat" class={secondaryCtaBtn}>
              <Bot size={15} />
              <span>Talk to Operator Copilot</span>
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}
