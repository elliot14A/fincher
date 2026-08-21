---
name: gsd-discuss
description: >-
  Initiates the discuss stage of the GSD Core workflow. Captures implementation decisions,
  user intent, architectural choices, and boundaries into .agents/planning/phases/XX/CONTEXT.md before code is planned or written.
---

# GSD Discuss Skill

Use this skill when starting a new phase or feature to establish clarity, align on architecture, and prevent downstream rework.

---

## Workflow Steps

1. **Identify the Active Phase**:
   Read `.agents/planning/STATE.md` to identify the current milestone and phase directory (e.g. `.agents/planning/phases/01-foundation/`).

2. **Capture Key Design Decisions**:
   Interview the user or evaluate the active request for:
   - Specific component boundaries and responsibilities.
   - Database schema choices and migrations.
   - API endpoints, request/response formats, and SSE events.
   - Invariants and non-negotiables (e.g. Read-Only Agents, `FINCHER_{SERVICE}_*` env names, MCP isolation).

3. **Write `CONTEXT.md`**:
   Persist all decisions in `.agents/planning/phases/XX/CONTEXT.md` with:
   - Objective & Scope
   - Design Decisions Table
   - Architectural Constraints & Invariants
   - Open Questions / Risk Mitigations

4. **Update `STATE.md`**:
   Record the completion of the discussion phase in `.agents/planning/STATE.md`.
