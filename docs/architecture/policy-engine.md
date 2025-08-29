# OPA Policy Engine Design for Sentinel + CipherMesh

## Overview

This document outlines the design and implementation of the OPA-based policy engine for Sentinel + CipherMesh. The policy engine is responsible for evaluating rules that govern data handling, security responses, and system behavior.

## Architecture

### Core Components

1. **Policy Store**: Centralized storage for policy definitions
2. **Policy Evaluator**: OPA runtime for policy evaluation
3. **Policy Compiler**: Transforms policy definitions into Rego
4. **Version Manager**: Handles policy versioning and deployment
5. **Audit Logger**: Records all policy evaluation results

### Integration Flow

```
[Input Request]
       ↓
[Policy Resolver] → Determines applicable policies based on tenant, version, etc.
       ↓
[Policy Evaluator] → Executes OPA evaluation with input data
       ↓
[Decision Engine] → Processes OPA results and determines actions
       ↓
[Action Executor] → Executes determined actions (allow, block, encrypt, etc.)
```

## Policy Structure

### YAML Policy Definition

Policies are defined in YAML format for human readability:

```yaml
policy:
  id: "example-policy"
  name: "PII Protection Policy"
  version: "1.0.0"
  description: "Policy for protecting PII data"
  tags: ["pii", "gdpr", "hipaa"]
  severity: "high"
  enabled: true
  effective_from: "2023-01-01T00:00:00Z"
  effective_to: "2024-01-01T00:00:00Z"

  scope:
    tenants: ["*"] # or specific tenant IDs
    models: ["gpt-4", "claude-2"] # or ["*"] for all
    data_classes: ["pii", "phi"]

  rules:
    - id: "detect-ssn"
      description: "Detect and protect Social Security Numbers"
      type: "detection"
      priority: 100
      conditions:
        - field: "content"
          operator: "regex_match"
          value: "\\d{3}-\\d{2}-\\d{4}"
          confidence: 0.9
      actions:
        - type: "redact"
          method: "fpe"
          format: "XXX-XX-NNNN"

    - id: "block-jailbreak"
      description: "Block known jailbreak attempts"
      type: "security"
      priority: 200
      conditions:
        - field: "embedding_similarity"
          operator: "greater_than"
          value: 0.85
          against: "jailbreak_signatures"
      actions:
        - type: "block"
          reason: "jailbreak_attempt"

    - id: "encrypt-sensitive-output"
      description: "Encrypt sensitive model outputs"
      type: "response"
      priority: 150
      conditions:
        - field: "response.sensitivity_score"
          operator: "greater_than"
          value: 0.8
      actions:
        - type: "encrypt"
          algorithm: "aes-256-gcm"
```

### Rego Translation

The YAML policies are translated into Rego for OPA evaluation:

```rego
package sentinel.policy.example

# Rule: detect-ssn
detect_ssn[decision] {
    # Input data
    input.content
    input.tenant_id

    # Condition matching
    re_match(`\d{3}-\d{2}-\d{4}`, input.content)

    # Decision
    decision := {
        "rule_id": "detect-ssn",
        "matched": true,
        "confidence": 0.9,
        "actions": [
            {
                "type": "redact",
                "method": "fpe",
                "format": "XXX-XX-NNNN"
            }
        ]
    }
}

# Rule: block-jailbreak
block_jailbreak[decision] {
    # Input data
    input.embedding_similarity
    input.tenant_id

    # Condition matching
    input.embedding_similarity > 0.85

    # Decision
    decision := {
        "rule_id": "block-jailbreak",
        "matched": true,
        "confidence": input.embedding_similarity,
        "actions": [
            {
                "type": "block",
                "reason": "jailbreak_attempt"
            }
        ]
    }
}

# Rule: encrypt-sensitive-output
encrypt_sensitive_output[decision] {
    # Input data
    input.response.sensitivity_score
    input.tenant_id

    # Condition matching
    input.response.sensitivity_score > 0.8

    # Decision
    decision := {
        "rule_id": "encrypt-sensitive-output",
        "matched": true,
        "confidence": input.response.sensitivity_score,
        "actions": [
            {
                "type": "encrypt",
                "algorithm": "aes-256-gcm"
            }
        ]
    }
}

# Main policy evaluation
evaluate = decisions {
    # Collect all rule decisions
    detect_ssn_decisions := [d | d := detect_ssn[_]]
    block_jailbreak_decisions := [d | d := block_jailbreak[_]]
    encrypt_decisions := [d | d := encrypt_sensitive_output[_]]

    # Combine all decisions
    decisions := array.concat(
        array.concat(detect_ssn_decisions, block_jailbreak_decisions),
        encrypt_decisions
    )
}
```

