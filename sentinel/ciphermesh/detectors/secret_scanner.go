package detectors

import (
	"context"
	"encoding/base64"
	"math"
	"strings"
	"time"
	"unicode"
)

// SecretScannerDetector detects potential secrets and credentials
type SecretScannerDetector struct {
	name       string
	dataType   string
	confidence float64
}

// NewSecretScannerDetector creates a new secret scanner detector
func NewSecretScannerDetector() *SecretScannerDetector {
	return &SecretScannerDetector{
		name:       "secret_scanner",
		dataType:   "credentials",
		confidence: 0.95,
	}
}

// Detect identifies potential secrets in the provided text
func (ssd *SecretScannerDetector) Detect(ctx context.Context, text string) ([]DetectionResult, error) {
	var results []DetectionResult

	// Check for high entropy strings that might be secrets
	entropyResults := ssd.detectHighEntropyStrings(text)
	results = append(results, entropyResults...)

	// Check for common secret patterns
	patternResults := ssd.detectSecretPatterns(text)
	results = append(results, patternResults...)

	// Check for base64 encoded secrets
	base64Results := ssd.detectBase64Secrets(text)
	results = append(results, base64Results...)

	return results, nil
}

// detectHighEntropyStrings finds strings with high entropy that might be secrets
func (ssd *SecretScannerDetector) detectHighEntropyStrings(text string) []DetectionResult {
	var results []DetectionResult

	// Split text into potential tokens
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || strings.ContainsRune(",.!?;:()[]{}\"", r)
	})

	for _, token := range tokens {
		// Skip short tokens
		if len(token) < 8 {
			continue
		}

		// Calculate entropy
		entropy := calculateShannonEntropy(token)

		// If entropy is high enough, it might be a secret
		if entropy > 4.0 {
			// Find position of token in original text
			start := strings.Index(text, token)
			if start >= 0 {
				result := DetectionResult{
					ID:         generateID(),
					Type:       ssd.dataType,
					Subtype:    "high_entropy",
					Confidence: ssd.confidence * (entropy / 8.0), // Normalize entropy to 0-1 scale
					Start:      start,
					End:        start + len(token),
					Text:       token,
					Context:    getContext(text, start, start+len(token), 50),
					DetectedAt: ssd.getCurrentTime(),
				}
				results = append(results, result)
			}
		}
	}

	return results
}

// detectSecretPatterns finds common secret patterns
func (ssd *SecretScannerDetector) detectSecretPatterns(text string) []DetectionResult {
	var results []DetectionResult

	// Common secret patterns
	patterns := []struct {
		subtype string
		pattern string
	}{
		{"api_key", `(?i)(api[_-]?key|apikey)["']?\s*[:=]\s*["']?[a-zA-Z0-9_\-]{32,}["']?`},
		{"password", `(?:password|passwd|pwd)["']?\s*[:=]\s*["']?[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]{8,}["']?`},
		{"secret", `(?:secret|token)["']?\s*[:=]\s*["']?[a-zA-Z0-9_\-]{32,}["']?`},
		{"private_key", `-----BEGIN[ A-Z]* PRIVATE KEY-----`},
		{"aws_key", `AKIA[0-9A-Z]{16}`},
		{"github_token", `ghp_[a-zA-Z0-9]{36}`},
		{"slack_token", `xox[baprs]-[0-9a-zA-Z]{10,48}`},
	}

	for _, p := range patterns {
		detector, err := NewRegexDetector(
			"secret_"+p.subtype,
			ssd.dataType,
			p.subtype,
			p.pattern,
			ssd.confidence,
			50,
		)
		if err != nil {
			continue
		}

		patternResults, _ := detector.Detect(context.Background(), text)
		results = append(results, patternResults...)
	}

	return results
}

// detectBase64Secrets finds potentially base64-encoded secrets
func (ssd *SecretScannerDetector) detectBase64Secrets(text string) []DetectionResult {
	var results []DetectionResult

	// Look for base64-like strings
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || strings.ContainsRune(",.!?;:()[]{}\"", r)
	})

	for _, token := range tokens {
		// Skip short tokens
		if len(token) < 20 {
			continue
		}

		// Check if it looks like base64
		if isBase64Like(token) {
			// Try to decode it
			if decoded, err := base64.StdEncoding.DecodeString(token); err == nil {
				// Check if decoded data has high entropy
				entropy := calculateShannonEntropy(string(decoded))
				if entropy > 3.0 {
					// Find position of token in original text
					start := strings.Index(text, token)
					if start >= 0 {
						result := DetectionResult{
							ID:         generateID(),
							Type:       ssd.dataType,
							Subtype:    "base64_secret",
							Confidence: ssd.confidence * (entropy / 8.0),
							Start:      start,
							End:        start + len(token),
							Text:       token,
							Context:    getContext(text, start, start+len(token), 50),
							DetectedAt: ssd.getCurrentTime(),
						}
						results = append(results, result)
					}
				}
			}
		}
	}

	return results
}

// GetType returns the type of data this detector identifies
func (ssd *SecretScannerDetector) GetType() string {
	return ssd.dataType
}

// GetName returns the name of this detector
func (ssd *SecretScannerDetector) GetName() string {
	return ssd.name
}

// getCurrentTime returns the current time (separated for testing)
func (ssd *SecretScannerDetector) getCurrentTime() time.Time {
	return time.Now()
}

// calculateShannonEntropy calculates the Shannon entropy of a string
func calculateShannonEntropy(input string) float64 {
	if len(input) == 0 {
		return 0
	}

	// Count frequency of each character
	freqMap := make(map[rune]float64)
	for _, char := range input {
		freqMap[char]++
	}

	// Calculate entropy
	var entropy float64
	inputLen := float64(len(input))
	for _, freq := range freqMap {
		probability := freq / inputLen
		entropy -= probability * log2(probability)
	}

	return entropy
}

// log2 calculates the base-2 logarithm
func log2(x float64) float64 {
	return math.Log(x) / math.Log(2)
}

// isBase64Like checks if a string looks like base64
func isBase64Like(s string) bool {
	// Base64 strings are typically multiples of 4 characters
	if len(s)%4 != 0 {
		return false
	}

	// Check if string contains only valid base64 characters
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	for _, char := range s {
		if !strings.ContainsRune(validChars, char) {
			return false
		}
	}

	return true
}

// getContext extracts context around a detected item
func getContext(text string, start, end, window int) string {
	contextStart := start - window
	if contextStart < 0 {
		contextStart = 0
	}

	contextEnd := end + window
	if contextEnd > len(text) {
		contextEnd = len(text)
	}

	return text[contextStart:contextEnd]
}
