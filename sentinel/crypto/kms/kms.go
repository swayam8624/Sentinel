package kms

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"sync"
	"time"
)

// KeyManager provides envelope encryption services
type KeyManager struct {
	masterKey []byte
}

// NewKeyManager creates a new KeyManager with the provided master key
func NewKeyManager(masterKey []byte) *KeyManager {
	return &KeyManager{
		masterKey: masterKey,
	}
}

// GenerateDataKey generates a new data encryption key
func (km *KeyManager) GenerateDataKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate data key: %w", err)
	}
	return key, nil
}

// EncryptDataKey encrypts a data key using the master key
func (km *KeyManager) EncryptDataKey(dataKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(km.masterKey)
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

	// Encrypt the data key
	ciphertext := aesGCM.Seal(nil, nonce, dataKey, nil)

	// Prepend nonce to ciphertext
	result := append(nonce, ciphertext...)
	return result, nil
}

// DecryptDataKey decrypts an encrypted data key using the master key
func (km *KeyManager) DecryptDataKey(encryptedDataKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(km.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Extract nonce (first 12 bytes)
	if len(encryptedDataKey) < 12 {
		return nil, fmt.Errorf("invalid encrypted data key")
	}
	nonce := encryptedDataKey[:12]
	ciphertext := encryptedDataKey[12:]

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the data key
	dataKey, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data key: %w", err)
	}

	return dataKey, nil
}

// EncryptData encrypts data using a data key
func (km *KeyManager) EncryptData(data, dataKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(dataKey)
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
	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	// Prepend nonce to ciphertext
	result := append(nonce, ciphertext...)
	return result, nil
}

// DecryptData decrypts data using a data key
func (km *KeyManager) DecryptData(encryptedData, dataKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Extract nonce (first 12 bytes)
	if len(encryptedData) < 12 {
		return nil, fmt.Errorf("invalid encrypted data")
	}
	nonce := encryptedData[:12]
	ciphertext := encryptedData[12:]

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// KeyManagementService defines the interface for KMS operations
type KeyManagementService interface {
	GenerateDataKey() ([]byte, []byte, error) // Returns (plaintext key, encrypted key, error)
	DecryptDataKey(encryptedKey []byte) ([]byte, error)
	Encrypt(plaintext, key []byte) ([]byte, error)
	Decrypt(ciphertext, key []byte) ([]byte, error)
}

// LocalKMS is a local implementation of KMS for development/testing
type LocalKMS struct {
	masterKey []byte
}

// NewLocalKMS creates a new local KMS instance
func NewLocalKMS(masterKey []byte) *LocalKMS {
	return &LocalKMS{
		masterKey: masterKey,
	}
}

// GenerateDataKey generates a new data key and encrypts it with the master key
func (k *LocalKMS) GenerateDataKey() ([]byte, []byte, error) {
	// Generate a random data key (256 bits)
	dataKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dataKey); err != nil {
		return nil, nil, fmt.Errorf("failed to generate data key: %w", err)
	}

	// Encrypt the data key with the master key
	encryptedKey, err := k.Encrypt(dataKey, k.masterKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt data key: %w", err)
	}

	return dataKey, encryptedKey, nil
}

// DecryptDataKey decrypts an encrypted data key with the master key
func (k *LocalKMS) DecryptDataKey(encryptedKey []byte) ([]byte, error) {
	return k.Decrypt(encryptedKey, k.masterKey)
}

// Encrypt encrypts plaintext using AES-GCM
func (k *LocalKMS) Encrypt(plaintext, key []byte) ([]byte, error) {
	// Ensure key is proper length
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: %d", len(key))
	}

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
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-GCM
func (k *LocalKMS) Decrypt(ciphertext, key []byte) ([]byte, error) {
	// Ensure key is proper length
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: %d", len(key))
	}

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

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// KMSClient interface for key management services
type KMSClient interface {
	// GenerateKey generates a new key
	GenerateKey(ctx context.Context, keySpec string) (KeyMetadata, error)

	// Encrypt encrypts data with a key
	Encrypt(ctx context.Context, keyID string, plaintext []byte) (EncryptOutput, error)

	// Decrypt decrypts data with a key
	Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error)

	// RotateKey rotates a key
	RotateKey(ctx context.Context, keyID string) (KeyMetadata, error)

	// GetKeyMetadata gets key metadata
	GetKeyMetadata(ctx context.Context, keyID string) (KeyMetadata, error)
}