## Policy Versioning

### Versioning Scheme

Policies follow semantic versioning (SemVer):

- MAJOR: Breaking changes to policy structure or behavior
- MINOR: Backward-compatible additions or improvements
- PATCH: Backward-compatible bug fixes

### Deployment Modes

1. **Audit Mode**: Policies evaluated but not enforced
2. **Shadow Mode**: Policies evaluated alongside existing policies
3. **Enforce Mode**: Policies fully enforced

### Version Pinning

Tenants can pin to specific policy versions:

```yaml
tenant_policy_config:
  tenant_id: "tenant-123"
  policy_pins:
    "pii-protection": "1.2.3"
    "security-policy": "2.1.0"
  default_mode: "enforce"
```

## Multi-Tenancy Support

### Tenant Isolation

1. **Policy Namespacing**:

   - Policies scoped to specific tenants
   - Default policies for all tenants
   - Tenant-specific policy overrides

2. **Data Separation**:
   - Tenant-specific policy storage
   - Isolated policy evaluation contexts
   - Cross-tenant policy sharing controls

### Policy Inheritance

```
[Default Policies] (System-wide)
        ↓
[Tenant Policies] (Tenant-specific overrides)
        ↓
[User Policies] (User-specific overrides)
```

## Policy Evaluation

### Input Data Structure

```json
{
  "request": {
    "id": "req-123",
    "tenant_id": "tenant-456",
    "user_id": "user-789",
    "timestamp": "2023-01-01T00:00:00Z",
    "model": "gpt-4",
    "messages": [
      {
        "role": "user",
        "content": "Hello, my SSN is 123-45-6789"
      }
    ]
  },
  "context": {
    "data_classes": ["pii"],
    "sensitivity_score": 0.95,
    "embedding_similarity": 0.87,
    "tools_allowed": true
  },
  "response": {
    "content": "I can't help with that SSN.",
    "sensitivity_score": 0.92
  }
}
```

### Evaluation Process

1. **Policy Resolution**:

   - Identify applicable policies based on tenant, model, data classes
   - Apply version pins and mode settings
   - Resolve policy dependencies

2. **OPA Evaluation**:

   - Convert input data to OPA format
   - Execute policy evaluation
   - Collect all decision results

3. **Decision Aggregation**:
   - Sort decisions by priority
   - Resolve conflicting actions
   - Determine final action

### Conflict Resolution

When multiple policies suggest conflicting actions:

1. **Priority-based**: Higher priority policies take precedence
2. **Deny-by-default**: Block/encrypt actions override allow actions
3. **Explicit overrides**: Explicit deny overrides other actions

## Policy Management API

### CRUD Operations

```http
# Create a new policy
POST /api/v1/policies
Content-Type: application/yaml

# Get a specific policy
GET /api/v1/policies/{policy_id}

# Update a policy
PUT /api/v1/policies/{policy_id}

# Delete a policy
DELETE /api/v1/policies/{policy_id}

# List policies
GET /api/v1/policies?tenant_id={tenant_id}&tags={tag1,tag2}
```

### Version Management

```http
# Create a new version
POST /api/v1/policies/{policy_id}/versions

# Get specific version
GET /api/v1/policies/{policy_id}/versions/{version}

# List versions
GET /api/v1/policies/{policy_id}/versions

# Set active version
PUT /api/v1/policies/{policy_id}/versions/{version}/active
```

### Deployment Control

