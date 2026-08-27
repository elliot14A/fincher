# Phase 06 Plan: Dedicated Agents (ADK Go v2 Graph)

## Work Units

1. **Dependencies (`go.mod`)**:
   * Add `google.golang.org/adk/v2`, `google.golang.org/genai`.
   * Config: `FINCHER_GEMINI_API_KEY`, `FINCHER_POLICY_JUDGE_MAX_RETRIES` (default `3`).

2. **Ent Schemas (`internal/turso/ent/schema/`)**:
   * `investigation_run.go`, `run_node_event.go`, `judge_verdict.go` per `CONTEXT.md` §3.
   * Extend `vendor.go`: `standard_rate_usd`, `rush_rate_usd`, `standard_turnaround_hours`,
     `rush_turnaround_hours` (`field.Float32`, `.Optional()` + sane defaults for existing rows).
   * Extend `pkg/domain/models/vendor.go` + `UpdateVendorInput` to match.
   * `go generate ./internal/turso/ent`.

3. **Shared Types (`pkg/domain/models/`)**:
   * `actionplan.go`: `ActionPlan{TitleID, Summary, Actions []Action}`,
     `Action{Type, TargetID, Reason, NewVendorID}`.
   * `evidence.go`: `HistorianEvidence`, `LineageEvidence`, `VendorEvidence` (accuracy,
     rate/turnaround tiers, open incidents).
   * `verdict.go`: `Verdict{Outcome (APPROVED|REJECTED|ESCALATE|YES|NO), Rationale}`.

4. **ADK Go v2 Graph Wiring (`internal/agent/graph.go`)**:
   * `NewInvestigationGraph(deps) (agent.Agent, error)`: wires
     `triage_judge → {historian, lineage} → join → optimizer → policy_judge → executor`
     using `workflow.NewFunctionNode` (Go-only nodes), `workflow.NewAgentNode` (LLM judges),
     `workflow.NewJoinNode` (historian/lineage fan-in), `workflow.StringRoute` for
     YES/NO and APPROVED/REJECTED/ESCALATE branches, `workflow.Concat` for the reject→retry
     back-edge capped at `FINCHER_POLICY_JUDGE_MAX_RETRIES`.
   * `NewVendorAllocationGraph(deps) (agent.Agent, error)`: wires
     `vendor_scoring → vendor_judge → executor`.

5. **Incident Graph Nodes (`internal/agent/`)**:
   * `triage_judge.go`: flash call, input = `CoalescedBatch` + premiere countdown context,
     output = `Verdict{YES|NO}`.
   * `historian.go`: pre-baked queries (`internal/clickhouse`) keyed by `defect_category`; MCP
     tool-calling (`pkg/mcp`) fallback for unmapped categories.
   * `lineage.go`: Go-only, wraps `internal/turso/dependencies` graph traversal.
   * `optimizer.go`: flash/pro call, input = joined evidence (+ prior rejection reason on
     retry), output = `ActionPlan`.
   * `policy_judge.go`: flash call, input = `ActionPlan` + evidence, output =
     `Verdict{APPROVED|REJECTED|ESCALATE}`; persists a `judge_verdict` row per attempt.
   * `executor.go`: single SQLite tx applying `ActionPlan.Actions`, SSE broadcast per action
     (`run_node_event` insert + push), emits resulting event via `internal/clickhouse`
     (or an events-write helper) back into ClickHouse.

6. **Vendor Allocation Nodes (`internal/agent/`)**:
   * `vendor_scoring.go`: Go-only, queries `internal/clickhouse.RecencyWeightedAccuracy` +
     `OpenIncidents` per candidate, joins SQLite `Vendor` rate cards, returns
     `[]VendorEvidence`.
   * `vendor_judge.go`: flash call, input = `[]VendorEvidence` + title urgency context, output
     = `{VendorID, Reasoning}`; always invoked, no threshold gate.
   * Reuses `executor.go` for the `VENDOR_ASSIGNED` mutation + emission.

7. **Budget Guard (net-new — not stubbed in Feature 05)**:
   * Feature 05's `internal/ingest` design (bounded queue + SQLite counter) was cut before
     implementation — see `.agents/reviews/14-event-ingestion-pipeline/`. Build this fresh here,
     Turso-persisted (a `budget_counters` row, not an in-memory counter) so the daily cap survives
     Cloud Run scale-to-zero between requests.
   * Both graphs check the guard as their entry `workflow.NewFunctionNode` before any LLM node;
     on cap/concurrency exceeded, route to a deterministic-fallback terminal node that logs
     and exits at $0 cost — never a silent drop.

8. **API (`internal/api/runs/`, `internal/api/investigations/`)**:
   * `internal/api/runs/`: `get.go` (`GET /api/runs/{id}`), `stream.go`
     (`GET /api/runs/{id}/stream`, Echo SSE writer polling/subscribing to new
     `run_node_event`/`judge_verdict` rows).
   * `internal/api/investigations/`: `create.go` (`POST /api/investigations`,
     `POST /api/vendor-assignments`) — `OPERATOR_FORCED`, calls the graph directly.
   * `internal/api/budget/`: `get.go` (`GET /api/budget`).
   * There is no `Handoff` interface to wire from Feature 05 — `internal/api/events` in Feature 05
     only writes to ClickHouse, so this phase's graphs are invoked directly from wherever the
     per-event trigger lives (this endpoint's own classify-and-dispatch step, still to be added
     here, or a separate small entrypoint — not yet decided which).
   * Register all in `internal/api/server.go`.

9. **Frontend (`web/src/features/runs/`)**:
   * `src/lib/hooks/useSSEStream.ts`: shared SSE consumption hook.
   * `src/features/runs/graph/runGraph.tsx` + `.css.ts` + `index.ts`: `@xyflow/react` canvas
     reusing the Lineage DAG visual pattern, subscribing to `useSSEStream`.
   * `src/features/runs/queryKeys.ts`, `queryOptions.ts` per the Query Key Separation Rule.
   * Judge nodes render verdict + rationale on hover/click.

10. **Verification**:
    * Unit tests per node with a stub Gemini client (deterministic fixture responses) —
      `historian`, `lineage`, `vendor_scoring` fully testable without any live model call.
    * Integration test: fixture `CoalescedBatch` through the full incident graph, asserting the
      executor's SQLite mutations and emitted events.
    * Integration test: fixture allocation event through the vendor graph, asserting
      `VENDOR_ASSIGNED` + package vendor_id update.
    * Live demo run against real Gemini (manual, not CI) replaying Scenario A/B from the
      brainstorm (German dub disaster, midnight flood).
    * `go test ./... -race`; frontend `bun run typecheck`, `biome check src`, browser
      walkthrough showing nodes lighting up in real time.

## Update on Completion
* Flip Feature 06 checklist items in `STATE.md`, advance `Active Milestone` to Feature 07
  (Docent Conversational Assistant).
