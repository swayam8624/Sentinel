package integration

import (
	"testing"
	"time"

	"github.com/sentinel-platform/sentinel/adapters"
	"github.com/sentinel-platform/sentinel/adapters/openai"
)

// TestAdapterCreation tests creating different LLM adapters
func TestAdapterCreation(t *testing.T) {
	// Test OpenAI adapter creation
	oa := openai.NewOpenAIAdapter("test-key", "https://api.openai.com/v1", 30*time.Second)

	if oa == nil {
		t.Error("Failed to create OpenAI adapter")
	}

	// Test capabilities
	capabilities := oa.GetCapabilities()
	if capabilities == nil {
		t.Error("OpenAI adapter should return capabilities")
	}

	// Verify all capabilities are true (as per implementation)
	if !capabilities.Streaming {
		t.Error("OpenAI adapter should support streaming")
	}

	if !capabilities.FunctionCalls {
		t.Error("OpenAI adapter should support function calls")
	}

	if !capabilities.Embeddings {
		t.Error("OpenAI adapter should support embeddings")
	}

	if !capabilities.ModelInfo {
		t.Error("OpenAI adapter should support model info")
	}

	if !capabilities.RateLimiting {
		t.Error("OpenAI adapter should support rate limiting")
	}
}

// TestAdapterConfigValidation tests adapter configuration validation
func TestAdapterConfigValidation(t *testing.T) {
	// Test OpenAI adapter with valid config
	oa := openai.NewOpenAIAdapter("test-key", "https://api.openai.com/v1", 30*time.Second)

	err := oa.ValidateConfig()
	if err != nil {
		t.Errorf("OpenAI adapter config validation failed: %v", err)
	}

	// Test OpenAI adapter with empty API key (should fail)
	oaEmptyKey := openai.NewOpenAIAdapter("", "https://api.openai.com/v1", 30*time.Second)

	err = oaEmptyKey.ValidateConfig()
	if err == nil {
		t.Error("OpenAI adapter config validation should fail with empty API key")
	}
}

// TestChatCompletionRequestStructure tests the chat completion request structure
func TestChatCompletionRequestStructure(t *testing.T) {
	req := &adapters.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []adapters.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: "Hello, how are you?",
			},
		},
		Temperature: 0.7,
		MaxTokens:   150,
		Stream:      false,
	}

	// Validate request structure
	if req.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected model gpt-3.5-turbo, got %s", req.Model)
	}

	if len(req.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(req.Messages))
	}

	if req.Messages[0].Role != "system" {
		t.Errorf("Expected first message role 'system', got %s", req.Messages[0].Role)
	}

	if req.Messages[1].Role != "user" {
		t.Errorf("Expected second message role 'user', got %s", req.Messages[1].Role)
	}

	if req.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", req.Temperature)
	}

	if req.MaxTokens != 150 {
		t.Errorf("Expected max tokens 150, got %d", req.MaxTokens)
	}

	if req.Stream != false {
		t.Errorf("Expected stream false, got %t", req.Stream)
	}
}

// TestChatCompletionResponseStructure tests the chat completion response structure
func TestChatCompletionResponseStructure(t *testing.T) {
	resp := &adapters.ChatCompletionResponse{
		ID:      "chatcmpl-1234567890",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-3.5-turbo",
		Choices: []adapters.Choice{
			{
				Index: 0,
				Message: adapters.Message{
					Role:    "assistant",
					Content: "I'm doing well, thank you for asking! How can I assist you today?",
				},
				FinishReason: "stop",
			},
		},
		Usage: adapters.Usage{
			PromptTokens:     15,
			CompletionTokens: 25,
			TotalTokens:      40,
		},
	}

	// Validate response structure
	if resp.ID == "" {
		t.Error("Response ID should not be empty")
	}

	if resp.Object != "chat.completion" {
		t.Errorf("Expected object 'chat.completion', got %s", resp.Object)
	}

	if resp.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got %s", resp.Model)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Role != "assistant" {
		t.Errorf("Expected assistant role, got %s", resp.Choices[0].Message.Role)
	}

	if resp.Usage.TotalTokens != 40 {
		t.Errorf("Expected total tokens 40, got %d", resp.Usage.TotalTokens)
	}
}
