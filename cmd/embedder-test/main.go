package main

import (
	"context"
	"fmt"

	"github.com/Sasikuttan2163/Telescope/internal/embed"
)

func main() {
	host := "localhost"
	port := 11434
	model := "bge-m3"
	mcpServerName := "example-mcp"
	toolNames := []string{"search", "summarize"}
	toolDescriptions := []string{
		"Search documents by keyword and return matching snippets.",
		"Summarize a document or list of snippets into a short overview.",
	}
	gpu := 0
	embeddings, err := embed.OllamaGetToolVector(context.Background(), host, port, &gpu, model, mcpServerName, toolNames, toolDescriptions)
	if err != nil {
		fmt.Printf("embed error: %v\n", err)
		return
	}

	for i, vec := range embeddings {
		fmt.Printf("tool=%s dims=%d first3=%v\n", toolNames[i], len(vec), vec[:min(3, len(vec))])
	}
}
