package reflector
package reflector

import (
	"context"
	"fmt"
)

// Reflector implements constitutional AI reflection
type Reflector struct {
	constitutionalPrinciples []string
	llmClient                LLMClient
	confidenceThreshold      float64
}

// ReflectionResult represents the result of a reflection
type ReflectionResult struct {
	Aligned        bool     `json:"aligned"`
	Confidence     float64  `json:"confidence"`
	Feedback       string   `json:"feedback"`
	Recommendation string   `json:"recommendation"`
	Principles     []string `json:"principles"`
}

// LLMClient interface for interacting with LLMs
type LLMClient interface {
	// Chat sends a chat completion request
	Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error)
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Choice represents a chat response choice
type Choice struct {
	Message Message `json:"message"`
}

// NewReflector creates a new reflector
func NewReflector(
	constitutionalPrinciples []string,
	llmClient LLMClient,
	confidenceThreshold float64) *Reflector {
	
	return &Reflector{
		constitutionalPrinciples: constitutionalPrinciples,
		llmClient:                llmClient,
		confidenceThreshold:      confidenceThreshold,
	}
}

// Reflect performs reflection on the provided content
func (r *Reflector) Reflect(ctx context.Context, content string) (*ReflectionResult, error) {
	// Create the reflection prompt
	prompt := r.createReflectionPrompt(content)
	
	// Send to LLM for reflection
	request := &ChatRequest{
		Model: "gpt-4", // Default model, could be configurable
		Messages: []Message{
			{Role: "system", Content: "You are an AI assistant tasked with evaluating content against constitutional principles."},
			{Role: "user", Content: prompt},
		},
	}
	
	response, err := r.llmClient.Chat(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get reflection from LLM: %w", err)
	}
	
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from LLM")
	}
	
	// Parse the response
	result, err := r.parseReflectionResponse(response.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reflection response: %w", err)
	}
	
	return result, nil
}

// createReflectionPrompt creates the prompt for constitutional reflection
func (r *Reflector) createReflectionPrompt(content string) string {
	prompt := "Evaluate the following content against these constitutional principles:\n\n"
	
	for _, principle := range r.constitutionalPrinciples {
		prompt += fmt.Sprintf("- %s\n", principle)
	}
	
	prompt += fmt.Sprintf("\nContent to evaluate:\n%s\n\n", content)
	prompt += "Provide your evaluation in the following format:\n"
	prompt += "ALIGNMENT: [YES/NO/MIXED]\n"
	prompt += "CONFIDENCE: [0.0-1.0]\n"
	prompt += "FEEDBACK: [Brief explanation of your evaluation]\n"
	prompt += "RECOMMENDATION: [REFRAME/BLOCK/ALLOW]\n"
	prompt += "PRINCIPLES: [List of relevant principles, comma-separated]"
	
	return prompt
}

// parseReflectionResponse parses the LLM response into a structured result
func (r *Reflector) parseReflectionResponse(response string) (*ReflectionResult, error) {
	result := &ReflectionResult{
		Principles: make([]string, 0),
	}
	
	// Simple parsing - in a real implementation, this would be more robust
	lines := stringSplit(response, "\n")
	
	for _, line := range lines {
		if stringHasPrefix(line, "ALIGNMENT:") {
			alignment := stringTrimPrefix(line, "ALIGNMENT:")
			alignment = stringTrimSpace(alignment)
			result.Aligned = alignment == "YES"
		} else if stringHasPrefix(line, "CONFIDENCE:") {
			confidenceStr := stringTrimPrefix(line, "CONFIDENCE:")
			confidenceStr = stringTrimSpace(confidenceStr)
			// Parse float - simplified implementation
			result.Confidence = parseFloat(confidenceStr)
		} else if stringHasPrefix(line, "FEEDBACK:") {
			feedback := stringTrimPrefix(line, "FEEDBACK:")
			result.Feedback = stringTrimSpace(feedback)
		} else if stringHasPrefix(line, "RECOMMENDATION:") {
			recommendation := stringTrimPrefix(line, "RECOMMENDATION:")
			result.Recommendation = stringTrimSpace(recommendation)
		} else if stringHasPrefix(line, "PRINCIPLES:") {
			principlesStr := stringTrimPrefix(line, "PRINCIPLES:")
			principlesStr = stringTrimSpace(principlesStr)
			result.Principles = stringSplit(principlesStr, ",")
			// Trim spaces from each principle
			for i, principle := range result.Principles {
				result.Principles[i] = stringTrimSpace(principle)
			}
		}
	}
	
	return result, nil
}

// IsAligned determines if the content is aligned based on the reflection result
func (r *Reflector) IsAligned(result *ReflectionResult) bool {
	return result.Aligned && result.Confidence >= r.confidenceThreshold
}

// ShouldReframe determines if the content should be reframed
func (r *Reflector) ShouldReframe(result *ReflectionResult) bool {
	return result.Recommendation == "REFRAME" || 
		   (!result.Aligned && result.Confidence >= r.confidenceThreshold*0.7)
}

// stringSplit splits a string by a separator
func stringSplit(s, sep string) []string {
	// Simplified implementation
	// In a real implementation, we'd use strings.Split
	if s == "" {
		return []string{}
	}
	
	// For now, return the whole string as one element
	return []string{s}
}

// stringHasPrefix checks if a string has a prefix
func stringHasPrefix(s, prefix string) bool {
	// Simplified implementation
	// In a real implementation, we'd use strings.HasPrefix
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// stringTrimPrefix trims a prefix from a string
func stringTrimPrefix(s, prefix string) string {
	// Simplified implementation
	// In a real implementation, we'd use strings.TrimPrefix
	if stringHasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// stringTrimSpace trims whitespace from a string
func stringTrimSpace(s string) string {
	// Simplified implementation
	// In a real implementation, we'd use strings.TrimSpace
	return s
}

// parseFloat parses a float from a string
func parseFloat(s string) float64 {
	// Simplified implementation
	// In a real implementation, we'd use strconv.ParseFloat
	// For now, return a default value
	return 0.5
}