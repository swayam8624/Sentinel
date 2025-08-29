# Deployment Guide

This guide covers deploying Sentinel + CipherMesh in various environments.

## Deployment Options

### 1. Proxy Mode (Recommended for most use cases)

Deploy Sentinel as a reverse proxy between your applications and LLM providers.

### 2. Sidecar Mode

Deploy Sentinel alongside your application containers for tighter integration.

### 3. SDK Mode

Integrate Sentinel directly into your applications using language-specific SDKs.

## System Requirements

### Minimum Requirements

- **CPU**: 2 cores
- **Memory**: 4 GB RAM
- **Storage**: 10 GB available disk space
- **Network**: Access to LLM providers

### Recommended Requirements

- **CPU**: 4+ cores
- **Memory**: 8+ GB RAM
- **Storage**: 50+ GB available disk space
- **Network**: Low-latency access to LLM providers

## Prerequisites

1. **Kubernetes** (1.20+) or **Docker** (20.10+)
2. **Cloud KMS/HSM** access for BYOK (optional but recommended)
3. **PostgreSQL** (12+) for metadata storage
4. **Redis** (6+) for caching
5. **Object Storage** (S3-compatible) for log storage

## Deployment Methods

### Helm Chart (Kubernetes)

1. Add the Sentinel Helm repository:

```bash
helm repo add sentinel https://sentinel.github.io/helm-charts
helm repo update
```

2. Create a values file (`values.yaml`):

```yaml
# Tenant configuration
tenant:
  id: "default"
  name: "Default Tenant"

# KMS configuration
kms:
  provider: "aws" # or "azure", "gcp"
  keyArn: "arn:aws:kms:region:account:key/id"

# Database configuration
database:
  host: "postgresql.example.com"
  port: 5432
  username: "sentinel"
  password: "secret"
  database: "sentinel"

# Redis configuration
redis:
  host: "redis.example.com"
  port: 6379
```

3. Install the chart:

```bash
helm install sentinel sentinel/sentinel -f values.yaml
```

### Docker Compose

1. Create a `docker-compose.yaml`:

```yaml
version: "3.8"
services:
  sentinel:
    image: sentinel/gateway:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://sentinel:secret@postgres:5432/sentinel
      - REDIS_URL=redis://redis:6379
      - KMS_PROVIDER=aws
      - KMS_KEY_ARN=arn:aws:kms:region:account:key/id
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_USER=sentinel
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=sentinel

  redis:
    image: redis:6
```

2. Start the services:

```bash
docker-compose up -d
```

## Configuration

### Environment Variables

| Variable       | Description                              | Required              |
| -------------- | ---------------------------------------- | --------------------- |
| `DATABASE_URL` | PostgreSQL connection string             | Yes                   |
| `REDIS_URL`    | Redis connection string                  | Yes                   |
| `KMS_PROVIDER` | KMS provider (aws, azure, gcp)           | Yes                   |
| `KMS_KEY_ARN`  | KMS key identifier                       | Yes                   |
| `LOG_LEVEL`    | Logging level (debug, info, warn, error) | No                    |
| `PORT`         | HTTP port to listen on                   | No (defaults to 8080) |

### Configuration Files

Sentinel supports YAML configuration files:

```yaml
# config.yaml
server:
  port: 8080
  logLevel: info

database:
  url: postgresql://sentinel:secret@postgres:5432/sentinel

redis:
  url: redis://redis:6379

kms:
  provider: aws
  keyArn: arn:aws:kms:region:account:key/id

ciphermesh:
  detectors:
    languages: [en, es, fr]
    enableOcr: false

sentinel:
  thresholds:
    violationSimilarity: 0.78
    reflectConfidence: 0.65
```

## Multi-Tenancy

To configure multiple tenants:

1. **Database Setup**: Each tenant requires separate database schemas or separate databases
2. **KMS Configuration**: Each tenant should have separate KMS keys
3. **Network Isolation**: Consider network-level isolation for high-security tenants

## Monitoring & Observability

### Metrics

Sentinel exposes Prometheus metrics at `/metrics`

Key metrics to monitor:

- `sentinel_request_duration_seconds`: Request latency
- `sentinel_violation_count`: Security violations detected
- `sentinel_redaction_count`: Data redaction operations
- `sentinel_error_count`: Error occurrences

### Logging

Logs are output in JSON format for easy parsing:

```json
{
  "timestamp": "2023-01-01T00:00:00Z",
  "level": "info",
  "message": "Request processed",
  "tenant_id": "default",
  "request_id": "abc123",
  "duration_ms": 150
}
```

### Tracing

OpenTelemetry tracing is available with the following spans:

- Request processing
- Policy evaluation
- Data redaction
- Security detection
- LLM provider calls

## Security Considerations

### Network Security

- Restrict access to admin endpoints
- Use TLS for all communications
- Implement network policies to limit pod-to-pod communication

### Data Security

- Enable BYOK for all encryption operations
- Regularly rotate encryption keys
- Implement proper backup and recovery procedures

### Access Control

- Use strong authentication for admin interfaces
- Implement role-based access control
- Regularly audit access logs

## Backup & Recovery

### Database Backup

Regularly backup the PostgreSQL database:

```bash
pg_dump -h postgres -U sentinel -d sentinel > sentinel_backup.sql
```

### Key Management

Follow your KMS provider's backup and recovery procedures for encryption keys.

### Configuration Backup

Version control all configuration files and deployment manifests.

## Troubleshooting

### Common Issues

1. **Connection Refused**: Check that all services are running and network policies allow communication

2. **Database Connection Failed**: Verify database credentials and connectivity

3. **KMS Access Denied**: Check IAM permissions for KMS access

4. **High Latency**: Monitor resource usage and consider scaling

### Logs

Check container logs for detailed error information:

```bash
kubectl logs <pod-name>
# or
docker-compose logs sentinel
```

## Upgrading

### Version Compatibility

- Minor versions: Backward compatible, safe to upgrade
- Major versions: May require migration steps

### Upgrade Process

1. Backup database and configuration
2. Review release notes for breaking changes
3. Update deployment manifests
4. Perform rolling upgrade (Kubernetes) or blue-green deployment
5. Verify functionality after upgrade
