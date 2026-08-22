# Fincher — System Requirements Specification (Workflow DAG Edition)

## 1. Functional Requirements

### ClickHouse Analytical History & MCP (`REQ-CH`)
* **REQ-CH-01**: ClickHouse manages denormalized event tables (`qc_events`, `delivery_events`, `asset_events`, `dependency_events`, `node_execution_events`).
* **REQ-CH-02**: ClickHouse provides 4 precomputed Materialized Views: `vendor_check_failure_rates`, `package_lineage`, `redelivery_counts`, `recent_master_changes`.
* **REQ-CH-03**: All runtime analytical reads by workflow nodes route strictly through the official `mcp-clickhouse` server over HTTP transport (`FINCHER_MCP_URL`).

### Turso / libSQL Application State (`REQ-TURSO`)
* **REQ-TURSO-01**: Turso stores the delivery domain: `titles`, `masters`, `packages`, `vendors`, `deliveries`, `dependencies`.
* **REQ-TURSO-02**: Turso stores the workflow execution domain: `workflow_definitions`, `workflow_runs`, `node_executions`, `node_inputs`, `node_outputs`, `query_log`, `decisions`, `executed_actions`, `notifications`, `budget_counters`.
* **REQ-TURSO-03**: Database schema is managed via type-safe Ent ORM schemas in `internal/turso/ent/schema/` with automated migrations.

### Workflow DAG & Node Palette (`REQ-DAG`)
* **REQ-DAG-01**: Workflows are directed acyclic graphs (DAG) composed from a fixed palette of 17 node types.
* **REQ-DAG-02**: Trigger & Gate: `schedule_trigger` and `delta_gate` (exits at zero LLM cost if no new events).
* **REQ-DAG-03**: Query Nodes: `vendor_reliability_query`, `lineage_query`, `redelivery_query`, `recent_master_change_query`, `time_to_premiere`.
* **REQ-DAG-04**: Agent Nodes (Gemini Flash): `vendor_risk_agent`, `dependency_impact_agent`, `assessment_agent` (Pro).
* **REQ-DAG-05**: Decision Node (Gemini Pro): Combines deterministic query inputs and agent synthesis to select a branch (`HOLD`, `RE_QC`, `RELEASE`, `NONE`) with recorded rationale.
* **REQ-DAG-06**: Action & Sink: `hold_delivery_action`, `request_requc_action`, `release_delivery_action`, `stakeholder_notice_action` (drafted notice), `release_notice_action`, `event_emitter`.

### Backend API & Real-Time Stream (`REQ-API`)
* **REQ-API-01**: Workflow management: `GET /api/workflows`, `GET /api/workflows/{id}`, `POST /api/workflows`, `PATCH /api/workflows/{id}`, `POST /api/workflows/{id}/run`.
* **REQ-API-02**: Run inspection: `GET /api/runs/{id}`, `GET /api/runs/{id}/stream` (SSE live step execution).
* **REQ-API-03**: Title catalog & CRUD: `GET /api/titles`, `GET /api/titles/{id}`, `POST /api/titles`, `PATCH /api/titles/{id}`, `DELETE /api/titles/{id}`.
* **REQ-API-04**: Master & Package & Vendor CRUD: `POST /api/masters`, `GET /api/masters`, `POST /api/packages`, `GET /api/packages`, `PATCH /api/packages/{id}`, `POST /api/vendors`, `GET /api/vendors`.
* **REQ-API-05**: Read-first Docent query: `POST /api/query`, `GET /api/query/{session}/stream`.
* **REQ-API-06**: Palette & Budget: `GET /api/node-palette`, `GET /api/budget`.
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
* **REQ-UI-09**: **Studio Pipeline Live Run Inspector**: Real-time SSE node execution visualizer, Gemini agent reasoning logs, and policy verification badges.
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
