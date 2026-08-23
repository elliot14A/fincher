# AGENTS.md — Fincher Engineering Context & Guidelines

This document defines the architectural rules, engineering standards, multi-agent invariants, and pair-programming contract for all AI assistants working on **Fincher**.

---

## 1. Engineering Ownership & AI Collaboration Contract

### Core Working Philosophy
> **Human establishes understanding, intent, architecture, invariants, and implementation.**  
> **AI assists with research, scaffolding, boilerplate, and structure ONLY when explicitly requested.**  
> **Human reviews, writes the core logic, and maintains 100% intellectual ownership of the system.**

### Guiding Principles

1. **Permission-Gated Implementation**:
   * The AI assistant must **never** write core domain logic, algorithms, or complex features without explicit user instruction.
   * The AI assistant's primary role is scaffolding, boilerplate definitions, type declarations, configuration templates, and research.
2. **Untrusted Until Verified**:
   * All AI-generated code is untrusted until reviewed, compiled, and verified against defined rules, unit tests, and integration checks.
3. **Event-Driven Deterministism**:
   * All actions emerge dynamically from incoming production events, historical evidence from ClickHouse MCP, multi-agent reasoning, and database-backed policies.
4. **No Unchecked AI Authority**:
   * Agents **never** mutate production state directly. Agents produce structured proposals (`ActionPlan`). The Go policy engine deterministically gates every action against the database policies table.

### AI Engineering Rules
1. **Think Before Proposing**:
   * State all assumptions explicitly before proposing code or boilerplate.
   * If requirements, event types, or concurrency lifecycles are ambiguous, surface the uncertainty explicitly.
2. **Simplicity First & Idiomatic Go / Modern Frontend**:
   * Implement the cleanest working solution. Strictly NO speculative abstractions, unnecessary generic wrappers, or unrequested framework overhead.
   * Keep package boundaries strict and dependencies minimal.
3. **Surgical Changes**:
   * Touch only what is necessary. Every modified line must trace directly back to the active request.
4. **Goal-Driven Verification**:
   * Frame every implementation with concrete, verifiable success criteria (unit tests, schema verification, MCP tool calls, policy evaluations, type checks, lint checks).

---

## 2. GSD Core Spec-Driven Development (`.agents/planning/`)

All engineering milestones, state transitions, and requirements are managed through the **GSD Core** spec-driven framework in `.agents/planning/`:

1. **`PROJECT.md`**: Foundational domain context, architectural boundaries, and core operating principles.
2. **`REQUIREMENTS.md`**: Deterministic requirement IDs (`REQ-INGEST-01`, `REQ-MCP-01`, `REQ-POLICY-01`, `REQ-UI-*`, etc.) with verifiable acceptance criteria.
3. **`ROADMAP.md`**: Multi-phase progression, deliverables, dependencies, and Definition of Done.
4. **`STATE.md`**: Live operational state, active phase pointer, blocker log, and key decision log.
5. **Phase Loop**: Every phase executes through the 5-step loop: `DISCUSS` -> `PLAN` -> `EXECUTE` -> `VERIFY` -> `SHIP`.

---

## 3. System Invariants & Autonomous Workflow Contract

### Invariant 1: Single Source of Mutation
* **AI Agents are strictly read-only.** Agents query ClickHouse via MCP and reason with Gemini. Agents **NEVER** issue SQL `INSERT`/`UPDATE`/`DELETE` or execute system mutations directly.
* All state mutations are performed exclusively by the Go application (`internal/executor`) upon positive policy verification.

### Invariant 2: Parallel Investigation & Independence
* The **Historian Agent** and **Dependency Agent** execute concurrently via Go goroutines / Google ADK Go orchestrator.
* Each agent operates with specific bounded system instructions and focused tool access.

### Invariant 3: Deterministic Policy Gating
* The Action Planner proposes *intent* (e.g., `HOLD_DELIVERY`, `INVALIDATE_PACKAGE`).
* The Policy Engine verifies concrete conditions against operational rules stored in the database `policies` table.
* AI never overrides policy thresholds.

