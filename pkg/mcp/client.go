package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/mcptoolset"
)

type Client struct {
	endpoint  string
	toolset   tool.Toolset
	mcpClient *mcp.Client
}

func NewClient(endpoint string) (*Client, error) {
	cleanEndpoint := strings.TrimRight(endpoint, "/")

	ts, err := mcptoolset.New(mcptoolset.Config{
		Endpoint: cleanEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing adk mcptoolset: %w", err)
	}

	officialClient := mcp.NewClient(&mcp.Implementation{
		Name:    "fincher",
		Version: "1.0.0",
	}, nil)

	return &Client{
		endpoint:  cleanEndpoint,
		toolset:   ts,
		mcpClient: officialClient,
	}, nil
}

func (c *Client) Toolset() tool.Toolset {
	return c.toolset
}

func (c *Client) Ping(ctx context.Context) error {
	transport := &mcp.StreamableClientTransport{
		Endpoint:             c.endpoint,
		DisableStandaloneSSE: true,
	}

	session, err := c.mcpClient.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("connecting to clickhouse mcp server: %w", err)
	}
	defer session.Close()

	_, err = session.ListTools(ctx, nil)
	if err != nil {
		return fmt.Errorf("listing tools from clickhouse mcp server: %w", err)
	}
	return nil
}

func (c *Client) RunQuery(ctx context.Context, sql string) (string, error) {
	transport := &mcp.StreamableClientTransport{
		Endpoint:             c.endpoint,
		DisableStandaloneSSE: true,
	}

	session, err := c.mcpClient.Connect(ctx, transport, nil)
	if err != nil {
		return "", fmt.Errorf("connecting to clickhouse mcp server: %w", err)
	}
	defer session.Close()

	callRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "run_query",
		Arguments: map[string]any{
			"query": strings.TrimSpace(sql),
		},
	})
	if err != nil {
		return "", fmt.Errorf("executing run_query via mcp: %w", err)
	}

	if callRes.IsError {
		var errMsgs []string
		for _, content := range callRes.Content {
			if tc, ok := content.(*mcp.TextContent); ok {
				errMsgs = append(errMsgs, tc.Text)
			}
		}
		return "", fmt.Errorf("mcp tool error: %s", strings.Join(errMsgs, "; "))
	}

	for _, content := range callRes.Content {
		if tc, ok := content.(*mcp.TextContent); ok {
			return tc.Text, nil
		}
	}
	return "", nil
}

func (c *Client) ListTools(ctx context.Context) ([]*mcp.Tool, error) {
	transport := &mcp.StreamableClientTransport{
		Endpoint:             c.endpoint,
		DisableStandaloneSSE: true,
	}

	session, err := c.mcpClient.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to clickhouse mcp server: %w", err)
	}
	defer session.Close()

	toolsRes, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("listing mcp tools: %w", err)
	}
	return toolsRes.Tools, nil
}
