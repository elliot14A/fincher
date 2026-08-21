# Phase 02 Context: Masters, Packages & Vendors

## Objectives
* Deliver relational entity graph for media mastering and localized deliverables: `masters`, `packages`, `vendors`.
* Establish immutable master versioning with `supersedes_version` chains and transactional title synchronization.
* Enforce relational DB constraints (Title 1──* Master, Title 1──* Package, Vendor 1──* Package).
* Implement pure functional actions taking `*ent.Client` and returning domain models (`Result[T]` / `Option[T]`).
* Expose REST endpoints with modular `routes.go` per entity under `internal/api/`.

## Architectural Decisions
1. **Masters in Turso / SQLite**:
   * Stored in Turso for instantaneous relational FK integrity against titles and packages, with sub-millisecond lookups.
   * `Master` creation transactionally updates the parent `Title.current_master_version`.
   * Masters are append-only editorial cuts (`V01`, `V05`, `V12`, `V13`).
2. **MediaPackage Component State & Staleness**:
   * Packages represent localized deliverables (`VIDEO`, `AUDIO`, `SUBTITLE`, `METADATA`).
   * Staleness against active master cut is evaluated via `Package.IsStaleAgainst(title.CurrentMasterVersion)`.
3. **Discrete Route Modules**:
   * Each entity package in `internal/api/{titles,masters,vendors,packages}` exports `RegisterRoutes(g *echo.Group, client *ent.Client)` consumed by `Server`.
