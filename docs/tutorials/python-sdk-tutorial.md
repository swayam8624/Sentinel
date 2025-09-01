# Sentinel Python SDK Tutorial

This tutorial will guide you through using the Sentinel Python SDK to secure your LLM applications.

## Prerequisites

- Python 3.8 or higher
- A running Sentinel gateway instance
- An API key for authentication

## Installation

Install the Sentinel Python SDK using pip:

```bash
pip install sentinel-sdk
```

## Basic Usage

### Initializing the Client

```python
from sentinel import SentinelClient

# Initialize the client
client = SentinelClient(
    base_url="https://your-sentinel-gateway.com",
    api_key="your-api-key"
)
```

### Sending Chat Completions

The Sentinel Python SDK is designed to be a drop-in replacement for the OpenAI SDK:

```python
from sentinel import SentinelClient

client = SentinelClient(
    base_url="https://your-sentinel-gateway.com",
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

## Advanced Usage

### Data Protection Features

Sentinel automatically protects your data through its security pipeline:

```python
from sentinel import SentinelClient

client = SentinelClient(
    base_url="https://your-sentinel-gateway.com",
    api_key="your-api-key"
)

# Sensitive data is automatically detected and protected
sensitive_prompt = "Process customer SSN: 123-45-6789"
response = client.chat_completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": sensitive_prompt}
    ]
)

# The response is also scanned for sensitive information
print(response.choices[0].message.content)
```

### Custom Policy Configuration

Configure custom security policies for your application:

```python
from sentinel import SentinelClient

client = SentinelClient(
    base_url="https://your-sentinel-gateway.com",
    api_key="your-api-key"
)

# Configure custom policies
policies = {
    "pii_detection": {
        "enabled": True,
        "languages": ["en", "es", "fr"],
        "action": "tokenize"
    },
    "prompt_filtering": {
        "enabled": True,
        "threshold": 0.75
    },
    "response_encryption": {
        "enabled": True,
        "algorithm": "AES-256-GCM"
    }
}

result = client.configure_policies(policies)
print(f"Policies configured: {result['success']}")
```

## Integration Examples

### LangChain Integration

```python
from langchain.llms import Sentinel
from langchain.prompts import PromptTemplate
from langchain.chains import LLMChain

# Initialize Sentinel LLM
llm = Sentinel(
    base_url="https://your-sentinel-gateway.com",
    api_key="your-api-key"
)

# Create a prompt template
template = "What is {subject}?"
prompt = PromptTemplate.from_template(template)

# Create a chain
chain = LLMChain(llm=llm, prompt=prompt)

# Run the chain
response = chain.run(subject="artificial intelligence")
print(response)
```

### LlamaIndex Integration

```python
from llama_index.llms import Sentinel
from llama_index import VectorStoreIndex, SimpleDirectoryReader

# Initialize Sentinel LLM
llm = Sentinel(
    base_url="https://your-sentinel-gateway.com",
    api_key="your-api-key"
)

# Load documents
documents = SimpleDirectoryReader("data").load_data()

# Create index with Sentinel LLM
index = VectorStoreIndex.from_documents(documents, llm=llm)

# Create query engine
query_engine = index.as_query_engine()

# Query the index
response = query_engine.query("What did the author do growing up?")
print(response)
```

## Multi-tenant Usage

Sentinel supports multi-tenant deployments with isolated security policies:

```python
from sentinel import SentinelClient

# Client for Tenant A
tenant_a_client = SentinelClient(
    base_url="https://your-sentinel-gateway.com",
    api_key="tenant-a-key"
)

# Client for Tenant B
tenant_b_client = SentinelClient(
    base_url="https://your-sentinel-gateway.com",
    api_key="tenant-b-key"
)

# Each tenant can have different policies
tenant_a_response = tenant_a_client.chat_completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello"}]
)

tenant_b_response = tenant_b_client.chat_completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello"}]
)
```

## Error Handling

Proper error handling is important for production applications:

```python
from sentinel import SentinelClient
import requests

client = SentinelClient(
    base_url="https://your-sentinel-gateway.com",
    api_key="your-api-key"
)

try:
    response = client.chat_completions.create(
        model="gpt-3.5-turbo",
        messages=[{"role": "user", "content": "Hello"}]
    )
    print(response.choices[0].message.content)
except requests.exceptions.RequestException as e:
    print(f"Network error: {e}")
except Exception as e:
    print(f"Unexpected error: {e}")
```

## Best Practices

1. **Always use environment variables** for sensitive configuration:

   ```python
   import os
   from sentinel import SentinelClient

   client = SentinelClient(
       base_url=os.getenv("SENTINEL_BASE_URL"),
       api_key=os.getenv("SENTINEL_API_KEY")
   )
   ```

2. **Handle timeouts appropriately**:

   ```python
   client = SentinelClient(
       base_url="https://your-sentinel-gateway.com",
       api_key="your-api-key",
       timeout=60  # 60 seconds
   )
   ```

3. **Implement retry logic** for production applications:

   ```python
   import time
   from sentinel import SentinelClient

   def robust_call(prompt, max_retries=3):
       client = SentinelClient(
           base_url="https://your-sentinel-gateway.com",
           api_key="your-api-key"
       )

       for attempt in range(max_retries):
           try:
               response = client.chat_completions.create(
                   model="gpt-3.5-turbo",
                   messages=[{"role": "user", "content": prompt}]
               )
               return response
           except Exception as e:
               if attempt == max_retries - 1:
                   raise e
               time.sleep(2 ** attempt)  # Exponential backoff
   ```

## Next Steps

- Check out the [API Reference](../api/) for detailed documentation
- Learn about [Security Features](../security/)
- Explore [Deployment Options](../deployment/)
- Review [Configuration Guide](../configuration/)

For support, please [open an issue](https://github.com/swayam8624/Sentinel/issues) on GitHub.
