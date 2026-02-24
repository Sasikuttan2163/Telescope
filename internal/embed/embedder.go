package embed

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Sasikuttan2163/Telescope/internal/types"
	"github.com/ollama/ollama/api"
)

func OllamaGetToolVector(ctx context.Context, host string, port int, numGpu *int, model string, mcpServerName string, tools *[]*types.Tool) ([][]float32, error) {
	embedTexts := make([]string, len(*tools))
	for i, tool := range *tools {
		embedStr := fmt.Sprintf("Provider: %s\nTool: %s\nDescription: %s", mcpServerName, tool.Name, tool.Description)
		tool.EmbedString = embedStr
		embedTexts[i] = embedStr
	}

	baseURL := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
	}

	client := api.NewClient(baseURL, http.DefaultClient)

	req := &api.EmbedRequest{
		Model: model,
		Input: embedTexts,
	}

	if numGpu != nil {
		req.Options = map[string]interface{}{
			"num_gpu": numGpu,
		}
	}

	resp, err := client.Embed(ctx, req)
	if err != nil {
		return nil, err
	}

	for i, tool := range *tools {
		if i < len(resp.Embeddings) {
			tool.Vector = resp.Embeddings[i]
		}
	}

	return resp.Embeddings, nil
}
