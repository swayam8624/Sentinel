package policy

import (
	"context"
	"fmt"
	"time"
)

// Policy represents a security policy
type Policy struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Rules       map[string]interface{} `json:"rules"`
	Enabled     bool                   `json:"enabled"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Engine handles policy management and evaluation
type Engine struct {
	policies map[string]*Policy
}

// NewEngine creates a new policy engine
func NewEngine() *Engine {
	return &Engine{
		policies: make(map[string]*Policy),
	}
}

// CreatePolicy creates a new policy
func (e *Engine) CreatePolicy(ctx context.Context, policy *Policy) error {
	if policy.ID == "" {
		return fmt.Errorf("policy ID is required")
	}

	if _, exists := e.policies[policy.ID]; exists {
		return fmt.Errorf("policy with ID %s already exists", policy.ID)
	}

	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()
	e.policies[policy.ID] = policy

	return nil
}

// GetPolicy retrieves a policy by ID
func (e *Engine) GetPolicy(ctx context.Context, id string) (*Policy, error) {
	policy, exists := e.policies[id]
	if !exists {
		return nil, fmt.Errorf("policy with ID %s not found", id)
	}

	return policy, nil
}

// UpdatePolicy updates an existing policy
func (e *Engine) UpdatePolicy(ctx context.Context, policy *Policy) error {
	if policy.ID == "" {
		return fmt.Errorf("policy ID is required")
	}

	existing, exists := e.policies[policy.ID]
	if !exists {
		return fmt.Errorf("policy with ID %s not found", policy.ID)
	}

	policy.CreatedAt = existing.CreatedAt
	policy.UpdatedAt = time.Now()
	e.policies[policy.ID] = policy

	return nil
}

// DeletePolicy deletes a policy by ID
func (e *Engine) DeletePolicy(ctx context.Context, id string) error {
	if _, exists := e.policies[id]; !exists {
		return fmt.Errorf("policy with ID %s not found", id)
	}

	delete(e.policies, id)
	return nil
}

// ListPolicies returns all policies
func (e *Engine) ListPolicies(ctx context.Context) ([]*Policy, error) {
	policies := make([]*Policy, 0, len(e.policies))
	for _, policy := range e.policies {
		policies = append(policies, policy)
	}

	return policies, nil
}

// Evaluate evaluates a request against all enabled policies
func (e *Engine) Evaluate(ctx context.Context, request interface{}) (*EvaluationResult, error) {
	// In a real implementation, this would evaluate the request against all enabled policies
	// For now, we'll return a simple result
	return &EvaluationResult{
		Allowed: true,
		Reason:  "No policies configured",
	}, nil
}

// EvaluationResult represents the result of a policy evaluation
type EvaluationResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}
