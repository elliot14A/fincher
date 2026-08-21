# Phase 03 Context: Deliveries & Lineage / Dependencies

## Objectives
* Deliver country release targeting via `Delivery` entity (`id`, `title_id`, `country`, `status`, `target_date`).
* Deliver package dependency graph via `Dependency` entity (`id`, `parent_id`, `child_id`, `dependency_type`).
* Implement cycle detection and recursive lineage graph traversals.
* Pure functional Turso actions returning domain models (`Result[T]`, `Option[T]`).
* REST API endpoints with modular `routes.go` per entity under `internal/api/`.

## Architectural Decisions
1. **Deliveries**:
   * Represents country/territory release target (`title_id`, `country`, `status`, `target_date`).
   * Status: `PENDING`, `READY_TO_SHIP`, `HOLD`, `SHIPPED`.
2. **Dependencies**:
   * Links parent component to dependent child component (e.g. Master Video -> Spanish Audio Dub, Spanish Audio -> Spanish Subtitles).
   * DependencyType: `AUDIO_SYNC`, `SUBTITLE_ALIGNMENT`, `MASTER_DERIVATION`.
   * Graph must be a Directed Acyclic Graph (DAG); cyclical dependencies are rejected at insertion time.
3. **Discrete Route Modules**:
   * `internal/api/deliveries/routes.go`
   * `internal/api/dependencies/routes.go`
