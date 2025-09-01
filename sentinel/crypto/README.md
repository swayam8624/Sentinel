# Sentinel Crypto Components

This directory contains all the cryptographic components for the Sentinel platform, implementing the security requirements specified in the SRS.

## Components Overview

### 1. HKDF (HMAC-based Key Derivation Function)

- **Location**: [hkdf/](hkdf/)
- **Purpose**: Key derivation using HKDF-SHA-512
- **Features**:
  - Per-message key derivation
  - Random salt generation
  - Contextual info parameter support
- **Specification**: RFC 5869

### 2. Nonce Management

- **Location**: [nonce/](nonce/)
- **Purpose**: Nonce generation and uniqueness enforcement
- **Features**:
  - Cryptographically secure nonce generation
  - Uniqueness tracking with expiration
  - Automatic cleanup of expired nonces
  - Collision detection and prevention

### 3. KMS (Key Management Service)

- **Location**: [kms/](kms/)
- **Purpose**: Key management for envelope encryption
- **Features**:
  - Data key generation and encryption
  - AES-GCM encryption/decryption
  - Envelope encryption pattern
  - Master key management

### 4. FPE (Format Preserving Encryption)

- **Location**: [fpe/](fpe/)
- **Purpose**: Format-preserving encryption for sensitive data
- **Features**:
  - FF3-1 implementation (simplified)
  - Credit card number validation (Luhn algorithm)
  - Format preservation for numeric data

### 5. Merkle Tree

- **Location**: [merkle/](merkle/)
- **Purpose**: Tamper-evident audit logs
- **Features**:
  - Merkle tree construction
  - Root hash calculation
  - Merkle proof generation
  - Proof verification

### 6. Token Vault

- **Location**: [vault/](vault/)
- **Purpose**: Secure storage for encrypted tokens
- **Features**:
  - AES-GCM encrypted storage
  - Token-based retrieval
  - Time-to-live enforcement
  - Access reason tracking

## Security Compliance

These implementations satisfy the cryptographic requirements specified in the SRS:

1. **HKDF-SHA-512 per-message key derivation** (B2, F1)
2. **AES-GCM nonce policy with uniqueness enforcement** (B2, F1)
3. **FF3-1 FPE for format-preserving encryption** (B2, C3)
4. **Envelop encryption with KMS integration** (F2)
5. **Tamper-evident logs with Merkle hash chains** (F4)
6. **Secure token storage with access tracking** (F3)

## Usage Instructions

### Running Tests

```bash
# Test all crypto components
go test ./sentinel/crypto/... -v

# Test individual components
go test ./sentinel/crypto/hkdf/... -v
go test ./sentinel/crypto/nonce/... -v
go test ./sentinel/crypto/kms/... -v
go test ./sentinel/crypto/fpe/... -v
go test ./sentinel/crypto/merkle/... -v
go test ./sentinel/crypto/vault/... -v
```

### Running Examples

```bash
# Run crypto example (HKDF + Nonce)
go run sentinel/crypto/example/main.go

# Run Merkle tree example
go run sentinel/crypto/merkle/example/main.go

# Run FPE example
go run sentinel/crypto/fpe/example/main.go

# Run vault example
go run sentinel/crypto/vault/example/main.go
```

## Component Details

### HKDF Implementation

The HKDF implementation provides secure key derivation using HMAC-SHA-512. It's used for deriving per-message keys from a master key, with random salts for each derivation.

### Nonce Management

The nonce manager ensures unique nonces for AES-GCM encryption, preventing cryptographic vulnerabilities that could arise from nonce reuse.

### KMS Implementation

The KMS implementation provides envelope encryption, where data keys are encrypted with a master key. In production, this would integrate with cloud KMS services.

### FPE Implementation

The FPE implementation provides format-preserving encryption for sensitive data like credit card numbers, ensuring the encrypted data maintains the same format as the original.

### Merkle Tree Implementation

The Merkle tree implementation provides tamper-evident audit logs, where any modification to log entries can be detected through cryptographic proofs.

### Token Vault Implementation

The token vault provides secure storage for encrypted tokens with TTL enforcement and access reason tracking for audit purposes.

## Next Steps

### Task F2: KMS/HSM Integration

- Implement cloud KMS integrations (AWS KMS, GCP KMS, Azure Key Vault)
- Implement BYOK flows
- Create key rotation mechanisms

## Design Decisions

1. **Local Implementations**: Most implementations are designed for local development and testing. In production, these would integrate with actual security services.

2. **Simplified Algorithms**: Some algorithms (like FF3-1) are simplified for demonstration. Production implementations would follow full specifications.

3. **Security First**: All implementations prioritize security and follow cryptographic best practices.

4. **Extensibility**: Components are designed to be extensible for future enhancements and integrations.
