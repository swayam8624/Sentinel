# Sentinel Tutorials

This directory contains comprehensive tutorials for using Sentinel with different programming languages and frameworks.

## Getting Started

Before diving into the tutorials, make sure you have:

1. A running Sentinel gateway instance
2. An API key for authentication
3. The appropriate SDK installed for your language

## Language-Specific Tutorials

### Python

- [Python SDK Tutorial](python-sdk-tutorial.md) - Complete guide to using the Sentinel Python SDK
- Topics covered:
  - Basic client initialization
  - Sending chat completions
  - Data protection features
  - Custom policy configuration
  - Integration with popular frameworks (LangChain, LlamaIndex)
  - Multi-tenant usage
  - Error handling and best practices

### Node.js

- [Node.js SDK Tutorial](nodejs-sdk-tutorial.md) - Complete guide to using the Sentinel Node.js SDK
- Topics covered:
  - Basic client initialization
  - Sending chat completions
  - Data protection features
  - Custom policy configuration
  - Integration with Express.js
  - Multi-tenant usage
  - Error handling and best practices

## Advanced Topics

### Security Features

- Data redaction and tokenization
- Format-preserving encryption (FF3-1)
- Semantic violation detection
- Constitutional AI reflection
- Prompt rewriting and ranking
- Tool/function call guarding

### Deployment Patterns

- Reverse-proxy mode
- SDK middleware integration
- Sidecar/gateway deployment
- Multi-tenant configurations
- Kubernetes deployments with Helm

### Policy Management

- OPA-style policy evaluation
- Policy versioning and canary deployments
- Multi-tenant policy management
- Audit trails and compliance reporting

## Examples Directory

For runnable examples, check out:

- [Python Examples](../../sdk/python/examples/)
- [Node.js Examples](../../sdk/nodejs/examples/) (when available)

## Contributing

If you'd like to contribute additional tutorials or improve existing ones:

1. Fork the repository
2. Create a new branch for your tutorial
3. Add your tutorial following the existing format
4. Submit a pull request

## Support

For issues, feature requests, or questions about these tutorials, please [open an issue](https://github.com/swayam8624/Sentinel/issues) on GitHub.
