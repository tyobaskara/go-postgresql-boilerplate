# Data Structure Documentation

## Overview

This document explains the data structures used in the authentication module and how they map to the database tables using GORM (Go Object Relational Mapper).

## Database Tables

### Users Table

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255),
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_logout_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**Note**: The `password` field is nullable to support users who authenticate via Google OAuth (who don't have passwords).

### Sessions Table

```sql
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

## Go Structs

### User Struct

```go
type User struct {
    ID           uuid.UUID `json:"id"`
    Email        string    `json:"email"`
    Name         string    `json:"name"`
    Password     string    `json:"-"`          // Hashed password, not exposed in JSON
    LastLoginAt  time.Time `json:"lastLoginAt"`
    LastLogoutAt time.Time `json:"lastLogoutAt"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
```

### Session Struct

```go
type Session struct {
    ID           uuid.UUID `json:"id"`
    UserID       uuid.UUID `json:"userId"`
    RefreshToken string    `json:"refreshToken"`
    ExpiresAt    time.Time `json:"expiresAt"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
```

### Request/Response Structs

```go
// SignupRequest represents the request payload for user signup
type SignupRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Name     string `json:"name" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// AuthToken represents the JWT token structure
type AuthToken struct {
    AccessToken  string    `json:"accessToken"`
    TokenType    string    `json:"tokenType"`
    ExpiresIn    int64     `json:"expiresIn"`
    RefreshToken string    `json:"refreshToken,omitempty"`
    ExpiresAt    time.Time `json:"expiresAt"`
}
```

## GORM Mapping

GORM (Go Object Relational Mapper) automatically handles the mapping between Go structs and database tables based on naming conventions:

### Table Naming Convention

- Go struct name is converted to plural form for table name
  - `User` struct → `users` table
  - `Session` struct → `sessions` table

### Column Naming Convention

- Go struct field names are converted to snake_case for column names
  - `ID` → `id`
  - `UserID` → `user_id`
  - `RefreshToken` → `refresh_token`
  - `CreatedAt` → `created_at`
  - `UpdatedAt` → `updated_at`

### Type Mapping

GORM automatically maps Go types to SQL types:
- `uuid.UUID` → `UUID`
- `string` → `VARCHAR(255)` or `TEXT`
- `time.Time` → `TIMESTAMP WITH TIME ZONE`
- `bool` → `BOOLEAN`
- `int` → `INTEGER`
- `float64` → `DOUBLE PRECISION`

### Special Fields

GORM recognizes special field names and handles them automatically:
- `ID` or `Id`: Primary key
- `CreatedAt`: Automatically set on creation
- `UpdatedAt`: Automatically updated on modification
- `DeletedAt`: For soft deletes (if using `gorm.Model`)

### JSON Tags

- `json:"-"`: Field is excluded from JSON serialization (used for password)
- `json:"fieldName"`: Custom JSON field name in camelCase
- `binding:"required,email"`: Gin validation tags

## Repository Implementation

The repository layer uses GORM to perform database operations:

```go
type userRepository struct {
    db *gorm.DB
}

// Create a new user
func (r *userRepository) Create(user *domain.User) error {
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    return r.db.Create(user).Error
}

// Find user by email
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Find user by ID
func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Update user
func (r *userRepository) Update(user *domain.User) error {
    user.UpdatedAt = time.Now()
    return r.db.Save(user).Error
}

// Delete user
func (r *userRepository) Delete(id uuid.UUID) error {
    return r.db.Delete(&domain.User{}, id).Error
}

// Get all users with pagination
func (r *userRepository) GetAll(page, limit int) ([]*domain.User, error) {
    var users []*domain.User
    offset := (page - 1) * limit
    err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error
    if err != nil {
        return nil, err
    }
    return users, nil
}
```

```go
type authRepository struct {
    db *gorm.DB
}

// Create a new session
func (r *authRepository) CreateSession(session *domain.Session) error {
    return r.db.Create(session).Error
}

// Get session by refresh token
func (r *authRepository) GetSessionByRefreshToken(refreshToken string) (*domain.Session, error) {
    var session domain.Session
    err := r.db.Where("refresh_token = ? AND expires_at > ?", refreshToken, time.Now()).First(&session).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

// Delete a session
func (r *authRepository) DeleteSession(id uuid.UUID) error {
    return r.db.Delete(&domain.Session{}, "id = ?", id).Error
}

// Delete all sessions for a user
func (r *authRepository) DeleteUserSessions(userID uuid.UUID) error {
    return r.db.Delete(&domain.Session{}, "user_id = ?", userID).Error
}
```

## Password Security

### Password Hashing

Passwords are hashed using bcrypt with the following configuration:
- **Algorithm**: bcrypt
- **Cost Factor**: 10 (default)
- **Salt**: Automatically generated by bcrypt

```go
// Hash password during signup
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return fmt.Errorf("failed to hash password: %w", err)
}

// Verify password during login
err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
if err != nil {
    return nil, errors.New("invalid password")
}
```

## Best Practices

1. **Struct Tags**
   - Use `json` tags for API serialization in camelCase
   - Use `json:"-"` to exclude sensitive fields from JSON
   - Use `binding` tags for input validation
   - Use `gorm` tags for custom GORM behavior if needed

2. **Indexes**
   - Create indexes for frequently queried columns
   - Create indexes for foreign key columns
   - Create indexes for columns used in WHERE clauses
   - Create unique indexes for email addresses

3. **Relationships**
   - Use foreign key constraints for data integrity
   - Use `ON DELETE CASCADE` for automatic cleanup
   - Define relationships in both directions

4. **Timestamps**
   - Always include `created_at` and `updated_at`
   - Let GORM handle timestamp updates
   - Use timezone-aware timestamps

5. **Error Handling**
   - Always check for errors from GORM operations
   - Use transactions for multiple operations
   - Handle unique constraint violations
   - Handle record not found scenarios

6. **Security**
   - Never store plain text passwords
   - Use bcrypt for password hashing
   - Exclude sensitive fields from JSON responses
   - Validate input data
   - Use prepared statements (handled by GORM)

7. **Data Types**
   - Use UUID for primary keys
   - Use appropriate string lengths
   - Use timezone-aware timestamps
   - Consider nullable fields for optional data 