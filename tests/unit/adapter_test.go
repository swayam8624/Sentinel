package unit

import (
	"testing"

	"github.com/sentinel-platform/sentinel/adapters/openai"
)

// TestChatCompletionRequest tests the ChatCompletionRequest structure
func TestChatCompletionRequest(t *testing.T) {
	req := &openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
			{
				Role:    "user",
				Content: "Hello, world!",
			},
		},
		Temperature: 0.7,
		MaxTokens:   100,
		Stream:      false,
	}

	if req.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected model gpt-3.5-turbo, got %s", req.Model)
	}

	if len(req.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(req.Messages))
	}

	if req.Messages[0].Content != "Hello, world!" {
		t.Errorf("Expected message content 'Hello, world!', got %s", req.Messages[0].Content)
	}
}

// TestChatCompletionResponse tests the ChatCompletionResponse structure
func TestChatCompletionResponse(t *testing.T) {
	resp := &openai.ChatCompletionResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   "gpt-3.5-turbo",
		Choices: []openai.Choice{
			{
				Index: 0,
				Message: openai.Message{
					Role:    "assistant",
					Content: "Hello! How can I help you today?",
				},
				FinishReason: "stop",
			},
		},
		Usage: openai.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}

	if resp.ID != "test-id" {
		t.Errorf("Expected ID test-id, got %s", resp.ID)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Hello! How can I help you today?" {
		t.Errorf("Expected response content 'Hello! How can I help you today?', got %s", resp.Choices[0].Message.Content)
	}

	if resp.Usage.TotalTokens != 30 {
		t.Errorf("Expected total tokens 30, got %d", resp.Usage.TotalTokens)
	}
}
