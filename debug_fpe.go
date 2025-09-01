package main

import (
	"fmt"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/fpe"
)

func main() {
	key := []byte{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	tweak := []byte("test-tweak")

	fpeInstance := fpe.New(key, tweak)

	// Test with a simple number
	plaintext := "123456789"
	fmt.Printf("Original: %s\n", plaintext)

	// Show the rotations
	baseRotation := key[0] % 10
	fmt.Printf("Base rotation: %d\n", baseRotation)

	for i := 0; i < len(plaintext); i++ {
		digit := plaintext[i] - '0'
		rotation := (baseRotation + byte(i)) % 10
		encryptedDigit := (digit + rotation) % 10
		fmt.Printf("Position %d: digit=%d, rotation=%d, encrypted=%d\n", i, digit, rotation, encryptedDigit)
	}

	ciphertext, err := fpeInstance.Encrypt(plaintext)
	if err != nil {
		fmt.Printf("Failed to encrypt: %v\n", err)
		return
	}
	fmt.Printf("Encrypted: %s\n", ciphertext)

	// Show the decryption step by step
	fmt.Println("Decryption steps:")
	for i := 0; i < len(ciphertext); i++ {
		digit := ciphertext[i] - '0'
		rotation := (baseRotation + byte(i)) % 10
		decryptedDigit := digit - rotation
		fmt.Printf("  Position %d: digit=%d, rotation=%d, diff=%d", i, digit, rotation, decryptedDigit)
		if decryptedDigit < 0 {
			decryptedDigit += 10
			fmt.Printf(" (adjusted to %d)", decryptedDigit)
		}
		fmt.Println()
	}

	decrypted, err := fpeInstance.Decrypt(ciphertext)
	if err != nil {
		fmt.Printf("Failed to decrypt: %v\n", err)
		return
	}
	fmt.Printf("Decrypted: %s\n", decrypted)

	fmt.Printf("Match: %t\n", plaintext == decrypted)

	// Show byte values
	fmt.Printf("Decrypted bytes: ")
	for i := 0; i < len(decrypted); i++ {
		fmt.Printf("%d ", decrypted[i])
	}
	fmt.Println()
}
