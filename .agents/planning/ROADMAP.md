# Fincher — Feature-Wise Project Roadmap

Build one complete feature slice end-to-end at a time:
**Domain Entity Models + Store / DB Ops + REST API + Preact UI + Verification Tests**

---

## Milestone Progression (Feature-by-Feature)

```text
[ Feature 01: Titles & Launch Calendar ]          (Completed - Backend)
                 │
                 ▼
[ Feature 02: Masters, Packages & Vendors ]       (Completed - Backend)
                 │
                 ▼
[ Feature 03: Deliveries & Dependencies ]         (Completed - Backend)
                 │
                 ▼
[ Feature 04: UI Scaffolding & Operations Console ] (Completed - Preact UI + Modals + Uploads)
                 │
                 ▼
[ Feature 05: ClickHouse Ingestion Pipeline ]     (Completed - Schema, MCP client, taxonomy, batch insert endpoint. Zero LLM cost.)
                 │
                 ▼
[ Feature 06: Dedicated Agents ]                  (Completed - ADK Go v2 graph + tool-driven two-phase agents, strict MCP reads, live SSE + xyflow viz. Hardened via 5-scenario live E2E.)
                 │
                 ▼
[ Feature 07: Chat Assistant ]                    (Active - Gemini Chat + Live Database & MCP Tools)
                 │
                 ▼
[ Feature 08: Smart Simulator & Comms Hub ]       (Dynamic Event Generator + Protected Seed Data + Dispatches)
                 │
                 ▼
[ Feature 09: Console UI Refinement ]             (Operator control buttons, runs/deliveries polish, landing/visual passes. Deferred UI work.)
```

> Feature 05 and 06 replace the earlier single "ClickHouse MCP & Multi-Agent Engine" milestone —
> split so the zero-LLM ingestion substrate (05) is a clean, independently testable phase before
> any agent/graph code (06) is written on top of it.

---

## Milestone Deliverables & Acceptance

### Feature 01: Titles & Launch Calendar (Completed)
* **Scope**: LUME title catalog management, premiere date countdowns, status tracking.
* **Deliverables**:
  - `pkg/domain/models/title.go`: Title model, status enums (`ON_TRACK`, `AT_RISK`, `HOLD`, `PROCESSING`, `SHIPPED`), validation, `HoursUntilPremiere`.
  - `internal/turso/titles/`: Create, Get, List (with status filter and server-side pagination), Update, Delete.
  - `internal/api/titles/`: `GET /api/titles`, `GET /api/titles/{id}`, `POST /api/titles`, `PATCH /api/titles/{id}`, `DELETE /api/titles/{id}`.
  - `internal/api/server.go`: Base Echo HTTP server setup with JSON middleware.
* **Verification**: Unit tests and HTTP tests verifying title CRUD and pagination.

---

### Feature 02: Masters, Packages & Vendors (Completed)
* **Scope**: Asset version tracking, audio/subtitle/video packages, vendor degradation registry.
* **Deliverables**:
  - `pkg/domain/models/master.go`, `package.go`, `vendor.go`.
  - `internal/turso/masters/`, `packages/`, `vendors/`.
  - `internal/api/masters/`, `packages/`, `vendors/`.
* **Verification**: Unit tests and HTTP tests for package staleness against masters and vendor lookups.

---

### Feature 03: Deliveries & Dependencies (Completed)
* **Scope**: Territory deliveries, blocking run IDs, and asset dependency DAG with cycle detection.
* **Deliverables**:
  - `pkg/domain/models/delivery.go`, `dependency.go`.
  - `internal/turso/deliveries/`, `dependencies/` (with cycle prevention and recursive lineage graph builder).
  - `internal/api/deliveries/`, `dependencies/` (`GET /api/dependencies/graph/{title_id}`).
* **Verification**: Territory delivery hold/release status mutations, dependency graph queries, cycle rejection tests.

---

### Feature 04: UI Scaffolding & Initial Operations Console (Completed)
* **Scope**: Setup `web/` workspace with Bun, Vite, Preact, TanStack Router, TanStack Query, Vanilla Extract design system, Hey API codegen, dedicated feature modals, SQLite BLOB image uploads (1MB limit), and centered layouts.
* **Deliverables**:
  - `web/`: Scaffolding with Preact, Bun, Vite, `@vanilla-extract/css`, `@tanstack/react-router`, `@tanstack/react-query`.
  - `web/src/styles/theme.css.ts`: Dark operations theme tokens.
  - `web/src/lib/api/generated/`: Auto-generated SDK from `openapi/swagger.json`.
  - `web/src/components/ui/`: `button/`, `badge/`, `modal/` (`Modal`, `DeleteModal`), `input/` (`FormField`, `TextInput`, `SelectInput`, `ImageUpload`).
  - `web/src/features/*/components/modals/`: Dedicated creation modals for Titles, Vendors, Deliveries, and Packages.
  - `internal/turso/uploads/` & `internal/api/uploads/`: SQLite BLOB storage with server-side MIME sniffing, strict 1MB cap, `nosniff` headers, and `DELETE`.
  - `web/src/routes/`: Centered empty states on X/Y axes across Titles, Vendors, Deliveries, and Runs.
