package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port     int    `mapstructure:"port"`
		LogLevel string `mapstructure:"logLevel"`
	} `mapstructure:"server"`
	Database struct {
		URL               string        `mapstructure:"url"`
		MaxConnections    int           `mapstructure:"maxConnections"`
		ConnectionTimeout time.Duration `mapstructure:"connectionTimeout"`
	} `mapstructure:"database"`
	Redis struct {
		URL          string `mapstructure:"url"`
		MaxRetries   int    `mapstructure:"maxRetries"`
		MinIdleConns int    `mapstructure:"minIdleConns"`
	} `mapstructure:"redis"`
}

func main() {
	// Initialize configuration
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register routes
	registerRoutes(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// initConfig initializes the application configuration
func initConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.logLevel", "info")
	viper.SetDefault("database.maxConnections", 20)
	viper.SetDefault("database.connectionTimeout", "30s")
	viper.SetDefault("redis.maxRetries", 3)
	viper.SetDefault("redis.minIdleConns", 5)

	// Set config file name and paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Enable environment variable support
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// registerRoutes registers all HTTP routes
func registerRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Sentinel Gateway is running",
		})
	})

	// Version endpoint
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "0.1.0",
			"commit":  "dev",
			"build":   time.Now().Format(time.RFC3339),
		})
	})

	// OpenAI-compatible chat completions endpoint
	router.POST("/v1/chat/completions", handleChatCompletions)

	// Admin endpoints
	admin := router.Group("/sentinel/admin")
	{
		admin.GET("/policies", handleGetPolicies)
		admin.POST("/policies", handleCreatePolicy)
		admin.GET("/tenants", handleGetTenants)
		admin.POST("/tenants", handleCreateTenant)
		admin.GET("/logs", handleGetLogs)
	}
}

// handleChatCompletions handles the OpenAI-compatible chat completions endpoint
func handleChatCompletions(c *gin.Context) {
	// Get tenant from header
	tenant := c.GetHeader("X-Tenant")
	if tenant == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "X-Tenant header is required",
		})
		return
	}

	// TODO: Implement CipherMesh redaction
	// TODO: Implement Sentinel security checks
	// TODO: Forward to LLM provider
	// TODO: Implement response processing

	c.JSON(http.StatusOK, gin.H{
		"id":      "chatcmpl-example",
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   "sentinel-guarded-model",
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": "This is a placeholder response from Sentinel Gateway. Implementation pending.",
				},
				"finish_reason": "stop",
			},
		},
		"usage": map[string]interface{}{
			"prompt_tokens":     0,
			"completion_tokens": 0,
			"total_tokens":      0,
		},
	})
}

// handleGetPolicies handles getting policies
func handleGetPolicies(c *gin.Context) {
	// TODO: Implement policy retrieval
	c.JSON(http.StatusOK, []map[string]interface{}{})
}

// handleCreatePolicy handles creating a policy
func handleCreatePolicy(c *gin.Context) {
	// TODO: Implement policy creation
	c.JSON(http.StatusCreated, map[string]interface{}{})
}

// handleGetTenants handles getting tenants
func handleGetTenants(c *gin.Context) {
	// TODO: Implement tenant retrieval
	c.JSON(http.StatusOK, []map[string]interface{}{})
}

// handleCreateTenant handles creating a tenant
func handleCreateTenant(c *gin.Context) {
	// TODO: Implement tenant creation
	c.JSON(http.StatusCreated, map[string]interface{}{})
}

// handleGetLogs handles getting logs
func handleGetLogs(c *gin.Context) {
	// TODO: Implement log retrieval
	c.JSON(http.StatusOK, []map[string]interface{}{})
}
