package main

import (
	"flag"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"graph-med-mcp/src/internal/config"
	"graph-med-mcp/src/internal/tools"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	config.MustLoad(*configFile, &c)

	mcpServer := server.NewMCPServer(
		"knowledge_graph_assistant",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	// setup tools
	mcpTools := []tools.McpTool{
		tools.NewDiseaseSubgraphTool(c),
	}

	for _, tool := range mcpTools {
		mcpServer.AddTool(tool.Definition(), tool.ToolHandlerFunc)
	}

	addr := fmt.Sprintf(":%d", c.Port)

	sseServer := server.NewSSEServer(mcpServer)
	if err := sseServer.Start(addr); err != nil {
		panic(err)
	}
}
