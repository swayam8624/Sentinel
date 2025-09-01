package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sentinel-platform/sentinel/sentinel/admin"
	"github.com/sentinel-platform/sentinel/sentinel/admin/api"
	"github.com/sentinel-platform/sentinel/sentinel/ciphermesh/detectors"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/kms"
	"github.com/sentinel-platform/sentinel/sentinel/policy"
)

func main() {
	// Create observability manager
	obs, err := admin.NewObservabilityManager()
	if err != nil {
		log.Fatalf("Failed to create observability manager: %v", err)
	}
	defer obs.Shutdown(context.Background())

	// Create a simple local KMS client for demonstration
	kmsClient := kms.NewLocalKMSClient()

	// Create mock components (in a real implementation, these would be properly initialized)
	policyEngine := policy.NewEngine()
	detectorMgr := detectors.NewDetectorManager()

	// Create admin API server
	server := api.NewServer(obs, policyEngine, kmsClient, detectorMgr)

	// Start the server
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Failed to start admin API server: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server gracefully
	log.Println("Shutting down admin API server...")

	if err := server.Stop(); err != nil {
		log.Fatalf("Failed to stop admin API server: %v", err)
	}

	log.Println("Admin API server exited")
}
