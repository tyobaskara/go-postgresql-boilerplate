package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/tyobaskara/maxwash-backend/internal/modules/auth/domain"
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

func (r *authRepository) GetUserSessions(userID uuid.UUID) ([]*domain.Session, error) {
	var sessions []*domain.Session
	err := r.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *authRepository) UpdateSession(session *domain.Session) error {
	return r.db.Save(session).Error
}

func (r *authRepository) DeleteSession(id uuid.UUID) error {
	return r.db.Delete(&domain.Session{}, "id = ?", id).Error
}

func (r *authRepository) DeleteUserSessions(userID uuid.UUID) error {
	return r.db.Delete(&domain.Session{}, "user_id = ?", userID).Error
} 