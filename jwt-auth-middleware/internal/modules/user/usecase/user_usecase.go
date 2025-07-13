package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tyobaskara/maxwash-backend/internal/modules/user/domain"
	"golang.org/x/crypto/bcrypt"
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

func (u *userUsecase) UpdateLastLoginAt(id uuid.UUID, loginTime time.Time) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	user.LastLoginAt = loginTime
	user.UpdatedAt = time.Now()
	return u.userRepo.Update(user)
}

func (u *userUsecase) UpdateLastLogoutAt(id uuid.UUID, logoutTime time.Time) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	user.LastLogoutAt = logoutTime
	user.UpdatedAt = time.Now()
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

// SignupUser handles user registration with password hashing
func (u *userUsecase) SignupUser(request *domain.SignupRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &domain.User{
		ID:        uuid.New(),
		Email:     request.Email,
		Name:      request.Name,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ValidatePassword validates user credentials and returns user if valid
func (u *userUsecase) ValidatePassword(email, password string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if user has password (for Google OAuth users who might not have password)
	if user.Password == "" {
		return nil, errors.New("user does not have password authentication enabled")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}