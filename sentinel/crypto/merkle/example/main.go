package main

import (
	"fmt"
	"log"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/merkle"
)

func main() {
	// Example audit log entries
	auditLogs := [][]byte{
		[]byte("User login: alice@example.com"),
		[]byte("PII detection: Credit card number found"),
		[]byte("Redaction applied: CC number masked"),
		[]byte("Request forwarded to OpenAI"),
		[]byte("Response received from OpenAI"),
		[]byte("Detokenization performed"),
		[]byte("Response sent to client"),
	}

	// Create Merkle tree from audit logs
	tree, err := merkle.NewMerkleTree(auditLogs)
	if err != nil {
		log.Fatalf("Failed to create Merkle tree: %v", err)
	}

	// Get root hash
	rootHash := tree.GetRootHash()
	fmt.Printf("Merkle Root Hash: %x\n", rootHash)

	// Generate proof for a specific log entry (e.g., the PII detection log)
	index := 1 // Index of "PII detection: Credit card number found"
	proof, err := tree.GenerateProof(index, auditLogs)
	if err != nil {
		log.Fatalf("Failed to generate proof: %v", err)
	}

	fmt.Printf("\nProof for log entry #%d:\n", index)
	for i, hash := range proof {
		fmt.Printf("  Proof %d: %x\n", i, hash)
	}

	// Verify the proof
	isValid := tree.VerifyProof(auditLogs[index], proof, index)
	fmt.Printf("\nProof verification result: %t\n", isValid)

	// Try to verify with incorrect data
	isValid = tree.VerifyProof([]byte("tampered data"), proof, index)
	fmt.Printf("Proof verification with tampered data: %t\n", isValid)
}