// KeyMetadata represents key metadata
type KeyMetadata struct {
	KeyID           string    `json:"key_id"`
	Arn             string    `json:"arn"`
	CreationDate    time.Time `json:"creation_date"`
	Description     string    `json:"description"`
	Enabled         bool      `json:"enabled"`
	ExpirationModel string    `json:"expiration_model"`
	KeyManager      string    `json:"key_manager"`
	KeyState        string    `json:"key_state"`
	KeyUsage        string    `json:"key_usage"`
	Origin          string    `json:"origin"`
}

// EncryptOutput represents encryption output
type EncryptOutput struct {
	CiphertextBlob []byte `json:"ciphertext_blob"`
	KeyID          string `json:"key_id"`
}

// LocalKMSClient implements KMSClient for local development
type LocalKMSClient struct {
	keys  map[string][]byte
	mutex sync.RWMutex
	keyID int
}

// NewLocalKMSClient creates a new local KMS client
func NewLocalKMSClient() *LocalKMSClient {
	return &LocalKMSClient{
		keys:  make(map[string][]byte),
		keyID: 0,
	}
}

// GenerateKey generates a new key
func (lk *LocalKMSClient) GenerateKey(ctx context.Context, keySpec string) (KeyMetadata, error) {
	lk.mutex.Lock()
	defer lk.mutex.Unlock()

	// Generate a random key
	key := make([]byte, 32) // 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to generate key: %w", err)
	}

	// Generate key ID
	keyID := fmt.Sprintf("local-key-%d", lk.keyID)
	lk.keyID++

	// Store key
	lk.keys[keyID] = key

	// Create metadata
	metadata := KeyMetadata{
		KeyID:        keyID,
		Arn:          fmt.Sprintf("arn:local:kms:::key/%s", keyID),
		CreationDate: time.Now(),
		Description:  "Local development key",
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "Enabled",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "AWS_KMS",
	}

	return metadata, nil
}

// Encrypt encrypts data with a key
func (lk *LocalKMSClient) Encrypt(ctx context.Context, keyID string, plaintext []byte) (EncryptOutput, error) {
	lk.mutex.RLock()
	key, exists := lk.keys[keyID]
	lk.mutex.RUnlock()

	if !exists {
		return EncryptOutput{}, fmt.Errorf("key not found: %s", keyID)
	}

	// Simple encryption - in a real implementation, we'd use proper encryption
	// This is just for demonstration
	ciphertext := make([]byte, len(plaintext))
	for i, b := range plaintext {
		ciphertext[i] = b ^ key[i%len(key)]
	}

	return EncryptOutput{
		CiphertextBlob: ciphertext,
		KeyID:          keyID,
	}, nil
}

// Decrypt decrypts data with a key
func (lk *LocalKMSClient) Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	lk.mutex.RLock()
	key, exists := lk.keys[keyID]
	lk.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("key not found: %s", keyID)
	}

	// Simple decryption - in a real implementation, we'd use proper decryption
	// This is just for demonstration
	plaintext := make([]byte, len(ciphertext))
	for i, b := range ciphertext {
		plaintext[i] = b ^ key[i%len(key)]
	}

	return plaintext, nil
}

// RotateKey rotates a key
func (lk *LocalKMSClient) RotateKey(ctx context.Context, keyID string) (KeyMetadata, error) {
	lk.mutex.Lock()
	defer lk.mutex.Unlock()

	// Generate a new key
	newKey := make([]byte, 32) // 256 bits
	_, err := rand.Read(newKey)
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to generate key: %w", err)
	}

	// Replace the key
	lk.keys[keyID] = newKey

	// Create metadata
	metadata := KeyMetadata{
		KeyID:        keyID,
		Arn:          fmt.Sprintf("arn:local:kms:::key/%s", keyID),
		CreationDate: time.Now(),
		Description:  "Local development key (rotated)",
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "Enabled",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "AWS_KMS",
	}

	return metadata, nil
}

