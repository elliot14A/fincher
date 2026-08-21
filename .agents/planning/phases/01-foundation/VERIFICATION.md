# Phase 01 — Foundation — Verification Report

**Date**: 2026-08-21
**Verdict**: **FAIL** (expected — phase is `IN_PROGRESS`, only Task 1 partially done)

---

## 1. Static Analysis & Compilation

| Check | Result |
|---|---|
| `go build ./...` | ✅ PASS |
| `go vet ./...` | ✅ PASS |
| `gofmt -l .` | ❌ FAIL — `pkg/domain/models/workflow.go` is not gofmt-clean (misaligned struct-tag/const columns) |

## 2. Verification Guard Script (`.agents/scripts/verify.sh`)

| Step | Result |
|---|---|
| [1/5] Planning artifacts present | ✅ PASS |
| [2/5] Seed data (`data/seed/media_events.json`) | ❌ FAIL — file doesn't exist (Task 5 not started) |
| [3/5] ClickHouse + MCP containers reachable | ⚠️ SKIPPED (script exits at step 2) |
| [4/5] Go build/vet/test | ⚠️ SKIPPED |
| [5/5] `FINCHER_*` env-var naming | ⚠️ SKIPPED — script only scans `internal/config/`, but the real config package is `internal/config/` (path mismatch, would always short-circuit as "no dir" even once implemented) |

Script halts on first failure by design (`set -e`), so steps 3–5 never ran.

## 3. Unit & Concurrency Test Suite

```
go test -v -race ./...
```
✅ PASS — 2 tests (`TestTitle_ValidationAndImminentLaunch`, `TestPackage_StaleMaster`), no race reports.

**Coverage gap**: `Delivery.Validate()`, `Event.Validate()`, `Config.Validate()`, and `WorkflowDefinition.Validate()` all exist but have zero test coverage.

## 4. Integration Diagnostics (ClickHouse MCP HTTP)

Not applicable yet — `pkg/mcp/client.go` (Task 3) hasn't been created. `docker compose config` validates the compose file syntactically (see issues below), but containers were not started for this pass.

## 5. Architectural Review

No engine/agent/action code exists yet (Milestone 2+), so the read-only-agent and deterministic-policy-gate invariants aren't testable at this phase.

---

## Issues Found

1. **`data/clickhouse/schema.sql` is a directory, not a file.** It's empty (`find` shows no children). Almost certainly created by an errant `mkdir -p data/clickhouse/schema.sql` instead of `mkdir -p data/clickhouse && touch schema.sql`. Also, per `ROADMAP.md` Milestone 1, schema should live in `migrations/clickhouse/*.sql` as individual per-table files — the current `data/clickhouse/` path doesn't match the documented target layout at all, and `migrations/` doesn't exist yet.

2. **`docker-compose.yml` has a dead top-level named volume.**
   ```yaml
   volumes:
     volumes:
   ```
   This declares an unused named volume literally called `volumes`. The clickhouse service actually bind-mounts `./volumes/clickhouse:/var/lib/clickhouse` (a host path, matching the new `.gitignore` entry `volumes/`), so the named-volume stanza is copy-paste cruft and should be deleted.

3. **`verify.sh` checks the wrong config path.** Step 5 greps `internal/config/` for `env:"FINCHER_...` tags, but the actual config lives in `internal/config/config.go` using Kong tags (`kong:"...env='FINCHER_...'"`) — different path *and* different tag key (`env:` vs `kong:"...env='...'"`). As written this check will silently no-op forever ("Config directory will enforce ... upon scaffolding") instead of ever validating the real file.

4. **`verify.sh` checks for a seed artifact that doesn't match the roadmap.** It requires `data/seed/media_events.json`, but `ROADMAP.md` / `STATE.md` Task 5 specifies `cmd/seed/main.go` + `data/seed/lume.go` generating ClickHouse rows and Turso records — not a static JSON file. One of the two specs is stale.

5. **`workflow.go` fails `gofmt`.** The `NodeType` const block has inconsistent alignment around the "Query Nodes" group (extra spaces before `NodeType =`), and it's the only non-gofmt-clean file in the diff.

6. **Test coverage is thin relative to surface area.** 4 of the 6 model/config files with `Validate()` methods (`Delivery`, `Event`, `Config`, `WorkflowDefinition`) have no tests at all — including the `WorkflowDefinition.Validate()` that's central to this phase's new DAG model.

7. **Empty scaffolding dirs** `cmd/fincher/` and `cmd/seed/` currently contain no Go files — expected at this point (Tasks 3–5 not started), just flagging so it isn't mistaken for a build gap once `go build ./...` starts including them.

None of these block continued work, but #1–#4 are worth fixing before Task 2 (migrations) and Task 5 (seeder) land, since both depend on the exact paths these issues touch.
