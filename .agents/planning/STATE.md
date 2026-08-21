# Fincher — Live Operational State & Milestone Pointer

## Current Status Pointer
* **Active Milestone**: Feature 03: Deliveries & Lineage / Dependencies
* **Active Phase**: `03-deliveries-lineage`
* **Phase Status**: COMPLETED
* **Next Milestone**: Feature 04: ClickHouse History & MCP Client
* **Timestamp**: 2026-08-21T22:28:00+05:30

---

## 1. Milestone & Feature Completion Checklist

- [x] **Feature 01: Titles & Launch Calendar**
  - [x] Domain model `Title`, validation tags, status enum (`ON_TRACK`, `AT_RISK`, `HOLD`, `PROCESSING`, `SHIPPED`).
  - [x] Ent schema `title.go`, automigrations.
  - [x] Turso functional actions (`Create`, `Get`, `List`, `Update`, `Delete`) returning `Result[T]`.
  - [x] Echo v4 REST endpoints (`POST /titles`, `GET /titles`, `GET /titles/:id`, `PATCH /titles/:id`, `DELETE /titles/:id`).
  - [x] Verification: CRUD test suite and HTTP lifecycle test passing with `-race`.

- [x] **Feature 02: Masters, Packages & Vendors**
  - [x] Ent schemas: `Master` (`supersedes_version`), `Vendor`, `MediaPackage` (`component`, `language`, `version`, `derived_from_master_version`, `status`).
  - [x] Domain models and validation in `pkg/domain/models/`.
  - [x] Master insertion transactionally synchronizes `Title.current_master_version`.
  - [x] Turso functional actions for `masters/`, `packages/`, `vendors/` returning pure domain models.
  - [x] Echo REST handlers with modular `routes.go` per entity.
  - [x] Verification: Unit tests for CRUD, FK constraints, staleness detection, and HTTP lifecycles passing with `-race`.

- [x] **Feature 03: Deliveries & Lineage / Dependencies**
  - [x] Ent schemas: `Delivery` (`country`, `status`, `target_date`), `Dependency` (`parent_id`, `child_id`).
  - [x] Domain models and validation in `pkg/domain/models/`.
  - [x] Pure functional actions in `internal/turso/deliveries/` and `internal/turso/dependencies/`.
  - [x] Echo REST endpoints: `GET /deliveries`, `GET /dependencies/graph/:title_id`.
  - [x] Verification: Lineage graph cycle detection tests and full HTTP lifecycle tests.

- [ ] **Feature 04: ClickHouse History & MCP Client**
  - [ ] ClickHouse schema: `qc_events`, `delivery_events`, `asset_events`, `dependency_events`.
  - [ ] ClickHouse materialized views.
  - [ ] Remote MCP HTTP client wrapper in `pkg/mcp/`.
  - [ ] Verification: Real MCP tool calls against running `mcp-clickhouse` container.

- [ ] **Feature 05: Core DAG Workflow Engine & 17-Node Palette**
  - [ ] DAG runner in `internal/workflow/dag/`.
  - [ ] Implementation of all 17 nodes.
  - [ ] Verification: Deterministic mock run passing in under 5 seconds.

- [ ] **Feature 06: Gemini LLM Runtime & ADK Agents**
  - [ ] Google GenAI / ADK Go orchestration in `internal/agent/`.
  - [ ] Agents: `VendorRiskAgent`, `DependencyImpactAgent`, `AssessmentAgent`, `DocentAgent`.
  - [ ] Verification: Structured JSON schema validation and fallback tests.

- [ ] **Feature 07: Pure Go Policy Engine & State Mutations**
  - [ ] Policy evaluator in `internal/policy/`.
  - [ ] Executor in `internal/executor/`.
  - [ ] Verification: Policy unit tests with deterministic outcomes.

- [ ] **Feature 08: Operations Console UI & Real-Time Stream**
  - [ ] Embedded web UI in `web/`.
  - [ ] SSE stream endpoint in `internal/api/stream.go`.
  - [ ] Verification: End-to-end run observed live via SSE stream.

- [ ] **Feature 09: Causal Media Seeder & Hackathon Final Polish**
  - [ ] Seeder binary in `cmd/seed/`.
  - [ ] 250k+ historical QC events into ClickHouse.
  - [ ] Hero demonstration run.
  - [ ] Verification: Full end-to-end validation script.
