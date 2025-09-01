# Sentinel Node.js SDK Tutorial

This tutorial will guide you through using the Sentinel Node.js SDK to secure your LLM applications.

## Prerequisites

- Node.js 14 or higher
- A running Sentinel gateway instance
- An API key for authentication

## Installation

Install the Sentinel Node.js SDK using npm:

```bash
npm install @yugenkairo/sentinel-sdk
```

## Basic Usage

### Initializing the Client

```javascript
const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

// Initialize the client
const client = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "your-api-key",
});
```

### Sending Chat Completions

The Sentinel Node.js SDK is designed to be a drop-in replacement for the OpenAI SDK:

```javascript
const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

const client = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "your-api-key",
});

// Send a chat completion request through Sentinel
async function sendChatCompletion() {
  try {
    const response = await client.chatCompletions.create({
      model: "gpt-3.5-turbo",
      messages: [{ role: "user", content: "Hello, world!" }],
    });

    console.log(response.choices[0].message.content);
  } catch (error) {
    console.error("Error:", error);
  }
}

sendChatCompletion();
```

## Advanced Usage

### Data Protection Features

Sentinel automatically protects your data through its security pipeline:

```javascript
const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

const client = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "your-api-key",
});

// Sensitive data is automatically detected and protected
async function processSensitiveData() {
  try {
    const sensitivePrompt = "Process customer SSN: 123-45-6789";

    const response = await client.chatCompletions.create({
      model: "gpt-3.5-turbo",
      messages: [{ role: "user", content: sensitivePrompt }],
    });

    // The response is also scanned for sensitive information
    console.log(response.choices[0].message.content);
  } catch (error) {
    console.error("Error:", error);
  }
}

processSensitiveData();
```

### Custom Policy Configuration

Configure custom security policies for your application:

```javascript
const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

const client = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "your-api-key",
});

// Configure custom policies
async function configurePolicies() {
  try {
    const policies = {
      piiDetection: {
        enabled: true,
        languages: ["en", "es", "fr"],
        action: "tokenize",
      },
      promptFiltering: {
        enabled: true,
        threshold: 0.75,
      },
      responseEncryption: {
        enabled: true,
        algorithm: "AES-256-GCM",
      },
    };

    const result = await client.configurePolicies(policies);
    console.log(`Policies configured: ${result.success}`);
  } catch (error) {
    console.error("Error:", error);
  }
}

configurePolicies();
```

## Integration Examples

### Express.js Integration

```javascript
const express = require("express");
const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

const app = express();
app.use(express.json());

const client = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "your-api-key",
});

app.post("/chat", async (req, res) => {
  try {
    const { messages } = req.body;

    const response = await client.chatCompletions.create({
      model: "gpt-3.5-turbo",
      messages: messages,
    });

    res.json(response);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

app.listen(3000, () => {
  console.log("Server running on port 3000");
});
```

### Multi-tenant Usage

Sentinel supports multi-tenant deployments with isolated security policies:

```javascript
const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

// Client for Tenant A
const tenantAClient = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "tenant-a-key",
});

// Client for Tenant B
const tenantBClient = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "tenant-b-key",
});

// Each tenant can have different policies
async function handleTenantRequests() {
  try {
    const tenantAResponse = await tenantAClient.chatCompletions.create({
      model: "gpt-3.5-turbo",
      messages: [{ role: "user", content: "Hello" }],
    });

    const tenantBResponse = await tenantBClient.chatCompletions.create({
      model: "gpt-3.5-turbo",
      messages: [{ role: "user", content: "Hello" }],
    });

    console.log("Tenant A:", tenantAResponse.choices[0].message.content);
    console.log("Tenant B:", tenantBResponse.choices[0].message.content);
  } catch (error) {
    console.error("Error:", error);
  }
}

handleTenantRequests();
```

## Error Handling

Proper error handling is important for production applications:

```javascript
const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

const client = new SentinelClient({
  baseUrl: "https://your-sentinel-gateway.com",
  apiKey: "your-api-key",
});

async function robustCall(prompt, maxRetries = 3) {
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      const response = await client.chatCompletions.create({
        model: "gpt-3.5-turbo",
        messages: [{ role: "user", content: prompt }],
      });
      return response;
    } catch (error) {
      if (attempt === maxRetries - 1) {
        throw error;
      }
      // Wait before retrying (exponential backoff)
      await new Promise((resolve) =>
        setTimeout(resolve, Math.pow(2, attempt) * 1000)
      );
    }
  }
}

// Usage
robustCall("Hello, world!")
  .then((response) => console.log(response.choices[0].message.content))
  .catch((error) => console.error("Final error:", error));
```

## Best Practices

1. **Always use environment variables** for sensitive configuration:

   ```javascript
   const { SentinelClient } = require("@yugenkairo/sentinel-sdk");

   const client = new SentinelClient({
     baseUrl: process.env.SENTINEL_BASE_URL,
     apiKey: process.env.SENTINEL_API_KEY,
   });
   ```

2. **Handle timeouts appropriately**:

   ```javascript
   const client = new SentinelClient({
     baseUrl: "https://your-sentinel-gateway.com",
     apiKey: "your-api-key",
     timeout: 60000, // 60 seconds
   });
   ```

3. **Implement proper logging**:

   ```javascript
   const client = new SentinelClient({
     baseUrl: "https://your-sentinel-gateway.com",
     apiKey: "your-api-key",
   });

   // Add logging middleware
   client.on("request", (config) => {
     console.log("Sending request:", config.method, config.url);
   });

   client.on("response", (response) => {
     console.log("Received response:", response.status, response.statusText);
   });
   ```

## Next Steps

- Check out the [API Reference](../api/) for detailed documentation
- Learn about [Security Features](../security/)
- Explore [Deployment Options](../deployment/)
- Review [Configuration Guide](../configuration/)

For support, please [open an issue](https://github.com/swayam8624/Sentinel/issues) on GitHub.
