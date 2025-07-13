package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AuthToken represents the JWT token structure
type AuthToken struct {
	AccessToken  string    `json:"accessToken"`
	TokenType    string    `json:"tokenType"`
	ExpiresIn    int64     `json:"expiresIn"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// GoogleUserInfo represents the user information from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verifiedEmail"`
	Name          string `json:"name"`
	GivenName     string `json:"givenName"`
	FamilyName    string `json:"familyName"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// Session represents a user's active session
type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// AuthRepository defines the interface for auth data access
type AuthRepository interface {
	CreateSession(session *Session) error
	GetSessionByRefreshToken(refreshToken string) (*Session, error)
	GetUserSessions(userID uuid.UUID) ([]*Session, error)
	UpdateSession(session *Session) error
	DeleteSession(id uuid.UUID) error
	DeleteUserSessions(userID uuid.UUID) error
}

// AuthUsecase defines the interface for auth business logic
type AuthUsecase interface {
	SignupWithPassword(ctx context.Context, email, name, password string) (*AuthToken, error)
	LoginWithPassword(ctx context.Context, email, password string) (*AuthToken, error)
	LoginWithGoogleIDToken(ctx context.Context, idToken string) (*AuthToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthToken, error)
	Logout(ctx context.Context, userID uuid.UUID, currentToken string) error
	ValidateToken(tokenString string) (*uuid.UUID, error)
} 