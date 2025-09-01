"""
Offline example for the Sentinel Python SDK demonstrating usage without a running service
"""

from sentinel import SentinelClient


def main():
    # Initialize the client
    client = SentinelClient(
        base_url="http://localhost:8080",  # This won't be used in this example
        api_key="your-api-key"
    )
    
    print("Sentinel Python SDK - Offline Example")
    print("=" * 40)
    print()
    
    # Show what the client looks like
    print("Client Configuration:")
    print(f"  Base URL: {client.base_url}")
    print(f"  API Key: {'*' * len(client.api_key) if client.api_key else 'None'}")
    print(f"  Timeout: {client.timeout} seconds")
    print()
    
    # Show available methods
    print("Available Methods:")
    print("  - client.sanitize_prompt(prompt)")
    print("  - client.process_response(response)")
    print("  - client.configure_policies(policies)")
    print("  - client.chat_completions.create(model, messages, ...)")
    print()
    
    # Show example of how to use chat completions
    print("Example Usage:")
    print("```python")
    print("from sentinel import SentinelClient")
    print()
    print("client = SentinelClient(")
    print('    base_url="https://your-sentinel-gateway.com",')
    print('    api_key="your-api-key"')
    print(")")
    print()
    print("response = client.chat_completions.create(")
    print('    model="gpt-3.5-turbo",')
    print('    messages=[{"role": "user", "content": "Hello, world!"}]')
    print(")")
    print("```")
    print()
    
    # Show what the response structure looks like
    print("Expected Response Structure:")
    print("{")
    print('  "id": "chatcmpl-123",')
    print('  "object": "chat.completion",')
    print('  "created": 1234567890,')
    print('  "model": "gpt-3.5-turbo",')
    print('  "choices": [{')
    print('    "index": 0,')
    print('    "message": {')
    print('      "role": "assistant",')
    print('      "content": "Hello! How can I help you today?"')
    print('    },')
    print('    "finish_reason": "stop"')
    print('  }],')
    print('  "usage": {')
    print('    "prompt_tokens": 10,')
    print('    "completion_tokens": 10,')
    print('    "total_tokens": 20')
    print('  }')
    print("}")
    print()
    
    print("For full documentation, visit: https://swayam8624.github.io/Sentinel/")


if __name__ == "__main__":
    main()