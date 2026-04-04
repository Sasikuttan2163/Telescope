package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/indexer"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var configPath *string

type ListAllToolsParams struct {
}

func ListAllTools(ctx context.Context, req *mcp.CallToolRequest, args ListAllToolsParams) (*mcp.CallToolResult, any, error) {
	var res []byte

	cfg, err := config.GetConfig(*configPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error loading config: " + err.Error()}},
			IsError: true,
		}, nil, nil
	}

	allTools, err := indexer.GetAllIndexedStars(ctx, cfg)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error fetching tools: " + err.Error()}},
			IsError: true,
		}, nil, nil
	}

	for i, v := range allTools {
		res = fmt.Append(res, fmt.Sprintf("Tool %d: %s\n", i, v.Name))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(res),
			},
		},
	}, nil, nil
}

func main() {
	configPath = flag.String("path", "test.json", "Specify path to the JSON configuration file")
	flag.Parse()
	config, err := config.GetConfig(*configPath)
	if err != nil {
		log.Fatal("Fatal error " + err.Error())
	}

	succ, errs := indexer.IndexAllStars(context.Background(), config)
	fmt.Printf("Successfully indexed: %d out of %d\n", succ, len(config.Stars))

	for _, err := range errs {
		fmt.Printf("%s", err.Error())
	}

	allTools, err := indexer.GetAllIndexedStars(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to retrieve indexed tools: " + err.Error())
	}

	for i, v := range allTools {
		fmt.Printf("Tool %d: %s\n", i, v.Name)
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "telescope",
			Version: "v0.0.1",
		},
		nil,
	)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listTools",
		Description: "Use this tool to view the full list of tools available",
	}, ListAllTools)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
