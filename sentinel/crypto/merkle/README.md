# Merkle Tree for Audit Logs

This directory contains the Merkle tree implementation for tamper-evident audit logs as specified in the Sentinel SRS.

## Overview

The Merkle tree implementation provides:

1. **Tamper-Evident Logs**: Any modification to audit log entries can be detected through Merkle proofs
2. **Efficient Verification**: Log entries can be verified without requiring the entire log set
3. **Compact Proofs**: Small proof sizes that scale logarithmically with the number of log entries

## Features

- Merkle tree construction from audit log entries
- Root hash calculation
- Merkle proof generation for individual log entries
- Proof verification to detect tampering
- Support for trees with any number of leaf nodes

## Usage

See the [example](example/main.go) for usage patterns of the Merkle tree implementation.

## Security Notes

1. **Hash Function**: Uses SHA-256 for cryptographic security
2. **Tree Structure**: Binary tree with duplicate nodes for odd-numbered levels
3. **Proof Security**: Proofs are cryptographically secure and can detect any modification to the proven log entry or any other log entry in the tree

## Testing

The implementation includes comprehensive unit tests:

```bash
go test ./merkle/... -v
```
