# Token Vault

This directory contains the token vault implementation for the Sentinel platform, providing secure storage for encrypted tokens and sensitive data.

## Overview

The vault implementation provides:

1. **Secure Storage**: AES-GCM encrypted storage for sensitive data
2. **Token Management**: Store, retrieve, and delete tokens with TTL
3. **Access Reason Tracking**: Record access reasons for audit purposes
4. **Automatic Expiration**: Automatic cleanup of expired tokens

## Features

- AES-GCM encryption for data at rest
- Token-based retrieval system
- Time-to-live (TTL) enforcement
- Access reason logging
- Thread-safe operations
- Automatic expiration handling

## Usage

See the [example](example/main.go) for usage patterns of the vault implementation.

## Security Notes

1. **Encryption**: Uses AES-GCM for authenticated encryption
2. **Key Management**: Requires a master key for encryption/decryption
3. **Nonce Management**: Automatically generates unique nonces for each encryption operation
4. **Access Control**: Tracks access reasons for audit purposes
5. **Data Expiration**: Automatically expires tokens after TTL

## Testing

The implementation includes comprehensive unit tests:

```bash
go test ./vault/... -v
```

## Implementation Details

The vault implementation includes:

1. **VaultEntry**: Represents a stored token with metadata
2. **Vault**: Main vault interface with thread-safe operations
3. **Encryption**: AES-GCM encryption with automatic nonce generation
4. **Expiration**: Automatic cleanup of expired entries
5. **Access Tracking**: Recording of access reasons for audit purposes

## Supported Operations

- **Store**: Store encrypted data with TTL and access reason
- **Retrieve**: Retrieve and decrypt data with access reason tracking
- **Delete**: Remove entries from the vault
- **ListTokens**: List all token IDs in the vault
- **GetEntryInfo**: Get metadata about an entry without decrypting

## Use Cases

- Credit card tokenization
- PII storage and retrieval
- Authentication token storage
- API key management
- Any sensitive data that needs secure storage with expiration
