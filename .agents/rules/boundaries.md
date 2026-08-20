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

4. **Storage Layer (`internal/store/`)**:
   - SQLite manages operational state, policy definitions, and audit logs in WAL mode with single-writer concurrency (`SetMaxOpenConns(1)`).
   - Analytical queries flow through the MCP HTTP client (`FINCHER_MCP_URL`).
