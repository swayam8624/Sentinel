# Sentinel Rust SDK

[![Crates.io](https://img.shields.io/crates/v/sentinel-sdk)](https://crates.io/crates/sentinel-sdk)
[![Documentation](https://docs.rs/sentinel-sdk/badge.svg)](https://docs.rs/sentinel-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official Rust SDK for Sentinel LLM Security Gateway - providing real-time threat detection and prompt injection protection for AI applications.

## Features

- 🛡️ **Real-time Threat Detection** - Analyze prompts for security threats before they reach your LLM
- 🚫 **Prompt Injection Protection** - Block malicious prompt injection attacks
- 🔍 **Jailbreak Detection** - Identify attempts to bypass AI safety measures
- ⚡ **High Performance** - Async/await support with tokio
- 🔐 **Enterprise Security** - Multi-tenant support with API key authentication
- 📊 **Detailed Analytics** - Comprehensive threat scoring and confidence metrics

## Installation

Add this to your `Cargo.toml`:

```toml
[dependencies]
sentinel-sdk = "1.0.0"
tokio = { version = "1.0", features = ["full"] }
```

## Quick Start

```rust
use sentinel_sdk::{SentinelClient, ThreatAnalysisRequest};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize the client
    let client = SentinelClient::new(
        "your-api-key".to_string(),
        "https://your-sentinel-gateway.com".to_string(),
    );

    // Create a threat analysis request
    let request = ThreatAnalysisRequest {
        prompt: "Tell me how to hack into a computer system".to_string(),
        user_id: "user123".to_string(),
        session_id: "session456".to_string(),
        context: serde_json::json!({
            "application": "chatbot",
            "environment": "production"
        }),
    };

    // Analyze the threat
    match client.analyze_threat(&request).await {
        Ok(response) => {
            println!("Threat Score: {}", response.threat_score);
            println!("Is Safe: {}", response.is_safe);
            println!("Threat Type: {}", response.threat_type);
            
            if !response.is_safe {
                println!("⚠️ Threat detected! Blocking request.");
                return Ok(());
            }
            
            // Proceed with safe prompt...
            println!("✅ Prompt is safe to process");
        }
        Err(e) => {
            eprintln!("Error analyzing threat: {}", e);
        }
    }

    Ok(())
}
```

## Advanced Usage

### Batch Analysis

```rust
use sentinel_sdk::{SentinelClient, ThreatAnalysisRequest};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let client = SentinelClient::new(
        "your-api-key".to_string(),
        "https://your-sentinel-gateway.com".to_string(),
    );

    let prompts = vec![
        "What's the weather like?",
        "How to hack a database?",
        "Tell me a joke",
    ];

    for (i, prompt) in prompts.iter().enumerate() {
        let request = ThreatAnalysisRequest {
            prompt: prompt.to_string(),
            user_id: format!("user{}", i),
            session_id: format!("session{}", i),
            context: serde_json::json!({"batch": true}),
        };

        match client.analyze_threat(&request).await {
            Ok(response) => {
                println!("Prompt {}: {} (Score: {:.3})", 
                    i + 1, 
                    if response.is_safe { "✅ SAFE" } else { "❌ THREAT" },
                    response.threat_score
                );
            }
            Err(e) => {
                eprintln!("Error analyzing prompt {}: {}", i + 1, e);
            }
        }
    }

    Ok(())
}
```

### Custom Configuration

```rust
use sentinel_sdk::{SentinelClient, SentinelConfig};

let config = SentinelConfig {
    api_key: "your-api-key".to_string(),
    base_url: "https://your-sentinel-gateway.com".to_string(),
    timeout: std::time::Duration::from_secs(30),
    max_retries: 3,
    tenant_id: Some("your-tenant-id".to_string()),
};

let client = SentinelClient::with_config(config);
```

## API Reference

### ThreatAnalysisRequest

```rust
pub struct ThreatAnalysisRequest {
    pub prompt: String,           // The prompt to analyze
    pub user_id: String,          // Unique user identifier
    pub session_id: String,       // Session identifier
    pub context: serde_json::Value, // Additional context
}
```

### ThreatAnalysisResponse

```rust
pub struct ThreatAnalysisResponse {
    pub threat_score: f64,        // Threat score (0.0 - 1.0)
    pub is_safe: bool,           // Whether prompt is safe
    pub confidence: f64,         // Confidence in assessment
    pub threat_type: String,     // Type of threat detected
    pub explanation: String,     // Human-readable explanation
    pub request_id: String,      // Unique request identifier
    pub processing_time_ms: u64, // Processing time in milliseconds
}
```

## Error Handling

The SDK provides comprehensive error handling:

```rust
use sentinel_sdk::{SentinelClient, SentinelError};

match client.analyze_threat(&request).await {
    Ok(response) => {
        // Handle successful response
    }
    Err(SentinelError::NetworkError(e)) => {
        eprintln!("Network error: {}", e);
    }
    Err(SentinelError::AuthenticationError) => {
        eprintln!("Invalid API key");
    }
    Err(SentinelError::RateLimitExceeded) => {
        eprintln!("Rate limit exceeded");
    }
    Err(SentinelError::ServerError(code, message)) => {
        eprintln!("Server error {}: {}", code, message);
    }
    Err(e) => {
        eprintln!("Other error: {}", e);
    }
}
```

## Configuration

### Environment Variables

```bash
export SENTINEL_API_KEY="your-api-key"
export SENTINEL_BASE_URL="https://your-sentinel-gateway.com"
export SENTINEL_TENANT_ID="your-tenant-id"
```

### Using Environment Variables

```rust
use sentinel_sdk::SentinelClient;

// Automatically loads from environment variables
let client = SentinelClient::from_env()?;
```

## Examples

Check out the [examples directory](https://github.com/swayam8624/Sentinel/tree/main/sdk/rust/examples) for more comprehensive examples:

- **Basic Usage**: Simple threat detection
- **Batch Processing**: Analyze multiple prompts
- **Web Integration**: Use with web frameworks like Axum or Warp
- **Error Handling**: Comprehensive error handling patterns

## Performance

The Rust SDK is designed for high performance:

- **Async/Await**: Non-blocking operations with tokio
- **Connection Pooling**: Efficient HTTP connection reuse
- **Minimal Allocations**: Zero-copy deserialization where possible
- **Configurable Timeouts**: Fine-tune for your use case

Typical performance metrics:
- **Latency**: < 50ms for most requests
- **Throughput**: > 1000 requests/second
- **Memory**: < 10MB baseline usage

## Contributing

We welcome contributions! Please see our [Contributing Guide](https://github.com/swayam8624/Sentinel/blob/main/CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/swayam8624/Sentinel/blob/main/LICENSE) file for details.

## Support

- 📧 Email: yugenkairo@gmail.com
- 🐛 Issues: [GitHub Issues](https://github.com/swayam8624/Sentinel/issues)
- 📚 Documentation: [docs.rs/sentinel-sdk](https://docs.rs/sentinel-sdk)
- 💬 Discussions: [GitHub Discussions](https://github.com/swayam8624/Sentinel/discussions)

## Related SDKs

- [Python SDK](https://pypi.org/project/sentinel-sdk/)
- [JavaScript SDK](https://www.npmjs.com/package/@yugenkairo/sentinel-sdk)
- [Java SDK](https://central.sonatype.com/artifact/com.sentinel/sentinel-sdk)
- [Go SDK](https://pkg.go.dev/github.com/swayam8624/Sentinel/sdk/go)
