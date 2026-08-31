---
name: gsd-verify
description: >-
  Executes comprehensive automated verification, code review, diagnostics, invariant checks,
  and integration tests for Fincher, compiling the evidence into .agents/planning/phases/XX/VERIFICATION.md.
---

# GSD Verify Skill

Use this skill to validate phase deliverables before shipping.

---

## Verification Pipeline

### 1. Static Analysis & Compilation
```bash
go build ./...
go vet ./...
```
*Verify zero compiler warnings, proper package boundaries, and strict type safety.*

### 2. Verification Guard Script
```bash
./.agents/scripts/verify.sh
```
*Checks:*
* `FINCHER_{SERVICE}_*` environment variable naming across configuration structs.
* MCP credential isolation (zero ClickHouse secrets in application memory).
* Planning state consistency (`PROJECT.md`, `REQUIREMENTS.md`, `ROADMAP.md`, `STATE.md`).

### 3. Unit & Concurrency Test Suite
```bash
go test -v -race ./...
```
*Validates:*
* Policy Engine evaluation rules against database `policies` table.
* SQLite schema migrations, foreign keys, and WAL mode single-writer concurrency.
* `context.Context` cancellation and goroutine leak safety.

### 4. Integration Diagnostics (ClickHouse MCP HTTP)
```bash
curl -sf http://127.0.0.1:8000/health
go test -v ./pkg/mcp/... -race
```
*Verifies read-only MCP queries and tool invocations.*

### 5. Architectural Review
* **Read-Only Agents**: Ensure AI sub-agents query strictly via MCP and NEVER perform direct SQL writes.
* **Policy Verification Gate**: Ensure all state mutations flow through policy verification and transactional software execution.
* **Closed-Loop Progression**: Verify executed actions emit downstream events.

### 6. Write Verification Report
Compile the results into `.agents/planning/phases/XX/VERIFICATION.md` with:
* Summary Verdict (`PASS` / `FAIL`)
* Package Test Results Table
* Invariant Verification Output
* Diagnostic Traces
