# Fincher Agentic Workflow & Closed-Loop Architecture

This document diagrams the complete multi-agent workflow, from production event trigger to policy gating, software execution, and closed-loop re-evaluation.

---

## Complete Multi-Agent Workflow

```mermaid
flowchart TD
    EV[Incoming Event: API / Simulator<br/>MASTER_UPDATED or QC_FAILED] --> ING[Ingestion Pipeline]
    ING -->|Insert Event| CH[(ClickHouse: fincher.qc_events / asset_events)]
    ING -->|Update State| SQ[(SQLite: fincher.db)]
    ING --> ORCH[Fincher Orchestrator]

    subgraph Parallel Investigation
        ORCH --> HIST[Historian Agent<br/>Query Vendor Track Record & Drift Trends via MCP]
        ORCH --> DEP[Dependency Agent<br/>Query Master Version Lineage & Affected Assets via MCP]
    end

    HIST -->|HistoricalEvidence| MERGE[Evidence Aggregator]
    DEP -->|DependencyEvidence| MERGE

    MERGE --> ANALYST[Incident Analyst Agent<br/>Synthesize Signals, Classify Severity, Attribute Root Cause]
    ANALYST -->|AnalystAssessment| BUNDLE[EvidenceBundle]

    BUNDLE --> PLANNER[Action Planner Agent<br/>Propose Structured ActionPlan: Invalidate, Re-QC, Hold]
    PLANNER -->|ActionPlan| POL[Policy Verification Judge<br/>Evaluate Operational Criteria & Blast Radius]

    POL --> DEC{Policy Decision}
    
    DEC -->|ALLOWED| EXEC[Go Executor]
    DEC -->|HUMAN_APPROVAL_REQUIRED| APPQ[Human Approval Queue]
    DEC -->|DENIED| AUDIT_DENY[Record Denial in Audit Log]

    APPQ -->|Operator Approves| EXEC
    APPQ -->|Operator Rejects| AUDIT_REJ[Record Rejection in Audit Log]

    EXEC -->|Mutate Delivery/Package State| SQ
    EXEC -->|Emit Downstream Event<br/>PACKAGE_INVALIDATED, RE_QC_REQUESTED, DELIVERY_RELEASED| CH
    EXEC -->|Broadcast State & Trace| SSE[SSE Live Stream -> Operations Console]

    EXEC -->|Closed-Loop Feedback| ORCH
```

---

## Sub-Agent Query Contracts (ClickHouse MCP)

```mermaid
sequenceDiagram
    autonumber
    actor Orchestrator
    participant Historian as Historian Agent (ADK Go)
    participant Dependency as Dependency Agent (ADK Go)
    participant MCP as ClickHouse MCP Server
    participant ClickHouse as ClickHouse Analytical DB

    par Parallel Query Phase
        Orchestrator->>Historian: Investigate Vendor & Historical Anomaly (VendorID, CheckName)
        Historian->>MCP: execute_query(SELECT failure_rate FROM qc_events WHERE vendor_id=...)
        MCP->>ClickHouse: Read-only SQL query
        ClickHouse-->>MCP: Aggregate result rows
        MCP-->>Historian: Vendor history & previous incident frequency

    and
        Orchestrator->>Dependency: Investigate Asset Lineage & Affected Languages (TitleID, PackageID)
        Dependency->>MCP: execute_query(SELECT package_id, version, master_version FROM asset_events WHERE ...)
        MCP->>ClickHouse: Read-only SQL query
        ClickHouse-->>MCP: Lineage dataset
        MCP-->>Dependency: Outdated downstream packages
    end

    Historian-->>Orchestrator: HistoricalEvidence
    Dependency-->>Orchestrator: DependencyEvidence
```
