package qdrantdb

import (
	"testing"

	"github.com/Sasikuttan2163/Telescope/internal/types"
)

func TestStructToMap(t *testing.T) {
	tool := &types.Tool{
		ServerName:   "test-server",
		Name:         "test-tool",
		Identifier:   "test-server::test-tool",
		Uuid:         "test-uuid",
		Description:  "A test tool",
		InputSchema:  map[string]interface{}{"type": "object"},
		OutputSchema: map[string]interface{}{"type": "object"},
		Vector:       []float32{1.0, 2.0, 3.0},
	}

	result := structToMap(tool, "Vector")

	if result["ServerName"] != "test-server" {
		t.Errorf("Expected ServerName 'test-server', got %v", result["ServerName"])
	}
	if result["Name"] != "test-tool" {
		t.Errorf("Expected Name 'test-tool', got %v", result["Name"])
	}
	if result["Identifier"] != "test-server::test-tool" {
		t.Errorf("Expected Identifier 'test-server::test-tool', got %v", result["Identifier"])
	}
	if result["Uuid"] != "test-uuid" {
		t.Errorf("Expected Uuid 'test-uuid', got %v", result["Uuid"])
	}
	if result["Description"] != "A test tool" {
		t.Errorf("Expected Description 'A test tool', got %v", result["Description"])
	}

	if _, exists := result["Vector"]; exists {
		t.Error("Vector should be excluded")
	}
}

func TestStructToMap_NoExclusions(t *testing.T) {
	tool := &types.Tool{
		Name:        "test-tool",
		ServerName:  "test-server",
		Description: "Test description",
		Vector:      []float32{1.0, 2.0},
	}

	result := structToMap(tool)

	if result["Name"] != "test-tool" {
		t.Errorf("Expected Name 'test-tool', got %v", result["Name"])
	}
}

func TestStructToMap_EmptyStruct(t *testing.T) {
	tool := &types.Tool{}

	result := structToMap(tool)

	if len(result) == 0 {
		t.Error("Expected non-empty map for empty struct")
	}
}
