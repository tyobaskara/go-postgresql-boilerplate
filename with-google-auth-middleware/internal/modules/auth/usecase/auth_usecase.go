package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/domain"
	userdomain "github.com/tyobaskara/jeki-backend/internal/modules/user/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type TokenConfig struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type AuthUsecaseConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	JWTSecret    string
	TokenConfig  TokenConfig
}

type authUsecase struct {
	authRepo    domain.AuthRepository
	userRepo    userdomain.UserRepository
	config      *oauth2.Config
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewAuthUsecase(
	authRepo domain.AuthRepository,
	userRepo userdomain.UserRepository,
	cfg AuthUsecaseConfig,
) domain.AuthUsecase {
	return &authUsecase{
		authRepo:   authRepo,
		userRepo:   userRepo,
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		jwtSecret:  []byte(cfg.JWTSecret),
		accessTTL:  cfg.TokenConfig.AccessTTL,
		refreshTTL: cfg.TokenConfig.RefreshTTL,
	}
}

func (u *authUsecase) LoginWithGoogle(code string) (*domain.AuthToken, error) {
	// Exchange code for token
	token, err := u.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// Get user info from Google
	client := u.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleUser domain.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	// Find or create user
	user, err := u.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new user
			user = &userdomain.User{
				ID:        uuid.New(),
				Email:     googleUser.Email,
				Name:      googleUser.Name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := u.userRepo.Create(user); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Create session
	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTTL),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := u.authRepo.CreateSession(session); err != nil {
		return nil, err
	}

	// Generate access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthToken{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(u.accessTTL.Seconds()),
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.accessTTL),
	}, nil
}

func (u *authUsecase) RefreshToken(refreshToken string) (*domain.AuthToken, error) {
	session, err := u.authRepo.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthToken{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(u.accessTTL.Seconds()),
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.accessTTL),
	}, nil
}

func (u *authUsecase) Logout(userID uuid.UUID) error {
	return u.authRepo.DeleteUserSessions(userID)
}

func (u *authUsecase) ValidateToken(token string) (*domain.AuthToken, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return u.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
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

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
} 