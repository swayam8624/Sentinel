# Key Management Service (KMS)

The KMS package provides a comprehensive key management solution for the Sentinel platform, supporting both local development and production cloud deployments.

## Overview

This package implements:

1. **Local KMS**: For development and testing
2. **Cloud KMS Integrations**: AWS KMS, GCP KMS, and Azure Key Vault
3. **Envelope Encryption**: For secure data encryption with key separation
4. **KMS Factory**: For creating KMS clients based on configuration

## Components

### 1. KeyManagementService Interface

The [KeyManagementService](file:///Users/swayamsingal/Desktop/Programming/Sentinel/sentinel/crypto/kms/kms.go#L14-L19) interface provides basic key management operations:

```go
type KeyManagementService interface {
    GenerateDataKey() ([]byte, []byte, error) // Returns (plaintext key, encrypted key, error)
    DecryptDataKey(encryptedKey []byte) ([]byte, error)
    Encrypt(plaintext, key []byte) ([]byte, error)
    Decrypt(ciphertext, key []byte) ([]byte, error)
}
```

### 2. KMSClient Interface

The [KMSClient](file:///Users/swayamsingal/Desktop/Programming/Sentinel/sentinel/crypto/kms/kms.go#L83-L92) interface provides cloud KMS operations:

```go
type KMSClient interface {
    GenerateKey(ctx context.Context, keySpec string) (KeyMetadata, error)
    Encrypt(ctx context.Context, keyID string, plaintext []byte) (EncryptOutput, error)
    Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error)
    RotateKey(ctx context.Context, keyID string) (KeyMetadata, error)
    GetKeyMetadata(ctx context.Context, keyID string) (KeyMetadata, error)
}
```

### 3. EnvelopeEncryption

The [EnvelopeEncryption](file:///Users/swayamsingal/Desktop/Programming/Sentinel/sentinel/crypto/kms/kms.go#L344-L347) struct provides envelope encryption functionality:

```go
type EnvelopeEncryption struct {
    kmsClient     KMSClient
    dataKeyLength int
}
```

## Cloud Provider Integrations

### AWS KMS

The [AWSKMSClient](file:///Users/swayamsingal/Desktop/Programming/Sentinel/sentinel/crypto/kms/aws_kms.go#L13-L16) implements the KMSClient interface for AWS KMS.

```go
// Create AWS KMS client
ctx := context.Background()
awsKMS, err := NewAWSKMSClient(ctx, "us-east-1")
if err != nil {
    log.Fatalf("Failed to create AWS KMS client: %v", err)
}

// Generate a new key
metadata, err := awsKMS.GenerateKey(ctx, "SYMMETRIC_DEFAULT")
if err != nil {
    log.Fatalf("Failed to generate key: %v", err)
}
```

### GCP KMS

The [GCPKMSClient](file:///Users/swayamsingal/Desktop/Programming/Sentinel/sentinel/crypto/kms/gcp_kms.go#L14-L16) implements the KMSClient interface for GCP KMS.

```go
// Create GCP KMS client
ctx := context.Background()
gcpKMS, err := NewGCPKMSClient(ctx, "path/to/credentials.json")
if err != nil {
    log.Fatalf("Failed to create GCP KMS client: %v", err)
}

// Generate a new key in a key ring
keyRingPath := "projects/PROJECT_ID/locations/global/keyRings/KEY_RING_ID"
metadata, err := gcpKMS.GenerateKey(ctx, keyRingPath)
if err != nil {
    log.Fatalf("Failed to generate key: %v", err)
}
```

### Azure Key Vault

The [AzureKMSClient](file:///Users/swayamsingal/Desktop/Programming/Sentinel/sentinel/crypto/kms/azure_kms.go#L14-L17) implements the KMSClient interface for Azure Key Vault.

```go
// Create Azure Key Vault client
ctx := context.Background()
vaultURL := "https://YOUR_VAULT_NAME.vault.azure.net/"
tenantID := "YOUR_TENANT_ID"
clientID := "YOUR_CLIENT_ID"
clientSecret := "YOUR_CLIENT_SECRET"

azureKMS, err := NewAzureKMSClient(ctx, vaultURL, tenantID, clientID, clientSecret)
if err != nil {
    log.Fatalf("Failed to create Azure KMS client: %v", err)
}

// Generate a new key
metadata, err := azureKMS.GenerateKey(ctx, "RSA_2048")
if err != nil {
    log.Fatalf("Failed to generate key: %v", err)
}
```

## KMS Factory

The KMS factory provides a unified way to create KMS clients based on configuration:

```go
// Local KMS configuration
localConfig := KMSConfig{
    Provider:       "local",
    LocalMasterKey: []byte("masterkey12345678901234567890123"),
}

// Cloud KMS configuration (example for AWS)
awsConfig := KMSConfig{
    Provider:  "aws",
    AWSRegion: "us-east-1",
}

ctx := context.Background()
kmsClient, err := NewKMSClient(ctx, localConfig)
if err != nil {
    log.Fatalf("Failed to create KMS client: %v", err)
}
```

## Usage Examples

### Basic Encryption/Decryption

```go
// Create a local KMS with a master key
masterKey := []byte("masterkey12345678901234567890123") // 32 bytes
kmsService := NewLocalKMS(masterKey)

// Generate a data key
plaintextKey, encryptedKey, err := kmsService.GenerateDataKey()
if err != nil {
    log.Fatalf("Failed to generate data key: %v", err)
}

// Encrypt some data
plaintext := []byte("This is a secret message!")
ciphertext, err := kmsService.Encrypt(plaintext, plaintextKey)
if err != nil {
    log.Fatalf("Failed to encrypt data: %v", err)
}

// Decrypt the data
decrypted, err := kmsService.Decrypt(ciphertext, plaintextKey)
if err != nil {
    log.Fatalf("Failed to decrypt data: %v", err)
}
```

### Envelope Encryption

```go
// Create a local KMS client for envelope encryption
localClient := NewLocalKMSClient()

// Create a key for envelope encryption
ctx := context.Background()
metadata, err := localClient.GenerateKey(ctx, "AES_256")
if err != nil {
    log.Fatalf("Failed to generate key: %v", err)
}

// Create envelope encryption service
envelope := NewEnvelopeEncryption(localClient, 32)

// Encrypt data using envelope encryption
plaintext := []byte("This is a secret message!")
ciphertext, err := envelope.EncryptWithEnvelope(ctx, metadata.KeyID, plaintext)
if err != nil {
    log.Fatalf("Failed to encrypt with envelope: %v", err)
}

// Decrypt data using envelope encryption
decrypted, err := envelope.DecryptWithEnvelope(ctx, metadata.KeyID, ciphertext)
if err != nil {
    log.Fatalf("Failed to decrypt with envelope: %v", err)
}
```

## Security Features

1. **Key Separation**: Data keys are separated from master keys
2. **Envelope Encryption**: Provides an additional layer of security
3. **Cloud Integration**: Leverages cloud provider security features
4. **Key Rotation**: Supports key rotation for enhanced security
5. **Audit Trail**: Key operations can be audited

## Testing

The package includes comprehensive unit tests:

```bash
go test ./sentinel/crypto/kms/... -v
```

## Examples

Run the examples to see the KMS in action:

```bash
# Basic KMS example
go run sentinel/crypto/kms/examples/basic/main.go

# Cloud KMS examples (documentation only)
go run sentinel/crypto/kms/examples/cloud/main.go
```
