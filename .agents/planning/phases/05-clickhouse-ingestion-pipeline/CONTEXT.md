# Phase 05 Context: ClickHouse Ingestion Pipeline

## Goal
Stand up the zero-LLM substrate every later agent depends on: ClickHouse schema migrations conforming to CNCF CloudEvents v1.0, two distinct ClickHouse access paths (direct database/sql client vs. agent-facing MCP client), the static CloudEvent taxonomy, and a direct batch ingestion endpoint. **No Gemini/ADK call happens anywhere in this phase.**

**Scope cut (2026-08-27)**: debounce/coalesce/dedup and the daily-budget/concurrency gate were designed in detail (see `.agents/reviews/14-event-ingestion-pipeline/`) and then deliberately dropped from this phase as unwarranted complexity for the actual scale involved (simulator bursts hit `/api/events` directly; no downstream Feature 06 graph exists yet to protect from stampede). The budget/concurrency gate is deferred to Feature 06 — when it's built, it must be Turso-persisted, not in-memory, to survive Cloud Run scale-to-zero between requests. Debounce/coalesce is not planned to return at all; Feature 06's judges fire per individual event.

## Objectives
* ClickHouse migrations reflecting the CloudEvents v1.0 specification (`fincher.events`, `fincher.qc`, `fincher.vendor_metrics`).
* `internal/clickhouse/`: direct, deterministic Go `database/sql` client for facts (accuracy queries, vendor rollups) — never touched by agent tool-calling code. Self-bootstraps (`CREATE DATABASE IF NOT EXISTS fincher`) so a cold/fresh ClickHouse instance works with zero manual setup.
* `pkg/mcp/`: agent-facing MCP client built on the official `modelcontextprotocol/go-sdk` + `google.golang.org/adk/v2/tool/mcptoolset` — the *only* ClickHouse access path available to LLM tool-calling code, with ClickHouse native `readonly=1` safety.
* `pkg/domain/models/`: CloudEvents v1.0 model, static taxonomy classifying every `type` into one of five categories, and typed per-event payload structs (`QCPayload`, `MasterCutPayload`, `SLABreachPayload`, `PackageInvalidatedPayload`). This is ontology (what an event structurally is), never a numeric business threshold.
* `POST /api/events`: batch ingestion entrypoint — accepts a JSON array of CloudEvents, validates, writes each directly to ClickHouse. No classification-based branching lives here.

## Architectural Decisions
1. **ClickHouse schema** (`migrations/clickhouse/`):
   - `001_events.sql`: root `fincher.events` (`MergeTree`, `partition by toYYYYMM(time)`, `order by (type, subject, time, id)`).
     *CloudEvents v1.0 Specification*: `id UUID`, `type LowCardinality(String)`, `source LowCardinality(String)`, `subject LowCardinality(String) default 'GLOBAL'`, `time DateTime64(3, 'UTC') codec(Delta(8), ZSTD(1))`, `data String`, `severity Enum8('INFO' = 1, 'WARN' = 2, 'CRITICAL' = 3)`, `datacontenttype LowCardinality(String) default 'application/json'`.
     *Subject as Title Slug*: `subject` holds the unique title slug (e.g. `'eclipse'`), or entity slug for non-title events. Title-agnostic events explicitly record `subject = 'GLOBAL'`.
     *Lowercase SQL*: All DDL and SQL queries are authored in lowercase.
   - `002_qc.sql`: flattened QC inspection log (`fincher.qc`) projecting `id as event_id`, `subject as title_slug`, `JSONExtractString(data, ...)` with defensive `transform()` fallback to `UNKNOWN` and `NaN`, filtering `where type = 'fincher.qc.completed'`.
   - `003_vendor_metrics.sql`: `SummingMergeTree` rollup (`fincher.vendor_metrics`) with `measured_status_count` and `measured_drift_count` grouping daily inspections per `vendor_id` and `component`.
2. **Two ClickHouse access paths, two trust boundaries**:
   - `internal/clickhouse/client.go`: standard Go `database/sql` connection pool with `async_insert=1&wait_for_async_insert=1`, used only by deterministic backend code.
   - MCP Agent Interface: Google ADK Go's native `mcptoolset` connects Gemini agents in Phase 06 directly to `mcp-clickhouse:8000`. ClickHouse operates in `readonly=1` mode, natively rejecting any mutation attempt at the engine level.
3. **Event taxonomy** (`pkg/domain/models/event.go`) — ontology, not policy:
   - `TELEMETRY`:
     * `fincher.vendor.heartbeat` (Title-agnostic, `subject = 'GLOBAL'`)
     * `fincher.package.download.started`, `fincher.package.download.progress` (Title-scoped, `subject = <title_slug>`)
     — never leaves the ClickHouse write path.
   - `ROUTINE_OUTCOME`:
     * `fincher.qc.completed{status:PASSED}` (Title-scoped, `subject = <title_slug>`) — the lab's own verdict, not a Fincher-invented threshold.
   - `ANOMALY_SIGNAL`:
     * `fincher.qc.completed{status:FAILED/WARNING}` (Title-scoped, `subject = <title_slug>`)
     * `fincher.audio.sync_drift` (Title-scoped, `subject = <title_slug>`)
     * `fincher.master.cut.revised` (Title-scoped, `subject = <title_slug>`)
     * `fincher.vendor.sla_breach` (Dual-scoped: Title-scoped if tied to a title package deadline; Title-agnostic `subject = 'GLOBAL'` if facility-wide degradation)
     * `fincher.package.invalidated` (Title-scoped, `subject = <title_slug>`)
     — reaches a triage judge one event at a time in Phase 06 (no batching).
   - `ALLOCATION_REQUEST`:
     * `fincher.title.created`, `fincher.package.required`, `fincher.vendor.reconform.dispatched` (All Title-scoped, `subject = <title_slug>`)
     — handed off one at a time in Phase 06.
   - `OPERATOR_FORCED`:
     * `fincher.operator.forced`, `fincher.investigation.triggered` (Title-scoped if targeted at a title; Title-agnostic `subject = 'GLOBAL'` if system audit) — immediate in Phase 06.
4. **Ingestion API** (`internal/api/events/`): `POST /api/events` accepts a JSON array of
   CloudEvents, rejects empty arrays and non-array bodies with 400, validates each event, and
   inserts it directly into ClickHouse via `internal/clickhouse/events`. Returns
   `{status: "ingested", count: N}` on success. Classification (`event.Classify()`) exists and
   is tested, but nothing in this phase branches on its result — that's Phase 06's job, once the
   agent graph exists to actually consume the category.

## Explicit Non-Goals (deferred to Phase 06, or dropped entirely)
* No Gemini/ADK Go calls, no judge nodes, no `ActionPlan`, no vendor scoring, no SSE streaming.
* No debounce/coalesce/dedup — dropped, not deferred. See scope-cut note above.
* No budget/concurrency gate — deferred to Phase 06, must be Turso-persisted when built.