### Invariant 4: Closed-Loop Event Progression
* Every executed action must emit a downstream event (e.g., `INVALIDATE_PACKAGE` -> emits `PACKAGE_INVALIDATED` -> creates re-QC job -> emits `QC_STARTED` -> `QC_PASSED` -> emits `DELIVERY_RELEASED`).
* Fincher listens to its own resulting events to drive the workflow to resolution (`READY_TO_SHIP` or `HOLD`) without manual human kicking.

### Invariant 5: Complete Auditability & Contract-First Sync
* Every step is written to SQLite and streamed in real-time to the UI via SSE.
* Backend OpenAPI (`openapi/swagger.json`) serves as the single source of truth for frontend types, Valibot schemas, and TanStack Query options via Hey API (`@hey-api/openapi-ts`).

---

## 4. Technology Stack & Ecosystem Standards

| Layer / Purpose | Technology / Package | Notes & Invariants |
| :--- | :--- | :--- |
| **Language (Backend)** | Go (1.24+) | Standard library idiomatic code, explicit error handling |
| **HTTP Framework** | `github.com/labstack/echo/v4` | REST endpoints for `/api/*`, OpenAPI spec serving at `/openapi.json`, SSE for live console |
| **AI Runtime** | Google ADK Go (`google.golang.org/adk/v2`) + Google GenAI (`google.golang.org/genai`) | Programmatic multi-agent orchestration, concurrency, structured schema output |
| **LLM Model** | Gemini 2.5 / Gemini 2.0 Flash / Pro | Structured JSON outputs, fast reasoning over historical analytical evidence |
| **Analytical DB** | ClickHouse | Historical event store: QC logs, asset updates, vendor track records, past incidents |
| **Agent Interface to DB** | Official ClickHouse MCP Server (`ghcr.io/clickhouse/mcp-clickhouse:latest`) | Remote MCP HTTP transport client (`pkg/mcp`). ClickHouse credentials isolated exclusively in MCP container. |
| **Application State & Policies DB** | Turso / SQLite + `go-sqlite3` | High-performance operational state & policies table with WAL mode |
| **Frontend Runtime** | **Preact + Vite + TypeScript** | Microscopic ~3kb UI runtime with `@preact/preset-vite` and `preact/compat` |
| **Frontend Styling** | **Vanilla Extract (`.css.ts`) + Recipes** | Zero-runtime CSS extraction, 100% type-safe design tokens (`theme.css.ts`) |
| **Frontend Routing** | **`@tanstack/react-router`** | Type-safe, file-based routing (`src/routes/`) with `@tanstack/router-plugin` |
| **Frontend Data & State** | **`@tanstack/react-query` + `@tanstack/react-db`** | Reactive client-side database collections (`src/db/`) with live SSE sync and optimistic updates |
| **Frontend Data Grids & Canvas**| **`@tanstack/react-table` + `@tanstack/react-virtual` + `@xyflow/react`** | 60fps virtualized territory matrices + interactive node Lineage DAG |
| **Frontend Icons & UI Feedback**| **`lucide-preact` + `sonner`** | Preact-native icons with zero React wrapper surface + toast alerts |
| **Frontend API Codegen** | **`@hey-api/openapi-ts` + `valibot`** | Auto-generates TypeScript SDK & Valibot validators in `src/lib/api/generated/` from backend `openapi/swagger.json` |
| **Frontend Tooling** | **Biome (`biome.json`)** | Sub-millisecond formatting and strict linting |
| **Container & Infra** | Docker Compose + Google Cloud Run + Nix Flake | Production multi-stage build embedding `web/dist` via Go `embed.FS` (Nix devShell: `bun`) |

---

## 5. Frontend (`web/`) Architecture & Strict Invariants

