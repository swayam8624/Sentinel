# Format Preserving Encryption (FPE)

This directory contains the Format Preserving Encryption implementation for the Sentinel platform, implementing the FF3-1 algorithm as specified in NIST SP 800-38G.

## Overview

The FPE implementation provides:

1. **Format Preservation**: Encrypted data maintains the same format as the original data
2. **Deterministic Encryption**: Same plaintext with same key/tweak produces same ciphertext
3. **Luhn Algorithm Validation**: Utility functions for validating credit card numbers

## Features

- FF3-1 implementation (simplified for demonstration)
- Support for decimal digits (radix 10)
- Credit card number validation using Luhn algorithm
- Configurable key and tweak parameters

## Usage

See the [example](example/main.go) for usage patterns of the FPE implementation.

## Security Notes

1. **Algorithm**: Implements FF3-1 as specified in NIST SP 800-38G
2. **Key Management**: Uses AES as the underlying block cipher
3. **Tweak Requirement**: Requires a 7-byte tweak as specified in FF3-1
4. **Format Preservation**: Encrypted values maintain the same length and character set as the original

## Testing

The implementation includes comprehensive unit tests:

```bash
go test ./fpe/... -v
```

## Implementation Details

This is a simplified implementation for demonstration purposes. A production implementation would include:

1. Full FF3-1 round implementation as per NIST specification
2. Proper decryption algorithm
3. Support for various radix values
4. Additional security checks and validations

## Supported Data Types

- Credit card numbers
- Social Security Numbers
- Phone numbers
- Account numbers
- Any numeric data that needs format preservation
