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
* **REQ-API-01**: Workflow management: `GET /workflows`, `GET /workflows/{id}`, `POST /workflows`, `PATCH /workflows/{id}`, `POST /workflows/{id}/run`.
* **REQ-API-02**: Run inspection: `GET /runs/{id}`, `GET /runs/{id}/stream` (SSE live step execution).
* **REQ-API-03**: Title catalog & CRUD: `GET /titles`, `GET /titles/{id}`, `POST /titles`, `PATCH /titles/{id}`, `DELETE /titles/{id}`.
* **REQ-API-04**: Master & Package & Vendor CRUD: `POST /masters`, `GET /masters`, `POST /packages`, `GET /packages`, `PATCH /packages/{id}`, `POST /vendors`, `GET /vendors`.
* **REQ-API-05**: Read-first Docent query: `POST /query`, `GET /query/{session}/stream`.
* **REQ-API-06**: Palette & Budget: `GET /node-palette`, `GET /budget`.
* **REQ-API-07**: Deliveries CRUD: `GET /deliveries`, `GET /deliveries/{id}`, `POST /deliveries`, `PATCH /deliveries/{id}`, `DELETE /deliveries/{id}`.
* **REQ-API-08**: Dependency graph & lineage: `POST /dependencies`, `GET /dependencies`, `GET /dependencies/graph/{title_id}`, `DELETE /dependencies/{id}`.

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
