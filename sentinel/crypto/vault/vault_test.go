package vault

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"
)

func TestVaultStoreRetrieve(t *testing.T) {
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)
	
	v := NewVault(encryptionKey)
	
	key := "test-key"
	data := []byte("test data")
	
	// Store data
	err := v.Store(key, data, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	// Retrieve data
	retrievedData, err := v.Retrieve(key, "test-reason")
	if err != nil {
		t.Fatalf("Failed to retrieve data: %v", err)
	}
	
	if !bytes.Equal(data, retrievedData) {
		t.Error("Retrieved data does not match stored data")
	}
}

func TestVaultRetrieveNonExistentKey(t *testing.T) {
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)
	
	v := NewVault(encryptionKey)
	
	_, err := v.Retrieve("non-existent-key", "test-reason")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
}

func TestVaultDelete(t *testing.T) {
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)
	
	v := NewVault(encryptionKey)
	
	key := "test-key"
	data := []byte("test data")
	
	// Store data
	err := v.Store(key, data, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	// Delete data
	v.Delete(key)
	
	// Try to retrieve deleted data
	_, err = v.Retrieve(key, "test-reason")
	if err == nil {
		t.Error("Expected error for deleted key")
	}
}

func TestVaultListEntries(t *testing.T) {
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)
	
	v := NewVault(encryptionKey)
	
	// Store some data
	err := v.Store("key1", []byte("data1"), 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	err = v.Store("key2", []byte("data2"), 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	// List entries
	entries := v.ListEntries()
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}

func TestVaultGetAccessLog(t *testing.T) {
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)
	
	v := NewVault(encryptionKey)
	
	key := "test-key"
	data := []byte("test data")
	
	// Store data
	err := v.Store(key, data, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	// Retrieve data
	_, err = v.Retrieve(key, "test-reason")
	if err != nil {
		t.Fatalf("Failed to retrieve data: %v", err)
	}
	
	// Get access log
	accessLog, err := v.GetAccessLog(key)
	if err != nil {
		t.Fatalf("Failed to get access log: %v", err)
	}
	
	if len(accessLog) != 1 {
		t.Errorf("Expected 1 access log entry, got %d", len(accessLog))
	}
	
	if accessLog[0].Reason != "test-reason" {
		t.Errorf("Expected reason 'test-reason', got '%s'", accessLog[0].Reason)
	}
}

func TestVaultExpiredEntry(t *testing.T) {
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)
	
	v := NewVault(encryptionKey)
	
	key := "test-key"
	data := []byte("test data")
	
	// Store data with short TTL
	err := v.Store(key, data, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	// Wait for expiration
	time.Sleep(10 * time.Millisecond)
	
	// Try to retrieve expired data
	_, err = v.Retrieve(key, "test-reason")
	if err == nil {
		t.Error("Expected error for expired key")
	}
}