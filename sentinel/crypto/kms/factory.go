package kms

import (
	"context"
	"fmt"
)

// KMSConfig holds configuration for KMS clients
type KMSConfig struct {
	Provider          string `mapstructure:"provider"`
	AWSRegion         string `mapstructure:"awsRegion"`
	GCPCredentials    string `mapstructure:"gcpCredentials"`
	AzureVaultURL     string `mapstructure:"azureVaultUrl"`
	AzureTenantID     string `mapstructure:"azureTenantId"`
	AzureClientID     string `mapstructure:"azureClientId"`
	AzureClientSecret string `mapstructure:"azureClientSecret"`
	LocalMasterKey    []byte `mapstructure:"localMasterKey"`
}

// NewKMSClient creates a KMS client based on configuration
func NewKMSClient(ctx context.Context, config KMSConfig) (KMSClient, error) {
	switch config.Provider {
	case "aws":
		return NewAWSKMSClient(ctx, config.AWSRegion)
	case "gcp":
		return NewGCPKMSClient(ctx, config.GCPCredentials)
	case "azure":
		return NewAzureKMSClient(ctx, config.AzureVaultURL, config.AzureTenantID, config.AzureClientID, config.AzureClientSecret)
	case "local":
		return NewLocalKMSClient(), nil
	default:
		return nil, fmt.Errorf("unsupported KMS provider: %s", config.Provider)
	}
}

// NewKMSService creates a KeyManagementService based on configuration
func NewKMSService(config KMSConfig) (KeyManagementService, error) {
	switch config.Provider {
	case "local":
		if len(config.LocalMasterKey) == 0 {
			// Generate a default master key for local development
			defaultKey := make([]byte, 32)
			for i := range defaultKey {
				defaultKey[i] = byte(i)
			}
			config.LocalMasterKey = defaultKey
		}
		return NewLocalKMS(config.LocalMasterKey), nil
	default:
		// For cloud providers, we use envelope encryption with the KMS client
		return nil, fmt.Errorf("envelope encryption service not implemented for provider: %s", config.Provider)
	}
}
