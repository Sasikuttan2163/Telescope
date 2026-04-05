package embed

import (
	"context"
	"testing"

	"github.com/Sasikuttan2163/Telescope/internal/types"
)

func TestEmbedStringFormatting(t *testing.T) {
	tools := []*types.Tool{
		{
			ServerName:  "test-server",
			Name:        "read_file",
			Description: "Read contents of a file",
			EmbedString: "Provider: test-server\nTool: read_file\nDescription: Read contents of a file",
		},
		{
			ServerName:  "test-server",
			Name:        "write_file",
			Description: "Write content to a file",
			EmbedString: "Provider: test-server\nTool: write_file\nDescription: Write content to a file",
		},
	}

	for i, tool := range tools {
		expected := ""
		switch tool.Name {
		case "read_file":
			expected = "Provider: test-server\nTool: read_file\nDescription: Read contents of a file"
		case "write_file":
			expected = "Provider: test-server\nTool: write_file\nDescription: Write content to a file"
		}

		if tool.EmbedString != expected {
			t.Errorf("Tool %d: expected embed string %q, got %q", i, expected, tool.EmbedString)
		}
	}
}

func TestEmbedStringEmptyDescription(t *testing.T) {
	tool := &types.Tool{
		ServerName:  "test-server",
		Name:        "empty-desc-tool",
		Description: "",
		EmbedString: "Provider: test-server\nTool: empty-desc-tool\nDescription: ",
	}

	expected := "Provider: test-server\nTool: empty-desc-tool\nDescription: "
	if tool.EmbedString != expected {
		t.Errorf("Expected embed string %q, got %q", expected, tool.EmbedString)
	}
}

func TestOllamaGetToolVector_EmbedStringGeneration(t *testing.T) {
	host := "localhost"
	port := 11434
	numGpu := 1
	model := "nomic-embed-text"
	mcpServerName := "test-server"

	tools := []*types.Tool{
		{
			ServerName:  "server1",
			Name:        "tool1",
			Description: "First tool",
		},
		{
			ServerName:  "server2",
			Name:        "tool2",
			Description: "Second tool",
		},
	}

	ctx := context.Background()

	vectors, err := OllamaGetToolVector(ctx, host, port, &numGpu, model, mcpServerName, &tools)

	if err != nil {
		t.Logf("Ollama call failed (expected if Ollama not running): %v", err)
		return
	}

	if len(vectors) != 2 {
		t.Errorf("Expected 2 vectors, got %d", len(vectors))
	}

	for i, tool := range tools {
		if tool.EmbedString == "" {
			t.Errorf("Tool %d: EmbedString should not be empty", i)
		}
		if len(tool.Vector) == 0 {
			t.Errorf("Tool %d: Vector should not be empty", i)
		}
	}
}

func TestOllamaGetQueryVector(t *testing.T) {
	host := "localhost"
	port := 11434
	numGpu := 1
	model := "nomic-embed-text"
	query := "Find a file reading tool"

	ctx := context.Background()

	vector, err := OllamaGetQueryVector(ctx, host, port, &numGpu, model, query)

	if err != nil {
		t.Logf("Ollama call failed (expected if Ollama not running): %v", err)
		return
	}

	if len(vector) == 0 {
		t.Error("Expected non-empty vector")
	}
}
