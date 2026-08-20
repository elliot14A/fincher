# Fincher — Project Roadmap

## Milestone Overview

```text
[ Phase 01: Foundation & Storage ] ──► [ Phase 02: MCP & Multi-Agent ]
                                                │
                                                ▼
[ Phase 04: UI & Live Stream ] ◄── [ Phase 03: Policy Engine & Executor ]
```

---

## Phases & Deliverables

### Phase 01: Core Architecture, Setup, Domain Models & Storage
* **Status**: `IN_PROGRESS`
* **Goal**: Establish durable project setup, canonical domain types, SQLite state and policy tables with WAL mode, and configuration parsed via Kong.
* **Deliverables**:
  - `pkg/domain/types.go`: Core domain types (Delivery, Package, QCResult, Incident, Action, Evidence).
  - `pkg/events/events.go`: Canonical event types and payload schemas.
  - `pkg/policy/schema.go`: Policy schema structs.
  - `internal/config/config.go`: Kong configuration adhering to `FINCHER_{SERVICE}_*`.
  - `internal/store/sqlite.go`: SQLite store with schema migrations (operational entities & policies table) and WAL mode.
  - `data/seed/media_events.json`: Realistic media operations dataset.
* **Definition of Done (DoD)**:
  - SQLite table migrations succeed and pass integration tests with race detector.
  - Configuration parses all `FINCHER_*` environment variables cleanly.

---

### Phase 02: ClickHouse MCP Client & Multi-Agent Investigation
* **Status**: `PENDING`
* **Goal**: Implement remote HTTP MCP client and Google ADK Go parallel sub-agents (Historian, Dependency, Incident Analyst, Action Planner).
* **Deliverables**:
  - `pkg/mcp/client.go`: Remote MCP client communicating over HTTP transport with `ghcr.io/clickhouse/mcp-clickhouse:latest`.
  - `internal/agent/historian.go`: Historian agent querying vendor failure rates via MCP.
  - `internal/agent/dependency.go`: Dependency agent querying asset version lineage via MCP.
  - `internal/agent/analyst.go`: Incident Analyst synthesizing evidence into severity/classification.
  - `internal/agent/planner.go`: Action Planner generating structured candidate `ActionPlan`.
  - `internal/agent/orchestrator.go`: Multi-agent parallel coordination via goroutines / ADK Go.
* **Definition of Done (DoD)**:
  - Live MCP integration tests pass against ClickHouse MCP container.
  - Parallel sub-agents execute concurrently and emit structured typed JSON outputs.

---

### Phase 03: Policy Engine, Executor & Closed-Loop Pipeline
* **Status**: `PENDING`
* **Goal**: Implement pure Go deterministic policy evaluator backed by SQLite policies, state executor, and closed-loop re-evaluation loop.
* **Deliverables**:
  - `internal/policy/engine.go`: Deterministic rule evaluator mapping DB-backed policies to action gating (`ALLOWED`, `DENIED`, `HUMAN_APPROVAL_REQUIRED`).
  - `internal/executor/executor.go`: State mutator executing permitted actions and emitting downstream events.
  - `internal/pipeline/pipeline.go`: Closed-loop orchestration pipeline driving workflows to resolution (`READY_TO_SHIP` or `HOLD`).
* **Definition of Done (DoD)**:
  - 100% unit test coverage across all policy rules.
  - End-to-end event transition: `MASTER_UPDATED` $\rightarrow$ `PACKAGE_INVALIDATED` $\rightarrow$ `RE_QC_REQUESTED` $\rightarrow$ `QC_PASSED` $\rightarrow$ `DELIVERY_RELEASED`.

---

### Phase 04: REST API, SSE Live Stream & Operations Console UI
* **Status**: `PENDING`
* **Goal**: Build Echo HTTP handlers for ingestion and approvals, SSE real-time broadcasting, and embedded cinematic dark Operations Console UI.
* **Deliverables**:
  - `internal/api/server.go`: Echo HTTP server with middleware and CORS.
  - `internal/api/events.go`: `POST /api/v1/events` ingestion handler.
  - `internal/api/approvals.go`: `POST /api/v1/approvals/:id/resolve` human review handler.
  - `internal/api/stream.go`: `GET /api/v1/stream` non-blocking SSE broadcaster.
  - `web/`: Cinematic Operations Console single-page application served via `embed.FS`.
  - `cmd/fincher/main.go`: Server binary entry point.
* **Definition of Done (DoD)**:
  - Ingestion endpoint validates and processes events end-to-end.
  - UI displays real-time agent investigations, policy decisions, and audit timeline without page refresh.
