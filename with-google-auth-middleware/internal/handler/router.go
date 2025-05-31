// Package handler implements the HTTP routing layer of the application
package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/tyobaskara/jeki-backend/internal/handler/v1"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/config"
	authhandler "github.com/tyobaskara/jeki-backend/internal/modules/auth/handler"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/middleware"
	authrepo "github.com/tyobaskara/jeki-backend/internal/modules/auth/repository"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/usecase"
	userhandler "github.com/tyobaskara/jeki-backend/internal/modules/user/handler"
	userrepo "github.com/tyobaskara/jeki-backend/internal/modules/user/repository"
	userusecase "github.com/tyobaskara/jeki-backend/internal/modules/user/usecase"
	"gorm.io/gorm"
)

// SetupRouter configures all the routes for the application
// It initializes the Gin router and registers all route handlers
// Returns:
//   - *gin.Engine: configured Gin router instance
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Auth module manual wiring
	authRepo := authrepo.NewAuthRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(
		authRepo,
		userRepo,
		usecase.AuthUsecaseConfig{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			JWTSecret:    cfg.JWTSecret,
			TokenConfig: usecase.TokenConfig{
				AccessTTL:  cfg.AccessTokenTTL,
				RefreshTTL: cfg.RefreshTokenTTL,
			},
		},
	)
	authHandler := authhandler.NewAuthHandler(authUsecase)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// User module manual wiring
	userUsecase := userusecase.NewUserUsecase(userRepo)
	userHandler := userhandler.NewUserHandler(userUsecase)

	// Setup router with handlers
	return v1.SetupRouter(userHandler, authHandler, authMiddleware)
} 