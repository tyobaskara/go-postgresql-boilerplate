package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tyobaskara/jeki-backend/internal/modules/user/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/users")

	{
		group.GET("", h.GetAllUsers)
		group.GET("/:id", h.GetUserByID)
		group.PUT("/:id", h.UpdateUser)
		group.DELETE("/:id", h.DeleteUser)
	}
}

// GetUserByID handles getting a user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userUsecase.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles user updates
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get existing user first
	existingUser, err := h.userUsecase.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse update request
	var updateRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if updateRequest.Name != "" {
		existingUser.Name = updateRequest.Name
	}
	if updateRequest.Email != "" {
		existingUser.Email = updateRequest.Email
	}
	if updateRequest.Password != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		existingUser.Password = string(hashedPassword)
	}

	if err := h.userUsecase.UpdateUser(existingUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingUser)
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userUsecase.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllUsers handles getting all users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page := 1
	limit := 10

	// Parse query parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	users, err := h.userUsecase.GetAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
} 