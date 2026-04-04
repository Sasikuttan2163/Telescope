package types

type Tool struct {
	Name         string
	Identifier   string
	Uuid         string
	Description  string
	EmbedString  string
	InputSchema  any
	OutputSchema any
	Vector       []float32
}
