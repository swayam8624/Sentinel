package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/kms"
)

func main() {
	// Example 1: Using Local KMS
	fmt.Println("=== Local KMS Example ===")
	localKMSExample()

	// Example 2: Using Envelope Encryption with Local KMS Client
	fmt.Println("\n=== Envelope Encryption Example ===")
	envelopeEncryptionExample()
}

func localKMSExample() {
	// Create a local KMS with a master key
	masterKey := []byte("masterkey12345678901234567890123") // 32 bytes
	kmsService := kms.NewLocalKMS(masterKey)

	// Generate a data key
	plaintextKey, encryptedKey, err := kmsService.GenerateDataKey()
	if err != nil {
		log.Fatalf("Failed to generate data key: %v", err)
	}

	fmt.Printf("Generated plaintext key: %x\n", plaintextKey)
	fmt.Printf("Generated encrypted key length: %d bytes\n", len(encryptedKey))

	// Decrypt the data key
	decryptedKey, err := kmsService.DecryptDataKey(encryptedKey)
	if err != nil {
		log.Fatalf("Failed to decrypt data key: %v", err)
	}

	fmt.Printf("Decrypted key matches original: %t\n", string(plaintextKey) == string(decryptedKey))

	// Encrypt some data
	plaintext := []byte("This is a secret message that needs to be protected!")
	ciphertext, err := kmsService.Encrypt(plaintext, plaintextKey)
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}

	fmt.Printf("Original plaintext: %s\n", plaintext)
	fmt.Printf("Encrypted data length: %d bytes\n", len(ciphertext))

	// Decrypt the data
	decrypted, err := kmsService.Decrypt(ciphertext, plaintextKey)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}

	fmt.Printf("Decrypted data matches original: %t\n", string(plaintext) == string(decrypted))
	fmt.Printf("Decrypted plaintext: %s\n", decrypted)
}

func envelopeEncryptionExample() {
	// Create a local KMS client for envelope encryption
	localClient := kms.NewLocalKMSClient()

	// Create a key for envelope encryption
	ctx := context.Background()
	metadata, err := localClient.GenerateKey(ctx, "AES_256")
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	fmt.Printf("Generated key ID: %s\n", metadata.KeyID)

	// Create envelope encryption service
	envelope := kms.NewEnvelopeEncryption(localClient, 32)

	// Encrypt data using envelope encryption
	plaintext := []byte("This is a secret message protected with envelope encryption!")
	ciphertext, err := envelope.EncryptWithEnvelope(ctx, metadata.KeyID, plaintext)
	if err != nil {
		log.Fatalf("Failed to encrypt with envelope: %v", err)
	}

	fmt.Printf("Original plaintext: %s\n", plaintext)
	fmt.Printf("Encrypted data length: %d bytes\n", len(ciphertext))

	// Decrypt data using envelope encryption
	decrypted, err := envelope.DecryptWithEnvelope(ctx, metadata.KeyID, ciphertext)
	if err != nil {
		log.Fatalf("Failed to decrypt with envelope: %v", err)
	}

	fmt.Printf("Decrypted data matches original: %t\n", string(plaintext) == string(decrypted))
	fmt.Printf("Decrypted plaintext: %s\n", decrypted)
}