* **Verification**: `bun run typecheck`, `biome check src`, `bun test`, and full browser navigation.

---

### Feature 05: ClickHouse Ingestion Pipeline (Completed — zero LLM cost)
* **Scope**: The substrate every agent depends on but which itself never calls a model — CNCF CloudEvents v1.0 ClickHouse schema, both ClickHouse access paths (agent-facing MCP + direct database/sql client), CloudEvents taxonomy, and a direct batch ingestion endpoint. Debounce/coalesce and the budget/rate-limit gate were scoped out of this phase entirely (2026-08-27) as unwarranted complexity at this scale — deferred to Feature 06, if/when needed, backed by Turso persistence (not in-memory) to survive Cloud Run scale-to-zero.
* **Deliverables**:
  - `migrations/clickhouse/001_events.sql` .. `003_vendor_metrics.sql`: CloudEvents v1.0 root event stream, QC inspection projection, vendor historical rollups.
  - `internal/clickhouse/`: direct deterministic `database/sql` client (`client.go`, self-bootstrapping `CREATE DATABASE IF NOT EXISTS`) + domain packages (`events/`, `vendors/`) — recency-weighted accuracy (120-day decay), vendor rollups. No LLM involved.
  - `pkg/mcp/`: agent-facing MCP client built on the official `modelcontextprotocol/go-sdk` + `google.golang.org/adk/v2/tool/mcptoolset`, connecting Gemini agents in Phase 06 to `mcp-clickhouse:8000` with ClickHouse native `readonly=1` safety.
  - `pkg/domain/models/event.go` + `event_payloads.go`: CloudEvents v1.0 struct, static category classifier (`TELEMETRY`, `ROUTINE_OUTCOME`, `ANOMALY_SIGNAL`, `ALLOCATION_REQUEST`, `OPERATOR_FORCED`), typed per-event payload structs.
  - `internal/api/events/`: `POST /api/events` — accepts a JSON array of CloudEvents, validates, inserts each directly into ClickHouse, returns `{status, count}`. No classification-based branching, no debounce, no dispatch — Feature 06 owns everything downstream of the ClickHouse write.
* **Verification**: ClickHouse migrations apply cleanly against the `docker-compose` container (including from a genuinely fresh/no-database instance); CloudEvent classifier + typed-payload table-driven tests; batch endpoint HTTP lifecycle tests (success, empty array, non-array, malformed JSON); MCP round-trip test against `mcp-clickhouse`.

---

### Feature 06: Dedicated Agents (ADK Go v2 Graph) (Completed — hardened via live E2E)
> Shipped as staged Go workflow graphs (`internal/agent/graph/`) rather than the literal ADK `workflow`
> node primitives originally sketched below. Agent ClickHouse reads route strictly through the official
> `mcp-clickhouse` `run_query` tool; vendor selection and incident remediation are tool-calling agents
> using a two-phase tool→schema pattern (`runToolThenSchema`) plus a read-only `query_turso` tool.
> Hardened and verified with a 5-scenario judge-grade live E2E (all PASS). See STATE.md 2026-09-08 log.

* **Scope**: Everything that actually calls Gemini — the incident-investigation graph and the vendor-allocation graph, both built as ADK Go v2 `workflow` graphs firing per individual `ANOMALY_SIGNAL`/`ALLOCATION_REQUEST` event from Feature 05 (no batching/coalescing upstream), plus the budget/concurrency gate (Turso-persisted, deferred from Feature 05 — see `REQ-AGENT-10`), the live SSE stream, and `@xyflow/react` visualization of the graph executing.
* **Deliverables**:
  - `internal/agent/graph.go`: ADK Go v2 wiring (`workflow.NewFunctionNode`, `workflow.NewAgentNode`, `workflow.Chain`/`Concat`, `workflow.NewJoinNode`, `workflowagent.New`).
  - `internal/agent/triage_judge.go`, `historian.go` (hybrid), `lineage.go` (Go-only), `optimizer.go`, `policy_judge.go` (bounded reject→revise, capped retries), `executor.go` (transactional SQLite + SSE + downstream event emission).
  - `internal/agent/vendor_scoring.go` (Go-only evidence assembly) + `vendor_judge.go` (always-fires allocation judge).
  - `internal/turso/ent/schema/`: extend `vendor.go` with `hourly_rate_usd` and `turnaround_hours`; new `run.go` (`title_slug`), `step.go`, and `wf_result.go` schemas.
  - `internal/turso/runs/`: CRUD actions for runs, steps, and results.
  - `internal/api/runs/`: `GET /api/runs`, `GET /api/runs/{id}`, `GET /api/runs/{id}/stream` (SSE step transition + result feed); `internal/api/investigations/`: operator-forced trigger endpoints.
  - `web/src/features/runs/`: live investigation graph reusing the `@xyflow/react` canvas, `useSSEStream` hook, judge verdict + rationale inline display.
