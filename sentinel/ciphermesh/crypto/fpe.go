package crypto

import (
	"fmt"
)

// FF3Cipher implements the FF3-1 format-preserving encryption algorithm
// This is a simplified implementation for demonstration purposes
type FF3Cipher struct {
	key   []byte
	tweak []byte
	radix int
}

// NewFF3Cipher creates a new FF3-1 cipher
func NewFF3Cipher(key, tweak []byte, radix int) *FF3Cipher {
	return &FF3Cipher{
		key:   key,
		tweak: tweak,
		radix: radix,
	}
}

// Encrypt encrypts plaintext using FF3-1
func (f *FF3Cipher) Encrypt(plaintext string) (string, error) {
	if len(plaintext) == 0 {
		return "", fmt.Errorf("plaintext cannot be empty")
	}

	// Validate that all characters are within the radix
	for _, char := range plaintext {
		digit := int(char - '0')
		if digit < 0 || digit >= f.radix {
			return "", fmt.Errorf("character %c is not valid for radix %d", char, f.radix)
		}
	}

	// Simplified FF3-1 implementation
	// In a real implementation, this would follow the NIST SP 800-38G specification
	return f.simplifiedEncrypt(plaintext)
}

// Decrypt decrypts ciphertext using FF3-1
func (f *FF3Cipher) Decrypt(ciphertext string) (string, error) {
	if len(ciphertext) == 0 {
		return "", fmt.Errorf("ciphertext cannot be empty")
	}

	// Validate that all characters are within the radix
	for _, char := range ciphertext {
		digit := int(char - '0')
		if digit < 0 || digit >= f.radix {
			return "", fmt.Errorf("character %c is not valid for radix %d", char, f.radix)
		}
	}

	// Simplified FF3-1 implementation
	// In a real implementation, this would follow the NIST SP 800-38G specification
	return f.simplifiedDecrypt(ciphertext)
}

// simplifiedEncrypt is a simplified encryption function for demonstration
func (f *FF3Cipher) simplifiedEncrypt(plaintext string) (string, error) {
	// This is not a real FF3-1 implementation
	// It's just a placeholder to show the interface
	result := ""

	// Simple transformation for demonstration
	for i, char := range plaintext {
		// Use key and tweak to influence the transformation
		keyByte := f.key[i%len(f.key)]
		tweakByte := f.tweak[i%len(f.tweak)]

		// Simple numeric transformation
		digit := int(char - '0')
		transformed := (digit + int(keyByte) + int(tweakByte)) % f.radix
		result += fmt.Sprintf("%d", transformed)
	}

	return result, nil
}

// simplifiedDecrypt is a simplified decryption function for demonstration
func (f *FF3Cipher) simplifiedDecrypt(ciphertext string) (string, error) {
	// This is not a real FF3-1 implementation
	// It's just a placeholder to show the interface
	result := ""

	// Simple reverse transformation for demonstration
	for i, char := range ciphertext {
		// Use key and tweak to influence the transformation
		keyByte := f.key[i%len(f.key)]
		tweakByte := f.tweak[i%len(f.tweak)]

		// Simple numeric transformation
		digit := int(char - '0')
		transformed := (digit - int(keyByte) - int(tweakByte)) % f.radix
		// Handle negative results
		if transformed < 0 {
			transformed += f.radix
		}
		result += fmt.Sprintf("%d", transformed)
	}

	return result, nil
}

// FF3Domain represents a domain for FF3 encryption
type FF3Domain struct {
	Name      string
	Radix     int
	MinLength int
	MaxLength int
	Charset   string
}

// Common FF3 domains
var (
	// Digits only (0-9)
	DigitsDomain = &FF3Domain{
		Name:      "digits",
		Radix:     10,
		MinLength: 1,
		MaxLength: 36,
		Charset:   "0123456789",
	}

	// Alphanumeric (0-9, A-Z)
	AlphanumericDomain = &FF3Domain{
		Name:      "alphanumeric",
		Radix:     36,
		MinLength: 1,
		MaxLength: 36,
		Charset:   "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}

	// Lowercase hex (0-9, a-f)
	LowercaseHexDomain = &FF3Domain{
		Name:      "lowercase_hex",
		Radix:     16,
		MinLength: 1,
		MaxLength: 36,
		Charset:   "0123456789abcdef",
	}
)

// ValidateDomain checks if text is valid for a domain
func (domain *FF3Domain) Validate(text string) error {
	if len(text) < domain.MinLength {
		return fmt.Errorf("text too short: %d < %d", len(text), domain.MinLength)
	}

	if len(text) > domain.MaxLength {
		return fmt.Errorf("text too long: %d > %d", len(text), domain.MaxLength)
	}

	// Check that all characters are in the charset
	for _, char := range text {
		if stringIndex(domain.Charset, string(char)) == -1 {
			return fmt.Errorf("character %c not in domain charset", char)
		}
	}

	return nil
}

// StringToNumeral converts a string to a numeral representation
func (domain *FF3Domain) StringToNumeral(text string) []int {
	numerals := make([]int, len(text))
	for i, char := range text {
		numerals[i] = stringIndex(domain.Charset, string(char))
	}
	return numerals
}

// NumeralToString converts a numeral representation to a string
func (domain *FF3Domain) NumeralToString(numerals []int) string {
	result := ""
	for _, numeral := range numerals {
		if numeral >= 0 && numeral < len(domain.Charset) {
			result += string(domain.Charset[numeral])
		}
	}
	return result
}

// stringIndex returns the index of substring in string, or -1 if not found
func stringIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
