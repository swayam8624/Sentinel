# Sentinel + CipherMesh

<p align="center">
  <img src="nope-zy.jpg" alt="Sentinel Banner" width="100%">
</p>

<p align="center">
  <strong>A self-healing LLM firewall with cryptographic data protection</strong>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
  <a href="https://goreportcard.com/report/github.com/swayam8624/Sentinel"><img src="https://goreportcard.com/badge/github.com/swayam8624/Sentinel" alt="Go Report Card"></a>
  <a href="https://github.com/swayam8624/Sentinel/releases"><img src="https://img.shields.io/github/v/release/swayam8624/Sentinel" alt="GitHub release"></a>
  <a href="https://swayam8624.github.io/Sentinel/"><img src="https://img.shields.io/badge/docs-latest-brightgreen.svg" alt="Documentation"></a>
  <a href="https://pypi.org/project/yugenkairo-sentinel-sdk/"><img src="https://img.shields.io/pypi/v/yugenkairo-sentinel-sdk.svg" alt="PyPI"></a>
  <a href="https://www.npmjs.com/package/@yugenkairo/sentinel-sdk"><img src="https://img.shields.io/npm/v/@yugenkairo/sentinel-sdk.svg" alt="npm"></a>
</p>

## 🔐 Sentinel: Enterprise-Grade LLM Security Gateway

Sentinel is a production-ready, drop-in gateway/SDK that shields upstream LLM providers from raw sensitive data while adding a self-healing security layer that detects, corrects, or cryptographically contains adversarial prompts.

CipherMesh provides advanced PII/tokenization and cryptographic protection layers, ensuring your data remains secure while maintaining LLM functionality.

## 🌟 Key Features

### 🔐 Advanced Data Protection

- **Real-time Data Redaction**: Automatic detection and tokenization of sensitive information
- **Format-Preserving Encryption (FF3-1)**: Maintain data format while ensuring security
- **Reversible Detokenization**: Controlled access with policy gating
- **Multi-language PII Detection**: Support for 50+ languages and dialects
- **Code Secret Detection**: Automatic detection of API keys, passwords, and secrets

### 🛡️ Self-Healing Security Pipeline

- **Semantic Violation Detection**: Advanced AI-powered threat detection
- **Constitutional AI Reflection**: Ethical alignment and bias correction
- **Prompt Rewriting & Ranking**: Automatic sanitization of malicious prompts
- **Tool/Function Call Guarding**: Prevent unauthorized system access
- **Adaptive Learning**: Continuous improvement through feedback loops

### 🔌 Universal Provider Compatibility

- **Multi-Provider Adapters**: OpenAI, Anthropic, Mistral, Hugging Face, Ollama, and more
- **Reverse Proxy Mode**: Seamless integration with existing infrastructure
- **Language SDKs**: Native support for Python, Node.js, Java, and Go
- **Streaming Support**: Real-time inspection with zero latency impact
- **Multi-Tenancy**: Secure isolation for enterprise environments

### ⚙️ Advanced Policy Engine

- **OPA-Style Policy Evaluation**: Industry-standard policy management
- **Policy Versioning**: Safe deployment with canary rollouts
- **Multi-Tenant Policy Management**: Granular control for complex organizations
- **Audit Trails**: Comprehensive compliance reporting
- **Dynamic Policy Updates**: Real-time policy changes without downtime

### 🔐 Cryptographic Security

