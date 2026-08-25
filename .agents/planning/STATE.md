# Fincher — Live Operational State & Milestone Pointer

## Current Status Pointer
* **Active Milestone**: Feature 05: ClickHouse MCP & Multi-Agent Investigation Engine
* **Active Phase**: `05-clickhouse-multi-agent`
* **Phase Status**: READY TO PLAN & EXECUTE
* **Next Milestone**: Feature 06: Docent Conversational Assistant (Gemini Chat)
* **Timestamp**: 2026-08-24T21:08:00+05:30

---

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

- [ ] **Feature 05: ClickHouse MCP & Multi-Agent Investigation Engine**
  - [ ] ClickHouse schema: `qc_events`, `delivery_events`, `asset_events`, `vendor_performance`.
  - [ ] Remote MCP HTTP client wrapper in `pkg/mcp/` connecting to `mcp-clickhouse`.
  - [ ] Historian Sub-Agent querying ClickHouse MCP for vendor defect history and sync drift logs.
  - [ ] Lineage Sub-Agent checking affected market deliveries in SQLite.
  - [ ] Multi-Factor Vendor Optimizer balancing Speed, Quality, and Rates.
  - [ ] Verification: Live MCP tool calls against `mcp-clickhouse` and automated multi-agent run.

- [ ] **Feature 06: Docent Conversational Assistant (Gemini Chat)**
  - [ ] Natural language operator assistant streaming reasoning over ClickHouse & SQLite.
  - [ ] Frontend: Interactive chat view in `web/src/routes/index.tsx` with ClickHouse query citations.
  - [ ] Verification: Query accuracy and live citation verification.

- [ ] **Feature 07: Smart Simulator, Protected Seed Data & Communications Hub**
  - [ ] Context-aware dynamic event generator tab (`web/src/routes/simulator.tsx`).
  - [ ] Protected reference baseline seed data in `cmd/seed/main.go`.
  - [ ] Autonomous generation of Vendor SLA notices, Public Social broadcasts (X/Twitter), and Executive Briefings.
  - [ ] Production Go binary embedding `web/dist`.
  - [ ] Verification: Full unattended cold start and single-binary run.