```http
# Deploy policy to tenant
POST /api/v1/policies/{policy_id}/deploy
{
  "tenant_id": "tenant-123",
  "mode": "audit",
  "pin_version": "1.2.3"
}

# Undeploy policy from tenant
DELETE /api/v1/policies/{policy_id}/deploy?tenant_id={tenant_id}

# Get deployment status
GET /api/v1/policies/{policy_id}/deploy?tenant_id={tenant_id}
```

## Policy Testing Framework

### Unit Testing

Rego unit tests for individual policies:

```rego
# Test for SSN detection
test_detect_ssn {
    input := {
        "content": "My SSN is 123-45-6789"
    }

    decisions := data.sentinel.policy.example.detect_ssn

    count(decisions) == 1
    decisions[0].matched == true
    decisions[0].actions[0].type == "redact"
}
```

### Integration Testing

End-to-end policy evaluation tests:

```go
func TestPIIPolicy(t *testing.T) {
    // Setup test policy
    policy := loadTestPolicy("pii-protection.yaml")

    // Setup test input
    input := PolicyInput{
        Request: Request{
            Content: "Hello, my SSN is 123-45-6789",
            TenantID: "test-tenant",
        },
        Context: Context{
            DataClasses: []string{"pii"},
        },
    }

    // Evaluate policy
    result := EvaluatePolicy(policy, input)

    // Assert expected actions
    assert.Equal(t, 1, len(result.Actions))
    assert.Equal(t, "redact", result.Actions[0].Type)
    assert.Equal(t, "fpe", result.Actions[0].Method)
}
```

## Performance Considerations

### Caching Strategy

1. **Policy Caching**:

   - Compiled Rego policies cached in memory
   - Cache invalidation on policy updates
   - Tenant-specific policy bundles

2. **Evaluation Caching**:
   - Cache results for identical inputs
   - Time-based expiration
   - Cache warming for frequently used policies

### Optimization Techniques

1. **Partial Evaluation**:

   - Pre-compute static portions of policies
   - Optimize for common evaluation paths
   - Reduce runtime evaluation overhead

2. **Indexing**:
   - Index policies by tenant, model, data class
   - Fast policy resolution
   - Efficient rule matching

## Security Considerations

### Policy Integrity

1. **Signing**:

   - Cryptographically sign policy definitions
   - Verify signatures before policy evaluation
   - Prevent policy tampering

2. **Access Control**:
   - RBAC for policy management
   - Audit logging for all policy changes
   - Approval workflows for critical policies

### Sandboxing

1. **Resource Limits**:

   - CPU and memory limits for policy evaluation
   - Timeout controls
   - Prevent denial of service

2. **Function Restrictions**:
   - Limit available built-in functions
   - Prevent external network access
   - Control file system access

## Monitoring and Observability

### Metrics Collection

1. **Evaluation Metrics**:

   - Policy evaluation latency
   - Decision distribution (allow, block, encrypt)
   - Rule match rates

2. **Performance Metrics**:
   - Cache hit/miss ratios
   - Memory usage
   - CPU utilization

### Logging

1. **Evaluation Logs**:

   - Input data (PII-safe)
   - Policy matches and actions
   - Evaluation duration

2. **Audit Logs**:
   - Policy changes
   - Deployment activities
   - Access control events

## Error Handling

### Policy Errors

1. **Syntax Errors**:

   - Validation during policy creation
   - Detailed error messages
   - Line number reporting

2. **Evaluation Errors**:
   - Graceful degradation
   - Fallback to default policies
   - Alerting on persistent errors

### Recovery Mechanisms

1. **Rollback**:

   - Automatic rollback on policy errors
   - Manual rollback capability
   - Version history retention

2. **Fallback**:
   - Default deny policies
   - Safe mode operation
   - Emergency bypass procedures

## Future Enhancements

### Machine Learning Integration

1. **Policy Learning**:

   - Automatically generate policies from data
   - Adaptive policy tuning
   - Anomaly detection

2. **Continuous Improvement**:
   - A/B testing for policies
   - Performance optimization
   - Feedback loops

### Advanced Policy Features

1. **Temporal Policies**:

   - Time-based policy activation
   - Scheduled policy changes
   - Event-driven policies

2. **Context-Aware Policies**:
   - Location-based policies
   - Device-based policies
   - Behavioral policies