### 1. Strict `camelCase` File & Directory Naming Rule
* **All files and folders across `web/src/` must be `camelCase`** (e.g. `calendarGrid.tsx`, `calendarGrid.css.ts`, `queryKeys.ts`, `queryOptions.ts`, `holdOverrideModal.tsx`).
* **Strictly NO kebab-case** (`calendar-grid.tsx` ❌) and **NO PascalCase files** (`CalendarGrid.tsx` ❌).
* **Exceptions**:
  - TanStack Router route parameter files dictated by framework conventions (`src/routes/$id.tsx`, `__root.tsx`, `index.tsx`).
  - Untouched machine-generated Hey API output in `src/lib/api/generated/` (`sdk.gen.ts`, `types.gen.ts`, `valibot.gen.ts`).

### 2. Deep Component Co-Location Rule
* Components **never** sit dumped in a flat directory.
* Every UI primitive or feature sub-component lives in its own dedicated directory alongside its styling and index barrel:
  * `src/components/ui/button/button.tsx` + `button.css.ts` + `index.ts`
  * `src/features/calendar/grid/calendarGrid.tsx` + `calendarGrid.css.ts` + `index.ts`
* Feature-local hooks live inside that feature's `hooks/` directory (e.g. `src/features/lineage/hooks/useLineageLayout.ts` + `index.ts`).
* Shared cross-feature hooks live in `src/lib/hooks/` + `index.ts`.
* Every subdirectory contains an `index.ts` barrel acting as the sole public export surface.

### 3. Query Key & Query Option Separation Rule
* Every feature maintains separated files:
  * `queryKeys.ts`: Deterministic query key factory.
  * `queryOptions.ts`: TanStack Query options consuming the auto-generated Hey API client from `src/lib/api/generated/`.

### 4. Zero-Runtime Styling & Zero Hardcoded CSS Rule
* All styling is authored in `*.css.ts` using `@vanilla-extract/css` and `@vanilla-extract/recipes` strictly consuming tokens from `src/styles/theme.css.ts`.
* **Strictly NO hardcoded CSS**: Never write raw hex values (`#ffffff` ❌), arbitrary `rgba()` strings, or hardcoded sizes directly in `.css.ts` or component JSX.
* **Strictly NO inline `style={{ ... }}` attributes in JSX**. All presentation logic belongs in co-located `*.css.ts` files or recipe variants.

### 5. Universal Pagination Standard (Gaur Standard)
* All paginated data grids, catalog lists, and table views must use the unified pagination contract and `<PaginationControls />` primitive from `src/components/ui/pagination/`.
* Standard envelope: `items`, `total_items`, `page`, `limit`, `total_pages`, `has_next_page`, `has_prev_page` (matching `.agents/rules/pagination.md`).

### 6. Runtime & Dev Integration Invariants
* **Preact/compat aliasing is mandatory** in `vite.config.ts` for all `@tanstack/react-*`, `@xyflow/react`, and `sonner` packages.
* **`lucide-preact`** is used directly (eliminating React SVG wrapper overhead).
* **Dev API Proxying**: `vite.config.ts` proxies `/api` to the Go backend (`http://localhost:8080`).

---

## 6. Architecture & Directory Boundaries

