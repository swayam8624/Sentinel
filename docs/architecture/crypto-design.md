# Cryptographic Design for Sentinel + CipherMesh

## Overview

This document outlines the cryptographic design for Sentinel + CipherMesh, focusing on the implementation of HKDF-SHA-512 for key derivation, AES-256-GCM for authenticated encryption, and FF3-1 for format-preserving encryption.

## Key Derivation Function (HKDF)

### Algorithm Selection

We use HKDF (RFC 5869) with SHA-512 as the underlying hash function for key derivation. This provides:

- Strong pseudorandom key generation
- Resistance to length extension attacks
- Domain separation through info parameters
- Salt for uniqueness per operation

### Implementation Details

```go
// HKDF parameters
const (
    HKDFHash = crypto.SHA512
    SaltLength = 32  // 256 bits
    InfoPrefix = "sentinel-v1"
)

// Key derivation function
func DeriveKey(masterKey []byte, salt []byte, info string, length int) ([]byte, error) {
    hkdf := hkdf.New(HKDFHash.New, masterKey, salt, []byte(InfoPrefix+info))
    derivedKey := make([]byte, length)
    _, err := io.ReadFull(hkdf, derivedKey)
    if err != nil {
        return nil, err
    }
    return derivedKey, nil
}
```

### Key Derivation Process

1. **Per-Message Key Derivation**:

   - Salt: 32-byte random value generated for each message
   - Info: Includes "message-key" and token count as associated data
   - Output: 32-byte data encryption key (DEK)

2. **Per-Tenant Key Derivation**:

   - Salt: Tenant-specific value derived from tenant ID
   - Info: Includes "tenant-key" and timestamp
   - Output: Tenant-specific key encryption key (KEK)

3. **Per-Field FPE Key Derivation**:
   - Salt: Field-type specific value
   - Info: Includes "fpe-key" and field identifier
   - Output: Field-specific FPE key

### Security Considerations

- Random salt ensures uniqueness even with the same master key
- Info parameter provides domain separation
- SHA-512 provides 256-bit security level against collision attacks
- Key derivation is computationally expensive enough to deter brute force

## Authenticated Encryption (AES-256-GCM)

### Algorithm Selection

AES-256-GCM provides:

- Confidentiality through AES encryption
- Integrity through GHASH authentication
- Nonce misuse resistance (when properly implemented)
- Single-pass encryption and authentication

### Implementation Details

```go
// AES-GCM parameters
const (
    KeyLength = 32    // 256 bits
    NonceLength = 12  // 96 bits
    TagLength = 16    // 128 bits
)

// Encryption function
func EncryptAESGCM(key, nonce, plaintext, additionalData []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    aesGCM, err := cipher.NewGCMWithTagSize(block, TagLength)
    if err != nil {
        return nil, err
    }

    ciphertext := aesGCM.Seal(nil, nonce, plaintext, additionalData)
    return ciphertext, nil
}

// Decryption function
func DecryptAESGCM(key, nonce, ciphertext, additionalData []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    aesGCM, err := cipher.NewGCMWithTagSize(block, TagLength)
    if err != nil {
        return nil, err
    }

    plaintext, err := aesGCM.Open(nil, nonce, ciphertext, additionalData)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}
```

### Nonce Management

Critical for security to ensure nonce uniqueness:

1. **Generation Strategy**:

   - 12-byte random nonce for each encryption operation
   - Cryptographically secure random number generator
   - Storage of used nonces for duplicate detection

2. **Nonce Storage**:

   - Persistent storage for nonce tracking
   - Efficient lookup mechanisms
   - Automatic cleanup of expired nonces

3. **Collision Handling**:
   - Detection of duplicate nonces
   - Automatic nonce regeneration
   - Alerting on nonce collision events

### Additional Data

Using token count as additional authenticated data:

- Prevents tampering with token metadata
- Provides context for decryption
- Ensures integrity of associated information

## Format-Preserving Encryption (FF3-1)

### Algorithm Selection

FF3-1 (NIST SP 800-38G) provides:

- Format preservation for sensitive data
- Strong security guarantees
- Support for various character sets
- Resistance to statistical attacks

### Implementation Details

```go
// FF3-1 parameters
const (
    FF3KeyLength = 32  // 256 bits
    FF3TweakLength = 7 // 56 bits
)

// Encryption function
func EncryptFF3(key, tweak, plaintext []byte, radix int) ([]byte, error) {
    // Implementation of FF3-1 algorithm
    // ...
}

// Decryption function
func DecryptFF3(key, tweak, ciphertext []byte, radix int) ([]byte, error) {
    // Implementation of FF3-1 algorithm
    // ...
}
```

### Domain and Radix Definition

1. **Credit Card Numbers**:

   - Radix: 10 (digits 0-9)
   - Length: 13-19 characters
   - Character set: [0-9]

2. **Social Security Numbers**:

   - Radix: 10 (digits 0-9)
   - Length: 9 characters
   - Character set: [0-9]
   - Format: XXX-XX-XXXX

