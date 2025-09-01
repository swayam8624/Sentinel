package redaction

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// RedactionAction represents a redaction action to be performed
type RedactionAction struct {
	Type        string `json:"type"`        // "mask", "tokenize", "fpe", "encrypt", "drop"
	Format      string `json:"format"`      // Format pattern for FPE or masking
	MaskChar    string `json:"mask_char"`   // Character to use for masking
	PreserveDomain bool `json:"preserve_domain"` // Whether to preserve domain in email masking
}

// Redactor performs redaction actions on detected sensitive data
type Redactor struct {
	// Encryption key for tokenization
	encryptionKey []byte
}

// NewRedactor creates a new redactor with the provided encryption key
func NewRedactor(encryptionKey []byte) *Redactor {
	return &Redactor{
		encryptionKey: encryptionKey,
	}
}

// Redact performs the specified redaction action on the input text
func (r *Redactor) Redact(text string, action RedactionAction) (string, error) {
	switch action.Type {
	case "mask":
		return r.mask(text, action)
	case "tokenize":
		return r.tokenize(text, action)
	case "fpe":
		return r.formatPreservingEncrypt(text, action)
	case "encrypt":
		return r.encrypt(text, action)
	case "drop":
		return "", nil
	default:
		return "", fmt.Errorf("unknown redaction action: %s", action.Type)
	}
}

// mask replaces characters in the text with a mask character
func (r *Redactor) mask(text string, action RedactionAction) (string, error) {
	maskChar := "*"
	if action.MaskChar != "" {
		maskChar = action.MaskChar
	}
	
	// If preserving domain (for emails), only mask the username part
	if action.PreserveDomain && isEmail(text) {
		return maskEmail(text, maskChar)
	}
	
	// Simple masking - replace all characters with mask character
	masked := ""
	for range text {
		masked += maskChar
	}
	
	return masked, nil
}

// tokenize replaces the text with a reversible token
func (r *Redactor) tokenize(text string, action RedactionAction) (string, error) {
	// Generate a token based on the encrypted hash of the text
	token, err := r.generateToken(text)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	
	return token, nil
}

// formatPreservingEncrypt applies format-preserving encryption
func (r *Redactor) formatPreservingEncrypt(text string, action RedactionAction) (string, error) {
	// This is a simplified implementation
	// A real implementation would use FF3-1 or similar
	if action.Format == "" {
		return "", fmt.Errorf("format required for FPE")
	}
	
	// For demonstration, we'll just return a placeholder
	// In a real implementation, this would apply actual FPE
	return applyFormat(text, action.Format), nil
}

// encrypt applies standard encryption (AES-GCM)
func (r *Redactor) encrypt(text string, action RedactionAction) (string, error) {
	// Generate a new nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	// Create cipher
	block, err := aes.NewCipher(r.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Encrypt the text
	ciphertext := gcm.Seal(nil, nonce, []byte(text), nil)
	
	// Combine nonce and ciphertext for storage
	result := append(nonce, ciphertext...)
	
	// Return as base64 encoded string
	return encodeBase64(result), nil
}

// generateToken creates a reversible token for the text
func (r *Redactor) generateToken(text string) (string, error) {
	// Create cipher
	block, err := aes.NewCipher(r.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Generate a deterministic nonce based on the text
	nonce := generateDeterministicNonce(text, gcm.NonceSize())
	
	// Encrypt the text
	ciphertext := gcm.Seal(nil, nonce, []byte(text), nil)
	
	// Combine nonce and ciphertext
	result := append(nonce, ciphertext...)
	
	// Return as base64 encoded string with prefix
	return "token_" + encodeBase64(result), nil
}

// isEmail checks if the text is an email address
func isEmail(text string) bool {
	// Simple email check - contains @ and .
	return len(text) > 0 && 
		   stringContains(text, "@") && 
		   stringContains(text, ".") &&
		   stringIndex(text, "@") < stringIndex(text, ".")
}

// maskEmail masks only the username part of an email
func maskEmail(email, maskChar string) (string, error) {
	atIndex := stringIndex(email, "@")
	if atIndex == -1 {
		return "", fmt.Errorf("invalid email format")
	}
	
	username := email[:atIndex]
	domain := email[atIndex:]
	
	maskedUsername := ""
	for range username {
		maskedUsername += maskChar
	}
	
	return maskedUsername + domain, nil
}

// applyFormat applies a format pattern to the text
func applyFormat(text, format string) string {
	// This is a simplified implementation
	// A real implementation would preserve format more carefully
	result := ""
	textIndex := 0
	
	for _, char := range format {
		if char == 'N' || char == 'X' {
			// Use character from original text or '0' if exhausted
			if textIndex < len(text) {
				result += string(text[textIndex])
				textIndex++
			} else {
				result += "0"
			}
		} else {
			// Use format character
			result += string(char)
		}
	}
	
	return result
}

// generateDeterministicNonce creates a deterministic nonce for a given text
func generateDeterministicNonce(text string, size int) []byte {
	// In a real implementation, this would use a proper KDF
	// For now, we'll use a simplified approach
	nonce := make([]byte, size)
	for i := 0; i < size && i < len(text); i++ {
		nonce[i] = text[i]
	}
	
	// Pad with zeros if needed
	for i := len(text); i < size; i++ {
		nonce[i] = 0
	}
	
	return nonce
}

// encodeBase64 encodes bytes to base64 string
func encodeBase64(data []byte) string {
	// In a real implementation, we'd use encoding/base64
	// For now, we'll use a placeholder
	return fmt.Sprintf("base64_%x", data)
}

// stringContains checks if a string contains a substring
func stringContains(s, substr string) bool {
	return stringIndex(s, substr) != -1
}

// stringIndex returns the index of substring in string, or -1 if not found
func stringIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}package redaction
