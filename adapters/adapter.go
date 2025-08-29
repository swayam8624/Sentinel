package adapters

import (
	"context"
)

// LLMAdapter is the interface that all LLM adapters must implement
type LLMAdapter interface {
	// ChatCompletion sends a chat completion request
	ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

	// ChatCompletionStream sends a streaming chat completion request
	ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)

	// GetModelInfo gets information about a model
	GetModelInfo(ctx context.Context, modelID string) (*ModelInfo, error)

	// ValidateConfig validates the adapter configuration
	ValidateConfig() error

	// GetCapabilities returns the adapter capabilities
	GetCapabilities() *AdapterCapabilities
}

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents a chat completion response
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a response choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionStream is the interface for streaming responses
type ChatCompletionStream interface {
	// Recv receives the next chunk
	Recv() (*ChatCompletionStreamResponse, error)

	// Close closes the stream
	Close() error
}

// ChatCompletionStreamResponse represents a streaming response chunk
type ChatCompletionStreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice represents a streaming response choice
type StreamChoice struct {
	Index        int     `json:"index"`
	Delta        Message `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

// ModelInfo represents model information
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// AdapterCapabilities represents adapter capabilities
type AdapterCapabilities struct {
	Streaming     bool `json:"streaming"`
	FunctionCalls bool `json:"function_calls"`
	Embeddings    bool `json:"embeddings"`
	ModelInfo     bool `json:"model_info"`
	RateLimiting  bool `json:"rate_limiting"`
}

// AdapterFactory creates adapters
type AdapterFactory interface {
	// CreateAdapter creates a new adapter
	CreateAdapter(config *AdapterConfig) (LLMAdapter, error)
}

// AdapterConfig represents adapter configuration
type AdapterConfig struct {
	Type     string                 `json:"type"`
	Settings map[string]interface{} `json:"settings"`
}
