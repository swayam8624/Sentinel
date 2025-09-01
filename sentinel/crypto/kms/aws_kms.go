package kms

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// AWSKMSClient implements KMSClient for AWS KMS
type AWSKMSClient struct {
	client *kms.Client
	region string
}

// NewAWSKMSClient creates a new AWS KMS client
func NewAWSKMSClient(ctx context.Context, region string) (*AWSKMSClient, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create KMS client
	client := kms.NewFromConfig(cfg)

	return &AWSKMSClient{
		client: client,
		region: region,
	}, nil
}

// GenerateKey generates a new key in AWS KMS
func (ak *AWSKMSClient) GenerateKey(ctx context.Context, keySpec string) (KeyMetadata, error) {
	// Create key input
	input := &kms.CreateKeyInput{
		Description: aws.String("Sentinel Crypto Key"),
		KeyUsage:    types.KeyUsageTypeEncryptDecrypt,
		Origin:      types.OriginTypeAwsKms,
	}

	// Set key spec if provided
	if keySpec != "" {
		input.KeySpec = types.KeySpec(keySpec)
	} else {
		input.KeySpec = types.KeySpecSymmetricDefault
	}

	// Create the key
	result, err := ak.client.CreateKey(ctx, input)
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to create key: %w", err)
	}

	// Convert to our metadata format
	metadata := KeyMetadata{
		KeyID:           *result.KeyMetadata.KeyId,
		Arn:             *result.KeyMetadata.Arn,
		CreationDate:    *result.KeyMetadata.CreationDate,
		Description:     *result.KeyMetadata.Description,
		Enabled:         result.KeyMetadata.Enabled,
		ExpirationModel: string(result.KeyMetadata.ExpirationModel),
		KeyManager:      string(result.KeyMetadata.KeyManager),
		KeyState:        string(result.KeyMetadata.KeyState),
		KeyUsage:        string(result.KeyMetadata.KeyUsage),
		Origin:          string(result.KeyMetadata.Origin),
	}

	return metadata, nil
}

// Encrypt encrypts data with a key in AWS KMS
func (ak *AWSKMSClient) Encrypt(ctx context.Context, keyID string, plaintext []byte) (EncryptOutput, error) {
	// Create encrypt input
	input := &kms.EncryptInput{
		KeyId:     aws.String(keyID),
		Plaintext: plaintext,
	}

	// Encrypt the data
	result, err := ak.client.Encrypt(ctx, input)
	if err != nil {
		return EncryptOutput{}, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Return the result
	output := EncryptOutput{
		CiphertextBlob: result.CiphertextBlob,
		KeyID:          *result.KeyId,
	}

	return output, nil
}

// Decrypt decrypts data with a key in AWS KMS
func (ak *AWSKMSClient) Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	// Create decrypt input
	input := &kms.DecryptInput{
		KeyId:          aws.String(keyID),
		CiphertextBlob: ciphertext,
	}

	// Decrypt the data
	result, err := ak.client.Decrypt(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return result.Plaintext, nil
}

// RotateKey rotates a key in AWS KMS
func (ak *AWSKMSClient) RotateKey(ctx context.Context, keyID string) (KeyMetadata, error) {
	// Enable key rotation
	_, err := ak.client.EnableKeyRotation(ctx, &kms.EnableKeyRotationInput{
		KeyId: aws.String(keyID),
	})
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to enable key rotation: %w", err)
	}

	// Get key metadata
	result, err := ak.client.DescribeKey(ctx, &kms.DescribeKeyInput{
		KeyId: aws.String(keyID),
	})
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to describe key: %w", err)
	}

	// Convert to our metadata format
	metadata := KeyMetadata{
		KeyID:           *result.KeyMetadata.KeyId,
		Arn:             *result.KeyMetadata.Arn,
		CreationDate:    *result.KeyMetadata.CreationDate,
		Description:     *result.KeyMetadata.Description,
		Enabled:         result.KeyMetadata.Enabled,
		ExpirationModel: string(result.KeyMetadata.ExpirationModel),
		KeyManager:      string(result.KeyMetadata.KeyManager),
		KeyState:        string(result.KeyMetadata.KeyState),
		KeyUsage:        string(result.KeyMetadata.KeyUsage),
		Origin:          string(result.KeyMetadata.Origin),
	}

	return metadata, nil
}

// GetKeyMetadata gets key metadata from AWS KMS
func (ak *AWSKMSClient) GetKeyMetadata(ctx context.Context, keyID string) (KeyMetadata, error) {
	// Get key metadata
	result, err := ak.client.DescribeKey(ctx, &kms.DescribeKeyInput{
		KeyId: aws.String(keyID),
	})
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to describe key: %w", err)
	}

	// Convert to our metadata format
	metadata := KeyMetadata{
		KeyID:           *result.KeyMetadata.KeyId,
		Arn:             *result.KeyMetadata.Arn,
		CreationDate:    *result.KeyMetadata.CreationDate,
		Description:     *result.KeyMetadata.Description,
		Enabled:         result.KeyMetadata.Enabled,
		ExpirationModel: string(result.KeyMetadata.ExpirationModel),
		KeyManager:      string(result.KeyMetadata.KeyManager),
		KeyState:        string(result.KeyMetadata.KeyState),
		KeyUsage:        string(result.KeyMetadata.KeyUsage),
		Origin:          string(result.KeyMetadata.Origin),
	}

	return metadata, nil
}
