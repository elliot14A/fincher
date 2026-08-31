# Phase 06 Context: Dedicated Agents (ADK Go v2 Graph)

## Goal
Build the incident-investigation graph and the vendor-allocation graph as ADK Go v2
`workflow` graphs, each firing per individual `ANOMALY_SIGNAL` / `ALLOCATION_REQUEST` /
`OPERATOR_FORCED` event from Feature 05 (no batching/coalescing upstream — that design was cut,
see `.agents/reviews/14-event-ingestion-pipeline/`), using one reusable node primitive
throughout, then stream every node transition + judge verdict live to the operator console.

## Objectives
* Implement **Evidence → Judgment → Execution** as the one reusable node shape across both
  graphs: Go assembles deterministic evidence, a single-scoped Gemini judge renders a verdict
  + rationale (never a hardcoded threshold), Go executes the verdict transactionally.
* Wire both graphs with ADK Go v2 (`workflow.NewFunctionNode`, `workflow.NewAgentNode`,
  `workflow.Chain`/`Concat`, `workflow.NewJoinNode`, `workflowagent.New`).
* Stream node transitions + verdicts over SSE; render live in a reused `@xyflow/react` canvas.

## Architectural Decisions

### 1. Incident investigation graph
```
START → triage_judge ⚡ (route NO → end / YES → continue)
      → fan-out: historian ⚡hybrid ∥ lineage (Go-only) → join
      → optimizer ⚡ → ActionPlan
      → policy_judge ⚡ → APPROVED → executor
                        → REJECTED → loop to optimizer w/ reason
                          (capped FINCHER_POLICY_JUDGE_MAX_RETRIES, default 3)
                        → ESCALATE → executor writes HOLD + operator alert
```
* `triage_judge.go`: ⚡ flash, single call per `ANOMALY_SIGNAL` event, route `YES`/`NO`. No
  hardcoded severity threshold — reasons over the event + premiere-urgency context.
* `historian.go`: hybrid — `internal/clickhouse` pre-baked queries for known `defect_category`
  values; MCP tool-calling fallback (`pkg/mcp`) for novel/unclassified incidents.
* `lineage.go`: Go-only, reuses `internal/turso/dependencies` graph query. No LLM.
* `optimizer.go`: ⚡, synthesizes a typed `ActionPlan`
  (`pkg/domain/models/actionplan.go`: `HOLD_DELIVERY`, `INVALIDATE_PACKAGE`,
  `REASSIGN_VENDOR`, ...).
* `policy_judge.go`: ⚡, verdict `APPROVED | REJECTED | ESCALATE` + rationale, one
  `wf_results` row per attempt.
* `executor.go`: Go-only, single SQLite tx applying `ActionPlan` actions, SSE broadcast per
  action, emits the resulting event back into ClickHouse (closed loop, AGENTS.md Invariant 4).

### 2. Vendor allocation graph
```
TitleCreated / PackageRequired / VENDOR_RECONFORM_DISPATCHED  (one event at a time)
      → vendor_scoring (Go-only) → vendor_judge ⚡ (always fires) → executor
```
* `vendor_scoring.go`: merges `internal/clickhouse.RecencyWeightedAccuracy` +
  `internal/clickhouse.OpenIncidents` + SQLite `Vendor` rate cards into a `VendorEvidence`
  bundle per candidate.
* `vendor_judge.go`: ⚡, **always fires, no margin-threshold gate** — reasons freshly about
  cost/speed/quality tradeoffs given premiere urgency and any open incidents.
* `executor.go` (shared): applies the vendor assignment, emits `VENDOR_ASSIGNED`.

### 3. Ent schema additions (`internal/turso/ent/schema/`)
* `run.go`: `id`, `title_slug` (default `"GLOBAL"`), `trigger`, `status`, `started_at`, `ended_at`.
* `step.go`: `id`, `run_id` (FK), `name`, `status`, `started_at`, `ended_at` — the SSE replay/audit source.
* `wf_result.go`: `id`, `run_id` (FK), `step_id` (FK), `judge`, `outcome`, `rationale`, `attempt`.
* `vendor.go` extension: `hourly_rate_usd` (`field.Float`, default `0.0`), `turnaround_hours` (`field.Int`, default `24`).

### 4. API & SSE (`internal/api/runs/`, `internal/api/investigations/`)
* `GET /api/runs`, `GET /api/runs/{id}`, `GET /api/runs/{id}/stream` — SSE, one frame per `step` /
  `wf_result` insert.
* `POST /api/investigations`, `POST /api/vendor-assignments` — `OPERATOR_FORCED` manual
  triggers.
* `GET /api/budget` — exposes the budget/concurrency gate's state. **Not yet designed**: Feature
  05 dropped the in-memory `internal/ingest` design entirely (see
  `.agents/reviews/14-event-ingestion-pipeline/`); this phase must build the daily-cap +
  concurrency gate fresh, Turso-persisted (not in-memory) so it survives Cloud Run
  scale-to-zero between requests.

### 5. Frontend (`web/src/features/runs/`)
* Reuse the `@xyflow/react` canvas pattern from the existing Lineage DAG feature.
* New shared hook `src/lib/hooks/useSSEStream.ts` consuming `GET /api/runs/{id}/stream`.
* Nodes light up per transition; judge nodes show verdict + rationale inline.

## Explicit Constraints Carried From the Brainstorm
* No hardcoded severity thresholds, score margins, or fixed scoring weights anywhere in this
  feature — every judgment point is a scoped LLM call over Go-assembled evidence, every time.
* The Policy Judge's reject → revise loop is the **only** intentional cycle in either graph;
  it is capped, never unbounded, and the topology (not a rule table) is what stays
  deterministic.
* MCP is used only by agent tool-calling code (historian's novel-incident fallback). All other
  ClickHouse reads route through `internal/clickhouse` directly.
* Rate cards (SQLite `Vendor`) and historical performance (ClickHouse rollups) are never
  conflated — `vendor_scoring.go` explicitly merges both, from two different stores.

## Reference
* Full brainstorm/rationale captured in this session's chat log.
* `PROJECT.md` §3; `REQUIREMENTS.md` `REQ-AGENT`, `REQ-TURSO-04`, `REQ-API-01/02/06`,
  `REQ-UI-09`.
