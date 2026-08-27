# Phase 05 Plan: ClickHouse Ingestion Pipeline

## Work Units

1. **Config (`internal/config/`)**:
   * Kong-bound fields: `FINCHER_CLICKHOUSE_DSN` (default `127.0.0.1:9000`), `FINCHER_MCP_URL`
     (default `http://localhost:8000/mcp`), `FINCHER_MAX_CONCURRENT_AI_WORKFLOW_RUNS` (default `3`).
   * Domain constants: `DefaultMaxSyncDriftMs = 120.0`, `DefaultVendorDefectThreshold = 0.05`,
     `DefaultPremiereUrgentHours = 72.0`.

2. **ClickHouse Migrations (`migrations/clickhouse/`)**:
   * `001_events.sql`: CNCF CloudEvents v1.0 root stream (`id`, `type`, `source`, `subject`,
     `time`, `data`, `severity`, `datacontenttype`), `partition by toYYYYMM(time)`,
     `order by (type, subject, time, id)` with lowercase SQL keywords.
   * `002_qc.sql`: Materialized View `fincher.qc` projecting `id as event_id`, `subject as title_slug`,
     and metrics extracted from `data` where `type = 'fincher.qc.completed'`.
   * `003_vendor_metrics.sql`: SummingMergeTree `fincher.vendor_metrics` rolling up daily inspections
     per `vendor_id` and `component` where `type = 'fincher.qc.completed'`.
   * `internal/clickhouse/migrate.go`: applies `.sql` files in numeric filename order at startup.

3. **Direct ClickHouse Client (`internal/clickhouse/`)**:
   * `client.go`: establishes standard Go `database/sql` connection pool with
     `async_insert=1&wait_for_async_insert=1`.
   * `error.go`: maps ClickHouse driver errors into `domainerrors.Result[T]`.
   * `events/insert.go`: `Insert(ctx, db, event)` writes CloudEvents into `fincher.events`,
     defaulting empty subjects to `"GLOBAL"`.
   * `vendors/accuracy.go`: `RecencyWeightedAccuracy(ctx, db, vendorID, component)` computes
     120-day exponential decay pass rate (`REQ-CH-07`).
   * `vendors/metrics.go`: `GetMetrics(ctx, db, vendorID)` returns daily `MetricRow` slices.

4. **Agent-Facing MCP Integration (`google.golang.org/adk/tool/mcptoolset`)**:
   * Official ClickHouse MCP server (`mcp-clickhouse:8000`) is configured in Docker Compose.
   * Google ADK Go agents in Phase 06 connect to it natively via `mcptoolset.New(ctx, &mcptoolset.Config{Endpoint: cfg.MCPURL})`.
   * ClickHouse native `readonly = 1` mode strictly rejects any state mutation at the database engine level.

5. **CloudEvents Model & Taxonomy (`pkg/domain/models/event.go`)**:
   * CNCF CloudEvents v1.0 struct: `id`, `type`, `source`, `subject` (title slug), `time`, `data`,
     `severity`, `datacontenttype`.
   * Canonical event types: `fincher.qc.completed`, `fincher.audio.sync_drift`,
     `fincher.master.cut.revised`, `fincher.vendor.sla_breach`, `fincher.package.invalidated`,
     `fincher.vendor.heartbeat`, `fincher.title.created`, `fincher.package.required`, etc.
   * `Classify() EventCategory`: maps types to `TELEMETRY`, `ROUTINE_OUTCOME`, `ANOMALY_SIGNAL`,
     `ALLOCATION_REQUEST`, `OPERATOR_FORCED`.
   * Table-driven unit tests covering all categories and QC PASSED/FAILED branches.

6. **Dropped (2026-08-27): Debounce & Coalescing Engine, Budget & Concurrency Guard**.
   Designed in full in `.agents/reviews/14-event-ingestion-pipeline/20260827T171710Z.md`
   (`internal/ingest/`: `debounce.go`, `coalesce.go`, `budget.go`, a `budget_counters` SQLite
   table) then cut as unwarranted complexity for this scale — see `STATE.md` decision log.
   Debounce/coalesce/dedup will not return. The budget/concurrency gate is deferred to Phase 06,
   and must be Turso-persisted (not in-memory) when it's eventually built, to survive Cloud Run
   scale-to-zero between requests.

7. **Ingestion API (`internal/api/events/`)**:
   * `routes.go`, `create.go`: `POST /api/events` — binds a JSON array of CloudEvents, 400s on
     empty array or non-array body, validates and inserts each event sequentially via
     `internal/clickhouse/events.Insert`, returns `{status: "ingested", count: N}` on success.
     No classification-based branching, no debounce, no `Handoff` interface — Feature 06 owns
     everything downstream of the ClickHouse write.
   * Registered in `internal/api/server.go`, gated on a non-nil ClickHouse connection (server
     still boots and serves everything else if ClickHouse is unreachable at startup).
   * HTTP lifecycle tests: successful batch insert, empty array, single non-array object,
     malformed JSON.

8. **Typed Event Payloads (`pkg/domain/models/event_payloads.go`)**:
   * `QCPayload`, `MasterCutPayload`, `SLABreachPayload`, `PackageInvalidatedPayload` structs.
   * `UnmarshalPayload[T any](e *Event) (T, error)` generic helper — no `map[string]any`
     extractor helpers.

9. **Verification**:
   * `go test ./... -race` green.
   * End-to-end event ingestion test against live ClickHouse container, including from a
     genuinely fresh instance with no `fincher` database yet.
