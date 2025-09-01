package router

import (
	"context"
	"fmt"
)

// Router determines the appropriate action for a given request
type Router struct {
	policyEngine PolicyEngine
	mode         string // audit, enforce, silent
}

// RouteDecision represents a routing decision
type RouteDecision struct {
	Action          string                 `json:"action"`           // allow, reframe, encrypt, block
	Reason          string                 `json:"reason"`           // reason for the decision
	Confidence      float64                `json:"confidence"`       // confidence in the decision
	Recommendations []string               `json:"recommendations"`  // recommended actions
	Metadata        map[string]interface{} `json:"metadata"`         // additional metadata
}

// PolicyEngine interface for evaluating policies
type PolicyEngine interface {
	// Evaluate evaluates policies for the given context
	Evaluate(ctx context.Context, input *PolicyInput) (*PolicyOutput, error)
}

// PolicyInput represents input to the policy engine
type PolicyInput struct {
	TenantID       string                 `json:"tenant_id"`
	UserID         string                 `json:"user_id"`
	Model          string                 `json:"model"`
	Prompt         string                 `json:"prompt"`
	DetectionScore float64                `json:"detection_score"`
	Reflection     *ReflectionResult      `json:"reflection,omitempty"`
	Rewrite        *RewriteResult         `json:"rewrite,omitempty"`
	Context        map[string]interface{} `json:"context"`
}

// PolicyOutput represents output from the policy engine
type PolicyOutput struct {
	Decision       string                 `json:"decision"`
	Reason         string                 `json:"reason"`
	Confidence     float64                `json:"confidence"`
	Recommendations []string              `json:"recommendations"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// ReflectionResult represents a reflection result
type ReflectionResult struct {
	Aligned    bool    `json:"aligned"`
	Confidence float64 `json:"confidence"`
}

// RewriteResult represents a rewrite result
type RewriteResult struct {
	UserConfirmed bool `json:"user_confirmed"`
}

// NewRouter creates a new router
func NewRouter(policyEngine PolicyEngine, mode string) *Router {
	return &Router{
		policyEngine: policyEngine,
		mode:         mode,
	}
}

// Route determines the appropriate action for a request
func (r *Router) Route(ctx context.Context, input *RoutingInput) (*RouteDecision, error) {
	// Evaluate policies
	policyInput := &PolicyInput{
		TenantID:       input.TenantID,
		UserID:         input.UserID,
		Model:          input.Model,
		Prompt:         input.Prompt,
		DetectionScore: input.DetectionScore,
		Reflection:     input.Reflection,
		Rewrite:        input.Rewrite,
		Context:        input.Context,
	}
	
	policyOutput, err := r.policyEngine.Evaluate(ctx, policyInput)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policies: %w", err)
	}
	
	// Create route decision based on policy output and mode
	decision := &RouteDecision{
		Action:          policyOutput.Decision,
		Reason:          policyOutput.Reason,
		Confidence:      policyOutput.Confidence,
		Recommendations: policyOutput.Recommendations,
		Metadata:        policyOutput.Metadata,
	}
	
	// Modify decision based on mode
	switch r.mode {
	case "audit":
		// In audit mode, allow everything but log the decision
		decision.Action = "allow"
		decision.Reason = "Audit mode - original decision was: " + policyOutput.Decision
	case "silent":
		// In silent mode, follow policy but don't inform user
		// Decision remains as policy output
	}
	
	return decision, nil
}

// RoutingInput represents input to the router
type RoutingInput struct {
	TenantID       string                 `json:"tenant_id"`
	UserID         string                 `json:"user_id"`
	Model          string                 `json:"model"`
	Prompt         string                 `json:"prompt"`
	DetectionScore float64                `json:"detection_score"`
	Reflection     *ReflectionResult      `json:"reflection,omitempty"`
	Rewrite        *RewriteResult         `json:"rewrite,omitempty"`
	Context        map[string]interface{} `json:"context"`
}

// ShouldAllow determines if the request should be allowed
func (r *Router) ShouldAllow(decision *RouteDecision) bool {
	return decision.Action == "allow"
}

// ShouldReframe determines if the request should be reframed
func (r *Router) ShouldReframe(decision *RouteDecision) bool {
	return decision.Action == "reframe"
}

// ShouldEncrypt determines if the response should be encrypted
func (r *Router) ShouldEncrypt(decision *RouteDecision) bool {
	return decision.Action == "encrypt"
}

// ShouldBlock determines if the request should be blocked
func (r *Router) ShouldBlock(decision *RouteDecision) bool {
	return decision.Action == "block"
}

// GetRecommendedActions returns recommended actions based on the decision
func (r *Router) GetRecommendedActions(decision *RouteDecision) []string {
	return decision.Recommendations
}

// GetMetadata returns metadata from the decision
func (r *Router) GetMetadata(decision *RouteDecision) map[string]interface{} {
	return decision.Metadata
}package router
