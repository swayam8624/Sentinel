"""
Basic usage example for the Sentinel Python SDK
"""

from sentinel import SentinelClient


def main():
    # Initialize the client
    client = SentinelClient(
        base_url="http://localhost:8080",
        api_key="your-api-key"
    )
    
    # Example 1: Sanitize a prompt
    print("=== Prompt Sanitization Example ===")
    prompt = "Process sensitive data: 123-45-6789"
    sanitized = client.sanitize_prompt(prompt)
    print(f"Original prompt: {prompt}")
    print(f"Sanitized prompt: {sanitized['sanitizedPrompt']}")
    print()
    
    # Example 2: Process an LLM response
    print("=== Response Processing Example ===")
    response = "Here's the sensitive information: 123-45-6789"
    processed = client.process_response(response)
    print(f"Original response: {response}")
    print(f"Processed response: {processed['processedResponse']}")
    print()
    
    # Example 3: Chat completion (this would normally connect to the Sentinel gateway)
    print("=== Chat Completion Example ===")
    try:
        chat_response = client.chat_completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "user", "content": "Hello, world!"}
            ],
            temperature=0.7
        )
        print("Chat completion request would be sent to Sentinel gateway")
        print(f"Model: {chat_response.get('model', 'N/A')}")
        if 'choices' in chat_response:
            print(f"Response: {chat_response['choices'][0]['message']['content']}")
    except Exception as e:
        print(f"Note: This example requires a running Sentinel gateway - {type(e).__name__}")
    print()
    
    # Example 4: Configure policies
    print("=== Policy Configuration Example ===")
    policies = {
        "pii_detection": {
            "enabled": True,
            "languages": ["en", "es", "fr"],
            "action": "tokenize"
        },
        "prompt_filtering": {
            "enabled": True,
            "threshold": 0.75
        }
    }
    
    try:
        result = client.configure_policies(policies)
        print("Policies would be configured on the Sentinel gateway")
        print(f"Success: {result.get('success', 'N/A')}")
    except Exception as e:
        print(f"Note: This example requires a running Sentinel gateway - {type(e).__name__}")


if __name__ == "__main__":
    main()