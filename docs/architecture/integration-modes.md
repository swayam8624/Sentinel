# Integration Modes and Provider Adapters

## Integration Modes

Sentinel supports three primary integration modes to accommodate different deployment scenarios and requirements.

### 1. Proxy Mode (Reverse Proxy)

**Description**: Sentinel acts as a reverse proxy between client applications and LLM providers.

**Architecture**:

```
[Client Application] → [Sentinel Proxy] → [LLM Provider]
```

**Advantages**:

- Zero code changes required in client applications
- Centralized security and data protection
- Easy deployment and management
- Transparent to client applications
- Supports all LLM provider protocols

**Disadvantages**:

- Additional network hop
- Potential single point of failure
- May require network configuration changes

**Use Cases**:

- Legacy applications where code changes are difficult
- Organizations wanting centralized control
- Rapid prototyping and proof of concepts

**Implementation Details**:

- HTTP/HTTPS proxy server
- Protocol-compatible endpoints for major LLM providers
- Load balancing and failover capabilities
- Request/response interception and processing

### 2. SDK Mode (Library Integration)

**Description**: Sentinel functionality is integrated directly into applications via language-specific SDKs.

**Architecture**:

```
[Client Application + Sentinel SDK] → [LLM Provider]
```

**Advantages**:

- No additional network hops
- Fine-grained control over integration points
- Language-specific optimizations
- Better performance for latency-sensitive applications

**Disadvantages**:

- Requires code changes in client applications
- SDK maintenance across multiple languages
- Version compatibility management
- Distributed security controls

**Use Cases**:

- New application development
- Performance-critical applications
- Microservices architectures
- Applications with complex integration requirements

**Implementation Details**:

- Language-specific SDKs (Python, Node.js, Java, Go)
- Middleware wrappers for popular frameworks
- Configuration-driven behavior
- Local caching where appropriate

### 3. Sidecar Mode (Microservices)

**Description**: Sentinel runs as a sidecar container alongside application containers in a microservices architecture.

**Architecture**:

```
[App Container] ↔ [Sentinel Sidecar] → [LLM Provider]
```

**Advantages**:

- Isolation of concerns
- Independent scaling
- Kubernetes-native deployment
- Shared infrastructure between services

**Disadvantages**:

- Container orchestration dependency
- Resource overhead per service
- Network complexity
- Debugging challenges

**Use Cases**:

- Kubernetes-based microservices
- Service mesh environments
- Containerized applications
- Cloud-native architectures

**Implementation Details**:

- Lightweight container image
- gRPC/HTTP communication with main container
- Shared volume for configuration
- Health checks and monitoring integration

## Provider Adapters

Sentinel provides adapters for major LLM providers to ensure compatibility and optimal integration.

### 1. OpenAI Adapter

**Supported Endpoints**:

- `/v1/chat/completions`
- `/v1/completions`
- `/v1/embeddings`
- `/v1/models`

**Features**:

- Streaming support
- Function calling integration
- Model parameter mapping
- Rate limit handling

**Implementation**:

- Protocol-compatible request/response handling
- Streaming chunk processing
- Error mapping and handling
- Model-specific optimizations

### 2. Anthropic Adapter

**Supported Endpoints**:

- `/v1/complete`
- `/v1/messages`

**Features**:

- Claude-specific parameter handling
- Content filtering integration
- Custom stop sequences
- Model version management

**Implementation**:

- Anthropic API compatibility
- Prompt formatting for Claude models
- Response parsing and validation
- Error handling for Anthropic-specific errors

### 3. Mistral Adapter

**Supported Endpoints**:

- `/v1/chat/completions`
- `/v1/completions`
- `/v1/embeddings`

**Features**:

- Mistral model parameter support
- Function calling capabilities
- Streaming response handling
- Model-specific optimizations

**Implementation**:

- API endpoint compatibility
- Model parameter mapping
- Response format standardization
- Error handling and retry logic

### 4. Hugging Face Adapter

**Supported Endpoints**:

- `/api/tasks/*`
- `/api/models/*`
- Inference API endpoints

**Features**:

- Model hub integration
- Task-specific processing (text-generation, question-answering, etc.)
- Custom model support
- Pipeline integration

