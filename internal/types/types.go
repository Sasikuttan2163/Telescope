package types

type Tool struct {
	ServerName   string
	Name         string
	Identifier   string
	Uuid         string
	Description  string
	EmbedString  string
	InputSchema  any
	OutputSchema any
	Vector       []float32
}
