# Fincher — Feature-Wise Project Roadmap

Build one complete feature slice end-to-end at a time:
**Domain Entity Models + SQL Migrations + Store / DB Ops + REST API + Verification Tests**

---

## Milestone Progression (Feature-by-Feature)

```text
[ Feature 01: Titles & Launch Calendar ]
                 │
                 ▼
[ Feature 02: Masters, Packages & Vendors ]
                 │
                 ▼
[ Feature 03: Deliveries & Dependencies ]
                 │
                 ▼
[ Feature 04: ClickHouse History & MCP Client ]
                 │
                 ▼
[ Feature 05: Workflow Definitions & DAG Spec ]
                 │
                 ▼
[ Feature 06: Workflow Execution Engine & Runs ]
                 │
                 ▼
[ Feature 07: Node Palette & Multi-Agent Investigation ]
                 │
                 ▼
[ Feature 08: Docent Operator Assistant API ]
                 │
                 ▼
[ Feature 09: System Hardening & Hackathon Seeder ]
```

---

## Milestone Deliverables & Acceptance

### Feature 01: Titles & Launch Calendar
* **Scope**: LUME title catalog management, premiere date countdowns, status tracking.
* **Deliverables**:
  - `pkg/domain/models/title.go`: Title model, status enums (`ON_TRACK`, `AT_RISK`, `HOLD`, `PROCESSING`, `SHIPPED`), validation, `HoursUntilPremiere`.
  - `migrations/turso/001_titles.sql`: `titles` table.
  - `pkg/turso/titles.go`: Create, Get, List (with status filter), Update, Delete.
  - `internal/api/titles.go`: `GET /titles`, `GET /titles/{id}`, `POST /titles`, `PATCH /titles/{id}`, `DELETE /titles/{id}`.
  - `internal/api/server.go`: Base Echo HTTP server setup with JSON middleware.
* **Verification**: `go test -v ./... -race` verifying title CRUD and HTTP endpoints.

---

### Feature 02: Masters, Packages & Vendors
* **Scope**: Asset version tracking, audio/subtitle/video packages, vendor degradation registry.
* **Deliverables**:
  - `pkg/domain/models/master.go`, `package.go`, `vendor.go`.
  - `migrations/turso/002_masters.sql`, `003_packages.sql`, `004_vendors.sql`.
  - `pkg/turso/masters.go`, `packages.go`, `vendors.go`.
  - `internal/api/packages.go`, `internal/api/vendors.go`.
* **Verification**: Unit tests and HTTP tests for package staleness against masters and vendor lookups.

---

### Feature 03: Deliveries & Dependencies
* **Scope**: Territory deliveries, blocking run IDs, and asset dependency graphs.
* **Deliverables**:
  - `pkg/domain/models/delivery.go`, `dependency.go`.
  - `migrations/turso/005_deliveries.sql`, `006_dependencies.sql`.
  - `pkg/turso/deliveries.go`, `dependencies.go`.
  - `internal/api/deliveries.go`.
* **Verification**: Territory delivery hold/release status mutations and dependency queries.

---

### Feature 04: ClickHouse History & MCP Client
* **Scope**: ClickHouse schema migrations, 4 Materialized Views, and official MCP HTTP client.
* **Deliverables**:
  - `migrations/clickhouse/001_qc_events.sql` through `009_mv_recent_master_changes.sql`.
  - `pkg/mcp/client.go`: MCP JSON-RPC 2.0 HTTP transport client (`run_query`, `list_tables`).
  - `pkg/turso/query_log.go`: Query logging and provenance tracking.
* **Verification**: Live MCP connection test against `mcp-clickhouse` executing analytical queries against views.

---

### Feature 05: Workflow Definitions & DAG Spec
* **Scope**: GraphSpec models, workflow template catalog, and definition authoring.
* **Deliverables**:
  - `pkg/domain/models/workflow_def.go`: DAG node & edge spec validation.
  - `migrations/turso/007_workflow_definitions.sql`.
  - `pkg/turso/workflow_defs.go`.
  - `internal/api/workflows.go`: `GET /workflows`, `GET /workflows/{id}`, `POST /workflows`, `PATCH /workflows/{id}`, `GET /node-palette`.
* **Verification**: Validates DAG topology (acyclic, valid triggers, valid node palette).

---

### Feature 06: Workflow Execution Engine & Runs
* **Scope**: Single-request DAG runner, node dependency resolution, execution state recording.
* **Deliverables**:
  - `pkg/domain/models/run.go`: `WorkflowRun`, `NodeExecution`, `NodeInputs/Outputs`.
  - `migrations/turso/008_workflow_runs.sql`, `009_node_executions.sql`, `010_node_inputs_outputs.sql`.
  - `internal/engine/runner.go`: In-memory topological DAG runner.
  - `internal/api/runs.go`: `POST /workflows/{id}/run`, `GET /runs/{id}`, `GET /runs/{id}/stream` (SSE).
* **Verification**: Runs a minimal DAG (`schedule_trigger` $\rightarrow$ `delta_gate` $\rightarrow$ `event_emitter`) with full execution trail.

---

### Feature 07: Node Palette & Multi-Agent Investigation
* **Scope**: All 17 node types, parallel Gemini Flash agents, Pro Decision Node (disagreement logic), actions, and drafted notices.
* **Deliverables**:
  - `internal/engine/nodes/`: Implement all 17 nodes.
  - `migrations/turso/012_decisions.sql`, `013_executed_actions.sql`, `014_notifications.sql`.
  - `pkg/turso/decisions.go`, `notifications.go`.
* **Verification**: Full hero run (*Eclipse* V13 + Vendor A audio drift $\rightarrow$ HOLD decision + drafted notice).

---

### Feature 08: Docent Operator Assistant API
* **Scope**: Read-first conversational assistant answering launch status and explaining run decisions with SQL query citations.
* **Deliverables**:
  - `pkg/domain/models/assistant.go`.
  - `migrations/turso/016_query_sessions.sql`.
  - `internal/assistant/docent.go`: Context assembler and Gemini Pro narrator.
  - `internal/api/query.go`: `POST /query`, `GET /query/{session}/stream` (SSE).
* **Verification**: Accurate answers to "What's releasing this weekend?" and "Why was Eclipse held?".

---

### Feature 09: System Hardening & Hackathon Seeder
* **Scope**: Budget counter caps, scale-to-zero verification, and LUME demo causal seeder.
* **Deliverables**:
  - `migrations/turso/015_budget_counters.sql`.
  - `pkg/turso/budget.go`: Enforces daily model limit & kill switch.
  - `cmd/seed/main.go` & `data/seed/lume.go`: Seeds 250k+ ClickHouse QC rows, 5–7 titles, pre-resolved hero run.
  - `cmd/fincher/main.go`: Production unified server binary.
* **Verification**: Full unattended demo cold start and compliance verification.
