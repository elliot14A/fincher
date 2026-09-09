# Phase 07 Plan: Chat Assistant (Gemini Chat)

> **Core Principle**: Read-only analytical chat. Reuse the existing agent tool layer
> (`internal/agent/tools/` + `pkg/mcp`) and the proven SSE pattern (`internal/api/runs/stream.go`).
> Every quantitative answer is backed by a real tool call surfaced as an auditable citation.
> Strictly no mutations, no event emission, no workflow dispatch (Invariant 1).

---

## Work Units

### Unit 1: Broaden the read tools (extend `query_turso`, add `query_clickhouse`)
* **Target Files**:
  * `internal/agent/tools/turso.go` (extend the schema hint in `tursoQueryToolDescription`)
  * `internal/agent/tools/clickhouse.go` (new — raw `query_clickhouse` MCP read tool)
* **Details**:
  1. Extend `tursoQueryToolDescription` to include the workflow/audit tables so the assistant can
     answer "why" and "when":
     * `runs(id, title_slug, trigger, status, started_at, ended_at)`
     * `steps(id, run_id, name, status, started_at, ended_at, metadata)`
     * `wf_results(id, run_id, step_id, judge, outcome, rationale, attempt)` — note in the hint
       that `rationale` is the recorded judge reasoning (the "why").
     * Keep the existing `titles`, `media_packages`, `deliveries`, `vendors`, `masters` +
       `json_each` guidance. This is additive and safe: the SELECT/`LIMIT` guard in
       `RunTursoQuery` is unchanged, so the workflow agents that already use this tool are
       unaffected (they simply never query the new tables).
  2. Add `tools.NewClickhouseQueryTool(client *mcp.Client) (tool.Tool, error)` in a new
     `internal/agent/tools/clickhouse.go`:
     * `functiontool` named `query_clickhouse`, wrapping `mcp.RunQueryRows(ctx, client, args.SQL)`.
     * Description = a schema hint for the analytical tables (`fincher.events`, `fincher.qc`,
       `fincher.vendor_metrics`) with 2-3 example SELECTs (e.g. events for a title ordered by
       `occurred_at`, QC failures for a vendor, vendor rollups). Read-only is guaranteed by the
       MCP connection's `readonly=1`; the tool passes the SQL straight through.
     * Return `[]map[string]any` rows (same shape as `query_turso`) so the model handles both
       uniformly.
* **Verification**:
  * `go build ./internal/agent/tools/...` clean; existing workflow-agent tests still green
    (schema hint change is additive, no behavior change to `RunTursoQuery`).
  * Live smoke: `query_clickhouse` returns rows for a simple `SELECT count() FROM fincher.events`.

---

### Unit 2: Chat System Prompt
* **Target Files**:
  * `prompts/chat.md` (new)
  * `prompts/prompts.go` (add `//go:embed chat.md` + `var Chat string`)
* **Details**:
  1. Author a read-only system prompt establishing the assistant persona: Fincher's supply-chain
     operations chat for global film/TV delivery.
  2. Rules to encode:
     * Always call the tools to fetch real data before answering; NEVER fabricate numbers,
       vendor names, statuses, or counts.
     * If a tool returns no rows, state that plainly — do not invent.
     * Cite the queries used (the runtime attaches them; the prompt instructs the model to lean
       on real tool output).
     * Refuse/redirect any request to mutate state, hold/release deliveries, reassign vendors, or
       trigger workflows — those are out of scope; the assistant only explains and cites.
  3. Teach the two-store mental model explicitly:
     * `query_turso` = current truth (titles, deliveries, packages, vendors, masters) AND the
       agent audit trail (`runs`/`steps`/`wf_results.rationale` = WHY the agent decided things).
     * `query_clickhouse` = history (`fincher.events`, `fincher.qc`, `fincher.vendor_metrics` =
       what happened over time).
     * For "why/when" → cross-reference both stores. For "status/what/what premieres" → live Turso
       is usually enough. Prefer several small targeted queries over one giant join.
  4. Include worked example question→tool-plan pairs in the prompt so the model reliably reaches
     for the right store:
     * "why is X on HOLD" → `query_turso` (title/delivery/package status + latest run's
       `wf_results.rationale`) + `query_clickhouse` (the QC failure / `fincher.delivery.held` event).
     * "when did X go on HOLD" → `query_turso` run/step timestamps + `query_clickhouse` events by
       `occurred_at`.
     * "what premieres tomorrow" → `query_turso` on `titles.premiere_date` range.
     * "best vendor for German audio + track record" → `query_turso` vendors (json_each) +
       `query_clickhouse` `fincher.vendor_metrics`/`fincher.qc`.