- **BYOK/HSM Integration**: Bring your own keys for maximum security
- **Envelope Encryption**: AES-256-GCM for data at rest and in transit
- **Advanced Key Derivation**: HKDF-SHA-512 for secure key management
- **Tamper-Evident Audit Logs**: Merkle tree-based integrity verification
- **Cloud KMS Integration**: AWS KMS, Azure Key Vault, GCP KMS support
- **Protection Against All Crypto Attacks**: Comprehensive security against side-channel, replay, brute force, MITM, and other attacks

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌────────────────────┐
│   Application   │───▶│   Sentinel       │───▶│   LLM Provider     │
│   (Client)      │    │   Gateway        │    │   (OpenAI, etc.)   │
└─────────────────┘    │  ┌─────────────┐ │    └────────────────────┘
                       │  │  CipherMesh │ │              │
                       │  │  (Data      │ │              ▼
                       │  │  Protection)│ │    ┌────────────────────┐
                       │  └─────────────┘ │    │   Response         │
                       │  ┌─────────────┐ │    │   Processing       │
                       │  │  Policy     │ │    │   & Encryption     │
                       │  │  Engine     │ │    └────────────────────┘
                       │  └─────────────┘ │              │
                       │  ┌─────────────┐ │              ▼
                       │  │  Security   │ │    ┌────────────────────┐
                       │  │  Pipeline   │ │◀───│   Secure           │
                       │  │  (Detection,│ │    │   Response         │
                       │  │  Reflection,│ │    │   Return           │
                       │  │  Rewriting) │ │    └────────────────────┘
                       │  └─────────────┘ │
                       └──────────────────┘
```

## 🚀 Quick Start

### Using Docker Compose

```bash
# Clone the repository
git clone https://github.com/swayam8624/Sentinel.git
cd Sentinel

# Start the services
docker-compose up -d

# Access the gateway at http://localhost:8080
```

### Using Helm (Kubernetes)

```bash
# Add the Helm repository
helm repo add sentinel https://swayam8624.github.io/Sentinel/charts

# Update repository information
helm repo update

# Install the chart
helm install sentinel sentinel/sentinel
```

### Using Language SDKs

**Node.js:**

```bash
npm install @yugenkairo/sentinel-sdk
```

**Python:**

```bash
pip install yugenkairo-sentinel-sdk
```

## 📚 Comprehensive Documentation

- [📘 Software Requirements Specification](docs/srs.md)
- [🏗️ Architecture Decision Records](docs/adr/)
- [🔌 API Documentation](docs/api/)
- [🚀 Deployment Guide](docs/deployment/)
- [🛡️ Security Policy](docs/security.md)
- [🔐 Cryptographic Security](docs/security/crypto-security.md)
- [📦 Distribution Channels](DISTRIBUTION_CHANNELS.md)
- [📋 Threat Modeling](docs/threat-modeling/)
- [🎓 Tutorials](docs/tutorials/)

## 🛠️ Development Setup

1. **Install Prerequisites**:

   - Go 1.23 or later
   - Docker and Docker Compose
   - Helm 3.x (for Kubernetes deployment)

2. **Clone and Setup**:

   ```bash
   git clone https://github.com/swayam8624/Sentinel.git
   cd Sentinel
   make deps
   ```

3. **Build and Run**:

   ```bash
   make build
   make run
   ```

4. **Testing**:
   ```bash
   make test
   make test-integration
   make test-security
   ```

## 📦 Distribution Channels

Sentinel is available through multiple enterprise-grade distribution channels:

### Docker

Pre-built Docker images for seamless deployment:

```bash
docker pull sentinel/gateway:latest
```

### Helm Charts

Production-ready Kubernetes deployments:

```bash
helm repo add sentinel https://swayam8624.github.io/Sentinel/charts
helm install sentinel sentinel/sentinel
```

### Language SDKs

Native integration for your applications:

**Node.js:**

```bash
npm install @yugenkairo/sentinel-sdk
```

**Python:**

```bash
pip install yugenkairo-sentinel-sdk
```

### From Source

Build from the latest source code:

```bash
git clone https://github.com/swayam8624/Sentinel.git
cd Sentinel
make build
```

## 🤝 Contributing

We welcome contributions from the community! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Community Guidelines

- Follow our [Code of Conduct](CODE_OF_CONDUCT.md)
- Check [existing issues](https://github.com/swayam8624/Sentinel/issues) before creating new ones
- Review our [development practices](CONTRIBUTING.md#development-practices)

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For enterprise support, security disclosures, or professional services, please contact our team at [support@sentinel-platform.org](mailto:support@sentinel-platform.org).

---

<p align="center">
  <strong>Built with ❤️ for the AI security community</strong>
</p>
