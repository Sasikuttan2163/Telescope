package types

import (
	"testing"
)

func TestTool_Structure(t *testing.T) {
	tool := Tool{
		ServerName:   "test-server",
		Name:         "test-tool",
		Identifier:   "test-server::test-tool",
		Uuid:         "test-uuid",
		Description:  "A test tool",
		InputSchema:  map[string]interface{}{"type": "object"},
		OutputSchema: map[string]interface{}{"type": "object"},
		Vector:       []float32{1.0, 2.0, 3.0},
	}

	if tool.ServerName != "test-server" {
		t.Errorf("Expected ServerName 'test-server', got %s", tool.ServerName)
	}
	if tool.Name != "test-tool" {
		t.Errorf("Expected Name 'test-tool', got %s", tool.Name)
	}
	if tool.Identifier != "test-server::test-tool" {
		t.Errorf("Expected Identifier 'test-server::test-tool', got %s", tool.Identifier)
	}
	if tool.Uuid != "test-uuid" {
		t.Errorf("Expected Uuid 'test-uuid', got %s", tool.Uuid)
	}
	if tool.Description != "A test tool" {
		t.Errorf("Expected Description 'A test tool', got %s", tool.Description)
	}
	if len(tool.Vector) != 3 {
		t.Errorf("Expected Vector length 3, got %d", len(tool.Vector))
	}
}

func TestTool_VectorEmpty(t *testing.T) {
	tool := Tool{
		Name:   "empty-tool",
		Vector: []float32{},
	}

	if len(tool.Vector) != 0 {
		t.Errorf("Expected empty Vector, got length %d", len(tool.Vector))
	}
}
