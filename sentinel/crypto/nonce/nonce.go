package nonce

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// NonceManager manages cryptographically secure nonces with uniqueness enforcement
type NonceManager struct {
	usedNonces map[string]time.Time
	mutex      sync.RWMutex
	ttl        time.Duration
}

// NewNonceManager creates a new NonceManager with the specified TTL
func NewNonceManager(ttl time.Duration) *NonceManager {
	nm := &NonceManager{
		usedNonces: make(map[string]time.Time),
		ttl:        ttl,
	}
	go nm.cleanupExpired()
	return nm
}

// GenerateNonce generates a cryptographically secure nonce
func (nm *NonceManager) GenerateNonce(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("nonce size must be positive")
	}

	nonce := make([]byte, size)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return nonce, nil
}

// IsUnique checks if a nonce is unique and marks it as used
func (nm *NonceManager) IsUnique(nonce []byte) bool {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	nonceStr := string(nonce)
	now := time.Now()

	// Check if nonce was used before TTL expiration
	if usedAt, exists := nm.usedNonces[nonceStr]; exists {
		if now.Sub(usedAt) < nm.ttl {
			return false // Nonce was used recently
		}
	}

	// Mark nonce as used
	nm.usedNonces[nonceStr] = now
	return true
}

// cleanupExpired removes expired nonces from the map
func (nm *NonceManager) cleanupExpired() {
	ticker := time.NewTicker(nm.ttl)
	defer ticker.Stop()

	for range ticker.C {
		nm.mutex.Lock()
		now := time.Now()
		for nonce, usedAt := range nm.usedNonces {
			if now.Sub(usedAt) >= nm.ttl {
				delete(nm.usedNonces, nonce)
			}
		}
		nm.mutex.Unlock()
	}
}
