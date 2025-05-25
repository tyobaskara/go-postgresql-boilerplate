// Package entity defines the core domain entities of the application
package entity

import "time"

// User represents a user in the system
// This is the core domain entity that contains user information
type User struct {
	ID        uint      `json:"id"`         // Unique identifier for the user
	Email     string    `json:"email"`      // User's email address (used for login)
	Name      string    `json:"name"`       // User's full name
	CreatedAt time.Time `json:"created_at"` // Timestamp when the user was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the user was last updated
} 