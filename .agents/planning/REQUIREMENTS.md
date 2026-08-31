# Fincher — System Requirements Specification (Workflow DAG Edition)

## 1. Functional Requirements

### ClickHouse Analytical History & MCP (`REQ-CH`)
* **REQ-CH-01**: `migrations/clickhouse/001_events.sql` — root immutable `fincher.events` table (`MergeTree`, `partition by toYYYYMM(time)`, `order by (type, subject, time, id)`). Every ingested event is an append-only CNCF CloudEvents v1.0 record (`id`, `type`, `source`, `subject`, `time`, `data`, `severity`, `datacontenttype`) before any downstream processing.
  - *Subject as Entity Slug*: The `subject` attribute stores the unique title slug (e.g. `eclipse`) for title-scoped events, or entity slug for non-title events.
  - *Title-Agnostic Sentinel*: Events that are not tied to a specific media title (e.g. `fincher.vendor.heartbeat`, facility-wide capacity alerts) set `subject = 'GLOBAL'`. This avoids the storage and sorting penalties of `Nullable(String)` in `order by` while preventing ambiguity.
* **REQ-CH-02**: `migrations/clickhouse/002_qc.sql` — real-time Materialized View `fincher.qc` projecting from `fincher.events` where `type = 'fincher.qc.completed'` (`event_id`, `title_slug`, `package_id`, `vendor_id`, `component`, `language`, `status`, `sync_drift_ms`, `video_corruption_score`, `defect_category`, `inspector_agent`, `inspected_at`).
* **REQ-CH-03**: `migrations/clickhouse/003_vendor_metrics.sql` — SummingMergeTree Materialized View `fincher.vendor_metrics` projecting daily aggregates (`recorded_date`, `vendor_id`, `component`, `total_inspections`, `failed_inspections`, `warning_inspections`, `measured_status_count`, `total_sync_drift_ms`, `measured_drift_count`) where `type = 'fincher.qc.completed'`.
* **REQ-CH-04**: Real-time projections from `fincher.events` into `fincher.qc` and `fincher.vendor_metrics` operate automatically inside ClickHouse, avoiding backend-orchestrated ETL.
* **REQ-CH-05**: `005_incident_log.sql` — anomaly detection & resolution audit history, joined against by the Vendor Judge to detect open incidents on allocation candidates.
* **REQ-CH-06**: Agent-issued analytical queries route through the official `mcp-clickhouse` server via Google ADK Go's native `mcptoolset` (`FINCHER_MCP_URL`) with ClickHouse native `readonly = 1` security enforcement. Deterministic Go nodes query ClickHouse via standard `database/sql` without LLM overhead.
* **REQ-CH-07**: Accuracy is computed with recency-weighted exponential decay (half-life ~120 days) over `fincher.vendor_metrics`, not a flat lifetime average.

### Turso / libSQL Application State (`REQ-TURSO`)
* **REQ-TURSO-01**: Turso stores the delivery domain: `titles`, `masters`, `packages`, `vendors`, `deliveries`, `dependencies`.
* **REQ-TURSO-02**: Turso stores the agent workflow execution domain: `runs` (workflow run with `title_slug`), `steps` (per-node transitions for SSE replay/audit), `wf_results` (judge outcome + rationale, immutable).
* **REQ-TURSO-03**: Database schema is managed via type-safe Ent ORM schemas in `internal/turso/ent/schema/` with automated migrations.
* **REQ-TURSO-04**: `Vendor` schema extended with operational fields: `hourly_rate_usd` and `turnaround_hours` — quoted/current operational data, distinct from ClickHouse's historical performance rollups (`REQ-CH-03`).

