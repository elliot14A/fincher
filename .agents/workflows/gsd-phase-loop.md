# GSD Core — Phase Loop Workflow

This document illustrates the 5-step GSD Core development lifecycle:

```mermaid
flowchart TD
    subgraph Phase Lifecycle
        D[1. DISCUSS<br/>gsd-discuss<br/>Capture decisions in CONTEXT.md] --> P[2. PLAN<br/>gsd-plan<br/>Decompose tasks into PLAN.md]
        P --> E[3. EXECUTE<br/>gsd-execute<br/>Surgical code implementation & SUMMARY.md]
        E --> V[4. VERIFY<br/>gsd-verify<br/>Automated verification & VERIFICATION.md]
        V --> S[5. SHIP<br/>gsd-ship<br/>Update ROADMAP.md & STATE.md]
    end

    S -->|Next Phase| D
```

---

## Artifact Mapping

| Stage | Input Artifacts | Output Artifact | Key Responsibility |
| :--- | :--- | :--- | :--- |
| **Discuss** | `PROJECT.md`, `REQUIREMENTS.md` | `phases/XX/CONTEXT.md` | Align on intent, choices, and architectural boundaries. |
| **Plan** | `CONTEXT.md`, `REQUIREMENTS.md` | `phases/XX/PLAN.md` | Break work into atomic, verifiable tasks. |
| **Execute** | `PLAN.md` | `phases/XX/SUMMARY.md` | Implement code and unit tests with fresh context. |
| **Verify** | Code, `PLAN.md` | `phases/XX/VERIFICATION.md` | Run compilation, vet, race tests, and diagnostics. |
| **Ship** | `VERIFICATION.md` | Updated `ROADMAP.md` & `STATE.md` | Archive phase and advance active milestone pointer. |
