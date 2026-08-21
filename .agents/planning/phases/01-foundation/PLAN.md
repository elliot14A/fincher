# Phase 01: Milestone 0 (Runtime Spike) & Milestone 1 (Data Foundation & Migrations) — Plan

## Tasks

### Task 1: Domain Models (Pure Go, Zero Dependencies)
* **Files**:
  - `pkg/domain/models/title.go` (Titles, Premiere Dates, Launch Statuses)
  - `pkg/domain/models/package.go` (Packages, Asset Versions, Vendors)
  - `pkg/domain/models/delivery.go` (Deliveries, Territory mapping)
  - `pkg/domain/models/events.go` (ClickHouse event schemas)
  - `pkg/domain/models/workflow.go` (DAG GraphSpec, Node Types, Run, NodeExecution, Decision, Action, Notification)
  - `internal/config/config.go` (Kong Config with `FINCHER_{SERVICE}_*`)
* **Details**: Update models to match the 17-node DAG palette and drop old CEL remnants.
* **Verification**: `go test -v ./pkg/domain/...` passes 100%.

### Task 2: ClickHouse & Turso SQL Migrations (Individual Files)
* **Files**:
  - `migrations/clickhouse/001_qc_events.sql`
  - `migrations/clickhouse/002_delivery_events.sql`
  - `migrations/clickhouse/003_asset_events.sql`
  - `migrations/clickhouse/004_dependency_events.sql`
  - `migrations/clickhouse/005_node_execution_events.sql`
  - `migrations/clickhouse/006_mv_vendor_check_failure_rates.sql`
  - `migrations/clickhouse/007_mv_package_lineage.sql`
  - `migrations/clickhouse/008_mv_redelivery_counts.sql`
  - `migrations/clickhouse/009_mv_recent_master_changes.sql`
  - `migrations/turso/001_titles.sql`
  - `migrations/turso/002_masters.sql`
  - `migrations/turso/003_packages.sql`
  - `migrations/turso/004_vendors.sql`
  - `migrations/turso/005_deliveries.sql`
  - `migrations/turso/006_dependencies.sql`
  - `migrations/turso/007_workflow_definitions.sql`
  - `migrations/turso/008_workflow_runs.sql`
  - `migrations/turso/009_node_executions.sql`
  - `migrations/turso/010_node_inputs_outputs.sql`
  - `migrations/turso/011_query_log.sql`
  - `migrations/turso/012_decisions.sql`
  - `migrations/turso/013_executed_actions.sql`
  - `migrations/turso/014_notifications.sql`
  - `migrations/turso/015_budget_counters.sql`
  - `migrations/turso/016_query_sessions.sql`
* **Details**: Clean standard migration files, one per table and view.
* **Verification**: SQL syntax validation.

### Task 3: ClickHouse MCP HTTP Client & Turso Client
* **Files**:
  - `pkg/mcp/client.go`, `pkg/mcp/client_test.go`
  - `internal/turso/client.go`, `internal/turso/client_test.go`
* **Details**: MCP client for `mcp-clickhouse` (`run_query`, `list_tables`). Turso client with migration runner and transactional execution.
* **Verification**: Live MCP connection test against local Docker container + Turso SQLite migration run.

### Task 4: Causal Seeder for LUME
* **Files**: `cmd/seed/main.go`, `data/seed/lume.go`
* **Details**: Seed 5–7 LUME titles, packages, degrading `vendor_a` historical QC logs (with enough volume for MVs), and the canonical Studio Pipeline DAG template.
* **Verification**: `go run ./cmd/seed` seeds both databases cleanly and ClickHouse MVs return calculated rates.
