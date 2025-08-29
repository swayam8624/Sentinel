# Cryptography Approach for Sentinel + CipherMesh

## Status

Accepted

## Context

We need to establish a strong cryptographic foundation for Sentinel + CipherMesh that:

1. Protects sensitive data both in transit and at rest
2. Implements secure key management with BYOK support
3. Provides format-preserving encryption for data utility
4. Ensures unique nonce generation to prevent replay attacks
5. Supports tenant isolation and data residency requirements
6. Complies with industry standards (NIST, FIPS)

We evaluated several approaches for implementing cryptography:

- Custom cryptographic implementations
- Using established libraries (OpenSSL, libsodium, etc.)
- Cloud KMS/HSM integration
- Hybrid approach combining libraries and cloud services

## Decision

We will implement a **hybrid cryptographic approach** using established libraries and cloud KMS/HSM services:

1. **Key Derivation**: HKDF-SHA-512 (RFC 5869) with per-message random salt
2. **Authenticated Encryption**: AES-256-GCM (NIST SP 800-38D) with unique nonces
3. **Format-Preserving Encryption**: FF3-1 (NIST SP 800-38G) for PII tokenization
4. **Key Management**: Cloud KMS/HSM integration for BYOK support
5. **Envelope Encryption**: Data encryption keys (DEKs) encrypted with key encryption keys (KEKs)
6. **Nonce Management**: Systematic nonce generation to ensure uniqueness

Implementation details:

- Use well-established cryptographic libraries (e.g., OpenSSL, BoringSSL)
- Integrate with major cloud KMS services (AWS KMS, Azure Key Vault, GCP KMS)
- Implement strict nonce management with counters or random generation
- Enforce key rotation policies
- Support split-knowledge for vault master keys (optional)

## Consequences

### Positive Consequences

- Leverages well-tested, industry-standard cryptographic algorithms
- Supports compliance with regulatory requirements (GDPR, HIPAA, etc.)
- Enables BYOK through cloud KMS/HSM integration
- Provides strong security guarantees with proper implementation
- Supports tenant isolation through separate key hierarchies
- Allows for key rotation and cryptographic agility

### Negative Consequences

- Complexity in implementing and managing key hierarchies
- Dependency on external KMS/HSM services
- Potential performance overhead from cryptographic operations
- Need for specialized knowledge in cryptographic implementation
- Risk of implementation errors leading to security vulnerabilities
