package unit

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/sentinel-platform/sentinel/proxy"
)

// Mock implementations for testing
type MockCipherMeshProcessor struct{}

func (m *MockCipherMeshProcessor) ProcessRequest(ctx context.Context, req *http.Request) error {
	return nil
}

func (m *MockCipherMeshProcessor) ProcessResponse(ctx context.Context, resp *http.Response) error {
	return nil
}

type MockSentinelProcessor struct{}

func (m *MockSentinelProcessor) ProcessRequest(ctx context.Context, req *http.Request) error {
	return nil
}

func (m *MockSentinelProcessor) ProcessResponse(ctx context.Context, resp *http.Response) error {
	return nil
}

type MockRateLimiter struct{}

func (m *MockRateLimiter) Allow(tenantID string) bool {
	return true
}

func (m *MockRateLimiter) GetRetryAfter(tenantID string) time.Duration {
	return 0
}

// TestReverseProxyCreation tests creating a new reverse proxy
func TestReverseProxyCreation(t *testing.T) {
	targetURL, err := url.Parse("https://api.openai.com")
	if err != nil {
		t.Fatalf("Failed to parse target URL: %v", err)
	}

	cipherMesh := &MockCipherMeshProcessor{}
	sentinel := &MockSentinelProcessor{}
	rateLimiter := &MockRateLimiter{}

	proxy := proxy.NewReverseProxy(targetURL, cipherMesh, sentinel, rateLimiter, 30*time.Second)

	if proxy == nil {
		t.Error("Failed to create reverse proxy")
	}

	if proxy.GetTargetURL().String() != targetURL.String() {
		t.Errorf("Expected target URL %s, got %s", targetURL.String(), proxy.GetTargetURL().String())
	}
}

// TestReverseProxyTargetURLModification tests modifying the target URL
func TestReverseProxyTargetURLModification(t *testing.T) {
	targetURL, err := url.Parse("https://api.openai.com")
	if err != nil {
		t.Fatalf("Failed to parse target URL: %v", err)
	}

	newTargetURL, err := url.Parse("https://api.anthropic.com")
	if err != nil {
		t.Fatalf("Failed to parse new target URL: %v", err)
	}

	cipherMesh := &MockCipherMeshProcessor{}
	sentinel := &MockSentinelProcessor{}
	rateLimiter := &MockRateLimiter{}

	proxy := proxy.NewReverseProxy(targetURL, cipherMesh, sentinel, rateLimiter, 30*time.Second)

	// Check initial URL
	if proxy.GetTargetURL().String() != targetURL.String() {
		t.Errorf("Expected initial target URL %s, got %s", targetURL.String(), proxy.GetTargetURL().String())
	}

	// Modify URL
	proxy.SetTargetURL(newTargetURL)

	// Check modified URL
	if proxy.GetTargetURL().String() != newTargetURL.String() {
		t.Errorf("Expected modified target URL %s, got %s", newTargetURL.String(), proxy.GetTargetURL().String())
	}
}
