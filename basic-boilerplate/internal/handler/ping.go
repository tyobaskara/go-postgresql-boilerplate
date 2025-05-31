package handler

import (
	"github.com/gin-gonic/gin"
)

// Ping handles the ping request
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
} 