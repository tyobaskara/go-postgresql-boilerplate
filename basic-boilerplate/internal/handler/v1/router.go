// Package v1 implements version 1 of the API endpoints
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tyobaskara/jeki-backend/internal/modules/user/handler"
)

// SetupRouter configures the router with all routes
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Register user routes
	userHandler.RegisterRoutes(router)

	return router
}
