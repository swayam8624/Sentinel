package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Cloud KMS Integration Examples ===")

	// Note: These examples require actual cloud credentials to run
	// They are provided for demonstration purposes

	// Example 1: AWS KMS (requires AWS credentials)
	fmt.Println("\n1. AWS KMS Example (requires AWS credentials):")
	awsKMSExample()

	// Example 2: GCP KMS (requires GCP credentials)
	fmt.Println("\n2. GCP KMS Example (requires GCP credentials):")
	gcpKMSExample()

	// Example 3: Azure Key Vault (requires Azure credentials)
	fmt.Println("\n3. Azure Key Vault Example (requires Azure credentials):")
	azureKMSExample()

	// Example 4: KMS Factory Usage
	fmt.Println("\n4. KMS Factory Example:")
	kmsFactoryExample()
}

func awsKMSExample() {
	fmt.Println("   To use AWS KMS, you would need to:")
	fmt.Println("   1. Set up AWS credentials (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)")
	fmt.Println("   2. Specify the AWS region")
	fmt.Println("   3. Have an existing KMS key or create one")
	fmt.Println("   ")
	fmt.Println("   Example code:")
	fmt.Println("   ctx := context.Background()")
	fmt.Println("   awsKMS, err := kms.NewAWSKMSClient(ctx, \"us-east-1\")")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to create AWS KMS client: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   ")
	fmt.Println("   // Generate a new key")
	fmt.Println("   metadata, err := awsKMS.GenerateKey(ctx, \"SYMMETRIC_DEFAULT\")")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to generate key: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   fmt.Printf(\"Generated AWS KMS key: %s\\n\", metadata.KeyID)")
}

func gcpKMSExample() {
	fmt.Println("   To use GCP KMS, you would need to:")
	fmt.Println("   1. Set up GCP credentials (GOOGLE_APPLICATION_CREDENTIALS)")
	fmt.Println("   2. Have an existing key ring or create one")
	fmt.Println("   3. Specify the key ring path")
	fmt.Println("   ")
	fmt.Println("   Example code:")
	fmt.Println("   ctx := context.Background()")
	fmt.Println("   gcpKMS, err := kms.NewGCPKMSClient(ctx, \"path/to/credentials.json\")")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to create GCP KMS client: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   ")
	fmt.Println("   // Generate a new key in a key ring")
	fmt.Println("   keyRingPath := \"projects/PROJECT_ID/locations/global/keyRings/KEY_RING_ID\"")
	fmt.Println("   metadata, err := gcpKMS.GenerateKey(ctx, keyRingPath)")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to generate key: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   fmt.Printf(\"Generated GCP KMS key: %s\\n\", metadata.KeyID)")
}

func azureKMSExample() {
	fmt.Println("   To use Azure Key Vault, you would need to:")
	fmt.Println("   1. Set up Azure credentials (tenant ID, client ID, client secret)")
	fmt.Println("   2. Have an existing key vault or create one")
	fmt.Println("   3. Specify the vault URL")
	fmt.Println("   ")
	fmt.Println("   Example code:")
	fmt.Println("   ctx := context.Background()")
	fmt.Println("   vaultURL := \"https://YOUR_VAULT_NAME.vault.azure.net/\"")
	fmt.Println("   tenantID := \"YOUR_TENANT_ID\"")
	fmt.Println("   clientID := \"YOUR_CLIENT_ID\"")
	fmt.Println("   clientSecret := \"YOUR_CLIENT_SECRET\"")
	fmt.Println("   ")
	fmt.Println("   azureKMS, err := kms.NewAzureKMSClient(ctx, vaultURL, tenantID, clientID, clientSecret)")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to create Azure KMS client: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   ")
	fmt.Println("   // Generate a new key")
	fmt.Println("   metadata, err := azureKMS.GenerateKey(ctx, \"RSA_2048\")")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to generate key: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   fmt.Printf(\"Generated Azure Key Vault key: %s\\n\", metadata.KeyID)")
}

func kmsFactoryExample() {
	fmt.Println("   Using the KMS factory to create clients based on configuration:")
	fmt.Println("   ")
	fmt.Println("   // Local KMS configuration")
	fmt.Println("   localConfig := kms.KMSConfig{")
	fmt.Println("       Provider:       \"local\",")
	fmt.Println("       LocalMasterKey: []byte(\"masterkey12345678901234567890123\"),")
	fmt.Println("   }")
	fmt.Println("   ")
	fmt.Println("   ctx := context.Background()")
	fmt.Println("   kmsClient, err := kms.NewKMSClient(ctx, localConfig)")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to create KMS client: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   ")
	fmt.Println("   // Use the client for encryption/decryption operations")
	fmt.Println("   metadata, err := kmsClient.GenerateKey(ctx, \"AES_256\")")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatalf(\"Failed to generate key: %v\", err)")
	fmt.Println("   }")
	fmt.Println("   fmt.Printf(\"Generated key: %s\\n\", metadata.KeyID)")
}
