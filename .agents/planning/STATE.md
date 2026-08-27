# Fincher — Live Operational State & Milestone Pointer

## Current Status Pointer
* **Active Milestone**: Feature 05: ClickHouse Ingestion Pipeline (zero LLM cost substrate)
* **Active Phase**: `05-clickhouse-ingestion-pipeline`
* **Phase Status**: PLANNED, READY TO EXECUTE
* **Next Milestone**: Feature 06: Dedicated Agents (ADK Go v2 graph)
* **Timestamp**: 2026-08-26T00:00:00+05:30

## Key Decision Log (2026-08-27)
* **Cut debounce/coalesce/dedup and the in-memory budget/concurrency gate from Feature 05
  entirely**, after a full design pass (Q&A + detailed plan documented in
  `.agents/reviews/14-event-ingestion-pipeline/20260827T171710Z.md`, now marked superseded).
  Rationale: this machinery was solving an agent-stampede/quota-exhaustion problem that doesn't
  exist yet in Feature 05 — there's no Feature 06 graph downstream to protect, and the primary
  event source (the burst simulator) hits `/api/events` directly with real HTTP requests, not
  something requiring in-memory aggregation. Building it now would have meant gating a
  `NoopHandoff` that does nothing — exactly the speculative infrastructure AGENTS.md's
  "Simplicity First" rule warns against.
* **What Feature 05 actually ships instead**: `POST /api/events` accepts a JSON array of
  CloudEvents and inserts each directly into ClickHouse via `internal/clickhouse/events`. No
  classification-based branching lives in this endpoint.
* **Debounce/coalesce/dedup is dropped, not deferred** — Feature 06's judges will fire per
  individual `ANOMALY_SIGNAL`/`ALLOCATION_REQUEST`/`OPERATOR_FORCED` event, one LLM call each,
  no batching stage ahead of them.
* **The daily-cap/concurrency budget gate is deferred to Feature 06**, and must be Turso-persisted
  (not in-memory) when built, since the service runs on Cloud Run with scale-to-zero — an
  in-memory counter would silently reset on every cold start.
* Reconciled `REQUIREMENTS.md` (`REQ-AGENT-01`, `REQ-AGENT-10`), `ROADMAP.md`, `PROJECT.md`, and
  both `phases/05-clickhouse-ingestion-pipeline/` and `phases/06-dedicated-agents/` docs to
  remove/correct now-stale `internal/ingest.*` references left over from the cut design.
