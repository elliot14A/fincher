# GSD Core — Spec-Driven Workflow Rules

These rules govern the Git. Ship. Done (GSD) workflow:

1. **Durable Markdown State**:
   - The `.planning/` directory is the single source of truth for project state.
   - Every phase maintains `CONTEXT.md`, `PLAN.md`, `SUMMARY.md`, and `VERIFICATION.md`.
   - `STATE.md` must be updated on every stage transition to ensure session continuity across AI turns.

2. **Strict Phase Loop**:
   - **Step 1: Discuss** (`gsd-discuss`) $\rightarrow$ Capture decisions in `CONTEXT.md`.
   - **Step 2: Plan** (`gsd-plan`) $\rightarrow$ Decompose into actionable tasks in `PLAN.md`.
   - **Step 3: Execute** (`gsd-execute`) $\rightarrow$ Implement surgically and document in `SUMMARY.md`.
   - **Step 4: Verify** (`gsd-verify`) $\rightarrow$ Run full test suite and diagnostics in `VERIFICATION.md`.
   - **Step 5: Ship** (`gsd-ship`) $\rightarrow$ Update `ROADMAP.md` and advance `STATE.md`.

3. **Anti-Context Rot**:
   - Decompose tasks so that each execution step is small, self-contained, and verified in clean context.

4. **Requirement Traceability**:
   - Every implemented feature or component must map directly to a requirement ID in `.planning/REQUIREMENTS.md` (e.g. `REQ-MCP-01`, `REQ-POLICY-01`).
