// Package handler implements the HTTP routing layer of the application
package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/tyobaskara/jeki-backend/internal/handler/v1"
)

// SetupRouter configures all the routes for the application
// It initializes the Gin router and registers all route handlers
// Returns:
//   - *gin.Engine: configured Gin router instance
func SetupRouter() *gin.Engine {
	// Langsung gunakan router dari v1
	return v1.SetupRouter()
} 