### Event Taxonomy & Agent Graph (`REQ-EVT`, `REQ-AGENT`)
* **REQ-EVT-01**: `pkg/domain/models/` defines CloudEvents v1.0 and a static taxonomy mapping `type` → category: `TELEMETRY`, `ROUTINE_OUTCOME`, `ANOMALY_SIGNAL`, `ALLOCATION_REQUEST`, `OPERATOR_FORCED`. This is ontology (what an event structurally is), never a numeric business threshold.
* **REQ-EVT-02**: `TELEMETRY` and `ROUTINE_OUTCOME` events are written to ClickHouse and never invoke a model.
* **REQ-AGENT-01**: **Dropped (2026-08-27).** Debounce/coalesce/dedup was planned but deliberately cut as unwarranted complexity for this scale — see `STATE.md` decision log. Each `ANOMALY_SIGNAL` event is handed to `triage_judge.go` individually, one call per event, no batching/`InvestigationRequest` aggregation step.
* **REQ-AGENT-02**: `internal/agent/triage_judge.go` (⚡ flash, single call per `ANOMALY_SIGNAL` event) — routes `YES`/`NO` on whether the event warrants full investigation. No hardcoded severity threshold.
* **REQ-AGENT-03**: `internal/agent/historian.go` (hybrid) — pre-baked deterministic ClickHouse queries for known defect categories; falls back to MCP tool-calling for novel/unclassified incidents.
* **REQ-AGENT-04**: `internal/agent/lineage.go` (Go-only, no LLM) — deterministic SQLite dependency DAG walk for affected downstream deliveries.
* **REQ-AGENT-05**: `internal/agent/optimizer.go` (⚡) — synthesizes a typed `ActionPlan` (`HOLD_DELIVERY`, `INVALIDATE_PACKAGE`, `REASSIGN_VENDOR`, ...) from joined evidence.
* **REQ-AGENT-06**: `internal/agent/policy_judge.go` (⚡) — verdict `APPROVED | REJECTED | ESCALATE` + rationale over the proposed `ActionPlan`. `REJECTED` loops back to the Optimizer with the rejection reason, capped at a configurable retry limit (default 3); cap exceeded → `ESCALATE`.
* **REQ-AGENT-07**: `internal/agent/vendor_scoring.go` (Go-only) — assembles per-candidate `VendorEvidence` (recency-decayed accuracy from ClickHouse + standard/rush rate & turnaround from SQLite + open incidents from `incident_log`).
* **REQ-AGENT-08**: `internal/agent/vendor_judge.go` (⚡) — fires on every `ALLOCATION_REQUEST` event (new title/package needing QC), reasons freshly about cost/speed/quality tradeoffs given premiere urgency; no fixed scoring-weight constants, no margin-threshold gate.
* **REQ-AGENT-09**: `internal/agent/executor.go` (Go-only) — applies any approved verdict (`ActionPlan` or vendor assignment) inside a single SQLite transaction, broadcasts each step over SSE, emits the resulting event back into ClickHouse.
* **REQ-AGENT-10**: **Deferred to Feature 06 (parked 2026-08-27).** Daily invocation cap (`FINCHER_DAILY_MODEL_CAP`) and concurrency semaphore (`FINCHER_MAX_CONCURRENT_AI_WORKFLOW_RUNS`) still need to gate Stage C triage before real Gemini calls start, but the exact mechanism (queue shape, eviction policy) is undecided — must be Turso-persisted (not in-memory) to survive Cloud Run scale-to-zero between requests. Not built in Feature 05.
* **REQ-AGENT-11**: All graphs are assembled with ADK Go v2 (`google.golang.org/adk/v2`, `workflow` package): `workflow.NewFunctionNode` for Go-only nodes, `workflow.NewAgentNode` for LLM judges, `workflow.Edge` + `workflow.StringRoute` for conditional routing, `workflow.NewJoinNode` for the historian/lineage fan-in.

