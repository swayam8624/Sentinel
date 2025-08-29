# Sentinel Python SDK

The official Python SDK for Sentinel - a self-healing LLM firewall with cryptographic data protection.

## Installation

```bash
pip install sentinel-sdk
```

## Usage

### Basic Usage

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

### Streaming

```python
# Streaming example
for chunk in sentinel.chat_stream(
    model="gpt-4",
    messages=[{"role": "user", "content": "Tell me a story"}]
):
    print(chunk.choices[0].delta.content or "", end="")
```

## Error Handling

```python
from sentinel import Sentinel, SecurityError, PolicyError

try:
    response = sentinel.chat(
        model="gpt-4",
        messages=[{"role": "user", "content": "Sensitive content"}]
    )
except SecurityError as e:
    print(f"Security violation: {e}")
except PolicyError as e:
    print(f"Policy violation: {e}")
except Exception as e:
    print(f"Other error: {e}")
```

## Configuration

The SDK can be configured with:

- `endpoint`: The Sentinel gateway URL
- `api_key`: API key for authentication
- `timeout`: Request timeout in seconds (default: 30)

## Features

- Automatic redaction of sensitive data
- Policy-based handling of different data classes
- Security violation detection and response
- Tool/function call guarding
- Streaming support
- Comprehensive error handling

## License

Apache 2.0
