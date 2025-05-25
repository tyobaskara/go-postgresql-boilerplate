// Package usecase implements the business logic layer of the application
package usecase

import (
	"context"

	"github.com/tyobaskara/jeki-backend/internal/domain/entity"
	"github.com/tyobaskara/jeki-backend/internal/domain/repository"
)

// UserUsecase implements the business logic for user operations
// It acts as an intermediary between the handlers and the repository
type UserUsecase struct {
	userRepo repository.UserRepository // Repository for user data operations
}

// NewUserUsecase creates a new instance of UserUsecase
// Parameters:
//   - userRepo: implementation of UserRepository interface
// Returns:
//   - *UserUsecase: new instance of UserUsecase
func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// CreateUser handles the business logic for creating a new user
// Parameters:
//   - ctx: context for the operation
//   - user: pointer to the user entity to be created
// Returns:
//   - error: any error that occurred during creation
func (uc *UserUsecase) CreateUser(ctx context.Context, user *entity.User) error {
	return uc.userRepo.Create(ctx, user)
}

// GetUserByID handles the business logic for retrieving a user by ID
// Parameters:
//   - ctx: context for the operation
//   - id: the user's ID to find
// Returns:
//   - *entity.User: pointer to the found user
//   - error: any error that occurred during retrieval
func (uc *UserUsecase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// GetUserByEmail handles the business logic for retrieving a user by email
// Parameters:
//   - ctx: context for the operation
//   - email: the user's email to find
// Returns:
//   - *entity.User: pointer to the found user
//   - error: any error that occurred during retrieval
func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return uc.userRepo.FindByEmail(ctx, email)
}

// UpdateUser handles the business logic for updating a user
// Parameters:
//   - ctx: context for the operation
//   - user: pointer to the user entity to be updated
// Returns:
//   - error: any error that occurred during update
func (uc *UserUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	return uc.userRepo.Update(ctx, user)
}

// DeleteUser handles the business logic for deleting a user
// Parameters:
//   - ctx: context for the operation
//   - id: the ID of the user to delete
// Returns:
//   - error: any error that occurred during deletion
func (uc *UserUsecase) DeleteUser(ctx context.Context, id uint) error {
	return uc.userRepo.Delete(ctx, id)
} 