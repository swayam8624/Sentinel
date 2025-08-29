# Sentinel + CipherMesh v0.1.0

## Initial Public Release

We're excited to announce the initial public release of Sentinel + CipherMesh, a self-healing LLM firewall with cryptographic data protection!

## What is Sentinel + CipherMesh?

Sentinel is a drop-in gateway/SDK that shields upstream LLM providers from raw sensitive data while adding a self-healing security layer that detects, corrects, or cryptographically contains adversarial prompts.

CipherMesh provides PII/tokenization and cryptographic protection layers.

## Key Features

### 🔐 Data Protection

- Real-time data redaction/tokenization
- Format-preserving encryption (FF3-1)
- Reversible detokenization with policy gating
- Multi-language PII detection

### 🛡️ Self-Healing Security

- Semantic violation detection
- Constitutional AI reflection
- Prompt rewriting and ranking
- Tool/function call guarding

### 🔄 Provider Agnostic

- Adapters for OpenAI, Anthropic, Mistral, Hugging Face, Ollama
- Reverse proxy mode with provider-compatible endpoints
- SDK middleware for Python and Node.js
- Streaming support with mid-stream inspection

### ⚙️ Policy Engine

- OPA-style policy evaluation
- Policy versioning and canary deployments
- Multi-tenant policy management
- Audit trails and compliance reporting

### 🔐 Cryptographic Security

- BYOK/HSM integration
- Envelope encryption with AES-256-GCM
- HKDF-SHA-512 key derivation
- Tamper-evident audit logs

## Distribution Channels

This release is available through multiple distribution channels:

### Docker

Pre-built Docker images are available on Docker Hub:

```bash
docker pull sentinel/gateway:latest
```

### Helm Charts

Deploy Sentinel on Kubernetes using Helm:

```bash
helm repo add sentinel https://sentinel-platform.github.io/helm-charts
helm install sentinel sentinel/sentinel
```

### Language SDKs

Integrate Sentinel directly into your applications:

**Node.js:**

```bash
npm install @sentinel-platform/sentinel-sdk
```

**Python:**

```bash
pip install sentinel-sdk
```

### From Source

Build from source:

```bash
git clone https://github.com/swayam8624/Sentinel.git
cd Sentinel
make build
```

## Getting Started

For detailed installation and usage instructions, please refer to our [documentation](https://github.com/swayam8624/Sentinel/docs).

## Feedback and Contributions

We welcome feedback, bug reports, and contributions! Please see our [Contributing Guide](https://github.com/swayam8624/Sentinel/CONTRIBUTING.md) for details on how to get involved.
