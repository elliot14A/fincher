# Fincher — Project Specification

## 1. Executive Summary & Core Value
**Fincher** is an event-driven, autonomous delivery-integrity operations engine for film and television post-production release pipelines. It continuously monitors production and QC events, conducts parallel historical investigations via the official ClickHouse Model Context Protocol (MCP) server, synthesizes evidence with Gemini, generates remediation action plans, and gates all state mutations behind deterministic database-backed operational policies.

> **Core Operating Principle**:  
> **AI investigates and plans. Database policies decide what is permitted. Software executes and emits downstream events.**

---

## 2. Problem Statement & Domain Context
In modern media distribution (streaming platforms, theatrical releases, linear broadcast):
* Upstream video/audio recuts (e.g. director trim +3 frames) frequently invalidate downstream localized assets (5.1 audio dubs, timed subtitles, forced narratives).
* Undetected version drift and vendor-specific sync defects cause broadcast delays, SLA penalties, and silent release failures.
* Traditional automation lacks historical context awareness (e.g. vendor failure rates, asset lineage), while pure LLM chatbots lack deterministic safety and operational execution authority.

Fincher solves this by creating a **closed-loop autonomous operations engine** combining:
1. Fast historical analytical querying (ClickHouse via MCP).
2. Parallel multi-agent investigation (Historian, Dependency, Incident Analyst, Action Planner).
3. Deterministic policy gating (pure Go evaluator backed by operational policies table).
4. Automated state mutations and closed-loop re-evaluation.

---

## 3. System Architecture & Component Boundaries

```text
               Production Event (REST API / Simulator)
                                │
                                ▼
                       [ Ingestion Layer ]
                                │
                ┌───────────────┴───────────────┐
                ▼                               ▼
       [ SQLite: fincher.db ]         [ ClickHouse (MCP) ]
   (State & Policy Tables)             (Historical Event Store)
                │
                ▼
      [ Fincher Orchestrator ]
                │
    ┌───────────┴───────────┐
    ▼                       ▼
[ Historian Agent ]   [ Dependency Agent ]
 (ClickHouse MCP)       (ClickHouse MCP)
    │                       │
    └───────────┬───────────┘
                ▼
     [ Evidence Aggregator ]
                │
                ▼
     [ Incident Analyst Agent ]
                │
                ▼
      [ Action Planner Agent ]
                │
                ▼
      [ Policy Engine ] (Evaluates rules from DB policies table)
                │
        ┌───────┴───────┐
        ▼               ▼
   [ ALLOWED ]   [ HUMAN APPROVAL ]
        │               │
        ▼               ▼
  [ Go Executor ]  [ Approval Queue (UI) ]
        │
        ▼
 Mutate State & Emit Downstream Event (PACKAGE_INVALIDATED, RE_QC, etc.)
        │
        ▼
  [ Closed-Loop Feedback to Orchestrator ]
```

---

## 4. Technical Stack & Invariants

| Component | Standard Technology | Architectural Rule |
| :--- | :--- | :--- |
| **Language** | Go (1.24+) | Standard library idioms, explicit error wrapping, `context.Context` everywhere |
| **Application State & Policies DB** | SQLite (`mattn/go-sqlite3`) | High-performance operational state & policies table in WAL mode with `SetMaxOpenConns(1)` |
| **Historical Analytics DB** | ClickHouse (`clickhouse:24.3-alpine`) | Time-series event store (`qc_events`, `asset_events`, `delivery_events`, `incident_events`) |
| **Agent DB Interface** | Official ClickHouse MCP Server (`ghcr.io/clickhouse/mcp-clickhouse:latest`) | Remote MCP HTTP transport (`/mcp`). ClickHouse credentials isolated strictly in MCP container |
| **AI Multi-Agent Runtime** | Google ADK Go (`google.golang.org/adk/v2`) + Google GenAI (`google.golang.org/genai`) | Code-first multi-agent orchestration, parallel goroutines, typed JSON schemas |
| **Policy Engine** | Pure Go Deterministic Evaluator | Maps database policies to candidate action gating without AI overrides |
| **HTTP Framework** | `github.com/labstack/echo/v4` | Event ingestion, approvals REST API, SSE real-time event streaming |
| **Configuration** | `github.com/alecthomas/kong` | Strict `FINCHER_{SERVICE}_{SETTING}` environment variable naming |

---

## 5. Non-Negotiable Standards
1. **Event-Driven Resolution**: All decisions emerge dynamically from events, ClickHouse history, and database policies.
2. **Operations Console**: Fincher is an autonomous operations engine, not a conversational chatbot.
3. **Read-Only Agents**: Agents are strictly read-only. State mutations are executed solely by the Go executor upon policy authorization.
4. **GSD Core Lifecycle**: Every phase follows the GSD phase loop (`DISCUSS` $\rightarrow$ `PLAN` $\rightarrow$ `EXECUTE` $\rightarrow$ `VERIFY` $\rightarrow$ `SHIP`) with persistent markdown state in `.planning/`.
