package kms

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

// AzureKMSClient implements KMSClient for Azure Key Vault
type AzureKMSClient struct {
	client   *azkeys.Client
	vaultURL string
}

// NewAzureKMSClient creates a new Azure Key Vault client
func NewAzureKMSClient(ctx context.Context, vaultURL, tenantID, clientID, clientSecret string) (*AzureKMSClient, error) {
	// Create credential
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	// Create key client
	client, err := azkeys.NewClient(vaultURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create key client: %w", err)
	}

	return &AzureKMSClient{
		client:   client,
		vaultURL: vaultURL,
	}, nil
}

// GenerateKey generates a new key in Azure Key Vault
func (ak *AzureKMSClient) GenerateKey(ctx context.Context, keySpec string) (KeyMetadata, error) {
	// Create key parameters
	keyType := azkeys.JSONWebKeyTypeRSA
	encryptOp := azkeys.JSONWebKeyOperationEncrypt
	decryptOp := azkeys.JSONWebKeyOperationDecrypt
	keySize := int32(2048)
	
	params := azkeys.CreateKeyParameters{
		Kty: &keyType,
		KeyOps: []*azkeys.JSONWebKeyOperation{
			&encryptOp,
			&decryptOp,
		},
		KeySize: &keySize,
	}

	// Generate key name
	keyName := fmt.Sprintf("sentinel-key-%d", time.Now().Unix())

	// Create the key
	result, err := ak.client.CreateKey(ctx, keyName, params, nil)
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to create key: %w", err)
	}

	// Convert to our metadata format
	keyID := string(*result.Key.KID)
	metadata := KeyMetadata{
		KeyID:        keyID,
		Arn:          keyID,
		CreationDate: time.Now(),
		Description:  "Sentinel Crypto Key",
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "Enabled",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "AZURE_KEY_VAULT",
	}

	return metadata, nil
}

// Encrypt encrypts data with a key in Azure Key Vault
func (ak *AzureKMSClient) Encrypt(ctx context.Context, keyID string, plaintext []byte) (EncryptOutput, error) {
	// Parse key ID to get key name
	// For Azure, keyID should be the full key identifier
	keyName := getKeyNameFromID(keyID)

	// Create encrypt parameters
	algorithm := azkeys.JSONWebKeyEncryptionAlgorithmRSAOAEP
	params := azkeys.KeyOperationsParameters{
		Algorithm: &algorithm,
		Value:     plaintext,
	}

	// Encrypt the data
	result, err := ak.client.Encrypt(ctx, keyName, "", params, nil)
	if err != nil {
		return EncryptOutput{}, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Return the result
	output := EncryptOutput{
		CiphertextBlob: result.Result,
		KeyID:          keyID,
	}

	return output, nil
}

// Decrypt decrypts data with a key in Azure Key Vault
func (ak *AzureKMSClient) Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	// Parse key ID to get key name
	keyName := getKeyNameFromID(keyID)

	// Create decrypt parameters
	algorithm := azkeys.JSONWebKeyEncryptionAlgorithmRSAOAEP
	params := azkeys.KeyOperationsParameters{
		Algorithm: &algorithm,
		Value:     ciphertext,
	}

	// Decrypt the data
	result, err := ak.client.Decrypt(ctx, keyName, "", params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return result.Result, nil
}

// RotateKey rotates a key in Azure Key Vault
func (ak *AzureKMSClient) RotateKey(ctx context.Context, keyID string) (KeyMetadata, error) {
	// Parse key ID to get key name
	keyName := getKeyNameFromID(keyID)

	// Rotate the key by creating a new version
	result, err := ak.client.RotateKey(ctx, keyName, nil)
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to rotate key: %w", err)
	}

	// Convert to our metadata format
	keyIDStr := string(*result.Key.KID)
	metadata := KeyMetadata{
		KeyID:        keyIDStr,
		Arn:          keyIDStr,
		CreationDate: time.Now(),
		Description:  "Sentinel Crypto Key (Rotated)",
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "Enabled",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "AZURE_KEY_VAULT",
	}

	return metadata, nil
}

// GetKeyMetadata gets key metadata from Azure Key Vault
func (ak *AzureKMSClient) GetKeyMetadata(ctx context.Context, keyID string) (KeyMetadata, error) {
	// Parse key ID to get key name
	keyName := getKeyNameFromID(keyID)

	// Get key metadata
	result, err := ak.client.GetKey(ctx, keyName, "", nil)
	if err != nil {
		return KeyMetadata{}, fmt.Errorf("failed to get key metadata: %w", err)
	}

	// Convert to our metadata format
	keyIDStr := string(*result.Key.KID)
	metadata := KeyMetadata{
		KeyID:        keyIDStr,
		Arn:          keyIDStr,
		CreationDate: time.Now(),
		Description:  "Sentinel Crypto Key",
		Enabled:      true,
		KeyManager:   "CUSTOMER",
		KeyState:     "Enabled",
		KeyUsage:     "ENCRYPT_DECRYPT",
		Origin:       "AZURE_KEY_VAULT",
	}

	return metadata, nil
}

// getKeyNameFromID extracts key name from full key identifier
func getKeyNameFromID(keyID string) string {
	// Simple implementation - in practice, you'd parse the URL properly
	// This is just for demonstration
	lastSlash := -1
	for i := len(keyID) - 1; i >= 0; i-- {
		if keyID[i] == '/' {
			lastSlash = i
			break
		}
	}

	if lastSlash != -1 && lastSlash < len(keyID)-1 {
		return keyID[lastSlash+1:]
	}

	return keyID
}