# Data Model Design for Sentinel + CipherMesh

## Overview

This document defines the core data models for Sentinel + CipherMesh, including the token vault, events, policies, and signatures. These models are designed to support multi-tenancy, security, and performance requirements.

## Core Entities

### 1. Tenant

Represents an organization or customer using the system.

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    region VARCHAR(100) NOT NULL,
    kms_key_arn TEXT,
    policy_version_pin UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_region ON tenants(region);
CREATE INDEX idx_tenants_status ON tenants(status);
```

**Fields:**

- `id`: Unique identifier for the tenant
- `name`: Human-readable name of the tenant
- `slug`: URL-friendly identifier for the tenant
- `status`: Current status (active, suspended, deleted)
- `region`: Data residency region
- `kms_key_arn`: ARN of the tenant's KMS key for BYOK
- `policy_version_pin`: Pinned policy version for the tenant
- `created_at`: Timestamp when the tenant was created
- `updated_at`: Timestamp when the tenant was last updated

### 2. Policy

Represents a set of rules for data handling and security.

```sql
CREATE TABLE policies (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    description TEXT,
    ruleset_json JSONB NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,  -- User ID reference
    UNIQUE(tenant_id, name, version)
);

-- Indexes
CREATE INDEX idx_policies_tenant ON policies(tenant_id);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_name ON policies(name);
CREATE INDEX idx_policies_version ON policies(version);
```

**Fields:**

- `id`: Unique identifier for the policy
- `tenant_id`: Reference to the tenant that owns this policy
- `name`: Human-readable name of the policy
- `version`: Semantic version of the policy (e.g., "1.0.0")
- `description`: Description of what the policy does
- `ruleset_json`: JSON representation of the policy rules
- `status`: Current status (draft, active, deprecated, archived)
- `created_at`: Timestamp when the policy was created
- `updated_at`: Timestamp when the policy was last updated
- `created_by`: Reference to the user who created the policy

### 3. Redaction Map (Token Vault)

Stores mappings between original sensitive data and redacted tokens.

```sql
CREATE TABLE redaction_maps (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    conversation_id UUID,
    token_id VARCHAR(255) NOT NULL,
    enc_value_aead BYTEA NOT NULL,  -- AEAD encrypted original value
    fpe_tweak BYTEA,  -- Tweak used for FPE, if applicable
    data_class VARCHAR(100) NOT NULL,  -- Type of data (pii, phi, pci, etc.)
    field_type VARCHAR(100),  -- Specific field type (ssn, credit_card, etc.)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    access_reason_code VARCHAR(100),  -- Reason for last access
    last_accessed_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_redaction_maps_tenant ON redaction_maps(tenant_id);
CREATE INDEX idx_redaction_maps_token ON redaction_maps(token_id);
CREATE INDEX idx_redaction_maps_conversation ON redaction_maps(conversation_id);
CREATE INDEX idx_redaction_maps_expires ON redaction_maps(expires_at);
CREATE INDEX idx_redaction_maps_class ON redaction_maps(data_class);
```

**Fields:**

- `id`: Unique identifier for the redaction map entry
- `tenant_id`: Reference to the tenant that owns this mapping
- `conversation_id`: Reference to the conversation this token belongs to
- `token_id`: The redacted token value
- `enc_value_aead`: The original value encrypted with AEAD
- `fpe_tweak`: Tweak used for FPE operations (if applicable)
- `data_class`: Classification of the data (pii, phi, pci, etc.)
- `field_type`: Specific type of field (ssn, credit_card, etc.)
- `created_at`: Timestamp when the mapping was created
- `expires_at`: Timestamp when the mapping expires
- `access_reason_code`: Reason code for the last access
- `last_accessed_at`: Timestamp when the mapping was last accessed

### 4. Event (Audit Log)

Records all system events for audit and compliance purposes.

```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    event_type VARCHAR(100) NOT NULL,
    ts TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    decision VARCHAR(50),  -- allow, block, encrypt, reframe
    scores JSONB,  -- Detection scores and metadata
    token_count INTEGER,
    ciphertext_len INTEGER,
    prev_hash BYTEA,  -- Previous event hash for chaining
    hash BYTEA NOT NULL,  -- This event's hash
    metadata JSONB,  -- Additional event-specific data
    request_id UUID,  -- Correlation ID for request tracing
    user_id UUID,  -- User who triggered the event
    ip_address INET,  -- IP address of the requester
    user_agent TEXT  -- User agent of the requester
);