* **Verification**:
  * File exists; content is read-only-scoped with explicit "no mutation / no fabrication" rules
    and the two-store guidance + worked examples.
  * `prompts.Chat` compiles and is referenced by `chat_agent.go` in Unit 3.

---

### Unit 3: Chat Agent (`internal/agent/chat_agent.go`)
* **Target Files**:
  * `internal/agent/chat_agent.go` (new)
  * `internal/agent/chat_agent_test.go` (new, optional — live E2E may substitute per AGENTS.md)
* **Details**:
  1. Define the citation + result types:
     ```go
     type ChatCitation struct {
         Label  string `json:"label"`  // the SQL text the tool ran
         Source string `json:"source"` // "sqlite" | "clickhouse"
         Tool   string `json:"tool"`   // "query_turso" | "query_clickhouse" | "projection" | "impact"
     }
     type ChatAnswer struct {
         Answer    string         `json:"answer"`
         Citations []ChatCitation `json:"citations"`
     }
     ```
  2. Implement `AnswerQuery(ctx, m model.LLM, tursoClient *ent.Client, tursoDB *sql.DB,
     mcpClient *mcp.Client, history []Turn, userMessage string) domainerrors.Result[*ChatAnswer]`:
     * Assemble the BROAD chat tool set explicitly (do NOT call `BuildAgentTools`, which is the
       narrow planner set): `tools.NewTursoQueryTool(tursoDB)` (extended schema),
       `tools.NewClickhouseQueryTool(mcpClient)`, plus `tools.NewProjectionTool(tursoClient)` and
       `tools.NewDeliveryImpactTool(tursoClient)` as convenience helpers. All read-only.
     * Construct an `llmagent` with `prompts.Chat` as the system instruction, the tools attached,
       and the conversation `history` rendered into the prompt (or ADK session messages).
     * Run the tool-calling loop; collect each tool invocation (tool name + SQL) into
       `[]ChatCitation`, mapping `query_clickhouse` → `clickhouse`, `query_turso`/projection/impact
       → `sqlite`.
     * Return `{ answer, citations }`. Free-form prose output — NO `OutputSchema` (avoids the ADK
       array-shrinking issue that only affects strict structured outputs; N/A for prose).
  3. Keep it single-phase (tool-agent that also answers). The two-phase `runToolThenSchema` helper
     is NOT needed here because the final output is unstructured text.
  4. Reuse `MapError`/`NewError` and the `domainerrors.Result` conventions already used across
     `internal/agent/`.
* **Verification**:
  * `go build ./internal/agent/...` clean.
  * Live: `AnswerQuery` returns a non-empty answer and ≥1 citation for a data question, and for a
    "why" question calls BOTH stores (a `query_turso` audit-trail read + a `query_clickhouse`
    history read) — verified via the endpoint in Unit 5 / E2E in Unit 7.

---

### Unit 4: Chat Session Store (`internal/api/chat/session.go`)
* **Target Files**:
  * `internal/api/chat/session.go` (new)
