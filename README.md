# Sentinel + CipherMesh

A self-healing LLM firewall with cryptographic data protection.

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/sentinel-platform/sentinel)](https://goreportcard.com/report/github.com/sentinel-platform/sentinel)

## Overview

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

## Components

- **Sentinel Gateway**: Core security pipeline with detection, reflection, rewriting, and encryption capabilities
- **CipherMesh**: Data detection and redaction engine with reversible tokenization
- **Policy Engine**: OPA-style policy evaluation system
- **Crypto Vault**: Key management and encryption services
- **Admin Console**: Policy, key, and tenant management interface

## Quick Start

### Using Docker Compose

1. Clone the repository:

```bash
git clone https://github.com/sentinel-platform/sentinel.git
cd sentinel
```

2. Start the services:

```bash
docker-compose up -d
```

3. The gateway will be available at `http://localhost:8080`

### Using Helm (Kubernetes)

1. Add the Helm repository:

```bash
helm repo add sentinel https://sentinel-platform.github.io/helm-charts
helm repo update
```

2. Install the chart:

```bash
helm install sentinel sentinel/sentinel
```

## Documentation

- [Software Requirements Specification](docs/srs.md)
- [Architecture Decision Records](docs/adr/)
- [API Documentation](docs/api/)
- [Deployment Guide](docs/deployment/)
- [Security Policy](docs/security.md)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Install Go 1.23 or later
2. Clone the repository
3. Install dependencies: `make deps`
4. Build the application: `make build`
5. Run the application: `make run`

## Distribution Channels

Sentinel is available through multiple distribution channels:

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

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
