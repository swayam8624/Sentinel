package kms

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestKeyManagerGenerateDataKey(t *testing.T) {
	masterKey := make([]byte, 32)
	rand.Read(masterKey)
	
	km := NewKeyManager(masterKey)
	
	dataKey, err := km.GenerateDataKey()
	if err != nil {
		t.Fatalf("Failed to generate data key: %v", err)
	}
	
	if len(dataKey) != 32 {
		t.Errorf("Expected data key length 32, got %d", len(dataKey))
	}
}

func TestKeyManagerEncryptDecryptDataKey(t *testing.T) {
	masterKey := make([]byte, 32)
	rand.Read(masterKey)
	
	km := NewKeyManager(masterKey)
	
	// Generate data key
	dataKey, err := km.GenerateDataKey()
	if err != nil {
		t.Fatalf("Failed to generate data key: %v", err)
	}
	
	// Encrypt data key
	encryptedDataKey, err := km.EncryptDataKey(dataKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data key: %v", err)
	}
	
	if len(encryptedDataKey) <= 12 {
		t.Error("Encrypted data key should be longer than nonce")
	}
	
	// Decrypt data key
	decryptedDataKey, err := km.DecryptDataKey(encryptedDataKey)
	if err != nil {
		t.Fatalf("Failed to decrypt data key: %v", err)
	}
	
	if !bytes.Equal(dataKey, decryptedDataKey) {
		t.Error("Decrypted data key does not match original")
	}
}

func TestKeyManagerEncryptDecryptData(t *testing.T) {
	masterKey := make([]byte, 32)
	rand.Read(masterKey)
	
	km := NewKeyManager(masterKey)
	
	// Generate data key
	dataKey, err := km.GenerateDataKey()
	if err != nil {
		t.Fatalf("Failed to generate data key: %v", err)
	}
	
	// Encrypt data
	data := []byte("This is test data for encryption")
	encryptedData, err := km.EncryptData(data, dataKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	
	if len(encryptedData) <= 12 {
		t.Error("Encrypted data should be longer than nonce")
	}
	
	// Decrypt data
	decryptedData, err := km.DecryptData(encryptedData, dataKey)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}
	
	if !bytes.Equal(data, decryptedData) {
		t.Error("Decrypted data does not match original")
	}
}

func TestKeyManagerDecryptInvalidDataKey(t *testing.T) {
	masterKey := make([]byte, 32)
	rand.Read(masterKey)
	
	km := NewKeyManager(masterKey)
	
	// Try to decrypt invalid data
	_, err := km.DecryptDataKey([]byte("invalid"))
	if err == nil {
		t.Error("Expected error when decrypting invalid data key")
	}
}

func TestKeyManagerDecryptInvalidData(t *testing.T) {
	masterKey := make([]byte, 32)
	rand.Read(masterKey)
	
	km := NewKeyManager(masterKey)
	
	// Generate data key
	dataKey, err := km.GenerateDataKey()
	if err != nil {
		t.Fatalf("Failed to generate data key: %v", err)
	}
	
	// Try to decrypt invalid data
	_, err = km.DecryptData([]byte("invalid"), dataKey)
	if err == nil {
		t.Error("Expected error when decrypting invalid data")
	}
}