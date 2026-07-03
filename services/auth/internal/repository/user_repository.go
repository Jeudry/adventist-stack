// Package repository implementa el acceso a datos del servicio auth sobre
// PostgreSQL (pgx).
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeudry/adventist-stack/services/auth/internal/domain"
	"github.com/Jeudry/adventist-stack/services/auth/internal/models"
)

// UserRepository persiste y consulta usuarios.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository crea el repositorio.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserta un usuario y devuelve el modelo de dominio resultante.
func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	const q = `
		INSERT INTO users (email, name, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, password_hash, role, created_at, updated_at`

	var m models.User
	err := r.db.QueryRow(ctx, q, u.Email, u.Name, u.PasswordHash, string(u.Role)).
		Scan(&m.ID, &m.Email, &m.Name, &m.PasswordHash, &m.Role, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("repository: create user: %w", err)
	}
	return m.ToDomain(), nil
}

// FindByEmail busca un usuario por email. Devuelve ErrUserNotFound si no existe.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	const q = `
		SELECT id, email, name, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1`

	var m models.User
	err := r.db.QueryRow(ctx, q, email).
		Scan(&m.ID, &m.Email, &m.Name, &m.PasswordHash, &m.Role, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("repository: find by email: %w", err)
	}
	return m.ToDomain(), nil
}

// ExistsByEmail indica si ya hay un usuario con ese email.
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	if err := r.db.QueryRow(ctx, q, email).Scan(&exists); err != nil {
		return false, fmt.Errorf("repository: exists by email: %w", err)
	}
	return exists, nil
}
