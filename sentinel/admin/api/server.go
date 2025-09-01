package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"

	"github.com/sentinel-platform/sentinel/sentinel/admin"
	"github.com/sentinel-platform/sentinel/sentinel/ciphermesh/detectors"
	"github.com/sentinel-platform/sentinel/sentinel/crypto/kms"
	"github.com/sentinel-platform/sentinel/sentinel/policy"
)

// Server represents the admin API server
type Server struct {
	engine        *gin.Engine
	observability *admin.ObservabilityManager
	handler       *APIHandler
	httpServer    *http.Server
}

// NewServer creates a new admin API server
func NewServer(
	obs *admin.ObservabilityManager,
	policyEngine *policy.Engine,
	kmsClient kms.KMSClient,
	detectorMgr *detectors.DetectorManager,
) *Server {
	// Set gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Create gin engine
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Create API handler
	handler := NewAPIHandler(obs, policyEngine, kmsClient, detectorMgr)

	return &Server{
		engine:        r,
		observability: obs,
		handler:       handler,
	}
}

// SetupRoutes sets up the API routes
func (s *Server) SetupRoutes() {
	// Health check endpoint
	s.engine.GET("/health", s.handler.GetHealthCheck)

	// Policy management endpoints
	policies := s.engine.Group("/policies")
	{
		policies.GET("", s.handler.GetAllPolicies)
		policies.POST("", s.handler.CreatePolicy)
		policies.GET("/:id", s.handler.GetPolicy)
		policies.PUT("/:id", s.handler.UpdatePolicy)
		policies.DELETE("/:id", s.handler.DeletePolicy)
	}

	// Key management endpoints
	keys := s.engine.Group("/keys")
	{
		keys.POST("", s.handler.GenerateKey)
	}

	// Metrics endpoint
	s.engine.GET("/metrics", s.handler.GetMetrics)

	// Audit logs endpoint
	s.engine.GET("/audit", s.handler.GetAuditLogs)
}

// Start starts the admin API server
func (s *Server) Start(addr string) error {
	// Setup routes
	s.SetupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting admin API server on %s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start admin API server: %v", err)
		}
	}()

	// Record server start in observability
	ctx := context.Background()
	s.observability.RecordMetric(ctx, "request.count", 1,
		attribute.String("event", "server_start"),
		attribute.String("address", addr))

	return nil
}

// Stop stops the admin API server
func (s *Server) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown admin API server: %w", err)
	}

	log.Println("Admin API server stopped")
	return nil
}
