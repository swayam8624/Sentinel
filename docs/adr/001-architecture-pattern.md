# Architecture Pattern for Sentinel + CipherMesh

## Status

Accepted

## Context

We need to decide on the core architecture pattern for the Sentinel + CipherMesh system. The system must:

1. Act as a gateway/proxy between applications and LLM providers
2. Provide data redaction and tokenization capabilities (CipherMesh)
3. Implement security detection and response mechanisms (Sentinel)
4. Support multiple integration modes (proxy, SDK, sidecar)
5. Be provider-agnostic with adapter support
6. Maintain low latency while providing strong security guarantees

We evaluated several architectural patterns:

- Monolithic architecture
- Microservices architecture
- Plugin/extension architecture
- Layered architecture with modular components

## Decision

We will adopt a **layered architecture with modular components** that can be deployed in different modes:

1. **Gateway/Proxy Layer**: Acts as a reverse proxy between clients and LLM providers
2. **CipherMesh Layer**: Handles data detection, redaction, and tokenization
3. **Sentinel Layer**: Provides security detection, reflection, and response capabilities
4. **Policy Engine Layer**: Evaluates OPA-style policies for data handling
5. **Crypto/Vault Layer**: Manages encryption, key management, and secure storage
6. **Observability Layer**: Provides logging, metrics, and monitoring

This architecture supports multiple deployment modes:

- **Proxy Mode**: Full gateway deployment
- **Sidecar Mode**: Lightweight deployment alongside services
- **SDK Mode**: Library integration directly in applications

Each layer will be implemented as modular components with well-defined interfaces, allowing for independent development, testing, and deployment.

## Consequences

### Positive Consequences

- Clear separation of concerns makes the system easier to understand and maintain
- Modular design allows for independent development and testing of components
- Supports multiple deployment modes without duplicating core logic
- Provider-agnostic design through adapter pattern
- Enables gradual rollout and feature development
- Facilitates horizontal scaling of individual components

### Negative Consequences

- Initial complexity in defining clean interfaces between layers
- Potential for increased latency due to multiple processing layers
- Need for careful dependency management between modules
- May require more initial setup compared to a monolithic approach
