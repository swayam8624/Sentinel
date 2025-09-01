# CipherMesh

CipherMesh is the data detection and redaction component of Sentinel. It provides:

1. Detection of sensitive data (PII, PHI, PCI, secrets)
2. Redaction/tokenization of sensitive data
3. Reversible detokenization with policy gating
4. Format-preserving encryption (FPE) for utility preservation

## Components

- **Detectors**: Identify sensitive data using regex, NER, and entropy-based scanners
- **Normalizers**: Canonicalize input text for consistent detection
- **Redactors**: Apply redaction actions (tokenize, FPE, mask, etc.)
- **Detokenizer**: Reverse redaction with policy checks
- **Token Vault**: Secure storage of token mappings

## Architecture

```
[Input]
   ↓
[Normalizer] → Unicode NFKC, zero-width char removal, homoglyph mapping
   ↓
[Detector] → Regex, NER, Secret Scanners
   ↓
[Redactor] → Tokenize/FPE/Mask based on policy
   ↓
[Output]
```

## Configuration

CipherMesh is configured through policies that define:

- Which detectors to enable
- What actions to take for different data classes
- Detokenization permissions
- Language support settings

## Data Classes

- **PII**: Personally Identifiable Information
- **PHI**: Protected Health Information
- **PCI**: Payment Card Information
- **Credentials**: API keys, passwords, tokens
- **Custom**: User-defined sensitive data patterns

## Redaction Actions

1. **Tokenize**: Replace with reversible tokens
2. **FPE**: Format-Preserving Encryption (FF3-1)
3. **Mask**: Replace with fixed characters
4. **Hash**: One-way hash transformation
5. **Drop**: Remove entirely
6. **Allow**: Permit without modification
