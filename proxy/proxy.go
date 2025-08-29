package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// ReverseProxy implements a reverse proxy for LLM providers
type ReverseProxy struct {
	proxy       *httputil.ReverseProxy
	targetURL   *url.URL
	cipherMesh  CipherMeshProcessor
	sentinel    SentinelProcessor
	rateLimiter RateLimiter
	timeout     time.Duration
}

// CipherMeshProcessor interface for CipherMesh processing
type CipherMeshProcessor interface {
	// ProcessRequest processes a request with CipherMesh
	ProcessRequest(ctx context.Context, req *http.Request) error
	
	// ProcessResponse processes a response with CipherMesh
	ProcessResponse(ctx context.Context, resp *http.Response) error
}

// SentinelProcessor interface for Sentinel processing
type SentinelProcessor interface {
	// ProcessRequest processes a request with Sentinel
	ProcessRequest(ctx context.Context, req *http.Request) error
	
	// ProcessResponse processes a response with Sentinel
	ProcessResponse(ctx context.Context, resp *http.Response) error
}

// RateLimiter interface for rate limiting
type RateLimiter interface {
	// Allow checks if a request is allowed
	Allow(tenantID string) bool
	
	// GetRetryAfter returns the time to wait before retrying
	GetRetryAfter(tenantID string) time.Duration
}

// NewReverseProxy creates a new reverse proxy
func NewReverseProxy(
	targetURL *url.URL,
	cipherMesh CipherMeshProcessor,
	sentinel SentinelProcessor,
	rateLimiter RateLimiter,
	timeout time.Duration) *ReverseProxy {
	
	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	return &ReverseProxy{
		proxy:       proxy,
		targetURL:   targetURL,
		cipherMesh:  cipherMesh,
		sentinel:    sentinel,
		rateLimiter: rateLimiter,
		timeout:     timeout,
	}
}

// ServeHTTP implements the http.Handler interface
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), rp.timeout)
	defer cancel()
	
	// Extract tenant ID from request
	tenantID := r.Header.Get("X-Tenant")
	if tenantID == "" {
		http.Error(w, "X-Tenant header is required", http.StatusBadRequest)
		return
	}
	
	// Check rate limits
	if !rp.rateLimiter.Allow(tenantID) {
		retryAfter := rp.rateLimiter.GetRetryAfter(tenantID)
		w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	
	// Process request with CipherMesh
	if err := rp.cipherMesh.ProcessRequest(ctx, r); err != nil {
		http.Error(w, fmt.Sprintf("CipherMesh processing failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Process request with Sentinel
	if err := rp.sentinel.ProcessRequest(ctx, r); err != nil {
		http.Error(w, fmt.Sprintf("Sentinel processing failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Modify request to forward to target
	r.URL.Host = rp.targetURL.Host
	r.URL.Scheme = rp.targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = rp.targetURL.Host
	
	// Forward the request
	rp.proxy.ServeHTTP(w, r)
}

// ProcessResponse processes a response
func (rp *ReverseProxy) ProcessResponse(w http.ResponseWriter, r *http.Request, resp *http.Response) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), rp.timeout)
	defer cancel()
	
	// Process response with Sentinel
	if err := rp.sentinel.ProcessResponse(ctx, resp); err != nil {
		return fmt.Errorf("Sentinel processing failed: %w", err)
	}
	
	// Process response with CipherMesh
	if err := rp.cipherMesh.ProcessResponse(ctx, resp); err != nil {
		return fmt.Errorf("CipherMesh processing failed: %w", err)
	}
	
	return nil
}

// GetTargetURL returns the target URL
func (rp *ReverseProxy) GetTargetURL() *url.URL {
	return rp.targetURL
}

// SetTargetURL sets the target URL
func (rp *ReverseProxy) SetTargetURL(targetURL *url.URL) {
	rp.targetURL = targetURL
	rp.proxy = httputil.NewSingleHostReverseProxy(targetURL)
}package proxy
