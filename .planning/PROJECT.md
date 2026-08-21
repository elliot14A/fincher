# Fincher — Project Specification (Workflow DAG Edition)

## 1. Executive Summary & Core Value
**Fincher** is an autonomous delivery-integrity workflow engine for **LUME**, a premier streaming service. The console is organized around a **launch calendar** of upcoming releases with firm premiere dates, so every incident carries real operational stakes (*"Will Eclipse ship across 40 territories in time for Friday's global premiere?"*).

Fincher runs a structured **Studio Pipeline DAG** (Directed Acyclic Graph) of specialized nodes (data queries, Gemini Flash/Pro agent investigators, hybrid decision nodes, concrete state actions, and drafted stakeholder dispatches).

> **Core Operating Principle**:  
> **ClickHouse holds deep historical telemetry. Scheduled DAG nodes query OLAP views via MCP, AI synthesizes evidence with deterministic numbers inside the Decision Node, and software executes state mutations.**

---

## 2. The Four-Beat Story & The Disagreement Money Shot
Every run is legible in ten seconds:
> **Delta detected → Agents query ClickHouse history via MCP → Decision Node synthesizes numbers + context → State action executed & stakeholder dispatch drafted.**

### The Hero Scenario: Disagreement Panel
* **Context**: *Eclipse* master bumps to V13 two days before a 40-territory launch. Vendor A's Spanish audio produces a borderline sync drift (110ms against a 120ms tolerance).
* **Naive Rule Automation**: PASS (110ms is within the static 120ms threshold).
* **Fincher Decision Node (Gemini Pro + ClickHouse MVs)**: HOLD. Combines quantitative inputs (36h to launch, 40 territories, redelivery count = 2) with historical evidence (Vendor A's rolling defect rate is escalating, upstream master superseded) to hold delivery and draft a stakeholder notice.

---

## 3. Architecture & Data Flow

```text
       Cloud Scheduler (Periodic ticks: every N mins)
                    │  HTTP POST /workflows/{id}/run
                    ▼
     ┌──────────────────────────────────────────────┐
     │          Single-Request DAG Runner           │
     │                                              │
     │ 1. [ schedule_trigger ]                      │
     │ 2. [ delta_gate ] (0 LLM cost exit if idle) │
     │       │ (if delta found)                     │
     │ 3. Parallel Query Nodes (ClickHouse MCP)    │
     │    - vendor_reliability_query                │
     │    - lineage_query                           │
     │    - redelivery_query                        │
     │    - recent_master_change_query              │
     │    - time_to_premiere (computed live)        │
     │       │                                      │
     │ 4. Parallel Agent Nodes (Gemini Flash)       │
     │    - vendor_risk_agent                       │
     │    - dependency_impact_agent                 │
     │       │                                      │
     │ 5. [ assessment_agent ] (Gemini Pro)         │
     │       │                                      │
     │ 6. [ decision_node ] (Gemini Pro)            │
     │    (Combines deterministic inputs + agents)  │
     │    Branches: HOLD / RE_QC / RELEASE / NONE   │
     │       │                                      │
     │ 7. Action & Notification Nodes               │
     │    - hold_delivery_action / release_action   │
     │    - stakeholder_notice_action (Flash draft) │
     │ 8. [ event_emitter ] (Sink to ClickHouse)    │
     └──────────────────────┬───────────────────────┘
                            ▼
     Turso (State, Runs, Node Executions, Notifications)
```

---

## 4. Technical Stack & Invariants

| Component | Standard Technology | Architectural Invariant |
| :--- | :--- | :--- |
| **Language** | Go (1.24+) | Standard library idioms, explicit error wrapping, `context.Context` propagation |
| **Runtime & Cloud** | Google Cloud Run (**cold / scale-to-zero**) | Stateless single-request DAG execution (~2–4s), $0 idle cost |
| **Application State Store** | Turso / libSQL (`@libsql/client` / Go driver) | Serverless HTTP SQLite storing titles, packages, runs, node executions, notifications |
| **Historical Analytics DB** | ClickHouse Cloud / Local Container (`24.3-alpine`) | 250k+ append-only QC/asset events + 4 Materialized Views |
| **Agent DB Interface** | Official ClickHouse MCP Server (`mcp-clickhouse`) | Remote MCP HTTP transport (`/mcp`). ClickHouse credentials isolated exclusively in MCP service |
| **AI Models** | Google GenAI SDK (`google.golang.org/genai` / ADK Go) | Gemini 2.5 Flash for query agents & draft notifications; Gemini 2.5 Pro for Assessment & Decision nodes |
| **Decisions** | Hybrid `decision_node` | Combines quantitative MV metrics + agent synthesis into fixed branches (`HOLD`, `RE_QC`, `RELEASE`, `NONE`) |
| **Operator Assistant** | Read-first Docent Assistant | Answers "What's releasing this weekend?" and explains past run decisions with SQL query citations |
| **HTTP API & SSE** | `github.com/labstack/echo/v4` | REST endpoints for titles, workflow runs, and SSE for real-time node stepping |
| **Frontend UI** | Preact Single-Page App | Launch Calendar + Studio Pipeline Visualizer + Disagreement Panel + Stakeholder Dispatch Preview |
| **Configuration** | `github.com/alecthomas/kong` | Strict `FINCHER_{SERVICE}_{SETTING}` environment variable bindings |

---

## 5. Non-Negotiable Rules
1. **ClickHouse via MCP Only**: All analytical history queries at runtime route through the official `mcp-clickhouse` server.
2. **Google Cloud AI Tooling Only**: Exclusively use Google Gemini SDKs.
3. **Delta Gate Cost Protection**: Routine ticks with no new events exit in $< 10\text{ms}$ at $0 LLM cost.
4. **Mocked Notifications**: Stakeholder emails/dispatches are drafted with realistic LLM content and stored as `drafted` for in-browser review.
5. **Acyclic Workflow Execution**: A run is a single DAG execution; closed loops (hold $\rightarrow$ re-QC $\rightarrow$ release) occur across successive runs.
