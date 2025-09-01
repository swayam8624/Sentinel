package kms

import (
	"context"
	"fmt"
	"time"

	gcpkms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/api/option"
)

// GCPKMSClient implements KMSClient for GCP KMS
type GCPKMSClient struct {
	client *gcpkms.KeyManagementClient
}

// NewGCPKMSClient creates a new GCP KMS client
func NewGCPKMSClient(ctx context.Context, credentialsFile string) (*GCPKMSClient, error) {
	// Create KMS client
	var client *gcpkms.KeyManagementClient
	var err error

	if credentialsFile != "" {
		client, err = gcpkms.NewKeyManagementClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		client, err = gcpkms.NewKeyManagementClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GCP KMS client: %w", err)
	}

	return &GCPKMSClient{
		client: client,
	}, nil
}

// GenerateKey generates a new key in GCP KMS
func (gk *GCPKMSClient) GenerateKey(ctx context.Context, keySpec string) (KeyMetadata, error) {
	// For GCP, keySpec should contain the full key name format
	// We'll parse it to get the parent and create a new key

	// Create key input
	key := &kmspb.CryptoKey{
		Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
		VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
			Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
		},
	}

	// Create the key
	result, err := gk.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      keySpec, // This should be the key ring path
		CryptoKeyId: fmt.Sprintf("sentinel-key-%d", time.Now().Unix()),
		CryptoKey:   key,
	})

	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to create key: %w", err)
	}

	// Convert to our metadata format
	metadata := KeyMetadata{
		KeyID:        result.Name,
		Arn:          result.Name,
		CreationDate: result.CreateTime.AsTime(),
		Description:  result.Purpose.String(),
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "ENABLED",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "GOOGLE_KMS",
	}

	return metadata, nil
}

// Encrypt encrypts data with a key in GCP KMS
func (gk *GCPKMSClient) Encrypt(ctx context.Context, keyID string, plaintext []byte) (EncryptOutput, error) {
	// Create encrypt input
	req := &kmspb.EncryptRequest{
		Name:      keyID,
		Plaintext: plaintext,
	}

	// Encrypt the data
	result, err := gk.client.Encrypt(ctx, req)
	if err != nil {
		return EncryptOutput{}, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Return the result
	output := EncryptOutput{
		CiphertextBlob: result.Ciphertext,
		KeyID:          keyID,
	}

	return output, nil
}

// Decrypt decrypts data with a key in GCP KMS
func (gk *GCPKMSClient) Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	// Create decrypt input
	req := &kmspb.DecryptRequest{
		Name:       keyID,
		Ciphertext: ciphertext,
	}

	// Decrypt the data
	result, err := gk.client.Decrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return result.Plaintext, nil
}

// RotateKey rotates a key in GCP KMS
func (gk *GCPKMSClient) RotateKey(ctx context.Context, keyID string) (KeyMetadata, error) {
	// Create a new key version to rotate the key
	_, err := gk.client.CreateCryptoKeyVersion(ctx, &kmspb.CreateCryptoKeyVersionRequest{
		Parent: keyID,
	})

	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to rotate key: %w", err)
	}

	// Get key metadata
	keyResult, err := gk.client.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{
		Name: keyID,
	})

	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to get key metadata: %w", err)
	}

	// Convert to our metadata format
	metadata := KeyMetadata{
		KeyID:        keyResult.Name,
		Arn:          keyResult.Name,
		CreationDate: keyResult.CreateTime.AsTime(),
		Description:  keyResult.Purpose.String(),
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "ENABLED",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "GOOGLE_KMS",
	}

	return metadata, nil
}

// GetKeyMetadata gets key metadata from GCP KMS
func (gk *GCPKMSClient) GetKeyMetadata(ctx context.Context, keyID string) (KeyMetadata, error) {
	// Get key metadata
	result, err := gk.client.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{
		Name: keyID,
	})

	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to get key metadata: %w", err)
	}

	// Convert to our metadata format
	metadata := KeyMetadata{
		KeyID:        result.Name,
		Arn:          result.Name,
		CreationDate: result.CreateTime.AsTime(),
		Description:  result.Purpose.String(),
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "ENABLED",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "GOOGLE_KMS",
	}

	return metadata, nil
}
