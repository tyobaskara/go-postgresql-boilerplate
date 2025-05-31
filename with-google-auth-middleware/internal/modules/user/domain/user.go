package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
// This is the core domain entity that contains user information
type User struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier for the user
	Email     string    `json:"email"`      // User's email address (used for login)
	Name      string    `json:"name"`       // User's full name
	CreatedAt time.Time `json:"created_at"` // Timestamp when the user was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the user was last updated
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	GetAll(page, limit int) ([]*User, error)
}

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	CreateUser(user *User) error
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uuid.UUID) error
	GetAllUsers(page, limit int) ([]*User, error)
} 