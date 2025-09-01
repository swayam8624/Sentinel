package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/vault"
)

func main() {
	// Example of vault usage
	fmt.Println("=== Vault Example ===")

	// Create a master key and vault
	masterKey := []byte("masterkey12345678901234567890123") // 32 bytes
	vault := vault.NewVault(masterKey)

	// Store some sensitive data
	tokenID := "cc-token-12345"
	creditCard := "4532015112830366"
	ttl := 24 * time.Hour
	accessReason := "initial storage"

	err := vault.Store(tokenID, []byte(creditCard), ttl, accessReason)
	if err != nil {
		log.Fatalf("Failed to store credit card: %v", err)
	}

	fmt.Printf("Stored credit card %s with token ID %s\n", creditCard, tokenID)

	// Retrieve the data
	retrievedData, err := vault.Retrieve(tokenID, "retrieval for transaction")
	if err != nil {
		log.Fatalf("Failed to retrieve credit card: %v", err)
	}

	fmt.Printf("Retrieved credit card: %s\n", string(retrievedData))

	// Get entry information
	info, err := vault.GetEntryInfo(tokenID)
	if err != nil {
		log.Fatalf("Failed to get entry info: %v", err)
	}

	fmt.Printf("Entry created at: %s\n", info.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Entry expires at: %s\n", info.ExpiresAt.Format(time.RFC3339))
	fmt.Printf("Last access reason: %s\n", info.AccessReason)

	// List all tokens
	tokens := vault.ListTokens()
	fmt.Printf("Total tokens in vault: %d\n", len(tokens))

	for _, token := range tokens {
		fmt.Printf("  - %s\n", token)
	}

	// Delete the token
	vault.Delete(tokenID)
	fmt.Printf("Deleted token %s\n", tokenID)

	// Try to retrieve deleted token (should fail)
	_, err = vault.Retrieve(tokenID, "post-deletion access")
	if err != nil {
		fmt.Printf("Expected error when retrieving deleted token: %v\n", err)
	}
}
