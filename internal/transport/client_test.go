package transport

import (
	"net/http"
	"testing"

	"github.com/Sasikuttan2163/Telescope/internal/config"
)

func TestGetStarByName_Found(t *testing.T) {
	stars := []config.StarConfig{
		{
			Name: "server1",
			ID:   "id1",
		},
		{
			Name: "server2",
			ID:   "id2",
		},
	}

	result := GetStarByName(stars, "server1")
	if result == nil {
		t.Fatal("Expected to find star 'server1'")
	}
	if result.Name != "server1" {
		t.Errorf("Expected name 'server1', got %s", result.Name)
	}
	if result.ID != "id1" {
		t.Errorf("Expected ID 'id1', got %s", result.ID)
	}
}

func TestGetStarByName_NotFound(t *testing.T) {
	stars := []config.StarConfig{
		{
			Name: "server1",
			ID:   "id1",
		},
		{
			Name: "server2",
			ID:   "id2",
		},
	}

	result := GetStarByName(stars, "nonexistent")
	if result != nil {
		t.Error("Expected nil for non-existent star")
	}
}

func TestGetStarByName_EmptyList(t *testing.T) {
	stars := []config.StarConfig{}

	result := GetStarByName(stars, "server1")
	if result != nil {
		t.Error("Expected nil for empty list")
	}
}

func TestTransportConfig_HTTP(t *testing.T) {
	transport := config.TransportConfig{
		Type: "http",
		HTTP: &config.HTTPTransportConfig{
			BaseURL: "http://localhost:8080/mcp",
			Headers: map[string]string{
				"Authorization": "Bearer token",
			},
		},
	}

	if transport.Type != "http" {
		t.Errorf("Expected type 'http', got %s", transport.Type)
	}
	if transport.HTTP == nil {
		t.Fatal("Expected HTTP config to be not nil")
	}
	if transport.HTTP.BaseURL != "http://localhost:8080/mcp" {
		t.Errorf("Expected BaseURL 'http://localhost:8080/mcp', got %s", transport.HTTP.BaseURL)
	}
}

func TestTransportConfig_Stdio(t *testing.T) {
	transport := config.TransportConfig{
		Type: "stdio",
		Stdio: &config.StdioTransportConfig{
			Command:    []string{"npx", "@mcp/server-filesystem"},
			Args:       []string{"/path"},
			Env:        []string{"DEBUG=true"},
			WorkingDir: "/home/user",
		},
	}

	if transport.Type != "stdio" {
		t.Errorf("Expected type 'stdio', got %s", transport.Type)
	}
	if transport.Stdio == nil {
		t.Fatal("Expected Stdio config to be not nil")
	}
	if len(transport.Stdio.Command) != 2 {
		t.Errorf("Expected 2 command elements, got %d", len(transport.Stdio.Command))
	}
}

func TestTransportConfig_SSE(t *testing.T) {
	transport := config.TransportConfig{
		Type: "sse",
		SSE: &config.SSETransportConfig{
			URL: "http://localhost:8080/sse",
		},
	}

	if transport.Type != "sse" {
		t.Errorf("Expected type 'sse', got %s", transport.Type)
	}
	if transport.SSE == nil {
		t.Fatal("Expected SSE config to be not nil")
	}
	if transport.SSE.URL != "http://localhost:8080/sse" {
		t.Errorf("Expected URL 'http://localhost:8080/sse', got %s", transport.SSE.URL)
	}
}

func TestHttpTripper_RoundTrip(t *testing.T) {
	tripper := &httpTripper{
		headers: map[string]string{
			"Authorization": "Bearer test-token",
			"X-Custom":      "custom-value",
		},
		mainTripper: http.DefaultTransport,
	}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tripper.RoundTrip(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.Header.Get("Authorization") != "Bearer test-token" {
		t.Errorf("Expected Authorization header to be set")
	}
	if req.Header.Get("X-Custom") != "custom-value" {
		t.Errorf("Expected X-Custom header to be set")
	}
}
