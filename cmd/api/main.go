package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tyobaskara/jeki-backend/internal/config"
	"github.com/tyobaskara/jeki-backend/internal/handler"
)

func main() {
	// Get environment from ENV variable, default to "dev"
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load configuration
	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize router
	router := handler.SetupRouter()

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s in %s environment", serverAddr, cfg.Environment)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 