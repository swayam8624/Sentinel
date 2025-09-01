package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/fpe"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/hkdf"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/kms"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/merkle"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/nonce"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/vault"
)

func main() {
	fmt.Println("Sentinel Crypto Components Demo")
	fmt.Println("==============================")

	// Demo HKDF
	demoHKDF()

	// Demo Nonce
	demoNonce()

	// Demo KMS
	demoKMS()

	// Demo FPE
	demoFPE()

	// Demo Merkle Tree
	demoMerkle()

	// Demo Vault
	demoVault()
}

func demoHKDF() {
	fmt.Println("\n1. HKDF Demo:")

	secret := []byte("my-secret-key")
	salt := []byte("my-salt")
	info := []byte("my-info")

	key, err := hkdf.DeriveKey(secret, salt, info, 32)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Derived key: %x\n", key)
}

func demoNonce() {
	fmt.Println("\n2. Nonce Management Demo:")

	nm := nonce.NewNonceManager(10 * time.Second)

	// Generate a nonce
	nonce1, err := nm.GenerateNonce(12)
	if err != nil {
		fmt.Printf("Error generating nonce: %v\n", err)
		return
	}

	fmt.Printf("Generated nonce: %x\n", nonce1)

	// Check uniqueness
	isUnique := nm.IsUnique(nonce1)
	fmt.Printf("Nonce is unique: %t\n", isUnique)

	// Try again with same nonce
	isUnique2 := nm.IsUnique(nonce1)
	fmt.Printf("Same nonce is unique: %t\n", isUnique2)
}

func demoKMS() {
	fmt.Println("\n3. KMS Demo:")

	// Generate master key
	masterKey := make([]byte, 32)
	rand.Read(masterKey)

	km := kms.NewKeyManager(masterKey)

	// Generate data key
	dataKey, err := km.GenerateDataKey()
	if err != nil {
		fmt.Printf("Error generating data key: %v\n", err)
		return
	}

	fmt.Printf("Generated data key: %x\n", dataKey)

	// Encrypt data key
	encryptedDataKey, err := km.EncryptDataKey(dataKey)
	if err != nil {
		fmt.Printf("Error encrypting data key: %v\n", err)
		return
	}

	fmt.Printf("Encrypted data key: %x\n", encryptedDataKey)

	// Decrypt data key
	decryptedDataKey, err := km.DecryptDataKey(encryptedDataKey)
	if err != nil {
		fmt.Printf("Error decrypting data key: %v\n", err)
		return
	}

	fmt.Printf("Decrypted data key: %x\n", decryptedDataKey)

	// Encrypt some data
	data := []byte("This is sensitive data that needs encryption")
	encryptedData, err := km.EncryptData(data, dataKey)
	if err != nil {
		fmt.Printf("Error encrypting data: %v\n", err)
		return
	}

	fmt.Printf("Encrypted data: %x\n", encryptedData)

	// Decrypt the data
	decryptedData, err := km.DecryptData(encryptedData, decryptedDataKey)
	if err != nil {
		fmt.Printf("Error decrypting data: %v\n", err)
		return
	}

	fmt.Printf("Decrypted data: %s\n", string(decryptedData))
}

func demoFPE() {
	fmt.Println("\n4. FPE Demo:")

	key := make([]byte, 16)
	rand.Read(key)
	tweak := []byte("my-tweak")

	fpeInstance := fpe.New(key, tweak)

	// Encrypt a credit card number
	ccNumber := "123456789"
	encryptedCC, err := fpeInstance.Encrypt(ccNumber)
	if err != nil {
		fmt.Printf("Error encrypting CC: %v\n", err)
		return
	}

	fmt.Printf("Original CC: %s\n", ccNumber)
	fmt.Printf("Encrypted CC: %s\n", encryptedCC)

	// Decrypt the credit card number
	decryptedCC, err := fpeInstance.Decrypt(encryptedCC)
	if err != nil {
		fmt.Printf("Error decrypting CC: %v\n", err)
		return
	}

	fmt.Printf("Decrypted CC: %s\n", decryptedCC)

	// Validate with Luhn check
	isValid := fpe.LuhnCheck("4532123456789010")
	fmt.Printf("CC is valid (Luhn): %t\n", isValid)
}

func demoMerkle() {
	fmt.Println("\n5. Merkle Tree Demo:")

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

	fmt.Printf("Merkle root hash: %x\n", tree.RootHash())

	// Generate proof for first leaf
	proof, err := tree.GenerateProof(0)
	if err != nil {
		fmt.Printf("Error generating proof: %v\n", err)
		return
	}

	fmt.Printf("Proof for leaf 0: %d elements\n", len(proof))

	// Verify proof
	isValid := tree.VerifyProof(data[0], proof, tree.RootHash())
	fmt.Printf("Proof is valid: %t\n", isValid)
}

func demoVault() {
	fmt.Println("\n6. Token Vault Demo:")

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

	err := v.Store(key, data, 1*time.Hour)
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
}
