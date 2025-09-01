package main

import (
	"fmt"
	"log"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/fpe"
)

func main() {
	// Example of FPE usage for credit card numbers
	fmt.Println("=== FPE Example ===")

	// Create a key and tweak
	key := []byte("examplekey123456") // 16 bytes for AES-128
	tweak := []byte("tweak12")        // 7 bytes as required by FF3-1

	// Create FF3-1 instance
	ff31, err := fpe.NewFF31(key, tweak)
	if err != nil {
		log.Fatalf("Failed to create FF3-1 instance: %v", err)
	}

	// Example credit card number
	creditCard := "4532015112830366"

	fmt.Printf("Original credit card: %s\n", creditCard)

	// Encrypt the credit card number
	encrypted, err := ff31.Encrypt(creditCard, 10) // Radix 10 for decimal digits
	if err != nil {
		log.Fatalf("Failed to encrypt credit card: %v", err)
	}

	fmt.Printf("Encrypted credit card: %s\n", encrypted)

	// Validate using Luhn algorithm
	isValid := fpe.ValidateCreditCardLuhn(creditCard)
	fmt.Printf("Original credit card is valid: %t\n", isValid)

	// Note: In a real implementation, the decrypted value would match the original
	// For this simplified example, we're just demonstrating the API

	// Example of validating other types of numbers
	fmt.Println("\n=== Luhn Validation Examples ===")

	examples := []string{
		"4532015112830366", // Valid Visa
		"5555555555554444", // Valid Mastercard
		"1234567890123456", // Invalid
		"4000056655665556", // Valid Visa debit
	}

	for _, example := range examples {
		isValid := fpe.ValidateCreditCardLuhn(example)
		fmt.Printf("Card %s is valid: %t\n", example, isValid)
	}
}
