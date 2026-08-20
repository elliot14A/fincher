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
2. **Simplicity First & Idiomatic Go**:
   * Implement the cleanest working solution. Strictly NO speculative abstractions, unnecessary generic wrappers, or unrequested framework overhead.
   * Keep package boundaries strict and dependencies minimal.
3. **Surgical Changes**:
   * Touch only what is necessary. Every modified line must trace directly back to the active request.
4. **Goal-Driven Verification**:
   * Frame every implementation with concrete, verifiable success criteria (unit tests, schema verification, MCP tool calls, policy evaluations).

---

## 2. GSD Core Spec-Driven Development (`.planning/`)

All engineering milestones, state transitions, and requirements are managed through the **GSD Core** spec-driven framework in `.planning/`:

1. **`PROJECT.md`**: Foundational domain context, architectural boundaries, and core operating principles.
2. **`REQUIREMENTS.md`**: Deterministic requirement IDs (`REQ-INGEST-01`, `REQ-MCP-01`, `REQ-POLICY-01`, etc.) with verifiable acceptance criteria.
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

### Invariant 5: Complete Auditability
* Every step (Event -> MCP Queries -> Agent Observations -> Merged Evidence -> Action Plan -> Policy Reason -> Execution Result) is written to the SQLite audit log and streamed in real-time to the UI.

---

## 4. Technology Stack & Ecosystem Standards

| Layer / Purpose | Technology / Package | Notes & Invariants |
| :--- | :--- | :--- |
| **Language** | Go (1.24+) | Standard library idiomatic code, explicit error handling |
| **HTTP Framework** | `github.com/labstack/echo/v4` | REST endpoints for event ingestion, incident lifecycle, human approval, SSE for live console |
| **AI Runtime** | Google ADK Go (`google.golang.org/adk/v2`) + Google GenAI (`google.golang.org/genai`) | Programmatic multi-agent orchestration, concurrency, structured schema output |
| **LLM Model** | Gemini 2.5 / Gemini 2.0 Flash / Pro | Structured JSON outputs, fast reasoning over historical analytical evidence |
| **Analytical DB** | ClickHouse | Historical event store: QC logs, asset updates, vendor track records, past incidents |
| **Agent Interface to DB** | Official ClickHouse MCP Server (`ghcr.io/clickhouse/mcp-clickhouse:latest`) | Remote MCP HTTP transport client (`pkg/mcp`). ClickHouse credentials isolated exclusively in MCP container. |
| **Application State & Policies DB** | SQLite + `go-sqlite3` | High-performance operational state & policies table with WAL mode |
| **Policy Engine** | Pure Go Deterministic Evaluator (`internal/policy`) | Evaluates candidate actions against database `policies` table |
| **Config & CLI** | `github.com/alecthomas/kong` | Struct-tag driven environment variables strictly adhering to `FINCHER_{SERVICE}_{SETTING}` |
| **Logging** | `log/slog` (Standard Library) | Structured JSON logs with tracing context (`trace_id`) |
| **Frontend UI** | Embedded Single-Page App (React / SolidJS + Vanilla CSS) | Dark cinematic Operations Console served directly from the Go binary via `embed.FS` |
| **Container & Infra** | Docker Compose + Google Cloud Run + Nix Flake | Production container + local Nix devShell + official ClickHouse MCP HTTP runner |

---

## 5. Architecture & Component Boundaries