### Backend API & Real-Time Stream (`REQ-API`)
* **REQ-API-01**: Operator-forced trigger (`OPERATOR_FORCED` category): `POST /api/investigations` (title-scoped) and `POST /api/vendor-assignments` (package-scoped). There is no user-editable workflow/DAG builder — the agent graph is fixed in Go code (`internal/agent/`), not stored or edited via API.
* **REQ-API-02**: Run inspection: `GET /api/runs/{id}`, `GET /api/runs/{id}/stream` (SSE live node-transition + judge-verdict feed for the xyflow graph visualization).
* **REQ-API-03**: Title catalog & CRUD: `GET /api/titles`, `GET /api/titles/{id}`, `POST /api/titles`, `PATCH /api/titles/{id}`, `DELETE /api/titles/{id}`.
* **REQ-API-04**: Master & Package & Vendor CRUD: `POST /api/masters`, `GET /api/masters`, `POST /api/packages`, `GET /api/packages`, `PATCH /api/packages/{id}`, `POST /api/vendors`, `GET /api/vendors`.
* **REQ-API-05**: Read-first Docent query: `POST /api/query`, `GET /api/query/{session}/stream`.
* **REQ-API-06**: Budget status: `GET /api/budget` (daily invocation cap usage, concurrency semaphore state, for the operator console budget widget).
* **REQ-API-07**: Deliveries CRUD: `GET /api/deliveries`, `GET /api/deliveries/{id}`, `POST /api/deliveries`, `PATCH /api/deliveries/{id}`, `DELETE /api/deliveries/{id}`.
* **REQ-API-08**: Dependency graph & lineage: `POST /api/dependencies`, `GET /api/dependencies`, `GET /api/dependencies/graph/{title_id}`, `DELETE /api/dependencies/{id}`.
* **REQ-API-09**: Code-First OpenAPI Specification: `GET /openapi.json` served live and embedded in Go binary.

### Frontend Operations Console (`REQ-UI`)
* **REQ-UI-01**: Single-Page App built with **Preact + Vite + TypeScript**, served directly from Go binary via `embed.FS` (with dev API proxy).
* **REQ-UI-02**: **Zero-Runtime Styling**: All styles authored in `*.css.ts` using `@vanilla-extract/css` and `@vanilla-extract/recipes` bound to `src/styles/theme.css.ts` dark operations tokens.
* **REQ-UI-03**: **Strict camelCase & Co-location**: All files and directories are `camelCase`. Every component and sub-component has its own directory with `*.tsx`, `*.css.ts`, and `index.ts`.
* **REQ-UI-04**: **Contract-First Codegen**: Type-safe SDK and Valibot schemas generated directly from backend `openapi/swagger.json` via `@hey-api/openapi-ts`.
* **REQ-UI-05**: **Reactive Client Database**: `@tanstack/react-db` collections sync incoming SSE events (`PACKAGE_INVALIDATED`, `DELIVERY_HELD`) directly to local state for 60fps UI reactivity.
* **REQ-UI-06**: **Launch Calendar**: Dynamic visual timeline of LUME titles with days-to-premiere countdowns and overall status pills.
* **REQ-UI-07**: **Interactive Lineage DAG**: `@xyflow/react` canvas showing parent-child asset derivations, staleness indicators, and component nodes.
* **REQ-UI-08**: **Territory Delivery Matrix**: High-density virtualized table (`@tanstack/react-table` + `@tanstack/react-virtual`) displaying 40+ country delivery statuses with hold override actions.
* **REQ-UI-09**: **Live Investigation Graph**: Reuses the `@xyflow/react` DAG canvas (`REQ-UI-07`) to render `internal/agent/` graph nodes lighting up in real time as the SSE stream (`REQ-API-02`) reports transitions; each judge node surfaces its verdict + rationale inline.
* **REQ-UI-10**: **Docent Assistant Interface**: Natural language operator drawer querying ClickHouse with SQL query citations.

### Causal Seeder & LUME World (`REQ-SEED`)
* **REQ-SEED-01**: Idempotent seeder creates 5–7 LUME launch titles (*Eclipse*, *Atlas*, *Orbit*, *Meridian*, *Vantage*) with dynamic premiere dates relative to `now()`.
* **REQ-SEED-02**: Causal vendor datasets (`vendor_a` degrading audio trend, `vendor_b` reliable subtitles, `vendor_c` mixed).
* **REQ-SEED-03**: Pre-populates 250k+ historical QC records and a pre-resolved hero demonstration run.

---

## 2. Non-Functional & Safety Requirements
* **REQ-GO-01**: Idiomatic Go 1.24+ passing `go test ./... -race` with 0 warnings.
* **REQ-COLD-01**: Runs cleanly on cold Cloud Run (`min-instances: 0`) with DAG execution under 5 seconds.
* **REQ-BUDGET-01**: Hard daily model invocation limit and kill switch in `budget_counters` to protect the $100 credit.
* **REQ-MCP-01**: ClickHouse credentials isolated exclusively in MCP service.
* **REQ-LINT-01**: Frontend code must pass `biome check src` and `tsc --noEmit` with 0 errors.
