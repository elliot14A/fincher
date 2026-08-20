# Fincher — System Requirements Specification

## 1. Functional Requirements

### Ingestion & Event Model (`REQ-INGEST`)
* **REQ-INGEST-01**: The system must provide a REST endpoint `POST /api/v1/events` to ingest standard media lifecycle events (`MASTER_UPDATED`, `ASSET_UPDATED`, `PACKAGE_CREATED`, `QC_FAILED`, `QC_PASSED`, etc.).
* **REQ-INGEST-02**: All incoming events must be validated against the canonical `ProductionEvent` schema.
* **REQ-INGEST-03**: Events must be persisted synchronously to SQLite and ClickHouse (via MCP / store writer) before triggering investigation.

### Application State & Storage (`REQ-STATE`)
* **REQ-STATE-01**: SQLite must manage operational entities: `deliveries`, `packages`, `qc_results`, `incidents`, `human_approvals`, and `audit_logs`.
* **REQ-STATE-02**: SQLite must operate in WAL mode (`_journal_mode=WAL`) with foreign key enforcement and single-writer concurrency (`SetMaxOpenConns(1)`).
* **REQ-STATE-03**: State changes must be atomic, idempotent, and emit traceable audit log entries.

### ClickHouse MCP Integration (`REQ-MCP`)
* **REQ-MCP-01**: The Go backend must act as a remote MCP client communicating over HTTP transport (`FINCHER_MCP_URL`) with the official `ghcr.io/clickhouse/mcp-clickhouse:latest` container.
* **REQ-MCP-02**: ClickHouse database credentials must remain strictly isolated inside the MCP service environment.
* **REQ-MCP-03**: AI agents must access ClickHouse strictly through MCP read-only tools (`run_query`, `list_tables`, `list_databases`).

### Multi-Agent Investigation & Planning (`REQ-AGENT`)
* **REQ-AGENT-01**: The **Historian Agent** and **Dependency Agent** must execute in parallel goroutines upon receiving an anomaly event.
* **REQ-AGENT-02**: Historian Agent must query historical vendor defect rates and past check anomalies via ClickHouse MCP.
* **REQ-AGENT-03**: Dependency Agent must query master version lineage and identify all downstream localized packages affected by upstream changes via ClickHouse MCP.
* **REQ-AGENT-04**: The **Incident Analyst Agent** must synthesize findings into an `AnalystAssessment` with classification, severity, confidence, and root-cause analysis.
* **REQ-AGENT-05**: The **Action Planner Agent** must propose a structured `ActionPlan` with discrete candidate remediation actions (`INVALIDATE_PACKAGE`, `CREATE_RE_QC`, `HOLD_DELIVERY`, `RELEASE_DELIVERY`, `CREATE_INCIDENT`, `ESCALATE_VENDOR`).
* **REQ-AGENT-06**: All agent outputs must adhere strictly to structured JSON schemas using Google GenAI / ADK Go.

### Policy Gating (`REQ-POLICY`)
* **REQ-POLICY-01**: The policy engine must be a pure Go deterministic evaluator querying active operational rules from the database `policies` table.
* **REQ-POLICY-02**: Candidate actions must evaluate into `ALLOWED`, `DENIED`, or `HUMAN_APPROVAL_REQUIRED`.
* **REQ-POLICY-03**: AI proposals must never bypass or override policy thresholds.

### State Execution & Closed-Loop Feedback (`REQ-EXEC`)
* **REQ-EXEC-01**: The Go Executor must execute state mutations strictly for actions verified as `ALLOWED` (or approved by a human operator).
* **REQ-EXEC-02**: Every executed action must emit a downstream event (`PACKAGE_INVALIDATED`, `QC_STARTED`, `DELIVERY_RELEASED`, etc.).
* **REQ-EXEC-03**: Fincher must consume its own resulting downstream events, driving the delivery lifecycle to terminal resolution (`READY_TO_SHIP` or `HOLD`) without manual human intervention.

### Live Streaming & Operations Console (`REQ-STREAM` / `REQ-UI`)
* **REQ-STREAM-01**: The system must provide a Server-Sent Events (SSE) endpoint `GET /api/v1/stream` broadcasting real-time events, agent reasoning stages, policy decisions, and state updates.
* **REQ-STREAM-02**: The Operations Console UI must be served directly from the Go binary via `embed.FS`.

---

## 2. Non-Functional Requirements

* **REQ-ENV-01**: All environment variables must adhere to `FINCHER_{SERVICE}_{SETTING}`.
* **REQ-PERF-01**: Ingestion API response latency must be $< 20\text{ms}$.
* **REQ-AUDIT-01**: 100% of agent queries, decisions, and mutations must be recorded in immutable SQLite audit logs with correlated `trace_id`.
* **REQ-GO-01**: Codebase must pass `go build ./...`, `go vet ./...`, and race-detector tests with 0 warnings.