* **Details**:
  1. In-memory session store: `type store struct { mu sync.Mutex; sessions map[string]*session }`.
  2. `type session struct { ID string; History []agent.Turn; Answer string; Citations []agent.ChatCitation; Status string /* running|done|error */; Err string; UpdatedAt time.Time }`.
  3. Methods: `New()`, `Get(id)`, `Append(id, turn)`, `SetResult(id, answer, citations)`,
     `SetError(id, msg)`, and a background TTL sweep (drop sessions idle > 30m).
  4. Ephemeral by design (documented tradeoff): a chat conversation need not survive Cloud Run
     scale-to-zero, unlike the budget gate. No Turso persistence.
* **Verification**:
  * Unit test: create → append → set result → get returns the right state; concurrent access is
    race-clean (`go test -race ./internal/api/chat/...`).

---

### Unit 5: Chat API + SSE (`internal/api/chat/`)
* **Target Files**:
  * `internal/api/chat/create.go` (new) — `POST /api/chat`
  * `internal/api/chat/stream.go` (new) — `GET /api/chat/:session/stream`
  * `internal/api/chat/routes.go` (new)
  * `internal/api/server.go` (wire lazy route registration; add deps if needed)
  * `openapi/generate.go` regen → `openapi/swagger.json`
* **Details**:
  1. `POST /api/chat` — body `{ message string, session_id?: string }`:
     * Guard: if model or MCP is unavailable, return 503 with a clear message (mirror the
       existing "AI runtime not initialized" guard pattern in `events/create.go`).
     * Look up or create the session; append the user turn; set status `running`.
     * Launch `agent.AnswerQuery(...)` in a background goroutine with a bounded `context.WithTimeout`
       (~60s); on completion call `store.SetResult` / `store.SetError`.
     * Return `{ session_id }` immediately (202/200).
  2. `GET /api/chat/:session/stream` — SSE, mirroring `internal/api/runs/stream.go`:
     * Set `text/event-stream`, `no-cache`, `keep-alive`, `X-Accel-Buffering: no`; `Flush()` per frame.
     * Poll the session on a ticker (~200-300ms). Emit `event: token` frames as the answer text
       grows (or a single `event: message` when complete if incremental tokens aren't captured),
       then a terminal `event: done` with `{ answer, citations }`. On error, `event: error`.
     * Terminate on `status == done|error` or `ctx.Done()`; bail after N consecutive missing-session
       reads (same guard shape as runs stream).
  3. `routes.go`: `RegisterRoutes(g, tursoClient, tursoDB, mcpClient, modelProvider)`; register
     `POST ""` and `GET "/:session/stream"`.
  4. Wire into `server.go`'s `registerRoutes` (lazy `sync.Once`), guarded on model+MCP+tursoDB set.
  5. Add swag annotations to both handlers; regenerate `openapi/swagger.json` (`go generate ./openapi`).
* **Verification**:
  * `go build ./...` + `go vet` + `gofmt` clean.
  * HTTP lifecycle: `POST /api/chat` returns a `session_id`; opening the stream yields a `done`
    frame with a non-empty answer + citations for a real data question (live, against ClickHouse+Turso).
  * `openapi/swagger.json` contains the two new paths (`spec_test.go` still passes).

---

### Unit 6: Frontend Chat Feature Slice + Route Wiring
* **Target Files**:
  * `web/src/features/chat/queryKeys.ts` (new)
  * `web/src/features/chat/queryOptions.ts` (new — consumes regenerated Hey API SDK)
  * `web/src/features/chat/components/messageList/` (`messageList.tsx` + `.css.ts` + `index.ts`)
  * `web/src/features/chat/components/citationPills/` (`citationPills.tsx` + `.css.ts` + `index.ts`)
  * `web/src/features/chat/components/composer/` (`composer.tsx` + `.css.ts` + `index.ts`)
  * `web/src/features/chat/hooks/useChatStream.ts` (+ `index.ts`) OR
    `web/src/lib/hooks/useSSEStream.ts` (shared)
  * `web/src/features/chat/index.ts` (barrel)
  * `web/src/routes/chat.tsx` (wire to live endpoint)
  * `web/src/styles/routes/chat.css.ts` (extend as needed)
  * `web/openapi-ts.config.ts` regen (`bun run generate:api`) after Unit 5
* **Details**:
  1. Regenerate the Hey API SDK from the updated `openapi/swagger.json` (`bun run generate:api`) so
     `postChat` (and any types) exist.
  2. `useChatStream` (or shared `useSSEStream`): given a `session_id`, open an `EventSource` to
     `/api/chat/:session/stream`, accumulate `token`/`message` frames into streaming answer text,
     resolve on `done` (capturing citations), surface `error`. Clean up on unmount.
  3. `chat.tsx`: replace the no-op `handleSubmit` — on submit, `POST /api/chat` with the message +
     current `session_id`, then subscribe via the hook. Render:
     * Conversation history (user + assistant bubbles), streaming assistant text.
     * `citationPills` under each assistant answer (Database icon → clickhouse, Terminal icon →
       sqlite), reusing the visual language from the landing page (`chatCitationPill`).
     * Keep the suggested-prompt buttons pre-filling the composer.
     * Loading/streaming state, empty state, and error toast (`sonner`).
  4. Strict rules: camelCase files, co-location with `*.css.ts` + `index.ts` barrels, zero hardcoded
     CSS (tokens only), zero inline `style={{}}`.
* **Verification**:
  * `bun run typecheck` + `biome check src` clean.
  * `bun run build` succeeds.
  * Browser: submit a question → streaming answer renders → citation pills show → multi-turn
    follow-up retains context.

---

### Unit 7: Live Verification & Citation Accuracy
* **Target Files**:
  * `tmp/eval/chat.py` (gitignored harness, optional) or manual browser + curl checks
* **Details**:
  1. Fire representative queries against the running server (real Gemini + MCP + ClickHouse + Turso):
     * "Why is <title> on HOLD and what's the blast radius?"
     * "Which vendor has the best historical accuracy for German audio?"
     * "Are the German dubs ready for <title>'s premiere?"
     * "What happened when <title>'s master bumped to V02?"
  2. Assert: answers are non-empty, grounded, and each surfaced citation corresponds to a tool call
     that actually ran (cross-check against server behavior / returned citation SQL).
  3. Assert read-only: a mutation-style request ("put X on hold") is politely refused/redirected,
     and NO state changes (verify title/delivery/package rows unchanged after the request).
* **Verification**:
  * All representative queries return grounded, cited answers.
  * Mutation-style prompt produces a refusal and zero DB writes (confirmed via API reads before/after).

---

## Cross-Cutting Constraints (apply to every unit)
* **Read-only (Invariant 1)** — no writes, no event emission, no run/workflow dispatch anywhere in
  the chat path. The tool layer already enforces SELECT/WITH-only + ClickHouse `readonly=1`.
* **No fabrication** — every number/name/status in an answer traces to a real tool call citation.
* **Broaden the surface, reuse the plumbing** — reuse `RunTursoQuery`'s SELECT guard,
  `pkg/mcp.RunQueryRows`, the projection/impact helpers, and the `runs/stream.go` SSE shape; but
  give the chat agent BROADER tools than the workflow agents (extended `query_turso` schema + new
  raw `query_clickhouse`). Do NOT reuse the narrow planner-only `BuildAgentTools`/`query_analytics`.
* **Single model config** — `cfg.GeminiModel` + `cfg.GeminiOptions` (regional Vertex) already wired;
  no new model configuration.
* **Frontend invariants** — camelCase, deep co-location + `index.ts` barrels, zero-runtime Vanilla
  Extract with theme tokens only, Hey API codegen as the type source of truth.

## Definition of Done
* Operator can hold a multi-turn conversation at `/chat` that returns grounded, cited answers
  streamed live from Gemini over real ClickHouse + Turso data, and the assistant refuses any mutation
  request. Backend `go build`/`vet`/`gofmt`/tests green; frontend `typecheck`/`biome`/`build` green;
  `openapi/swagger.json` regenerated with the two new chat paths.
