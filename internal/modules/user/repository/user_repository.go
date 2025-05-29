package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/tyobaskara/jeki-backend/internal/modules/user/domain"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Name,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query,
		user.Email,
		user.Name,
		time.Now(),
		user.ID,
	)
	return err
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	return err
} 