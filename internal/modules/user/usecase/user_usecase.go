package usecase

import (
	"github.com/google/uuid"
	"github.com/tyobaskara/jeki-backend/internal/modules/user/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	return u.userRepo.Create(user)
}

func (u *userUsecase) GetUserByID(id uuid.UUID) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) GetUserByEmail(email string) (*domain.User, error) {
	return u.userRepo.FindByEmail(email)
}

func (u *userUsecase) UpdateUser(user *domain.User) error {
	return u.userRepo.Update(user)
}

func (u *userUsecase) DeleteUser(id uuid.UUID) error {
	return u.userRepo.Delete(id)
}

func (u *userUsecase) GetAllUsers(page, limit int) ([]*domain.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return u.userRepo.GetAll(page, limit)
}