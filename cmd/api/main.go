package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/tyobaskara/jeki-backend/internal/config"
	v1 "github.com/tyobaskara/jeki-backend/internal/handler/v1"
	"github.com/tyobaskara/jeki-backend/internal/modules/user/handler"
	"github.com/tyobaskara/jeki-backend/internal/modules/user/repository"
	"github.com/tyobaskara/jeki-backend/internal/modules/user/usecase"
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

	// Initialize database connection
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// Initialize router
	router := v1.SetupRouter(userHandler)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s in %s environment", serverAddr, cfg.Environment)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
} 