* Fixed the ClickHouse cold-start bootstrap bug flagged in two prior code reviews
  (`.agents/reviews/13-.../` review #4 equivalent) — `internal/clickhouse/client.go` now runs
  `CREATE DATABASE IF NOT EXISTS fincher` against the `default` database before connecting to
  `fincher`, verified live against a container with the database dropped.
* Replaced the hand-rolled JSON-RPC/SSE `pkg/mcp` client with the official
  `github.com/modelcontextprotocol/go-sdk/mcp` + `google.golang.org/adk/v2/tool/mcptoolset`,
  resolving the doc/code mismatch flagged in `EVENT_INGESTION_DESIGN_AND_DOUBTS.md`.

## Key Decision Log (2026-08-26)
* Split the former single "ClickHouse MCP & Multi-Agent Engine" milestone into two clean
  phases: **05 (ingestion pipeline, no LLM)** and **06 (dedicated agents, ADK Go v2 graph)**.
  Docent moved to Feature 07, Simulator/Comms moved to Feature 08.
* Dropped the earlier "Workflow DAG Edition" execution model (Cloud Scheduler polling,
  fixed 17-node palette, `decision_node` with 4 hardcoded branches, strictly acyclic runs).
  Replaced with an event-driven, push-based model: static taxonomy filter → debounce/coalesce
  → scoped LLM judges → transactional executor. See `PROJECT.md` §3 and `REQUIREMENTS.md`
  (`REQ-EVT`, `REQ-AGENT`).
* Dropped all invented numeric business thresholds (sync-drift ms tolerance, vendor score
  margins, fixed scoring weights). Every judgment point is a single-scoped LLM judge call
  over Go-assembled evidence, never a hardcoded constant. The Policy Judge's reject→revise
  loop is bounded (configurable, default cap 3) — the one deliberate exception to "no loops
  within a run," justified as bounded self-correction rather than unbounded retries.
* Vendor rate cards (standard/rush rate + turnaround) are SQLite operational data
  (`Vendor` entity), distinct from ClickHouse's historical realized-performance rollups.
* ClickHouse migration review #3 (`.agents/reviews/13-clickhouse-schema-and-invocation-design/20260826T171940Z.md`)
  found `001_events.sql` **fails to apply as committed** — `TTL occurred_at + INTERVAL 2 YEAR
  DELETE` throws `BAD_TTL_EXPRESSION` because `occurred_at` is `DateTime64(3, 'UTC')` and
  ClickHouse's `TTL` clause requires `DateTime`/`Date` (fix: wrap in `toDateTime(occurred_at)`).
  Also found `003_vendor_metrics.sql` counts malformed/`UNKNOWN`-status QC events into
  `total_inspections` without excluding them from the accuracy denominator, silently scoring
  vendors with unparseable telemetry as 100% accurate — verified live via test insert, not yet
  fixed (user implements fixes manually per the review-only contract).

## 1. Milestone & Feature Completion Checklist

- [x] **Feature 01: Titles & Launch Calendar**
  - [x] Domain model `Title`, validation tags, status enum (`ON_TRACK`, `AT_RISK`, `HOLD`, `PROCESSING`, `SHIPPED`).
  - [x] Ent schema `title.go`, automigrations.
  - [x] Turso functional actions (`Create`, `Get`, `List`, `Update`, `Delete`) returning `Result[T]`.
  - [x] Echo v4 REST endpoints (`POST /api/titles`, `GET /api/titles`, `GET /api/titles/:id`, `PATCH /api/titles/:id`, `DELETE /api/titles/:id`).
  - [x] Verification: CRUD test suite and HTTP lifecycle test passing with `-race`.

- [x] **Feature 02: Masters, Packages & Vendors**
  - [x] Ent schemas: `Master` (`supersedes_version`), `Vendor`, `MediaPackage` (`component`, `language`, `version`, `derived_from_master_version`, `status`).
  - [x] Domain models and validation in `pkg/domain/models/`.
  - [x] Master insertion transactionally synchronizes `Title.current_master_version`.
  - [x] Turso functional actions for `masters/`, `packages/`, `vendors/` returning pure domain models.
  - [x] Echo REST handlers with modular `routes.go` per entity under `/api/*`.
  - [x] Verification: Unit tests for CRUD, FK constraints, staleness detection, and HTTP lifecycles passing with `-race`.

- [x] **Feature 03: Deliveries & Lineage / Dependencies**
  - [x] Ent schemas: `Delivery` (`country`, `status`, `target_date`), `Dependency` (`parent_id`, `child_id`).
  - [x] Domain models and validation in `pkg/domain/models/`.
  - [x] Pure functional actions in `internal/turso/deliveries/` and `internal/turso/dependencies/`.
  - [x] Echo REST endpoints: `GET /api/deliveries`, `GET /api/dependencies/graph/:title_id`.
  - [x] Verification: Lineage graph cycle detection tests and full HTTP lifecycle tests.

- [x] **Feature 04: UI Scaffolding & Operations Console**
  - [x] Initialized `web/` workspace with Bun, Vite, Preact, TypeScript, Biome, `@vanilla-extract/css`.
  - [x] Co-located UI primitives: `button/`, `badge/`, `modal/` (`Modal`, `DeleteModal`), `input/` (`FormField`, `TextInput`, `SelectInput`, `ImageUpload`).
  - [x] Dedicated entity creation modals in `src/features/*/components/modals/` for Titles, Vendors, Deliveries, and Packages.
  - [x] SQLite BLOB image storage (`internal/turso/uploads/`, `internal/api/uploads/`) with strict 1MB size limit and server-side MIME sniffing.
  - [x] Centered layouts and empty states on X and Y axes across all console views.
  - [x] Replaced 'territories' with 'markets', removed em dashes, and removed 'DAG' backend terminology from UI.
  - [x] Removed `+ New Package` button from Runs page.
  - [x] Verification: `bun test`, `bun run typecheck`, `biome check src`, `go test -race ./...` all passing.

- [ ] **Feature 05: ClickHouse Ingestion Pipeline** (zero LLM cost — see `.agents/planning/phases/05-clickhouse-ingestion-pipeline/`)
  - [x] ClickHouse migrations `001_events.sql`, `002_qc.sql`, `003_vendor_metrics.sql` conforming to CNCF CloudEvents v1.0 with subject as title slug and lowercase SQL.
  - [x] `internal/clickhouse/`: direct deterministic `database/sql` client + 120-day decay accuracy and metrics queries in domain packages.
  - [x] Agent-facing MCP integration: official ClickHouse MCP server (`mcp-clickhouse`) wired directly via Google ADK Go's native `mcptoolset` in Phase 06 with ClickHouse native `readonly=1` safety.
  - [x] `pkg/domain/models/event.go` & `event_payloads.go`: CloudEvents v1.0 schema, category classifier, and typed payloads (`QCPayload`, `MasterCutPayload`, `SLABreachPayload`, `PackageInvalidatedPayload`).
  - [x] `internal/api/events/`: `POST /api/events` batch ingestion endpoint (enforces CloudEvent array only, inserts directly into ClickHouse).
  - [ ] Rate Limiting: Parked for later (to be implemented with Turso persistent store to survive Cloud Run scale-to-zero).
  - [x] Verification: full pipeline tests against live ClickHouse container (`internal/api/events`, `internal/clickhouse/events`, `internal/clickhouse/vendors`, `pkg/mcp`).

- [ ] **Feature 06: Dedicated Agents (ADK Go v2 Graph)** (see `.agents/planning/phases/06-dedicated-agents/`)
  - [ ] `internal/agent/graph.go`: ADK Go v2 workflow wiring.
  - [ ] `triage_judge.go`, `historian.go`, `lineage.go`, `optimizer.go`, `policy_judge.go`, `executor.go`.
  - [ ] `vendor_scoring.go`, `vendor_judge.go` + `Vendor` ent schema rate-card extension.
  - [ ] `internal/api/runs/` SSE stream + `internal/api/investigations/` operator triggers.
  - [ ] `web/src/features/runs/`: live `@xyflow/react` investigation graph.
  - [ ] Verification: stubbed-LLM node tests, fixture full-graph run, live Gemini demo run, frontend typecheck/biome.

- [ ] **Feature 07: Docent Conversational Assistant (Gemini Chat)**
  - [ ] Natural language operator assistant streaming reasoning over ClickHouse & SQLite.
  - [ ] Frontend: Interactive chat view in `web/src/routes/index.tsx` with ClickHouse query citations.
  - [ ] Verification: Query accuracy and live citation verification.

- [ ] **Feature 08: Smart Simulator, Protected Seed Data & Communications Hub**
  - [ ] Context-aware dynamic event generator tab (`web/src/routes/simulator.tsx`).
  - [ ] Protected reference baseline seed data in `cmd/seed/main.go`.
  - [ ] Autonomous generation of Vendor SLA notices, Public Social broadcasts (X/Twitter), and Executive Briefings.
  - [ ] Production Go binary embedding `web/dist`.
  - [ ] Verification: Full unattended cold start and single-binary run.
