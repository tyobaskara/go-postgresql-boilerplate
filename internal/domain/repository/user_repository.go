// Package repository defines the interfaces for data access layer
package repository

import (
	"context"

	"github.com/tyobaskara/jeki-backend/internal/domain/entity"
)

// UserRepository defines the interface for user data operations
// This interface follows the Repository pattern for data access abstraction
type UserRepository interface {
	// Create persists a new user to the database
	// Parameters:
	//   - ctx: context for the operation
	//   - user: pointer to the user entity to be created
	// Returns:
	//   - error: any error that occurred during creation
	Create(ctx context.Context, user *entity.User) error

	// FindByID retrieves a user by their ID
	// Parameters:
	//   - ctx: context for the operation
	//   - id: the user's ID to find
	// Returns:
	//   - *entity.User: pointer to the found user
	//   - error: any error that occurred during retrieval
	FindByID(ctx context.Context, id uint) (*entity.User, error)

	// FindByEmail retrieves a user by their email address
	// Parameters:
	//   - ctx: context for the operation
	//   - email: the user's email to find
	// Returns:
	//   - *entity.User: pointer to the found user
	//   - error: any error that occurred during retrieval
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update modifies an existing user in the database
	// Parameters:
	//   - ctx: context for the operation
	//   - user: pointer to the user entity to be updated
	// Returns:
	//   - error: any error that occurred during update
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user from the database
	// Parameters:
	//   - ctx: context for the operation
	//   - id: the ID of the user to delete
	// Returns:
	//   - error: any error that occurred during deletion
	Delete(ctx context.Context, id uint) error
} 