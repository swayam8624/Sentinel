package fpe

import (
	"crypto/rand"
	"testing"
)

func TestFPEEncryptDecrypt(t *testing.T) {
	key := []byte{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	tweak := []byte("test-tweak")

	fpe := New(key, tweak)

	// Test with a simple number
	plaintext := "123456789"
	ciphertext, err := fpe.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if len(ciphertext) != len(plaintext) {
		t.Errorf("Ciphertext length mismatch: expected %d, got %d", len(plaintext), len(ciphertext))
	}

	// Decrypt
	decrypted, err := fpe.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if plaintext != decrypted {
		t.Errorf("Decryption mismatch: expected %s, got %s", plaintext, decrypted)
	}
}

func TestFPEEncryptNonNumeric(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	tweak := []byte("test-tweak")

	fpe := New(key, tweak)

	// Test with non-numeric input
	_, err := fpe.Encrypt("abc123")
	if err == nil {
		t.Error("Expected error for non-numeric input")
	}
}

func TestFPEDecryptNonNumeric(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	tweak := []byte("test-tweak")

	fpe := New(key, tweak)

	// Test with non-numeric input
	_, err := fpe.Decrypt("abc123")
	if err == nil {
		t.Error("Expected error for non-numeric input")
	}
}

func TestFPELuhnCheck(t *testing.T) {
	// Test valid credit card number (this one passes Luhn check)
	validCC := "4000056655665556"
	if !LuhnCheck(validCC) {
		t.Error("Expected valid credit card number to pass Luhn check")
	}

	// Test invalid credit card number
	invalidCC := "1234567890123456"
	if LuhnCheck(invalidCC) {
		t.Error("Expected invalid credit card number to fail Luhn check")
	}
}
