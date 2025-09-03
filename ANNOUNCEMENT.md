# 🎉 Sentinel + CipherMesh v0.1.0 - Initial Public Release! 🎉

We're excited to announce the initial public release of **Sentinel + CipherMesh**, a self-healing LLM firewall with cryptographic data protection!

## What is Sentinel + CipherMesh?

Sentinel is a drop-in gateway/SDK that shields upstream LLM providers from raw sensitive data while adding a self-healing security layer that detects, corrects, or cryptographically contains adversarial prompts.

CipherMesh provides PII/tokenization and cryptographic protection layers.

## 🔥 Key Features

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

## 🚀 Getting Started

### Docker

```bash
docker run -p 8080:8080 sentinel/gateway:0.1.0
```

### Helm (Kubernetes)

```bash
helm repo add sentinel https://swayam8624.github.io/Sentinel/charts
helm install sentinel sentinel/sentinel
```

### SDKs

**Node.js:**

```bash
npm install @yugenkairo/sentinel-sdk
```

**Python:**

```bash
pip install sentinel-sdk
```

**Rust:**

```bash
cargo add sentinel-sdk
```

**Java:**

```xml
<dependency>
    <groupId>com.sentinel</groupId>
    <artifactId>sentinel-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

**Go:**

```bash
go get github.com/swayam8624/Sentinel/sdk/go
```

## 📚 Documentation

Complete documentation is available at: https://swayam8624.github.io/Sentinel/

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](https://github.com/swayam8624/Sentinel/blob/main/CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](https://github.com/swayam8624/Sentinel/blob/main/LICENSE) file for details.

---

**Repository:** https://github.com/swayam8624/Sentinel  
**Documentation:** https://swayam8624.github.io/Sentinel/  
**Docker Images:** sentinel/gateway on Docker Hub  
**Helm Charts:** https://swayam8624.github.io/Sentinel/charts  
**Node.js SDK:** @yugenkairo/sentinel-sdk on npm  
**Python SDK:** sentinel-sdk on PyPI  
**Rust SDK:** sentinel-sdk on crates.io  
**Java SDK:** com.sentinel:sentinel-sdk on Maven Central  
**Go SDK:** github.com/swayam8624/Sentinel/sdk/go

Join us in building the future of secure LLM applications! 🚀
