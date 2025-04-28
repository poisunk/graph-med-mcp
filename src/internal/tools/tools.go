package tools

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type McpTool interface {
	Name() string
	Definition() mcp.Tool
	ToolHandlerFunc(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
