---
name: gsd-plan
description: >-
  Creates and verifies an implementation plan for the active phase in .planning/phases/XX/PLAN.md.
  Decomposes tasks into small, verifiable chunks and verifies them against project invariants.
---

# GSD Plan Skill

Use this skill to research, structure, and decompose implementation tasks for the current phase.

---

## Workflow Steps

1. **Read Inputs**:
   - `.planning/PROJECT.md`
   - `.planning/REQUIREMENTS.md`
   - `.planning/phases/XX/CONTEXT.md`

2. **Decompose Tasks**:
   Break the phase down into discrete, atomic tasks:
   - Target files (types, logic, tests)
   - Scope and implementation details
   - Concrete verification criteria (unit tests, CLI runs, HTTP responses)

3. **Verify Plan Feasibility**:
   - Check against system invariants (Read-only agents, MCP credential isolation, deterministic policies).
   - Ensure tasks fit within a clean execution context without cognitive bloat.

4. **Write `PLAN.md`**:
   Save to `.planning/phases/XX/PLAN.md`.

5. **Update `STATE.md`**:
   Update active checklist and status in `.planning/STATE.md`.
