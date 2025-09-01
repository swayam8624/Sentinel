package redaction

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// TokenVault stores mappings between tokens and original values
type TokenVault struct {
	// In-memory storage for demonstration
	// In a real implementation, this would be a database
	storage map[string]*VaultEntry
	mutex   sync.RWMutex
}

// VaultEntry represents a stored token mapping
type VaultEntry struct {
	TokenID        string    `json:"token_id"`
	EncryptedValue []byte    `json:"encrypted_value"` // AEAD encrypted original value
	Tweak          []byte    `json:"tweak,omitempty"` // Tweak used for FPE, if applicable
	DataClass      string    `json:"data_class"`      // Type of data (pii, phi, pci, etc.)
	FieldType      string    `json:"field_type"`      // Specific field type (ssn, credit_card, etc.)
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	AccessCount    int       `json:"access_count"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
}

// NewTokenVault creates a new token vault
func NewTokenVault() *TokenVault {
	return &TokenVault{
		storage: make(map[string]*VaultEntry),
	}
}

// Store stores a mapping between a token and its original value
func (tv *TokenVault) Store(tokenID string, originalValue string, dataClass, fieldType string, ttl time.Duration, encryptionKey []byte) error {
	tv.mutex.Lock()
	defer tv.mutex.Unlock()

	// Encrypt the original value
	encryptedValue, err := encryptWithAEAD([]byte(originalValue), encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	// Create vault entry
	entry := &VaultEntry{
		TokenID:        tokenID,
		EncryptedValue: encryptedValue,
		DataClass:      dataClass,
		FieldType:      fieldType,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(ttl),
		AccessCount:    0,
	}

	// Store the entry
	tv.storage[tokenID] = entry

	return nil
}

// Retrieve retrieves the original value for a token
func (tv *TokenVault) Retrieve(tokenID string, encryptionKey []byte) (string, error) {
	tv.mutex.Lock()
	defer tv.mutex.Unlock()

	// Check if token exists
	entry, exists := tv.storage[tokenID]
	if !exists {
		return "", fmt.Errorf("token not found: %s", tokenID)
	}

	// Check if token has expired
	if time.Now().After(entry.ExpiresAt) {
		// Remove expired token
		delete(tv.storage, tokenID)
		return "", fmt.Errorf("token expired: %s", tokenID)
	}

	// Decrypt the value
	decryptedValue, err := decryptWithAEAD(entry.EncryptedValue, encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt value: %w", err)
	}

	// Update access tracking
	entry.AccessCount++
	entry.LastAccessedAt = time.Now()

	return string(decryptedValue), nil
}

// StoreWithTweak stores a mapping with an FPE tweak
func (tv *TokenVault) StoreWithTweak(tokenID string, originalValue string, tweak []byte, dataClass, fieldType string, ttl time.Duration, encryptionKey []byte) error {
	tv.mutex.Lock()
	defer tv.mutex.Unlock()

	// Encrypt the original value
	encryptedValue, err := encryptWithAEAD([]byte(originalValue), encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	// Create vault entry
	entry := &VaultEntry{
		TokenID:        tokenID,
		EncryptedValue: encryptedValue,
		Tweak:          tweak,
		DataClass:      dataClass,
		FieldType:      fieldType,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(ttl),
		AccessCount:    0,
	}

	// Store the entry
	tv.storage[tokenID] = entry

	return nil
}

// GetTweak retrieves the tweak for a token
func (tv *TokenVault) GetTweak(tokenID string) ([]byte, error) {
	tv.mutex.RLock()
	defer tv.mutex.RUnlock()

	// Check if token exists
	entry, exists := tv.storage[tokenID]
	if !exists {
		return nil, fmt.Errorf("token not found: %s", tokenID)
	}

	// Check if token has expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %s", tokenID)
	}

	return entry.Tweak, nil
}

// Delete removes a token from the vault
func (tv *TokenVault) Delete(tokenID string) error {
	tv.mutex.Lock()
	defer tv.mutex.Unlock()

	// Check if token exists
	_, exists := tv.storage[tokenID]
	if !exists {
		return fmt.Errorf("token not found: %s", tokenID)
	}

	// Remove the entry
	delete(tv.storage, tokenID)

	return nil
}

// CleanupExpired removes expired tokens from the vault
func (tv *TokenVault) CleanupExpired() int {
	tv.mutex.Lock()
	defer tv.mutex.Unlock()

	now := time.Now()
	count := 0

	for tokenID, entry := range tv.storage {
		if now.After(entry.ExpiresAt) {
			delete(tv.storage, tokenID)
			count++
		}
	}

	return count
}

// GetStats returns statistics about the vault
func (tv *TokenVault) GetStats() map[string]interface{} {
	tv.mutex.RLock()
	defer tv.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_entries"] = len(tv.storage)

	// Count by data class
	classCounts := make(map[string]int)
	for _, entry := range tv.storage {
		classCounts[entry.DataClass]++
	}
	stats["by_data_class"] = classCounts

	return stats
}

// Export exports all vault entries (for backup/migration)
func (tv *TokenVault) Export() ([]byte, error) {
	tv.mutex.RLock()
	defer tv.mutex.RUnlock()

	// Convert to JSON
	data, err := json.Marshal(tv.storage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vault data: %w", err)
	}

	return data, nil
}

// Import imports vault entries (for restore/migration)
func (tv *TokenVault) Import(data []byte) error {
	tv.mutex.Lock()
	defer tv.mutex.Unlock()

	// Parse JSON
	var entries map[string]*VaultEntry
	err := json.Unmarshal(data, &entries)
	if err != nil {
		return fmt.Errorf("failed to unmarshal vault data: %w", err)
	}

	// Replace storage
	tv.storage = entries

	return nil
}

// encryptWithAEAD encrypts data using AES-GCM
func encryptWithAEAD(plaintext, key []byte) ([]byte, error) {
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Prepend nonce to ciphertext
	result := append(nonce, ciphertext...)

	return result, nil
}

// decryptWithAEAD decrypts data using AES-GCM
func decryptWithAEAD(data, key []byte) ([]byte, error) {
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
