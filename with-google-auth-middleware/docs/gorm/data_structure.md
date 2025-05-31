# Data Structure Documentation

## Overview

This document explains the data structures used in the authentication module and how they map to the database tables using GORM (Go Object Relational Mapper).

## Database Tables

### Users Table

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

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
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Session Struct

```go
type Session struct {
    ID           uuid.UUID `json:"id"`
    UserID       uuid.UUID `json:"user_id"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
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
- `string` → `TEXT`
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

## Repository Implementation

The repository layer uses GORM to perform database operations:

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

## Best Practices

1. **Struct Tags**
   - Use `json` tags for API serialization
   - Use `gorm` tags for custom GORM behavior if needed
   - Use `validate` tags for input validation

2. **Indexes**
   - Create indexes for frequently queried columns
   - Create indexes for foreign key columns
   - Create indexes for columns used in WHERE clauses

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