* **Verification**: per-node unit tests with a stub Gemini client; full-graph integration test against a fixture `ANOMALY_SIGNAL`/`ALLOCATION_REQUEST` event; live demo run against real Gemini; frontend `bun run typecheck`, `biome check src`, and browser walkthrough showing nodes lighting up in real time.

---

### Feature 07: Chat Assistant (Gemini Chat) (Active)
* **Scope**: A read-only natural-language operator chat that answers supply-chain questions by
  reasoning over live SQLite operational state and ClickHouse analytical history, using the SAME tool
  layer the Feature 06 agents already use (`query_turso` + MCP `run_query` analytics). Strictly read-only
  (Invariant 1) — the chat assistant never mutates state or triggers workflows; it explains and cites.
* **Deliverables**:
  - `internal/api/chat/`: `POST /api/chat` (submit a message, returns `{session_id}`), and
    `GET /api/chat/:session/stream` (SSE streaming assistant tokens/answer + structured tool-call
    citations). Session state kept in-memory (ephemeral, scale-to-zero acceptable for a chat session).
  - `internal/agent/chat_agent.go`: a Gemini tool-calling agent (reuses `tools.NewTursoQueryTool`,
    `tools.NewAnalyticsTool`, and title/projection tools) with a read-only system prompt; emits each
    tool call as an auditable citation.
  - `web/src/features/chat/`: feature slice (`queryKeys.ts`, `queryOptions.ts`, message list +
    composer + citation pills components, co-located `*.css.ts`, `index.ts` barrel) and a shared
    `web/src/lib/hooks/useSSEStream.ts` consuming the chat stream.
  - `web/src/routes/chat.tsx`: wire the existing static shell to the live streaming endpoint with
    conversation history, suggested prompts, and inline ClickHouse/SQLite citations.
* **Verification**: live interactive queries against real ClickHouse + Turso asserting accurate answers
  and that surfaced citations correspond to real executed tool calls; frontend `bun run typecheck`,
  `biome check src`, and a browser walkthrough of a multi-turn conversation with streaming + citations.

---

### Feature 09: Console UI Refinement (Deferred UI Work)
* **Scope**: The operator-facing console polish consciously deferred out of Feature 06 — explicit
  control buttons and visual refinement, built on backend endpoints that already work end-to-end.
* **Deliverables**:
  - Operator controls: "Release Hold" (`PATCH /api/deliveries/:id`) and "Re-run Workflow"
    (`POST /api/events` / `POST /api/runs`) button-level triggers in `runs.tsx` / `deliveries.tsx`.
  - Runs and deliveries console visual/interaction polish; additional landing/visual passes beyond the
    Feature 07-era hero clarity work already shipped.
* **Verification**: `bun run typecheck`, `biome check src`, and browser walkthroughs of each operator
  control performing the correct backend mutation and reflecting live via SSE.

---

### Feature 08: Smart Simulator, Protected Seed Data & Communications Hub
* **Scope**: Context-aware dynamic event simulator tab (`/simulator`), protected baseline reference seed data, and automated multi-channel stakeholder dispatches.
* **Deliverables**:
  - `cmd/seed/main.go`: Seed protected system titles (*Eclipse*, *Atlas*) and top vendors into ClickHouse & SQLite.
  - `web/src/routes/simulator.tsx`: Dynamic event generator allowing operators to trigger master revisions, QC failures on audio drift, or vendor re-conform deliveries.
  - `internal/agent/comms.go`: Autonomous generation of Vendor SLA notices, Public Social broadcasts (X/Twitter), and Executive Briefings.
  - `web/src/features/runs/`: Live workflow activity stream and communications view.
* **Verification**: Unattended demo cold start, dynamic event generation, and full closed-loop recovery.
