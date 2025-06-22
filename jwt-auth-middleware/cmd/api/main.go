package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tyobaskara/jeki-backend/internal/config"
	v1 "github.com/tyobaskara/jeki-backend/internal/handler/v1"
	authconfig "github.com/tyobaskara/jeki-backend/internal/modules/auth/config"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/handler"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/middleware"
	authrepo "github.com/tyobaskara/jeki-backend/internal/modules/auth/repository"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/usecase"
	userhandler "github.com/tyobaskara/jeki-backend/internal/modules/user/handler"
	userrepo "github.com/tyobaskara/jeki-backend/internal/modules/user/repository"
	userusecase "github.com/tyobaskara/jeki-backend/internal/modules/user/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Convert config to auth config
	authCfg := authconfig.NewConfig(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.JWTSecret,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)

	// Auth module manual wiring
	authRepo := authrepo.NewAuthRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(
		authRepo,
		userRepo,
		usecase.AuthUsecaseConfig{
			ClientID:     authCfg.GoogleClientID,
			ClientSecret: authCfg.GoogleClientSecret,
			JWTSecret:    authCfg.JWTSecret,
			TokenConfig: usecase.TokenConfig{
				AccessTTL:  authCfg.AccessTokenTTL,
				RefreshTTL: authCfg.RefreshTokenTTL,
			},
		},
	)
	authHandler := handler.NewAuthHandler(authUsecase)
	authMiddleware := middleware.NewAuthMiddleware(authCfg.JWTSecret)

	// User module manual wiring
	userUsecase := userusecase.NewUserUsecase(userRepo)
	userHandler := userhandler.NewUserHandler(userUsecase)

	// Initialize router
	router := v1.SetupRouter(userHandler, authHandler, authMiddleware)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s in %s environment", serverAddr, cfg.Environment)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	return db, nil
} 