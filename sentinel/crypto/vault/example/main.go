package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/vault"
)

func main() {
	// Generate encryption key
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)

	// Create vault
	v := vault.NewVault(encryptionKey)

	// Start cleanup goroutine
	v.StartCleanup(5 * time.Second)

	// Store some data
	key := "my-secret-token"
	data := []byte("This is a secret token value")
	ttl := 1 * time.Hour

	err := v.Store(key, data, ttl)
	if err != nil {
		fmt.Printf("Error storing data: %v\n", err)
		return
	}

	fmt.Printf("Stored data with key: %s\n", key)

	// Retrieve data
	retrievedData, err := v.Retrieve(key, "demo-access")
	if err != nil {
		fmt.Printf("Error retrieving data: %v\n", err)
		return
	}

	fmt.Printf("Retrieved data: %s\n", string(retrievedData))

	// Check access log
	accessLog, err := v.GetAccessLog(key)
	if err != nil {
		fmt.Printf("Error getting access log: %v\n", err)
		return
	}

	fmt.Printf("Access log entries: %d\n", len(accessLog))
	if len(accessLog) > 0 {
		fmt.Printf("Last access reason: %s\n", accessLog[len(accessLog)-1].Reason)
	}
}
