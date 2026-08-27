# Fincher — Project Specification (Workflow DAG Edition)

## 1. Executive Summary & Core Value
**Fincher** is an autonomous delivery-integrity workflow engine for **LUME**, a premier streaming service. The console is organized around a **launch calendar** of upcoming releases with firm premiere dates, so every incident carries real operational stakes (*"Will Eclipse ship across 40 territories in time for Friday's global premiere?"*).

Fincher runs a structured **Studio Pipeline DAG** (Directed Acyclic Graph) of specialized nodes (data queries, Gemini Flash/Pro agent investigators, hybrid decision nodes, concrete state actions, and drafted stakeholder dispatches).

> **Core Operating Principle**:  
> **ClickHouse holds deep historical telemetry. Scheduled DAG nodes query OLAP views via MCP, AI synthesizes evidence with deterministic numbers inside the Decision Node, and software executes state mutations.**

---

## 2. The Four-Beat Story & The Disagreement Money Shot
Every run is legible in ten seconds:
> **Anomaly settles → Evidence assembled from ClickHouse/SQLite → Judge renders a verdict with rationale → State action executed & stakeholder dispatch drafted.**

### The Hero Scenario: Disagreement Panel
* **Context**: *Eclipse* master bumps to V13 two days before a 40-territory launch. Vendor A's Spanish audio produces a borderline sync drift (110ms).
* **Naive Rule Automation**: PASS (110ms sits under a hardcoded tolerance — the exact kind of invented threshold Fincher deliberately does not use).
* **Fincher Policy Judge (Gemini + ClickHouse MVs, no hardcoded threshold)**: HOLD. Reasons freshly over quantitative facts (36h to launch, 40 territories, redelivery count = 2, Vendor A's rolling defect rate escalating, upstream master superseded) and renders a verdict + rationale — the same 110ms reading on a non-primary territory with no premiere pressure would reasonably render RELEASE.

---

## 3. Architecture & Data Flow

Fincher is **event-driven and push-based**, not scheduler-polled. Every event lands in ClickHouse immediately; a lightweight classifier decides, per individual event, whether it's worth a model's attention — see `REQUIREMENTS.md` (`REQ-EVT`, `REQ-AGENT`) for the full node contract. Debounce/coalesce/batching was evaluated and deliberately dropped (2026-08-27) as unwarranted complexity at this scale — see `STATE.md` decision log; the daily-cap/concurrency budget gate is deferred to Feature 06, Turso-persisted, not built in-memory.

```text
Ingest → ClickHouse (events, immutable)
   │
   ▼
[Stage A: CloudEvent Filter] (pkg/domain/models, static, 0 LLM cost)
   TELEMETRY / ROUTINE OUTCOME → drop (write-only)
   ANOMALY SIGNAL              → straight to Stage C, one event at a time
   ALLOCATION REQUEST          → straight to evidence+judgment
   OPERATOR-FORCED             → immediate
   │
   ▼
[Stage C: Triage Judge] ⚡ single flash call per event, no hardcoded severity thresholds
   → route NO: logged, done (0 further cost)
   → route YES: continue
   │
   ▼
[Evidence fan-out]  Historian (hybrid: Go pre-baked queries + MCP tool-calling
                     for novel cases) ∥ Lineage (Go-only dependency walk)
   │
   ▼ join
[Optimizer ⚡] synthesizes ActionPlan from evidence
   │
   ▼
[Policy Judge ⚡] verdict: APPROVED / REJECTED (loop back to Optimizer, capped
                  retries, configurable up to 3) / ESCALATE (→ HOLD + alert)
   │
   ▼ APPROVED
[Executor] (Go-only) transactional SQLite mutation → SSE broadcast → emits
           resulting event back into ClickHouse (closed loop, AGENTS.md Inv. 4)
```

### Parallel path: vendor allocation (new title / package needing QC)
```text
TitleCreated / PackageRequired
   │
   ▼
[vendorScoringFn] (Go-only) merges SQLite rate cards (standard/rush rate +
                   turnaround) with ClickHouse recency-weighted accuracy
   │
   ▼
[Vendor Judge ⚡] ALWAYS fires (one flash call, no margin threshold gate) —
                  reasons about cost/speed/quality tradeoffs fresh given
                  premiere urgency and any open incidents on the candidates
   │
   ▼
[Executor] → VENDOR_ASSIGNED
```

Both paths share the same primitive: **Go assembles evidence → one scoped LLM judge renders a verdict → Go executes it transactionally.** No node in either graph invents a numeric policy threshold; the graph topology (not a rule table) is what stays deterministic.

---

## 4. Technical Stack & Invariants

| Component | Standard Technology | Architectural Invariant |
| :--- | :--- | :--- |
| **Language (Backend)** | Go (1.24+) | Standard library idioms, explicit error wrapping, `context.Context` propagation |
| **Runtime & Cloud** | Google Cloud Run (**cold / scale-to-zero**) | Stateless single-request DAG execution (~2–4s), $0 idle cost |
| **Application State Store** | Turso / libSQL (`@libsql/client` / Go driver) | Serverless HTTP SQLite storing titles, packages, runs, node executions, notifications |
| **Historical Analytics DB** | ClickHouse Cloud / Local Container (`24.3-alpine`) | 250k+ append-only QC/asset events + 4 Materialized Views |
| **Agent DB Interface** | Official ClickHouse MCP Server (`mcp-clickhouse`) | Remote MCP HTTP transport (`/mcp`). ClickHouse credentials isolated exclusively in MCP container |
| **AI Models** | Google GenAI SDK (`google.golang.org/genai`) + ADK Go v2 graph engine (`google.golang.org/adk/v2`, `workflow` package) | Gemini Flash for triage/vendor/policy judges & draft notifications; Gemini Pro reserved for the Optimizer's `ActionPlan` synthesis on complex incidents |
| **Decisions** | Evidence → Judgment → Execution primitive (`internal/agent/*`) | Go nodes assemble deterministic evidence; single-scoped LLM judge nodes render a verdict + rationale with no hardcoded thresholds; Go executes transactionally. Graph topology (ADK Go `workflow.Edge` routing), not a rule table, is what stays deterministic |
| **Operator Assistant** | Read-first Docent Assistant | Answers "What's releasing this weekend?" and explains past run decisions with SQL query citations |
| **HTTP API & SSE** | `github.com/labstack/echo/v4` | REST endpoints under `/api/*`, Swagger JSON spec serving, SSE for real-time node stepping |
| **Frontend UI Runtime** | **Preact + Vite + TypeScript** | Microscopic ~3kb UI footprint with `@preact/preset-vite` and `preact/compat` |
| **Frontend Styling** | **Vanilla Extract (`.css.ts`) + Recipes** | Zero-runtime CSS extraction, 100% type-safe design tokens (`theme.css.ts`) |
| **Frontend State & DB** | **`@tanstack/react-query` + `@tanstack/react-db`** | Reactive client-side database collections (`src/db/`) with live SSE sync and optimistic updates |
| **Frontend Routing** | **`@tanstack/react-router`** | Type-safe, file-based routing with `@tanstack/router-plugin` |
| **Frontend Data Grids** | **`@tanstack/react-table` + `@tanstack/react-virtual`** | 60fps virtualized territory matrices and package feeds |
| **Frontend DAG Canvas** | **`@xyflow/react`** | Interactive node canvas for Title Lineage DAG and Workflow Execution Inspector |
| **API Codegen & Validation** | **`@hey-api/openapi-ts` + `valibot`** | Auto-generates TypeScript SDK & Valibot validators from backend `openapi/swagger.json` |
| **Configuration** | `github.com/alecthomas/kong` | Strict `FINCHER_{SERVICE}_{SETTING}` environment variable bindings |

---

## 5. Non-Negotiable Rules
1. **ClickHouse via MCP Only (for agents)**: All *agent-issued* analytical history queries route through the official `mcp-clickhouse` server. Deterministic Go nodes (e.g. `vendorScoringFn`) may query ClickHouse directly — MCP is the AI-facing safety boundary, not the only access path.
2. **Google Cloud AI Tooling Only**: Exclusively use Google Gemini SDKs (via ADK Go v2).
3. **Zero-Cost Mechanical Filter**: The Stage A taxonomy filter runs in-memory Go, at $0 LLM cost, before any model is invoked. Only `ANOMALY_SIGNAL`/`ALLOCATION_REQUEST`/`OPERATOR_FORCED` events reach a judge, one event at a time — no batching/coalescing stage.
4. **No Invented Numeric Thresholds**: Fincher does not hardcode business-judgment constants (sync-drift ms tolerances, score margins, weighting coefficients). Deterministic code computes and presents *facts*; a scoped LLM judge renders the *verdict* over those facts, every time, with rationale.
5. **Mocked Notifications**: Stakeholder emails/dispatches are drafted with realistic LLM content and stored as `drafted` for in-browser review.
6. **Bounded Self-Correction, Not Unbounded Loops**: The Policy Judge's reject → revise cycle is capped (configurable, default up to 3) within a single investigation run. Beyond the cap, the run terminates in `ESCALATE` → `HOLD` + operator alert, never a silent retry storm.
7. **Strict camelCase & Co-located Frontend**: All frontend files and directories are `camelCase`. Every component and feature sub-component is co-located with its `*.css.ts` styling and `index.ts` barrel export.
