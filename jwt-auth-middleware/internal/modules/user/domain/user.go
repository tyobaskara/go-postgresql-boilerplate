package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
// This is the core domain entity that contains user information
type User struct {
	ID           uuid.UUID `json:"id"`            // Unique identifier for the user
	Email        string    `json:"email"`         // User's email address (used for login)
	Name         string    `json:"name"`          // User's full name
	Password     string    `json:"-"`             // User's password (hashed, not exposed in JSON)
	LastLoginAt  time.Time `json:"lastLoginAt"`   // Timestamp when the user last logged in
	LastLogoutAt time.Time `json:"lastLogoutAt"`  // Timestamp when the user last logged out
	CreatedAt    time.Time `json:"createdAt"`     // Timestamp when the user was created
	UpdatedAt    time.Time `json:"updatedAt"`     // Timestamp when the user was last updated
}

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

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	UpdateLastLoginAt(id uuid.UUID, loginTime time.Time) error
	UpdateLastLogoutAt(id uuid.UUID, logoutTime time.Time) error
	Delete(id uuid.UUID) error
	GetAll(page, limit int) ([]*User, error)
}

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	CreateUser(user *User) error
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	UpdateLastLoginAt(id uuid.UUID, loginTime time.Time) error
	UpdateLastLogoutAt(id uuid.UUID, logoutTime time.Time) error
	DeleteUser(id uuid.UUID) error
	GetAllUsers(page, limit int) ([]*User, error)
} 