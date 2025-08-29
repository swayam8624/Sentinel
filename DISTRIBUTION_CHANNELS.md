# Sentinel + CipherMesh Distribution Channels

This document outlines all the distribution channels available for Sentinel + CipherMesh.

## 1. Docker Images

Pre-built Docker images are available on Docker Hub:

```bash
docker pull sentinel/gateway:latest
```

### Supported Tags

- `latest` - Latest stable release
- `vX.Y.Z` - Specific version releases
- `sha-<commit>` - Specific commit builds

### Usage

```bash
# Run Sentinel gateway
docker run -p 8080:8080 sentinel/gateway:latest

# Run with custom configuration
docker run -p 8080:8080 -v /path/to/config.yaml:/config.yaml sentinel/gateway:latest
```

## 2. Helm Charts

Deploy Sentinel on Kubernetes using Helm:

```bash
helm repo add sentinel https://swayam8624.github.io/Sentinel/charts
helm install sentinel sentinel/sentinel
```

### Chart Repository

The Helm charts are hosted on GitHub Pages:
https://swayam8624.github.io/Sentinel/charts

### Configuration

See [Helm Chart README](charts/sentinel/README.md) for detailed configuration options.

## 3. Language SDKs

Integrate Sentinel directly into your applications.

### Node.js SDK

```bash
npm install @sentinel-platform/sentinel-sdk
```

GitHub: [sdk/nodejs](sdk/nodejs)

### Python SDK

```bash
pip install sentinel-sdk
```

GitHub: [sdk/python](sdk/python)

## 4. Source Code

Build from source:

```bash
git clone https://github.com/swayam8624/Sentinel.git
cd Sentinel
make build
```

### Requirements

- Go 1.23 or later
- Docker (for containerization)
- Kubernetes (for Helm deployment)

## 5. GitHub Releases

Download pre-built binaries from GitHub Releases:
https://github.com/swayam8624/Sentinel/releases

## 6. Documentation

Documentation is available on GitHub Pages:
https://swayam8624.github.io/Sentinel

## 7. CI/CD Pipeline

The project uses GitHub Actions for CI/CD:

- Automated testing on every push
- Docker image building and publishing
- Helm chart packaging
- Release creation

See [.github/workflows/ci.yml](.github/workflows/ci.yml) for details.

## 8. Development Setup

1. Clone the repository
2. Install dependencies: `make deps`
3. Build the application: `make build`
4. Run tests: `make test`
5. Run the application: `make run`

## 9. Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.