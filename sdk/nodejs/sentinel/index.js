const axios = require('axios');

class Sentinel {
  /**
   * Initialize the Sentinel client
   * @param {Object} options - Configuration options
   * @param {string} options.endpoint - The Sentinel gateway endpoint
   * @param {string} options.apiKey - The API key for authentication
   * @param {number} options.timeout - Request timeout in milliseconds
   */
  constructor(options = {}) {
    this.endpoint = options.endpoint?.replace(/\/$/, '') || 'http://localhost:8080';
    this.apiKey = options.apiKey;
    this.timeout = options.timeout || 30000;
    
    this.client = axios.create({
      baseURL: this.endpoint,
      timeout: this.timeout,
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
        'Content-Type': 'application/json'
      }
    });
  }

  /**
   * Send a chat completion request through Sentinel
   * @param {Object} options - Request options
   * @param {string} options.model - The model to use
   * @param {Array} options.messages - List of messages
   * @param {boolean} options.stream - Whether to stream the response
   * @param {Object} options.metadata - Additional metadata for processing
   * @returns {Promise<Object>} The chat completion response
   */
  async chat(options = {}) {
    const { model, messages, stream = false, metadata } = options;
    
    const payload = {
      model,
      messages,
      stream
    };
    
    if (metadata) {
      payload.metadata = metadata;
    }
    
    try {
      const response = await this.client.post('/v1/chat/completions', payload);
      return response.data;
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Send a streaming chat completion request
   * @param {Object} options - Request options
   * @param {string} options.model - The model to use
   * @param {Array} options.messages - List of messages
   * @param {Object} options.metadata - Additional metadata for processing
   * @returns {AsyncGenerator} Chunks of the streaming response
   */
  async* chatStream(options = {}) {
    const { model, messages, metadata } = options;
    
    try {
      const response = await this.client.post('/v1/chat/completions', {
        model,
        messages,
        stream: true,
        metadata
      }, {
        responseType: 'stream'
      });
      
      for await (const chunk of this._streamResponse(response.data)) {
        yield chunk;
      }
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Process a streaming response
   * @private
   */
  async* _streamResponse(stream) {
    let buffer = '';
    
    for await (const chunk of stream) {
      buffer += chunk.toString();
      
      const lines = buffer.split('\n');
      buffer = lines.pop() || '';
      
      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.substring(6);
          if (data !== '[DONE]') {
            yield JSON.parse(data);
          }
        }
      }
    }
    
    // Process any remaining data in buffer
    if (buffer.startsWith('data: ')) {
      const data = buffer.substring(6);
      if (data !== '[DONE]') {
        yield JSON.parse(data);
      }
    }
  }

  /**
   * Handle error responses
   * @private
   */
  _handleError(error) {
    if (error.response) {
      const { status, data } = error.response;
      const errorCode = data?.error?.code || 'unknown';
      const errorMessage = data?.error?.message || 'Unknown error';
      
      switch (errorCode) {
        case 'security_violation':
          throw new SecurityError(errorMessage);
        case 'policy_violation':
          throw new PolicyError(errorMessage);
        default:
          throw new SentinelError(`API error (${status}): ${errorMessage}`);
      }
    } else if (error.request) {
      throw new NetworkError('Network error: No response received');
    } else {
      throw new SentinelError(`Request failed: ${error.message}`);
    }
  }
}

class SentinelError extends Error {
  constructor(message) {
    super(message);
    this.name = 'SentinelError';
  }
}

class SecurityError extends SentinelError {
  constructor(message) {
    super(message);
    this.name = 'SecurityError';
  }
}

class PolicyError extends SentinelError {
  constructor(message) {
    super(message);
    this.name = 'PolicyError';
  }
}

class NetworkError extends SentinelError {
  constructor(message) {
    super(message);
    this.name = 'NetworkError';
  }
}

module.exports = {
  Sentinel,
  SentinelError,
  SecurityError,
  PolicyError,
  NetworkError
};