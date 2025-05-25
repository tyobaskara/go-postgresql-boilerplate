// Package v1 implements version 1 of the API endpoints
package v1

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the router with all routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// User routes
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", CreateUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
	}

	return router
}
