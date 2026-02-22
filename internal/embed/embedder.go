package embed

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

func getEmbedString(mcpServerName string, toolName string, toolDescription string) string {
	return fmt.Sprintf("Tool: %s\nProvider: %s\nDescription: %s", toolName, mcpServerName, toolDescription)
}

func ollamaGetToolVector(host string, port int, model string, mcpServerName string, toolName string, toolDescription string) ([][]float32, error) {
	embed_text := getEmbedString(mcpServerName, toolName, toolDescription)

	baseURL := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
	}

	client := api.NewClient(baseURL, http.DefaultClient)

	req := &api.EmbedRequest{
		Model: model,
		Input: embed_text,
	}

	resp, err := client.Embed(context.Background(), req)

	if err != nil {
		return nil, err
	}

	return resp.Embeddings, err
}

func OllamaGetMCPVectors(host string, port int, model string, mcpServerName string, tools map[string]string) {

}
