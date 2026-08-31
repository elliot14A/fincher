# Fincher

**Autonomous Delivery-Integrity Operations Engine for Media Production**

Fincher is an event-driven operations workflow engine designed for film and television post-production release pipelines. It continuously consumes media lifecycle events, investigates historical context via ClickHouse using the official ClickHouse Model Context Protocol (MCP) server, coordinates a multi-agent Gemini investigation workflow, generates remediation plans, and verifies all actions through a policy judge before execution.

---

## Core Philosophy

> **AI investigates and plans. Scoped judges verify. Software executes.**

Fincher is an autonomous operations engine operating a closed-loop operations lifecycle.

---

## Workflow Architecture

```text
               Production Event (API / Simulator)
                              │
                              ▼
                     [ Ingestion Layer ] ───► Write Event to ClickHouse & SQLite
                              │
                              ▼
                  [ Fincher Orchestrator ]
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
   [ Historian Agent ] [ Dependency Agent ] [ Incident Analyst ]
     (ClickHouse MCP)    (ClickHouse MCP)   (Context & Graph)
          │                   │                   │
          └───────────────────┼───────────────────┘
                              ▼
                   [ Evidence Aggregator ]
                              │
                              ▼
                    [ Action Planner ]
               (Proposes candidate actions)
                              │
                              ▼
                [ Policy Verification Judge ]
            (Evaluates candidate actions vs operational criteria)
                               │
                ┌──────────────┴──────────────┐
                ▼                             ▼
         [ APPROVED ]                  [ HUMAN APPROVAL ]
                │                             │
                ▼                             ▼
        [ Go Executor ]              [ Approval Queue (UI) ]
                │                             │ (Upon human review)
                ▼                             │
       Mutate State & Emit                    ▼
           New Event ─────────────────────────┘
                │
                ▼
       [ Closed Loop ] ──► Feed back to Orchestrator (Re-evaluates until terminal state)
```

---

## Technology Stack

* **Backend / Orchestrator**: Go 1.24+ (`github.com/labstack/echo/v4`)
* **AI Runtime**: Google ADK Go (`google.golang.org/adk/v2`) + Google GenAI (`google.golang.org/genai`)
* **Historical Analytics DB**: ClickHouse
* **Agent-to-DB Interface**: Official ClickHouse Model Context Protocol (MCP) Server
* **Application State DB**: SQLite (WAL mode)
* **Policy Verifier**: Scoped Gemini Judge with bounded self-correction loop
* **Operations Console**: Embedded Web UI served via Go `embed.FS`

---

## Quick Start (Local Development)

### Prerequisites
* Go 1.24+
* Docker & Docker Compose

### 1. Environment Setup
```bash
cp .env.example .env
# Set FINCHER_GEMINI_API_KEY in .env
```

### 2. Start Infrastructure
```bash
docker compose up -d clickhouse mcp-clickhouse
```

### 3. Verify System Setup & MCP Connection
```bash
./.agents/scripts/verify.sh
```

---

## Event Vocabulary

* **Production**: `MASTER_UPDATED`, `ASSET_UPDATED`, `PACKAGE_CREATED`, `PACKAGE_INVALIDATED`, `PACKAGE_APPROVED`, `PACKAGE_REDELIVERED`
* **QC**: `QC_STARTED`, `QC_COMPLETED`, `QC_PASSED`, `QC_FAILED`
* **Workflow**: `DELIVERY_HELD`, `DELIVERY_RELEASED`, `RE_QC_REQUESTED`, `INCIDENT_CREATED`, `INCIDENT_RESOLVED`
