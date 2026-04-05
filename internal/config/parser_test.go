package config

import (
	"os"
	"testing"
)

func TestGetConfig_ValidConfig(t *testing.T) {
	content := `{
		"qdrant": {
			"host": "http://localhost",
			"port": 6333,
			"collection_name": "test_collection"
		},
		"ollama": {
			"host": "localhost",
			"port": 11434,
			"model": "nomic-embed-text",
			"embed_dim": 768,
			"num_gpu": 1
		},
		"stars": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"name": "test-star",
				"transport": {
					"type": "http",
					"http": {
						"base_url": "http://localhost:8080/mcp"
					}
				},
				"enabled": true,
				"timeout_seconds": 30,
				"max_retries": 3
			}
		],
		"debug_level": "info"
	}`

	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	config, err := GetConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Qdrant.Host != "http://localhost" {
		t.Errorf("Expected Qdrant.Host to be 'http://localhost', got %s", config.Qdrant.Host)
	}
	if config.Qdrant.Port != 6333 {
		t.Errorf("Expected Qdrant.Port to be 6333, got %d", config.Qdrant.Port)
	}
	if config.Qdrant.CollectionName != "test_collection" {
		t.Errorf("Expected Qdrant.CollectionName to be 'test_collection', got %s", config.Qdrant.CollectionName)
	}

	if config.Ollama.Host != "localhost" {
		t.Errorf("Expected Ollama.Host to be 'localhost', got %s", config.Ollama.Host)
	}
	if config.Ollama.Port != 11434 {
		t.Errorf("Expected Ollama.Port to be 11434, got %d", config.Ollama.Port)
	}
	if config.Ollama.Model != "nomic-embed-text" {
		t.Errorf("Expected Ollama.Model to be 'nomic-embed-text', got %s", config.Ollama.Model)
	}
	if config.Ollama.EmbedDim != 768 {
		t.Errorf("Expected Ollama.EmbedDim to be 768, got %d", config.Ollama.EmbedDim)
	}

	if len(config.Stars) != 1 {
		t.Fatalf("Expected 1 star, got %d", len(config.Stars))
	}
	star := config.Stars[0]
	if star.Name != "test-star" {
		t.Errorf("Expected star name to be 'test-star', got %s", star.Name)
	}
	if !star.Enabled {
		t.Errorf("Expected star to be enabled")
	}
}

func TestGetConfig_InvalidJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("invalid json"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	_, err = GetConfig(tmpFile.Name())
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestGetConfig_FileNotFound(t *testing.T) {
	_, err := GetConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("Expected error for non-existent file")
	}
}
