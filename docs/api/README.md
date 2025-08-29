# API Documentation

Sentinel provides RESTful APIs for integration with applications and management of the system.

## API Endpoints

### Gateway Endpoints

These endpoints are compatible with major LLM providers:

#### Chat Completions

```
POST /v1/chat/completions
```

**Headers:**

- `X-Tenant` (required): Tenant identifier
- `X-Policy-Version` (optional): Specific policy version to use
- `Authorization` (required): Bearer token for authentication

**Request Body:**

```json
{
  "model": "gpt-4",
  "messages": [
    {
      "role": "user",
      "content": "Hello, how are you?"
    }
  ],
  "stream": false,
  "metadata": {
    "data_classes": ["pii", "pci"],
    "tools_allowed": true
  }
}
```

**Response (Non-streaming):**

```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-4",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "I'm doing well, thank you for asking!"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}
```

#### Streaming Response

For streaming requests, the response follows the Server-Sent Events format:

```
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"I'm"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4","choices":[{"index":0,"delta":{"content":" doing well, thank you for asking!"},"finish_reason":"stop"}]}

data: [DONE]
```

### Admin Endpoints

#### Policies

```
GET    /sentinel/admin/policies
POST   /sentinel/admin/policies
GET    /sentinel/admin/policies/{id}
PUT    /sentinel/admin/policies/{id}
DELETE /sentinel/admin/policies/{id}
```

#### Tenants

```
GET    /sentinel/admin/tenants
POST   /sentinel/admin/tenants
GET    /sentinel/admin/tenants/{id}
PUT    /sentinel/admin/tenants/{id}
DELETE /sentinel/admin/tenants/{id}
```

#### Keys

```
GET    /sentinel/admin/keys
POST   /sentinel/admin/keys/rotate
GET    /sentinel/admin/keys/{id}
```

#### Logs

```
GET    /sentinel/admin/logs
GET    /sentinel/admin/logs/{id}
POST   /sentinel/admin/logs/export
```

## Authentication

All admin endpoints require authentication via:

1. **Bearer Token Authentication**

   - Header: `Authorization: Bearer <token>`
   - Tokens issued by the admin console

2. **API Keys**
   - Header: `X-API-Key: <key>`
   - For programmatic access

## Error Responses

All API endpoints return standardized error responses:

```json
{
  "error": {
    "code": "policy_violation",
    "message": "Request blocked due to security policy",
    "details": {
      "violation_type": "jailbreak_attempt",
      "score": 0.95
    }
  }
}
```

## Rate Limiting

API requests are subject to rate limiting:

- **Default Limit**: 100 requests per minute per tenant
- **Burst Limit**: 200 requests per minute per tenant
- **Response**: 429 Too Many Requests when limit exceeded

## Versioning

API versions are specified in the URL path:

- `/v1/` - Current stable version
- `/v2/` - Next version (when available)

Backward compatibility is maintained within major versions.
