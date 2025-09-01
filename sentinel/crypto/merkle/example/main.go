package main

import (
	"fmt"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/merkle"
)

func main() {
	// Create some data
	data := [][]byte{
		[]byte("Transaction 1"),
		[]byte("Transaction 2"),
		[]byte("Transaction 3"),
		[]byte("Transaction 4"),
	}

	// Create Merkle tree
	tree, err := merkle.NewMerkleTree(data)
	if err != nil {
		fmt.Printf("Error creating Merkle tree: %v\n", err)
		return
	}

	rootHash := tree.RootHash()
	fmt.Printf("Merkle root hash: %x\n", rootHash)

	// Generate proof for first leaf
	proof, err := tree.GenerateProof(0)
	if err != nil {
		fmt.Printf("Error generating proof: %v\n", err)
		return
	}

	fmt.Printf("Proof for leaf 0: %d elements\n", len(proof))

	// Verify proof
	isValid := tree.VerifyProof(data[0], proof, rootHash)
	fmt.Printf("Proof is valid: %t\n", isValid)

	// Try to verify with wrong data
	wrongData := []byte("Wrong transaction")
	isValidWrong := tree.VerifyProof(wrongData, proof, rootHash)
	fmt.Printf("Proof with wrong data is valid: %t\n", isValidWrong)
}
