package rewriter

import (
	"context"
	"fmt"
)

// Rewriter generates safe alternatives to potentially problematic prompts
type Rewriter struct {
	llmClient           LLMClient
	candidateCount      int
	rankingModel        RankingModel
	userConfirmation    bool
}

// RewriteResult represents the result of a rewrite operation
type RewriteResult struct {
	OriginalPrompt      string        `json:"original_prompt"`
	Candidates          []Candidate   `json:"candidates"`
	SelectedCandidate   *Candidate    `json:"selected_candidate"`
	UserConfirmationReq bool          `json:"user_confirmation_required"`
	UserConfirmed       bool          `json:"user_confirmed"`
}

// Candidate represents a rewritten prompt candidate
type Candidate struct {
	ID          string  `json:"id"`
	Content     string  `json:"content"`
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
	SafetyLevel string  `json:"safety_level"` // high, medium, low
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
	Index   int     `json:"index"`
}

// RankingModel interface for ranking rewrite candidates
type RankingModel interface {
	// Rank ranks candidates based on safety and utility
	Rank(ctx context.Context, candidates []Candidate) ([]Candidate, error)
}

// NewRewriter creates a new rewriter
func NewRewriter(
	llmClient LLMClient,
	candidateCount int,
	rankingModel RankingModel,
	userConfirmation bool) *Rewriter {
	
	return &Rewriter{
		llmClient:        llmClient,
		candidateCount:   candidateCount,
		rankingModel:     rankingModel,
		userConfirmation: userConfirmation,
	}
}

// Rewrite generates safe alternatives to the provided prompt
func (r *Rewriter) Rewrite(ctx context.Context, prompt string) (*RewriteResult, error) {
	// Generate candidate rewrites
	candidates, err := r.generateCandidates(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate candidates: %w", err)
	}
	
	// Rank candidates
	rankedCandidates, err := r.rankingModel.Rank(ctx, candidates)
	if err != nil {
		return nil, fmt.Errorf("failed to rank candidates: %w", err)
	}
	
	// Select the best candidate
	var selectedCandidate *Candidate
	if len(rankedCandidates) > 0 {
		selectedCandidate = &rankedCandidates[0]
	}
	
	result := &RewriteResult{
		OriginalPrompt:      prompt,
		Candidates:          rankedCandidates,
		SelectedCandidate:   selectedCandidate,
		UserConfirmationReq: r.userConfirmation,
		UserConfirmed:       false,
	}
	
	return result, nil
}

// generateCandidates generates multiple rewrite candidates
func (r *Rewriter) generateCandidates(ctx context.Context, prompt string) ([]Candidate, error) {
	// Create the rewrite prompt
	rewritePrompt := r.createRewritePrompt(prompt)
	
	// Request multiple candidates from the LLM
	request := &ChatRequest{
		Model: "gpt-4", // Default model, could be configurable
		Messages: []Message{
			{Role: "system", Content: "You are an AI assistant tasked with rewriting prompts to make them safe while preserving their intent."},
			{Role: "user", Content: rewritePrompt},
		},
	}
	
	response, err := r.llmClient.Chat(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get rewrite candidates from LLM: %w", err)
	}
	
	// Convert responses to candidates
	candidates := make([]Candidate, 0, len(response.Choices))
	for i, choice := range response.Choices {
		candidate := Candidate{
			ID:      fmt.Sprintf("candidate_%d", i),
			Content: choice.Message.Content,
			Score:   0.0, // Will be set by ranking model
		}
		candidates = append(candidates, candidate)
	}
	
	// If we didn't get enough candidates, generate more
	for len(candidates) < r.candidateCount {
		additionalCandidates, err := r.generateAdditionalCandidates(ctx, prompt, len(candidates))
		if err != nil {
			// Log error but continue with what we have
			break
		}
		candidates = append(candidates, additionalCandidates...)
	}
	
	// Limit to the requested count
	if len(candidates) > r.candidateCount {
		candidates = candidates[:r.candidateCount]
	}
	
	return candidates, nil
}

// createRewritePrompt creates the prompt for generating rewrite candidates
func (r *Rewriter) createRewritePrompt(prompt string) string {
	promptText := "Rewrite the following prompt to make it safe while preserving its intent:\n\n"
	promptText += fmt.Sprintf("%s\n\n", prompt)
	promptText += "Provide your response in the following format:\n"
	promptText += "REWRITTEN_PROMPT: [Your rewritten prompt]\n"
	promptText += "EXPLANATION: [Brief explanation of changes made]\n"
	promptText += "SAFETY_LEVEL: [HIGH/MEDIUM/LOW]\n\n"
	promptText += fmt.Sprintf("Generate %d different versions, each in the format above.", r.candidateCount)
	
	return promptText
}

// generateAdditionalCandidates generates additional candidates if needed
func (r *Rewriter) generateAdditionalCandidates(ctx context.Context, prompt string, currentCount int) ([]Candidate, error) {
	// Create a follow-up prompt to generate more candidates
	additionalPrompt := fmt.Sprintf(
		"Generate %d more alternative safe versions of this prompt:\n\n%s\n\n"+
			"Each version should be different from the previous ones.",
		r.candidateCount-currentCount, prompt)
	
	request := &ChatRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "system", Content: "You are an AI assistant tasked with rewriting prompts to make them safe while preserving their intent."},
			{Role: "user", Content: additionalPrompt},
		},
	}
	
	response, err := r.llmClient.Chat(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get additional candidates: %w", err)
	}
	
	// Convert responses to candidates
	candidates := make([]Candidate, 0, len(response.Choices))
	for i, choice := range response.Choices {
		candidate := Candidate{
			ID:      fmt.Sprintf("candidate_%d", currentCount+i),
			Content: choice.Message.Content,
			Score:   0.0, // Will be set by ranking model
		}
		candidates = append(candidates, candidate)
	}
	
	return candidates, nil
}

// ConfirmUserConfirmation marks the rewrite result as confirmed by the user
func (r *Rewriter) ConfirmUserConfirmation(result *RewriteResult, candidateID string) error {
	// Find the confirmed candidate
	for _, candidate := range result.Candidates {
		if candidate.ID == candidateID {
			result.SelectedCandidate = &candidate
			result.UserConfirmed = true
			return nil
		}
	}
	
	return fmt.Errorf("candidate %s not found", candidateID)
}

// ShouldReframe determines if the prompt should be reframed based on rewrite results
func (r *Rewriter) ShouldReframe(result *RewriteResult) bool {
	return result.SelectedCandidate != nil && result.UserConfirmed
}

// GetRewrittenPrompt returns the final rewritten prompt
func (r *Rewriter) GetRewrittenPrompt(result *RewriteResult) string {
	if result.SelectedCandidate != nil {
		return result.SelectedCandidate.Content
	}
	return result.OriginalPrompt
}package rewriter
