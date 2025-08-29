package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIAdapter implements the LLMAdapter interface for OpenAI
type OpenAIAdapter struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
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

// ChatCompletionStream represents a streaming chat completion
type ChatCompletionStream struct {
	response *http.Response
	reader   *StreamReader
}

// StreamReader reads streaming responses
type StreamReader struct {
	reader io.Reader
	buffer []byte
}

// NewOpenAIAdapter creates a new OpenAI adapter
func NewOpenAIAdapter(apiKey, baseURL string, timeout time.Duration) *OpenAIAdapter {
	return &OpenAIAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// ChatCompletion sends a chat completion request to OpenAI
func (oa *OpenAIAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Convert request to JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", oa.baseURL+"/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+oa.apiKey)

	// Send request
	httpResp, err := oa.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// Check status code
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned status %d", httpResp.StatusCode)
	}

	// Parse response
	var resp ChatCompletionResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// ChatCompletionStream sends a streaming chat completion request to OpenAI
func (oa *OpenAIAdapter) ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionStream, error) {
	// Set streaming flag
	req.Stream = true

	// Convert request to JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", oa.baseURL+"/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+oa.apiKey)

	// Send request
	httpResp, err := oa.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check status code
	if httpResp.StatusCode != http.StatusOK {
		httpResp.Body.Close()
		return nil, fmt.Errorf("OpenAI API returned status %d", httpResp.StatusCode)
	}

	// Create stream
	stream := &ChatCompletionStream{
		response: httpResp,
		reader: &StreamReader{
			reader: httpResp.Body,
			buffer: make([]byte, 0, 4096),
		},
	}

	return stream, nil
}

// Recv receives the next chunk from the stream
func (s *ChatCompletionStream) Recv() (*ChatCompletionStreamResponse, error) {
	// This is a simplified implementation
	// In a real implementation, we'd parse the SSE format
	return nil, fmt.Errorf("streaming not fully implemented")
}

// Close closes the stream
func (s *ChatCompletionStream) Close() error {
	if s.response != nil && s.response.Body != nil {
		return s.response.Body.Close()
	}
	return nil
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

// GetModelInfo gets information about a model
func (oa *OpenAIAdapter) GetModelInfo(ctx context.Context, modelID string) (*ModelInfo, error) {
	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", oa.baseURL+"/models/"+modelID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+oa.apiKey)

	// Send request
	httpResp, err := oa.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// Check status code
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned status %d", httpResp.StatusCode)
	}

	// Parse response
	var resp ModelInfo
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// ModelInfo represents model information
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ValidateConfig validates the adapter configuration
func (oa *OpenAIAdapter) ValidateConfig() error {
	if oa.apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	if oa.baseURL == "" {
		oa.baseURL = "https://api.openai.com/v1"
	}

	return nil
}

// GetCapabilities returns the adapter capabilities
func (oa *OpenAIAdapter) GetCapabilities() *AdapterCapabilities {
	return &AdapterCapabilities{
		Streaming:     true,
		FunctionCalls: true,
		Embeddings:    true,
		ModelInfo:     true,
		RateLimiting:  true,
	}
}

// AdapterCapabilities represents adapter capabilities
type AdapterCapabilities struct {
	Streaming     bool `json:"streaming"`
	FunctionCalls bool `json:"function_calls"`
	Embeddings    bool `json:"embeddings"`
	ModelInfo     bool `json:"model_info"`
	RateLimiting  bool `json:"rate_limiting"`
}
