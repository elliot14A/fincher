# Phase 01: ClickHouse MCP Spike, Analytical Schemas & Turso State Store — Context

## Objective
Establish the foundational data layer and infrastructure for Fincher v2:
1. Denormalized ClickHouse OLAP schema with Materialized Views optimized for vendor trends, asset lineage, and failure rates.
2. Official ClickHouse MCP client connectivity validation over HTTP transport.
3. Turso / libSQL application store with Launch Calendar tables (`titles`, `packages`, `deliveries`), `policies` (CEL rules), `workflows`, `work_queue`, and `audit_trail`.
4. Synthetic causal seed dataset modeling LUME streaming titles (*Eclipse*, *Atlas*, *Orbit*, etc.) and vendors (`vendor_a` degradation).
5. Canonical Go domain models and Kong configuration (`FINCHER_{SERVICE}_*`).

---

## Design Decisions
1. **Module & Language**: `github.com/elliot14A/fincher`, Go 1.24+ standard library conventions.
2. **State Store**: Turso / libSQL (via modern Go client compatible with both local SQLite file / in-memory for testing and remote Turso HTTP database).
3. **Historical Store**: ClickHouse accessed by agents strictly via official `mcp-clickhouse` container (`http://127.0.0.1:8000/mcp`).
4. **Policy Format**: CEL expression strings stored in `policies` table with rule IDs (`R-017`, `R-021`).
5. **Causal Data Core**: Deterministic seed with realistic volume (historical QC/delivery events) and causal hooks (*Eclipse* V13 imminent premiere, `vendor_a` Spanish audio borderline drift).

---

## Invariants Checklist
- [ ] Read-Only Agent Invariant (Agents only read ClickHouse via MCP)
- [ ] MCP Credential Isolation Invariant (No ClickHouse credentials inside agent code)
- [ ] Mandatory `FINCHER_*` Env Naming Invariant
- [ ] Cold Cloud Run Ready (Serverless state and durable work queue)
- [ ] Deterministic CEL Policy Rule Storage
