# Phase 01: Core Architecture, Setup, Domain Models & Storage — Plan

## Tasks

### Task 1: Scaffolding Domain Models & Event Vocabulary
* **Files**: `pkg/domain/types.go`, `pkg/events/events.go`, `pkg/policy/schema.go`, `pkg/telemetry/logger.go`
* **Details**: Define canonical types for Deliveries, Packages, QC Results, Incidents, Actions, Evidence Bundles, and Curated Events.
* **Verification**: Compiles cleanly with zero compiler warnings.

### Task 2: Configuration & Logging Scaffolding
* **Files**: `internal/config/config.go`, `pkg/telemetry/logger.go`
* **Details**: Implement Kong struct-tag parser binding to `FINCHER_{SERVICE}_*` environment variables and structured `slog` logger with `trace_id` support.
* **Verification**: Struct tags pass automated regex check for `FINCHER_*` format.

### Task 3: SQLite Storage Implementation & Tests
* **Files**: `internal/store/sqlite.go`, `internal/store/sqlite_test.go`
* **Details**: Implement SQLite table migrations and CRUD operations for deliveries, packages, qc results, incidents, human approvals, and audit logs.
* **Verification**: Run `go test -v ./internal/store/... -race` verifying 100% pass rate in WAL mode.

### Task 4: Phase Verification Report
* **Files**: `.planning/phases/01-foundation/SUMMARY.md`, `.planning/phases/01-foundation/VERIFICATION.md`
* **Details**: Run `./.agents/scripts/verify.sh` and document test evidence.
