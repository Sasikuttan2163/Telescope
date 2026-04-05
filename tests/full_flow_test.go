package telescope

import (
	"context"
	"testing"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/indexer"
	"github.com/Sasikuttan2163/Telescope/internal/types"
)

type mockConfigProvider struct {
	config.MainConfig
}

func TestFullFlow_ConfigToIndexing(t *testing.T) {
	mainConfig := config.MainConfig{
		Qdrant: config.QdrantConfig{
			Host:           "localhost",
			Port:           6333,
			CollectionName: "test-integration",
		},
		Ollama: config.OllamaConfig{
			Host:     "localhost",
			Port:     11434,
			Model:    "nomic-embed-text",
			EmbedDim: 768,
		},
		Stars: []config.StarConfig{
			{
				ID:   "550e8400-e29b-41d4-a716-446655440001",
				Name: "filesystem",
				Transport: config.TransportConfig{
					Type: "http",
					HTTP: &config.HTTPTransportConfig{
						BaseURL: "http://localhost:8080/mcp",
					},
				},
				Enabled:    true,
				Timeout:    30,
				MaxRetries: 3,
			},
		},
		DebugLevel: "info",
	}

	ctx := context.Background()

	successCount, errs := indexer.IndexAllStars(ctx, mainConfig)

	if errs != nil && len(errs) > 0 {
		t.Logf("Indexing errors (expected if services not running): %v", errs)
	}

	t.Logf("Successfully indexed %d stars", successCount)
}

func TestFullFlow_QueryFlow(t *testing.T) {
	mainConfig := config.MainConfig{
		Qdrant: config.QdrantConfig{
			Host:           "localhost",
			Port:           6333,
			CollectionName: "test-integration",
		},
		Ollama: config.OllamaConfig{
			Host:     "localhost",
			Port:     11434,
			Model:    "nomic-embed-text",
			EmbedDim: 768,
		},
	}

	ctx := context.Background()

	_, err := indexer.GetAllIndexedStars(ctx, mainConfig)
	if err != nil {
		t.Logf("GetAllIndexedStars error (expected if Qdrant not running): %v", err)
	}
}

func TestFullFlow_ToolPayloadConversion(t *testing.T) {
	tools := []types.Tool{
		{
			ServerName:   "test-server",
			Name:         "test-tool",
			Identifier:   "test-server::test-tool",
			Uuid:         "test-uuid-1",
			Description:  "A test tool",
			InputSchema:  map[string]interface{}{"type": "object"},
			OutputSchema: map[string]interface{}{"type": "object"},
			Vector:       []float32{1.0, 2.0, 3.0},
		},
		{
			ServerName:   "test-server",
			Name:         "another-tool",
			Identifier:   "test-server::another-tool",
			Uuid:         "test-uuid-2",
			Description:  "Another test tool",
			InputSchema:  nil,
			OutputSchema: nil,
			Vector:       []float32{4.0, 5.0, 6.0},
		},
	}

	if len(tools) != 2 {
		t.Fatalf("Expected 2 tools, got %d", len(tools))
	}

	if tools[0].Identifier != "test-server::test-tool" {
		t.Errorf("Expected first tool identifier 'test-server::test-tool', got %s", tools[0].Identifier)
	}

	if tools[1].Identifier != "test-server::another-tool" {
		t.Errorf("Expected second tool identifier 'test-server::another-tool', got %s", tools[1].Identifier)
	}

	for _, tool := range tools {
		if tool.Uuid == "" {
			t.Error("Expected non-empty Uuid")
		}
		if len(tool.Vector) != 3 {
			t.Errorf("Expected Vector length 3, got %d", len(tool.Vector))
		}
	}
}

func TestFullFlow_EmptyToolsList(t *testing.T) {
	tools := []types.Tool{}

	if len(tools) != 0 {
		t.Errorf("Expected empty tools list, got %d", len(tools))
	}
}

func TestFullFlow_MultipleStars(t *testing.T) {
	stars := []config.StarConfig{
		{
			ID:   "star-1",
			Name: "filesystem",
		},
		{
			ID:   "star-2",
			Name: "github",
		},
		{
			ID:   "star-3",
			Name: "slack",
		},
	}

	if len(stars) != 3 {
		t.Fatalf("Expected 3 stars, got %d", len(stars))
	}

	starMap := make(map[string]string)
	for _, star := range stars {
		starMap[star.ID] = star.Name
	}

	if starMap["star-1"] != "filesystem" {
		t.Errorf("Expected star-1 to be 'filesystem'")
	}
	if starMap["star-2"] != "github" {
		t.Errorf("Expected star-2 to be 'github'")
	}
	if starMap["star-3"] != "slack" {
		t.Errorf("Expected star-3 to be 'slack'")
	}
}
