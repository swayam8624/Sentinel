package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sentinel-platform/sentinel/sentinel/admin"
	"github.com/sentinel-platform/sentinel/sentinel/admin/api"
	"github.com/sentinel-platform/sentinel/sentinel/admin/audit"
	"github.com/sentinel-platform/sentinel/sentinel/admin/penetration"
	"github.com/sentinel-platform/sentinel/sentinel/admin/performance"
	"github.com/sentinel-platform/sentinel/sentinel/ciphermesh/detectors"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/kms"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/merkle"
	"github.com/sentinel-platform/sentinel/sentinel/policy"
	"github.com/sentinel-platform/sentinel/sentinel/sentinel/detector"
)

func main() {
	fmt.Println("Sentinel Admin Console Demo")
	fmt.Println("==========================")

	// Create observability manager
	obs, err := admin.NewObservabilityManager()
	if err != nil {
		log.Fatalf("Failed to create observability manager: %v", err)
	}
	defer obs.Shutdown(context.Background())

	// Create Merkle tree for audit logging
	merkleTree := merkle.NewMerkleTree()

	// Create audit framework
	auditFramework := audit.NewAuditFramework(merkleTree, obs)

	// Create KMS client
	kmsClient := kms.NewLocalKMSClient()

	// Create policy engine
	policyEngine := policy.NewEngine()

	// Create detector manager
	detectorMgr := detectors.NewDetectorManager()

	// Create violation detector (mock implementation)
	violationDetector := &detector.ViolationDetector{}

	// Create performance optimizer
	perfOptimizer := performance.NewPerformanceOptimizer(obs)

	// Create penetration test framework
	penTestFramework := penetration.NewPenTestFramework(auditFramework, detectorMgr, violationDetector)

	// Log startup event
	ctx := context.Background()
	auditFramework.LogEvent(ctx, "system_startup", "admin_console",
		"Sentinel Admin Console started",
		map[string]interface{}{
			"version":    "1.0.0",
			"components": []string{"observability", "audit", "kms", "policy", "detectors"},
		}, "info")

	// Apply performance optimizations
	perfOptimizer.ApplyOptimizations()

	// Run penetration tests
	testSuite := penTestFramework.CreateStandardTestSuite()
	penTestFramework.LoadTestSuite(testSuite)

	testResults, err := penTestFramework.RunAllTests(ctx)
	if err != nil {
		log.Printf("Error running penetration tests: %v", err)
	} else {
		report := penTestFramework.GenerateTestReport(testResults)
		fmt.Println("\nPenetration Test Report:")
		fmt.Println(report)

		// Log test results
		auditFramework.LogEvent(ctx, "pen_test_completed", "penetration_testing",
			"Penetration tests completed",
			map[string]interface{}{
				"total_tests": len(testResults),
				"passed":      countPassedTests(testResults),
			}, "info")
	}

	// Create admin API server
	server := api.NewServer(obs, policyEngine, kmsClient, detectorMgr)

	// Setup and start the server
	server.SetupRoutes()

	fmt.Println("\nAdmin API server is ready to start on :8080")
	fmt.Println("Press Ctrl+C to stop the server")

	// In a real implementation, we would start the server here
	// server.Start(":8080")

	// For demo purposes, we'll just show what would be available
	fmt.Println("\nAvailable API Endpoints:")
	fmt.Println("  GET  /health     - Health check")
	fmt.Println("  GET  /policies   - List all policies")
	fmt.Println("  POST /policies   - Create a new policy")
	fmt.Println("  GET  /policies/{id} - Get a specific policy")
	fmt.Println("  PUT  /policies/{id} - Update a policy")
	fmt.Println("  DELETE /policies/{id} - Delete a policy")
	fmt.Println("  POST /keys       - Generate a new key")
	fmt.Println("  GET  /metrics    - Get system metrics")
	fmt.Println("  GET  /audit      - Get audit logs")

	// Generate optimization report
	fmt.Println("\n" + perfOptimizer.GetOptimizationReport())

	// Generate compliance report
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30) // Last 30 days
	complianceReport, err := auditFramework.GenerateComplianceReport(ctx, startTime, endTime)
	if err != nil {
		log.Printf("Error generating compliance report: %v", err)
	} else {
		fmt.Printf("\nCompliance Report:\n")
		fmt.Printf("  Period: %s to %s\n", complianceReport.PeriodStart.Format("2006-01-02"), complianceReport.PeriodEnd.Format("2006-01-02"))
		fmt.Printf("  Status: %s\n", complianceReport.ComplianceStatus)
		fmt.Printf("  Total Events: %d\n", complianceReport.Summary.TotalEvents)
		fmt.Printf("  Critical Findings: %d\n", complianceReport.Summary.CriticalFindings)
		fmt.Printf("  High Findings: %d\n", complianceReport.Summary.HighFindings)
	}

	fmt.Println("\nDemo completed. In a real implementation, the admin server would continue running.")
}

// countPassedTests counts the number of passed tests
func countPassedTests(results []penetration.TestResult) int {
	count := 0
	for _, result := range results {
		if result.Passed {
			count++
		}
	}
	return count
}
