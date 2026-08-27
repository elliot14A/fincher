package mcp_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/elliot14A/fincher/pkg/mcp"
)

func TestClient_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mcp.NewClient("http://localhost:8000/mcp")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// 1. Verify ADK Toolset is initialized
	ts := client.Toolset()
	if ts == nil {
		t.Fatal("expected ADK Toolset to be non-nil")
	}
	t.Logf("ADK Toolset initialized: %s", ts.Name())

	// 2. Ping MCP server
	if err := client.Ping(ctx); err != nil {
		t.Skip("skipping integration: MCP server ping failed:", err)
		return
	}

	// 3. List discovered tools via official SDK
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
	if len(tools) == 0 {
		t.Fatal("expected at least 1 tool discovered")
	}

	foundRunQuery := false
	for _, tool := range tools {
		t.Logf("discovered MCP tool: %s", tool.Name)
		if tool.Name == "run_query" {
			foundRunQuery = true
		}
	}
	if !foundRunQuery {
		t.Fatal("expected run_query tool to be discovered")
	}

	// 4. Run read-only query
	res, err := client.RunQuery(ctx, "select count() as total from fincher.events")
	if err != nil {
		t.Fatalf("RunQuery failed: %v", err)
	}
	t.Logf("RunQuery result: %s", res)

	// 5. Mutation query should fail natively via ClickHouse readonly=1
	_, err = client.RunQuery(ctx, "drop table fincher.events")
	if err == nil {
		t.Fatal("expected drop table to fail against read-only MCP server")
	}
	if !strings.Contains(err.Error(), "readonly mode") && !strings.Contains(err.Error(), "READONLY") {
		t.Logf("mutation query failed as expected: %v", err)
	}
}
