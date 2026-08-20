# Environment & Configuration Rules

These rules govern configuration and secret isolation across Fincher:

1. **Mandatory `FINCHER_{SERVICE}_{SETTING}` Prefix**:
   - Every environment variable parsed by Fincher must strictly follow the `FINCHER_{SERVICE}_{SETTING}` convention (e.g., `FINCHER_SERVER_PORT`, `FINCHER_MCP_URL`, `FINCHER_GEMINI_API_KEY`, `FINCHER_SQLITE_PATH`, `FINCHER_SIMULATOR_ENABLE`).

2. **MCP Credential Isolation**:
   - ClickHouse analytical database credentials (`CLICKHOUSE_HOST`, `CLICKHOUSE_PORT`, `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`) must remain **exclusively** inside the MCP service container in `docker-compose.yml`.
   - The Go backend and AI agents communicate with ClickHouse solely via the remote MCP HTTP transport client (`FINCHER_MCP_URL`).

3. **No Hardcoded Defaults for Secrets**:
   - API keys and tokens must never have default values in code or configuration structs.
