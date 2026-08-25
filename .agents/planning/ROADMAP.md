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
[ Feature 05: ClickHouse MCP & Multi-Agent Engine ] (Active - Historian, Lineage, Vendor Optimizer)
                 │
                 ▼
[ Feature 06: Docent Conversational Assistant ]   (Gemini Chat + Live Database & MCP Tools)
                 │
                 ▼
[ Feature 07: Smart Simulator & Comms Hub ]       (Dynamic Event Generator + Protected Seed Data + Dispatches)
```

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

### Feature 05: ClickHouse MCP & Multi-Agent Investigation Engine (Next Up)
* **Scope**: ClickHouse schema migrations, MCP HTTP client integration, Historian Sub-Agent, Lineage Sub-Agent, and Multi-Factor Vendor Optimizer balancing Speed, Quality, and Rates.
* **Deliverables**:
  - `pkg/mcp/client.go`: MCP JSON-RPC 2.0 HTTP transport client (`run_query`, `list_tables`).
  - `internal/agent/historian.go`: Queries ClickHouse for vendor defect rates, audio sync drift telemetry, and past turnarounds.
  - `internal/agent/lineage.go`: Traces SQLite dependency trees for affected market deliveries.
  - `internal/agent/optimizer.go`: Evaluates trade-offs (Speed vs Quality vs Cost) to recommend optimal vendor re-routing.
  - `internal/agent/orchestrator.go`: Coordinates multi-agent parallel execution with Google GenAI / Gemini.
* **Verification**: Live MCP connection test against `mcp-clickhouse` and automated multi-agent run with mock/real incident triggers.

---

### Feature 06: Docent Conversational Assistant (Gemini Chat)
* **Scope**: Natural language operator assistant with tool access to ClickHouse MCP and SQLite live state.
* **Deliverables**:
  - `internal/api/assistant/`: Chat endpoint streaming Gemini reasoning and citations via SSE.
  - `web/src/routes/index.tsx`: Interactive chat workspace with suggested prompts, ClickHouse query citations, and instant context awareness.
* **Verification**: Interactive query verification asserting ClickHouse analytical citations and accurate operational state responses.

---

### Feature 07: Smart Simulator, Protected Seed Data & Communications Hub
* **Scope**: Context-aware dynamic event simulator tab (`/simulator`), protected baseline reference seed data, and automated multi-channel stakeholder dispatches.
* **Deliverables**:
  - `cmd/seed/main.go`: Seed protected system titles (*Eclipse*, *Atlas*) and top vendors into ClickHouse & SQLite.
  - `web/src/routes/simulator.tsx`: Dynamic event generator allowing operators to trigger master revisions, QC failures on audio drift, or vendor re-conform deliveries.
  - `internal/agent/comms.go`: Autonomous generation of Vendor SLA notices, Public Social broadcasts (X/Twitter), and Executive Briefings.
  - `web/src/features/runs/`: Live workflow activity stream and communications view.
* **Verification**: Unattended demo cold start, dynamic event generation, and full closed-loop recovery.
