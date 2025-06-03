package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/domain"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// LoginWithGoogle handles Google OAuth login for mobile applications only
// @Summary Login with Google (Mobile)
// @Description Authenticate user with Google OAuth using ID token
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id_token formData string true "Google ID token"
// @Success 200 {object} domain.AuthToken
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/google [post]
func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	idToken := c.PostForm("id_token")
	if idToken == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ID token is required",
		})
		return
	}

	token, err := h.authUsecase.LoginWithGoogleIDToken(c.Request.Context(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Failed to authenticate with Google",
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token query string true "Refresh token"
// @Success 200 {object} domain.AuthToken
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Refresh token is required",
		})
		return
	}

	token, err := h.authUsecase.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate all user sessions
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	if err := h.authUsecase.Logout(c.Request.Context(), userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Successfully logged out",
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// RegisterRoutes registers all auth routes
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/auth")
	{
		group.POST("/google", h.LoginWithGoogle)
		group.POST("/refresh", h.RefreshToken)
		group.POST("/logout", h.Logout)
	}
} 