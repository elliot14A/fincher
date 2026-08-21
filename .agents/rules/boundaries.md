# Architectural Boundaries & Responsibilities

These rules define strict package boundaries across Fincher:

1. **AI Agents (`internal/agent/`)**:
   - Strictly **read-only** analytical tools querying ClickHouse exclusively via the official MCP HTTP client (`pkg/mcp`).
   - Must **never** hold direct database credentials or issue SQL mutations (`INSERT`, `UPDATE`, `DELETE`).
   - Must produce typed JSON outputs adhering strictly to Google GenAI / ADK Go schemas.

2. **Policy Engine (`internal/policy/`)**:
   - Pure Go deterministic rule evaluator.
   - Authorizes or denies candidate actions against operational policies stored in the database.
   - Never overridden by AI outputs.

3. **Software Executor (`internal/executor/`)**:
   - Sole authorized mutator of operational state in SQLite.
   - Executes only actions verified as `ALLOWED` by the Policy Engine (or approved by an authorized operator).
   - Emits resulting downstream events for closed-loop progression.

4. **Storage Layer (`internal/turso/`)**:
   - SQLite manages operational state, policy definitions, and audit logs in WAL mode with single-writer concurrency (`SetMaxOpenConns(1)`).
   - Analytical queries flow through the MCP HTTP client (`FINCHER_MCP_URL`).
   - Owns the Ent client and generated code (`internal/turso/ent/`) as a private implementation detail — nothing outside `internal/turso/` may import it directly except through the store's own action functions.

## `pkg/` vs `internal/` Placement Rule

Fincher is a single-binary application (`cmd/fincher`), not a library published for external consumption. The test for where new code goes is not "is this shared across our own layers?" — it's **"could a project outside this module have any legitimate reason to import it?"**

- **`internal/`** — the default. Anything that is an implementation detail of *this* application: the ORM/driver (Ent, Turso/SQLite), config parsing, the store layer, API handlers, policy engine, executor, agents. If every current importer lives inside this module and there's no plan to publish the package standalone, it belongs here — compiler-enforced privacy costs nothing and prevents accidental coupling from outside.
- **`pkg/`** — the exception. Reserve it only for genuinely reusable, dependency-light *contracts*: domain types and vocabulary (`pkg/domain/models`, `pkg/domain/errors`), the MCP client wrapper (`pkg/mcp`), event constants (`pkg/events`), telemetry helpers. Before adding something new here, ask: "would I be comfortable if another Anthropic/internal Go project `go get`'d just this package?" If the answer requires also pulling in Ent, Turso, or Kong config wiring, it isn't a `pkg/`-shaped package — split the contract out or move it to `internal/`.
- **Generated code** (Ent client/schema/migrations) is always an implementation detail of whatever package owns the database, never a public contract on its own — it lives under `internal/turso/ent/`, not `pkg/ent/`.
- Config parsing (Kong structs, `FINCHER_*` env binding) is application-private wiring — it lives under `internal/config/`, not `pkg/domain/config/`.