```text
fincher/
├── cmd/
│   ├── fincher/                  # Main unified server binary (API, Orchestrator, UI, Simulator)
│   └── seed/                     # Historical dataset generator for ClickHouse & SQLite
│
├── internal/                     # Private domain logic (strict boundaries)
│   ├── api/                      # REST handlers, SSE stream, middleware
│   ├── agent/                    # Multi-agent orchestrator & Google ADK Go sub-agents
│   ├── policy/                   # Deterministic Policy Engine (evaluates DB policies table)
│   ├── executor/                 # Software Execution Engine (state mutations & downstream events)
│   ├── simulator/                # Production media event generator
│   ├── store/                    # Database access layer (SQLite state & policies)
│   └── config/                   # Configuration parsing via Kong
│
├── pkg/                          # Shared contracts, wire types, MCP integration
│   ├── domain/                   # Core types: Event, Delivery, Package, Incident, Action, Evidence
│   ├── events/                   # Event vocabulary & constants
│   ├── mcp/                      # ClickHouse MCP client wrapper for ADK Go agents
│   └── telemetry/                # Structured logging & audit trail helpers
│
├── web/                          # Embedded Operations Console UI
├── data/
│   ├── clickhouse/               # ClickHouse DDL schemas
│   └── seed/                     # Synthetic media historical event seeds
│
├── .agents/                      # AGY project customizations & workflow runbooks
│   ├── plugins/fincher-dev/      # Plugin bundle specification
│   ├── skills/                   # GSD loop skills (gsd-discuss, gsd-plan, gsd-execute, gsd-verify, gsd-ship, mcp-inspect)
│   ├── rules/                    # Rules (boundaries, code-quality, concurrency, environment, gsd-rules)
│   └── workflows/                # Workflow guides (engine-loop, checklist, gsd-phase-loop)
│
├── .planning/                    # GSD Core durable spec & state management
│   ├── PROJECT.md                # Project architecture, constraints, value
│   ├── REQUIREMENTS.md           # Requirement specifications (REQ-*)
│   ├── ROADMAP.md                # Phase roadmap & milestone deliverables
│   ├── STATE.md                  # Active operational phase state tracker
│   └── phases/                   # Phase execution artifacts (CONTEXT, PLAN, SUMMARY, VERIFICATION)
│       └── 01-foundation/
│
├── docker-compose.yml            # Local development orchestration (ClickHouse, MCP)
├── flake.nix                     # Nix devShell (Go, SQLite, ClickHouse CLI, Node.js)
├── Dockerfile                    # Production multi-stage build for Google Cloud Run
└── .golangci.yml                 # Strict Go linter & guardrail configuration
```

---

## 6. AGY & GSD Core Skills

The `.agents/skills/` directory equips the AI assistant with GSD Core loop skills and specialized inspectors:

* **`gsd-discuss`**: Captures decisions, user intent, and boundaries into `.planning/phases/XX/CONTEXT.md`.
* **`gsd-plan`**: Decomposes phase tasks into atomic, verifiable work units in `.planning/phases/XX/PLAN.md`.
* **`gsd-execute`**: Surgically implements tasks with fresh context and logs execution in `SUMMARY.md`.
* **`gsd-verify`**: Executes comprehensive automated verification and diagnostics in `VERIFICATION.md`.
* **`gsd-ship`**: Archives the phase, updates `ROADMAP.md` and advances `STATE.md`.
* **`mcp-inspect`**: Inspects and validates ClickHouse MCP HTTP connectivity and tools.

---

## 7. Engineering Quality Standards

1. **Event-Driven Deterministism**:
   State transitions and actions emerge strictly from structured events, historical evidence from ClickHouse MCP, and database-backed policies.
2. **Operations Console**:
   Fincher is an autonomous operations engine, not a conversational chatbot.
3. **Strict Idiomatic Go**:
   * No unnecessary reflection, deep inheritance, or generic wrapper boilerplate.
   * Use explicit error handling with wrapped errors: `fmt.Errorf("evaluating policy %s: %w", policyID, err)`.
   * Use `context.Context` propagation across all HTTP handlers, agent steps, MCP tool calls, and DB transactions.
4. **Structured JSON Output from LLMs**:
   Enforce strict JSON schema parsing on all Gemini outputs using structured schemas in the Google GenAI / ADK Go SDK.
5. **No Placeholders in UI or Code**:
   Build functional, high-polish components with realistic media production datasets and real-time event updates.
6. **Strict `FINCHER_{SERVICE}_*` Environment Naming**:
   All environment variables must follow the `FINCHER_{SERVICE}_{SETTING}` convention (e.g., `FINCHER_SERVER_PORT`, `FINCHER_MCP_URL`, `FINCHER_GEMINI_API_KEY`).
7. **MCP Credential Isolation**:
   ClickHouse database connection credentials must remain exclusively in the Docker Compose MCP service configuration. Go backend and ADK agents query ClickHouse strictly via the remote MCP HTTP endpoint.
