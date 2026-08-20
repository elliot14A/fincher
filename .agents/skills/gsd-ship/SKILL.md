---
name: gsd-ship
description: >-
  Finalizes and archives the active phase, updates ROADMAP.md milestone statuses,
  refreshes STATE.md, and prepares the repository for the next development loop.
---

# GSD Ship Skill

Use this skill once all phase tasks and verification reports pass cleanly.

---

## Workflow Steps

1. **Verify Completion**:
   Ensure `.planning/phases/XX/VERIFICATION.md` exists and reports `PASS` on all criteria.

2. **Update `ROADMAP.md`**:
   Change phase status from `IN_PROGRESS` to `COMPLETED`.

3. **Update `STATE.md`**:
   - Record completion timestamp.
   - Advance active phase pointer to the next scheduled phase in `ROADMAP.md`.

4. **Git Commit Preparation**:
   Summarize the shipped phase for the operator with a concise changelog.
