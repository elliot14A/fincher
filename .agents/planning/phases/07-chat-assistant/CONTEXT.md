# Phase 07 Context: Chat Assistant (Gemini Chat)

## Goal
Ship a genuinely smart, read-only, natural-language operator chat that can answer ANY
supply-chain question by freely combining two data sources: live SQLite operational state
(what IS) via a broad `query_turso` tool, and ClickHouse analytical history (what HAPPENED)
via a raw `query_clickhouse` MCP tool. It streams its answer over SSE and surfaces every tool
call it ran as an auditable citation.

The bar is open-ended question answering, not a fixed menu. It must handle:
* **Why did this happen?** — read `wf_results.rationale` + `steps` (the agent decision trail in
  Turso) joined with the triggering `fincher.events`/`fincher.qc` rows in ClickHouse.
* **When did this title go on HOLD?** — cross-reference `runs`/`steps` timestamps + the
  `fincher.delivery.held` / QC-failure events in ClickHouse.
* **What is X's status / why?** — live `titles`/`deliveries`/`media_packages` rows plus the
  history that produced them.
* **What premieres tomorrow / this week?** — `titles.premiere_date` range queries in Turso.
* **Which vendor is best for German audio? What's their track record?** — `vendors` rate cards
  (Turso) + `fincher.vendor_metrics` / `fincher.qc` history (ClickHouse).

## Objectives
* Add a Gemini tool-calling agent (`internal/agent/chat_agent.go`) with a read-only system
  prompt and a BROAD tool set: an extended `query_turso` (full operational schema INCLUDING
  runs/steps/wf_results — the "why"), a raw `query_clickhouse` MCP read tool (arbitrary
  historical SELECTs, not the narrow pre-baked helper), plus the existing projection/impact
  helpers for convenience. The agent decides which to call, and typically calls several across
  both stores to answer one question.
* Expose it over a small two-endpoint API (`internal/api/chat/`): submit a message, then
  stream the answer + citations via SSE — mirroring the proven `internal/api/runs/stream.go`
  pattern.
* Wire the existing static `web/src/routes/chat.tsx` shell to the live endpoint via a new
  `web/src/features/chat/` slice and a shared `useSSEStream` hook, rendering streaming text,
  multi-turn history, and citation pills.

## Architectural Decisions

### 1. Read-only, no mutations, no workflow triggers (Invariant 1)
* The chat assistant is analytical intelligence ONLY. It never issues writes, never dispatches
  actions, never enqueues events or runs. It answers and cites.
* Both query tools are read-only by construction: `query_turso` is SELECT/WITH-only with a
  `LIMIT` cap; `query_clickhouse` goes through the MCP `run_query` against a ClickHouse
  connection pinned to `readonly=1`. The read safety is enforced at the tool layer.

### 2. Tools — the two broad workhorses + helpers
The chat agent gets a PURPOSE-BUILT, broad tool set (do NOT just reuse `BuildAgentTools`, which
exposes narrow, planner-specific tools and omits the audit trail):
* **`query_turso` (extended schema)** — the existing `tools.NewTursoQueryTool`, but its schema
  hint MUST be extended to include the workflow/audit tables so the assistant can answer "why":
  - `runs(id, title_slug, trigger, status, started_at, ended_at)`
  - `steps(id, run_id, name, status, started_at, ended_at, metadata)`
  - `wf_results(id, run_id, step_id, judge, outcome, rationale, attempt)` ← `rationale` is the
    recorded judge reasoning = the literal "why did the agent decide this".
  - plus the existing `titles`, `media_packages`, `deliveries`, `vendors`, `masters`.
* **`query_clickhouse` (new raw read tool)** — a thin `functiontool` wrapping
  `mcp.RunQueryRows(ctx, client, sql)` with a schema hint for the analytical tables:
  - `fincher.events` (CloudEvents stream: type, subject=title slug, occurred_at, data JSON)
  - `fincher.qc` (QC inspections: vendor, component, status, drift, occurred_at)
  - `fincher.vendor_metrics` (rollups: vendor accuracy, turnaround, defect rates)
  This replaces the narrow, pre-baked `query_analytics` for chat — the assistant writes its own
  SELECTs so it can answer arbitrary historical questions, not just 3 fixed ones.
* **Convenience helpers (optional, pure reads)**: `tools.NewProjectionTool` (title readiness /
  premiere buffer) and `tools.NewDeliveryImpactTool` (blast radius) — handy shortcuts the agent
  may call for common readiness/impact questions.

### 3. Backend agent (`internal/agent/chat_agent.go`)
```
AnswerQuery(ctx, model, tursoClient, tursoDB, mcpClient, history, userMessage)
  → llmagent with read-only system prompt (prompts/chat.md) + [query_turso, query_clickhouse,
    projection, impact]
  → run the tool-calling loop; the agent commonly calls BOTH stores for one question
    (e.g. live status from Turso + the events/QC history from ClickHouse that explains it)
  → capture every tool call (name + SQL) as a Citation
  → return { answer string, citations []Citation }
```
* Citations are captured from the agent's tool-call transcript (tool name + the SQL text),
  typed as `{ label, source: "sqlite" | "clickhouse", tool }` to match the landing page's
  citation UI vocabulary.
* Single-phase is fine here (unlike the planner/selector two-phase fix): the chat agent's final
  output is free-form prose, not a strict array schema, so ADK's `OutputSchema` array-shrinking
  fragility does not apply.