-- Indexes
CREATE INDEX idx_events_tenant ON events(tenant_id);
CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_ts ON events(ts);
CREATE INDEX idx_events_decision ON events(decision);
CREATE INDEX idx_events_request ON events(request_id);
```

**Fields:**

- `id`: Unique identifier for the event
- `tenant_id`: Reference to the tenant this event belongs to
- `event_type`: Type of event (redaction, detection, policy_eval, etc.)
- `ts`: Timestamp of the event
- `decision`: Security decision made (allow, block, encrypt, reframe)
- `scores`: JSON object containing detection scores and metadata
- `token_count`: Number of tokens processed in this event
- `ciphertext_len`: Length of ciphertext if encryption occurred
- `prev_hash`: Hash of the previous event for chaining
- `hash`: Hash of this event for integrity
- `metadata`: Additional event-specific data
- `request_id`: Correlation ID for request tracing
- `user_id`: User who triggered the event
- `ip_address`: IP address of the requester
- `user_agent`: User agent of the requester

### 5. Signature

Stores known adversarial patterns for detection.

```sql
CREATE TABLE signatures (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    vector_id VARCHAR(255) NOT NULL,
    embedding VECTOR(768),  -- Assuming 768-dimensional embeddings
    label VARCHAR(100) NOT NULL,  -- Type of signature (jailbreak, injection, etc.)
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN NOT NULL DEFAULT true,
    confidence_threshold FLOAT  -- Threshold for matching this signature
);

-- Indexes
CREATE INDEX idx_signatures_tenant ON signatures(tenant_id);
CREATE INDEX idx_signatures_label ON signatures(label);
CREATE INDEX idx_signatures_active ON signatures(active);
CREATE INDEX idx_signatures_embedding ON signatures USING ivfflat (embedding vector_cosine_ops);
```

**Fields:**

- `id`: Unique identifier for the signature
- `tenant_id`: Reference to the tenant this signature belongs to
- `vector_id`: Identifier for the vector in the vector store
- `embedding`: Vector embedding of the signature pattern
- `label`: Type of signature (jailbreak, injection, etc.)
- `description`: Human-readable description of the signature
- `created_at`: Timestamp when the signature was created
- `updated_at`: Timestamp when the signature was last updated
- `active`: Whether the signature is currently active
- `confidence_threshold`: Threshold for matching this signature

## Relationship Diagram

```
Tenants
│
├── Policies
│   └── Policy Versions (via version field)
│
├── Redaction Maps
│   ├── Conversation Context
│   ├── Data Classes
│   └── Field Types
│
├── Events
│   ├── Event Types
│   ├── Security Decisions
│   └── Detection Scores
│
└── Signatures
    ├── Signature Labels
    └── Confidence Thresholds