```text
fincher/
├── cmd/
│   ├── fincher/                  # Main unified server binary (API, Orchestrator, UI, Simulator)
│   └── seed/                     # Historical dataset generator for ClickHouse & SQLite
│
├── internal/                     # Private Go domain logic (strict boundaries)
│   ├── api/                      # REST handlers (/api/*), OpenAPI serving, SSE stream
│   ├── agent/                    # Multi-agent orchestrator & Google ADK Go sub-agents
│   ├── policy/                   # Deterministic Policy Engine (evaluates DB policies table)
│   ├── executor/                 # Software Execution Engine (state mutations & downstream events)
│   ├── simulator/                # Production media event generator
│   ├── turso/                    # Database access layer (Turso/SQLite state & policies)
│   └── config/                   # Configuration parsing via Kong
│
├── pkg/                          # Shared Go contracts, wire types, MCP integration
│   ├── domain/                   # Core types: Event, Delivery, Package, Incident, Action, Evidence
│   ├── events/                   # Event vocabulary & constants
│   ├── mcp/                      # ClickHouse MCP client wrapper for ADK Go agents
│   └── logger/                   # Structured JSON logging
│
├── openapi/                      # Canonical backend OpenAPI contract
│   ├── generate.go               # Swag generator directive (pinned v1.16.4)
│   ├── swagger.json              # Auto-generated OpenAPI JSON specification
│   ├── spec.go                   # Go embed.FS wrapper
│   └── spec_test.go              # Specification content verification test
│
├── web/                          # Preact Operations Console UI
│   ├── package.json              # Minimal type-safe dependencies (Bun)
│   ├── vite.config.ts            # Preact + TanStack Router + Vanilla Extract + /api dev proxy
│   ├── biome.json                # Biome formatting & linting configuration
│   ├── openapi-ts.config.ts      # Hey API SDK & Valibot schema generator config
│   ├── public/                   # Static assets & favicon
│   └── src/
│       ├── main.tsx              # Root entry & providers (Query, DB)
│       ├── app.css.ts            # Global styling reset & font bindings
│       ├── vite-env.d.ts         # Vite client typing & ImportMetaEnv
│       ├── styles/               # Design tokens (theme.css.ts, tokens.ts)
│       ├── db/                   # TanStack DB reactive client collections
│       ├── routes/               # TanStack file-based routes (__root.tsx, index.tsx, etc.)
│       ├── features/             # Feature slices (calendar, lineage, deliveries, vendors, runs, docent, layout)
│       ├── components/           # Co-located UI primitives (ui/button, ui/modal, feedback/skeletonLoader)
│       └── lib/                  # Shared singletons, API client & cross-feature hooks
│           ├── queryClient.ts    # TanStack QueryClient instance
│           ├── dbClient.ts       # TanStack DB instance
│           ├── api/              # Configured Hey API fetch client & generated SDK
│           │   ├── client.ts     # Configured /api client
│           │   └── generated/    # sdk.gen.ts, types.gen.ts, valibot.gen.ts
│           └── hooks/            # Cross-feature hooks (useSSEStream, useDebounce)
│
├── data/
│   ├── clickhouse/               # ClickHouse DDL schemas & materialized views
│   └── seed/                     # Synthetic media historical event seeds
│
├── .agents/                      # AGY project customizations & workflow runbooks
│   ├── planning/                 # GSD Core durable spec & state management (PROJECT, REQUIREMENTS, ROADMAP, STATE)
│   ├── reviews/                  # Verified review reports & architecture plans
│   ├── rules/                    # Rules (boundaries, code-quality, frontend, environment)
│   └── skills/                   # GSD loop skills (gsd-discuss, gsd-plan, gsd-execute, gsd-verify, gsd-ship)
│
├── docker-compose.yml            # Local development orchestration (ClickHouse, MCP)
├── flake.nix                     # Nix devShell (Go, SQLite, ClickHouse CLI, Bun)
└── Dockerfile                    # Production multi-stage build embedding web/dist into Go binary
```

---

## 7. GSD Core Skills & Workflows

The `.agents/skills/` directory equips AI assistants with standardized GSD Core lifecycle skills:

* **`gsd-discuss`**: Captures decisions, user intent, and boundaries into `.agents/planning/phases/XX/CONTEXT.md`.
* **`gsd-plan`**: Decomposes phase tasks into atomic, verifiable work units in `.agents/planning/phases/XX/PLAN.md`.
* **`gsd-execute`**: Surgically implements tasks with fresh context and logs execution in `SUMMARY.md`.
* **`gsd-verify`**: Executes comprehensive automated verification and diagnostics in `VERIFICATION.md`.
* **`gsd-ship`**: Archives the phase, updates `ROADMAP.md` and advances `STATE.md`.
* **`mcp-inspect`**: Inspects and validates ClickHouse MCP HTTP connectivity and tools.
