---
name: mcp-inspect
description: >-
  Inspects, tests, and validates the official ClickHouse MCP HTTP server connection,
  read-only analytical queries, table schema introspections, and tool call responses
  for the Historian and Dependency sub-agents. Use when debugging or extending MCP integration.
---

# ClickHouse MCP Inspector

Use this procedure to verify and debug communication with the official ClickHouse Model Context Protocol (MCP) server container (`ghcr.io/clickhouse/mcp-clickhouse:latest`).

---

## Safety Checks
* **Read-Only Rule**: Agents must ONLY execute read-only queries (`SELECT`, `SHOW TABLES`, `DESCRIBE`). Never execute `INSERT`, `ALTER`, `DROP`, or `TRUNCATE`.
* **Credential Isolation**: Database credentials reside exclusively in the Docker Compose MCP service configuration. Go agents communicate strictly over HTTP transport at `FINCHER_MCP_URL`.

---

## MCP Inspection Steps

### 1. Check Container Health
Verify ClickHouse and MCP containers:
```bash
docker compose ps
curl -s http://localhost:8123/ping
curl -s http://localhost:8000/health
```

### 2. Verify MCP Protocol Handshake
Send initialization request to verify capabilities:
```bash
curl -s -X POST http://localhost:8000/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"mcp-inspector","version":"1.0.0"}}}'
```

### 3. List Discovered MCP Tools
Query available tools (`run_query`, `list_tables`, `list_databases`):
```bash
curl -s -X POST http://localhost:8000/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
```

### 4. Execute Read-Only Analytical Query via Go Client
Run the live integration test suite:
```bash
go test -v ./pkg/mcp/... -race
```

### 5. Verify Output Sanitization
Ensure query responses returned to Gemini / ADK Go are properly formatted JSON objects, bounded in size, and structured without raw unparsed SQL leaks.
