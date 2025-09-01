package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sentinel/crypto/kms"
	"sentinel/sentinel/detector"

	"go.opentelemetry.io/otel/attribute"
)

// SentinelCryptoIntegration handles the integration between KMS and Sentinel security pipeline
type SentinelCryptoIntegration struct {
	kmsClient     kms.KMSClient
	observability interface { // Using interface to avoid circular dependencies
		RecordMetric(context.Context, string, float64, ...attribute.KeyValue)
		RecordCryptoOperation(context.Context, string, time.Duration, ...attribute.KeyValue)
		RecordViolation(context.Context, string, ...attribute.KeyValue)
	}
	currentKeyID string
}

// NewSentinelCryptoIntegration creates a new Sentinel crypto integration handler
func NewSentinelCryptoIntegration(kmsClient kms.KMSClient, obs interface{}) *SentinelCryptoIntegration {
	return &SentinelCryptoIntegration{
		kmsClient:     kmsClient,
		observability: obs,
	}
}

// EnsureKey ensures a key is available for encryption operations
func (sci *SentinelCryptoIntegration) EnsureKey(ctx context.Context) error {
	if sci.currentKeyID != "" {
		// Check if key is still valid
		_, err := sci.kmsClient.GetKeyMetadata(ctx, sci.currentKeyID)
		if err == nil {
			return nil // Key is valid
		}
	}

	// Generate a new key
	metadata, err := sci.kmsClient.GenerateKey(ctx, "AES_256")
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	sci.currentKeyID = metadata.KeyID

	// Record the key generation
	if sci.observability != nil {
		sci.observability.RecordMetric(ctx, "request.count", 1,
			attribute.String("event", "key_generated"),
			attribute.String("key_id", metadata.KeyID))
	}

	return nil
}

// EncryptViolationData encrypts violation detection data for secure storage
func (sci *SentinelCryptoIntegration) EncryptViolationData(ctx context.Context, data []byte) ([]byte, error) {
	// Ensure we have a key
	if err := sci.EnsureKey(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure key: %w", err)
	}

	// Record the crypto operation
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		if sci.observability != nil {
			sci.observability.RecordCryptoOperation(ctx, "encrypt_violation", duration)
		}
	}()

	// Encrypt the data
	output, err := sci.kmsClient.Encrypt(ctx, sci.currentKeyID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt violation data: %w", err)
	}

	return output.CiphertextBlob, nil
}

// DecryptViolationData decrypts violation detection data
func (sci *SentinelCryptoIntegration) DecryptViolationData(ctx context.Context, encryptedData []byte) ([]byte, error) {
	// Ensure we have a key
	if err := sci.EnsureKey(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure key: %w", err)
	}

	// Record the crypto operation
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		if sci.observability != nil {
			sci.observability.RecordCryptoOperation(ctx, "decrypt_violation", duration)
		}
	}()

	// Decrypt the data
	plaintext, err := sci.kmsClient.Decrypt(ctx, sci.currentKeyID, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt violation data: %w", err)
	}

	return plaintext, nil
}

// LogEncryptedViolation securely logs a violation detection result
func (sci *SentinelCryptoIntegration) LogEncryptedViolation(ctx context.Context, result *detector.DetectionResult) error {
	// Convert result to JSON
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal detection result: %w", err)
	}

	// Encrypt the data
	encryptedData, err := sci.EncryptViolationData(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to encrypt violation data: %w", err)
	}

	// In a real implementation, this would be stored in a secure audit log
	// For now, we'll just record the event
	if sci.observability != nil {
		sci.observability.RecordViolation(ctx, result.ViolationType,
			attribute.String("status", "encrypted_and_logged"),
			attribute.Float64("score", result.Score),
			attribute.Float64("confidence", result.Confidence))
	}

	return nil
}

// RotateKey rotates the current encryption key
func (sci *SentinelCryptoIntegration) RotateKey(ctx context.Context) error {
	if sci.currentKeyID == "" {
		return fmt.Errorf("no current key to rotate")
	}

	// Record the crypto operation
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		if sci.observability != nil {
			sci.observability.RecordCryptoOperation(ctx, "rotate_key", duration)
		}
	}()

	// Rotate the key
	metadata, err := sci.kmsClient.RotateKey(ctx, sci.currentKeyID)
	if err != nil {
		return fmt.Errorf("failed to rotate key: %w", err)
	}

	// Update the current key ID
	sci.currentKeyID = metadata.KeyID

	// Record the key rotation
	if sci.observability != nil {
		sci.observability.RecordMetric(ctx, "request.count", 1,
			attribute.String("event", "key_rotated"),
			attribute.String("key_id", metadata.KeyID))
	}

	return nil
}
