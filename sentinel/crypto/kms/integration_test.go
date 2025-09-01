package kms

import (
	"context"
	"testing"
)

func TestKMSIntegration(t *testing.T) {
	ctx := context.Background()
	plaintext := []byte("This is a test message for KMS integration testing")

	// Test 1: Local KMS Service
	t.Run("LocalKMS", func(t *testing.T) {
		masterKey := []byte("masterkey12345678901234567890123") // 32 bytes
		kmsService := NewLocalKMS(masterKey)

		// Generate data key
		plaintextKey, encryptedKey, err := kmsService.GenerateDataKey()
		if err != nil {
			t.Fatalf("Failed to generate data key: %v", err)
		}

		// Decrypt data key
		decryptedKey, err := kmsService.DecryptDataKey(encryptedKey)
		if err != nil {
			t.Fatalf("Failed to decrypt data key: %v", err)
		}

		if string(plaintextKey) != string(decryptedKey) {
			t.Error("Decrypted key does not match original plaintext key")
		}

		// Encrypt data
		ciphertext, err := kmsService.Encrypt(plaintext, plaintextKey)
		if err != nil {
			t.Fatalf("Failed to encrypt data: %v", err)
		}

		// Decrypt data
		decrypted, err := kmsService.Decrypt(ciphertext, plaintextKey)
		if err != nil {
			t.Fatalf("Failed to decrypt data: %v", err)
		}

		if string(plaintext) != string(decrypted) {
			t.Error("Decrypted plaintext does not match original")
		}
	})

	// Test 2: Local KMS Client
	t.Run("LocalKMSClient", func(t *testing.T) {
		kmsClient := NewLocalKMSClient()

		// Generate key
		metadata, err := kmsClient.GenerateKey(ctx, "AES_256")
		if err != nil {
			t.Fatalf("Failed to generate key: %v", err)
		}

		// Encrypt data
		encryptOutput, err := kmsClient.Encrypt(ctx, metadata.KeyID, plaintext)
		if err != nil {
			t.Fatalf("Failed to encrypt data: %v", err)
		}

		// Decrypt data
		decrypted, err := kmsClient.Decrypt(ctx, metadata.KeyID, encryptOutput.CiphertextBlob)
		if err != nil {
			t.Fatalf("Failed to decrypt data: %v", err)
		}

		if string(plaintext) != string(decrypted) {
			t.Error("Decrypted plaintext does not match original")
		}

		// Get key metadata
		_, err = kmsClient.GetKeyMetadata(ctx, metadata.KeyID)
		if err != nil {
			t.Fatalf("Failed to get key metadata: %v", err)
		}

		// Rotate key
		_, err = kmsClient.RotateKey(ctx, metadata.KeyID)
		if err != nil {
			t.Fatalf("Failed to rotate key: %v", err)
		}
	})

	// Test 3: Envelope Encryption with Local KMS Client
	t.Run("EnvelopeEncryption", func(t *testing.T) {
		kmsClient := NewLocalKMSClient()

		// Generate key
		metadata, err := kmsClient.GenerateKey(ctx, "AES_256")
		if err != nil {
			t.Fatalf("Failed to generate key: %v", err)
		}

		// Create envelope encryption service
		envelope := NewEnvelopeEncryption(kmsClient, 32)

		// Encrypt with envelope
		ciphertext, err := envelope.EncryptWithEnvelope(ctx, metadata.KeyID, plaintext)
		if err != nil {
			t.Fatalf("Failed to encrypt with envelope: %v", err)
		}

		// Decrypt with envelope
		decrypted, err := envelope.DecryptWithEnvelope(ctx, metadata.KeyID, ciphertext)
		if err != nil {
			t.Fatalf("Failed to decrypt with envelope: %v", err)
		}

		if string(plaintext) != string(decrypted) {
			t.Error("Decrypted plaintext does not match original")
		}
	})

	// Test 4: KMS Factory
	t.Run("KMSFactory", func(t *testing.T) {
		// Test local KMS client creation
		localConfig := KMSConfig{
			Provider: "local",
		}

		kmsClient, err := NewKMSClient(ctx, localConfig)
		if err != nil {
			t.Fatalf("Failed to create KMS client: %v", err)
		}

		// Generate key
		metadata, err := kmsClient.GenerateKey(ctx, "AES_256")
		if err != nil {
			t.Fatalf("Failed to generate key: %v", err)
		}

		// Encrypt data
		encryptOutput, err := kmsClient.Encrypt(ctx, metadata.KeyID, plaintext)
		if err != nil {
			t.Fatalf("Failed to encrypt data: %v", err)
		}

		// Decrypt data
		decrypted, err := kmsClient.Decrypt(ctx, metadata.KeyID, encryptOutput.CiphertextBlob)
		if err != nil {
			t.Fatalf("Failed to decrypt data: %v", err)
		}

		if string(plaintext) != string(decrypted) {
			t.Error("Decrypted plaintext does not match original")
		}
	})
}