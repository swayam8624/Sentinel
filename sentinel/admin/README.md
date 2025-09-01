# Admin Console & APIs

The admin component provides web-based and API-based management interfaces for Sentinel + CipherMesh.

## Components

- **Web Console**: UI for policy management, tenant configuration, and monitoring
- **Admin APIs**: RESTful APIs for programmatic administration
- **Audit System**: Tamper-evident logging and reporting
- **Alerting**: Notification system for security events

## Web Console Features

### Policy Management

- Create, edit, and version policies
- Test policies with sample data
- Deploy policies with canary rollout
- View policy evaluation results

### Tenant Management

- Create and configure tenants
- Set data residency preferences
- Manage tenant-specific policies
- Configure rate limits and quotas

### Key Management

- View key status and rotation schedules
- Initiate key rotations
- Configure BYOK settings
- Monitor key usage

### Monitoring & Alerts

- Real-time dashboard of system metrics
- Security event timeline
- Violation detection reports
- Performance metrics

### Audit & Compliance

- Tamper-evident log viewer
- Export audit trails
- DSAR/DSR support tools
- Compliance reporting

## Admin APIs

### Policy APIs

```
GET    /api/v1/policies
POST   /api/v1/policies
GET    /api/v1/policies/{id}
PUT    /api/v1/policies/{id}
DELETE /api/v1/policies/{id}
```

### Tenant APIs

```
GET    /api/v1/tenants
POST   /api/v1/tenants
GET    /api/v1/tenants/{id}
PUT    /api/v1/tenants/{id}
DELETE /api/v1/tenants/{id}
```

### Key Management APIs

```
GET    /api/v1/keys
POST   /api/v1/keys/rotate
GET    /api/v1/keys/{id}
```

### Audit APIs

```
GET    /api/v1/audit/logs
GET    /api/v1/audit/logs/{id}
POST   /api/v1/audit/export
```

## Authentication & Authorization

- **Role-Based Access Control (RBAC)**:

  - Admin: Full access
  - Security Analyst: Policy and audit access
  - Operator: Monitoring and alerts
  - Read-Only: View-only access

- **Authentication Methods**:
  - SSO integration (SAML, OAuth)
  - API keys for programmatic access
  - mTLS for service-to-service communication

## Multi-tenancy

The admin console supports multi-tenancy with:

- **Hard isolation**: Tenant data separated at the database level
- **Scoped views**: Users only see data for their tenant
- **Cross-tenant management**: Admin users can manage multiple tenants
