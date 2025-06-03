# Authentication Module

This module provides authentication functionality using Google OAuth 2.0 (ID token for mobile) and JWT tokens.

## Features

- Google OAuth 2.0 authentication (ID token/mobile flow)
- JWT token-based session management
- Refresh token mechanism
- Session timeout
- Secure logout

## Configuration

The module requires the following environment variables:

```env
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
JWT_SECRET=your_jwt_secret
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=7d
```

## Database Migrations

### Using Makefile (Recommended)

The easiest way to run migrations is using the Makefile commands from the root project directory:

```bash
# Run migrations locally
make migrate-up-local

# Rollback migrations locally
make migrate-down-local

# Run migrations in Docker
make migrate-up-docker

# Rollback migrations in Docker
make migrate-down-docker
```

You can customize the database connection by overriding the default values:

```bash
# Example with custom database configuration
DB_HOST=localhost \
DB_PORT=5432 \
DB_NAME=jeki \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_SSL_MODE=disable \
make migrate-up-local
```

### Manual Migration

If you prefer to run migrations manually, you can use the migrate CLI directly:

```bash
# Local database
migrate -path internal/modules/auth/repository/migrations \
        -database "postgres://postgres:postgres@localhost:5432/jeki?sslmode=disable" \
        up

# Docker database
docker run --network jeki-network \
    -v $(pwd)/internal/modules/auth/repository/migrations:/migrations \
    migrate/migrate \
    -path /migrations \
    -database "postgres://postgres:postgres@db:5432/jeki?sslmode=disable" \
    up
```

The migration will create the following table structure:
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

## API Endpoints

### Google OAuth Login

```http
POST /v1/auth/google?code={authorization_code}
```

Response:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "refresh_token_here",
    "expires_at": "2024-03-21T12:00:00Z"
}
```

### Refresh Token

```http
POST /v1/auth/refresh?refresh_token={refresh_token}
```

Response:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "refresh_token_here",
    "expires_at": "2024-03-21T12:00:00Z"
}
```

### Logout

```http
POST /v1/auth/logout
Authorization: Bearer {access_token}
```

Response:
```json
{
    "message": "Successfully logged out"
}
```

## Usage

1. Initialize the module in your main application:

```go
authHandler, authMiddleware := auth.InitializeAuthModule(db, authConfig)
```

2. Set up the routes:

```go
// Create v1 router group
v1 := router.Group("/v1")

// Register auth routes
authHandler.RegisterRoutes(v1)

// Add auth middleware to v1 group
v1.Use(authMiddleware.AuthRequired())
{
    // Your protected routes here
    // Example: /v1/users
}
```

## Security Considerations

1. Always use HTTPS in production
2. Keep your JWT secret secure and rotate it periodically
3. Set appropriate token expiration times
4. Implement rate limiting for auth endpoints
5. Monitor for suspicious activity
6. Keep dependencies up to date

## Implementation Details

### Directory Structure

```
internal/modules/auth/
├── config/         # Configuration (JWT secret, OAuth settings)
├── handler/        # HTTP handlers for auth endpoints
├── middleware/     # JWT validation middleware
├── repository/     # Database operations
├── usecase/        # Business logic
└── domain/         # Interfaces and models
```

### Key Components

1. **AuthHandler**
   - Handles HTTP requests for auth endpoints
   - Implements Google OAuth login
   - Manages token refresh and logout

2. **AuthMiddleware**
   - Validates JWT tokens
   - Extracts user information
   - Protects routes

3. **AuthUsecase**
   - Implements business logic
   - Manages token generation
   - Handles session management

4. **AuthRepository**
   - Manages session storage
   - Handles database operations

## Usage Example

```go
// Initialize auth module
authConfig := config.NewConfig(
    googleClientID,
    googleClientSecret,
    googleRedirectURL,
    jwtSecret,
    accessTokenTTL,
    refreshTokenTTL,
)

// Setup dependencies
authRepo := repository.NewAuthRepository(db)
userRepo := userrepo.NewUserRepository(db)
authUsecase := usecase.NewAuthUsecase(
    authRepo,
    userRepo,
    usecase.AuthUsecaseConfig{
        ClientID:     authConfig.GoogleClientID,
        ClientSecret: authConfig.GoogleClientSecret,
        RedirectURL:  authConfig.GoogleRedirectURL,
        JWTSecret:    authConfig.JWTSecret,
        TokenConfig: usecase.TokenConfig{
            AccessTTL:  authConfig.AccessTokenTTL,
            RefreshTTL: authConfig.RefreshTokenTTL,
        },
    },
)
authHandler := handler.NewAuthHandler(authUsecase)
authMiddleware := middleware.NewAuthMiddleware(authConfig.JWTSecret)
```

## Error Handling

The module uses standard HTTP status codes:

- 200: Success
- 400: Bad Request (invalid input)
- 401: Unauthorized (invalid/missing token)
- 403: Forbidden (insufficient permissions)
- 500: Internal Server Error

Error responses follow this format:
```json
{
    "error": "Error message description"
}
```

## Testing

Run the tests using:
```bash
go test ./internal/modules/auth/...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request 