# Sentinel Node.js SDK

The official Node.js SDK for Sentinel - a self-healing LLM firewall with cryptographic data protection.

## Installation

```bash
npm install @sentinel/sdk
```

## Usage

### Basic Usage

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

### Streaming

```javascript
// Streaming example
for await (const chunk of sentinel.chatStream({
  model: "gpt-4",
  messages: [{ role: "user", content: "Tell me a story" }],
})) {
  process.stdout.write(chunk.choices[0].delta.content || "");
}
```

## Error Handling

```javascript
const { Sentinel, SecurityError, PolicyError } = require("@sentinel/sdk");

try {
  const response = await sentinel.chat({
    model: "gpt-4",
    messages: [{ role: "user", content: "Sensitive content" }],
  });
} catch (error) {
  if (error instanceof SecurityError) {
    console.log(`Security violation: ${error.message}`);
  } else if (error instanceof PolicyError) {
    console.log(`Policy violation: ${error.message}`);
  } else {
    console.log(`Other error: ${error.message}`);
  }
}
```

## Configuration

The SDK can be configured with:

- `endpoint`: The Sentinel gateway URL
- `apiKey`: API key for authentication
- `timeout`: Request timeout in milliseconds (default: 30000)

## Features

- Automatic redaction of sensitive data
- Policy-based handling of different data classes
- Security violation detection and response
- Tool/function call guarding
- Streaming support
- Comprehensive error handling

## License

Apache 2.0
