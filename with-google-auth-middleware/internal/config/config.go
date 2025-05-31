// Package config provides configuration management for the application
package config

import (
	"fmt"  // Package fmt implements formatted I/O with functions similar to C's printf and scanf
	"os"   // Package os provides a platform-independent interface to operating system functionality
	"sync" // Package sync provides basic synchronization primitives such as mutual exclusion locks
	"time" // Package time provides functionality for measuring and displaying time

	"github.com/joho/godotenv" // Package godotenv loads environment variables from .env files
)

// Config holds all configuration for the application
// This struct defines the structure of our application's configuration
type Config struct {
	Environment        string        // The current environment (e.g., "development", "production")
	ServerPort         string        // The port number where the server will listen
	DBHost             string        // Database host address
	DBPort             string        // Database port number
	DBUser             string        // Database username
	DBPassword         string        // Database password
	DBName             string        // Database name
	GoogleClientID     string        // Google OAuth client ID
	GoogleClientSecret string        // Google OAuth client secret
	GoogleRedirectURL  string        // Google OAuth redirect URL
	JWTSecret          string        // JWT secret key
	AccessTokenTTL     time.Duration // Access token time to live
	RefreshTokenTTL    time.Duration // Refresh token time to live
	// Add other configuration fields as needed
}

// Global variables for singleton pattern implementation
var (
	cfg  *Config        // The single instance of Config that will be used throughout the application
	once sync.Once      // sync.Once ensures that the initialization code runs only once
	mu   sync.RWMutex   // RWMutex provides mutual exclusion lock with reader/writer semantics
)

// LoadConfig loads the configuration based on the environment.
// It uses singleton pattern with thread safety.
// Parameters:
//   - env: string representing the environment (e.g., "development", "production")
// Returns:
//   - *Config: pointer to the configuration struct
//   - error: any error that occurred during loading
func LoadConfig(env string) (*Config, error) {
	var err error
	// sync.Once ensures that the initialization code runs only once, even if called multiple times
	once.Do(func() {
		// Load the appropriate .env file based on environment
		envFile := fmt.Sprintf(".env.%s", env)
		if err = godotenv.Load(envFile); err != nil {
			// Try to load default .env file if environment specific file is not found
			if err = godotenv.Load(); err != nil {
				// Jika file .env tidak ditemukan, abaikan error dan lanjutkan
				// karena environment variable sudah di-set oleh Docker Compose
				err = nil
			}
		}

		// Initialize the configuration with values from environment variables
		// getEnv function will use default values if environment variables are not set
		cfg = &Config{
			Environment:        env,
			ServerPort:         getEnv("SERVER_PORT", "8080"),
			DBHost:             getEnv("DB_HOST", "localhost"),
			DBPort:             getEnv("DB_PORT", "5432"),
			DBUser:             getEnv("DB_USER", "postgres"),
			DBPassword:         getEnv("DB_PASSWORD", ""),
			DBName:             getEnv("DB_NAME", "jeki"),
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
			JWTSecret:          getEnv("JWT_SECRET", ""),
			AccessTokenTTL:     time.Duration(getEnvAsInt("ACCESS_TOKEN_TTL", 15)) * time.Minute,
			RefreshTokenTTL:    time.Duration(getEnvAsInt("REFRESH_TOKEN_TTL", 7*24)) * time.Hour,
		}

		// Validate the configuration
		if err = cfg.validate(); err != nil {
			err = fmt.Errorf("configuration validation failed: %v", err)
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetConfig returns the current configuration in a thread-safe manner
// This function uses a read lock to allow multiple concurrent reads
func GetConfig() *Config {
	mu.RLock()         // Acquire a read lock
	defer mu.RUnlock() // Ensure the lock is released when the function returns
	return cfg
}

// getEnv gets an environment variable or returns a default value
// Parameters:
//   - key: string representing the environment variable name
//   - defaultValue: string to return if the environment variable is not set
// Returns:
//   - string: the value of the environment variable or the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key) // Get the value of the environment variable
	if value == "" {
		return defaultValue // Return default value if environment variable is not set
	}
	return value
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

// validate performs validation on the configuration
// This method checks if required configuration values are set
// Returns:
//   - error: any validation error that occurred
func (c *Config) validate() error {
	// Check if database password is set
	if c.DBPassword == "" {
		return fmt.Errorf("database password is required")
	}
	// Check if database name is set
	if c.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
} 