# Phase 03 Plan: Deliveries & Lineage / Dependencies

## Work Units
1. **Ent Schemas (`internal/turso/ent/schema/`)**:
   * `delivery.go`: `id`, `title_id`, `country`, `status` (`PENDING`, `READY_TO_SHIP`, `HOLD`, `SHIPPED`), `target_date`, timestamps.
   * `dependency.go`: `id`, `parent_id`, `child_id`, `dependency_type` (`AUDIO_SYNC`, `SUBTITLE_ALIGNMENT`, `MASTER_DERIVATION`), timestamps.
   * Run `go generate ./internal/turso/ent`.
2. **Domain Models (`pkg/domain/models/`)**:
   * `delivery.go`: `Delivery`, `UpdateDeliveryInput`, validation.
   * `dependency.go`: `Dependency`, `DependencyType`, `LineageNode`, `LineageGraph`.
3. **Turso Store Actions (`internal/turso/`)**:
   * `internal/turso/deliveries/`: `create.go`, `get.go`, `list.go`, `update.go`, `delete.go`.
   * `internal/turso/dependencies/`: `create.go` (with cycle prevention), `get.go`, `list.go`, `graph.go`, `delete.go`.
4. **Echo REST API (`internal/api/`)**:
   * `internal/api/deliveries/routes.go` and handlers (`POST`, `GET`, `PATCH`, `DELETE`).
   * `internal/api/dependencies/routes.go` and handlers (`POST`, `GET /dependencies`, `GET /dependencies/graph/:title_id`, `DELETE`).
   * Register in `internal/api/server.go`.
5. **Verification**:
   * Unit tests for CRUD, cycle detection, graph resolution, and HTTP lifecycle tests with race detector.
