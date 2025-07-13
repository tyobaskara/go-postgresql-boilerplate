package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tyobaskara/maxwash-backend/internal/modules/user/domain"
)

type AuthMiddleware struct {
	jwtSecret []byte
	userRepo  domain.UserRepository
}

func NewAuthMiddleware(jwtSecret string, userRepo domain.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		userRepo:  userRepo,
	}
}

// AuthRequired is a middleware that checks for a valid JWT token
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check if token is expired
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					c.Abort()
					return
				}
			}

			// Get user ID from claims
			if sub, ok := claims["sub"].(string); ok {
				userID, err := uuid.Parse(sub)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
					c.Abort()
					return
				}

				// Get user from database to check logout timestamp
				user, err := m.userRepo.FindByID(userID)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
					c.Abort()
					return
				}

				// Check if token was issued before last logout (IMMEDIATE INVALIDATION)
				if !user.LastLogoutAt.IsZero() {
					if iat, ok := claims["iat"].(float64); ok {
						tokenIssuedAt := time.Unix(int64(iat), 0)
						if user.LastLogoutAt.After(tokenIssuedAt) {
							c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalidated by logout"})
							c.Abort()
							return
						}
					}
				}

				c.Set("userId", userID)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
	}
} 