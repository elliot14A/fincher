# Phase 03 Summary: Deliveries & Lineage / Dependencies

## Summary of Completed Work
* **Deliveries (`internal/turso/deliveries/`, `internal/api/deliveries/`)**:
  * Implemented pure functional Turso actions for country release targeting (`Create`, `Get`, `List`, `Update`, `Delete`).
  * Created REST endpoints (`POST /deliveries`, `GET /deliveries`, `GET /deliveries/:id`, `PATCH /deliveries/:id`, `DELETE /deliveries/:id`).
* **Lineage & Dependencies (`internal/turso/dependencies/`, `internal/api/dependencies/`)**:
  * Implemented pure functional Turso actions with DAG cycle prevention logic.
  * Implemented recursive `GetLineageGraph` tree builder resolving parent-child component relationships for a title.
  * Created REST endpoints (`POST /dependencies`, `GET /dependencies`, `GET /dependencies/graph/:title_id`, `DELETE /dependencies/:id`).
* **Restructuring**:
  * Moved `pkg/turso` and `pkg/ent` to `internal/turso` and `internal/turso/ent`.
  * Moved `pkg/domain/config` to `internal/config`.
* **Testing & Verification**:
  * Unit tests and HTTP lifecycle integration tests passing 100% with race detector (`-race`).
  * Invariant verification script `./.agents/scripts/verify.sh` passes 100%.
