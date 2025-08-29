package performance

import (
	"testing"
	"time"

	"github.com/sentinel-platform/sentinel/adapters/openai"
)

// TestAdapterCreationPerformance tests the performance of adapter creation
func TestAdapterCreationPerformance(t *testing.T) {
	// Measure the time it takes to create adapters
	start := time.Now()

	for i := 0; i < 1000; i++ {
		_ = openai.NewOpenAIAdapter("test-key", "https://api.openai.com/v1", 30*time.Second)
	}

	elapsed := time.Since(start)

	// This should be fast - less than 100ms for 1000 creations
	if elapsed > 100*time.Millisecond {
		t.Errorf("Adapter creation took too long: %v", elapsed)
	}

	t.Logf("Created 1000 adapters in %v", elapsed)
}

// TestConcurrentAdapterUsage tests concurrent usage of adapters
func TestConcurrentAdapterUsage(t *testing.T) {
	adapter := openai.NewOpenAIAdapter("test-key", "https://api.openai.com/v1", 30*time.Second)

	// Test concurrent capability checks
	start := time.Now()

	// Run 100 concurrent capability checks
	numGoroutines := 100
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			// This is a simple operation that should be fast
			_ = adapter.GetCapabilities()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	elapsed := time.Since(start)

	// This should complete quickly
	if elapsed > 1*time.Second {
		t.Errorf("Concurrent operations took too long: %v", elapsed)
	}

	t.Logf("Completed %d concurrent operations in %v", numGoroutines, elapsed)
}

// BenchmarkAdapterCreation benchmarks adapter creation performance
func BenchmarkAdapterCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = openai.NewOpenAIAdapter("test-key", "https://api.openai.com/v1", 30*time.Second)
	}
}

// BenchmarkCapabilityCheck benchmarks capability check performance
func BenchmarkCapabilityCheck(b *testing.B) {
	adapter := openai.NewOpenAIAdapter("test-key", "https://api.openai.com/v1", 30*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = adapter.GetCapabilities()
	}
}

// TestTimeoutHandling tests timeout handling
func TestTimeoutHandling(t *testing.T) {
	// Create adapter with very short timeout
	adapter := openai.NewOpenAIAdapter("test-key", "https://api.openai.com/v1", 1*time.Millisecond)

	// Measure how quickly validation completes
	start := time.Now()

	err := adapter.ValidateConfig()

	elapsed := time.Since(start)

	// Validation should be fast regardless of timeout setting
	if elapsed > 100*time.Millisecond {
		t.Errorf("Validation took too long: %v", elapsed)
	}

	// Validation should still work with short timeout
	if err != nil {
		t.Errorf("Validation failed unexpectedly: %v", err)
	}
}
