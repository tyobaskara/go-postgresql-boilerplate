package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/tyobaskara/jeki-backend/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateSession(session *domain.Session) error {
	return r.db.Create(session).Error
}

func (r *authRepository) GetSessionByRefreshToken(refreshToken string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.Where("refresh_token = ? AND expires_at > ?", refreshToken, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *authRepository) DeleteSession(id uuid.UUID) error {
	return r.db.Delete(&domain.Session{}, "id = ?", id).Error
}

func (r *authRepository) DeleteUserSessions(userID uuid.UUID) error {
	return r.db.Delete(&domain.Session{}, "user_id = ?", userID).Error
} 