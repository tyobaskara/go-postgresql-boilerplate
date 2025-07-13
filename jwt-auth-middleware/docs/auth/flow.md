# Authentication Flow Documentation

## Overview

Dokumen ini menjelaskan alur autentikasi yang diimplementasikan dalam sistem maxwash backend. Sistem mendukung dua jenis autentikasi:

1. **Password-based Authentication** - Login dengan email dan password
2. **Google OAuth Authentication** - Login menggunakan Google OAuth

## Arsitektur Autentikasi

### Komponen Utama

1. **AuthHandler** - Menangani HTTP requests untuk autentikasi
2. **AuthUsecase** - Logika bisnis autentikasi
3. **AuthRepository** - Akses data autentikasi
4. **AuthMiddleware** - Middleware untuk protected routes
5. **UserRepository** - Akses data user

### Flow Diagram

```
Client Request
    ↓
AuthHandler
    ↓
AuthUsecase (Business Logic)
    ↓
AuthRepository/UserRepository
    ↓
Database (PostgreSQL)
```

## Password-based Authentication

### 1. User Registration (Signup)

**Endpoint**: `POST /v1/auth/signup`

**Request Body**:
```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "securepassword123"
}
```

**Response**:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tokenType": "Bearer",
  "expiresIn": 3600,
  "refreshToken": "abc123def456ghi789...",
  "expiresAt": "2024-01-01T12:00:00Z"
}
```

**Flow**:
1. Validasi input (email, name, password)
2. Hash password menggunakan bcrypt
3. Simpan user ke database
4. Generate JWT access token dan refresh token
5. Simpan session ke database
6. Return token response

### 2. User Login

**Endpoint**: `POST /v1/auth/login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**: Same as signup response

**Flow**:
1. Validasi input (email, password)
2. Cari user berdasarkan email
3. Verifikasi password menggunakan bcrypt
4. Update lastLoginAt timestamp
5. Generate JWT access token dan refresh token
6. Simpan session ke database
7. Return token response

## Google OAuth Authentication

### 1. Google OAuth Login

**Endpoint**: `POST /v1/auth/google`

**Request Body** (form-data):
```
id_token: <google_id_token>
```

**Response**: Same as password login response

**Flow**:
1. Validasi Google ID token
2. Extract user info dari Google
3. Cari atau buat user berdasarkan email
4. Update lastLoginAt timestamp
5. Generate JWT access token dan refresh token
6. Simpan session ke database
7. Return token response

## Token Management

### 1. Token Refresh

**Endpoint**: `POST /v1/auth/refresh?refreshToken=<token>`

**Response**:
```json
{
  "accessToken": "new_accessToken_here...",
  "tokenType": "Bearer",
  "expiresIn": 3600,
  "expiresAt": "2024-01-01T12:00:00Z"
}
```

**Flow**:
1. Validasi refresh token
2. Cek apakah token masih valid di database
3. Generate access token baru
4. Update session expiresAt
5. Return new access token

### 2. Logout

**Endpoint**: `POST /v1/auth/logout`

**Headers**: `Authorization: Bearer <accessToken>`

**Response**:
```json
{
  "message": "Successfully logged out"
}
```

**Flow**:
1. Validasi access token
2. Extract user ID dari token
3. Hapus semua session user dari database
4. Blacklist current token
5. Return success message

## Protected Routes

### Middleware Implementation

Semua protected routes menggunakan `AuthMiddleware`:

```go
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from Authorization header
        // Validate JWT token
        // Set userId in context
        // Continue to handler
    }
}
```

### Protected Endpoints

- `POST /v1/auth/logout` - Logout user
- `GET /v1/users` - Get all users
- `GET /v1/users/:id` - Get user by ID
- `PUT /v1/users/:id` - Update user
- `DELETE /v1/users/:id` - Delete user

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255), -- Nullable untuk Google OAuth users
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_logout_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Sessions Table

```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## JWT Token Structure

### Access Token Payload

```json
{
  "sub": "userId",
  "exp": 1640995200,
  "iat": 1640991600,
  "type": "access"
}
```

### Refresh Token Payload

```json
{
  "sub": "userId",
  "exp": 1641078000,
  "iat": 1640991600,
  "type": "refresh",
  "sessionId": "session_uuid"
}
```

## Error Handling

### Common Error Responses

**400 Bad Request**:
```json
{
  "error": "Invalid request body: email is required"
}
```

**401 Unauthorized**:
```json
{
  "error": "Invalid credentials"
}
```

**409 Conflict**:
```json
{
  "error": "User with this email already exists"
}
```

**500 Internal Server Error**:
```json
{
  "error": "Internal server error"
}
```

## Security Considerations

### 1. Password Security
- Password di-hash menggunakan bcrypt dengan cost factor 10
- Password minimum 6 karakter
- Password tidak pernah disimpan dalam plain text

### 2. Token Security
- Access token expires dalam 1 jam
- Refresh token expires dalam 24 jam
- Token disimpan dengan aman di database
- Implementasi token blacklisting

### 3. Session Management
- Multiple sessions per user
- Automatic cleanup expired sessions
- Session invalidation pada logout

### 4. Input Validation
- Validasi email format
- Validasi password strength
- Sanitasi input data
- Rate limiting pada auth endpoints

## Configuration

### Environment Variables

```env
# JWT Configuration
JWT_SECRET=your_jwt_secret_key
ACCESS_TOKEN_TTL=1h
REFRESH_TOKEN_TTL=24h

# Google OAuth Configuration
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=maxwash
```

### Token Configuration

```go
type TokenConfig struct {
    AccessTTL  time.Duration // Default: 1 hour
    RefreshTTL time.Duration // Default: 24 hours
}
```

## Testing

### Unit Tests

```bash
# Run auth tests
go test ./internal/modules/auth/...

# Run with coverage
go test -cover ./internal/modules/auth/...
```

### Integration Tests

```bash
# Test auth endpoints
curl -X POST http://localhost:8080/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User","password":"password123"}'
```

## Monitoring & Logging

### Log Levels
- **DEBUG**: Detailed authentication flow
- **INFO**: Successful login/logout events
- **WARN**: Failed authentication attempts
- **ERROR**: System errors and exceptions

### Metrics
- Login success/failure rates
- Token refresh frequency
- Session duration statistics
- Error rate monitoring

## Troubleshooting

### Common Issues

1. **Token Expired**: Refresh token atau login ulang
2. **Invalid Credentials**: Cek email dan password
3. **Database Connection**: Cek koneksi database
4. **Google OAuth**: Validasi client ID dan secret

### Debug Steps

1. Check server logs for detailed error messages
2. Verify database connection and schema
3. Test JWT token validation
4. Monitor authentication flow with debug logs 