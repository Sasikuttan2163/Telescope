package indexer

import (
	"testing"

	"github.com/qdrant/go-client/qdrant"
)

func TestQdrantPayloadToMap(t *testing.T) {
	payload := map[string]*qdrant.Value{
		"ServerName": {
			Kind: &qdrant.Value_StringValue{StringValue: "test-server"},
		},
		"Uuid": {
			Kind: &qdrant.Value_StringValue{StringValue: "test-uuid"},
		},
		"Count": {
			Kind: &qdrant.Value_IntegerValue{IntegerValue: 42},
		},
		"Rating": {
			Kind: &qdrant.Value_DoubleValue{DoubleValue: 3.14},
		},
		"Active": {
			Kind: &qdrant.Value_BoolValue{BoolValue: true},
		},
	}

	result := qdrantPayloadToMap(payload)

	if result["ServerName"] != "test-server" {
		t.Errorf("Expected ServerName 'test-server', got %v", result["ServerName"])
	}
	if result["Uuid"] != "test-uuid" {
		t.Errorf("Expected Uuid 'test-uuid', got %v", result["Uuid"])
	}
	if result["Count"].(int64) != 42 {
		t.Errorf("Expected Count 42, got %v", result["Count"])
	}
	if result["Rating"].(float64) != 3.14 {
		t.Errorf("Expected Rating 3.14, got %v", result["Rating"])
	}
	if result["Active"] != true {
		t.Errorf("Expected Active true, got %v", result["Active"])
	}
}

func TestQdrantValueToAny_String(t *testing.T) {
	value := &qdrant.Value{
		Kind: &qdrant.Value_StringValue{StringValue: "test-string"},
	}

	result := qdrantValueToAny(value)
	if result != "test-string" {
		t.Errorf("Expected 'test-string', got %v", result)
	}
}

func TestQdrantValueToAny_Integer(t *testing.T) {
	value := &qdrant.Value{
		Kind: &qdrant.Value_IntegerValue{IntegerValue: 123},
	}

	result := qdrantValueToAny(value)
	if result != int64(123) {
		t.Errorf("Expected 123, got %v", result)
	}
}

func TestQdrantValueToAny_Double(t *testing.T) {
	value := &qdrant.Value{
		Kind: &qdrant.Value_DoubleValue{DoubleValue: 1.23},
	}

	result := qdrantValueToAny(value)
	if result != 1.23 {
		t.Errorf("Expected 1.23, got %v", result)
	}
}

func TestQdrantValueToAny_Bool(t *testing.T) {
	value := &qdrant.Value{
		Kind: &qdrant.Value_BoolValue{BoolValue: false},
	}

	result := qdrantValueToAny(value)
	if result != false {
		t.Errorf("Expected false, got %v", result)
	}
}

func TestQdrantValueToAny_List(t *testing.T) {
	value := &qdrant.Value{
		Kind: &qdrant.Value_ListValue{
			ListValue: &qdrant.ListValue{
				Values: []*qdrant.Value{
					{Kind: &qdrant.Value_StringValue{StringValue: "a"}},
					{Kind: &qdrant.Value_StringValue{StringValue: "b"}},
					{Kind: &qdrant.Value_StringValue{StringValue: "c"}},
				},
			},
		},
	}

	result := qdrantValueToAny(value)
	list, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any type")
	}
	if len(list) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(list))
	}
	if list[0] != "a" {
		t.Errorf("Expected first element 'a', got %v", list[0])
	}
}

func TestQdrantValueToAny_Struct(t *testing.T) {
	value := &qdrant.Value{
		Kind: &qdrant.Value_StructValue{
			StructValue: &qdrant.Struct{
				Fields: map[string]*qdrant.Value{
					"name":  {Kind: &qdrant.Value_StringValue{StringValue: "test"}},
					"count": {Kind: &qdrant.Value_IntegerValue{IntegerValue: 5}},
				},
			},
		},
	}

	result := qdrantValueToAny(value)
	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any type")
	}
	if m["name"] != "test" {
		t.Errorf("Expected name 'test', got %v", m["name"])
	}
	if m["count"] != int64(5) {
		t.Errorf("Expected count 5, got %v", m["count"])
	}
}

func TestGetPayload_RetrievedPoint(t *testing.T) {
	point := &qdrant.RetrievedPoint{
		Payload: map[string]*qdrant.Value{
			"Name": {Kind: &qdrant.Value_StringValue{StringValue: "test-point"}},
		},
	}

	result := getPayload(point)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result["Name"] != "test-point" {
		t.Errorf("Expected Name 'test-point', got %v", result["Name"])
	}
}

func TestGetPayload_ScoredPoint(t *testing.T) {
	point := &qdrant.ScoredPoint{
		Payload: map[string]*qdrant.Value{
			"Name": {Kind: &qdrant.Value_StringValue{StringValue: "scored-point"}},
		},
	}

	result := getPayload(point)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result["Name"] != "scored-point" {
		t.Errorf("Expected Name 'scored-point', got %v", result["Name"])
	}
}

func TestGetPayload_UnknownType(t *testing.T) {
	result := getPayload("unknown-type")
	if result != nil {
		t.Error("Expected nil for unknown type")
	}
}
