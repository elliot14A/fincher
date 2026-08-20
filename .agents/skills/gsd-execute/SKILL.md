---
name: gsd-execute
description: >-
  Executes the tasks defined in .planning/phases/XX/PLAN.md with fresh context, surgical edits,
  and compiles execution records into .planning/phases/XX/SUMMARY.md.
---

# GSD Execute Skill

Use this skill to implement the planned code changes sequentially.

---

## Workflow Steps

1. **Read `PLAN.md`**:
   Load `.planning/phases/XX/PLAN.md` and identify the next incomplete task.

2. **Execute Surgically**:
   - Write boilerplate, type definitions, and component logic strictly as planned.
   - Respect package boundaries and Go idioms (explicit error wrapping, context propagation).
   - Write corresponding unit tests concurrently.

3. **Incremental Verification**:
   - Run compilation and targeted tests after each task:
     ```bash
     go build ./...
     go test -v ./... -race
     ```

4. **Write `SUMMARY.md`**:
   Document the completed tasks, modified files, deviations (if any), and key outcomes in `.planning/phases/XX/SUMMARY.md`.

5. **Update `STATE.md`**:
   Check off completed items in `.planning/STATE.md`.
