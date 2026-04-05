package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/types"
	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type httpTripper struct {
	mainTripper http.RoundTripper
	headers     map[string]string
}

func (ht *httpTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range ht.headers {
		req.Header.Set(k, v)
	}
	return ht.mainTripper.RoundTrip(req)
}

func FetchToolsOfStar(ctx context.Context, star config.StarConfig) ([]*types.Tool, error) {
	httpClient := &http.Client{
		Transport: &httpTripper{
			headers:     star.Transport.HTTP.Headers,
			mainTripper: http.DefaultTransport,
		},
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "telescope-client",
		Version: "v0.0.1",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint:   star.Transport.HTTP.BaseURL,
		HTTPClient: httpClient,
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", star.ID, err)
		return nil, fmt.Errorf("Failed to connect to MCP %s: %s", star.ID, star.Name)
	}
	defer session.Close()

	if session.InitializeResult().Capabilities.Tools != nil {
		res, err := session.ListTools(ctx, nil)
		if err != nil {
			log.Printf("Failed to list tools for %s: %v", star.ID, err)
			return nil, err
		}
		fmt.Println("Found", len(res.Tools), "tools in", star.Name)

		tools := make([]*types.Tool, len(res.Tools))
		for i, mcpTool := range res.Tools {
			ident := fmt.Sprintf("%s::%s", star.Name, mcpTool.Name)
			tools[i] = &types.Tool{
				ServerName:   star.Name,
				Name:         mcpTool.Name,
				Description:  mcpTool.Description,
				Identifier:   ident,
				Uuid:         uuid.NewSHA1(uuid.NameSpaceURL, []byte(ident)).String(),
				InputSchema:  mcpTool.InputSchema,
				OutputSchema: mcpTool.OutputSchema,
			}
		}
		return tools, err
	}

	return nil, fmt.Errorf("MCP %s: %s does not support tools", star.ID, star.Name)
}

func CallToolOnStar(ctx context.Context, star config.StarConfig, toolName string, input json.RawMessage) (*mcp.CallToolResult, error) {
	httpClient := &http.Client{
		Transport: &httpTripper{
			headers:     star.Transport.HTTP.Headers,
			mainTripper: http.DefaultTransport,
		},
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "telescope-client",
		Version: "v0.0.1",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint:   star.Transport.HTTP.BaseURL,
		HTTPClient: httpClient,
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MCP %s: %w", star.ID, err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: input,
	}

	result, err := session.CallTool(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s on %s: %w", toolName, star.ID, err)
	}

	return result, nil
}

func GetStarByName(stars []config.StarConfig, serverName string) *config.StarConfig {
	for _, star := range stars {
		if star.Name == serverName {
			return &star
		}
	}
	return nil
}
