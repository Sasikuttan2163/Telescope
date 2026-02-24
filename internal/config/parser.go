package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type MainConfig struct {
	Qdrant     QdrantConfig `json:"qdrant" validate:"required"`
	Ollama     OllamaConfig `json:"ollama" validate:"required"`
	Stars      []StarConfig `json:"stars" validate:"required"`
	DebugLevel string       `json:"debug_level" validate:"required"`
}

type OllamaConfig struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Model    string `json:"model" validate:"required"`
	EmbedDim uint64 `json:"embed_dim" validate:"required"`
	NumGpu   int    `json:""num_gpu,omitempty"`
}

type StarConfig struct {
	ID   string `json:"id" validate:"required,uuid4"`
	Name string `json:"name" validate:"required"` // Human-readable name

	Transport TransportConfig `json:"transport" validate:"required"`

	Enabled     bool          `json:"enabled" validate:"boolean"`
	Timeout     time.Duration `json:"timeout_seconds" validate:"min=1"`
	MaxRetries  int           `json:"max_retries" validate:"min=0,max=5"`
	HealthCheck *HealthConfig `json:"health_check,omitempty"`

	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`
}

type TransportConfig struct {
	Type string `json:"type" validate:"oneof=http stdio sse"` // Transport type

	HTTP *HTTPTransportConfig `json:"http,omitempty" validate:"omitempty"`

	Stdio *StdioTransportConfig `json:"stdio,omitempty" validate:"omitempty"`

	SSE *SSETransportConfig `json:"sse,omitempty" validate:"omitempty"`
}

type HTTPTransportConfig struct {
	BaseURL string            `json:"base_url" validate:"required,url"` // Full MCP endpoint
	Headers map[string]string `json:"headers,omitempty"`                // Static headers
}

type StdioTransportConfig struct {
	Command    []string `json:"command" validate:"required,min=1"` // ["npx", "@mcp/server-filesystem"]
	Args       []string `json:"args,omitempty"`
	Env        []string `json:"env,omitempty"` // Environment variables for process
	WorkingDir string   `json:"working_dir,omitempty"`
}

type SSETransportConfig struct {
	URL string `json:"url" validate:"required,url"`
}

type HealthConfig struct {
	Endpoint string        `json:"endpoint,omitempty"`
	Interval time.Duration `json:"interval_seconds" validate:"min=5"`
}

type QdrantConfig struct {
	Host           string `json:"host" validate:"required,url"`
	Port           int    `json:"port" validate:"required,url"`
	CollectionName string `json:"collection_name" validate:"required"`
}

func GetConfig(configFile string) (MainConfig, error) {
	content, _ := os.ReadFile(configFile)
	var starConfig MainConfig
	err := json.Unmarshal(content, &starConfig)

	fmt.Println(starConfig)
	return starConfig, err
}

func main() {
	GetConfig("test.json")
}
