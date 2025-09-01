package main

import (
	"crypto/rand"
	"fmt"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/fpe"
)

func main() {
	// Generate a key
	key := make([]byte, 16)
	rand.Read(key)
	tweak := []byte("my-tweak")

	fpeInstance := fpe.New(key, tweak)

	// Encrypt a credit card number
	ccNumber := "4532015112830366"
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
	isValid := luhnCheck("4532015112830366")
	fmt.Printf("CC is valid (Luhn): %t\n", isValid)
}

// luhnCheck validates a number using the Luhn algorithm
func luhnCheck(number string) bool {
	sum := 0
	alt := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		sum += digit
		alt = !alt
	}

	return sum%10 == 0
}