// GetKeyMetadata gets key metadata
func (lk *LocalKMSClient) GetKeyMetadata(ctx context.Context, keyID string) (KeyMetadata, error) {
	lk.mutex.RLock()
	_, exists := lk.keys[keyID]
	lk.mutex.RUnlock()

	if !exists {
		return KeyMetadata{}, fmt.Errorf("key not found: %s", keyID)
	}

	// Create metadata
	metadata := KeyMetadata{
		KeyID:        keyID,
		Arn:          fmt.Sprintf("arn:local:kms:::key/%s", keyID),
		CreationDate: time.Now(),
		Description:  "Local development key",
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "Enabled",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "AWS_KMS",
	}

	return metadata, nil
}

// EnvelopeEncryption implements envelope encryption using KMS
type EnvelopeEncryption struct {
	kmsClient     KMSClient
	dataKeyLength int
}

// NewEnvelopeEncryption creates a new envelope encryption instance
func NewEnvelopeEncryption(kmsClient KMSClient, dataKeyLength int) *EnvelopeEncryption {
	return &EnvelopeEncryption{
		kmsClient:     kmsClient,
		dataKeyLength: dataKeyLength,
	}
}

// EncryptWithEnvelope encrypts data using envelope encryption
func (ee *EnvelopeEncryption) EncryptWithEnvelope(ctx context.Context, keyID string, plaintext []byte) ([]byte, error) {
	// Generate a data key
	dataKey := make([]byte, ee.dataKeyLength)
	_, err := rand.Read(dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate data key: %w", err)
	}

	// Encrypt the data with the data key using AES-GCM
	encryptedData, err := encryptWithAESGCM(plaintext, dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Encrypt the data key with the KMS key
	encryptOutput, err := ee.kmsClient.Encrypt(ctx, keyID, dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data key: %w", err)
	}

	// Combine encrypted data key length, encrypted data key and encrypted data
	// We need to store the length of the encrypted data key to properly extract it during decryption
	dataKeyLen := len(encryptOutput.CiphertextBlob)

	// Create a buffer to hold the length (4 bytes for uint32)
	lenBuf := make([]byte, 4)
	// Convert length to bytes (big endian)
	lenBuf[0] = byte(dataKeyLen >> 24)
	lenBuf[1] = byte(dataKeyLen >> 16)
	lenBuf[2] = byte(dataKeyLen >> 8)
	lenBuf[3] = byte(dataKeyLen)

	// Combine everything
	result := append(lenBuf, encryptOutput.CiphertextBlob...)
	result = append(result, encryptedData...)

	return result, nil
}

// DecryptWithEnvelope decrypts data using envelope encryption
func (ee *EnvelopeEncryption) DecryptWithEnvelope(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	// Extract encrypted data key length (first 4 bytes)
	if len(ciphertext) < 4 {
		return nil, fmt.Errorf("ciphertext too short to contain length header")
	}

	dataKeyLen := int(ciphertext[0])<<24 | int(ciphertext[1])<<16 | int(ciphertext[2])<<8 | int(ciphertext[3])

	// Check if we have enough data
	if len(ciphertext) < 4+dataKeyLen {
		return nil, fmt.Errorf("ciphertext too short to contain encrypted data key")
	}

	// Extract encrypted data key and encrypted data
	encryptedDataKey := ciphertext[4 : 4+dataKeyLen]
	encryptedData := ciphertext[4+dataKeyLen:]

	// Decrypt the data key with the KMS key
	dataKey, err := ee.kmsClient.Decrypt(ctx, keyID, encryptedDataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data key: %w", err)
	}

	// Decrypt the data with the data key using AES-GCM
	plaintext, err := decryptWithAESGCM(encryptedData, dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// encryptWithAESGCM encrypts data with AES-GCM
func encryptWithAESGCM(plaintext, key []byte) ([]byte, error) {
	// Ensure key is proper length
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: %d", len(key))
	}

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
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decryptWithAESGCM decrypts data with AES-GCM
func decryptWithAESGCM(ciphertext, key []byte) ([]byte, error) {
	// Ensure key is proper length
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: %d", len(key))
	}

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

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
