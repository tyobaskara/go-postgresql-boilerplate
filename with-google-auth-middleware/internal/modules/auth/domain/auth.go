package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AuthToken represents the JWT token structure
type AuthToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GoogleUserInfo represents the user information from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// Session represents a user's active session
type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AuthRepository defines the interface for auth data access
type AuthRepository interface {
	CreateSession(session *Session) error
	GetSessionByRefreshToken(refreshToken string) (*Session, error)
	DeleteSession(id uuid.UUID) error
	DeleteUserSessions(userID uuid.UUID) error
}

// AuthUsecase defines the interface for auth business logic
type AuthUsecase interface {
	LoginWithGoogleIDToken(ctx context.Context, idToken string) (*AuthToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthToken, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	ValidateToken(ctx context.Context, token string) (*AuthToken, error)
} 