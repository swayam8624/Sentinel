package fpe

import (
	"fmt"
)

// FPE implements a simplified format preserving encryption
type FPE struct {
	key   []byte
	tweak []byte
}

// New creates a new FPE instance
func New(key, tweak []byte) *FPE {
	return &FPE{
		key:   key,
		tweak: tweak,
	}
}

// Encrypt encrypts a numeric string while preserving its format
func (f *FPE) Encrypt(plaintext string) (string, error) {
	// Validate input - must be numeric
	for _, r := range plaintext {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("plaintext must contain only digits")
		}
	}

	// Simple encryption - rotate each digit by a position-dependent amount
	result := make([]byte, len(plaintext))

	// Use first byte of key as base rotation amount
	baseRotation := int(f.key[0]) % 10

	for i := 0; i < len(plaintext); i++ {
		digit := int(plaintext[i] - '0')
		// Position-dependent rotation
		rotation := (baseRotation + i) % 10
		encryptedDigit := (digit + rotation) % 10
		result[i] = byte(encryptedDigit + '0')
	}

	return string(result), nil
}

// Decrypt decrypts a numeric string
func (f *FPE) Decrypt(ciphertext string) (string, error) {
	// Validate input - must be numeric
	for _, r := range ciphertext {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("ciphertext must contain only digits")
		}
	}

	// Simple decryption - reverse the position-dependent rotation
	result := make([]byte, len(ciphertext))

	// Use first byte of key as base rotation amount
	baseRotation := int(f.key[0]) % 10

	for i := 0; i < len(ciphertext); i++ {
		digit := int(ciphertext[i] - '0')
		// Position-dependent rotation
		rotation := (baseRotation + i) % 10
		// Reverse rotation
		decryptedDigit := digit - rotation
		// Handle negative results properly
		for decryptedDigit < 0 {
			decryptedDigit += 10
		}
		result[i] = byte(decryptedDigit + '0')
	}

	return string(result), nil
}

// LuhnCheck validates a number using the Luhn algorithm
func LuhnCheck(number string) bool {
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
