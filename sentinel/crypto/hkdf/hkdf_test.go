package hkdf

import (
	"testing"
)

func TestHKDFDeriveKey(t *testing.T) {
	secret := []byte("test-secret")
	salt := []byte("test-salt")
	info := []byte("test-info")

	hkdf := New(secret, salt, info)
	key, err := hkdf.DeriveKey(32)
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}

	// Test deterministic output
	key2, err := hkdf.DeriveKey(32)
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	for i := range key {
		if key[i] != key2[i] {
			t.Error("HKDF should produce deterministic output")
		}
	}
}

func TestHKDFDeriveKeyInvalidLength(t *testing.T) {
	secret := []byte("test-secret")
	salt := []byte("test-salt")
	info := []byte("test-info")

	hkdf := New(secret, salt, info)
	_, err := hkdf.DeriveKey(0)
	if err == nil {
		t.Error("Expected error for zero length key")
	}

	_, err = hkdf.DeriveKey(-1)
	if err == nil {
		t.Error("Expected error for negative length key")
	}
}

func TestHKDFConvenienceFunction(t *testing.T) {
	secret := []byte("test-secret")
	salt := []byte("test-salt")
	info := []byte("test-info")

	key, err := DeriveKey(secret, salt, info, 32)
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}
}