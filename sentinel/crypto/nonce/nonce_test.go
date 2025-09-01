package nonce

import (
	"testing"
	"time"
)

func TestNonceManagerGenerateNonce(t *testing.T) {
	nm := NewNonceManager(10 * time.Second)
	
	nonce, err := nm.GenerateNonce(12)
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}
	
	if len(nonce) != 12 {
		t.Errorf("Expected nonce length 12, got %d", len(nonce))
	}
}

func TestNonceManagerGenerateNonceInvalidSize(t *testing.T) {
	nm := NewNonceManager(10 * time.Second)
	
	_, err := nm.GenerateNonce(0)
	if err == nil {
		t.Error("Expected error for zero size nonce")
	}
	
	_, err = nm.GenerateNonce(-1)
	if err == nil {
		t.Error("Expected error for negative size nonce")
	}
}

func TestNonceManagerIsUnique(t *testing.T) {
	nm := NewNonceManager(10 * time.Second)
	
	nonce, err := nm.GenerateNonce(12)
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}
	
	// First use should be unique
	isUnique := nm.IsUnique(nonce)
	if !isUnique {
		t.Error("First nonce use should be unique")
	}
	
	// Second use should not be unique
	isUnique = nm.IsUnique(nonce)
	if isUnique {
		t.Error("Second nonce use should not be unique")
	}
}

func TestNonceManagerConcurrency(t *testing.T) {
	nm := NewNonceManager(10 * time.Second)
	
	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			nonce, err := nm.GenerateNonce(12)
			if err != nil {
				t.Errorf("Failed to generate nonce: %v", err)
			}
			
			isUnique := nm.IsUnique(nonce)
			if !isUnique {
				t.Error("Nonce should be unique")
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}