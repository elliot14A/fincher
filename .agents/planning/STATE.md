# Fincher — Live Operational State & Milestone Pointer

## Current Status Pointer
* **Active Milestone**: Feature 04: UI Scaffolding & Initial Operations Console
* **Active Phase**: `04-ui-scaffolding-console`
* **Phase Status**: READY TO EXECUTE
* **Next Milestone**: Feature 05: ClickHouse History & MCP Client
* **Timestamp**: 2026-08-22T14:15:00+05:30

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

- [ ] **Feature 04: UI Scaffolding & Initial Operations Console**
  - [ ] Initialize `web/` with Bun, Vite, Preact, TypeScript, Biome.
  - [ ] Configure `vite.config.ts` (Preact preset, Vanilla Extract plugin, TanStack Router plugin, `/api` proxy).
  - [ ] Configure `openapi-ts.config.ts` and generate SDK in `src/lib/api/generated/`.
  - [ ] Implement design tokens in `src/styles/theme.css.ts` and `src/app.css.ts`.
  - [ ] Build co-located UI primitives in `src/components/ui/` (`button/`, `badge/`, `modal/`, `table/`, `input/`, `drawer/`).
  - [ ] Implement `src/features/layout/` and root `src/routes/__root.tsx`.
  - [ ] Implement Launch Calendar (`src/features/calendar/`, `src/routes/index.tsx`).
  - [ ] Implement Vendor Directory (`src/features/vendors/`, `src/routes/vendors.tsx`).
  - [ ] Implement Territory Delivery Matrix (`src/features/deliveries/`, `src/routes/deliveries.tsx`).
  - [ ] Implement Lineage DAG Visualizer (`src/features/lineage/`, `src/routes/lineage/$id.tsx`).
  - [ ] Verification: `bun run typecheck`, `bun run lint`, full route navigation in browser.

- [ ] **Feature 05: ClickHouse History & MCP Client**
  - [ ] ClickHouse schema: `qc_events`, `delivery_events`, `asset_events`, `dependency_events`.
  - [ ] ClickHouse materialized views.
  - [ ] Remote MCP HTTP client wrapper in `pkg/mcp/`.
  - [ ] Verification: Real MCP tool calls against running `mcp-clickhouse` container.

- [ ] **Feature 06: Workflow Definitions & DAG Spec**
  - [ ] DAG runner spec in `internal/turso/workflow_defs/`.
  - [ ] GraphSpec models and `/api/workflows` endpoints.
  - [ ] Verification: Topological validity tests.

- [ ] **Feature 07: Execution Engine & Live SSE Runs (Backend + UI)**
  - [ ] In-memory topological DAG runner in `internal/engine/`.
  - [ ] Live SSE streaming endpoint `/api/runs/:id/stream`.
  - [ ] Frontend: `web/src/features/runs/`, `web/src/routes/runs/index.tsx`, `web/src/routes/runs/$id.tsx`.
  - [ ] Verification: Real-time step visualizer in browser.

- [ ] **Feature 08: 17-Node Palette & Multi-Agent Investigation (Backend + UI)**
  - [ ] 17 node types in `internal/engine/nodes/`.
  - [ ] Gemini Flash query agents & Gemini Pro Decision Node.
  - [ ] Frontend: Disagreement Panel highlighting Naive Rule vs Fincher Decision.
  - [ ] Verification: Hero scenario (*Eclipse* V13 + Vendor A audio drift $\rightarrow$ HOLD).

- [ ] **Feature 09: Docent NLQ Assistant & Causal Seeder (Backend + UI)**
  - [ ] Natural language operator assistant in `internal/assistant/docent.go`.
  - [ ] Frontend: `web/src/features/docent/` chat drawer.
  - [ ] Causal dataset seeder in `cmd/seed/`.
  - [ ] Production Go binary embedding `web/dist`.
  - [ ] Verification: Full unattended cold start and single-binary run.
