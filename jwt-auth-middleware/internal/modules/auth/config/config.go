package config

import "time"

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret         string
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
}

func NewConfig(
	googleClientID string,
	googleClientSecret string,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Config {
	return &Config{
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		JWTSecret:         jwtSecret,
		AccessTokenTTL:    accessTokenTTL,
		RefreshTokenTTL:   refreshTokenTTL,
	}
} 