* The system prompt must actively teach the "combine both stores" behavior: live state answers
  the *what/when/status*, the ClickHouse events/QC history + Turso `wf_results.rationale` answer
  the *why*. The agent should join them rather than answering from one source when a question
  spans both.

### 4. API & SSE (`internal/api/chat/`)
* `POST /api/chat` — body `{ message, session_id? }`. Creates/looks up an in-memory session,
  appends the user turn, kicks the agent asynchronously, returns `{ session_id }` immediately.
* `GET /api/chat/:session/stream` — SSE. Streams `event: token` frames (incremental answer
  text) then a terminal `event: done` frame carrying the final answer + `citations`. On error,
  `event: error`. Mirror the header/flush/`ctx.Done()` shape of `internal/api/runs/stream.go`
  (Content-Type `text/event-stream`, `X-Accel-Buffering: no`, `Flush()` per frame).
* Session store: in-memory `map[string]*session` guarded by a mutex, TTL-swept. Ephemeral is
  acceptable — a chat conversation does not need to survive Cloud Run scale-to-zero, unlike the
  budget gate. (Documented tradeoff, not an oversight.)
* Register via `internal/api/server.go`'s lazy `routeOnce`, guarded on model + MCP availability
  like the other agent-backed routes.

### 5. Frontend (`web/src/features/chat/` + `web/src/routes/chat.tsx`)
* Feature slice (strict camelCase + co-location, one `index.ts` barrel):
  - `queryKeys.ts`, `queryOptions.ts` (consume generated Hey API SDK for `POST /api/chat`).
  - `components/messageList/` (user + assistant bubbles, streaming answer), `components/
    citationPills/`, `components/composer/` — each `*.tsx` + `*.css.ts` + `index.ts`.
  - `hooks/useChatStream.ts` OR reuse a new shared `web/src/lib/hooks/useSSEStream.ts`
    consuming the chat SSE stream (token accumulation + done/error handling).
* `chat.tsx`: replace the no-op `handleSubmit` with a real submit → `POST /api/chat` →
  subscribe to the stream; render conversation history, streaming assistant text, and citation
  pills; keep the existing suggested-prompt buttons (they pre-fill the composer).
* Zero hardcoded CSS; all styling from `theme.css.ts` tokens. Reuse the citation-pill visual
  language already established on the landing page (`chatCitationPill` in `landing.css.ts`).

### 6. Prompt (`prompts/chat.md`)
* System prompt establishing: you are Fincher's read-only supply-chain chat assistant; use the
  tools to fetch real data before answering; NEVER fabricate numbers; cite the queries you ran;
  refuse/redirect any request to mutate state or trigger workflows (out of scope — read only).
* Teach the two-store mental model explicitly:
  - Turso (`query_turso`) = current truth: titles, deliveries, packages, vendors, masters, AND
    the agent audit trail (`runs`/`steps`/`wf_results.rationale`) that records WHY decisions were
    made.
  - ClickHouse (`query_clickhouse`) = history: the raw `fincher.events` stream, `fincher.qc`
    inspections, and `fincher.vendor_metrics` rollups — what actually happened over time.
  - For "why/when" questions, cross-reference both. For "status/what" questions, live Turso is
    usually enough. Prefer several small targeted queries over one giant join.
* Include worked example question→tool-plan pairs in the prompt (why-on-hold, when-held,
  premieres-tomorrow, vendor-track-record) so the model reliably reaches for the right store.

## Explicit Constraints
* **Strictly read-only** — no `INSERT`/`UPDATE`/`DELETE`, no event emission, no run dispatch.
  This is the headline invariant for the chat assistant (Invariant 1).
* **No fabricated data** — every quantitative claim must come from a real tool call, surfaced as
  a citation. If a tool returns nothing, say so; do not invent.
* **Reuse the plumbing, broaden the surface** — the SSE pattern (`internal/api/runs/stream.go`),
  the SELECT-guard/`RunTursoQuery`, and `pkg/mcp.RunQueryRows` already exist and are tested. This
  phase reuses that plumbing but deliberately gives the chat agent BROADER tools than the workflow
  agents (extended `query_turso` schema + a new raw `query_clickhouse`), because open-ended chat
  needs open-ended reads, not the planner's narrow pre-baked helpers.
* **Model config** — uses the single `cfg.GeminiModel` (`gemini-2.5-flash`) + `GeminiOptions`
  regional endpoint already wired in `cmd/fincher/main.go`. No new model config.

## Reference
* `PROJECT.md` §Invariant 1 (single source of mutation / read-only agents).
* `REQUIREMENTS.md`: `REQ-API-05` (read-only chat assistant), `REQ-UI-10` (chat assistant
  interface with SQL citations), `REQ-UI-01/02/03/04` (Preact/Vanilla-Extract/camelCase/codegen).
* Existing building blocks: `internal/agent/tools/turso.go` (`RunTursoQuery` + SELECT guard,
  extend its schema hint), `pkg/mcp/query.go` (`RunQueryRows` for the new `query_clickhouse`
  tool), `internal/agent/tools/projection.go` + `impact.go` (convenience helpers),
  `internal/api/runs/stream.go` (SSE pattern), `web/src/routes/chat.tsx` (static shell),
  `web/src/styles/routes/chat.css.ts`.
* Ent audit-trail tables (for the "why"): `runs`/`steps`/`wf_results` in
  `internal/turso/ent/schema/` — `wf_results.rationale` holds the recorded judge reasoning.
