# Sentinel Python SDK

Python SDK for Sentinel - A self-healing LLM firewall with cryptographic data protection.

## Installation

```bash
pip install sentinel-sdk
```

## Usage

```python
from sentinel import SentinelClient

# Initialize the client
client = SentinelClient(
    base_url="http://localhost:8080",
    api_key="your-api-key"
)

# Send a chat completion request through Sentinel
response = client.chat_completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello, world!"}
    ]
)

print(response.choices[0].message.content)
```

## Features

- **Security**: All requests are processed through Sentinel's security pipeline
- **Compatibility**: Drop-in replacement for OpenAI SDK
- **Multi-tenant**: Support for tenant isolation
- **Observability**: Built-in metrics and tracing

## Documentation

For full documentation, visit [https://github.com/swayam8624/Sentinel/docs](https://github.com/swayam8624/Sentinel/docs)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
