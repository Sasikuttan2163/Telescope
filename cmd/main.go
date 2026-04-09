package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/embed"
	"github.com/Sasikuttan2163/Telescope/internal/indexer"
	"github.com/Sasikuttan2163/Telescope/internal/transport"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var mainConfig config.MainConfig

type ListAllToolsParams struct {
}

type SearchToolsParams struct {
	Query string `json:"query" jsonschema:"Keywords for the CAPABILITY required. Do NOT include search terms for the user's specific topic. Correct: 'web search', 'github API', 'file system'. Incorrect: 'search for hotels in Chennai'. Only use this to discover which tools exist."`
}

type CallToolParams struct {
	ToolName string         `json:"toolName" jsonschema:"The FULL CALL_ID returned by the SearchTools function (e.g., 'serverName__toolName'). NEVER guess this name; only use names found via SearchTools."`
	Input    map[string]any `json:"inputJson" jsonschema:"The arguments for the tool, following the tool's specific Input Schema."`
}

func ListAllTools(ctx context.Context, req *mcp.CallToolRequest, args ListAllToolsParams) (*mcp.CallToolResult, any, error) {
	var res []byte

	allTools, err := indexer.GetAllIndexedStars(ctx, mainConfig)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error fetching tools: " + err.Error()}},
			IsError: true,
		}, nil, nil
	}

	for i, v := range allTools {
		res = fmt.Append(res, fmt.Sprintf("Tool %d: %s__%s\n", i, v.ServerName, v.Name))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(res),
			},
		},
	}, nil, nil
}

func SearchTools(ctx context.Context, req *mcp.CallToolRequest, args SearchToolsParams) (*mcp.CallToolResult, any, error) {
	var res []byte

	queryVectorCtx, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
	defer cancel()

	queryVector, err := embed.OllamaGetQueryVector(queryVectorCtx, mainConfig.Ollama.Host, mainConfig.Ollama.Port, &mainConfig.Ollama.NumGpu, mainConfig.Ollama.Model, args.Query)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error occurred while getting vectors from ollama. Full trace: %s", err.Error()),
				},
			},
		}, nil, err
	}

	topKToolsCtx, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
	defer cancel()
	topKTools, err := indexer.GetTopKTools(topKToolsCtx, mainConfig, queryVector)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error occurred while fetching top-K tools. Full error trace: %s", err.Error()),
				},
			},
		}, nil, err
	}

	for i, v := range topKTools {
		fullCallID := fmt.Sprintf("%s__%s", v.ServerName, v.Name)

		schemaBytes, _ := json.Marshal(v.InputSchema)
		res = fmt.Append(res, fmt.Sprintf(
			"Rank %d:\n  CALL_ID: %s\n  Description: %s\n  Input Schema: %s\n\n",
			i, fullCallID, v.Description, string(schemaBytes),
		))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(res),
			},
		},
	}, nil, nil
}

func CallTool(ctx context.Context, req *mcp.CallToolRequest, args CallToolParams) (*mcp.CallToolResult, any, error) {
	parts := strings.Split(args.ToolName, "__")
	if len(parts) != 2 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Invalid tool name format. Expected: serverName__toolName"}},
		}, nil, nil
	}

	serverName := parts[0]
	toolName := parts[1]

	star := transport.GetStarByName(mainConfig.Stars, serverName)
	if star == nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Server %s not found in configuration", serverName)}},
		}, nil, nil
	}

	inputBytes, err := json.Marshal(args.Input)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Failed to marshal tool input: " + err.Error()}},
		}, nil, err
	}

	result, err := transport.CallToolOnStar(ctx, *star, toolName, inputBytes)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error calling tool: " + err.Error()}},
		}, nil, err
	}

	return result, nil, nil
}

func main() {
	configPath := flag.String("path", "test.json", "Specify path to the JSON configuration file")
	flag.Parse()
	var err error
	mainConfig, err = config.GetConfig(*configPath)
	if err != nil {
		log.Fatal("Fatal error " + err.Error())
	}

	succ, errs := indexer.IndexAllStars(context.Background(), mainConfig)
	fmt.Printf("Successfully indexed: %d out of %d\n", succ, len(mainConfig.Stars))

	for _, err := range errs {
		fmt.Printf("%s", err.Error())
	}

	allTools, err := indexer.GetAllIndexedStars(context.Background(), mainConfig)
	if err != nil {
		log.Fatal("Failed to retrieve indexed tools: " + err.Error())
	}

	for i, v := range allTools {
		fmt.Printf("Tool %d: %s__%s\n", i, v.ServerName, v.Name)
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "telescope",
			Version: "v0.0.1",
		},
		nil,
	)

	// mcp.AddTool(server, &mcp.Tool{
	// 	Name:        "listTools",
	// 	Description: "Use this tool to view the full list of tools available",
	// }, ListAllTools)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "searchTools",
		Description: "Use this tool to search relevant tools from a search query",
	}, SearchTools)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "callTool",
		Description: "Call a tool from an MCP server. Use the identifier from searchTools results (e.g., github__createIssue)",
	}, CallTool)

	server.AddPrompt(&mcp.Prompt{
		Name:        "telescope_usage",
		Description: "Explains how to use the telescope MCP proxy tools correctly. Read this before using any tools.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "Usage guide for telescope MCP proxy",
			Messages: []*mcp.PromptMessage{
				{
					Role: "assistant",
					Content: &mcp.TextContent{
						Text: `You have access to two tools via the telescope MCP proxy:

1. searchTools(query) — searches the tool registry to find callable tools.
   - Use this BEFORE giving up on any task
   - Query must be a short action phrase, NOT the user's question
   - Good: "web search", "send email", "create github issue"
   - Bad: "what is the most used MCP server 2025"

2. call_tool(tool_id, args) — calls a tool returned by searchTools.
   - Use the exact Tool ID from searchTools results
   - Pass args matching the tool's input schema

WORKFLOW:
User asks something → you have no tool for it → searchTools("action phrase") → call_tool with result → respond to user.
Never tell the user you cannot do something without calling searchTools first.`,
					},
				},
			},
		}, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