3. **Phone Numbers**:
   - Radix: 10 (digits 0-9)
   - Length: 10 characters
   - Character set: [0-9]
   - Format: (XXX) XXX-XXXX

### Tweak Generation

Unique per field to ensure semantic security:

1. **Per-Field Tweak**:

   - Combination of tenant ID, field type, and timestamp
   - 56-bit value derived from these components
   - Regular rotation for long-term security

2. **Tweak Storage**:
   - Associated with encrypted values
   - Used for deterministic decryption
   - Protected with AEAD encryption

### Security Properties

- Pseudorandom permutation over specified domain
- Resistance to frequency analysis
- Semantic security under chosen plaintext attacks
- Format preservation without information leakage

## Key Management Architecture

### Envelope Encryption

1. **Data Encryption Keys (DEKs)**:

   - Generated per message or per data item
   - Used for AES-GCM encryption
   - Encrypted with Key Encryption Keys (KEKs)

2. **Key Encryption Keys (KEKs)**:
   - Tenant-specific keys
   - Stored encrypted with BYOK/HSM keys
   - Regular rotation policy

### BYOK/HSM Integration

1. **Key Generation**:

   - Master keys generated in HSMs
   - Key export prevention
   - Hardware-based key generation

2. **Key Usage**:

   - KEK encryption/decryption in HSM
   - Key wrapping/unwrapping
   - Audit logging of key operations

3. **Key Rotation**:
   - Automated rotation schedules
   - Manual rotation triggers
   - Graceful transition periods

### Key Hierarchy

```
HSM Master Key (BYOK)
└── Tenant KEKs (encrypted with HSM key)
    └── Message DEKs (encrypted with tenant KEK)
        └── Data (encrypted with message DEK)
```

## Cryptographic Implementation Requirements

### Compliance

1. **NIST Standards**:

   - AES-256 (FIPS 197)
   - SHA-512 (FIPS 180-4)
   - FF3-1 (NIST SP 800-38G)
   - HKDF (NIST SP 800-56C)

2. **Industry Best Practices**:
   - OWASP Cryptographic Storage Cheat Sheet
   - PCI DSS requirements for encryption
   - HIPAA requirements for data protection

### Security Controls

1. **Key Isolation**:

   - Per-tenant key separation
   - No key sharing between tenants
   - Access control to key materials

2. **Zero Plaintext Storage**:

   - No plaintext in logs
   - No plaintext in memory longer than necessary
   - Secure memory handling

3. **Audit and Monitoring**:
   - Key usage logging
   - Cryptographic operation monitoring
   - Anomaly detection for crypto operations

### Performance Considerations

1. **Efficient Algorithms**:

   - Hardware-accelerated AES where available
   - Optimized HKDF implementation
   - Parallel processing where possible

2. **Caching Strategies**:

   - Key derivation result caching
   - Encrypted value caching
   - Balance between security and performance

3. **Resource Management**:
   - Memory allocation for crypto operations
   - Connection pooling for HSM operations
   - Concurrency control for crypto resources

## Testing and Validation

### Cryptographic Testing

1. **Algorithm Validation**:

   - NIST test vectors for AES, SHA-512
   - FF3-1 test vectors from NIST
   - HKDF test cases from RFC 5869

2. **Implementation Testing**:

   - Known answer tests
   - Monte Carlo tests for randomness
   - Boundary condition testing

3. **Integration Testing**:
   - End-to-end encryption/decryption
   - Key rotation scenarios
   - Multi-tenant isolation verification

### Security Testing

1. **Penetration Testing**:

   - Cryptographic attack simulations
   - Key recovery attempts
   - Side-channel analysis

2. **Compliance Testing**:

   - FIPS 140-2 validation
   - PCI DSS assessment
   - HIPAA compliance verification

3. **Code Review**:
   - Cryptographic implementation review
   - Third-party library security audit
   - Secure coding practice verification

## Error Handling and Recovery

### Cryptographic Errors

1. **Decryption Failures**:

   - Authentication tag verification failures
   - Padding errors
   - Key derivation failures

2. **Key Management Errors**:

   - Key not found
   - Key decryption failures
   - HSM communication errors

3. **Recovery Procedures**:
   - Graceful degradation
   - Fallback mechanisms
   - Alerting and logging

### Security Incident Response

1. **Key Compromise**:

   - Immediate key rotation
   - Affected data identification
   - Impact assessment

2. **Cryptographic Vulnerability**:
   - Algorithm migration procedures
   - Backward compatibility maintenance
   - Customer notification processes

## Future Considerations

### Post-Quantum Cryptography

1. **Migration Path**:

   - Hybrid classical/quantum-resistant algorithms
   - Key management for multiple algorithm sets
   - Performance impact assessment

2. **Implementation Timeline**:
   - Monitoring NIST PQC standardization
   - Prototype implementations
   - Gradual deployment strategy

### Advanced Cryptographic Features

1. **Homomorphic Encryption**:

   - For computation on encrypted data
   - Performance optimization research
   - Use case identification

2. **Secure Multi-Party Computation**:
   - For collaborative analytics
   - Privacy-preserving computations
   - Implementation complexity evaluation
