# Phase 02 Plan: Masters, Packages & Vendors

## Work Units
1. **Ent Schemas (`internal/turso/ent/schema/`)**:
   * `master.go`: Fields `id`, `title_id`, `version`, `supersedes_version`, `created_at`.
   * `vendor.go`: Fields `id`, `name`, `specialty`, timestamps.
   * `media_package.go`: Fields `id`, `title_id`, `component`, `language`, `version`, `vendor_id`, `derived_from_master_version`, `redelivery_count`, `status`, timestamps.
2. **Domain Models (`pkg/domain/models/`)**:
   * `master.go`, `package.go`, `vendor.go` with validation tags and partial update input models (`UpdateVendorInput`, `UpdatePackageInput`).
   * Staleness method `Package.IsStaleAgainst(activeMasterVersion)`.
3. **Turso Store Actions (`internal/turso/`)**:
   * Pure functional actions for `masters/`, `packages/`, `vendors/` returning `Result[T]` with domain model mapping.
   * Master creation transactionally updates parent `Title.current_master_version`.
   * Error mapper `MapEntError` handles both foreign key and unique constraint failures distinctly.
4. **Echo REST API (`internal/api/`)**:
   * Handlers and `routes.go` for `masters/`, `packages/`, `vendors/`.
   * Modular registration in `internal/api/server.go`.
5. **Verification**:
   * Unit tests for all CRUD actions, FK constraints, staleness detection, and HTTP lifecycles.