```

## Multi-Tenancy Considerations

### Data Isolation

1. **Hard Separation**:

   - All tables include `tenant_id` as a foreign key
   - Queries always filter by `tenant_id`
   - Database-level separation for high-security tenants

2. **Indexing Strategy**:
   - All primary indexes include `tenant_id`
   - Secondary indexes prefixed with `tenant_id`
   - Partitioning by tenant where appropriate

### Performance Optimization

1. **Caching**:

   - Tenant-specific cache partitions
   - Policy caching with version invalidation
   - Hot token mapping cache

2. **Query Optimization**:
   - Composite indexes for tenant-specific queries
   - Materialized views for common aggregations
   - Read replicas for audit queries

## Security Considerations

### Data Protection

1. **Encryption at Rest**:

   - All sensitive data encrypted with tenant-specific keys
   - AEAD encryption for redaction maps
   - Envelope encryption for all data

2. **Access Controls**:

   - Row-level security policies
   - Role-based access controls
   - Audit logging for all data access

3. **Data Minimization**:
   - Only store necessary information
   - Automatic purging of expired data
   - Anonymization where possible

### Audit and Compliance

1. **Immutable Logs**:

   - Event hash chaining for integrity
   - Append-only storage for critical logs
   - Regular Merkle root anchoring

2. **Compliance Features**:
   - GDPR right to erasure support
   - HIPAA audit trail requirements
   - PCI DSS logging requirements

## Scalability Considerations

### Horizontal Scaling

1. **Sharding Strategy**:

   - Tenant-based sharding
   - Geographic sharding for global deployments
   - Consistent hashing for even distribution

2. **Read Scaling**:
   - Read replicas for audit and reporting
   - Caching layers for hot data
   - Asynchronous processing for non-critical operations

### Performance Optimization

1. **Indexing**:

   - Composite indexes for common query patterns
   - Partial indexes for filtered data
   - Covering indexes for frequently accessed fields

2. **Partitioning**:
   - Time-based partitioning for events
   - Tenant-based partitioning for all tables
   - Automatic partition management

## Data Lifecycle Management

### Retention Policies

1. **Redaction Maps**:

   - Configurable TTL based on policy
   - Automatic cleanup of expired mappings
   - Extension mechanisms for compliance requirements

2. **Events**:

   - Configurable retention per tenant
   - Archival to cold storage
   - Compliance-driven retention periods

3. **Signatures**:
   - Version history retention
   - Archive inactive signatures
   - Periodic review and cleanup

### Data Archival

1. **Cold Storage**:

   - Move old data to cost-effective storage
   - Maintain query access for compliance
   - Automated archival processes

2. **Data Deletion**:
   - Secure deletion of expired data
   - Compliance with right to erasure
   - Audit trail of deletion activities

## Monitoring and Observability

### Metrics Collection

1. **Database Metrics**:

   - Query performance by tenant
   - Storage utilization
   - Connection pool statistics

2. **Business Metrics**:
   - Redaction rates by data class
   - Policy evaluation outcomes
   - Security incident frequency

### Alerting

1. **Performance Alerts**:

   - Slow query detection
   - Storage capacity warnings
   - Connection pool exhaustion

2. **Security Alerts**:
   - Unusual access patterns
   - Policy bypass attempts
   - Data exfiltration indicators

## Implementation Considerations

### Technology Stack

1. **Database**:

   - PostgreSQL with extensions (pgvector for embeddings)
   - Connection pooling (PgBouncer)
   - Read replicas for scaling

2. **Caching**:

   - Redis for hot data
   - Tenant-specific cache namespaces
   - Cache invalidation strategies

3. **Search**:
   - Elasticsearch for audit log search
   - Full-text search capabilities
   - Aggregation and analytics

### Data Migration

1. **Schema Evolution**:

   - Backward-compatible schema changes
   - Migration scripts for breaking changes
   - Rollback procedures

2. **Data Migration**:
   - Online migration tools
   - Data validation procedures
   - Performance impact minimization

## Future Extensions

### Advanced Analytics

1. **Machine Learning Features**:

   - Anomaly detection in access patterns
   - Predictive policy tuning
   - Automated signature generation

2. **Graph Database Integration**:
   - Relationship analysis between entities
   - Complex policy dependency mapping
   - Investigation workflows

### Enhanced Security

1. **Blockchain Integration**:

   - Immutable audit trail anchoring
   - Smart contract-based policy enforcement
   - Decentralized key management

2. **Confidential Computing**:
   - Hardware-based security enclaves
   - Secure multi-party computation
   - Zero-trust architecture
