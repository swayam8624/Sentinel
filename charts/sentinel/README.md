# Sentinel Helm Chart

This Helm chart deploys Sentinel, a self-healing LLM firewall with cryptographic data protection.

## Prerequisites

- Kubernetes 1.16+
- Helm 3.0+

## Installing the Chart

To install the chart with the release name `sentinel`:

```bash
helm repo add sentinel https://swayam8624.github.io/Sentinel/charts
helm install sentinel sentinel/sentinel
```

## Configuration

The following table lists the configurable parameters of the Sentinel chart and their default values.

| Parameter           | Description             | Default                           |
| ------------------- | ----------------------- | --------------------------------- |
| `tenant.id`         | Tenant ID               | `"default"`                       |
| `tenant.name`       | Tenant name             | `"Default Tenant"`                |
| `image.repository`  | Image repository        | `"swayamsingal/sentinel-gateway"` |
| `image.tag`         | Image tag               | `"latest"`                        |
| `image.pullPolicy`  | Image pull policy       | `"IfNotPresent"`                  |
| `service.type`      | Kubernetes service type | `"ClusterIP"`                     |
| `service.port`      | Service port            | `8080`                            |
| `ingress.enabled`   | Enable ingress          | `false`                           |
| `database.host`     | Database host           | `"postgresql"`                    |
| `database.port`     | Database port           | `5432`                            |
| `database.username` | Database username       | `"sentinel"`                      |
| `database.password` | Database password       | `"sentinel"`                      |
| `database.database` | Database name           | `"sentinel"`                      |
| `redis.host`        | Redis host              | `"redis"`                         |
| `redis.port`        | Redis port              | `6379`                            |
| `kms.provider`      | KMS provider            | `"local"`                         |
| `kms.keyPath`       | KMS key path            | `"/etc/sentinel/keys"`            |

## Example Values

```yaml
# values.yaml
tenant:
  id: "my-tenant"
  name: "My Tenant"

image:
  repository: swayamsingal/sentinel-gateway
  tag: v0.1.0

service:
  type: LoadBalancer
  port: 8080

ingress:
  enabled: true
  hosts:
    - host: sentinel.example.com
      paths:
        - path: /
          pathType: ImplementationSpecific

database:
  host: my-postgresql.example.com
  port: 5432
  username: sentinel
  password: mysecretpassword
  database: sentinel

redis:
  host: my-redis.example.com
  port: 6379

kms:
  provider: aws
  keyPath: arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012
```

## Upgrading

To upgrade the chart:

```bash
helm upgrade sentinel sentinel/sentinel -f values.yaml
```

## Uninstalling

To uninstall/delete the release:

```bash
helm delete sentinel
```
