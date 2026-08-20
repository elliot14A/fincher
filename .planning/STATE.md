# Fincher — Project State

## Current Position
* **Active Milestone**: `Milestone 01: System Foundation & Setup`
* **Active Phase**: `01-foundation`
* **Phase Status**: `IN_PROGRESS`
* **Last Updated**: `2026-08-18`

---

## Active Phase Checklist (`01-foundation`)

- [ ] **Setup & Workflow Revamp**: Adopt GSD Core spec-driven development with durable `.planning/` state artifacts.
- [ ] **MCP Infrastructure**: Run official ClickHouse MCP server (`ghcr.io/clickhouse/mcp-clickhouse:latest`) over HTTP transport with isolated credentials in Docker Compose.
- [ ] **Policies Storage**: Define operational policies table in SQLite schema.
- [ ] **Realistic Datasets**: Establish synthetic media operations events in `data/seed/media_events.json`.
- [ ] **Canonical Domain Models**: Re-scaffold `pkg/domain/types.go` and `pkg/events/vocabulary.go`.
- [ ] **Configuration Layer**: Implement `internal/config/config.go` with strict `FINCHER_{SERVICE}_*` Kong struct tags.
- [ ] **SQLite State Storage**: Implement `internal/store/sqlite.go` with WAL mode and unit tests.
- [ ] **Phase 01 Verification**: Complete Phase 01 verification suite and generate `VERIFICATION.md`.

---

## Key Decisions Record
1. **GSD Core Workflow**: All development follows the 5-step phase loop (`DISCUSS` $\rightarrow$ `PLAN` $\rightarrow$ `EXECUTE` $\rightarrow$ `VERIFY` $\rightarrow$ `SHIP`) to prevent context rot.
2. **ClickHouse MCP via HTTP Transport**: The Go application communicates with ClickHouse strictly via the remote MCP HTTP transport endpoint (`FINCHER_MCP_URL=http://127.0.0.1:8000/mcp`), isolating credentials in `docker-compose.yml`.
3. **Environment Naming Standard**: Every environment variable adheres to `FINCHER_{SERVICE}_{SETTING}`.
4. **Single Source of Mutation**: AI agents are strictly read-only analytical tools; state mutations are performed solely by the Go executor upon deterministic policy verification.

---

## Blocker Log
* *None. Local ClickHouse (`:8123`) and official ClickHouse MCP server (`:8000`) containers are healthy and verified.*

---

## Next Steps
1. Create `01-foundation` phase execution artifacts (`.planning/phases/01-foundation/CONTEXT.md`, `PLAN.md`).
2. Scaffold canonical domain contracts (`pkg/domain`, `pkg/events`, `pkg/policy`).
3. Scaffold configuration (`internal/config`) and SQLite store (`internal/store`).
