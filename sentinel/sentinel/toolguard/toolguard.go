package toolguard

import (
	"context"
	"fmt"
)

// ToolGuard manages tool and function call permissions
type ToolGuard struct {
	policyEngine PolicyEngine
	mode         string // audit, enforce, silent
}

// ToolPermission represents permissions for a specific tool
type ToolPermission struct {
	ToolName    string   `json:"tool_name"`
	Allowed     bool     `json:"allowed"`
	Conditions  []string `json:"conditions"`
	Restrictions []string `json:"restrictions"`
}

// PolicyEngine interface for evaluating tool policies
type PolicyEngine interface {
	// EvaluateTool evaluates tool permissions for the given context
	EvaluateTool(ctx context.Context, input *ToolPolicyInput) (*ToolPolicyOutput, error)
}

// ToolPolicyInput represents input to the tool policy engine
type ToolPolicyInput struct {
	TenantID    string                 `json:"tenant_id"`
	UserID      string                 `json:"user_id"`
	ToolName    string                 `json:"tool_name"`
	ToolInput   map[string]interface{} `json:"tool_input"`
	Context     map[string]interface{} `json:"context"`
	Violation   bool                   `json:"violation"`
}

// ToolPolicyOutput represents output from the tool policy engine
type ToolPolicyOutput struct {
	Allowed      bool                   `json:"allowed"`
	Reason       string                 `json:"reason"`
	Conditions   []string              `json:"conditions"`
	Restrictions []string              `json:"restrictions"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// NewToolGuard creates a new ToolGuard
func NewToolGuard(policyEngine PolicyEngine, mode string) *ToolGuard {
	return &ToolGuard{
		policyEngine: policyEngine,
		mode:         mode,
	}
}

// CheckToolPermission checks if a tool call is permitted
func (tg *ToolGuard) CheckToolPermission(ctx context.Context, input *ToolCheckInput) (*ToolPermission, error) {
	// Evaluate tool policies
	policyInput := &ToolPolicyInput{
		TenantID:  input.TenantID,
		UserID:    input.UserID,
		ToolName:  input.ToolName,
		ToolInput: input.ToolInput,
		Context:   input.Context,
		Violation: input.Violation,
	}
	
	policyOutput, err := tg.policyEngine.EvaluateTool(ctx, policyInput)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate tool policies: %w", err)
	}
	
	// Create tool permission based on policy output and mode
	permission := &ToolPermission{
		ToolName:     input.ToolName,
		Allowed:      policyOutput.Allowed,
		Conditions:   policyOutput.Conditions,
		Restrictions: policyOutput.Restrictions,
	}
	
	// Modify permission based on mode
	switch tg.mode {
	case "audit":
		// In audit mode, allow everything but log the decision
		permission.Allowed = true
	case "silent":
		// In silent mode, follow policy but don't inform user
		// Permission remains as policy output
	}
	
	return permission, nil
}

// ToolCheckInput represents input to the tool permission check
type ToolCheckInput struct {
	TenantID  string                 `json:"tenant_id"`
	UserID    string                 `json:"user_id"`
	ToolName  string                 `json:"tool_name"`
	ToolInput map[string]interface{} `json:"tool_input"`
	Context   map[string]interface{} `json:"context"`
	Violation bool                   `json:"violation"`
}

// IsToolAllowed determines if a tool is allowed based on permission
func (tg *ToolGuard) IsToolAllowed(permission *ToolPermission) bool {
	return permission.Allowed
}

// GetToolConditions returns conditions for tool usage
func (tg *ToolGuard) GetToolConditions(permission *ToolPermission) []string {
	return permission.Conditions
}

// GetToolRestrictions returns restrictions for tool usage
func (tg *ToolGuard) GetToolRestrictions(permission *ToolPermission) []string {
	return permission.Restrictions
}

// ApplyToolRestrictions applies restrictions to tool input
func (tg *ToolGuard) ApplyToolRestrictions(input map[string]interface{}, restrictions []string) map[string]interface{} {
	// Apply restrictions to the input
	restrictedInput := make(map[string]interface{})
	
	for key, value := range input {
		// Check if this key is restricted
		restricted := false
		for _, restriction := range restrictions {
			if restriction == key || restriction == "*" {
				restricted = true
				break
			}
		}
		
		if !restricted {
			restrictedInput[key] = value
		}
		// If restricted, we simply don't include it in the restricted input
	}
	
	return restrictedInput
}

// Lockdown disables all tools when a violation is detected
func (tg *ToolGuard) Lockdown() *ToolPermission {
	return &ToolPermission{
		ToolName:    "*",
		Allowed:     false,
		Conditions:  []string{"violation_detected"},
		Restrictions: []string{"*"},
	}
}

// IsLockdown determines if the system is in lockdown mode
func (tg *ToolGuard) IsLockdown(permission *ToolPermission) bool {
	return permission.ToolName == "*" && !permission.Allowed
}package toolguard
