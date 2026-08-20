# Phase 01: Core Architecture, Setup, Domain Models & Storage — Context

## Objective
Establish the durable foundation for Fincher: Go module initialization, canonical domain models, event vocabulary, Kong configuration with `FINCHER_{SERVICE}_*` tags, SQLite state store with schema migrations and WAL mode, and declarative operational policies.

---

## Design Decisions
1. **Module Name**: `github.com/elliot14A/fincher`
2. **Go Version**: Go 1.24+ (or Go 1.26 on dev machine)
3. **Database Architecture**:
   - SQLite for application transactional state (`deliveries`, `packages`, `qc_results`, `incidents`, `human_approvals`, `audit_logs`).
   - SQLite connection uses `_journal_mode=WAL`, `_busy_timeout=5000`, `_foreign_keys=ON`, and `SetMaxOpenConns(1)`.
4. **Configuration Tagging**:
   - Uses Kong struct tags with explicit `env:"FINCHER_{SERVICE}_{SETTING}"` bindings.
   - Database credentials are not loaded into application config; agent operations communicate strictly via `FINCHER_MCP_URL`.
5. **Database-Backed Policies**:
   - `policies` table in SQLite stores operational rule constraints (`rule_id`, `condition_type`, `threshold`, `action`, `requires_approval`, `enabled`).

---

## Invariants Checklist
- [ ] Read-Only Agent Invariant
- [ ] MCP Credential Isolation Invariant
- [ ] Mandatory `FINCHER_*` Env Naming Invariant
- [ ] Pure Go Deterministic Policy Gating Invariant
