package crypto

import (
	"context"
	"fmt"
	"time"

	"sentinel/ciphermesh/detectors"

	"go.opentelemetry.io/otel/attribute"
)

// CryptoIntegration handles the integration between crypto components and CipherMesh
type CryptoIntegration struct {
	fpe           *FF3Cipher
	observability interface { // Using interface to avoid circular dependencies
		RecordMetric(context.Context, string, float64, ...attribute.KeyValue)
		RecordCryptoOperation(context.Context, string, time.Duration, ...attribute.KeyValue)
	}
}

// NewCryptoIntegration creates a new crypto integration handler
func NewCryptoIntegration(fpe *FF3Cipher, obs interface{}) *CryptoIntegration {
	return &CryptoIntegration{
		fpe:           fpe,
		observability: obs,
	}
}

// TokenizeDetectionResult applies format-preserving encryption to a detection result
func (ci *CryptoIntegration) TokenizeDetectionResult(ctx context.Context, result detectors.DetectionResult, tweak []byte) (string, error) {
	// Record the crypto operation
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		if ci.observability != nil {
			ci.observability.RecordCryptoOperation(ctx, "fpe_tokenize", duration,
				attribute.String("data_type", result.Type),
				attribute.String("data_subtype", result.Subtype))
		}
	}()

	// Apply FPE to the detected text
	encrypted, err := ci.fpe.Encrypt(result.Text, tweak)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt detection result: %w", err)
	}

	return encrypted, nil
}

// DetokenizeData applies format-preserving decryption to tokenized data
func (ci *CryptoIntegration) DetokenizeData(ctx context.Context, tokenized string, tweak []byte) (string, error) {
	// Record the crypto operation
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		if ci.observability != nil {
			ci.observability.RecordCryptoOperation(ctx, "fpe_detokenize", duration)
		}
	}()

	// Apply FPE decryption
	decrypted, err := ci.fpe.Decrypt(tokenized, tweak)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt tokenized data: %w", err)
	}

	return decrypted, nil
}

// ProcessDetectionResults applies tokenization to all detection results
func (ci *CryptoIntegration) ProcessDetectionResults(ctx context.Context, results []detectors.DetectionResult, tweak []byte) (map[string]string, error) {
	tokenizedData := make(map[string]string)

	for _, result := range results {
		// Only tokenize certain types of data (e.g., PII, PCI)
		if result.Type == "pii" || result.Type == "pci" || result.Type == "phi" {
			tokenized, err := ci.TokenizeDetectionResult(ctx, result, tweak)
			if err != nil {
				return nil, fmt.Errorf("failed to tokenize result %s: %w", result.ID, err)
			}

			tokenizedData[result.ID] = tokenized
		}
	}

	return tokenizedData, nil
}
