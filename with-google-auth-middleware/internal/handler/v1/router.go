// Package v1 implements version 1 of the API endpoints
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/handler"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/middleware"
	userhandler "github.com/tyobaskara/jeki-backend/internal/modules/user/handler"
)

// SetupRouter configures the router with all routes
func SetupRouter(userHandler *userhandler.UserHandler, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Register auth routes
		authHandler.RegisterRoutes(v1)

		// Protected routes
		v1.Use(authMiddleware.AuthRequired())
		{
			// Register user routes
			userHandler.RegisterRoutes(v1)
		}
	}

	return router
}
