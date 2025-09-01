# Sentinel SDK

The Sentinel SDK provides language-specific libraries for integrating Sentinel protection directly into applications.

## Supported Languages

- Python
- Node.js
- Java (planned)
- Go (planned)

## Installation

### Python

```bash
pip install sentinel-sdk
```

### Node.js

```bash
npm install @sentinel/sdk
```

## Usage

### Python Example

```python
from sentinel import Sentinel

# Initialize the client
sentinel = Sentinel(
    endpoint="https://sentinel.example.com",
    api_key="your-api-key"
)

# Send a chat request through Sentinel
response = sentinel.chat(
    model="gpt-4",
    messages=[
        {"role": "user", "content": "Hello, how are you?"}
    ],
    metadata={
        "data_classes": ["pii"],
        "tools_allowed": True
    }
)

print(response.choices[0].message.content)
```

### Node.js Example

```javascript
const { Sentinel } = require("@sentinel/sdk");

// Initialize the client
const sentinel = new Sentinel({
  endpoint: "https://sentinel.example.com",
  apiKey: "your-api-key",
});

// Send a chat request through Sentinel
const response = await sentinel.chat({
  model: "gpt-4",
  messages: [{ role: "user", content: "Hello, how are you?" }],
  metadata: {
    dataClasses: ["pii"],
    toolsAllowed: true,
  },
});

console.log(response.choices[0].message.content);
```

## Configuration

The SDK can be configured with:

- **Endpoint**: The Sentinel gateway URL
- **Authentication**: API key or other auth methods
- **Default Policies**: Tenant or policy version pins
- **Timeouts**: Request timeout settings
- **Retry Logic**: Retry policies for failed requests

## Features

### Data Protection

- Automatic redaction of sensitive data
- Policy-based handling of different data classes
- Secure transmission to LLM providers

### Security

- Violation detection and response
- Tool/function call guarding
- Encryption of unsafe outputs

### Observability

- Request/response logging (PII-safe)
- Performance metrics
- Error reporting

## Error Handling

The SDK provides comprehensive error handling:

- **SentinelError**: Base error class
- **PolicyError**: Policy evaluation failures
- **SecurityError**: Security violation responses
- **NetworkError**: Connectivity issues
- **ValidationError**: Input validation failures

## Streaming Support

Both streaming and non-streaming modes are supported:

```python
# Streaming example
for chunk in sentinel.chat_stream(
    model="gpt-4",
    messages=[{"role": "user", "content": "Tell me a story"}]
):
    print(chunk.choices[0].delta.content or "", end="")
```

## Middleware Integration

The SDK can be integrated as middleware with popular frameworks:

- **Express.js** (Node.js)
- **Flask/FastAPI** (Python)
- **Spring Boot** (Java - planned)