**Implementation**:

- Hugging Face API compatibility
- Model loading and caching
- Task-specific request/response handling
- Resource management for model loading

### 5. Ollama Adapter

**Supported Endpoints**:

- `/api/generate`
- `/api/chat`
- `/api/embeddings`

**Features**:

- Local model support
- Model management commands
- Streaming responses
- Custom model integration

**Implementation**:

- Ollama API compatibility
- Local model discovery and management
- Resource monitoring for local inference
- Error handling for local processing issues

## Adapter Architecture

### Common Interface

All adapters implement a common interface:

```go
type LLMAdapter interface {
    // Send a chat completion request
    ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

    // Send a streaming chat completion request
    ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)

    // Get model information
    GetModelInfo(ctx context.Context, modelID string) (*ModelInfo, error)

    // Validate adapter configuration
    ValidateConfig(config *AdapterConfig) error

    // Get adapter capabilities
    GetCapabilities() *AdapterCapabilities
}
```

### Request/Response Models

Standardized request/response models ensure consistency across adapters:

```go
type ChatCompletionRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    Temperature float32   `json:"temperature,omitempty"`
    MaxTokens   int       `json:"max_tokens,omitempty"`
    Stream      bool      `json:"stream,omitempty"`
    // ... other common fields
}

type ChatCompletionResponse struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage"`
}
```

### Streaming Interface

Streaming responses are handled through a common interface:

```go
type ChatCompletionStream interface {
    // Receive the next chunk
    Recv() (*ChatCompletionStreamResponse, error)

    // Close the stream
    Close() error
}
```

## Configuration Management

Each adapter supports configuration through a standardized approach:

```yaml
adapters:
  openai:
    enabled: true
    api_key: "${OPENAI_API_KEY}"
    base_url: "https://api.openai.com/v1"
    timeout: "30s"
    rate_limit:
      requests_per_minute: 60
      burst_limit: 10

  anthropic:
    enabled: true
    api_key: "${ANTHROPIC_API_KEY}"
    base_url: "https://api.anthropic.com/v1"
    timeout: "30s"

  mistral:
    enabled: true
    api_key: "${MISTRAL_API_KEY}"
    base_url: "https://api.mistral.ai/v1"
    timeout: "30s"

  huggingface:
    enabled: true
    api_key: "${HUGGINGFACE_API_KEY}"
    base_url: "https://api-inference.huggingface.co"
    timeout: "60s"

  ollama:
    enabled: true
    base_url: "http://localhost:11434"
    timeout: "120s"
```

## Error Handling and Retry Logic

Adapters implement consistent error handling:

1. **Retry Logic**:

   - Exponential backoff for rate limits
   - Circuit breaker pattern for persistent failures
   - Configurable retry attempts and timeouts

2. **Error Mapping**:

   - Standardized error types across adapters
   - Provider-specific error code mapping
   - Contextual error information

3. **Fallback Mechanisms**:
   - Model fallback within provider
   - Provider fallback when configured
   - Graceful degradation

## Performance Considerations

1. **Connection Pooling**:

   - HTTP connection reuse
   - gRPC connection management
   - Resource cleanup

2. **Caching**:

   - Model information caching
   - Rate limit state tracking
   - Response caching where appropriate

3. **Resource Management**:
   - Memory-efficient streaming
   - Concurrent request handling
   - Resource cleanup and monitoring

## Security Considerations

1. **Authentication**:

   - Secure storage of API keys
   - Token rotation support
   - Credential isolation per tenant

2. **Data Protection**:

   - TLS encryption for all communications
   - Request/response inspection
   - Secure handling of sensitive data

3. **Rate Limiting**:
   - Provider rate limit compliance
   - Fair usage policies
   - Quota management

## Monitoring and Observability

1. **Metrics**:

   - Request latency per adapter
   - Error rates and types
   - Rate limit events
   - Resource utilization

2. **Tracing**:

   - Distributed tracing context propagation
   - Adapter-specific span attributes
   - Performance bottleneck identification

3. **Logging**:
   - Structured logging with context
   - PII-safe log content
   - Debug and audit log levels
