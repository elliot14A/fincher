# Fincher — Feature-Wise Project Roadmap

Build one complete feature slice end-to-end at a time:
**Domain Entity Models + SQL Migrations + Store / DB Ops + REST API + Verification Tests**

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
[ Feature 04: UI Scaffolding & Initial Console ]  (In-Flight: Setup web/, Calendar, Vendors, Matrix, Lineage)
                 │
                 ▼
[ Feature 05: ClickHouse History & MCP Client ]
                 │
                 ▼
[ Feature 06: Workflow Definitions & DAG Spec ]
                 │
                 ▼
[ Feature 07: Execution Engine & SSE Runs ]       (Backend + Live UI Run Inspector)
                 │
                 ▼
[ Feature 08: 17-Node Palette & Multi-Agent Investigation ] (Backend + Disagreement Panel)
                 │
                 ▼
[ Feature 09: Docent NLQ Assistant & Causal Media Seeder ]   (Backend + Docent Drawer + Polish)
```

---

## Milestone Deliverables & Acceptance

### Feature 01: Titles & Launch Calendar (Completed)
* **Scope**: LUME title catalog management, premiere date countdowns, status tracking.
* **Deliverables**:
  - `pkg/domain/models/title.go`: Title model, status enums (`ON_TRACK`, `AT_RISK`, `HOLD`, `PROCESSING`, `SHIPPED`), validation, `HoursUntilPremiere`.
  - `internal/turso/titles/`: Create, Get, List (with status filter), Update, Delete.
  - `internal/api/titles/`: `GET /api/titles`, `GET /api/titles/{id}`, `POST /api/titles`, `PATCH /api/titles/{id}`, `DELETE /api/titles/{id}`.
  - `internal/api/server.go`: Base Echo HTTP server setup with JSON middleware.
* **Verification**: `go test -v ./... -race` verifying title CRUD and HTTP endpoints.

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

### Feature 04: UI Scaffolding & Initial Operations Console
* **Scope**: Setup `web/` workspace with Bun, Vite, Preact, TanStack DB, TanStack Router, Vanilla Extract design system, Hey API codegen, and build the UI views for the completed backend entities.
* **Deliverables**:
  - `web/`: Scaffolding with Preact, Bun, Vite, `@vanilla-extract/css`, `@vanilla-extract/recipes`, `@tanstack/react-router`, `@tanstack/react-db`.
  - `web/src/styles/theme.css.ts`: Dark operations theme tokens.
  - `web/src/lib/api/generated/`: Auto-generated SDK from `openapi/swagger.json`.
  - `web/src/components/ui/`: Co-located UI primitives (`button/`, `badge/`, `modal/`, `table/`, `input/`, `drawer/`).
  - `web/src/features/layout/`: `appShell/`, `navigationSidebar/`, `operationsTopbar/`.
  - `web/src/features/calendar/` & `web/src/routes/index.tsx`: Launch Calendar with title countdowns.
  - `web/src/features/vendors/` & `web/src/routes/vendors.tsx`: Vendor scorecard registry.
  - `web/src/features/deliveries/` & `web/src/routes/deliveries.tsx`: Territory Delivery Matrix with hold overrides.
  - `web/src/features/lineage/` & `web/src/routes/lineage/$id.tsx`: Interactive XYFlow Lineage DAG.
* **Verification**: `bun run typecheck`, `bun run lint`, and full navigation across Calendar, Vendors, Deliveries, and Lineage.

---

### Feature 05: ClickHouse History & MCP Client
* **Scope**: ClickHouse schema migrations, 4 Materialized Views, and official MCP HTTP client.
* **Deliverables**:
  - `data/clickhouse/001_qc_events.sql` through `009_mv_recent_master_changes.sql`.
  - `pkg/mcp/client.go`: MCP JSON-RPC 2.0 HTTP transport client (`run_query`, `list_tables`).
  - `internal/turso/query_log/`: Query logging and provenance tracking.
* **Verification**: Live MCP connection test against `mcp-clickhouse` executing analytical queries against views.

---

### Feature 06: Workflow Definitions & DAG Spec
* **Scope**: GraphSpec models, workflow template catalog, and definition authoring.
* **Deliverables**:
  - `pkg/domain/models/workflow_def.go`: DAG node & edge spec validation.
  - `internal/turso/workflow_defs/`.
  - `internal/api/workflows/`: `GET /api/workflows`, `GET /api/workflows/{id}`, `POST /api/workflows`, `PATCH /api/workflows/{id}`, `GET /api/node-palette`.
* **Verification**: Validates DAG topology (acyclic, valid triggers, valid node palette).

---

### Feature 07: Execution Engine & SSE Runs (Backend + UI)
* **Scope**: Single-request DAG runner, execution recording, live SSE streaming, and live Run Inspector UI.
* **Deliverables**:
  - **Backend**: `pkg/domain/models/run.go`, `internal/engine/runner.go`, `internal/api/runs/` (`POST /api/workflows/{id}/run`, `GET /api/runs/{id}/stream`).
  - **Frontend**: `web/src/features/runs/`, `web/src/routes/runs/index.tsx`, `web/src/routes/runs/$id.tsx` (Live Studio Pipeline run inspector with real-time SSE stepping).
* **Verification**: Minimal DAG execution trail streamed via SSE and visualized live in browser.

---

### Feature 08: 17-Node Palette & Multi-Agent Investigation (Backend + UI)
* **Scope**: All 17 node types, parallel Gemini Flash agents, Pro Decision Node (disagreement logic), actions, and drafted notices.
* **Deliverables**:
  - **Backend**: `internal/engine/nodes/`, `internal/turso/{decisions,notifications,executed_actions}/`.
  - **Frontend**: Disagreement Panel inside `web/src/features/runs/` + Stakeholder Dispatch drafted notice modal.
* **Verification**: Hero scenario (*Eclipse* V13 + Vendor A audio drift $\rightarrow$ HOLD decision + drafted notice) observed live in UI.

---

### Feature 09: Docent NLQ Assistant & Causal Media Seeder (Backend + UI)
* **Scope**: Read-first natural language operator assistant, budget counters, causal media seeder, and single-binary Go embed.
* **Deliverables**:
  - **Backend**: `internal/assistant/docent.go`, `internal/turso/budget/`, `cmd/seed/main.go`, `cmd/fincher/main.go`.
  - **Frontend**: `web/src/features/docent/` (interactive natural language query drawer with ClickHouse citations).
* **Verification**: Full unattended demo cold start, Docent query verification, and single-binary distribution test.
