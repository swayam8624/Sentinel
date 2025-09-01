package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"sync"
	"time"
)

// Vault provides secure storage for encrypted tokens
type Vault struct {
	encryptionKey []byte
	entries       map[string]*VaultEntry
	mutex         sync.RWMutex
}

// VaultEntry represents an entry in the vault
type VaultEntry struct {
	Data      []byte         `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `json:"expires_at"`
	AccessLog []AccessRecord `json:"access_log"`
}

// AccessRecord represents an access record for audit purposes
type AccessRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason"`
}

// NewVault creates a new Vault instance
func NewVault(encryptionKey []byte) *Vault {
	return &Vault{
		encryptionKey: encryptionKey,
		entries:       make(map[string]*VaultEntry),
	}
}

// Store stores data in the vault with a TTL
func (v *Vault) Store(key string, data []byte, ttl time.Duration) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	// Encrypt the data
	encryptedData, err := v.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Create entry
	entry := &VaultEntry{
		Data:      encryptedData,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
		AccessLog: make([]AccessRecord, 0),
	}

	// Store entry
	v.entries[key] = entry

	return nil
}

// Retrieve retrieves data from the vault
func (v *Vault) Retrieve(key, accessReason string) ([]byte, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	// Check if entry exists
	entry, exists := v.entries[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	// Check if entry has expired
	if time.Now().After(entry.ExpiresAt) {
		delete(v.entries, key)
		return nil, fmt.Errorf("entry has expired")
	}

	// Log access
	entry.AccessLog = append(entry.AccessLog, AccessRecord{
		Timestamp: time.Now(),
		Reason:    accessReason,
	})

	// Decrypt the data
	decryptedData, err := v.decrypt(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decryptedData, nil
}

// Delete removes an entry from the vault
func (v *Vault) Delete(key string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	delete(v.entries, key)
}

// ListEntries lists all non-expired entries in the vault
func (v *Vault) ListEntries() []string {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	var keys []string
	now := time.Now()

	for key, entry := range v.entries {
		if now.Before(entry.ExpiresAt) {
			keys = append(keys, key)
		}
	}

	return keys
}

// GetAccessLog retrieves the access log for a key
func (v *Vault) GetAccessLog(key string) ([]AccessRecord, error) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	entry, exists := v.entries[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	return entry.AccessLog, nil
}

// encrypt encrypts data using AES-GCM
func (v *Vault) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(v.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// Prepend nonce to ciphertext
	result := append(nonce, ciphertext...)
	return result, nil
}

// decrypt decrypts data using AES-GCM
func (v *Vault) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(v.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Extract nonce (first 12 bytes)
	if len(ciphertext) < 12 {
		return nil, fmt.Errorf("invalid ciphertext")
	}
	nonce := ciphertext[:12]
	encryptedData := ciphertext[12:]

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// cleanupExpired removes expired entries from the vault
func (v *Vault) cleanupExpired() {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	now := time.Now()
	for key, entry := range v.entries {
		if now.After(entry.ExpiresAt) {
			delete(v.entries, key)
		}
	}
}

// StartCleanup starts a background goroutine to clean up expired entries
func (v *Vault) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			v.cleanupExpired()
		}
	}()
}
