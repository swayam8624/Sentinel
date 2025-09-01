# Policy Engine

The policy engine evaluates OPA-style rules to make decisions about data handling, security responses, and system behavior.

## Overview

The policy engine uses Rego (OPA's policy language) to define rules that govern:

- Data classification and redaction actions
- Security violation responses
- Tool and function call permissions
- Detokenization permissions
- Tenant-specific configurations

## Policy Structure

Policies are defined in YAML/JSON format and can include:

```yaml
policy:
  id: "example-policy"
  version: "1.0.0"
  description: "Example policy for handling PII data"
  rules:
    - id: "pii-handling"
      description: "How to handle PII data"
      conditions:
        - data_class: "pii"
          action: "fpe"
      permissions:
        detokenize:
          roles: ["analyst", "admin"]
```

## Policy Evaluation

The policy engine evaluates policies in the following order:

1. **Tenant-specific policies** (if applicable)
2. **Version-pinned policies** (if specified)
3. **Default policies**

## Policy Versioning

Policies support semantic versioning with the following capabilities:

- **Audit Mode**: Test policies without enforcement
- **Canary Deployment**: Roll out to subset of traffic
- **Enforce Mode**: Full policy enforcement
- **Rollback**: Revert to previous policy versions

## Multi-tenancy

Policies support multi-tenancy through:

- Tenant-scoped policy definitions
- Policy version pinning per tenant
- Shadow evaluation (compare multiple policy versions)
