# Fincher — Project State

## Current Position
* **Active Milestone**: `Feature 01: Titles & Launch Calendar`
* **Active Phase**: `01-titles`
* **Phase Status**: `COMPLETED`
* **Last Updated**: `2026-08-21`

---

## Active Phase Checklist (`01-titles`)

- [x] **Step 1**: Create `pkg/domain/models/title.go` (Domain Model, Status Enums, Validator integration, Launch Calculations).
- [x] **Step 2**: Create `pkg/domain/config/config.go` & `config_test.go` (Config Struct with `FINCHER_{SERVICE}_*` bindings & validation tests).
- [x] **Step 3**: Create Ent schema in `pkg/ent/schema/title.go` and generate type-safe ORM entities in `pkg/ent/`.
- [x] **Step 4**: Implement `pkg/turso/client.go` & pure functional action functions in `pkg/turso/titles/{create,get,list,update,delete}.go` returning Rust-style `Result[T]` and `Option[T]`.
- [x] **Step 5**: Implement `internal/api/server.go` & pure functional Echo handlers in `internal/api/titles/{create,get,list,update,delete}.go`.
- [x] **Step 6**: Complete integration tests `go test -v ./... -race` verifying full HTTP lifecycle (POST, GET, LIST with filter, PATCH, DELETE, 404).

---

## Key Decisions Record
1. **Ent ORM & Type-Safe Operations**: Migrated from raw SQL strings to `pkg/ent/` with zero manual SQL and automated migrations.
2. **Pure Functional Action Functions**: Turso queries and API handlers are pure functions taking dependencies (`*ent.Client`) as parameters.
3. **Rust-Style `Result[T]` & `Option[T]`**: Error handling across domain and store layers uses `pkg/domain/error.go`.
4. **Decoupled API Error Package**: `internal/api/errors/errors.go` translates domain errors to HTTP JSON responses without import cycles.
5. **Clean Comment Standards**: Minimal, grounded code comments without verbose package names.

---

## Next Milestone
**Feature 02: Masters, Packages & Vendors** (Ent Schemas + Actions + REST API + Verification Tests).
