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
[ Feature 06: Workflow Execution Engine & SSE Runs ]
                 │
                 ▼
[ Feature 07: 17-Node Palette & Multi-Agent Investigation ]
                 │
                 ▼
[ Feature 08: Operations Console UI & Docent Assistant ]
                 │
                 ▼
[ Feature 09: Causal Media Seeder & Final Polish ]
```

---

## Milestone Deliverables & Acceptance

### Feature 01: Titles & Launch Calendar
* **Scope**: LUME title catalog management, premiere date countdowns, status tracking.
* **Deliverables**:
  - `pkg/domain/models/title.go`: Title model, status enums (`ON_TRACK`, `AT_RISK`, `HOLD`, `PROCESSING`, `SHIPPED`), validation, `HoursUntilPremiere`.
  - `internal/turso/titles/`: Create, Get, List (with status filter), Update, Delete.
  - `internal/api/titles/`: `GET /api/titles`, `GET /api/titles/{id}`, `POST /api/titles`, `PATCH /api/titles/{id}`, `DELETE /api/titles/{id}`.
  - `internal/api/server.go`: Base Echo HTTP server setup with JSON middleware.
* **Verification**: `go test -v ./... -race` verifying title CRUD and HTTP endpoints.

---

### Feature 02: Masters, Packages & Vendors
* **Scope**: Asset version tracking, audio/subtitle/video packages, vendor degradation registry.
* **Deliverables**:
  - `pkg/domain/models/master.go`, `package.go`, `vendor.go`.
  - `internal/turso/masters/`, `packages/`, `vendors/`.
  - `internal/api/masters/`, `packages/`, `vendors/`.
* **Verification**: Unit tests and HTTP tests for package staleness against masters and vendor lookups.

---

### Feature 03: Deliveries & Dependencies
* **Scope**: Territory deliveries, blocking run IDs, and asset dependency DAG with cycle detection.
* **Deliverables**:
  - `pkg/domain/models/delivery.go`, `dependency.go`.
  - `internal/turso/deliveries/`, `dependencies/` (with cycle prevention and recursive lineage graph builder).
  - `internal/api/deliveries/`, `dependencies/` (`GET /api/dependencies/graph/{title_id}`).
* **Verification**: Territory delivery hold/release status mutations, dependency graph queries, cycle rejection tests.

---

### Feature 04: ClickHouse History & MCP Client
* **Scope**: ClickHouse schema migrations, 4 Materialized Views, and official MCP HTTP client.
* **Deliverables**:
  - `data/clickhouse/001_qc_events.sql` through `009_mv_recent_master_changes.sql`.
  - `pkg/mcp/client.go`: MCP JSON-RPC 2.0 HTTP transport client (`run_query`, `list_tables`).
  - `internal/turso/query_log/`: Query logging and provenance tracking.
* **Verification**: Live MCP connection test against `mcp-clickhouse` executing analytical queries against views.

---

### Feature 05: Workflow Definitions & DAG Spec
* **Scope**: GraphSpec models, workflow template catalog, and definition authoring.
* **Deliverables**:
  - `pkg/domain/models/workflow_def.go`: DAG node & edge spec validation.
  - `internal/turso/workflow_defs/`.
  - `internal/api/workflows/`: `GET /api/workflows`, `GET /api/workflows/{id}`, `POST /api/workflows`, `PATCH /api/workflows/{id}`, `GET /api/node-palette`.
* **Verification**: Validates DAG topology (acyclic, valid triggers, valid node palette).

---

### Feature 06: Workflow Execution Engine & SSE Runs
* **Scope**: Single-request DAG runner, node dependency resolution, execution state recording, live SSE streaming.
* **Deliverables**:
  - `pkg/domain/models/run.go`: `WorkflowRun`, `NodeExecution`, `NodeInputs/Outputs`.
  - `internal/engine/runner.go`: In-memory topological DAG runner.
  - `internal/api/runs/`: `POST /api/workflows/{id}/run`, `GET /api/runs/{id}`, `GET /api/runs/{id}/stream` (SSE live step execution).
* **Verification**: Runs a minimal DAG (`schedule_trigger` $\rightarrow$ `delta_gate` $\rightarrow$ `event_emitter`) with full execution trail.

---

### Feature 07: 17-Node Palette & Multi-Agent Investigation
* **Scope**: All 17 node types, parallel Gemini Flash agents, Pro Decision Node (disagreement logic), actions, and drafted notices.
* **Deliverables**:
  - `internal/engine/nodes/`: Implement all 17 nodes.
  - `internal/turso/decisions/`, `notifications/`, `executed_actions/`.
* **Verification**: Full hero run (*Eclipse* V13 + Vendor A audio drift $\rightarrow$ HOLD decision + drafted notice).

---

### Feature 08: Operations Console UI & Docent Operator Assistant
* **Scope**: Preact Operations Console UI with zero-runtime Vanilla Extract styling, co-located camelCase component architecture, live SSE DAG visualizer, and Docent natural language assistant.
* **Deliverables**:
  - `web/`: Preact + Vite + TanStack DB + TanStack Router + Vanilla Extract + Hey API.
  - `web/src/styles/theme.css.ts`: Dark operations theme tokens.
  - `web/src/features/calendar/`: Launch Calendar with title countdowns and master cuts.
  - `web/src/features/lineage/`: Interactive XYFlow Lineage DAG.
  - `web/src/features/deliveries/`: Territory Delivery Matrix with hold overrides.
  - `web/src/features/vendors/`: Vendor scorecards and reliability charts.
  - `web/src/features/runs/`: Live SSE Studio Pipeline run inspector.
  - `web/src/features/docent/`: Read-first natural language ClickHouse query assistant.
  - `internal/api/query/`: `POST /api/query`, `GET /api/query/{session}/stream` (SSE).
* **Verification**: `bun run typecheck`, `bun run lint`, and end-to-end verification of embedded UI in Go server binary.

---

### Feature 09: Causal Media Seeder & Hackathon Final Polish
* **Scope**: Budget counter caps, scale-to-zero verification, and LUME demo causal seeder.
* **Deliverables**:
  - `internal/turso/budget/`: Enforces daily model limit & kill switch.
  - `cmd/seed/main.go` & `data/seed/lume.go`: Seeds 250k+ ClickHouse QC rows, 5–7 titles, pre-resolved hero run.
  - `cmd/fincher/main.go`: Production unified server binary serving API + embedded UI.
* **Verification**: Full unattended demo cold start and compliance verification.
