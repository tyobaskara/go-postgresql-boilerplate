package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/domain"
	userdomain "github.com/tyobaskara/jeki-backend/internal/modules/user/domain"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// Custom errors
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token has expired")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrGoogleAuthFailed = errors.New("failed to authenticate with Google")
)

type TokenConfig struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type AuthUsecaseConfig struct {
	ClientID     string
	ClientSecret string
	JWTSecret    string
	TokenConfig  TokenConfig
}

// GoogleClient interface for mocking in tests
type GoogleClient interface {
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*domain.GoogleUserInfo, error)
}

type googleClient struct {
	config *oauth2.Config
}

func (c *googleClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*domain.GoogleUserInfo, error) {
	client := c.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}
	defer resp.Body.Close()

	var googleUser domain.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}
	return &googleUser, nil
}

type authUsecase struct {
	authRepo     domain.AuthRepository
	userRepo     userdomain.UserRepository
	jwtSecret    []byte
	accessTTL    time.Duration
	refreshTTL   time.Duration
}

func NewAuthUsecase(
	authRepo domain.AuthRepository,
	userRepo userdomain.UserRepository,
	cfg AuthUsecaseConfig,
) domain.AuthUsecase {
	return &authUsecase{
		authRepo:   authRepo,
		userRepo:   userRepo,
		jwtSecret:  []byte(cfg.JWTSecret),
		accessTTL:  cfg.TokenConfig.AccessTTL,
		refreshTTL: cfg.TokenConfig.RefreshTTL,
	}
}

func (u *authUsecase) LoginWithGoogleIDToken(ctx context.Context, idToken string) (*domain.AuthToken, error) {
	// Verify the ID token
	tokenInfo, err := u.verifyGoogleIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}

	// Find or create user
	user, err := u.userRepo.FindByEmail(tokenInfo.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &userdomain.User{
				ID:        uuid.New(),
				Email:     tokenInfo.Email,
				Name:      tokenInfo.Name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := u.userRepo.Create(user); err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}
	}

	// Generate tokens
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := u.authRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return u.createAuthToken(accessToken, refreshToken), nil
}

func (u *authUsecase) verifyGoogleIDToken(ctx context.Context, idToken string) (*domain.GoogleUserInfo, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: idToken})
	client := oauth2.NewClient(ctx, tokenSource)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo domain.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthToken, error) {
	session, err := u.authRepo.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	user, err := u.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Generate new access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return u.createAuthToken(accessToken, refreshToken), nil
}

func (u *authUsecase) Logout(ctx context.Context, userID uuid.UUID) error {
	if err := u.authRepo.DeleteUserSessions(userID); err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	return nil
}

func (u *authUsecase) ValidateToken(ctx context.Context, token string) (*domain.AuthToken, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return u.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrTokenExpired
		}
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidUserID, err)
	}

	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &domain.AuthToken{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(u.accessTTL.Seconds()),
		ExpiresAt:   time.Now().Add(u.accessTTL),
	}, nil
}

func (u *authUsecase) generateAccessToken(user *userdomain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(u.accessTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}

func (u *authUsecase) createAuthToken(accessToken, refreshToken string) *domain.AuthToken {
	return &domain.AuthToken{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(u.accessTTL.Seconds()),
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.accessTTL),
	}
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
} 